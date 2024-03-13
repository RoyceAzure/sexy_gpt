package gapi

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/RoyceAzure/sexy_gpt/account_service/api/gapi/token"
	db "github.com/RoyceAzure/sexy_gpt/account_service/repository/db/sqlc"
	sso "github.com/RoyceAzure/sexy_gpt/account_service/repository/sso"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/pb"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util/converter"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util/gpt_error"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util/random"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util/validate"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

/*
des:

	ctx加入db message
	將code, msg轉換成status.Errorf
	msg 包入AuthDTOResponse
*/
func processAuthResponse(ctx context.Context, code codes.Code, msg string, err error) (*pb.AuthDTOResponse, error) {
	res := pb.AuthDTOResponse{Message: msg}
	if code != codes.OK {
		if err == nil {
			err = fmt.Errorf(msg)
		}
		util.NewOutGoingMetaDataKV(ctx, util.DBMSGKey, err.Error())
		return &res, status.Errorf(code, msg)
	}
	return &res, nil
}

func (s *Server) Login(ctx context.Context, req *pb.LoginRequset) (*pb.AuthDTOResponse, error) {
	if violations := valideteLoginReqV2(req); violations != nil {
		msg, err := violations.ToJson()
		if err != nil {
			msg = "Invalidate argument"
		}
		return processAuthResponse(ctx, codes.InvalidArgument, msg, err)
	}

	user, err := s.Service.IsValidateUser(ctx, req.GetEmail())
	if err != nil {
		if errors.Is(err, gpt_error.ErrNotFound) {
			return processAuthResponse(ctx, codes.NotFound, "wrong email or password", err)
		} else if errors.Is(err, gpt_error.ErrUnauthicated) {
			return processAuthResponse(ctx, codes.PermissionDenied, "user email is not vertified", err)
		} else {
			return processAuthResponse(ctx, codes.Internal, "internal err", err)
		}
	}

	if err := util.CheckPassword(req.GetPassword(), user.HashedPassword); err != nil {
		return processAuthResponse(ctx, codes.PermissionDenied, "wrong email or password", err)
	}

	userId := user.UserID.Bytes
	roldeId := user.RoleID.Bytes
	tokenSubject := token.NewTokenSubject(user.Email, userId, roldeId)

	var refreshToken string
	refreshToken, err = s.Service.LoginCreateSession(ctx, userId, s.tokenMaker, s.config)
	if err != nil {
		processAuthResponse(ctx, codes.Internal, "internal err", err)
	}

	accessToken, accessPayLoad, err := s.tokenMaker.CreateToken(tokenSubject, s.config.AUTH_AUDIENCE, s.config.AUTH_ISSUER, s.config.AccessTokenDuration)
	if err != nil {
		return processAuthResponse(ctx, codes.Internal, "internal err", err)
	}

	res := ConvertAuthDTO("login successed", ConvertUserDTO(user), ConvertToken(accessPayLoad, accessToken, refreshToken))
	if err != nil {
		return processAuthResponse(ctx, codes.Internal, "internal err", err)
	}
	return res, nil
}

/*
Logout 不管執行正確性

*/

func (s *Server) Logout(ctx context.Context, req *pb.LogoutRequset) (*pb.AuthDTOResponse, error) {
	payload, err := s.authorizUser(ctx)
	if err != nil {
		return processAuthResponse(ctx, codes.Unauthenticated, err.Error(), err)
	}

	userId := payload.UserId

	session, err := s.dao.GetSessionByUserId(ctx, pgtype.UUID{
		Bytes: userId,
		Valid: true,
	})
	if err != nil {
		if err.Error() == gpt_error.DB_ERR_NOT_FOUND.Error() {
			return processAuthResponse(ctx, codes.NotFound, "already logout", err)
		}
		return processAuthResponse(ctx, codes.NotFound, "internal err", err)
	} else {
		_, err = s.dao.DeleteSession(ctx, pgtype.UUID{
			Bytes: session.ID.Bytes,
			Valid: true,
		})
		if err != nil {
			return processAuthResponse(ctx, codes.NotFound, "internal err", err)
		}
	}

	res := ConvertAuthDTO("logout successed", nil, nil)
	if err != nil {
		return processAuthResponse(ctx, codes.Internal, "internal err", err)
	}
	return res, nil
}

func (s *Server) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequset) (*pb.AuthDTOResponse, error) {
	payload, err := s.authorizUser(ctx)
	if err != nil {
		return processAuthResponse(ctx, codes.Unauthenticated, err.Error(), err)
	}

	session, err := s.dao.GetSessionByUserId(ctx, pgtype.UUID{
		Bytes: payload.UserId,
		Valid: true,
	})

	if err != nil || time.Now().After(session.ExpiredAt.Time) || session.IsBlocked {
		s.dao.DeleteSession(ctx, session.ID)
		return processAuthResponse(ctx, codes.Unauthenticated, "please login again", err)
	}

	if session.RefreshToken != req.GetRefreshToken() {
		return processAuthResponse(ctx, codes.Unauthenticated, "invalid refresh token", err)
	}

	tokenSubject := token.NewTokenSubject(payload.Email, payload.UserId, payload.RoleId)
	accessToken, accessPayLoad, err := s.tokenMaker.CreateToken(tokenSubject, s.config.AUTH_AUDIENCE, s.config.AUTH_ISSUER, s.config.AccessTokenDuration)
	if err != nil {
		return processAuthResponse(ctx, codes.Internal, "internal err", err)
	}

	res := ConvertAuthDTO("refresh token successed", nil, ConvertToken(accessPayLoad, accessToken, ""))
	if err != nil {
		return processAuthResponse(ctx, codes.Internal, "internal err", err)
	}
	return res, nil
}

