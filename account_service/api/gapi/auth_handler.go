package gapi

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/RoyceAzure/sexy_gpt/account_service/api/gapi/token"
	sso "github.com/RoyceAzure/sexy_gpt/account_service/repository/SSO"
	db "github.com/RoyceAzure/sexy_gpt/account_service/repository/db/sqlc"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/pb"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util/gpt_error"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util/random"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util/validate"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
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
	if violations := valideteLoginReq(req); violations != nil {
		return nil, gpt_error.InvalidArgumentError(violations)
	}

	user, err := s.Service.IsValidateUser(ctx, req.GetEmail())
	if err != nil {
		if errors.Is(err, gpt_error.ErrNotFound) {
			return processAuthResponse(ctx, codes.NotFound, "user not exists", err)
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
		return processAuthResponse(ctx, codes.Unauthenticated, "invalid token", err)
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
	// httpClient, err := google.DefaultClient(ctx)
	// if err != nil {
	// 	return processAuthResponse(ctx, codes.Internal, "internal err", err)
	// }

	tokenInfo, err := sso.VertifyGoogleSSO(ctx, req.GetIdToken())
	if err != nil {
		return processAuthResponse(ctx, codes.Unauthenticated, "invalid token", err)
	}

	userEmail := tokenInfo.Email
	user, err := s.dao.GetUserDTOByEmail(ctx, userEmail)
	if err != nil {
		if err.Error() == gpt_error.DB_ERR_NOT_FOUND.Error() {
			// return processAuthResponse(ctx, codes.NotFound, "user not exists", err)
			//創建使用者
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
		}
		return processAuthResponse(ctx, codes.Internal, "internal err", err)
	}

	//直接修改為已經驗證
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

	var refreshToken string
	refreshToken, err = s.Service.LoginCreateSession(ctx, userId, s.tokenMaker, s.config)
	if err != nil {
		processAuthResponse(ctx, codes.Internal, "internal err", err)
	}

	accessToken, accessPayLoad, err := s.tokenMaker.CreateToken(tokenSubject, s.config.AUTH_AUDIENCE, s.config.AUTH_ISSUER, s.config.AccessTokenDuration)
	if err != nil {
		return processAuthResponse(ctx, codes.Internal, "internal err", err)
	}

	res := ConvertAuthDTO("login successed", ConvertUserDTO(&user), ConvertToken(accessPayLoad, accessToken, refreshToken))
	if err != nil {
		return processAuthResponse(ctx, codes.Internal, "internal err", err)
	}
	return res, nil
}

func valideteLoginReq(req *pb.LoginRequset) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validate.ValidateEmptyString(req.GetEmail()); err != nil {
		violations = append(violations, validate.FieldViolation("email", err))
	}
	if err := validate.ValidateEmptyString(req.GetPassword()); err != nil {
		violations = append(violations, validate.FieldViolation("password", err))
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
