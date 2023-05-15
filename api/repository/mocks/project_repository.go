// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	models "github.com/caraml-dev/mlp/api/models"
)

// ProjectRepository is an autogenerated mock type for the ProjectRepository type
type ProjectRepository struct {
	mock.Mock
}

// Get provides a mock function with given fields: projectID
func (_m *ProjectRepository) Get(projectID models.ID) (*models.Project, error) {
	ret := _m.Called(projectID)

	var r0 *models.Project
	if rf, ok := ret.Get(0).(func(models.ID) *models.Project); ok {
		r0 = rf(projectID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Project)
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

// GetByName provides a mock function with given fields: projectName
func (_m *ProjectRepository) GetByName(projectName string) (*models.Project, error) {
	ret := _m.Called(projectName)

	var r0 *models.Project
	if rf, ok := ret.Get(0).(func(string) *models.Project); ok {
		r0 = rf(projectName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Project)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(projectName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListProjects provides a mock function with given fields: name
func (_m *ProjectRepository) ListProjects(name string) ([]*models.Project, error) {
	ret := _m.Called(name)

	var r0 []*models.Project
	if rf, ok := ret.Get(0).(func(string) []*models.Project); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Project)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: project
func (_m *ProjectRepository) Save(project *models.Project) (*models.Project, error) {
	ret := _m.Called(project)

	var r0 *models.Project
	if rf, ok := ret.Get(0).(func(*models.Project) *models.Project); ok {
		r0 = rf(project)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Project)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*models.Project) error); ok {
		r1 = rf(project)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewProjectRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewProjectRepository creates a new instance of ProjectRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewProjectRepository(t mockConstructorTestingTNewProjectRepository) *ProjectRepository {
	mock := &ProjectRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
