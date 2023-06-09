// Code generated by mockery v2.26.1. DO NOT EDIT.

package txtest

import (
	txsql "github.com/sklyar/go-transact/txsql"
	mock "github.com/stretchr/testify/mock"
)

// Stmt is an autogenerated mock type for the Stmt type
type Stmt struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *Stmt) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Exec provides a mock function with given fields: args
func (_m *Stmt) Exec(args ...interface{}) (txsql.Result, error) {
	var _ca []interface{}
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 txsql.Result
	var r1 error
	if rf, ok := ret.Get(0).(func(...interface{}) (txsql.Result, error)); ok {
		return rf(args...)
	}
	if rf, ok := ret.Get(0).(func(...interface{}) txsql.Result); ok {
		r0 = rf(args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(txsql.Result)
		}
	}

	if rf, ok := ret.Get(1).(func(...interface{}) error); ok {
		r1 = rf(args...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Query provides a mock function with given fields: args
func (_m *Stmt) Query(args ...interface{}) (txsql.Rows, error) {
	var _ca []interface{}
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 txsql.Rows
	var r1 error
	if rf, ok := ret.Get(0).(func(...interface{}) (txsql.Rows, error)); ok {
		return rf(args...)
	}
	if rf, ok := ret.Get(0).(func(...interface{}) txsql.Rows); ok {
		r0 = rf(args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(txsql.Rows)
		}
	}

	if rf, ok := ret.Get(1).(func(...interface{}) error); ok {
		r1 = rf(args...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// QueryRow provides a mock function with given fields: args
func (_m *Stmt) QueryRow(args ...interface{}) txsql.Row {
	var _ca []interface{}
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 txsql.Row
	if rf, ok := ret.Get(0).(func(...interface{}) txsql.Row); ok {
		r0 = rf(args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(txsql.Row)
		}
	}

	return r0
}

type mockConstructorTestingTNewStmt interface {
	mock.TestingT
	Cleanup(func())
}

// NewStmt creates a new instance of Stmt. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewStmt(t mockConstructorTestingTNewStmt) *Stmt {
	mock := &Stmt{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
