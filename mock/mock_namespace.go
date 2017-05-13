// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/johnny-morrice/godless (interfaces: RowConsumer)

package mock_godless

import (
	gomock "github.com/golang/mock/gomock"
	godless "github.com/johnny-morrice/godless"
)

// Mock of RowConsumer interface
type MockRowConsumer struct {
	ctrl     *gomock.Controller
	recorder *_MockRowConsumerRecorder
}

// Recorder for MockRowConsumer (not exported)
type _MockRowConsumerRecorder struct {
	mock *MockRowConsumer
}

func NewMockRowConsumer(ctrl *gomock.Controller) *MockRowConsumer {
	mock := &MockRowConsumer{ctrl: ctrl}
	mock.recorder = &_MockRowConsumerRecorder{mock}
	return mock
}

func (_m *MockRowConsumer) EXPECT() *_MockRowConsumerRecorder {
	return _m.recorder
}

func (_m *MockRowConsumer) Accept(_param0 godless.RowName, _param1 godless.Row) {
	_m.ctrl.Call(_m, "Accept", _param0, _param1)
}

func (_mr *_MockRowConsumerRecorder) Accept(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Accept", arg0, arg1)
}
