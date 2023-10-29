// Code generated by MockGen. DO NOT EDIT.
// Source: auth-service/db/repo (interfaces: Repo)

// Package mockrepo is a generated GoMock package.
package mockrepo

import (
	db "auth-service/db/gen"
	context "context"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockRepo is a mock of Repo interface.
type MockRepo struct {
	ctrl     *gomock.Controller
	recorder *MockRepoMockRecorder
}

// MockRepoMockRecorder is the mock recorder for MockRepo.
type MockRepoMockRecorder struct {
	mock *MockRepo
}

// NewMockRepo creates a new mock instance.
func NewMockRepo(ctrl *gomock.Controller) *MockRepo {
	mock := &MockRepo{ctrl: ctrl}
	mock.recorder = &MockRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepo) EXPECT() *MockRepoMockRecorder {
	return m.recorder
}

// CountUsers mocks base method.
func (m *MockRepo) CountUsers(arg0 context.Context) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountUsers", arg0)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountUsers indicates an expected call of CountUsers.
func (mr *MockRepoMockRecorder) CountUsers(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountUsers", reflect.TypeOf((*MockRepo)(nil).CountUsers), arg0)
}

// CreateOAuthToken mocks base method.
func (m *MockRepo) CreateOAuthToken(arg0 context.Context, arg1 *db.CreateOAuthTokenParams) (*db.OauthToken, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOAuthToken", arg0, arg1)
	ret0, _ := ret[0].(*db.OauthToken)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateOAuthToken indicates an expected call of CreateOAuthToken.
func (mr *MockRepoMockRecorder) CreateOAuthToken(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOAuthToken", reflect.TypeOf((*MockRepo)(nil).CreateOAuthToken), arg0, arg1)
}

// CreateOTPAuth mocks base method.
func (m *MockRepo) CreateOTPAuth(arg0 context.Context, arg1 *db.CreateOTPAuthParams) (*db.OtpAuthentication, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOTPAuth", arg0, arg1)
	ret0, _ := ret[0].(*db.OtpAuthentication)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateOTPAuth indicates an expected call of CreateOTPAuth.
func (mr *MockRepoMockRecorder) CreateOTPAuth(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOTPAuth", reflect.TypeOf((*MockRepo)(nil).CreateOTPAuth), arg0, arg1)
}

// CreateUser mocks base method.
func (m *MockRepo) CreateUser(arg0 context.Context, arg1 *db.CreateUserParams) (*db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", arg0, arg1)
	ret0, _ := ret[0].(*db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockRepoMockRecorder) CreateUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockRepo)(nil).CreateUser), arg0, arg1)
}

