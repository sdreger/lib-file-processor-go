// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/sdreger/lib-file-processor-go/domain/lang (interfaces: Store)

// Package lang is a generated GoMock package.
package lang

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// Upsert mocks base method.
func (m *MockStore) Upsert(arg0 context.Context, arg1 string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Upsert", arg0, arg1)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Upsert indicates an expected call of Upsert.
func (mr *MockStoreMockRecorder) Upsert(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Upsert", reflect.TypeOf((*MockStore)(nil).Upsert), arg0, arg1)
}
