/*
 * MLP API
 *
 * API Guide for accessing MLP API
 *
 * API version: 0.4.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package client

import (
	"time"
)

type Secret struct {
	ID              int32     `json:"id,omitempty"`
	Name            string    `json:"name"`
	Data            string    `json:"data"`
	SecretStorageID int32     `json:"secret_storage_id,omitempty"`
	CreatedAt       time.Time `json:"created_at,omitempty"`
	UpdatedAt       time.Time `json:"updated_at,omitempty"`
}