// CreateUserTx mocks base method.
func (m *MockRepo) CreateUserTx(arg0 context.Context, arg1 int32, arg2 *db.CreateUserParams) (*db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUserTx", arg0, arg1, arg2)
	ret0, _ := ret[0].(*db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUserTx indicates an expected call of CreateUserTx.
func (mr *MockRepoMockRecorder) CreateUserTx(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUserTx", reflect.TypeOf((*MockRepo)(nil).CreateUserTx), arg0, arg1, arg2)
}

// DeleteOAuthToken mocks base method.
func (m *MockRepo) DeleteOAuthToken(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteOAuthToken", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteOAuthToken indicates an expected call of DeleteOAuthToken.
func (mr *MockRepoMockRecorder) DeleteOAuthToken(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteOAuthToken", reflect.TypeOf((*MockRepo)(nil).DeleteOAuthToken), arg0, arg1)
}

// DeleteOTPAuthByID mocks base method.
func (m *MockRepo) DeleteOTPAuthByID(arg0 context.Context, arg1 int32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteOTPAuthByID", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteOTPAuthByID indicates an expected call of DeleteOTPAuthByID.
func (mr *MockRepoMockRecorder) DeleteOTPAuthByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteOTPAuthByID", reflect.TypeOf((*MockRepo)(nil).DeleteOTPAuthByID), arg0, arg1)
}

// DeleteOTPAuthByPhone mocks base method.
func (m *MockRepo) DeleteOTPAuthByPhone(arg0 context.Context, arg1 *db.DeleteOTPAuthByPhoneParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteOTPAuthByPhone", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteOTPAuthByPhone indicates an expected call of DeleteOTPAuthByPhone.
func (mr *MockRepoMockRecorder) DeleteOTPAuthByPhone(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteOTPAuthByPhone", reflect.TypeOf((*MockRepo)(nil).DeleteOTPAuthByPhone), arg0, arg1)
}

// ExecTx mocks base method.
func (m *MockRepo) ExecTx(arg0 context.Context, arg1 func(*db.Queries) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExecTx", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ExecTx indicates an expected call of ExecTx.
func (mr *MockRepoMockRecorder) ExecTx(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExecTx", reflect.TypeOf((*MockRepo)(nil).ExecTx), arg0, arg1)
}

// Exists mocks base method.
func (m *MockRepo) Exists(arg0 context.Context, arg1 ...string) (bool, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Exists", varargs...)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exists indicates an expected call of Exists.
func (mr *MockRepoMockRecorder) Exists(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exists", reflect.TypeOf((*MockRepo)(nil).Exists), varargs...)
}

// Get mocks base method.
func (m *MockRepo) Get(arg0 context.Context, arg1 string) (interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockRepoMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockRepo)(nil).Get), arg0, arg1)
}

// GetOAuthToken mocks base method.
func (m *MockRepo) GetOAuthToken(arg0 context.Context, arg1 string) (*db.OauthToken, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOAuthToken", arg0, arg1)
	ret0, _ := ret[0].(*db.OauthToken)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOAuthToken indicates an expected call of GetOAuthToken.
func (mr *MockRepoMockRecorder) GetOAuthToken(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOAuthToken", reflect.TypeOf((*MockRepo)(nil).GetOAuthToken), arg0, arg1)
}

// GetOTPAuth mocks base method.
func (m *MockRepo) GetOTPAuth(arg0 context.Context, arg1 *db.GetOTPAuthParams) (*db.OtpAuthentication, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOTPAuth", arg0, arg1)
	ret0, _ := ret[0].(*db.OtpAuthentication)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOTPAuth indicates an expected call of GetOTPAuth.
func (mr *MockRepoMockRecorder) GetOTPAuth(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOTPAuth", reflect.TypeOf((*MockRepo)(nil).GetOTPAuth), arg0, arg1)
}

// GetUser mocks base method.
func (m *MockRepo) GetUser(arg0 context.Context, arg1 string) (*db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", arg0, arg1)
	ret0, _ := ret[0].(*db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockRepoMockRecorder) GetUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockRepo)(nil).GetUser), arg0, arg1)
}

// GetUserByEmail mocks base method.
func (m *MockRepo) GetUserByEmail(arg0 context.Context, arg1 *string) (*db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByEmail", arg0, arg1)
	ret0, _ := ret[0].(*db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByEmail indicates an expected call of GetUserByEmail.
func (mr *MockRepoMockRecorder) GetUserByEmail(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByEmail", reflect.TypeOf((*MockRepo)(nil).GetUserByEmail), arg0, arg1)
}

// GetUserByPhone mocks base method.
func (m *MockRepo) GetUserByPhone(arg0 context.Context, arg1 string) (*db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByPhone", arg0, arg1)
	ret0, _ := ret[0].(*db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByPhone indicates an expected call of GetUserByPhone.
func (mr *MockRepoMockRecorder) GetUserByPhone(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByPhone", reflect.TypeOf((*MockRepo)(nil).GetUserByPhone), arg0, arg1)
}

// HasOTPAuthValid mocks base method.
func (m *MockRepo) HasOTPAuthValid(arg0 context.Context, arg1 *db.HasOTPAuthValidParams) (int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasOTPAuthValid", arg0, arg1)
	ret0, _ := ret[0].(int32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HasOTPAuthValid indicates an expected call of HasOTPAuthValid.
func (mr *MockRepoMockRecorder) HasOTPAuthValid(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasOTPAuthValid", reflect.TypeOf((*MockRepo)(nil).HasOTPAuthValid), arg0, arg1)
}

// ListUsers mocks base method.
func (m *MockRepo) ListUsers(arg0 context.Context, arg1 *db.ListUsersParams) ([]*db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListUsers", arg0, arg1)
	ret0, _ := ret[0].([]*db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListUsers indicates an expected call of ListUsers.
func (mr *MockRepoMockRecorder) ListUsers(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListUsers", reflect.TypeOf((*MockRepo)(nil).ListUsers), arg0, arg1)
}

// Set mocks base method.
func (m *MockRepo) Set(arg0 context.Context, arg1 string, arg2 interface{}, arg3 time.Duration) (interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Set indicates an expected call of Set.
func (mr *MockRepoMockRecorder) Set(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockRepo)(nil).Set), arg0, arg1, arg2, arg3)
}

// UpdateUser mocks base method.
func (m *MockRepo) UpdateUser(arg0 context.Context, arg1 *db.UpdateUserParams) (*db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", arg0, arg1)
	ret0, _ := ret[0].(*db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockRepoMockRecorder) UpdateUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockRepo)(nil).UpdateUser), arg0, arg1)
}