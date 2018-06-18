// Code generated by mockery v1.0.0. DO NOT EDIT.
package mocks

import http "net/http"
import mock "github.com/stretchr/testify/mock"
import newrelic "github.com/newrelic/go-agent"
import time "time"

// Application is an autogenerated mock type for the Application type
type Application struct {
	mock.Mock
}

// RecordCustomEvent provides a mock function with given fields: eventType, params
func (_m *Application) RecordCustomEvent(eventType string, params map[string]interface{}) error {
	ret := _m.Called(eventType, params)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, map[string]interface{}) error); ok {
		r0 = rf(eventType, params)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RecordCustomMetric provides a mock function with given fields: name, value
func (_m *Application) RecordCustomMetric(name string, value float64) error {
	ret := _m.Called(name, value)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, float64) error); ok {
		r0 = rf(name, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Shutdown provides a mock function with given fields: timeout
func (_m *Application) Shutdown(timeout time.Duration) {
	_m.Called(timeout)
}

// StartTransaction provides a mock function with given fields: name, w, r
func (_m *Application) StartTransaction(name string, w http.ResponseWriter, r *http.Request) newrelic.Transaction {
	ret := _m.Called(name, w, r)

	var r0 newrelic.Transaction
	if rf, ok := ret.Get(0).(func(string, http.ResponseWriter, *http.Request) newrelic.Transaction); ok {
		r0 = rf(name, w, r)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(newrelic.Transaction)
		}
	}

	return r0
}

// WaitForConnection provides a mock function with given fields: timeout
func (_m *Application) WaitForConnection(timeout time.Duration) error {
	ret := _m.Called(timeout)

	var r0 error
	if rf, ok := ret.Get(0).(func(time.Duration) error); ok {
		r0 = rf(timeout)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}