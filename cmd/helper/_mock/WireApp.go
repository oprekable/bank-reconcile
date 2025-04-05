// Code generated by mockery v2.53.3. DO NOT EDIT.

package _mock

import (
	appcontext "github.com/oprekable/bank-reconcile/internal/app/appcontext"
	cconfig "github.com/oprekable/bank-reconcile/internal/app/component/cconfig"

	clogger "github.com/oprekable/bank-reconcile/internal/app/component/clogger"

	context "context"

	core "github.com/oprekable/bank-reconcile/internal/app/err/core"

	csqlite "github.com/oprekable/bank-reconcile/internal/app/component/csqlite"

	embed "embed"

	mock "github.com/stretchr/testify/mock"
)

// WireApp is an autogenerated mock type for the WireApp type
type WireApp struct {
	mock.Mock
}

// Execute provides a mock function with given fields: ctx, embedFS, appName, tz, errType, isShowLog, dBPath
func (_m *WireApp) Execute(ctx context.Context, embedFS *embed.FS, appName cconfig.AppName, tz cconfig.TimeZone, errType []core.ErrorType, isShowLog clogger.IsShowLog, dBPath csqlite.DBPath) (*appcontext.AppContext, func(), error) {
	ret := _m.Called(ctx, embedFS, appName, tz, errType, isShowLog, dBPath)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 *appcontext.AppContext
	var r1 func()
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, *embed.FS, cconfig.AppName, cconfig.TimeZone, []core.ErrorType, clogger.IsShowLog, csqlite.DBPath) (*appcontext.AppContext, func(), error)); ok {
		return rf(ctx, embedFS, appName, tz, errType, isShowLog, dBPath)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *embed.FS, cconfig.AppName, cconfig.TimeZone, []core.ErrorType, clogger.IsShowLog, csqlite.DBPath) *appcontext.AppContext); ok {
		r0 = rf(ctx, embedFS, appName, tz, errType, isShowLog, dBPath)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*appcontext.AppContext)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *embed.FS, cconfig.AppName, cconfig.TimeZone, []core.ErrorType, clogger.IsShowLog, csqlite.DBPath) func()); ok {
		r1 = rf(ctx, embedFS, appName, tz, errType, isShowLog, dBPath)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(func())
		}
	}

	if rf, ok := ret.Get(2).(func(context.Context, *embed.FS, cconfig.AppName, cconfig.TimeZone, []core.ErrorType, clogger.IsShowLog, csqlite.DBPath) error); ok {
		r2 = rf(ctx, embedFS, appName, tz, errType, isShowLog, dBPath)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// NewWireApp creates a new instance of WireApp. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewWireApp(t interface {
	mock.TestingT
	Cleanup(func())
}) *WireApp {
	mock := &WireApp{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
