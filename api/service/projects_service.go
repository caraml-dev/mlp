package service

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/exp/slices"

	"github.com/pkg/errors"

	"github.com/caraml-dev/mlp/api/repository"

	"github.com/caraml-dev/mlp/api/models"
	"github.com/caraml-dev/mlp/api/pkg/authz/enforcer"
)

type ProjectsService interface {
	ListProjects(ctx context.Context, name string, user string) ([]*models.Project, error)
	CreateProject(ctx context.Context, project *models.Project) (*models.Project, error)
	UpdateProject(ctx context.Context, project *models.Project) (*models.Project, error)
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
	authEnabled bool) (ProjectsService, error) {
	if strings.TrimSpace(mlflowURL) == "" {
		return nil, errors.New("default mlflow tracking url should be provided")
	}

	return &projectsService{
		projectRepository:           projectRepository,
		defaultMlflowTrackingServer: mlflowURL,
		authEnforcer:                authEnforcer,
		authEnabled:                 authEnabled,
	}, nil
}

type projectsService struct {
	projectRepository           repository.ProjectRepository
	defaultMlflowTrackingServer string
	authEnforcer                enforcer.Enforcer
	authEnabled                 bool
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

func (service *projectsService) UpdateProject(ctx context.Context, project *models.Project) (*models.Project, error) {
	if service.authEnabled {
		err := service.updateAuthorizationPolicy(ctx, project)
		if err != nil {
			return nil, fmt.Errorf("error while updating authorization policy for project %s", project.Name)
		}
	}

	return service.save(project)
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

func projectReaderRole(project *models.Project) string {
	return fmt.Sprintf("mlp.projects.%d.reader", project.ID)
}

func projectAdminRole(project *models.Project) string {
	return fmt.Sprintf("mlp.projects.%d.administrator", project.ID)
}

func rolesWithReadOnlyAccess(project *models.Project) []string {
	predefinedRoles := []string{
		"mlp.projects.reader",
	}
	return append(predefinedRoles, projectReaderRole(project))
}

func rolesWithAdminAccess(project *models.Project) []string {
	predefinedRoles := []string{
		"mlp.administrator",
	}
	return append(predefinedRoles, projectAdminRole(project))
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
	for _, role := range rolesWithReadOnlyAccess(project) {
		updateRequest.UpdateRolePermissions(role, readPermissions(project))
	}
	if project.Administrators != nil {
		updateRequest.UpdateRoleMembers(projectAdminRole(project), project.Administrators)
	} else {
		updateRequest.UpdateRoleMembers(projectAdminRole(project), []string{})
	}

	for _, role := range rolesWithAdminAccess(project) {
		updateRequest.UpdateRolePermissions(role, adminPermissions(project))
	}
	if project.Readers != nil {
		updateRequest.UpdateRoleMembers(projectReaderRole(project), project.Readers)
	} else {
		updateRequest.UpdateRoleMembers(projectReaderRole(project), []string{})
	}

	return service.authEnforcer.UpdateAuthorization(ctx, updateRequest)
}

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
		if slices.Contains([]string{"mlp.administrator", "mlp.projects.reader"}, role) {
			return projects, nil
		}
	}
	if err != nil {
		return nil, err
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
