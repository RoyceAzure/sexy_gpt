package gapi

import (
	"context"
	"errors"
	"fmt"
	"time"

	db "github.com/RoyceAzure/sexy_gpt/account_service/repository/db/sqlc"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/pb"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util/converter"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util/gpt_error"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util/validate"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
)

/*
 */
func processResponse(ctx context.Context, code codes.Code, msg string, err error) (*pb.UserDTOResponse, error) {
	res := pb.UserDTOResponse{Message: msg}
	if code != codes.OK {
		if err == nil {
			err = fmt.Errorf(msg)
		}
		util.NewOutGoingMetaDataKV(ctx, util.DBMSGKey, err.Error())
		return &res, status.Errorf(code, msg)
	}
	return &res, nil
}

func processResponses(ctx context.Context, code codes.Code, msg string, err error) (*pb.UserDTOsResponse, error) {
	res := pb.UserDTOsResponse{Message: msg}
	if code != codes.OK {
		util.NewOutGoingMetaDataKV(ctx, util.DBMSGKey, err.Error())
		return &res, status.Errorf(code, msg)
	}
	return &res, nil
}

/*
將發送認證信提出事務之外
*/
func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserDTOResponse, error) {
	var response pb.UserDTOResponse
	if violations := validateCreateUserReqV2(req); violations != nil {
		msg, err := violations.ToJson()
		if err != nil {
			msg = "Invalidate argument"
		}
		return processResponse(ctx, codes.InvalidArgument, msg, err)
	}

	hashPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return processResponse(ctx, codes.Internal, "internal err", err)
	}

	arg := db.CreateUserTxParms{
		Arg: &db.CreateUserParams{
			UserName:       req.GetUserName(),
			Email:          req.GetEmail(),
			HashedPassword: hashPassword,
			IsInternal:     false,
			CrUser:         "SYSTEM",
		},
	}

	TxResult, err := s.dao.CreateUserTx(ctx, &arg)
	if err != nil {
		if pqErr, ok := err.(*pgconn.PgError); ok {
			switch pqErr.Code {
			case gpt_error.PgErr_UniqueViolation:
				return processResponse(ctx, codes.AlreadyExists, "user already exists", err)
			}
		}
		return processResponse(ctx, codes.Internal, "internal err", err)
	}

	err = s.Service.SendVertifyEmail(ctx, TxResult.User.UserID.Bytes, TxResult.User.Email)
	if err != nil {
		return processResponse(ctx, codes.Internal, "internal err", err)
	}

	response.Data = &pb.UserDTO{
		UserId:   converter.ConvertXByte2UUID(TxResult.User.UserID.Bytes),
		UserName: TxResult.User.UserName,
		Email:    TxResult.User.Email,
		RoleName: TxResult.Role.RoleName,
	}

	if err != nil {
		return processResponse(ctx, codes.Internal, "internal err", err)
	}
	return &response, nil
}

/*
Internal 可以查所有  其餘只能查自己
*/
func (s *Server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserDTOResponse, error) {
	payload, err := s.authorizUser(ctx)
	if err != nil {
		if errors.Is(err, gpt_error.ErrInternal) {
			return processResponse(ctx, codes.Internal, "internal err", err)
		} else {
			return processResponse(ctx, codes.Unauthenticated, err.Error(), err)
		}
	}

	if violations := validateGetUserReqV2(req); violations != nil {
		// return nil, gpt_error.InvalidArgumentError(violations)
		msg, err := violations.ToJson()
		if err != nil {
			msg = "Invalidate argument"
		}
		return processResponse(ctx, codes.InvalidArgument, msg, err)
	}

	var uid uuid.UUID
	if req.GetUserId() == "" {
		uid = payload.UserId
	} else {
		if !payload.IsInternal {
			uid = payload.UserId
		} else {
			uid, _ = uuid.Parse(req.GetUserId())
		}
	}

	user, err := s.dao.GetUserDTO(ctx, pgtype.UUID{
		Bytes: uid,
		Valid: true,
	})

	if err != nil {
		if err.Error() == gpt_error.DB_ERR_NOT_FOUND.Error() {
			return processResponse(ctx, codes.NotFound, "user not found", err)
		}
		return processResponse(ctx, codes.Internal, "internal err", err)
	}
	var response pb.UserDTOResponse
	response.Data = ConvertUserDTO(&user)

	return &response, nil
}

