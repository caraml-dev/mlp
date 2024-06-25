package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/caraml-dev/mlp/api/log"
	"github.com/caraml-dev/mlp/api/models"
	apperror "github.com/caraml-dev/mlp/api/pkg/errors"
	"github.com/caraml-dev/mlp/api/pkg/webhooks"
)

type ProjectsController struct {
	*AppContext
}

func (c *ProjectsController) ListProjects(r *http.Request, vars map[string]string, _ interface{}) *Response {
	projects, err := c.ProjectsService.ListProjects(r.Context(), vars["name"], vars["user"])
	if err != nil {
		log.Errorf("error fetching projects: %s", err)
		return FromError(err)
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
		log.Errorf("invalid request body %v", body)
		return BadRequest("Unable to parse request body as project")
	}

	existingProject, err := c.ProjectsService.FindByName(project.Name)
	if existingProject != nil {
		log.Errorf("project %s already exists", project.Name)
		return BadRequest(fmt.Sprintf("Project %s already exists", project.Name))
	}

	if err != nil && !errors.Is(err, &apperror.NotFoundError{}) {
		log.Errorf("error fetching project with name %s: %s", project.Name, err)
		return InternalServerError(err.Error())
	}

	user := vars["user"]
	project.Administrators = addRequester(user, project.Administrators)
	project, err = c.ProjectsService.CreateProject(r.Context(), project)
	var webhookError *webhooks.WebhookError
	if err != nil {
		// NOTE: Here we are checking if the error is a WebhookError
		// This is to improve the error message shared with the user,
		// since the current logic creates the Project first before firing
		// the ProjectCreatedEvent webhook.
		if errors.As(err, &webhookError) {
			err := fmt.Errorf(`Project %s was created, 
			but not all webhooks were correctly invoked. 
			Some additional resources may not have been created successfully : %s`, project.Name, err)
			return FromError(err)
		}
		log.Errorf("error creating project %s: %s", project.Name, err)
		return FromError(err)
	}

	return Created(project)
}

func (c *ProjectsController) UpdateProject(r *http.Request, vars map[string]string, body interface{}) *Response {
	projectID, _ := models.ParseID(vars["project_id"])
	project, err := c.ProjectsService.FindByID(projectID)
	if err != nil {
		log.Errorf("error fetching project with id %s: %s", projectID, err)
		return FromError(err)
	}

	newProject, ok := body.(*models.Project)
	if !ok {
		log.Errorf("invalid request body %v", body)
		return BadRequest("Unable to parse request body as project")
	}

	project.Administrators = newProject.Administrators
	project.Readers = newProject.Readers
	project.Team = newProject.Team
	project.Stream = newProject.Stream
	project.Labels = newProject.Labels
	updatedProject, response, err := c.ProjectsService.UpdateProject(r.Context(), project)
	if err != nil {
		log.Errorf("error updating project %s: %s", project.Name, err)
		return FromError(err)
	}

	if response != nil {
		return Ok(response)
	}

	return Ok(updatedProject)
}

func (c *ProjectsController) GetProject(r *http.Request, vars map[string]string, body interface{}) *Response {
	projectID, _ := models.ParseID(vars["project_id"])
	project, err := c.ProjectsService.FindByID(projectID)
	if err != nil {
		log.Errorf("error fetching project with id %s: %s", projectID, err)
		return FromError(err)
	}

	return Ok(project)
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
