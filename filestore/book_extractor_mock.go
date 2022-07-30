// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/sdreger/lib-file-processor-go/filestore (interfaces: BookExtractor)

// Package filestore is a generated GoMock package.
package filestore

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockBookExtractor is a mock of BookExtractor interface.
type MockBookExtractor struct {
	ctrl     *gomock.Controller
	recorder *MockBookExtractorMockRecorder
}

// MockBookExtractorMockRecorder is the mock recorder for MockBookExtractor.
type MockBookExtractorMockRecorder struct {
	mock *MockBookExtractor
}

// NewMockBookExtractor creates a new mock instance.
func NewMockBookExtractor(ctrl *gomock.Controller) *MockBookExtractor {
	mock := &MockBookExtractor{ctrl: ctrl}
	mock.recorder = &MockBookExtractorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBookExtractor) EXPECT() *MockBookExtractorMockRecorder {
	return m.recorder
}

// ExtractZipFile mocks base method.
func (m *MockBookExtractor) ExtractZipFile(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExtractZipFile", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ExtractZipFile indicates an expected call of ExtractZipFile.
func (mr *MockBookExtractorMockRecorder) ExtractZipFile(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExtractZipFile", reflect.TypeOf((*MockBookExtractor)(nil).ExtractZipFile), arg0, arg1)
}