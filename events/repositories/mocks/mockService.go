// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/bubo-py/McK/events/service (interfaces: BusinessLogicInterface)

// Package mocks is a generated GoMock package.
package mocks

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

// AddEvent mocks base method.
func (m *MockBusinessLogicInterface) AddEvent(arg0 context.Context, arg1 types.Event) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddEvent", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddEvent indicates an expected call of AddEvent.
func (mr *MockBusinessLogicInterfaceMockRecorder) AddEvent(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddEvent", reflect.TypeOf((*MockBusinessLogicInterface)(nil).AddEvent), arg0, arg1)
}

// DeleteEvent mocks base method.
func (m *MockBusinessLogicInterface) DeleteEvent(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteEvent", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteEvent indicates an expected call of DeleteEvent.
func (mr *MockBusinessLogicInterfaceMockRecorder) DeleteEvent(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteEvent", reflect.TypeOf((*MockBusinessLogicInterface)(nil).DeleteEvent), arg0, arg1)
}

// GetEvent mocks base method.
func (m *MockBusinessLogicInterface) GetEvent(arg0 context.Context, arg1 int64) (types.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEvent", arg0, arg1)
	ret0, _ := ret[0].(types.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEvent indicates an expected call of GetEvent.
func (mr *MockBusinessLogicInterfaceMockRecorder) GetEvent(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEvent", reflect.TypeOf((*MockBusinessLogicInterface)(nil).GetEvent), arg0, arg1)
}

// GetEvents mocks base method.
func (m *MockBusinessLogicInterface) GetEvents(arg0 context.Context, arg1 types.Filters) ([]types.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEvents", arg0, arg1)
	ret0, _ := ret[0].([]types.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEvents indicates an expected call of GetEvents.
func (mr *MockBusinessLogicInterfaceMockRecorder) GetEvents(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEvents", reflect.TypeOf((*MockBusinessLogicInterface)(nil).GetEvents), arg0, arg1)
}

// UpdateEvent mocks base method.
func (m *MockBusinessLogicInterface) UpdateEvent(arg0 context.Context, arg1 types.Event, arg2 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateEvent", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateEvent indicates an expected call of UpdateEvent.
func (mr *MockBusinessLogicInterfaceMockRecorder) UpdateEvent(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateEvent", reflect.TypeOf((*MockBusinessLogicInterface)(nil).UpdateEvent), arg0, arg1, arg2)
}
