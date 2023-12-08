// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	context "context"

	artifact "github.com/caraml-dev/mlp/api/pkg/artifact"

	mock "github.com/stretchr/testify/mock"
)

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// DeleteArtifact provides a mock function with given fields: ctx, url
func (_m *Service) DeleteArtifact(ctx context.Context, url string) error {
	ret := _m.Called(ctx, url)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, url)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ParseURL provides a mock function with given fields: gsURL
func (_m *Service) ParseURL(gsURL string) (*artifact.URL, error) {
	ret := _m.Called(gsURL)

	var r0 *artifact.URL
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*artifact.URL, error)); ok {
		return rf(gsURL)
	}
	if rf, ok := ret.Get(0).(func(string) *artifact.URL); ok {
		r0 = rf(gsURL)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*artifact.URL)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(gsURL)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReadArtifact provides a mock function with given fields: ctx, url
func (_m *Service) ReadArtifact(ctx context.Context, url string) ([]byte, error) {
	ret := _m.Called(ctx, url)

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]byte, error)); ok {
		return rf(ctx, url)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []byte); ok {
		r0 = rf(ctx, url)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, url)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WriteArtifact provides a mock function with given fields: ctx, url, content
func (_m *Service) WriteArtifact(ctx context.Context, url string, content string) error {
	ret := _m.Called(ctx, url, content)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, url, content)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewService interface {
	mock.TestingT
	Cleanup(func())
}

// NewService creates a new instance of Service. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewService(t mockConstructorTestingTNewService) *Service {
	mock := &Service{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
