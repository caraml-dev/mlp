package api

import (
	"net/http"

	"github.com/gojek/mlp/api/api"
	"github.com/gojek/mlp/api/models/v2"
)

type ApplicationsController struct {
	Apps []models.Application
}

func (c *ApplicationsController) ListApplications(_ *http.Request, _ map[string]string, _ interface{}) *api.Response {
	return api.Ok(c.Apps)
}

func (c *ApplicationsController) Routes() []api.Route {
	return []api.Route{
		{
			http.MethodGet,
			"/applications",
			nil,
			c.ListApplications,
			"ListApplications",
		},
	}
}
