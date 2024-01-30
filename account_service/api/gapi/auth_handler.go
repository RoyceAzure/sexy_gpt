package gapi

import (
	"context"
	"time"

	"github.com/RoyceAzure/sexy_gpt/account_service/api/gapi/token"
	db "github.com/RoyceAzure/sexy_gpt/account_service/repository/db/sqlc"
	logger "github.com/RoyceAzure/sexy_gpt/account_service/repository/logger_distributor"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/pb"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util/gpt_error"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util/validate"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) Login(ctx context.Context, req *pb.LoginRequset) (*pb.LoginResponse, error) {
	if violations := valideteLoginReq(req); violations != nil {
		return nil, gpt_error.InvalidArgumentError(violations)
	}

	user, err := s.dao.GetUserDTOByEmail(ctx, req.GetEmail())
	if err != nil {
		if err.Error() == gpt_error.DB_ERR_NOT_FOUND.Error() {
			return nil, s.HandleAPIError(codes.NotFound, err, "user not exists")
		}
		logger.Logger.Error().Err(err).Msg("failed to get user")

		return nil, s.HandleAPIError(codes.Internal, err, "something wrong")
	}

	if !user.IsEmailVerified {
		return nil, s.HandleAPIError(codes.PermissionDenied, err, "user email is not vertified")

	}

	if err := util.CheckPassword(req.GetPassword(), user.HashedPassword); err != nil {
		return nil, s.HandleAPIError(codes.PermissionDenied, err, "wrong password")
	}

	userId := user.UserID.Bytes
	roldeId := user.RoleID.Bytes
	tokenSubject := token.NewTokenSubject(user.Email, userId, roldeId)
	oldSession, err := s.dao.GetSessionByUserId(ctx, user.UserID)

	var refreshToken string
	var refreshPayLoad *token.TokenPayload

	if err != nil {
		if err.Error() == gpt_error.DB_ERR_NOT_FOUND.Error() {
			refreshToken, refreshPayLoad, err = s.tokenMaker.CreateToken(nil, "refresh", s.config.AUTH_ISSUER, s.config.RefreshTokenDuration)
			if err != nil {

				return nil, s.HandleAPIError(codes.PermissionDenied, err, "failed to create token")
			}

			_, err = s.dao.CreateSession(ctx, db.CreateSessionParams{
				UserID: pgtype.UUID{
					Bytes: userId,
					Valid: true,
				},
				RefreshToken: refreshToken,
				UserAgent:    "todo",
				ClientIp:     "todo",
				ExpiredAt: pgtype.Timestamptz{
					Time:  refreshPayLoad.ExpiredAt,
					Valid: true,
				},
			})
			if err != nil {
				return nil, s.HandleAPIError(codes.PermissionDenied, err, "failed to create session")
			}
		} else {
			_, err := s.dao.DeleteSession(ctx, oldSession.ID)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "%s", err)
			}
		}
	} else {
		if time.Now().After(oldSession.ExpiredAt.Time) {
			_, err := s.dao.DeleteSession(ctx, oldSession.ID)
			if err != nil {
				return nil, s.HandleAPIError(codes.Internal, err, "internal err")
			}
			refreshToken, refreshPayLoad, err := s.tokenMaker.CreateToken(nil, "refresh", s.config.AUTH_ISSUER, s.config.RefreshTokenDuration)
			if err != nil {
				return nil, s.HandleAPIError(codes.Internal, err, "internal err")
			}

			_, err = s.dao.CreateSession(ctx, db.CreateSessionParams{
				UserID: pgtype.UUID{
					Bytes: userId,
					Valid: true,
				},
				RefreshToken: refreshToken,
				UserAgent:    "todo",
				ClientIp:     "todo",
				ExpiredAt: pgtype.Timestamptz{
					Time:  refreshPayLoad.ExpiredAt,
					Valid: true,
				},
			})
			if err != nil {
				return nil, s.HandleAPIError(codes.PermissionDenied, err, "internal err")
			}
		}
		refreshToken = oldSession.RefreshToken
	}

	accessToken, accessPayLoad, err := s.tokenMaker.CreateToken(tokenSubject, s.config.AUTH_AUDIENCE, s.config.AUTH_ISSUER, s.config.AccessTokenDuration)
	if err != nil {
		return nil, s.HandleAPIError(codes.PermissionDenied, err, "internal err")
	}
	res := &pb.LoginResponse{
		Data: convertUserDTO(user),
		TokenData: &pb.Token{
			Audience:     accessPayLoad.Audience,
			Issuer:       accessPayLoad.Issuer,
			IssureAt:     timestamppb.New(accessPayLoad.IssuedAt),
			ExpiredAt:    timestamppb.New(accessPayLoad.ExpiredAt),
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}
	return res, nil
}

/*
Logout 不管執行正確性

*/

func (s *Server) Logout(ctx context.Context, req *pb.LogoutRequset) (*pb.LogoutResponse, error) {
	payload, err := s.authorizUser(ctx)
	if err != nil {
		return nil, s.HandleAPIError(codes.Unauthenticated, err, "failed to logout")
	}

	userId := payload.UserId

	session, err := s.dao.GetSessionByUserId(ctx, pgtype.UUID{
		Bytes: userId,
		Valid: true,
	})
	if err != nil {
		if err.Error() == gpt_error.DB_ERR_NOT_FOUND.Error() {
			return nil, s.HandleAPIError(codes.NotFound, err, "already logout")
		}
		return nil, s.HandleAPIError(codes.Internal, err, "failed to logout")
	} else {
		_, err = s.dao.DeleteSession(ctx, pgtype.UUID{
			Bytes: session.ID.Bytes,
			Valid: true,
		})
		if err != nil {
			return nil, s.HandleAPIError(codes.Internal, err, "failed to logout")
		}
	}
	return &pb.LogoutResponse{}, nil
}
func (s *Server) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequset) (*pb.RefreshTokenResponse, error) {
	payload, err := s.authorizUser(ctx)
	if err != nil {
		return nil, s.HandleAPIError(codes.Unauthenticated, err, "failed to refresh token")
	}

	session, err := s.dao.GetSessionByUserId(ctx, pgtype.UUID{
		Bytes: payload.UserId,
		Valid: true,
	})

	if err != nil || time.Now().After(session.ExpiredAt.Time) || session.IsBlocked {
		s.dao.DeleteSession(ctx, session.ID)
		return &pb.RefreshTokenResponse{}, status.Errorf(codes.Canceled, "please login again")
	}

	if session.RefreshToken != req.GetRefreshToken() {
		return &pb.RefreshTokenResponse{}, status.Errorf(codes.Canceled, "invalid token")
	}

	tokenSubject := token.NewTokenSubject(payload.Email, payload.UserId, payload.RoleId)
	accessToken, accessPayLoad, err := s.tokenMaker.CreateToken(tokenSubject, s.config.AUTH_AUDIENCE, s.config.AUTH_ISSUER, time.Hour*1)
	if err != nil {
		return nil, s.HandleAPIError(codes.Internal, err, "failed to create token")
	}

	res := &pb.RefreshTokenResponse{
		Audience:    accessPayLoad.Audience,
		Issuer:      accessPayLoad.Issuer,
		IssureAt:    timestamppb.New(accessPayLoad.IssuedAt),
		ExpiredAt:   timestamppb.New(accessPayLoad.ExpiredAt),
		AccessToken: accessToken,
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