func (s *Server) GetUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.UserDTOsResponse, error) {
	payload, err := s.authorizUser(ctx)
	if err != nil {
		if errors.Is(err, gpt_error.ErrInternal) {
			return processResponses(ctx, codes.Internal, "internal err", err)
		} else {
			return processResponses(ctx, codes.Unauthenticated, err.Error(), err)
		}
	}
	if !payload.IsInternal {
		return processResponses(ctx, codes.PermissionDenied, "user is not internal", err)
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

	var response pb.UserDTOsResponse
	response.Data = ConvertUserDTOs(users)

	return &response, nil
}

func (s *Server) GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.UserDTOResponse, error) {
	payload, err := s.authorizUser(ctx)
	if err != nil {
		if errors.Is(err, gpt_error.ErrInternal) {
			return processResponse(ctx, codes.Internal, "internal err", err)
		} else {
			return processResponse(ctx, codes.Unauthenticated, err.Error(), err)
		}
	}

	email := req.GetEmail()
	if email == "" {
		email = payload.Email
	} else if !payload.IsInternal {
		email = payload.Email
	}

	user, err := s.dao.GetUserDTOByEmail(ctx, email)

	if err != nil {
		if err.Error() == gpt_error.DB_ERR_NOT_FOUND.Error() {
			return processResponse(ctx, codes.NotFound, "user not exists", err)
		}
		return processResponse(ctx, codes.Internal, "internal err", err)
	}

	var response pb.UserDTOResponse
	response.Data = ConvertUserDTO(&user)

	return &response, nil
}

/*
only allow user updated user_name
*/
func (s *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserDTOResponse, error) {
	payload, err := s.authorizUser(ctx)
	if err != nil {
		if errors.Is(err, gpt_error.ErrInternal) {
			return processResponse(ctx, codes.Internal, "internal err", err)
		} else {
			return processResponse(ctx, codes.Unauthenticated, err.Error(), err)
		}
	}
	if violations := validateUpdateUserReqV2(req); violations != nil {
		msg, err := violations.ToJson()
		if err != nil {
			msg = "Invalidate argument"
		}
		return processResponse(ctx, codes.InvalidArgument, msg, err)
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
		UserName: pgtype.Text{
			String: req.GetUserName(),
			Valid:  true,
		},
		UpDate: pgtype.Timestamptz{
			Time:  time.Now().UTC(),
			Valid: true,
		},
	}

	updatedUser, err := s.dao.UpdateUser(ctx, arg)
	if err != nil {
		if err.Error() == gpt_error.DB_ERR_NOT_FOUND.Error() {
			return processResponse(ctx, codes.NotFound, "user not exists", err)
		}
		return processResponse(ctx, codes.Internal, "internal err", err)
	}

	var response pb.UserDTOResponse
	response.Data = &pb.UserDTO{
		UserId:   converter.ConvertXByte2UUID(updatedUser.UserID.Bytes),
		UserName: updatedUser.UserName,
		Email:    updatedUser.Email,
	}

	return &response, nil

}

func (s *Server) UpdateUserPas(ctx context.Context, req *pb.UpdateUserPasRequest) (*pb.UserDTOResponse, error) {
	payload, err := s.authorizUser(ctx)
	if err != nil {
		if errors.Is(err, gpt_error.ErrInternal) {
			return processResponse(ctx, codes.Internal, "internal err", err)
		} else {
			return processResponse(ctx, codes.Unauthenticated, err.Error(), err)
		}
	}
	if violations := validateUpdateUserPasReqV2(req); violations != nil {
		msg, err := violations.ToJson()
		if err != nil {
			msg = "Invalidate argument"
		}
		return processResponse(ctx, codes.InvalidArgument, msg, err)
	}

	uid, _ := uuid.Parse(req.GetUserId())
	if payload.UserId != uid {
		return processResponse(ctx, codes.Unauthenticated, "unauthenticated err", nil)
	}

	hashPas, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return processResponse(ctx, codes.Internal, "internal err", err)
	}

	arg := db.UpdateUserParams{
		UserID: pgtype.UUID{
			Bytes: uid,
			Valid: true,
		},
		HashedPassword: pgtype.Text{
			String: hashPas,
			Valid:  true,
		},
		UpDate: pgtype.Timestamptz{
			Time:  time.Now().UTC(),
			Valid: true,
		},
	}

	_, err = s.dao.UpdateUser(ctx, arg)
	if err != nil {
		if err.Error() == gpt_error.DB_ERR_NOT_FOUND.Error() {
			return processResponse(ctx, codes.NotFound, "user not exists", err)
		}
		return processResponse(ctx, codes.Internal, "internal err", err)
	}

	var response pb.UserDTOResponse
	response.Message = "successed updated password, please re login"
	return &response, nil
}

