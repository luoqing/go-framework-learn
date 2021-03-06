// Code generated by MockGen. DO NOT EDIT.
// Source: sink.go

// Package gee is a generated GoMock package.
package gee

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockSinker is a mock of Sinker interface.
type MockSinker struct {
	ctrl     *gomock.Controller
	recorder *MockSinkerMockRecorder
}

// MockSinkerMockRecorder is the mock recorder for MockSinker.
type MockSinkerMockRecorder struct {
	mock *MockSinker
}

// NewMockSinker creates a new mock instance.
func NewMockSinker(ctrl *gomock.Controller) *MockSinker {
	mock := &MockSinker{ctrl: ctrl}
	mock.recorder = &MockSinkerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSinker) EXPECT() *MockSinkerMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockSinker) Get(key string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", key)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockSinkerMockRecorder) Get(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockSinker)(nil).Get), key)
}

// Set mocks base method.
func (m *MockSinker) Set(key string, value []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", key, value)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockSinkerMockRecorder) Set(key, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockSinker)(nil).Set), key, value)
}
