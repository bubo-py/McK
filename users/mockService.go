// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/bubo-py/McK/users/service (interfaces: BusinessLogicInterface)

// Package mocks is a generated GoMock package.
package users

import (
	context "context"
	reflect "reflect"

	types "github.com/bubo-py/McK/types"
	gomock "github.com/golang/mock/gomock"
)

// MockBusinessLogicInterface is a mock of BusinessLogicInterface interface.
type MockBusinessLogicInterface struct {
	ctrl     *gomock.Controller
	recorder *MockBusinessLogicInterfaceMockRecorder
}

// MockBusinessLogicInterfaceMockRecorder is the mock recorder for MockBusinessLogicInterface.
type MockBusinessLogicInterfaceMockRecorder struct {
	mock *MockBusinessLogicInterface
}

// NewMockBusinessLogicInterface creates a new mock instance.
func NewMockBusinessLogicInterface(ctrl *gomock.Controller) *MockBusinessLogicInterface {
	mock := &MockBusinessLogicInterface{ctrl: ctrl}
	mock.recorder = &MockBusinessLogicInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBusinessLogicInterface) EXPECT() *MockBusinessLogicInterfaceMockRecorder {
	return m.recorder
}

// AddUser mocks base method.
func (m *MockBusinessLogicInterface) AddUser(arg0 context.Context, arg1 types.User) (types.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUser", arg0, arg1)
	ret0, _ := ret[0].(types.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddUser indicates an expected call of AddUser.
func (mr *MockBusinessLogicInterfaceMockRecorder) AddUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUser", reflect.TypeOf((*MockBusinessLogicInterface)(nil).AddUser), arg0, arg1)
}

// DeleteUser mocks base method.
func (m *MockBusinessLogicInterface) DeleteUser(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockBusinessLogicInterfaceMockRecorder) DeleteUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockBusinessLogicInterface)(nil).DeleteUser), arg0, arg1)
}

// GetUserByLogin mocks base method.
func (m *MockBusinessLogicInterface) GetUserByLogin(arg0 context.Context, arg1 string) (types.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByLogin", arg0, arg1)
	ret0, _ := ret[0].(types.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByLogin indicates an expected call of GetUserByLogin.
func (mr *MockBusinessLogicInterfaceMockRecorder) GetUserByLogin(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByLogin", reflect.TypeOf((*MockBusinessLogicInterface)(nil).GetUserByLogin), arg0, arg1)
}

// LoginUser mocks base method.
func (m *MockBusinessLogicInterface) LoginUser(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoginUser", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// LoginUser indicates an expected call of LoginUser.
func (mr *MockBusinessLogicInterfaceMockRecorder) LoginUser(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoginUser", reflect.TypeOf((*MockBusinessLogicInterface)(nil).LoginUser), arg0, arg1, arg2)
}

// UpdateUser mocks base method.
func (m *MockBusinessLogicInterface) UpdateUser(arg0 context.Context, arg1 types.User, arg2 int64) (types.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", arg0, arg1, arg2)
	ret0, _ := ret[0].(types.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockBusinessLogicInterfaceMockRecorder) UpdateUser(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockBusinessLogicInterface)(nil).UpdateUser), arg0, arg1, arg2)
}
