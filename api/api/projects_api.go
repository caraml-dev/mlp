package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/caraml-dev/mlp/api/log"
	"github.com/caraml-dev/mlp/api/models"
	"github.com/caraml-dev/mlp/api/pkg/authz/enforcer"
	"github.com/jinzhu/gorm"
)

type ProjectsController struct {
	*AppContext
}

func (c *ProjectsController) ListProjects(r *http.Request, vars map[string]string, _ interface{}) *Response {
	projects, err := c.ProjectsService.ListProjects(vars["name"])
	if err != nil {
		return InternalServerError(err.Error())
	}

	user := vars["user"]
	projects, err = c.filterAuthorizedProjects(user, projects, enforcer.ActionRead)
	if err != nil {
		return InternalServerError(err.Error())
	}

	return Ok(projects)
}

func (c *ProjectsController) CreateProject(r *http.Request, vars map[string]string, body interface{}) *Response {
	userAgent := strings.ToLower(r.Header.Get("User-Agent"))
	if strings.Contains(userAgent, "swagger") {
		return Forbidden("Project creation from SDK is disabled. Use the MLP console to create a project.")
	}

	project, ok := body.(*models.Project)
	if !ok {
		return BadRequest("Unable to parse request body as project")
	}

	existingProject, err := c.ProjectsService.FindByName(project.Name)
	if existingProject != nil {
		return BadRequest(fmt.Sprintf("Project %s already exists", project.Name))
	}

	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return InternalServerError(err.Error())
	}

	user := vars["user"]
	project.Administrators = addRequester(user, project.Administrators)
	project, err = c.ProjectsService.CreateProject(project)
	if err != nil {
		return InternalServerError(err.Error())
	}

	return Created(project)
}

func (c *ProjectsController) UpdateProject(r *http.Request, vars map[string]string, body interface{}) *Response {
	projectID, _ := models.ParseID(vars["project_id"])
	project, err := c.ProjectsService.FindByID(projectID)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return NotFound(fmt.Sprintf("Project id %s not found", projectID))
		}

		return InternalServerError(err.Error())
	}

	newProject, ok := body.(*models.Project)
	if !ok {
		return BadRequest("Unable to parse request body as project")
	}

	project.Administrators = newProject.Administrators
	project.Readers = newProject.Readers
	project.Team = newProject.Team
	project.Stream = newProject.Stream
	project.Labels = newProject.Labels
	project, err = c.ProjectsService.UpdateProject(project)
	if err != nil {
		log.Errorf("unable to update project %s %v", projectID, err)
		return InternalServerError(fmt.Sprintf("Unable to update project %s", projectID))
	}

	return Ok(project)
}

func (c *ProjectsController) GetProject(r *http.Request, vars map[string]string, body interface{}) *Response {
	projectID, _ := models.ParseID(vars["project_id"])
	project, err := c.ProjectsService.FindByID(projectID)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return NotFound(fmt.Sprintf("Project id %s not found", projectID))
		}

		return InternalServerError(err.Error())
	}

	return Ok(project)
}

func (c *ProjectsController) filterAuthorizedProjects(
	user string,
	projects []*models.Project,
	action string,
) ([]*models.Project, error) {
	if c.AuthorizationEnabled {
		projectIds := make([]string, 0)
		allowedProjects := make([]*models.Project, 0)
		projectMap := make(map[string]*models.Project)
		for _, project := range projects {
			projectID := fmt.Sprintf("projects:%s", project.ID)
			projectIds = append(projectIds, projectID)
			projectMap[projectID] = project
		}

		allowedProjectIds, err := c.Enforcer.FilterAuthorizedResource(user, projectIds, action)
		if err != nil {
			return nil, err
		}

		for _, projectID := range allowedProjectIds {
			allowedProjects = append(allowedProjects, projectMap[projectID])
		}

		return allowedProjects, nil
	}

	return projects, nil
}

func (c *ProjectsController) Routes() []Route {
	return []Route{
		{
			http.MethodGet,
			"/projects/{project_id:[0-9]+}",
			nil,
			c.GetProject,
			"GetProject",
		},
		{
			http.MethodGet,
			"/projects",
			nil,
			c.ListProjects,
			"ListProjects",
		},
		{
			http.MethodPost,
			"/projects",
			models.Project{},
			c.CreateProject,
			"CreateProject",
		},
		{
			http.MethodPut,
			"/projects/{project_id:[0-9]+}",
			models.Project{},
			c.UpdateProject,
			"UpdateProject",
		},
	}
}

// addRequester add requester to users slice if it doesn't exists
func addRequester(requester string, users []string) []string {
	for _, user := range users {
		if user == requester {
			return users
		}
	}

	return append(users, requester)
}
