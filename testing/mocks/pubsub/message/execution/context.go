// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/go-foreman/foreman/pubsub/message/execution (interfaces: MessageExecutionCtx,MessageExecutionCtxFactory)

// Package execution is a generated GoMock package.
package execution

import (
	context "context"
	reflect "reflect"

	log "github.com/go-foreman/foreman/log"
	endpoint "github.com/go-foreman/foreman/pubsub/endpoint"
	message "github.com/go-foreman/foreman/pubsub/message"
	execution "github.com/go-foreman/foreman/pubsub/message/execution"
	gomock "github.com/golang/mock/gomock"
)

// MockMessageExecutionCtx is a mock of MessageExecutionCtx interface.
type MockMessageExecutionCtx struct {
	ctrl     *gomock.Controller
	recorder *MockMessageExecutionCtxMockRecorder
}

// MockMessageExecutionCtxMockRecorder is the mock recorder for MockMessageExecutionCtx.
type MockMessageExecutionCtxMockRecorder struct {
	mock *MockMessageExecutionCtx
}

// NewMockMessageExecutionCtx creates a new mock instance.
func NewMockMessageExecutionCtx(ctrl *gomock.Controller) *MockMessageExecutionCtx {
	mock := &MockMessageExecutionCtx{ctrl: ctrl}
	mock.recorder = &MockMessageExecutionCtxMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMessageExecutionCtx) EXPECT() *MockMessageExecutionCtxMockRecorder {
	return m.recorder
}

// Context mocks base method.
func (m *MockMessageExecutionCtx) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockMessageExecutionCtxMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockMessageExecutionCtx)(nil).Context))
}

// Logger mocks base method.
func (m *MockMessageExecutionCtx) Logger() log.Logger {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Logger")
	ret0, _ := ret[0].(log.Logger)
	return ret0
}

// Logger indicates an expected call of Logger.
func (mr *MockMessageExecutionCtxMockRecorder) Logger() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Logger", reflect.TypeOf((*MockMessageExecutionCtx)(nil).Logger))
}

// Message mocks base method.
func (m *MockMessageExecutionCtx) Message() *message.ReceivedMessage {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Message")
	ret0, _ := ret[0].(*message.ReceivedMessage)
	return ret0
}

// Message indicates an expected call of Message.
func (mr *MockMessageExecutionCtxMockRecorder) Message() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Message", reflect.TypeOf((*MockMessageExecutionCtx)(nil).Message))
}

// Return mocks base method.
func (m *MockMessageExecutionCtx) Return(arg0 ...endpoint.DeliveryOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Return", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Return indicates an expected call of Return.
func (mr *MockMessageExecutionCtxMockRecorder) Return(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Return", reflect.TypeOf((*MockMessageExecutionCtx)(nil).Return), arg0...)
}

// Send mocks base method.
func (m *MockMessageExecutionCtx) Send(arg0 *message.OutcomingMessage, arg1 ...endpoint.DeliveryOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Send", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockMessageExecutionCtxMockRecorder) Send(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockMessageExecutionCtx)(nil).Send), varargs...)
}

// Valid mocks base method.
func (m *MockMessageExecutionCtx) Valid() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Valid")
	ret0, _ := ret[0].(bool)
	return ret0
}

// Valid indicates an expected call of Valid.
func (mr *MockMessageExecutionCtxMockRecorder) Valid() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Valid", reflect.TypeOf((*MockMessageExecutionCtx)(nil).Valid))
}

// MockMessageExecutionCtxFactory is a mock of MessageExecutionCtxFactory interface.
type MockMessageExecutionCtxFactory struct {
	ctrl     *gomock.Controller
	recorder *MockMessageExecutionCtxFactoryMockRecorder
}

// MockMessageExecutionCtxFactoryMockRecorder is the mock recorder for MockMessageExecutionCtxFactory.
type MockMessageExecutionCtxFactoryMockRecorder struct {
	mock *MockMessageExecutionCtxFactory
}

// NewMockMessageExecutionCtxFactory creates a new mock instance.
func NewMockMessageExecutionCtxFactory(ctrl *gomock.Controller) *MockMessageExecutionCtxFactory {
	mock := &MockMessageExecutionCtxFactory{ctrl: ctrl}
	mock.recorder = &MockMessageExecutionCtxFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMessageExecutionCtxFactory) EXPECT() *MockMessageExecutionCtxFactoryMockRecorder {
	return m.recorder
}

// CreateCtx mocks base method.
func (m *MockMessageExecutionCtxFactory) CreateCtx(arg0 context.Context, arg1 *message.ReceivedMessage) execution.MessageExecutionCtx {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCtx", arg0, arg1)
	ret0, _ := ret[0].(execution.MessageExecutionCtx)
	return ret0
}

// CreateCtx indicates an expected call of CreateCtx.
func (mr *MockMessageExecutionCtxFactoryMockRecorder) CreateCtx(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCtx", reflect.TypeOf((*MockMessageExecutionCtxFactory)(nil).CreateCtx), arg0, arg1)
}
