// Code generated by MockGen. DO NOT EDIT.
// Source: deploy_account.go
//
// Generated by this command:
//
//	mockgen -source=deploy_account.go -destination=mock/deploy_account.go -package=mock -typed
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

// MockIDeployAccount is a mock of IDeployAccount interface.
type MockIDeployAccount struct {
	ctrl     *gomock.Controller
	recorder *MockIDeployAccountMockRecorder
}

// MockIDeployAccountMockRecorder is the mock recorder for MockIDeployAccount.
type MockIDeployAccountMockRecorder struct {
	mock *MockIDeployAccount
}

// NewMockIDeployAccount creates a new mock instance.
func NewMockIDeployAccount(ctrl *gomock.Controller) *MockIDeployAccount {
	mock := &MockIDeployAccount{ctrl: ctrl}
	mock.recorder = &MockIDeployAccountMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIDeployAccount) EXPECT() *MockIDeployAccountMockRecorder {
	return m.recorder
}

// CursorList mocks base method.
func (m *MockIDeployAccount) CursorList(ctx context.Context, id, limit uint64, order storage0.SortOrder, cmp storage0.Comparator) ([]*storage.DeployAccount, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CursorList", ctx, id, limit, order, cmp)
	ret0, _ := ret[0].([]*storage.DeployAccount)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CursorList indicates an expected call of CursorList.
func (mr *MockIDeployAccountMockRecorder) CursorList(ctx, id, limit, order, cmp any) *IDeployAccountCursorListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CursorList", reflect.TypeOf((*MockIDeployAccount)(nil).CursorList), ctx, id, limit, order, cmp)
	return &IDeployAccountCursorListCall{Call: call}
}

// IDeployAccountCursorListCall wrap *gomock.Call
type IDeployAccountCursorListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IDeployAccountCursorListCall) Return(arg0 []*storage.DeployAccount, arg1 error) *IDeployAccountCursorListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IDeployAccountCursorListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.DeployAccount, error)) *IDeployAccountCursorListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IDeployAccountCursorListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.DeployAccount, error)) *IDeployAccountCursorListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Filter mocks base method.
func (m *MockIDeployAccount) Filter(ctx context.Context, flt []storage.DeployAccountFilter, opts ...storage.FilterOption) ([]storage.DeployAccount, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, flt}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Filter", varargs...)
	ret0, _ := ret[0].([]storage.DeployAccount)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Filter indicates an expected call of Filter.
func (mr *MockIDeployAccountMockRecorder) Filter(ctx, flt any, opts ...any) *IDeployAccountFilterCall {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, flt}, opts...)
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Filter", reflect.TypeOf((*MockIDeployAccount)(nil).Filter), varargs...)
	return &IDeployAccountFilterCall{Call: call}
}

// IDeployAccountFilterCall wrap *gomock.Call
type IDeployAccountFilterCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IDeployAccountFilterCall) Return(arg0 []storage.DeployAccount, arg1 error) *IDeployAccountFilterCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IDeployAccountFilterCall) Do(f func(context.Context, []storage.DeployAccountFilter, ...storage.FilterOption) ([]storage.DeployAccount, error)) *IDeployAccountFilterCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IDeployAccountFilterCall) DoAndReturn(f func(context.Context, []storage.DeployAccountFilter, ...storage.FilterOption) ([]storage.DeployAccount, error)) *IDeployAccountFilterCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetByID mocks base method.
func (m *MockIDeployAccount) GetByID(ctx context.Context, id uint64) (*storage.DeployAccount, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*storage.DeployAccount)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockIDeployAccountMockRecorder) GetByID(ctx, id any) *IDeployAccountGetByIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockIDeployAccount)(nil).GetByID), ctx, id)
	return &IDeployAccountGetByIDCall{Call: call}
}

// IDeployAccountGetByIDCall wrap *gomock.Call
type IDeployAccountGetByIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IDeployAccountGetByIDCall) Return(arg0 *storage.DeployAccount, arg1 error) *IDeployAccountGetByIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IDeployAccountGetByIDCall) Do(f func(context.Context, uint64) (*storage.DeployAccount, error)) *IDeployAccountGetByIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IDeployAccountGetByIDCall) DoAndReturn(f func(context.Context, uint64) (*storage.DeployAccount, error)) *IDeployAccountGetByIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// IsNoRows mocks base method.
func (m *MockIDeployAccount) IsNoRows(err error) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsNoRows", err)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsNoRows indicates an expected call of IsNoRows.
func (mr *MockIDeployAccountMockRecorder) IsNoRows(err any) *IDeployAccountIsNoRowsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsNoRows", reflect.TypeOf((*MockIDeployAccount)(nil).IsNoRows), err)
	return &IDeployAccountIsNoRowsCall{Call: call}
}

