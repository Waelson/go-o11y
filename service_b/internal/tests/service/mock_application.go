// Code generated by MockGen. DO NOT EDIT.
// Source: ../internal/service/application.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	reflect "reflect"

	model "github.com/Waelson/go-o11y/internal/model"
	gomock "github.com/golang/mock/gomock"
)

// MockApplicationService is a mock of ApplicationService interface.
type MockApplicationService struct {
	ctrl     *gomock.Controller
	recorder *MockApplicationServiceMockRecorder
}

// MockApplicationServiceMockRecorder is the mock recorder for MockApplicationService.
type MockApplicationServiceMockRecorder struct {
	mock *MockApplicationService
}

// NewMockApplicationService creates a new mock instance.
func NewMockApplicationService(ctrl *gomock.Controller) *MockApplicationService {
	mock := &MockApplicationService{ctrl: ctrl}
	mock.recorder = &MockApplicationServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockApplicationService) EXPECT() *MockApplicationServiceMockRecorder {
	return m.recorder
}

// GetTemperature mocks base method.
func (m *MockApplicationService) GetTemperature(cep string) (model.ApplicationResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTemperature", cep)
	ret0, _ := ret[0].(model.ApplicationResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTemperature indicates an expected call of GetTemperature.
func (mr *MockApplicationServiceMockRecorder) GetTemperature(cep interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTemperature", reflect.TypeOf((*MockApplicationService)(nil).GetTemperature), cep)
}