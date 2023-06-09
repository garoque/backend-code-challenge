// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/database/transaction/transaction.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	entity "github.com/garoque/backend-code-challenge-snapfi/internal/entity"
	gomock "github.com/golang/mock/gomock"
)

// MockDabataseTransactionInterface is a mock of DabataseTransactionInterface interface.
type MockDabataseTransactionInterface struct {
	ctrl     *gomock.Controller
	recorder *MockDabataseTransactionInterfaceMockRecorder
}

// MockDabataseTransactionInterfaceMockRecorder is the mock recorder for MockDabataseTransactionInterface.
type MockDabataseTransactionInterfaceMockRecorder struct {
	mock *MockDabataseTransactionInterface
}

// NewMockDabataseTransactionInterface creates a new mock instance.
func NewMockDabataseTransactionInterface(ctrl *gomock.Controller) *MockDabataseTransactionInterface {
	mock := &MockDabataseTransactionInterface{ctrl: ctrl}
	mock.recorder = &MockDabataseTransactionInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDabataseTransactionInterface) EXPECT() *MockDabataseTransactionInterfaceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockDabataseTransactionInterface) Create(ctx context.Context, transaction *entity.Transaction) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, transaction)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockDabataseTransactionInterfaceMockRecorder) Create(ctx, transaction interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockDabataseTransactionInterface)(nil).Create), ctx, transaction)
}

// ReadAll mocks base method.
func (m *MockDabataseTransactionInterface) ReadAll(ctx context.Context) ([]entity.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadAll", ctx)
	ret0, _ := ret[0].([]entity.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadAll indicates an expected call of ReadAll.
func (mr *MockDabataseTransactionInterfaceMockRecorder) ReadAll(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadAll", reflect.TypeOf((*MockDabataseTransactionInterface)(nil).ReadAll), ctx)
}

// ReadBalance mocks base method.
func (m *MockDabataseTransactionInterface) ReadBalance(ctx context.Context, userId string) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadBalance", ctx, userId)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadBalance indicates an expected call of ReadBalance.
func (mr *MockDabataseTransactionInterfaceMockRecorder) ReadBalance(ctx, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadBalance", reflect.TypeOf((*MockDabataseTransactionInterface)(nil).ReadBalance), ctx, userId)
}

// UpdateBalanceUser mocks base method.
func (m *MockDabataseTransactionInterface) UpdateBalanceUser(ctx context.Context, userId string, value float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBalanceUser", ctx, userId, value)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateBalanceUser indicates an expected call of UpdateBalanceUser.
func (mr *MockDabataseTransactionInterfaceMockRecorder) UpdateBalanceUser(ctx, userId, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBalanceUser", reflect.TypeOf((*MockDabataseTransactionInterface)(nil).UpdateBalanceUser), ctx, userId, value)
}

// UpdateState mocks base method.
func (m *MockDabataseTransactionInterface) UpdateState(ctx context.Context, state entity.StatesTransaction, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateState", ctx, state, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateState indicates an expected call of UpdateState.
func (mr *MockDabataseTransactionInterfaceMockRecorder) UpdateState(ctx, state, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateState", reflect.TypeOf((*MockDabataseTransactionInterface)(nil).UpdateState), ctx, state, id)
}
