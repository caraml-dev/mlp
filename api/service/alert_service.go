package service

import (
	"context"

	"github.com/caraml-dev/mlp/api/log"
	"github.com/caraml-dev/mlp/api/pkg/alert"
)

type AlertService interface {
	List(ctx context.Context) error
	// TODO(arief): Alert data as input parameter
	Create(ctx context.Context) error
	Update(ctx context.Context) error
}

func NewAlertService(alert alert.Alert) (AlertService, error) {
	return &alertService{
		AlertClient: alert,
	}, nil
}

type alertService struct {
	AlertClient alert.Alert
}

func (service *alertService) List(ctx context.Context) error {
	log.Infof("alertService.List")
	service.AlertClient.List(ctx)
	return nil
}

func (service *alertService) Create(ctx context.Context) error {
	log.Infof("alertService.Create")
	return service.AlertClient.Create(ctx)
}

func (service *alertService) Update(ctx context.Context) error {
	log.Infof("alertService.Update")
	return service.AlertClient.Update(ctx)
}
