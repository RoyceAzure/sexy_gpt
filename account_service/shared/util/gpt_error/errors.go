package gpt_error

import (
	"errors"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInValidatePreConditionOp = errors.New("invalidated pre conditional op err")
	ErrUnauthicated             = errors.New("unauthicated err")
	ErrInternal                 = errors.New("internal err")
	ErrInvalidArgument          = errors.New("invalid argument err")
	ErrNotFound                 = errors.New("not found err")
	ErrUnavailable              = errors.New("service unavailable err")
	ErrInvalidSession           = errors.New("session invalidated")
)

var (
	DB_ERR_NOT_FOUND = errors.New("no rows in result set")
)

const (
	ForeignKeyViolation   = "foreign_key_violation"
	PgErr_UniqueViolation = "23505"

	DEFAULT_PAGE            = 1
	DEFAULT_PAGE_SIZE       = 10
	SELL                    = "sell"
	BUY                     = "buy"
	AuthorizationHeaderKey  = "authorization"
	AuthorizationTypeBearer = "bearer"
	AuthorizationPayloadKey = "authorization_payload"
)

func FieldViolation(field string, err error) *errdetails.BadRequest_FieldViolation {
	return &errdetails.BadRequest_FieldViolation{
		Field:       field,
		Description: err.Error(),
	}
}

// pacakge list of errdetails.BadRequest_FieldViolation use errdetails.BadRequest,
// use *status.status set code, with badrequest
// then use *status.status to generate error string
func InvalidArgumentError(violations []*errdetails.BadRequest_FieldViolation) error {
	badRequest := &errdetails.BadRequest{FieldViolations: violations}
	statusInvalid := status.New(codes.InvalidArgument, "Invalid parameters")
	statusDetails, err := statusInvalid.WithDetails(badRequest)
	if err != nil {
		return statusInvalid.Err()
	}
	return statusDetails.Err()
}

// API
func APIUnauthticatedError(err error) error {
	if err == nil {
		return status.Error(codes.Unauthenticated, "")
	}
	return status.Error(codes.Unauthenticated, err.Error())
}

func APIInternalError(err error) error {
	if err == nil {
		return status.Error(codes.Unauthenticated, "")
	}
	return status.Error(codes.Internal, err.Error())
}

func APIInValidateOperation(err error) error {
	if err == nil {
		return status.Error(codes.FailedPrecondition, "")
	}
	return status.Error(codes.FailedPrecondition, err.Error())
}
