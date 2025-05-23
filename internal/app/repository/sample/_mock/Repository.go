// Code generated by mockery v2.53.3. DO NOT EDIT.

package _mock

import (
	context "context"

	sample "github.com/oprekable/bank-reconcile/internal/app/repository/sample"
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// Close provides a mock function with no fields
func (_m *Repository) Close() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Close")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetTrx provides a mock function with given fields: ctx
func (_m *Repository) GetTrx(ctx context.Context) ([]sample.TrxData, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetTrx")
	}

	var r0 []sample.TrxData
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]sample.TrxData, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []sample.TrxData); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]sample.TrxData)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Post provides a mock function with given fields: ctx
func (_m *Repository) Post(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Post")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Pre provides a mock function with given fields: ctx, listBank, startDate, toDate, limitTrxData, matchPercentage
func (_m *Repository) Pre(ctx context.Context, listBank []string, startDate time.Time, toDate time.Time, limitTrxData int64, matchPercentage int) error {
	ret := _m.Called(ctx, listBank, startDate, toDate, limitTrxData, matchPercentage)

	if len(ret) == 0 {
		panic("no return value specified for Pre")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []string, time.Time, time.Time, int64, int) error); ok {
		r0 = rf(ctx, listBank, startDate, toDate, limitTrxData, matchPercentage)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewRepository creates a new instance of Repository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *Repository {
	mock := &Repository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
