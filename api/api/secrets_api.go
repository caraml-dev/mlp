package api

import (
	"fmt"
	"net/http"

	"github.com/jinzhu/copier"

	"github.com/caraml-dev/mlp/api/log"
	"github.com/caraml-dev/mlp/api/models"
)

type SecretsController struct {
	*AppContext
}

func (c *SecretsController) GetSecret(r *http.Request, vars map[string]string, _ interface{}) *Response {
	projectID, _ := models.ParseID(vars["project_id"])
	secretID, _ := models.ParseID(vars["secret_id"])
	if projectID <= 0 || secretID <= 0 {
		log.Errorf("invalid id, secret_id: %d, project_id: %d", secretID, projectID)
		return BadRequest("project_id and secret_id are not valid")
	}

	secret, err := c.SecretService.FindByID(secretID)
	if err != nil {
		log.Errorf("error fetching secret with ID %d: %s", secretID, err)
		return FromError(err)
	}

	return Ok(secret)
}

func (c *SecretsController) CreateSecret(r *http.Request, vars map[string]string, body interface{}) *Response {
	projectID, _ := models.ParseID(vars["project_id"])
	_, err := c.ProjectsService.FindByID(projectID)
	if err != nil {
		log.Errorf("error fetching project with ID: %d", projectID)
		return NotFound(fmt.Sprintf("Project with given `project_id: %d` not found", projectID))
	}

	secret, ok := body.(*models.Secret)
	secret.ProjectID = projectID
	if !ok || !secret.IsValidForInsertion() {
		log.Errorf("invalid request body: %v", body)
		return BadRequest("Invalid request body")
	}

	// check that the secret storage id exists if users specify it
	if secret.SecretStorageID != nil {
		_, err = c.SecretStorageService.FindByID(*secret.SecretStorageID)
		if err != nil {
			log.Errorf("error fetching secret storage with ID %d: %s", *secret.SecretStorageID, err)
			return FromError(err)
		}
	}

	// creates
	secret, err = c.SecretService.Create(secret)
	if err != nil {
		log.Errorf("Failed creating new secret: %s", err)
		return FromError(err)
	}

	return Created(secret)
}

func (c *SecretsController) UpdateSecret(r *http.Request, vars map[string]string, body interface{}) *Response {
	updateRequest, ok := body.(*models.Secret)
	if !ok {
		return BadRequest("Invalid request body")
	}
	projectID, _ := models.ParseID(vars["project_id"])
	secretID, _ := models.ParseID(vars["secret_id"])
	if projectID <= 0 || secretID <= 0 {
		log.Errorf("invalid id, secret_id: %d, project_id: %d", secretID, projectID)
		return BadRequest("project_id and secret_id are not valid")
	}

	secret, err := c.SecretService.FindByID(secretID)
	if err != nil {
		return FromError(err)
	}

	// check that the secret storage id exists if users specify it
	if updateRequest.SecretStorageID != nil {
		_, err := c.SecretStorageService.FindByID(*updateRequest.SecretStorageID)
		if err != nil {
			log.Errorf("error fetching secret storage with ID %d: %s", *updateRequest.SecretStorageID, err)
			return FromError(err)
		}
	}

	err = copier.CopyWithOption(secret, updateRequest, copier.Option{IgnoreEmpty: true})
	if err != nil {
		log.Errorf("Failed copy secret with %s", err)
		return InternalServerError(err.Error())
	}

	updatedSecret, err := c.SecretService.Update(secret)
	if err != nil {
		log.Errorf("Failed update secret with %s", err)
		return FromError(err)
	}
	return Ok(updatedSecret)
}

func (c *SecretsController) DeleteSecret(r *http.Request, vars map[string]string, _ interface{}) *Response {
	projectID, _ := models.ParseID(vars["project_id"])
	secretID, _ := models.ParseID(vars["secret_id"])
	if projectID <= 0 || secretID <= 0 {
		log.Errorf("invalid id, secret_id: %d, project_id: %d", secretID, projectID)
		return BadRequest("project_id and secret_id are not valid")
	}

	if err := c.SecretService.Delete(secretID); err != nil {
		log.Errorf("error deleting secret with id %v", err)
		return InternalServerError(err.Error())
	}
	return NoContent()
}

func (c *SecretsController) ListSecret(r *http.Request, vars map[string]string, body interface{}) *Response {
	projectID, _ := models.ParseID(vars["project_id"])
	_, err := c.ProjectsService.FindByID(projectID)
	if err != nil {
		log.Errorf("error fetching project with ID %d: %s", projectID, err)
		return FromError(err)
	}

	secrets, err := c.SecretService.List(projectID)
	if err != nil {
		log.Errorf("error retrieving secret from project id %s: %s", projectID, err)
		return FromError(err)
	}
	return Ok(secrets)
}

func (c *SecretsController) Routes() []Route {
	return []Route{
		{
			http.MethodGet,
			"/projects/{project_id:[0-9]+}/secrets",
			nil,
			c.ListSecret,
			"ListSecret",
		},
		{
			http.MethodGet,
			"/projects/{project_id:[0-9]+}/secrets/{secret_id:[0-9]+}",
			nil,
			c.GetSecret,
			"GetSecret",
		},
		{
			http.MethodPost,
			"/projects/{project_id:[0-9]+}/secrets",
			models.Secret{},
			c.CreateSecret,
			"CreateSecret",
		},
		{
			http.MethodPatch,
			"/projects/{project_id:[0-9]+}/secrets/{secret_id}",
			models.Secret{},
			c.UpdateSecret,
			"UpdateSecret",
		},
		{
			http.MethodDelete,
			"/projects/{project_id:[0-9]+}/secrets/{secret_id}",
			nil,
			c.DeleteSecret,
			"DeleteSecret",
		},
	}
}
