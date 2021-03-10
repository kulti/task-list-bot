// Code generated by MockGen. DO NOT EDIT.
// Source: processor.go

// Package processor_test is a generated GoMock package.
package processor_test

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	time "time"
)

// MockStore is a mock of store interface
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// CurrentSprintDump mocks base method
func (m *MockStore) CurrentSprintDump() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CurrentSprintDump")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CurrentSprintDump indicates an expected call of CurrentSprintDump
func (mr *MockStoreMockRecorder) CurrentSprintDump() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CurrentSprintDump", reflect.TypeOf((*MockStore)(nil).CurrentSprintDump))
}

// CreateNewSprint mocks base method
func (m *MockStore) CreateNewSprint(begin, end time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateNewSprint", begin, end)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateNewSprint indicates an expected call of CreateNewSprint
func (mr *MockStoreMockRecorder) CreateNewSprint(begin, end interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNewSprint", reflect.TypeOf((*MockStore)(nil).CreateNewSprint), begin, end)
}

// CreateTask mocks base method
func (m *MockStore) CreateTask(text string, points int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTask", text, points)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateTask indicates an expected call of CreateTask
func (mr *MockStoreMockRecorder) CreateTask(text, points interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTask", reflect.TypeOf((*MockStore)(nil).CreateTask), text, points)
}

// DoneTask mocks base method
func (m *MockStore) DoneTask(id int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DoneTask", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DoneTask indicates an expected call of DoneTask
func (mr *MockStoreMockRecorder) DoneTask(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DoneTask", reflect.TypeOf((*MockStore)(nil).DoneTask), id)
}

// BurnTaskPoints mocks base method
func (m *MockStore) BurnTaskPoints(id, burnt int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BurnTaskPoints", id, burnt)
	ret0, _ := ret[0].(error)
	return ret0
}

// BurnTaskPoints indicates an expected call of BurnTaskPoints
func (mr *MockStoreMockRecorder) BurnTaskPoints(id, burnt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BurnTaskPoints", reflect.TypeOf((*MockStore)(nil).BurnTaskPoints), id, burnt)
}
