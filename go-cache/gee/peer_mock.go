// Code generated by MockGen. DO NOT EDIT.
// Source: peer.go

// Package gee is a generated GoMock package.
package gee

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockPeerPicker is a mock of PeerPicker interface.
type MockPeerPicker struct {
	ctrl     *gomock.Controller
	recorder *MockPeerPickerMockRecorder
}

// MockPeerPickerMockRecorder is the mock recorder for MockPeerPicker.
type MockPeerPickerMockRecorder struct {
	mock *MockPeerPicker
}

// NewMockPeerPicker creates a new mock instance.
func NewMockPeerPicker(ctrl *gomock.Controller) *MockPeerPicker {
	mock := &MockPeerPicker{ctrl: ctrl}
	mock.recorder = &MockPeerPickerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPeerPicker) EXPECT() *MockPeerPickerMockRecorder {
	return m.recorder
}

// PickPeer mocks base method.
func (m *MockPeerPicker) PickPeer(key string) (PeerGetter, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PickPeer", key)
	ret0, _ := ret[0].(PeerGetter)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PickPeer indicates an expected call of PickPeer.
func (mr *MockPeerPickerMockRecorder) PickPeer(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PickPeer", reflect.TypeOf((*MockPeerPicker)(nil).PickPeer), key)
}

// MockPeerGetter is a mock of PeerGetter interface.
type MockPeerGetter struct {
	ctrl     *gomock.Controller
	recorder *MockPeerGetterMockRecorder
}

// MockPeerGetterMockRecorder is the mock recorder for MockPeerGetter.
type MockPeerGetterMockRecorder struct {
	mock *MockPeerGetter
}

// NewMockPeerGetter creates a new mock instance.
func NewMockPeerGetter(ctrl *gomock.Controller) *MockPeerGetter {
	mock := &MockPeerGetter{ctrl: ctrl}
	mock.recorder = &MockPeerGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPeerGetter) EXPECT() *MockPeerGetterMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockPeerGetter) Get(group, key string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", group, key)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockPeerGetterMockRecorder) Get(group, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockPeerGetter)(nil).Get), group, key)
}

// Set mocks base method.
func (m *MockPeerGetter) Set(group, key string, value []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", group, key, value)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockPeerGetterMockRecorder) Set(group, key, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockPeerGetter)(nil).Set), group, key, value)
}
