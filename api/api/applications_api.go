package api

import "net/http"

type ApplicationsController struct {
	*AppContext
}

func (c *ApplicationsController) ListApplications(r *http.Request, _ map[string]string, _ interface{}) *ApiResponse {
	applications, err := c.ApplicationService.List()
	if err != nil {
		return InternalServerError(err.Error())
	}
	return Ok(applications)

}