// IDeployAccountIsNoRowsCall wrap *gomock.Call
type IDeployAccountIsNoRowsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IDeployAccountIsNoRowsCall) Return(arg0 bool) *IDeployAccountIsNoRowsCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IDeployAccountIsNoRowsCall) Do(f func(error) bool) *IDeployAccountIsNoRowsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IDeployAccountIsNoRowsCall) DoAndReturn(f func(error) bool) *IDeployAccountIsNoRowsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// LastID mocks base method.
func (m *MockIDeployAccount) LastID(ctx context.Context) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LastID", ctx)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LastID indicates an expected call of LastID.
func (mr *MockIDeployAccountMockRecorder) LastID(ctx any) *IDeployAccountLastIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastID", reflect.TypeOf((*MockIDeployAccount)(nil).LastID), ctx)
	return &IDeployAccountLastIDCall{Call: call}
}

// IDeployAccountLastIDCall wrap *gomock.Call
type IDeployAccountLastIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IDeployAccountLastIDCall) Return(arg0 uint64, arg1 error) *IDeployAccountLastIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IDeployAccountLastIDCall) Do(f func(context.Context) (uint64, error)) *IDeployAccountLastIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IDeployAccountLastIDCall) DoAndReturn(f func(context.Context) (uint64, error)) *IDeployAccountLastIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// List mocks base method.
func (m *MockIDeployAccount) List(ctx context.Context, limit, offset uint64, order storage0.SortOrder) ([]*storage.DeployAccount, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, limit, offset, order)
	ret0, _ := ret[0].([]*storage.DeployAccount)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockIDeployAccountMockRecorder) List(ctx, limit, offset, order any) *IDeployAccountListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockIDeployAccount)(nil).List), ctx, limit, offset, order)
	return &IDeployAccountListCall{Call: call}
}

// IDeployAccountListCall wrap *gomock.Call
type IDeployAccountListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IDeployAccountListCall) Return(arg0 []*storage.DeployAccount, arg1 error) *IDeployAccountListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IDeployAccountListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.DeployAccount, error)) *IDeployAccountListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IDeployAccountListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.DeployAccount, error)) *IDeployAccountListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Save mocks base method.
func (m_2 *MockIDeployAccount) Save(ctx context.Context, m *storage.DeployAccount) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Save", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockIDeployAccountMockRecorder) Save(ctx, m any) *IDeployAccountSaveCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockIDeployAccount)(nil).Save), ctx, m)
	return &IDeployAccountSaveCall{Call: call}
}

// IDeployAccountSaveCall wrap *gomock.Call
type IDeployAccountSaveCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IDeployAccountSaveCall) Return(arg0 error) *IDeployAccountSaveCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IDeployAccountSaveCall) Do(f func(context.Context, *storage.DeployAccount) error) *IDeployAccountSaveCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IDeployAccountSaveCall) DoAndReturn(f func(context.Context, *storage.DeployAccount) error) *IDeployAccountSaveCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Update mocks base method.
func (m_2 *MockIDeployAccount) Update(ctx context.Context, m *storage.DeployAccount) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Update", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockIDeployAccountMockRecorder) Update(ctx, m any) *IDeployAccountUpdateCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockIDeployAccount)(nil).Update), ctx, m)
	return &IDeployAccountUpdateCall{Call: call}
}

// IDeployAccountUpdateCall wrap *gomock.Call
type IDeployAccountUpdateCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IDeployAccountUpdateCall) Return(arg0 error) *IDeployAccountUpdateCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IDeployAccountUpdateCall) Do(f func(context.Context, *storage.DeployAccount) error) *IDeployAccountUpdateCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IDeployAccountUpdateCall) DoAndReturn(f func(context.Context, *storage.DeployAccount) error) *IDeployAccountUpdateCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
