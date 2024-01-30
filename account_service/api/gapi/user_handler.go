package gapi

import (
	"context"

	db "github.com/RoyceAzure/sexy_gpt/account_service/repository/db/sqlc"
	logger "github.com/RoyceAzure/sexy_gpt/account_service/repository/logger_distributor"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/pb"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util/converter"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util/gpt_error"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util/validate"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	if violations := validateCreateUserReq(req); violations != nil {
		return nil, gpt_error.InvalidArgumentError(violations)
	}

	hashPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, s.HandleAPIError(codes.Internal, err, "internal err")
	}

	arg := db.CreateUserParams{
		UserName:       req.GetUserName(),
		Email:          req.GetEmail(),
		HashedPassword: hashPassword,
		IsInternal:     false,
		CrUser:         "SYSTEM",
	}

	TxResult, err := s.dao.CreateUserTx(ctx, &arg)
	if err != nil {
		if pqErr, ok := err.(*pgconn.PgError); ok {
			switch pqErr.Code {
			case gpt_error.PgErr_UniqueViolation:
				return nil, s.HandleAPIError(codes.AlreadyExists, err, "user already exists")
			}
		}
		return nil, s.HandleAPIError(codes.Internal, err, "internal err")
	}

	res := &pb.CreateUserResponse{
		Data: &pb.UserDTO{
			UserId:   converter.ConvertXByte2UUID(TxResult.User.UserID.Bytes),
			UserName: TxResult.User.UserName,
			Email:    TxResult.User.Email,
			RoleName: TxResult.Role.RoleName,
		},
	}

	return res, nil
}

/*
Internal 可以查所有  其餘只能查自己
*/
func (s *Server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	payload, err := s.authorizUser(ctx)
	if err != nil {
		if ge, ok := gpt_error.FromError(err); ok {
			if ge.Is(&gpt_error.ErrInternal) {
				return nil, s.HandleAPIError(codes.Internal, err, "internal err")
			} else if ge.Is(&gpt_error.ErrInvalidSession) {
				return nil, s.HandleAPIError(codes.Unauthenticated, err, "please re login")
			}
		}
		return nil, s.HandleAPIError(codes.Internal, err, "internal err")
	}

	if violations := validateGetUserReq(req); violations != nil {
		return nil, gpt_error.InvalidArgumentError(violations)
	}

	var uid uuid.UUID
	if payload.IsInternal {
		uid, _ = uuid.Parse(req.GetUserId())
	} else {
		uid = payload.UserId
	}

	user, err := s.dao.GetUserDTO(ctx, pgtype.UUID{
		Bytes: uid,
		Valid: true,
	})

	if err != nil {
		if err.Error() == gpt_error.DB_ERR_NOT_FOUND.Error() {
			return nil, status.Errorf(codes.NotFound, "user not exists")
		}
		logger.Logger.Error().Err(err).Msg("failed to get user")
		return nil, s.HandleAPIError(codes.Internal, err, "internal err")
	}

	res := &pb.GetUserResponse{
		Data: convertUserDTO(user),
	}

	return res, nil
}

/*
詳細err 紀錄log
*/

func (s *Server) GetUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	payload, err := s.authorizUser(ctx)
	if err != nil {
		if ge, ok := gpt_error.FromError(err); ok {
			if ge.Is(&gpt_error.ErrInternal) {
				return nil, s.HandleAPIError(codes.Internal, err, "internal err")
			} else if ge.Is(&gpt_error.ErrInvalidSession) {
				return nil, s.HandleAPIError(codes.Unauthenticated, err, "please re login")
			}
		}
		return nil, s.HandleAPIError(codes.Internal, err, "internal err")
	}
	if !payload.IsInternal {
		return nil, gpt_error.APIUnauthticatedError(err)
	}

	pgsize := req.GetPageSize()
	page := req.GetPage()

	if pgsize <= 0 {
		pgsize = 10
	}

	if page <= 0 {
		page = 1
	}

	users, err := s.dao.GetUsersDTO(ctx, db.GetUsersDTOParams{
		Limit:  pgsize,
		Offset: (page - 1) * pgsize,
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get users")
	}

	res := &pb.GetUsersResponse{
		Data: convertUserDTOs(users),
	}

	return res, nil
}

func (s *Server) GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.GetUserByEmailResponse, error) {
	payload, err := s.authorizUser(ctx)
	if err != nil {
		if ge, ok := gpt_error.FromError(err); ok {
			if ge.Is(&gpt_error.ErrInternal) {
				return nil, s.HandleAPIError(codes.Internal, err, "internal err")
			} else if ge.Is(&gpt_error.ErrInvalidSession) {
				return nil, s.HandleAPIError(codes.Unauthenticated, err, "please re login")
			}
		}
		return nil, s.HandleAPIError(codes.Internal, err, "internal err")
	}
	if violations := validateGetUserByEmailReq(req); violations != nil {
		return nil, gpt_error.InvalidArgumentError(violations)
	}

	var email string
	if payload.IsInternal {
		email = req.GetEmail()
	} else {
		email = payload.Email
	}

	user, err := s.dao.GetUserDTOByEmail(ctx, email)

	if err != nil {
		if err.Error() == gpt_error.DB_ERR_NOT_FOUND.Error() {
			return nil, status.Errorf(codes.NotFound, "user not exists")
		}
		logger.Logger.Error().Err(err).Msg("failed to get user by email")
		return nil, s.HandleAPIError(codes.Internal, err, "internal err")
	}

	res := &pb.GetUserByEmailResponse{
		Data: convertUserDTO(user),
	}

	return res, nil
}

