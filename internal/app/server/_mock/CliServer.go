// Code generated by mockery v2.53.3. DO NOT EDIT.

package _mock

import (
	mock "github.com/stretchr/testify/mock"
	errgroup "golang.org/x/sync/errgroup"
)

// CliServer is an autogenerated mock type for the CliServer type
type CliServer struct {
	mock.Mock
}

// Name provides a mock function with no fields
func (_m *CliServer) Name() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Name")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Shutdown provides a mock function with no fields
func (_m *CliServer) Shutdown() {
	_m.Called()
}

// Start provides a mock function with given fields: eg
func (_m *CliServer) Start(eg *errgroup.Group) {
	_m.Called(eg)
}

// NewCliServer creates a new instance of CliServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCliServer(t interface {
	mock.TestingT
	Cleanup(func())
}) *CliServer {
	mock := &CliServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
