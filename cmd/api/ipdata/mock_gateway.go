// Code generated by MockGen. DO NOT EDIT.
// Source: gateway.go

// Package ipdata is a generated GoMock package.
package ipdata

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockGateway is a mock of Gateway interface.
type MockGateway struct {
	ctrl     *gomock.Controller
	recorder *MockGatewayMockRecorder
}

// MockGatewayMockRecorder is the mock recorder for MockGateway.
type MockGatewayMockRecorder struct {
	mock *MockGateway
}

// NewMockGateway creates a new mock instance.
func NewMockGateway(ctrl *gomock.Controller) *MockGateway {
	mock := &MockGateway{ctrl: ctrl}
	mock.recorder = &MockGatewayMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGateway) EXPECT() *MockGatewayMockRecorder {
	return m.recorder
}

// GetDataFromIP mocks base method.
func (m *MockGateway) GetDataFromIP(ctx context.Context, ip string) (IpData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDataFromIP", ctx, ip)
	ret0, _ := ret[0].(IpData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDataFromIP indicates an expected call of GetDataFromIP.
func (mr *MockGatewayMockRecorder) GetDataFromIP(ctx, ip interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDataFromIP", reflect.TypeOf((*MockGateway)(nil).GetDataFromIP), ctx, ip)
}

// GetIpCountByCountyName mocks base method.
func (m *MockGateway) GetIpCountByCountyName(ctx context.Context, countryName string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetIpCountByCountyName", ctx, countryName)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetIpCountByCountyName indicates an expected call of GetIpCountByCountyName.
func (mr *MockGatewayMockRecorder) GetIpCountByCountyName(ctx, countryName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIpCountByCountyName", reflect.TypeOf((*MockGateway)(nil).GetIpCountByCountyName), ctx, countryName)
}

// GetIspIpsByCountryCode mocks base method.
func (m *MockGateway) GetIspIpsByCountryCode(ctx context.Context, countryCode string, limit int) ([]IspIpCount, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetIspIpsByCountryCode", ctx, countryCode, limit)
	ret0, _ := ret[0].([]IspIpCount)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetIspIpsByCountryCode indicates an expected call of GetIspIpsByCountryCode.
func (mr *MockGatewayMockRecorder) GetIspIpsByCountryCode(ctx, countryCode, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIspIpsByCountryCode", reflect.TypeOf((*MockGateway)(nil).GetIspIpsByCountryCode), ctx, countryCode, limit)
}

// GetTopISPFromSwitzerland mocks base method.
func (m *MockGateway) GetTopISPFromSwitzerland(ctx context.Context) ([]IspIpCount, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTopISPFromSwitzerland", ctx)
	ret0, _ := ret[0].([]IspIpCount)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTopISPFromSwitzerland indicates an expected call of GetTopISPFromSwitzerland.
func (mr *MockGatewayMockRecorder) GetTopISPFromSwitzerland(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTopISPFromSwitzerland", reflect.TypeOf((*MockGateway)(nil).GetTopISPFromSwitzerland), ctx)
}
