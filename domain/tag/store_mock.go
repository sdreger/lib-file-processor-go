// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/sdreger/lib-file-processor-go/domain/tag (interfaces: Store)

// Package tag is a generated GoMock package.
package tag

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

// ReplaceBookTags mocks base method.
func (m *MockStore) ReplaceBookTags(arg0 context.Context, arg1 int64, arg2 []int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReplaceBookTags", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReplaceBookTags indicates an expected call of ReplaceBookTags.
func (mr *MockStoreMockRecorder) ReplaceBookTags(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReplaceBookTags", reflect.TypeOf((*MockStore)(nil).ReplaceBookTags), arg0, arg1, arg2)
}

// UpsertAll mocks base method.
func (m *MockStore) UpsertAll(arg0 context.Context, arg1 []string) ([]int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertAll", arg0, arg1)
	ret0, _ := ret[0].([]int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpsertAll indicates an expected call of UpsertAll.
func (mr *MockStoreMockRecorder) UpsertAll(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertAll", reflect.TypeOf((*MockStore)(nil).UpsertAll), arg0, arg1)
}