func ConvertUser(user db.User) *pb.User {
	return &pb.User{
		UserId:   converter.ConvertXByte2UUID(user.UserID.Bytes),
		UserName: user.UserName,
		Email:    user.Email,
	}
}

func ConvertUserDTO(user *db.UserRoleView) *pb.UserDTO {
	if user == nil {
		return nil
	}
	return &pb.UserDTO{
		UserId:   converter.ConvertXByte2UUID(user.UserID.Bytes),
		UserName: user.UserName,
		Email:    user.Email,
		RoleName: user.RoleName.String,
	}
}

func ConvertUserDTOs(users []db.UserRoleView) (res []*pb.UserDTO) {
	for _, user := range users {
		res = append(res, ConvertUserDTO(&user))
	}
	return res
}

func convertUserPb(user db.UserRoleView) (*anypb.Any, error) {
	dto := ConvertUserDTO(&user)
	return anypb.New(dto)
}

func convertUserPbs(users []db.UserRoleView) (res []*anypb.Any, err error) {
	for _, user := range users {
		pb, err := convertUserPb(user)
		if err != nil {
			return res, err
		}
		res = append(res, pb)
	}
	return res, err
}

func convertUsers(users []db.User) (res []*pb.User) {
	for _, user := range users {
		res = append(res, ConvertUser(user))
	}
	return res
}

func validateCreateUserReqV2(req *pb.CreateUserRequest) *gpt_error.Message {
	violations := &gpt_error.Message{}
	if err := validate.ValidateEmailFormat(req.GetEmail()); err != nil {
		gpt_error.AddErrFieldString(violations, gpt_error.NewErrField("email", err))
	}

	if err := validate.ValidateStrongPas(req.GetPassword()); err != nil {
		gpt_error.AddErrFieldString(violations, gpt_error.NewErrField("password", err...))
	}

	if err := validate.ValidateEmptyString(req.GetUserName()); err != nil {
		gpt_error.AddErrFieldString(violations, gpt_error.NewErrField("user_name", err))
	}
	if len(violations.ErrMessage) == 0 {
		return nil
	}
	return violations
}

func validateGetUserReqV2(req *pb.GetUserRequest) *gpt_error.Message {
	violations := &gpt_error.Message{}
	if req.GetUserId() != "" {
		if err := validate.ValidateUUID(req.GetUserId()); err != nil {
			gpt_error.AddErrFieldString(violations, gpt_error.NewErrField("user_id", err))
		}
	}
	if len(violations.ErrMessage) == 0 {
		return nil
	}
	return violations
}

func validateUpdateUserReqV2(req *pb.UpdateUserRequest) *gpt_error.Message {
	violations := &gpt_error.Message{}
	if err := validate.ValidateEmptyString(req.GetUserId()); err != nil {
		gpt_error.AddErrFieldString(violations, gpt_error.NewErrField("user_id", err))
	}
	if err := validate.ValidateUUID(req.GetUserId()); err != nil {
		gpt_error.AddErrFieldString(violations, gpt_error.NewErrField("user_id", err))
	}
	if err := validate.ValidateEmptyString(req.GetUserName()); err != nil {
		gpt_error.AddErrFieldString(violations, gpt_error.NewErrField("user_name", err))
	}
	if len(violations.ErrMessage) == 0 {
		return nil
	}
	return violations
}

func validateUpdateUserPasReqV2(req *pb.UpdateUserPasRequest) *gpt_error.Message {
	violations := &gpt_error.Message{}
	if err := validate.ValidateEmptyString(req.GetUserId()); err != nil {
		gpt_error.AddErrFieldString(violations, gpt_error.NewErrField("user_id", err))
	}
	if err := validate.ValidateUUID(req.GetUserId()); err != nil {
		gpt_error.AddErrFieldString(violations, gpt_error.NewErrField("user_id", err))
	}
	if err := validate.ValidateStrongPas(req.GetPassword()); err != nil {
		gpt_error.AddErrFieldString(violations, gpt_error.NewErrField("password", err...))
	}
	if len(violations.ErrMessage) == 0 {
		return nil
	}
	return violations
}
