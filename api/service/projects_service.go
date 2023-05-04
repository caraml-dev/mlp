package service

import (
	"fmt"
	"strings"

	"github.com/caraml-dev/mlp/api/repository"
	"github.com/pkg/errors"

	"github.com/caraml-dev/mlp/api/models"
	"github.com/caraml-dev/mlp/api/pkg/authz/enforcer"
)

type ProjectsService interface {
	ListProjects(name string) ([]*models.Project, error)
	CreateProject(project *models.Project) (*models.Project, error)
	UpdateProject(project *models.Project) (*models.Project, error)
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

const (
	ProjectSubResources = "projects:%s:**"
	ProjectResources    = "projects:%s"
)

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

func (service *projectsService) CreateProject(project *models.Project) (*models.Project, error) {
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
		err = service.upsertAuthorizationPolicy(project)
		if err != nil {
			return nil, fmt.Errorf("error while creating authorization policy for project %s", project.Name)
		}
	}

	return project, nil
}

func (service *projectsService) ListProjects(name string) (projects []*models.Project, err error) {
	return service.projectRepository.ListProjects(name)
}

func (service *projectsService) UpdateProject(project *models.Project) (*models.Project, error) {
	if service.authEnabled {
		err := service.upsertAuthorizationPolicy(project)
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

func (service *projectsService) upsertAuthorizationPolicy(project *models.Project) error {
	// create administrators policy
	adminRole, err := service.upsertAdministratorsRole(project)
	if err != nil {
		return err
	}
	err = service.upsertAdministratorsPolicy(adminRole, project)
	if err != nil {
		return err
	}

	// create readers policy
	readersRole, err := service.upsertReadersRole(project)
	if err != nil {
		return err
	}
	err = service.upsertReadersPolicy(readersRole, project)
	if err != nil {
		return err
	}

	return nil
}

func (service *projectsService) upsertReadersRole(project *models.Project) (string, error) {
	roleName := fmt.Sprintf("%s-%s", project.Name, "readers")
	role, err := service.authEnforcer.UpsertRole(roleName, project.Readers)
	if err != nil {
		return "", err
	}
	return role.ID, nil
}

func (service *projectsService) upsertAdministratorsRole(project *models.Project) (string, error) {
	roleName := fmt.Sprintf("%s-%s", project.Name, "administrators")
	policy, err := service.authEnforcer.UpsertRole(roleName, project.Administrators)
	if err != nil {
		return "", err
	}
	return policy.ID, nil
}

func (service *projectsService) upsertAdministratorsPolicy(role string, project *models.Project) error {
	subResources := fmt.Sprintf(ProjectSubResources, project.ID)
	resource := fmt.Sprintf(ProjectResources, project.ID)
	nameResource := fmt.Sprintf(ProjectResources, project.Name)
	policyName := fmt.Sprintf("%s-administrators-policy", project.Name)
	_, err := service.authEnforcer.UpsertPolicy(
		policyName,
		[]string{role},
		[]string{},
		[]string{resource, subResources, nameResource},
		[]string{enforcer.ActionAll})
	return err
}

func (service *projectsService) upsertReadersPolicy(role string, project *models.Project) error {
	subResources := fmt.Sprintf(ProjectSubResources, project.ID)
	resource := fmt.Sprintf(ProjectResources, project.ID)
	nameResource := fmt.Sprintf(ProjectResources, project.Name)
	policyName := fmt.Sprintf("%s-readers-policy", project.Name)
	_, err := service.authEnforcer.UpsertPolicy(
		policyName,
		[]string{role},
		[]string{},
		[]string{resource, subResources, nameResource},
		[]string{enforcer.ActionRead})
	return err
}
