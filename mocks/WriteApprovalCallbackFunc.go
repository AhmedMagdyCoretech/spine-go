// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	api "github.com/enbility/spine-go/api"
	mock "github.com/stretchr/testify/mock"
)

// WriteApprovalCallbackFunc is an autogenerated mock type for the WriteApprovalCallbackFunc type
type WriteApprovalCallbackFunc struct {
	mock.Mock
}

type WriteApprovalCallbackFunc_Expecter struct {
	mock *mock.Mock
}

func (_m *WriteApprovalCallbackFunc) EXPECT() *WriteApprovalCallbackFunc_Expecter {
	return &WriteApprovalCallbackFunc_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: msg
func (_m *WriteApprovalCallbackFunc) Execute(msg *api.Message) {
	_m.Called(msg)
}

// WriteApprovalCallbackFunc_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type WriteApprovalCallbackFunc_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - msg *api.Message
func (_e *WriteApprovalCallbackFunc_Expecter) Execute(msg interface{}) *WriteApprovalCallbackFunc_Execute_Call {
	return &WriteApprovalCallbackFunc_Execute_Call{Call: _e.mock.On("Execute", msg)}
}

func (_c *WriteApprovalCallbackFunc_Execute_Call) Run(run func(msg *api.Message)) *WriteApprovalCallbackFunc_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*api.Message))
	})
	return _c
}

func (_c *WriteApprovalCallbackFunc_Execute_Call) Return() *WriteApprovalCallbackFunc_Execute_Call {
	_c.Call.Return()
	return _c
}

func (_c *WriteApprovalCallbackFunc_Execute_Call) RunAndReturn(run func(*api.Message)) *WriteApprovalCallbackFunc_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// NewWriteApprovalCallbackFunc creates a new instance of WriteApprovalCallbackFunc. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewWriteApprovalCallbackFunc(t interface {
	mock.TestingT
	Cleanup(func())
}) *WriteApprovalCallbackFunc {
	mock := &WriteApprovalCallbackFunc{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
