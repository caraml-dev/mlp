package api

import (
	"fmt"
	"net/http"

	"github.com/gojek/mlp/api/log"
	"github.com/gojek/mlp/api/models"
)

type SecretsController struct {
	*AppContext
}

func (c *SecretsController) CreateSecret(r *http.Request, vars map[string]string, body interface{}) *Response {
	projectID, _ := models.ParseID(vars["project_id"])
	_, err := c.ProjectsService.FindByID(projectID)
	if err != nil {
		log.Warnf("Project with id: %d not found", projectID)
		return NotFound(fmt.Sprintf("Project with given `project_id: %d` not found", projectID))
	}

	secret, ok := body.(*models.Secret)
	secret.ProjectID = projectID
	if !ok || !secret.IsValidForInsertion() {
		return BadRequest("Invalid request body")
	}

	secret, err = c.SecretService.Save(secret)
	if err != nil {
		log.Errorf("Failed create new secret with %v", err)
		return InternalServerError(err.Error())
	}

	return Created(secret)
}

func (c *SecretsController) UpdateSecret(r *http.Request, vars map[string]string, body interface{}) *Response {
	secret, ok := body.(*models.Secret)
	if !ok {
		return BadRequest("Invalid request body")
	}
	projectID, _ := models.ParseID(vars["project_id"])
	secretID, _ := models.ParseID(vars["secret_id"])
	if projectID <= 0 || secretID <= 0 {
		log.Warnf("ID: %d or project_id not valid", secretID, projectID)
		return BadRequest("project_id and secret_id is not valid")
	}
	existingSecret, err := c.SecretService.FindByIDAndProjectID(secretID, projectID)
	if err != nil {
		log.Errorf("Unable to find secret with %v", err)
		return InternalServerError(err.Error())
	}
	if existingSecret == nil {
		log.Warnf("Secret with id: %d and project_id not found", secretID, projectID)
		return NotFound(fmt.Sprintf("Secret with given `secret_id: %d` and `project_id: %d` not found", secretID, projectID))
	}

	existingSecret.CopyValueFrom(secret)
	if !existingSecret.IsValidForMutation() {
		log.Warnf("Unable to update secret because secret: %v is not valid", existingSecret)
		return BadRequest("Invalid request body")
	}

	existingSecret, err = c.SecretService.Save(existingSecret)
	if err != nil {
		log.Errorf("Failed update secret with %v", err)
		return InternalServerError(err.Error())
	}
	return Ok(existingSecret)
}

func (c *SecretsController) DeleteSecret(r *http.Request, vars map[string]string, _ interface{}) *Response {
	projectID, _ := models.ParseID(vars["project_id"])
	secretID, _ := models.ParseID(vars["secret_id"])
	if projectID <= 0 || secretID <= 0 {
		log.Warnf("ID: %d or project_id not valid", secretID, projectID)
		return BadRequest("project_id and secret_id is not valid")
	}

	if err := c.SecretService.Delete(secretID, projectID); err != nil {
		log.Errorf("Failed delete secret with %v", err)
		return InternalServerError(err.Error())
	}
	return NoContent()
}

func (c *SecretsController) ListSecret(r *http.Request, vars map[string]string, body interface{}) *Response {
	projectID, _ := models.ParseID(vars["project_id"])
	_, err := c.ProjectsService.FindByID(projectID)
	if err != nil {
		log.Warnf("Project with id: %d not found", projectID)
		return NotFound(fmt.Sprintf("Project with given `project_id: %d` not found", projectID))
	}

	secrets, err := c.SecretService.ListSecret(projectID)
	if err != nil {
		log.Errorf("Failed retrieving secret from project id %s: %v", projectID, err)
		return InternalServerError(err.Error())
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
