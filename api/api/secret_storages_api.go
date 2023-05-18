package api

import (
	"net/http"

	"github.com/caraml-dev/mlp/api/log"
	"github.com/caraml-dev/mlp/api/models"
)

type SecretStoragesController struct {
	*AppContext
}

// ListSecretStorage lists all secret storage for a project including the global secret storage
func (c *SecretStoragesController) ListSecretStorage(_ *http.Request,
	vars map[string]string,
	_ interface{}) *Response {

	projectID, _ := models.ParseID(vars["project_id"])
	if projectID <= 0 {
		log.Errorf("project_id is not valid: %d", projectID)
		return BadRequest("project_id is not valid")
	}

	_, err := c.ProjectsService.FindByID(projectID)
	if err != nil {
		log.Errorf("error fetching project with ID: %d", projectID)
		return FromError(err)
	}

	secretStorages, err := c.SecretStorageService.List(projectID)
	if err != nil {
		log.Errorf("error fetching secret storages for project with ID: %d", projectID)
		return FromError(err)
	}

	return Ok(secretStorages)
}

// GetSecretStorage gets a secret storage given specified project_id and secret_storage_id
func (c *SecretStoragesController) GetSecretStorage(_ *http.Request,
	vars map[string]string,
	_ interface{}) *Response {

	projectID, _ := models.ParseID(vars["project_id"])
	secretStorageID, _ := models.ParseID(vars["secret_storage_id"])
	if projectID <= 0 || secretStorageID <= 0 {
		log.Errorf("invalid id, secret_storage_id: %d, project_id: %d", secretStorageID, projectID)
		return BadRequest("project_id and secret_storage_id are not valid")
	}

	secretStorage, err := c.SecretStorageService.FindByID(secretStorageID)
	if err != nil {
		log.Errorf("error fetching secret storage with ID: %d", secretStorageID)
		return FromError(err)
	}

	return Ok(secretStorage)
}

// CreateSecretStorage creates a secret storage for a project
func (c *SecretStoragesController) CreateSecretStorage(_ *http.Request,
	vars map[string]string,
	body interface{}) *Response {

	projectID, _ := models.ParseID(vars["project_id"])
	project, err := c.ProjectsService.FindByID(projectID)
	if err != nil {
		log.Errorf("error fetching project with ID: %d", projectID)
		return FromError(err)
	}

	secretStorage, ok := body.(*models.SecretStorage)
	if !ok {
		log.Errorf("invalid request body: %v", body)
		return BadRequest("Invalid body")
	}
	secretStorage.ProjectID = &projectID
	secretStorage.Project = project

	err = secretStorage.ValidateForCreation()
	if err != nil {
		log.Errorf("invalid secret storage request: %s", err)
		return BadRequest(err.Error())
	}

	secretStorage, err = c.SecretStorageService.Create(secretStorage)
	if err != nil {
		log.Errorf("error creating secret storage: %s", err)
		return FromError(err)
	}

	return Created(secretStorage)
}

// UpdateSecretStorage updates a secret storage given specified project_id and secret_storage_id
// Note: cannot update global secret storage
func (c *SecretStoragesController) UpdateSecretStorage(_ *http.Request,
	vars map[string]string,
	body interface{}) *Response {

	projectID, _ := models.ParseID(vars["project_id"])
	secretStorageID, _ := models.ParseID(vars["secret_storage_id"])
	if projectID <= 0 || secretStorageID <= 0 {
		log.Errorf("invalid id, secret_storage_id: %d, project_id: %d", secretStorageID, projectID)
		return BadRequest("project_id and secret_storage_id are not valid")
	}

	secretStorage, err := c.SecretStorageService.FindByID(secretStorageID)
	if err != nil {
		log.Errorf("error fetching secret storage with ID: %d", secretStorageID)
		return FromError(err)
	}

	if secretStorage.Scope == models.GlobalSecretStorageScope {
		log.Errorf("cannot update global secret storage")
		return BadRequest("cannot update global secret storage")
	}

	updateRequest, ok := body.(*models.SecretStorage)
	if !ok {
		log.Errorf("invalid request body: %v", body)
		return BadRequest("invalid body")
	}

	err = secretStorage.MergeValue(updateRequest)
	if err != nil {
		log.Errorf("error merging secret storage: %s", err)
		return InternalServerError(err.Error())
	}

	err = secretStorage.ValidateForMutation()
	if err != nil {
		log.Errorf("invalid secret storage request: %s", err)
		return BadRequest(err.Error())
	}

	secretStorage, err = c.SecretStorageService.Update(secretStorage)
	if err != nil {
		log.Errorf("error updating secret storage: %s", err)
		return FromError(err)
	}

	return Ok(secretStorage)
}

func (c *SecretStoragesController) DeleteSecretStorage(_ *http.Request,
	vars map[string]string,
	_ interface{}) *Response {

	projectID, _ := models.ParseID(vars["project_id"])
	secretStorageID, _ := models.ParseID(vars["secret_storage_id"])
	if projectID <= 0 || secretStorageID <= 0 {
		log.Errorf("invalid id, secret_storage_id: %d, project_id: %d", secretStorageID, projectID)
		return BadRequest("project_id and secret_id are not valid")
	}

	err := c.SecretStorageService.Delete(secretStorageID)
	if err != nil {
		log.Errorf("error deleting secret storage with ID: %d", secretStorageID)
		return FromError(err)
	}

	return NoContent()
}

func (c *SecretStoragesController) Routes() []Route {
	return []Route{
		{
			http.MethodGet,
			"/projects/{project_id:[0-9]+}/secret_storages",
			nil,
			c.ListSecretStorage,
			"ListSecretStorage",
		},
		{
			http.MethodGet,
			"/projects/{project_id:[0-9]+}/secret_storages/{secret_storage_id:[0-9]+}",
			nil,
			c.GetSecretStorage,
			"GetSecretStorage",
		},
		{
			http.MethodPost,
			"/projects/{project_id:[0-9]+}/secret_storages",
			models.SecretStorage{},
			c.CreateSecretStorage,
			"CreateSecretStorage",
		},
		{
			http.MethodPatch,
			"/projects/{project_id:[0-9]+}/secret_storages/{secret_storage_id:[0-9]+}",
			models.SecretStorage{},
			c.UpdateSecretStorage,
			"UpdateSecretStorage",
		},
		{
			http.MethodDelete,
			"/projects/{project_id:[0-9]+}/secret_storages/{secret_storage_id:[0-9]+}",
			nil,
			c.DeleteSecretStorage,
			"DeleteSecretStorage",
		},
	}
}
