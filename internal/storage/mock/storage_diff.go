// Code generated by MockGen. DO NOT EDIT.
// Source: storage_diff.go
//
// Generated by this command:
//
//	mockgen -source=storage_diff.go -destination=mock/storage_diff.go -package=mock -typed
//
// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	storage "github.com/dipdup-io/starknet-indexer/internal/storage"
	storage0 "github.com/dipdup-net/indexer-sdk/pkg/storage"
	gomock "go.uber.org/mock/gomock"
)

// MockIStorageDiff is a mock of IStorageDiff interface.
type MockIStorageDiff struct {
	ctrl     *gomock.Controller
	recorder *MockIStorageDiffMockRecorder
}

// MockIStorageDiffMockRecorder is the mock recorder for MockIStorageDiff.
type MockIStorageDiffMockRecorder struct {
	mock *MockIStorageDiff
}

// NewMockIStorageDiff creates a new mock instance.
func NewMockIStorageDiff(ctrl *gomock.Controller) *MockIStorageDiff {
	mock := &MockIStorageDiff{ctrl: ctrl}
	mock.recorder = &MockIStorageDiffMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIStorageDiff) EXPECT() *MockIStorageDiffMockRecorder {
	return m.recorder
}

// CursorList mocks base method.
func (m *MockIStorageDiff) CursorList(ctx context.Context, id, limit uint64, order storage0.SortOrder, cmp storage0.Comparator) ([]*storage.StorageDiff, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CursorList", ctx, id, limit, order, cmp)
	ret0, _ := ret[0].([]*storage.StorageDiff)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CursorList indicates an expected call of CursorList.
func (mr *MockIStorageDiffMockRecorder) CursorList(ctx, id, limit, order, cmp any) *IStorageDiffCursorListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CursorList", reflect.TypeOf((*MockIStorageDiff)(nil).CursorList), ctx, id, limit, order, cmp)
	return &IStorageDiffCursorListCall{Call: call}
}

// IStorageDiffCursorListCall wrap *gomock.Call
type IStorageDiffCursorListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IStorageDiffCursorListCall) Return(arg0 []*storage.StorageDiff, arg1 error) *IStorageDiffCursorListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IStorageDiffCursorListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.StorageDiff, error)) *IStorageDiffCursorListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IStorageDiffCursorListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.StorageDiff, error)) *IStorageDiffCursorListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Filter mocks base method.
func (m *MockIStorageDiff) Filter(ctx context.Context, flt []storage.StorageDiffFilter, opts ...storage.FilterOption) ([]storage.StorageDiff, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, flt}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Filter", varargs...)
	ret0, _ := ret[0].([]storage.StorageDiff)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Filter indicates an expected call of Filter.
func (mr *MockIStorageDiffMockRecorder) Filter(ctx, flt any, opts ...any) *IStorageDiffFilterCall {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, flt}, opts...)
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Filter", reflect.TypeOf((*MockIStorageDiff)(nil).Filter), varargs...)
	return &IStorageDiffFilterCall{Call: call}
}

// IStorageDiffFilterCall wrap *gomock.Call
type IStorageDiffFilterCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IStorageDiffFilterCall) Return(arg0 []storage.StorageDiff, arg1 error) *IStorageDiffFilterCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IStorageDiffFilterCall) Do(f func(context.Context, []storage.StorageDiffFilter, ...storage.FilterOption) ([]storage.StorageDiff, error)) *IStorageDiffFilterCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IStorageDiffFilterCall) DoAndReturn(f func(context.Context, []storage.StorageDiffFilter, ...storage.FilterOption) ([]storage.StorageDiff, error)) *IStorageDiffFilterCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetByID mocks base method.
func (m *MockIStorageDiff) GetByID(ctx context.Context, id uint64) (*storage.StorageDiff, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*storage.StorageDiff)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockIStorageDiffMockRecorder) GetByID(ctx, id any) *IStorageDiffGetByIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockIStorageDiff)(nil).GetByID), ctx, id)
	return &IStorageDiffGetByIDCall{Call: call}
}

// IStorageDiffGetByIDCall wrap *gomock.Call
type IStorageDiffGetByIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IStorageDiffGetByIDCall) Return(arg0 *storage.StorageDiff, arg1 error) *IStorageDiffGetByIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IStorageDiffGetByIDCall) Do(f func(context.Context, uint64) (*storage.StorageDiff, error)) *IStorageDiffGetByIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IStorageDiffGetByIDCall) DoAndReturn(f func(context.Context, uint64) (*storage.StorageDiff, error)) *IStorageDiffGetByIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetOnBlock mocks base method.
func (m *MockIStorageDiff) GetOnBlock(ctx context.Context, height, contractId uint64, key []byte) (storage.StorageDiff, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOnBlock", ctx, height, contractId, key)
	ret0, _ := ret[0].(storage.StorageDiff)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOnBlock indicates an expected call of GetOnBlock.
func (mr *MockIStorageDiffMockRecorder) GetOnBlock(ctx, height, contractId, key any) *IStorageDiffGetOnBlockCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOnBlock", reflect.TypeOf((*MockIStorageDiff)(nil).GetOnBlock), ctx, height, contractId, key)
	return &IStorageDiffGetOnBlockCall{Call: call}
}

