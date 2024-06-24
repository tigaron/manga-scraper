// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/providers.go
//
// Generated by this command:
//
//	mockgen -package mock -destination internal/service/mock/providers.go -source internal/service/providers.go ProviderRepository
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	internal "fourleaves.studio/manga-scraper/internal"
	gomock "go.uber.org/mock/gomock"
)

// MockProviderRepository is a mock of ProviderRepository interface.
type MockProviderRepository struct {
	ctrl     *gomock.Controller
	recorder *MockProviderRepositoryMockRecorder
}

// MockProviderRepositoryMockRecorder is the mock recorder for MockProviderRepository.
type MockProviderRepositoryMockRecorder struct {
	mock *MockProviderRepository
}

// NewMockProviderRepository creates a new mock instance.
func NewMockProviderRepository(ctrl *gomock.Controller) *MockProviderRepository {
	mock := &MockProviderRepository{ctrl: ctrl}
	mock.recorder = &MockProviderRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProviderRepository) EXPECT() *MockProviderRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockProviderRepository) Create(ctx context.Context, params internal.ProviderParams) (internal.Provider, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, params)
	ret0, _ := ret[0].(internal.Provider)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockProviderRepositoryMockRecorder) Create(ctx, params any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockProviderRepository)(nil).Create), ctx, params)
}

// Delete mocks base method.
func (m *MockProviderRepository) Delete(ctx context.Context, slug string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, slug)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockProviderRepositoryMockRecorder) Delete(ctx, slug any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockProviderRepository)(nil).Delete), ctx, slug)
}

// Find mocks base method.
func (m *MockProviderRepository) Find(ctx context.Context, slug string) (internal.Provider, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", ctx, slug)
	ret0, _ := ret[0].(internal.Provider)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find.
func (mr *MockProviderRepositoryMockRecorder) Find(ctx, slug any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockProviderRepository)(nil).Find), ctx, slug)
}

// FindAll mocks base method.
func (m *MockProviderRepository) FindAll(ctx context.Context, order internal.SortOrder) ([]internal.Provider, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll", ctx, order)
	ret0, _ := ret[0].([]internal.Provider)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAll indicates an expected call of FindAll.
func (mr *MockProviderRepositoryMockRecorder) FindAll(ctx, order any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockProviderRepository)(nil).FindAll), ctx, order)
}

// FindBC mocks base method.
func (m *MockProviderRepository) FindBC(ctx context.Context, slug string) (internal.ProviderBC, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindBC", ctx, slug)
	ret0, _ := ret[0].(internal.ProviderBC)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindBC indicates an expected call of FindBC.
func (mr *MockProviderRepositoryMockRecorder) FindBC(ctx, slug any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindBC", reflect.TypeOf((*MockProviderRepository)(nil).FindBC), ctx, slug)
}

// Update mocks base method.
func (m *MockProviderRepository) Update(ctx context.Context, params internal.ProviderParams) (internal.Provider, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, params)
	ret0, _ := ret[0].(internal.Provider)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockProviderRepositoryMockRecorder) Update(ctx, params any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockProviderRepository)(nil).Update), ctx, params)
}
