// Code generated by MockGen. DO NOT EDIT.
// Source: cmd/handler/pois.go

// Package handler is a generated GoMock package.
package handler

import (
	data "poi-service/cmd/data"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockPoiHandler is a mock of PoiHandler interface.
type MockPoiHandler struct {
	ctrl     *gomock.Controller
	recorder *MockPoiHandlerMockRecorder
}

// MockPoiHandlerMockRecorder is the mock recorder for MockPoiHandler.
type MockPoiHandlerMockRecorder struct {
	mock *MockPoiHandler
}

// NewMockPoiHandler creates a new mock instance.
func NewMockPoiHandler(ctrl *gomock.Controller) *MockPoiHandler {
	mock := &MockPoiHandler{ctrl: ctrl}
	mock.recorder = &MockPoiHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPoiHandler) EXPECT() *MockPoiHandlerMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockPoiHandler) Create(poi *data.Poi) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", poi)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockPoiHandlerMockRecorder) Create(poi interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockPoiHandler)(nil).Create), poi)
}

// Delete mocks base method.
func (m *MockPoiHandler) Delete(id data.Id) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockPoiHandlerMockRecorder) Delete(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockPoiHandler)(nil).Delete), id)
}

// Get mocks base method.
func (m *MockPoiHandler) Get(id data.Id) (data.Poi, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", id)
	ret0, _ := ret[0].(data.Poi)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockPoiHandlerMockRecorder) Get(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockPoiHandler)(nil).Get), id)
}

// Search mocks base method.
func (m *MockPoiHandler) Search(pos data.SearchArea) (data.Pois, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Search", pos)
	ret0, _ := ret[0].(data.Pois)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Search indicates an expected call of Search.
func (mr *MockPoiHandlerMockRecorder) Search(pos interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Search", reflect.TypeOf((*MockPoiHandler)(nil).Search), pos)
}

// Update mocks base method.
func (m *MockPoiHandler) Update(idToUpdate data.Id, updatedPoi *data.Poi) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", idToUpdate, updatedPoi)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockPoiHandlerMockRecorder) Update(idToUpdate, updatedPoi interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockPoiHandler)(nil).Update), idToUpdate, updatedPoi)
}
