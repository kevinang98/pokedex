// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	echo "github.com/labstack/echo/v4"

	mock "github.com/stretchr/testify/mock"

	models "pokedex/src/models"
)

// Middleware is an autogenerated mock type for the Middleware type
type Middleware struct {
	mock.Mock
}

// GenerateToken provides a mock function with given fields: id, username, role
func (_m *Middleware) GenerateToken(id int, username string, role string) (string, error) {
	ret := _m.Called(id, username, role)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(int, string, string) (string, error)); ok {
		return rf(id, username, role)
	}
	if rf, ok := ret.Get(0).(func(int, string, string) string); ok {
		r0 = rf(id, username, role)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(int, string, string) error); ok {
		r1 = rf(id, username, role)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ParseJwtToken provides a mock function with given fields: c
func (_m *Middleware) ParseJwtToken(c echo.Context) (int, string, string) {
	ret := _m.Called(c)

	var r0 int
	var r1 string
	var r2 string
	if rf, ok := ret.Get(0).(func(echo.Context) (int, string, string)); ok {
		return rf(c)
	}
	if rf, ok := ret.Get(0).(func(echo.Context) int); ok {
		r0 = rf(c)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(echo.Context) string); ok {
		r1 = rf(c)
	} else {
		r1 = ret.Get(1).(string)
	}

	if rf, ok := ret.Get(2).(func(echo.Context) string); ok {
		r2 = rf(c)
	} else {
		r2 = ret.Get(2).(string)
	}

	return r0, r1, r2
}

// ValidateIntPokedexBody provides a mock function with given fields: poke
func (_m *Middleware) ValidateIntPokedexBody(poke []int) error {
	ret := _m.Called(poke)

	var r0 error
	if rf, ok := ret.Get(0).(func([]int) error); ok {
		r0 = rf(poke)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ValidatePokedexBody provides a mock function with given fields: poke
func (_m *Middleware) ValidatePokedexBody(poke []models.Pokedex) error {
	ret := _m.Called(poke)

	var r0 error
	if rf, ok := ret.Get(0).(func([]models.Pokedex) error); ok {
		r0 = rf(poke)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewMiddleware interface {
	mock.TestingT
	Cleanup(func())
}

// NewMiddleware creates a new instance of Middleware. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMiddleware(t mockConstructorTestingTNewMiddleware) *Middleware {
	mock := &Middleware{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}