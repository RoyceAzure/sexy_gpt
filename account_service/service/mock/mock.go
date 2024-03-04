package mock_service

import (
	context "context"
	reflect "reflect"

	token "github.com/RoyceAzure/sexy_gpt/account_service/api/gapi/token"
	db "github.com/RoyceAzure/sexy_gpt/account_service/repository/db/sqlc"
	config "github.com/RoyceAzure/sexy_gpt/account_service/shared/util/config"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockIService is a mock of IService interface.
type MockIService struct {
	ctrl     *gomock.Controller
	recorder *MockIServiceMockRecorder
}

// MockIServiceMockRecorder is the mock recorder for MockIService.
type MockIServiceMockRecorder struct {
	mock *MockIService
}

// NewMockIService creates a new mock instance.
func NewMockIService(ctrl *gomock.Controller) *MockIService {
	mock := &MockIService{ctrl: ctrl}
	mock.recorder = &MockIServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIService) EXPECT() *MockIServiceMockRecorder {
	return m.recorder
}

// IsUserLogin mocks base method.
func (m *MockIService) IsUserLogin(arg0 context.Context, arg1 uuid.UUID) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsUserLogin", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsUserLogin indicates an expected call of IsUserLogin.
func (mr *MockIServiceMockRecorder) IsUserLogin(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsUserLogin", reflect.TypeOf((*MockIService)(nil).IsUserLogin), arg0, arg1)
}

// IsValidateUser mocks base method.
func (m *MockIService) IsValidateUser(arg0 context.Context, arg1 string) (*db.UserRoleView, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsValidateUser", arg0, arg1)
	ret0, _ := ret[0].(*db.UserRoleView)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsValidateUser indicates an expected call of IsValidateUser.
func (mr *MockIServiceMockRecorder) IsValidateUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsValidateUser", reflect.TypeOf((*MockIService)(nil).IsValidateUser), arg0, arg1)
}

// LoginCreateSession mocks base method.
func (m *MockIService) LoginCreateSession(arg0 context.Context, arg1 uuid.UUID, arg2 token.Maker, arg3 config.Config) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoginCreateSession", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoginCreateSession indicates an expected call of LoginCreateSession.
func (mr *MockIServiceMockRecorder) LoginCreateSession(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoginCreateSession", reflect.TypeOf((*MockIService)(nil).LoginCreateSession), arg0, arg1, arg2, arg3)
}

// SendVertifyEmail mocks base method.
func (m *MockIService) SendVertifyEmail(arg0 context.Context, arg1 uuid.UUID, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendVertifyEmail", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendVertifyEmail indicates an expected call of SendVertifyEmail.
func (mr *MockIServiceMockRecorder) SendVertifyEmail(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendVertifyEmail", reflect.TypeOf((*MockIService)(nil).SendVertifyEmail), arg0, arg1, arg2)
}
