package service

import (
	wh "github.com/caraml-dev/mlp/api/pkg/webhooks"
)

type ServiceEvents struct{}

const (
	ProjectServiceType wh.ServiceType = "project"
)

const (
	ProjectCreatedEvent wh.EventType = "OnProjectCreated"
	ProjectUpdatedEvent wh.EventType = "OnProjectUpdated"
)

var EventList = []wh.EventType{
	ProjectCreatedEvent,
	ProjectUpdatedEvent,
}
