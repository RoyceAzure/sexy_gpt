package gpt_error

import (
	"errors"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type gError struct {
	Type    string
	Code    int
	Message string
}

func NewGerror(code int, msg string) gError {
	return gError{
		Code:    code,
		Message: msg,
	}
}

func FromError(err error) (*gError, bool) {
	if err == nil {
		return nil, true
	}

	if ge, ok := err.(gError); ok {
		temp := NewGerror(ge.Code, ge.Message)
		return &temp, true
	}
	return nil, false
}

func (g gError) Err(err error) gError {
	msg := g.Message + "," + err.Error()
	return gError{
		Code:    g.Code,
		Message: msg,
	}
}

func (g gError) ErrStr(err string) gError {
	msg := g.Message + "," + err
	return gError{
		Code:    g.Code,
		Message: msg,
	}
}

func (g gError) Error() string {
	return g.Message
}

func (g gError) Is(err *gError) bool {
	return g.Code == err.Code
}

const (
	NotFound           = 404
	InvalidArgument    = 400
	PreConditionFailed = 412
	InternalError      = 500
	Unavailable        = 503
)

var (
	ErrInValidatePreConditionOp = NewGerror(PreConditionFailed, "invalidated pre conditional op err")
	ErrInternal                 = NewGerror(InternalError, "imternal err")
	ErrInvalidArgument          = NewGerror(InvalidArgument, "invalid argument err")
	ErrNotFound                 = NewGerror(NotFound, "not found err")
	ErrUnavailable              = NewGerror(Unavailable, "service unavailable err")
)

var (
	ErrUserNotEsixts   = errors.New("user not exists or invalid password")
	ErrInvalidPassword = errors.New("wrong password")
	ERR_NOT_FOUND      = errors.New("no rows in result set")
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
		return status.Error(codes.Unauthenticated, "")
	}
	return status.Error(codes.FailedPrecondition, err.Error())
}