// IStorageDiffGetOnBlockCall wrap *gomock.Call
type IStorageDiffGetOnBlockCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IStorageDiffGetOnBlockCall) Return(arg0 storage.StorageDiff, arg1 error) *IStorageDiffGetOnBlockCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IStorageDiffGetOnBlockCall) Do(f func(context.Context, uint64, uint64, []byte) (storage.StorageDiff, error)) *IStorageDiffGetOnBlockCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IStorageDiffGetOnBlockCall) DoAndReturn(f func(context.Context, uint64, uint64, []byte) (storage.StorageDiff, error)) *IStorageDiffGetOnBlockCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// IsNoRows mocks base method.
func (m *MockIStorageDiff) IsNoRows(err error) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsNoRows", err)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsNoRows indicates an expected call of IsNoRows.
func (mr *MockIStorageDiffMockRecorder) IsNoRows(err any) *IStorageDiffIsNoRowsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsNoRows", reflect.TypeOf((*MockIStorageDiff)(nil).IsNoRows), err)
	return &IStorageDiffIsNoRowsCall{Call: call}
}

// IStorageDiffIsNoRowsCall wrap *gomock.Call
type IStorageDiffIsNoRowsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IStorageDiffIsNoRowsCall) Return(arg0 bool) *IStorageDiffIsNoRowsCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IStorageDiffIsNoRowsCall) Do(f func(error) bool) *IStorageDiffIsNoRowsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IStorageDiffIsNoRowsCall) DoAndReturn(f func(error) bool) *IStorageDiffIsNoRowsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// LastID mocks base method.
func (m *MockIStorageDiff) LastID(ctx context.Context) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LastID", ctx)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LastID indicates an expected call of LastID.
func (mr *MockIStorageDiffMockRecorder) LastID(ctx any) *IStorageDiffLastIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastID", reflect.TypeOf((*MockIStorageDiff)(nil).LastID), ctx)
	return &IStorageDiffLastIDCall{Call: call}
}

// IStorageDiffLastIDCall wrap *gomock.Call
type IStorageDiffLastIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IStorageDiffLastIDCall) Return(arg0 uint64, arg1 error) *IStorageDiffLastIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IStorageDiffLastIDCall) Do(f func(context.Context) (uint64, error)) *IStorageDiffLastIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IStorageDiffLastIDCall) DoAndReturn(f func(context.Context) (uint64, error)) *IStorageDiffLastIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// List mocks base method.
func (m *MockIStorageDiff) List(ctx context.Context, limit, offset uint64, order storage0.SortOrder) ([]*storage.StorageDiff, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, limit, offset, order)
	ret0, _ := ret[0].([]*storage.StorageDiff)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockIStorageDiffMockRecorder) List(ctx, limit, offset, order any) *IStorageDiffListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockIStorageDiff)(nil).List), ctx, limit, offset, order)
	return &IStorageDiffListCall{Call: call}
}

// IStorageDiffListCall wrap *gomock.Call
type IStorageDiffListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IStorageDiffListCall) Return(arg0 []*storage.StorageDiff, arg1 error) *IStorageDiffListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IStorageDiffListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.StorageDiff, error)) *IStorageDiffListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IStorageDiffListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.StorageDiff, error)) *IStorageDiffListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Save mocks base method.
func (m_2 *MockIStorageDiff) Save(ctx context.Context, m *storage.StorageDiff) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Save", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockIStorageDiffMockRecorder) Save(ctx, m any) *IStorageDiffSaveCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockIStorageDiff)(nil).Save), ctx, m)
	return &IStorageDiffSaveCall{Call: call}
}

// IStorageDiffSaveCall wrap *gomock.Call
type IStorageDiffSaveCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IStorageDiffSaveCall) Return(arg0 error) *IStorageDiffSaveCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IStorageDiffSaveCall) Do(f func(context.Context, *storage.StorageDiff) error) *IStorageDiffSaveCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IStorageDiffSaveCall) DoAndReturn(f func(context.Context, *storage.StorageDiff) error) *IStorageDiffSaveCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Update mocks base method.
func (m_2 *MockIStorageDiff) Update(ctx context.Context, m *storage.StorageDiff) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Update", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockIStorageDiffMockRecorder) Update(ctx, m any) *IStorageDiffUpdateCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockIStorageDiff)(nil).Update), ctx, m)
	return &IStorageDiffUpdateCall{Call: call}
}

// IStorageDiffUpdateCall wrap *gomock.Call
type IStorageDiffUpdateCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IStorageDiffUpdateCall) Return(arg0 error) *IStorageDiffUpdateCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IStorageDiffUpdateCall) Do(f func(context.Context, *storage.StorageDiff) error) *IStorageDiffUpdateCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IStorageDiffUpdateCall) DoAndReturn(f func(context.Context, *storage.StorageDiff) error) *IStorageDiffUpdateCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
