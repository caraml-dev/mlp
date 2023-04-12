package api

import (
	"net/http"

	"github.com/caraml-dev/mlp/api/api"
	"github.com/caraml-dev/mlp/api/models/v2"
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
			Method:  http.MethodGet,
			Path:    "/applications",
			Handler: c.ListApplications,
			Name:    "ListApplications",
		},
	}
}
