// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/alphaonly/gomartv2/internal/domain/user (interfaces: Storage)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	user "github.com/alphaonly/gomartv2/internal/domain/user"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockStorage is a mock of Storage interface
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// GetUser mocks base method
func (m *MockStorage) GetUser(arg0 context.Context, arg1 string) (*user.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", arg0, arg1)
	ret0, _ := ret[0].(*user.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser
func (mr *MockStorageMockRecorder) GetUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockStorage)(nil).GetUser), arg0, arg1)
}

// SaveUser mocks base method
func (m *MockStorage) SaveUser(arg0 context.Context, arg1 *user.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveUser", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveUser indicates an expected call of SaveUser
func (mr *MockStorageMockRecorder) SaveUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveUser", reflect.TypeOf((*MockStorage)(nil).SaveUser), arg0, arg1)
}