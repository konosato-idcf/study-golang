// Code generated by MockGen. DO NOT EDIT.
// Source: ./echo-OJT-webApp/user/models.go

// Package mock_user is a generated GoMock package.
package user

import (
	"context"
	"github.com/golang/mock/gomock"
	models2 "github.com/konosato-idcf/study-golang/echo-OJT-webApp/app/user/infra/models"
	"reflect"
)

// MockUsersInterface is a mock of UsersInterface interface
type MockUsersInterface struct {
	ctrl     *gomock.Controller
	recorder *MockUsersInterfaceMockRecorder
}

// MockUsersInterfaceMockRecorder is the mock recorder for MockUsersInterface
type MockUsersInterfaceMockRecorder struct {
	mock *MockUsersInterface
}

// NewMockUsersInterface creates a new mock instance
func NewMockUsersInterface(ctrl *gomock.Controller) *MockUsersInterface {
	mock := &MockUsersInterface{ctrl: ctrl}
	mock.recorder = &MockUsersInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUsersInterface) EXPECT() *MockUsersInterfaceMockRecorder {
	return m.recorder
}

// All mocks base method
func (m *MockUsersInterface) All(arg0 context.Context) (models2.UserSlice, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "All", arg0)
	ret0, _ := ret[0].(models2.UserSlice)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// All indicates an expected call of All
func (mr *MockUsersInterfaceMockRecorder) All(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "All", reflect.TypeOf((*MockUsersInterface)(nil).All), arg0)
}

// Create mocks base method
func (m *MockUsersInterface) Create(arg0 context.Context, arg1 *User) (*User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(*User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockUsersInterfaceMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUsersInterface)(nil).Create), arg0, arg1)
}

// FindById mocks base method
func (m *MockUsersInterface) FindById(arg0 context.Context, arg1 int) (*User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindById", arg0, arg1)
	ret0, _ := ret[0].(*User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindById indicates an expected call of FindById
func (mr *MockUsersInterfaceMockRecorder) FindById(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindById", reflect.TypeOf((*MockUsersInterface)(nil).FindById), arg0, arg1)
}

// Update mocks base method
func (m *MockUsersInterface) Update(arg0 context.Context, arg1 *User) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update
func (mr *MockUsersInterfaceMockRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockUsersInterface)(nil).Update), arg0, arg1)
}

// Delete mocks base method
func (m *MockUsersInterface) Delete(arg0 context.Context, arg1 *User) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete
func (mr *MockUsersInterfaceMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockUsersInterface)(nil).Delete), arg0, arg1)
}
