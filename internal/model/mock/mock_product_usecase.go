// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/irvankadhafi/erajaya-product-service/internal/model (interfaces: ProductUsecase)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/irvankadhafi/erajaya-product-service/internal/model"
)

// MockProductUsecase is a mock of ProductUsecase interface.
type MockProductUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockProductUsecaseMockRecorder
}

// MockProductUsecaseMockRecorder is the mock recorder for MockProductUsecase.
type MockProductUsecaseMockRecorder struct {
	mock *MockProductUsecase
}

// NewMockProductUsecase creates a new mock instance.
func NewMockProductUsecase(ctrl *gomock.Controller) *MockProductUsecase {
	mock := &MockProductUsecase{ctrl: ctrl}
	mock.recorder = &MockProductUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProductUsecase) EXPECT() *MockProductUsecaseMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockProductUsecase) Create(arg0 context.Context, arg1 model.CreateProductInput) (*model.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(*model.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockProductUsecaseMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockProductUsecase)(nil).Create), arg0, arg1)
}

// FindAllByIDs mocks base method.
func (m *MockProductUsecase) FindAllByIDs(arg0 context.Context, arg1 []int64) []*model.Product {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllByIDs", arg0, arg1)
	ret0, _ := ret[0].([]*model.Product)
	return ret0
}

// FindAllByIDs indicates an expected call of FindAllByIDs.
func (mr *MockProductUsecaseMockRecorder) FindAllByIDs(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllByIDs", reflect.TypeOf((*MockProductUsecase)(nil).FindAllByIDs), arg0, arg1)
}

// FindByID mocks base method.
func (m *MockProductUsecase) FindByID(arg0 context.Context, arg1 int64) (*model.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", arg0, arg1)
	ret0, _ := ret[0].(*model.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockProductUsecaseMockRecorder) FindByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockProductUsecase)(nil).FindByID), arg0, arg1)
}

// Search mocks base method.
func (m *MockProductUsecase) Search(arg0 context.Context, arg1 model.ProductSearchCriteria) ([]*model.Product, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Search", arg0, arg1)
	ret0, _ := ret[0].([]*model.Product)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Search indicates an expected call of Search.
func (mr *MockProductUsecaseMockRecorder) Search(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Search", reflect.TypeOf((*MockProductUsecase)(nil).Search), arg0, arg1)
}
