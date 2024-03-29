// Code generated by mockery v2.27.1. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	models "github.com/caraml-dev/mlp/api/models"
)

// SecretStorageRepository is an autogenerated mock type for the SecretStorageRepository type
type SecretStorageRepository struct {
	mock.Mock
}

// Delete provides a mock function with given fields: id
func (_m *SecretStorageRepository) Delete(id models.ID) error {
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
func (_m *SecretStorageRepository) Get(id models.ID) (*models.SecretStorage, error) {
	ret := _m.Called(id)

	var r0 *models.SecretStorage
	var r1 error
	if rf, ok := ret.Get(0).(func(models.ID) (*models.SecretStorage, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(models.ID) *models.SecretStorage); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.SecretStorage)
		}
	}

	if rf, ok := ret.Get(1).(func(models.ID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetGlobal provides a mock function with given fields: name
func (_m *SecretStorageRepository) GetGlobal(name string) (*models.SecretStorage, error) {
	ret := _m.Called(name)

	var r0 *models.SecretStorage
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*models.SecretStorage, error)); ok {
		return rf(name)
	}
	if rf, ok := ret.Get(0).(func(string) *models.SecretStorage); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.SecretStorage)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: projectID
func (_m *SecretStorageRepository) List(projectID models.ID) ([]*models.SecretStorage, error) {
	ret := _m.Called(projectID)

	var r0 []*models.SecretStorage
	var r1 error
	if rf, ok := ret.Get(0).(func(models.ID) ([]*models.SecretStorage, error)); ok {
		return rf(projectID)
	}
	if rf, ok := ret.Get(0).(func(models.ID) []*models.SecretStorage); ok {
		r0 = rf(projectID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.SecretStorage)
		}
	}

	if rf, ok := ret.Get(1).(func(models.ID) error); ok {
		r1 = rf(projectID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListAll provides a mock function with given fields:
func (_m *SecretStorageRepository) ListAll() ([]*models.SecretStorage, error) {
	ret := _m.Called()

	var r0 []*models.SecretStorage
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]*models.SecretStorage, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []*models.SecretStorage); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.SecretStorage)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListGlobal provides a mock function with given fields:
func (_m *SecretStorageRepository) ListGlobal() ([]*models.SecretStorage, error) {
	ret := _m.Called()

	var r0 []*models.SecretStorage
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]*models.SecretStorage, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []*models.SecretStorage); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.SecretStorage)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: secretStorage
func (_m *SecretStorageRepository) Save(secretStorage *models.SecretStorage) (*models.SecretStorage, error) {
	ret := _m.Called(secretStorage)

	var r0 *models.SecretStorage
	var r1 error
	if rf, ok := ret.Get(0).(func(*models.SecretStorage) (*models.SecretStorage, error)); ok {
		return rf(secretStorage)
	}
	if rf, ok := ret.Get(0).(func(*models.SecretStorage) *models.SecretStorage); ok {
		r0 = rf(secretStorage)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.SecretStorage)
		}
	}

	if rf, ok := ret.Get(1).(func(*models.SecretStorage) error); ok {
		r1 = rf(secretStorage)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewSecretStorageRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewSecretStorageRepository creates a new instance of SecretStorageRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewSecretStorageRepository(t mockConstructorTestingTNewSecretStorageRepository) *SecretStorageRepository {
	mock := &SecretStorageRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