func (s *Server) SSOGoogleLogin(ctx context.Context, req *pb.GoogleIDTokenRequest) (*pb.AuthDTOResponse, error) {
	tokenInfo, err := sso.VertifyGoogleSSOIDToken(ctx, req.GetToken(), s.config.FRONTEND_KEY)
	if err != nil {
		return processAuthResponse(ctx, codes.Unauthenticated, "invalid token", err)
	}

	userEmail := tokenInfo.Email
	user, err := s.dao.GetUserDTOByEmail(ctx, userEmail)
	if err != nil {
		//case user not exists : create user
		if err.Error() == gpt_error.DB_ERR_NOT_FOUND.Error() {

			userName := strings.Split(userEmail, "@")[0]

			randomPas := random.RandomString(15)
			hashPassword, err := util.HashPassword(randomPas)
			if err != nil {
				return processAuthResponse(ctx, codes.Internal, "internal err", err)
			}

			arg := db.CreateUserTxParms{
				Arg: &db.CreateUserParams{
					UserName:       userName,
					Email:          userEmail,
					HashedPassword: hashPassword,
					IsInternal:     false,
					CrUser:         "SYSTEM",
					SsoIdentifer: pgtype.Text{
						String: "google",
						Valid:  true,
					},
				},
			}

			TxResult, err := s.dao.CreateUserTx(ctx, &arg)
			if err != nil {
				return processAuthResponse(ctx, codes.Internal, "internal err", err)
			}

			user = db.UserRoleView{
				UserID:          TxResult.User.UserID,
				UserName:        TxResult.User.UserName,
				IsEmailVerified: TxResult.User.IsEmailVerified,
				Email:           TxResult.User.Email,
			}
		} else {
			return processAuthResponse(ctx, codes.Internal, "internal err", err)
		}
	}

	// set IsEmailVerified true
	if !user.IsEmailVerified {
		_, err = s.dao.UpdateUser(ctx, db.UpdateUserParams{
			UserID: user.UserID,
			IsEmailVerified: pgtype.Bool{
				Bool:  true,
				Valid: true,
			},
		})
		if err != nil {
			return processAuthResponse(ctx, codes.Internal, "internal err", err)
		}
	}

	userId := user.UserID.Bytes
	roldeId := user.RoleID.Bytes
	tokenSubject := token.NewTokenSubject(user.Email, userId, roldeId)

	//create login session
	var refreshToken string
	refreshToken, err = s.Service.LoginCreateSession(ctx, userId, s.tokenMaker, s.config)
	if err != nil {
		processAuthResponse(ctx, codes.Internal, "internal err", err)
	}

	//create token
	accessToken, accessPayLoad, err := s.tokenMaker.CreateToken(tokenSubject, s.config.AUTH_AUDIENCE, s.config.AUTH_ISSUER, s.config.AccessTokenDuration)
	if err != nil {
		return processAuthResponse(ctx, codes.Internal, "internal err", err)
	}

	res := ConvertAuthDTO("login successed", ConvertUserDTO(&user), ConvertToken(accessPayLoad, accessToken, refreshToken))
	return res, nil
}

func (s *Server) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.AuthDTOResponse, error) {
	payload, err := s.tokenMaker.VertifyToken(req.GetAccessToken())
	if err != nil {
		return processAuthResponse(ctx, codes.Unauthenticated, err.Error(), err)
	}

	user, err := s.Service.IsValidateUser(ctx, payload.Subject.UPN)
	if err != nil {
		return processAuthResponse(ctx, codes.Unauthenticated, err.Error(), err)
	}

	_, err = s.Service.IsUserLogin(ctx, payload.Subject.UserId)
	if err != nil {
		return processAuthResponse(ctx, codes.Unauthenticated, err.Error(), err)
	}

	userDTO := pb.UserDTO{
		UserId:   converter.ConvertXByte2UUID(user.UserID.Bytes),
		UserName: user.UserName,
		RoleName: user.RoleName.String,
		Email:    user.Email,
		RoleId:   converter.ConvertXByte2UUID(user.RoleID.Bytes),
	}
	res := ConvertAuthDTO("validatd access token successed", &userDTO, nil)

	return res, nil
}

func valideteLoginReqV2(req *pb.LoginRequset) *gpt_error.Message {
	violations := &gpt_error.Message{}
	if err := validate.ValidateEmptyString(req.GetEmail()); err != nil {
		gpt_error.AddErrFieldString(violations, gpt_error.NewErrField("email", err))
	}
	if err := validate.ValidateEmptyString(req.GetPassword()); err != nil {
		gpt_error.AddErrFieldString(violations, gpt_error.NewErrField("password", err))
	}
	if len(violations.ErrMessage) == 0 {
		return nil
	}
	return violations
}

func ConvertAuthDTO(message string, userDTO *pb.UserDTO, tokenDTO *pb.Token) *pb.AuthDTOResponse {
	return &pb.AuthDTOResponse{
		Message: message,
		User:    userDTO,
		Token:   tokenDTO,
	}
}

func ConvertToken(token *token.TokenPayload, accessToken string, refreshToken string) *pb.Token {
	return &pb.Token{
		Audience:     token.Audience,
		Issuer:       token.Issuer,
		IssureAt:     timestamppb.New(token.IssuedAt),
		ExpiredAt:    timestamppb.New(token.ExpiredAt),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
