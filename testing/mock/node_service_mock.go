// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1 (interfaces: NodeClient)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
	grpc "google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// MockNodeClient is a mock of NodeClient interface
type MockNodeClient struct {
	ctrl     *gomock.Controller
	recorder *MockNodeClientMockRecorder
}

// MockNodeClientMockRecorder is the mock recorder for MockNodeClient
type MockNodeClientMockRecorder struct {
	mock *MockNodeClient
}

// NewMockNodeClient creates a new mock instance
func NewMockNodeClient(ctrl *gomock.Controller) *MockNodeClient {
	mock := &MockNodeClient{ctrl: ctrl}
	mock.recorder = &MockNodeClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockNodeClient) EXPECT() *MockNodeClientMockRecorder {
	return m.recorder
}

// GetGenesis mocks base method
func (m *MockNodeClient) GetGenesis(arg0 context.Context, arg1 *emptypb.Empty, arg2 ...grpc.CallOption) (*eth.Genesis, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetGenesis", varargs...)
	ret0, _ := ret[0].(*eth.Genesis)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGenesis indicates an expected call of GetGenesis
func (mr *MockNodeClientMockRecorder) GetGenesis(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGenesis", reflect.TypeOf((*MockNodeClient)(nil).GetGenesis), varargs...)
}

// GetHost mocks base method
func (m *MockNodeClient) GetHost(arg0 context.Context, arg1 *emptypb.Empty, arg2 ...grpc.CallOption) (*eth.HostData, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetHost", varargs...)
	ret0, _ := ret[0].(*eth.HostData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetHost indicates an expected call of GetHost
func (mr *MockNodeClientMockRecorder) GetHost(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHost", reflect.TypeOf((*MockNodeClient)(nil).GetHost), varargs...)
}

// GetPeer mocks base method
func (m *MockNodeClient) GetPeer(arg0 context.Context, arg1 *eth.PeerRequest, arg2 ...grpc.CallOption) (*eth.Peer, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetPeer", varargs...)
	ret0, _ := ret[0].(*eth.Peer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPeer indicates an expected call of GetPeer
func (mr *MockNodeClientMockRecorder) GetPeer(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPeer", reflect.TypeOf((*MockNodeClient)(nil).GetPeer), varargs...)
}

// GetSyncStatus mocks base method
func (m *MockNodeClient) GetSyncStatus(arg0 context.Context, arg1 *emptypb.Empty, arg2 ...grpc.CallOption) (*eth.SyncStatus, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetSyncStatus", varargs...)
	ret0, _ := ret[0].(*eth.SyncStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSyncStatus indicates an expected call of GetSyncStatus
func (mr *MockNodeClientMockRecorder) GetSyncStatus(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSyncStatus", reflect.TypeOf((*MockNodeClient)(nil).GetSyncStatus), varargs...)
}

// GetVersion mocks base method
func (m *MockNodeClient) GetVersion(arg0 context.Context, arg1 *emptypb.Empty, arg2 ...grpc.CallOption) (*eth.Version, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetVersion", varargs...)
	ret0, _ := ret[0].(*eth.Version)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVersion indicates an expected call of GetVersion
func (mr *MockNodeClientMockRecorder) GetVersion(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVersion", reflect.TypeOf((*MockNodeClient)(nil).GetVersion), varargs...)
}

// ListImplementedServices mocks base method
func (m *MockNodeClient) ListImplementedServices(arg0 context.Context, arg1 *emptypb.Empty, arg2 ...grpc.CallOption) (*eth.ImplementedServices, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListImplementedServices", varargs...)
	ret0, _ := ret[0].(*eth.ImplementedServices)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListImplementedServices indicates an expected call of ListImplementedServices
func (mr *MockNodeClientMockRecorder) ListImplementedServices(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListImplementedServices", reflect.TypeOf((*MockNodeClient)(nil).ListImplementedServices), varargs...)
}

// ListPeers mocks base method
func (m *MockNodeClient) ListPeers(arg0 context.Context, arg1 *emptypb.Empty, arg2 ...grpc.CallOption) (*eth.Peers, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListPeers", varargs...)
	ret0, _ := ret[0].(*eth.Peers)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListPeers indicates an expected call of ListPeers
func (mr *MockNodeClientMockRecorder) ListPeers(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListPeers", reflect.TypeOf((*MockNodeClient)(nil).ListPeers), varargs...)
}
