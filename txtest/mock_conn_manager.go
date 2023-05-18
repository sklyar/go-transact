// Code generated by mockery v2.26.1. DO NOT EDIT.

package txtest

import (
	context "context"
	sql "database/sql"
	time "time"

	mock "github.com/stretchr/testify/mock"
)

// ConnManager is an autogenerated mock type for the ConnManager type
type ConnManager struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *ConnManager) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Conn provides a mock function with given fields: ctx
func (_m *ConnManager) Conn(ctx context.Context) (*sql.Conn, error) {
	ret := _m.Called(ctx)

	var r0 *sql.Conn
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*sql.Conn, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *sql.Conn); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sql.Conn)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetConnMaxIdleTime provides a mock function with given fields: d
func (_m *ConnManager) SetConnMaxIdleTime(d time.Duration) {
	_m.Called(d)
}

// SetConnMaxLifetime provides a mock function with given fields: d
func (_m *ConnManager) SetConnMaxLifetime(d time.Duration) {
	_m.Called(d)
}

// SetMaxIdleConns provides a mock function with given fields: n
func (_m *ConnManager) SetMaxIdleConns(n int) {
	_m.Called(n)
}

// SetMaxOpenConns provides a mock function with given fields: n
func (_m *ConnManager) SetMaxOpenConns(n int) {
	_m.Called(n)
}

type mockConstructorTestingTNewConnManager interface {
	mock.TestingT
	Cleanup(func())
}

// NewConnManager creates a new instance of ConnManager. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewConnManager(t mockConstructorTestingTNewConnManager) *ConnManager {
	mock := &ConnManager{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
