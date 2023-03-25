// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/auth/auth.go

// Package auth is a generated GoMock package.
package auth

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockAuthenticator is a mock of Authenticator interface.
type MockAuthenticator struct {
	ctrl     *gomock.Controller
	recorder *MockAuthenticatorMockRecorder
}

// MockAuthenticatorMockRecorder is the mock recorder for MockAuthenticator.
type MockAuthenticatorMockRecorder struct {
	mock *MockAuthenticator
}

// NewMockAuthenticator creates a new mock instance.
func NewMockAuthenticator(ctrl *gomock.Controller) *MockAuthenticator {
	mock := &MockAuthenticator{ctrl: ctrl}
	mock.recorder = &MockAuthenticatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthenticator) EXPECT() *MockAuthenticatorMockRecorder {
	return m.recorder
}

// AuthenticateToken mocks base method.
func (m *MockAuthenticator) AuthenticateToken(arg0 context.Context, arg1 string) (*Claims, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AuthenticateToken", arg0, arg1)
	ret0, _ := ret[0].(*Claims)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AuthenticateToken indicates an expected call of AuthenticateToken.
func (mr *MockAuthenticatorMockRecorder) AuthenticateToken(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthenticateToken", reflect.TypeOf((*MockAuthenticator)(nil).AuthenticateToken), arg0, arg1)
}