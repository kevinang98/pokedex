// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	models "pokedex/src/models"

	mock "github.com/stretchr/testify/mock"
)

// UserService is an autogenerated mock type for the UserService type
type UserService struct {
	mock.Mock
}

// LoginUser provides a mock function with given fields: username, password
func (_m *UserService) LoginUser(username string, password string) (*models.User, error) {
	ret := _m.Called(username, password)

	var r0 *models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (*models.User, error)); ok {
		return rf(username, password)
	}
	if rf, ok := ret.Get(0).(func(string, string) *models.User); ok {
		r0 = rf(username, password)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(username, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RegisterUser provides a mock function with given fields: username, password, role
func (_m *UserService) RegisterUser(username string, password string, role string) error {
	ret := _m.Called(username, password, role)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string) error); ok {
		r0 = rf(username, password, role)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateCapturedPokemonUser provides a mock function with given fields: id, pid
func (_m *UserService) UpdateCapturedPokemonUser(id int, pid string) error {
	ret := _m.Called(id, pid)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, string) error); ok {
		r0 = rf(id, pid)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewUserService interface {
	mock.TestingT
	Cleanup(func())
}

// NewUserService creates a new instance of UserService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUserService(t mockConstructorTestingTNewUserService) *UserService {
	mock := &UserService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}