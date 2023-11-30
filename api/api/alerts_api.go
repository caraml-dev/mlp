package api

import (
	http "net/http"

	"github.com/caraml-dev/mlp/api/log"
)

type AlertController struct {
	*AppContext
}

func (c *AlertController) ListAlerts(r *http.Request, vars map[string]string, _ interface{}) *Response {
	ctx := r.Context()

	err := c.AlertService.List(ctx)
	if err != nil {
		log.Errorf("error fetching projects: %s", err)
		return FromError(err)
	}

	return Ok(nil)
}

func (c *AlertController) CreateAlert(r *http.Request, vars map[string]string, _ interface{}) *Response {
	ctx := r.Context()

	log.Infof("alertController.CreateAlert")
	err := c.AlertService.Create(ctx)
	if err != nil {
		log.Errorf("error fetching projects: %s", err)
		return FromError(err)
	}

	return Ok(nil)
}

func (c *AlertController) UpdateAlert(r *http.Request, vars map[string]string, _ interface{}) *Response {
	ctx := r.Context()

	log.Infof("alertController.UpdateAlert")
	err := c.AlertService.Update(ctx)
	if err != nil {
		log.Errorf("error fetching projects: %s", err)
		return FromError(err)
	}

	return Ok(nil)
}

func (c *AlertController) Routes() []Route {
	return []Route{
		{
			http.MethodGet,
			"/alerts",
			nil,
			c.ListAlerts,
			"ListAlerts",
		},
		{
			http.MethodPost,
			"/alerts",
			nil,
			c.CreateAlert,
			"CreateAlert",
		},
		{
			http.MethodPut,
			"/alerts",
			nil,
			c.UpdateAlert,
			"UpdateAlert",
		},
	}
}
