package gapi

import (
	"context"
	"fmt"
	"time"

	db "github.com/RoyceAzure/sexy_gpt/account_service/repository/db/sqlc"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/pb"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util/gpt_error"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util/validate"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func processVertifyEmailResponse(ctx context.Context, code codes.Code, msg string, err error) (*pb.VertifyEmailResponse, error) {
	res := pb.VertifyEmailResponse{Message: msg}
	if code != codes.OK {
		if err == nil {
			err = fmt.Errorf(msg)
		}
		util.NewOutGoingMetaDataKV(ctx, util.DBMSGKey, err.Error())
		return &res, status.Errorf(code, msg)
	}
	return &res, nil
}

/*
安全考量  就算沒有使用者也會回傳成功
*/
func (s *Server) SendVertifyEmai(ctx context.Context, req *pb.SendVertifyEmailRequset) (*pb.VertifyEmailResponse, error) {
	if violations := validateSendVertifyEmailReq(req); violations != nil {
		return nil, gpt_error.InvalidArgumentError(violations)
	}

	user, err := s.dao.GetUserByEmail(ctx, req.GetEmail())
	if err != nil {
		if err.Error() == gpt_error.DB_ERR_NOT_FOUND.Error() {
			return processVertifyEmailResponse(ctx, codes.OK, "successed send vertify email", fmt.Errorf("user not exists"))
		}
		return processVertifyEmailResponse(ctx, codes.Internal, "please re login", err)
	}

	if user.IsEmailVerified {
		return processVertifyEmailResponse(ctx, codes.FailedPrecondition, "user already vertified email", err)
	}

	err = s.Service.SendVertifyEmail(ctx, user.UserID.Bytes, user.Email)
	if err != nil {
		return processVertifyEmailResponse(ctx, codes.Internal, "internal err", err)
	}

	return processVertifyEmailResponse(ctx, codes.OK, "successed, please wait to receive vertify email", err)
}

func (s *Server) VertifyEmail(ctx context.Context, req *pb.VertifyEmailRequset) (*pb.VertifyEmailResponse, error) {
	if violations := validateVertifyEmailReq(req); violations != nil {
		return nil, gpt_error.InvalidArgumentError(violations)
	}

	userId, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return processVertifyEmailResponse(ctx, codes.Internal, "internal err", err)
	}

	email, err := s.dao.GetVertifyEmailByUserIdAndCode(ctx, db.GetVertifyEmailByUserIdAndCodeParams{
		UserID: pgtype.UUID{
			Bytes: userId,
			Valid: true,
		},
		SecretCode: req.GetSecretCode(),
	})

	if err != nil {
		if err.Error() == gpt_error.DB_ERR_NOT_FOUND.Error() {
			return processVertifyEmailResponse(ctx, codes.NotFound, "vertify email not exists", err)
		}
		return processVertifyEmailResponse(ctx, codes.Internal, "internal err", err)
	}

	if !email.IsValidated || email.IsUsed {
		return processVertifyEmailResponse(ctx, codes.FailedPrecondition, "vertify email is used or invalid", err)

	}

	if time.Now().UTC().After(email.ExpiredAt) {
		s.dao.UpdateVertifyEmail(ctx, db.UpdateVertifyEmailParams{
			ID: email.ID,
			IsValidated: pgtype.Bool{
				Bool:  false,
				Valid: true,
			},
		})
		return processVertifyEmailResponse(ctx, codes.FailedPrecondition, "vertify code has expired", fmt.Errorf("vertify code has expired"))
	}

	_, err = s.dao.UpdateVerifyEmailTx(ctx, db.VerifyEmailTxParams{
		ID:          email.ID,
		SecretCode:  req.GetSecretCode(),
		IsUsed:      true,
		IsValidated: false,
	})

	if err != nil {
		return processVertifyEmailResponse(ctx, codes.Internal, "internal err", err)
	}

	return processVertifyEmailResponse(ctx, codes.OK, "vertify email successed", nil)
}

func validateSendVertifyEmailReq(req *pb.SendVertifyEmailRequset) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validate.ValidateEmailFormat(req.GetEmail()); err != nil {
		violations = append(violations, validate.FieldViolation("email", err))
	}

	return violations
}

func validateVertifyEmailReq(req *pb.VertifyEmailRequset) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validate.ValidateEmptyString(req.GetUserId()); err != nil {
		violations = append(violations, validate.FieldViolation("user_id", err))
	}

	if err := validate.ValidateUUID(req.GetUserId()); err != nil {
		violations = append(violations, validate.FieldViolation("user_id", err))
	}

	if err := validate.ValidateEmptyString(req.GetSecretCode()); err != nil {
		violations = append(violations, validate.FieldViolation("secret_code", err))
	}

	return violations
}
