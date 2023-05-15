// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	models "github.com/caraml-dev/mlp/api/models"
)

// SecretRepository is an autogenerated mock type for the SecretRepository type
type SecretRepository struct {
	mock.Mock
}

// Delete provides a mock function with given fields: id
func (_m *SecretRepository) Delete(id models.ID) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(models.ID) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: id
func (_m *SecretRepository) Get(id models.ID) (*models.Secret, error) {
	ret := _m.Called(id)

	var r0 *models.Secret
	if rf, ok := ret.Get(0).(func(models.ID) *models.Secret); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Secret)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(models.ID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: projectID
func (_m *SecretRepository) List(projectID models.ID) ([]*models.Secret, error) {
	ret := _m.Called(projectID)

	var r0 []*models.Secret
	if rf, ok := ret.Get(0).(func(models.ID) []*models.Secret); ok {
		r0 = rf(projectID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Secret)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(models.ID) error); ok {
		r1 = rf(projectID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: secret
func (_m *SecretRepository) Save(secret *models.Secret) (*models.Secret, error) {
	ret := _m.Called(secret)

	var r0 *models.Secret
	if rf, ok := ret.Get(0).(func(*models.Secret) *models.Secret); ok {
		r0 = rf(secret)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Secret)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*models.Secret) error); ok {
		r1 = rf(secret)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewSecretRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewSecretRepository creates a new instance of SecretRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewSecretRepository(t mockConstructorTestingTNewSecretRepository) *SecretRepository {
	mock := &SecretRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
