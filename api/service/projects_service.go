package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"bytes"
	"encoding/json"
	"html/template"
	"net/http"

	"golang.org/x/exp/slices"

	"github.com/pkg/errors"

	"github.com/caraml-dev/mlp/api/repository"

	"github.com/caraml-dev/mlp/api/models"
	"github.com/caraml-dev/mlp/api/pkg/authz/enforcer"
	"github.com/caraml-dev/mlp/api/pkg/webhooks"
)

type ProjectsService interface {
	ListProjects(ctx context.Context, name string, user string) ([]*models.Project, error)
	CreateProject(ctx context.Context, project *models.Project) (*models.Project, error)
	UpdateProject(ctx context.Context, project *models.Project) (*models.Project, string, error)
	FindByID(projectID models.ID) (*models.Project, error)
	FindByName(projectName string) (*models.Project, error)
	MakeRequestPayload(project *models.Project, templateString string) (string, error)
	SendUpdateRequest(url string, payload string) (*http.Response, error)
	ProcessResponseURL(response *http.Response, templateString string) (string, error)
}

var reservedProjectName = map[string]bool{
	"infrastructure":     true,
	"kube-system":        true,
	"knative-serving":    true,
	"kfserving-system":   true,
	"knative-monitoring": true,
}

func NewProjectsService(
	mlflowURL string,
	projectRepository repository.ProjectRepository,
	authEnforcer enforcer.Enforcer,
	authEnabled bool,
	webhookManager webhooks.WebhookManager,
	updateProjectEndpoint string,
	updateProjectPayloadTemplate string,
	updateProjectResponseTemplate string,
) (ProjectsService, error) {
	if strings.TrimSpace(mlflowURL) == "" {
		return nil, errors.New("default mlflow tracking url should be provided")
	}

	return &projectsService{
		projectRepository:             projectRepository,
		defaultMlflowTrackingServer:   mlflowURL,
		authEnforcer:                  authEnforcer,
		authEnabled:                   authEnabled,
		webhookManager:              webhookManager,
		updateProjectEndpoint:         updateProjectEndpoint,
		updateProjectPayloadTemplate:  updateProjectPayloadTemplate,
		updateProjectResponseTemplate: updateProjectResponseTemplate,
	}, nil
}

type projectsService struct {
	projectRepository             repository.ProjectRepository
	defaultMlflowTrackingServer   string
	authEnforcer                  enforcer.Enforcer
	authEnabled                   bool
	webhookManager              webhooks.WebhookManager
	updateProjectEndpoint         string
	updateProjectPayloadTemplate  string
	updateProjectResponseTemplate string
}

func (service *projectsService) CreateProject(ctx context.Context, project *models.Project) (*models.Project, error) {
	if _, ok := reservedProjectName[project.Name]; ok {
		return nil, fmt.Errorf("unable to use reserved project name: %s", project.Name)
	}

	if strings.TrimSpace(project.MLFlowTrackingURL) == "" {
		project.MLFlowTrackingURL = service.defaultMlflowTrackingServer
	}

	project, err := service.save(project)
	if err != nil {
		return nil, fmt.Errorf("unable to create new project")
	}

	if service.authEnabled {
		err = service.updateAuthorizationPolicy(ctx, project)
		if err != nil {
			return nil, fmt.Errorf("error while creating authorization policy for project %s", project.Name)
		}
	}
	if service.webhookManager == nil || !service.webhookManager.IsEventConfigured(ProjectCreatedEvent) {
		return project, nil
	}

	err = service.webhookManager.InvokeWebhooks(ctx, ProjectCreatedEvent, project, func(p []byte) error {
		// Expects webhook output to be a project object
		var tmpproject models.Project
		if err := json.Unmarshal(p, &tmpproject); err != nil {
			return err
		}
		project, err = service.save(&tmpproject)
		if err != nil {
			return err
		}
		return nil
	}, webhooks.NoOpErrorHandler)
	if err != nil {
		return project,
			fmt.Errorf("error while invoking %s webhooks or on success callback function, err: %s",
				ProjectCreatedEvent, err.Error())
	}
	return project, nil
}

func (service *projectsService) ListProjects(ctx context.Context, name string, user string) (projects []*models.Project,
	err error) {
	allProjects, err := service.projectRepository.ListProjects(name)
	if err != nil {
		return nil, err
	}
	if service.authEnabled {
		return service.filterAuthorizedProjects(ctx, allProjects, user)
	}
	return allProjects, nil
}

func (service *projectsService) UpdateProject(ctx context.Context, project *models.Project) (*models.Project, string, error) {
	if service.authEnabled {
		err := service.updateAuthorizationPolicy(ctx, project)
		if err != nil {
			return nil, "", fmt.Errorf("error while updating authorization policy for project %s", project.Name)
		}
	}
	if service.webhookManager == nil || !service.webhookManager.IsEventConfigured(ProjectUpdatedEvent) {
		updatedProject, err := service.save(project)
	if err != nil {
		return nil, "", fmt.Errorf("error saving project %s: %w", project.Name, err)
	}

	payload, err := service.MakeRequestPayload(project, service.updateProjectPayloadTemplate)
	if err != nil {
		return nil, "", fmt.Errorf("error generating request payload: %w", err)
	}

	response, err := service.SendUpdateRequest(service.updateProjectEndpoint, payload)
	if err != nil {
		return nil, "", fmt.Errorf("error sending update request: %w", err)
	}
	defer response.Body.Close()

	responseMessage, err := service.ProcessResponseURL(response, service.updateProjectResponseTemplate)
	if err != nil {
		return nil, "", fmt.Errorf("error processing response template: %w", err)
	}

	return updatedProject, responseMessage, nil
	}
	err := service.webhookManager.InvokeWebhooks(ctx, ProjectUpdatedEvent, project, func(p []byte) error {
		// Expects webhook output to be a project object
		var tmpproject models.Project
		var err error
		if err := json.Unmarshal(p, &tmpproject); err != nil {
			return err
		}
		project, err = service.save(&tmpproject)
		if err != nil {
			return err
		}
		return nil
	}, webhooks.NoOpErrorHandler)
	if err != nil {
		return project,
			fmt.Errorf("error while invoking %s webhooks or on success callback function, err: %s",
				ProjectCreatedEvent, err.Error())
	}
	return project, nil
}