/*
only allow user updated email and password
*/
func (s *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	payload, err := s.authorizUser(ctx)
	if err != nil {
		if ge, ok := gpt_error.FromError(err); ok {
			if ge.Is(&gpt_error.ErrInternal) {
				return nil, s.HandleAPIError(codes.Internal, err, "internal err")
			} else if ge.Is(&gpt_error.ErrInvalidSession) {
				return nil, s.HandleAPIError(codes.Unauthenticated, err, "please re login")
			}
		}
		return nil, s.HandleAPIError(codes.Internal, err, "internal err")
	}
	if violations := validateUpdateUserReq(req); violations != nil {
		return nil, gpt_error.InvalidArgumentError(violations)
	}

	var uid uuid.UUID
	if payload.IsInternal {
		uid, _ = uuid.Parse(req.GetUserId())
	} else {
		uid = payload.UserId
	}

	arg := db.UpdateUserParams{
		UserID: pgtype.UUID{
			Bytes: uid,
			Valid: true,
		},
	}

	if req.UserName != nil {
		arg.UserName = pgtype.Text{
			String: req.GetUserName(),
			Valid:  true,
		}
	}

	if req.Password != nil {
		hashPas, err := util.HashPassword(req.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hash password")
		}
		arg.HashedPassword = pgtype.Text{
			String: hashPas,
			Valid:  true,
		}
	}

	if req.Password == nil && req.UserName == nil {
		return nil, status.Errorf(codes.InvalidArgument, "nothing change")
	}

	updatedUser, err := s.dao.UpdateUser(ctx, arg)
	if err != nil {
		if err.Error() == gpt_error.DB_ERR_NOT_FOUND.Error() {
			return nil, status.Errorf(codes.NotFound, "user not exists")
		}
		logger.Logger.Error().Err(err).Msg("failed to updated user")
		return nil, status.Errorf(codes.Internal, "internal err")
	}

	res := &pb.UpdateUserResponse{
		Data: convertUser(updatedUser),
	}

	return res, nil
}

func convertUser(user db.User) *pb.User {
	return &pb.User{
		UserId:   converter.ConvertXByte2UUID(user.UserID.Bytes),
		UserName: user.UserName,
		Email:    user.Email,
	}
}

func convertUserDTO(user db.UserRoleView) *pb.UserDTO {
	return &pb.UserDTO{
		UserId:   converter.ConvertXByte2UUID(user.UserID.Bytes),
		UserName: user.UserName,
		Email:    user.Email,
		RoleName: user.RoleName.String,
	}
}

func convertUserDTOs(users []db.UserRoleView) (res []*pb.UserDTO) {
	for _, user := range users {
		res = append(res, convertUserDTO(user))
	}
	return res
}

func convertUsers(users []db.User) (res []*pb.User) {
	for _, user := range users {
		res = append(res, convertUser(user))
	}
	return res
}

func validateCreateUserReq(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validate.ValidateEmailFormat(req.GetEmail()); err != nil {
		violations = append(violations, validate.FieldViolation("email", err))
	}

	if err := validate.ValidateStrongPas(req.GetPassword()); err != nil {
		violations = append(violations, validate.FieldViolations("password", err))
	}

	if err := validate.ValidateEmptyString(req.GetUserName()); err != nil {
		violations = append(violations, validate.FieldViolation("user_name", err))
	}
	return violations
}

func validateGetUserReq(req *pb.GetUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validate.ValidateEmptyString(req.GetUserId()); err != nil {
		violations = append(violations, validate.FieldViolation("users_id", err))
	}
	if err := validate.ValidateUUID(req.GetUserId()); err != nil {
		violations = append(violations, validate.FieldViolation("users_id", err))
	}

	return violations
}

func validateGetUserByEmailReq(req *pb.GetUserByEmailRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validate.ValidateEmptyString(req.GetEmail()); err != nil {
		violations = append(violations, validate.FieldViolation("email", err))
	}
	return violations
}

func validateUpdateUserReq(req *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validate.ValidateEmptyString(req.GetUserId()); err != nil {
		violations = append(violations, validate.FieldViolation("users_id", err))
	}
	if err := validate.ValidateUUID(req.GetUserId()); err != nil {
		violations = append(violations, validate.FieldViolation("users_id", err))
	}

	if email := req.GetEmail(); email != "" {
		if err := validate.ValidateEmailFormat(email); err != nil {
			violations = append(violations, validate.FieldViolation("email", err))
		}
	}

	if req.Password != nil {
		if err := validate.ValidateStrongPas(req.GetPassword()); err != nil {
			violations = append(violations, validate.FieldViolations("password", err))
		}
	}

	return violations
}
