package api

import (
	http "net/http"
)

type ApplicationsController struct {
	*AppContext
}

func (c *ApplicationsController) ListApplications(_ *http.Request, _ map[string]string, _ interface{}) *Response {
	applications, err := c.ApplicationService.List()
	if err != nil {
		return InternalServerError(err.Error())
	}
	return Ok(applications)
}

func (c *ApplicationsController) Routes() []Route {
	return []Route{
		{
			http.MethodGet,
			"/applications",
			nil,
			c.ListApplications,
			"ListApplications",
		},
	}
}