func (service *projectsService) FindByID(projectID models.ID) (*models.Project, error) {
	return service.projectRepository.Get(projectID)
}

func (service *projectsService) FindByName(projectName string) (*models.Project, error) {
	return service.projectRepository.GetByName(projectName)
}

func (service *projectsService) save(project *models.Project) (*models.Project, error) {
	if strings.TrimSpace(project.MLFlowTrackingURL) == "" {
		project.MLFlowTrackingURL = service.defaultMlflowTrackingServer
	}

	return service.projectRepository.Save(project)
}

func readPermissions(project *models.Project) []string {
	permissions := make([]string, 0)
	for _, method := range []string{"get"} {
		permissions = append(permissions, fmt.Sprintf("mlp.projects.%d.%s", project.ID, method))
	}
	return permissions
}

func adminPermissions(project *models.Project) []string {
	permissions := make([]string, 0)
	for _, method := range []string{"get", "put", "post", "patch", "delete"} {
		permissions = append(permissions, fmt.Sprintf("mlp.projects.%d.%s", project.ID, method))
	}
	return permissions
}

func (service *projectsService) updateAuthorizationPolicy(ctx context.Context, project *models.Project) error {
	updateRequest := enforcer.NewAuthorizationUpdateRequest()
	rolesWithReadOnlyAccess, err := enforcer.ParseProjectRoles([]string{
		enforcer.MLPProjectsReaderRole,
		enforcer.MLPProjectReaderRole,
	}, project)
	if err != nil {
		return err
	}
	for _, role := range rolesWithReadOnlyAccess {
		updateRequest.AddRolePermissions(role, readPermissions(project))
	}
	projectAdminRole, err := enforcer.ParseProjectRole(enforcer.MLPProjectAdminRole, project)
	if err != nil {
		return err
	}
	if project.Administrators != nil {
		updateRequest.SetRoleMembers(projectAdminRole, project.Administrators)
	} else {
		updateRequest.SetRoleMembers(projectAdminRole, []string{})
	}

	rolesWithAdminAccess, err := enforcer.ParseProjectRoles([]string{
		enforcer.MLPAdminRole,
		enforcer.MLPProjectAdminRole,
	}, project)
	if err != nil {
		return err
	}
	for _, role := range rolesWithAdminAccess {
		updateRequest.AddRolePermissions(role, adminPermissions(project))
	}
	projectReaderRole, err := enforcer.ParseProjectRole(enforcer.MLPProjectReaderRole, project)
	if err != nil {
		return err
	}
	if project.Readers != nil {
		updateRequest.SetRoleMembers(projectReaderRole, project.Readers)
	} else {
		updateRequest.SetRoleMembers(projectReaderRole, []string{})

	}

	return service.authEnforcer.UpdateAuthorization(ctx, updateRequest)
}

// TODO: Evaluate if we should retrieve all permissions granted to a user as opposed to just roles
func (service *projectsService) filterAuthorizedProjects(ctx context.Context, projects []*models.Project,
	user string) ([]*models.Project, error) {
	if user == "" {
		return nil, fmt.Errorf("authorization is enabled but user is not provided")
	}

	roles, err := service.authEnforcer.GetUserRoles(ctx, user)
	if err != nil {
		return nil, err
	}
	for _, role := range roles {
		if slices.Contains([]string{enforcer.MLPAdminRole, enforcer.MLPProjectsReaderRole}, role) {
			return projects, nil
		}
	}
	authorizedProjects := make([]*models.Project, 0)
	for _, project := range projects {
		if (project.Administrators != nil && slices.Contains(project.Administrators, user)) ||
			(project.Readers != nil && slices.Contains(project.Readers, user)) {
			authorizedProjects = append(authorizedProjects, project)
		}
	}
	return authorizedProjects, nil
}

func (service *projectsService) MakeRequestPayload(project *models.Project, templateString string) (string, error) {
	tmpl, err := template.New("requestPayload").Parse(templateString)
	if err != nil {
		return "", err
	}
	var payload bytes.Buffer
	if err := tmpl.Execute(&payload, project); err != nil {
		return "", err
	}
	return payload.String(), nil
}

func (service *projectsService) SendUpdateRequest(url string, payload string) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (service *projectsService) ProcessResponseURL(response *http.Response, templateString string) (string, error) {
	var data map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return "", err
	}

	tmpl, err := template.New("responseURL").Parse(templateString)
	if err != nil {
		return "", err
	}

	var responseText bytes.Buffer
	if err := tmpl.Execute(&responseText, data); err != nil {
		return "", err
	}
	return responseText.String(), nil
}
