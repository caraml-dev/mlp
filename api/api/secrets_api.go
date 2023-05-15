package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/caraml-dev/mlp/api/log"
	"github.com/caraml-dev/mlp/api/models"
	"github.com/jinzhu/gorm"
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

	// check that the secret storage id exists if users specify it
	if secret.SecretStorageID != nil {
		_, err = c.SecretStorageService.FindByID(*secret.SecretStorageID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return NotFound(fmt.Sprintf("Secret storage with given `secret_storage_id: %d` not found", *secret.SecretStorageID))
			}

			return InternalServerError(err.Error())
		}
	}

	// creates
	secret, err = c.SecretService.Create(secret)
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

	_, err := c.SecretService.FindByID(secretID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return NotFound(fmt.Sprintf("Secret with given `secret_id: %d` not found", secretID))
		}

		return InternalServerError(err.Error())
	}

	// check that the secret storage id exists if users specify it
	if secret.SecretStorageID != nil {
		_, err := c.SecretStorageService.FindByID(*secret.SecretStorageID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return NotFound(fmt.Sprintf("Secret storage with given `secret_storage_id: %d` not found", *secret.SecretStorageID))
			}

			return InternalServerError(err.Error())
		}
	}

	updatedSecret, err := c.SecretService.Update(secret)
	if err != nil {
		log.Errorf("Failed update secret with %v", err)
		return InternalServerError(err.Error())
	}
	return Ok(updatedSecret)
}

func (c *SecretsController) DeleteSecret(r *http.Request, vars map[string]string, _ interface{}) *Response {
	projectID, _ := models.ParseID(vars["project_id"])
	secretID, _ := models.ParseID(vars["secret_id"])
	if projectID <= 0 || secretID <= 0 {
		log.Warnf("ID: %d or project_id not valid", secretID, projectID)
		return BadRequest("project_id and secret_id is not valid")
	}

	if err := c.SecretService.Delete(secretID); err != nil {
		log.Errorf("Failed delete secret with id %v", err)
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

	secrets, err := c.SecretService.List(projectID)
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
