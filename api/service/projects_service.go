package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"bytes"
	"html/template"
	"net/http"

	"golang.org/x/exp/slices"

	"github.com/pkg/errors"

	"github.com/caraml-dev/mlp/api/repository"

	"github.com/caraml-dev/mlp/api/config"
	"github.com/caraml-dev/mlp/api/models"
	"github.com/caraml-dev/mlp/api/pkg/authz/enforcer"
	"github.com/caraml-dev/mlp/api/pkg/webhooks"
)

type ProjectsService interface {
	ListProjects(ctx context.Context, name string, user string) ([]*models.Project, error)
	CreateProject(ctx context.Context, project *models.Project) (*models.Project, error)
	UpdateProject(ctx context.Context, project *models.Project) (*models.Project, map[string]interface{}, error)
	FindByID(projectID models.ID) (*models.Project, error)
	FindByName(projectName string) (*models.Project, error)
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
	updateProjectConfig config.UpdateProjectConfig) (ProjectsService, error) {
	if strings.TrimSpace(mlflowURL) == "" {
		return nil, errors.New("default mlflow tracking url should be provided")
	}

	return &projectsService{
		projectRepository:           projectRepository,
		defaultMlflowTrackingServer: mlflowURL,
		authEnforcer:                authEnforcer,
		authEnabled:                 authEnabled,
		webhookManager:              webhookManager,
		updateProjectConfig:         updateProjectConfig,
	}, nil
}

type projectsService struct {
	projectRepository           repository.ProjectRepository
	defaultMlflowTrackingServer string
	authEnforcer                enforcer.Enforcer
	authEnabled                 bool
	webhookManager              webhooks.WebhookManager
	updateProjectConfig         config.UpdateProjectConfig
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

func (service *projectsService) UpdateProject(ctx context.Context, project *models.Project) (*models.Project,
	map[string]interface{}, error) {
	if service.authEnabled {
		err := service.updateAuthorizationPolicy(ctx, project)
		if err != nil {
			return nil, nil, fmt.Errorf("error while updating authorization policy for project %s", project.Name)
		}
	}

	areBlacklistedLabelsChanged, err := service.areBlacklistedLabelsChanged(project.ID, project.Labels,
		service.updateProjectConfig.LabelsBlacklist)
	if err != nil {
		return nil, nil, err
	}
	if areBlacklistedLabelsChanged {
		return nil, nil,
			fmt.Errorf("one or more labels are blacklisted or have been removed or changed values and cannot be updated")
	}

	if service.webhookManager != nil && service.webhookManager.IsEventConfigured(ProjectUpdatedEvent) {
		err = service.webhookManager.InvokeWebhooks(ctx, ProjectUpdatedEvent, project, func(p []byte) error {
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
			return project, nil,
				fmt.Errorf("error while invoking %s webhooks or on success callback function, err: %s",
					ProjectUpdatedEvent, err.Error())
		}
	} else {
		project, err = service.save(project)
		if err != nil {
			return nil, nil, err
		}
	}

	project, response, err := service.handleUpdateProjectRequest(project)
	if err != nil {
		return nil, nil, err
	}

	return project, response, nil
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

func (service *projectsService) handleUpdateProjectRequest(project *models.Project) (*models.Project,
	map[string]interface{}, error) {
	if service.updateProjectConfig.Endpoint == "" {
		return project, nil, nil
	}

	if service.updateProjectConfig.PayloadTemplate == "" || service.updateProjectConfig.ResponseTemplate == "" {
		return project, nil, nil
	}

	payload, err := generateRequestPayload(project, service.updateProjectConfig.PayloadTemplate)
	if err != nil {
		return nil, nil, fmt.Errorf("error generating request payload: %w", err)
	}

	resp, err := sendUpdateRequest(service.updateProjectConfig.Endpoint, payload)
	if err != nil {
		return nil, nil, fmt.Errorf("error sending update request: %w", err)
	}
	defer resp.Body.Close()

	response, err := processResponseTemplate(resp, service.updateProjectConfig.ResponseTemplate)
	if err != nil {
		return nil, nil, fmt.Errorf("error processing response template: %w", err)
	}

	return project, response, nil
}

func generateRequestPayload(project *models.Project, templateString string) (map[string]interface{}, error) {
	tmpl, err := template.New("requestPayload").Parse(templateString)
	if err != nil {
		return nil, err
	}
	var payload bytes.Buffer
	if err := tmpl.Execute(&payload, project); err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := json.Unmarshal(payload.Bytes(), &result); err != nil {
		return nil, err
	}

	return result, nil
}

func sendUpdateRequest(url string, payload map[string]interface{}) (*http.Response, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
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

func processResponseTemplate(response *http.Response, templateString string) (map[string]interface{}, error) {
	var data map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return nil, err
	}

	tmpl, err := template.New("responsePayload").Parse(templateString)
	if err != nil {
		return nil, err
	}

	var responseText bytes.Buffer
	if err := tmpl.Execute(&responseText, data); err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(responseText.Bytes(), &result); err != nil {
		return nil, err
	}

	return result, nil
}

// areBlacklistedLabelsChanged check if any key in labels is blacklisted
func (service *projectsService) areBlacklistedLabelsChanged(projectID models.ID, newLabels []models.Label, blacklist map[string]bool) (bool, error) {
	existingProject, err := service.FindByID(projectID)
	if err != nil {
		return false, fmt.Errorf("error fetching project with id %s: %w", projectID, err)
	}

	for _, existingLabel := range existingProject.Labels {
		if blacklist[existingLabel.Key] {
			found := false
			for _, newLabel := range newLabels {
				if newLabel.Key == existingLabel.Key {
					found = true
					if newLabel.Value != existingLabel.Value {
						return true, nil
					}
				}
			}
			if !found {
				return true, nil
			}
		}
	}

	return false, nil
}
