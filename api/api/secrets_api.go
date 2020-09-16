package api

import (
	"fmt"
	"net/http"

	"github.com/gojek/mlp/log"
	"github.com/gojek/mlp/models"
)

type SecretsController struct {
	*AppContext
}

func (c *SecretsController) CreateSecret(r *http.Request, vars map[string]string, body interface{}) *ApiResponse {
	projectId, _ := models.ParseId(vars["project_id"])
	_, err := c.ProjectsService.FindById(projectId)
	if err != nil {
		log.Warnf("Project with id: %d not found", projectId)
		return NotFound(fmt.Sprintf("Project with given `project_id: %d` not found", projectId))
	}

	secret, ok := body.(*models.Secret)
	secret.ProjectId = projectId
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

func (c *SecretsController) UpdateSecret(r *http.Request, vars map[string]string, body interface{}) *ApiResponse {
	secret, ok := body.(*models.Secret)
	if !ok {
		return BadRequest("Invalid request body")
	}
	projectId, _ := models.ParseId(vars["project_id"])
	secretId, _ := models.ParseId(vars["secret_id"])
	if projectId <= 0 || secretId <= 0 {
		log.Warnf("Id: %d or project_id not valid", secretId, projectId)
		return BadRequest("project_id and secret_id is not valid")
	}
	existingSecret, err := c.SecretService.FindByIdAndProjectId(secretId, projectId)
	if err != nil {
		log.Errorf("Unable to find secret with %v", err)
		return InternalServerError(err.Error())
	}
	if existingSecret == nil {
		log.Warnf("Secret with id: %d and project_id not found", secretId, projectId)
		return NotFound(fmt.Sprintf("Secret with given `secret_id: %d` and `project_id: %d` not found", secretId, projectId))
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

func (c *SecretsController) DeleteSecret(r *http.Request, vars map[string]string, _ interface{}) *ApiResponse {
	projectId, _ := models.ParseId(vars["project_id"])
	secretId, _ := models.ParseId(vars["secret_id"])
	if projectId <= 0 || secretId <= 0 {
		log.Warnf("Id: %d or project_id not valid", secretId, projectId)
		return BadRequest("project_id and secret_id is not valid")
	}

	if err := c.SecretService.Delete(secretId, projectId); err != nil {
		log.Errorf("Failed delete secret with %v", err)
		return InternalServerError(err.Error())
	}
	return NoContent()
}

func (c *SecretsController) ListSecret(r *http.Request, vars map[string]string, body interface{}) *ApiResponse {
	projectId, _ := models.ParseId(vars["project_id"])
	_, err := c.ProjectsService.FindById(projectId)
	if err != nil {
		log.Warnf("Project with id: %d not found", projectId)
		return NotFound(fmt.Sprintf("Project with given `project_id: %d` not found", projectId))
	}

	secrets, err := c.SecretService.ListSecret(projectId)
	if err != nil {
		log.Errorf("Failed retrieving secret from project id %s: %v", projectId, err)
		return InternalServerError(err.Error())
	}
	return Ok(secrets)
}
