// Code generated by MockGen. DO NOT EDIT.
// Source: cmd/auth/keyStore.go

// Package auth is a generated GoMock package.
package auth

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockJwkStore is a mock of JwkStore interface.
type MockJwkStore struct {
	ctrl     *gomock.Controller
	recorder *MockJwkStoreMockRecorder
}

// MockJwkStoreMockRecorder is the mock recorder for MockJwkStore.
type MockJwkStoreMockRecorder struct {
	mock *MockJwkStore
}

// NewMockJwkStore creates a new mock instance.
func NewMockJwkStore(ctrl *gomock.Controller) *MockJwkStore {
	mock := &MockJwkStore{ctrl: ctrl}
	mock.recorder = &MockJwkStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockJwkStore) EXPECT() *MockJwkStoreMockRecorder {
	return m.recorder
}

// GetJWK mocks base method.
func (m *MockJwkStore) GetJWK(kid, iss string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetJWK", kid, iss)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetJWK indicates an expected call of GetJWK.
func (mr *MockJwkStoreMockRecorder) GetJWK(kid, iss interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetJWK", reflect.TypeOf((*MockJwkStore)(nil).GetJWK), kid, iss)
}