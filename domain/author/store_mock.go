// Code generated by MockGen. DO NOT EDIT.
// Source: ./domain/author/types.go

// Package author is a generated GoMock package.
package author

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

// ReplaceBookAuthors mocks base method.
func (m *MockStore) ReplaceBookAuthors(ctx context.Context, bookID int64, authorIDs []int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReplaceBookAuthors", ctx, bookID, authorIDs)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReplaceBookAuthors indicates an expected call of ReplaceBookAuthors.
func (mr *MockStoreMockRecorder) ReplaceBookAuthors(ctx, bookID, authorIDs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReplaceBookAuthors", reflect.TypeOf((*MockStore)(nil).ReplaceBookAuthors), ctx, bookID, authorIDs)
}

// UpsertAll mocks base method.
func (m *MockStore) UpsertAll(ctx context.Context, authors []string) ([]int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertAll", ctx, authors)
	ret0, _ := ret[0].([]int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpsertAll indicates an expected call of UpsertAll.
func (mr *MockStoreMockRecorder) UpsertAll(ctx, authors interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertAll", reflect.TypeOf((*MockStore)(nil).UpsertAll), ctx, authors)
}
