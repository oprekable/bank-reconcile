// Code generated by mockery v2.53.3. DO NOT EDIT.

package _mock

import mock "github.com/stretchr/testify/mock"

// ISystemCalls is an autogenerated mock type for the ISystemCalls type
type ISystemCalls struct {
	mock.Mock
}

// Executable provides a mock function with no fields
func (_m *ISystemCalls) Executable() (string, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Executable")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func() (string, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FilepathDir provides a mock function with given fields: path
func (_m *ISystemCalls) FilepathDir(path string) string {
	ret := _m.Called(path)

	if len(ret) == 0 {
		panic("no return value specified for FilepathDir")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(path)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Getwd provides a mock function with no fields
func (_m *ISystemCalls) Getwd() (string, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Getwd")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func() (string, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewISystemCalls creates a new instance of ISystemCalls. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewISystemCalls(t interface {
	mock.TestingT
	Cleanup(func())
}) *ISystemCalls {
	mock := &ISystemCalls{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
