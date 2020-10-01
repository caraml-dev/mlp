// +build integration integration_local

package storage

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"

	"github.com/gojek/mlp/api/it/database"
	"github.com/gojek/mlp/api/models"
)

func TestProjectsService_SaveAndGet(t *testing.T) {
	tests := []struct {
		name    string
		project models.Project
	}{
		{
			"project with label",
			models.Project{
				Name:           "test_project",
				Administrators: []string{"user@example.com"},
				Readers:        []string{"user-2@example.com"},
				Team:           "dsp",
				Stream:         "dsp",
				Labels: models.Labels{
					{
						Key:   "labelKey",
						Value: "labelValue",
					},
				},
			},
		},
		{
			"project without label",
			models.Project{
				Name:           "test_project_with_label",
				Administrators: []string{"user@example.com"},
				Readers:        []string{"user-2@example.com"},
				Team:           "dsp",
				Stream:         "dsp",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			database.WithTestDatabase(t, func(t *testing.T, db *gorm.DB) {
				projectStorage := NewProjectStorage(db)

				saved, err := projectStorage.Save(&tt.project)
				assert.NoError(t, err)

				assert.Equal(t, tt.project.Name, saved.Name)
				assert.Equal(t, tt.project.Administrators, saved.Administrators)
				assert.Equal(t, tt.project.Readers, saved.Readers)
				assert.NotNil(t, saved.CreatedAt)
				assert.NotNil(t, saved.UpdatedAt)

				res, err := projectStorage.Get(saved.Id)
				assert.NoError(t, err)
				assert.Equal(t, tt.project.Name, res.Name)
				assert.Equal(t, tt.project.Administrators, res.Administrators)
				assert.Equal(t, tt.project.Readers, res.Readers)
				assert.Equal(t, tt.project.Team, res.Team)
				assert.Equal(t, tt.project.Stream, res.Stream)
				assert.EqualValues(t, tt.project.Labels, res.Labels)

				res, err = projectStorage.GetByName(saved.Name)
				assert.NoError(t, err)
				assert.Equal(t, tt.project.Name, res.Name)
				assert.Equal(t, tt.project.Administrators, res.Administrators)
				assert.Equal(t, tt.project.Readers, res.Readers)
				assert.Equal(t, tt.project.Team, res.Team)
				assert.Equal(t, tt.project.Stream, res.Stream)
				assert.EqualValues(t, tt.project.Labels, res.Labels)
			})
		})
	}
}

func TestProjectsService_List(t *testing.T) {
	database.WithTestDatabase(t, func(t *testing.T, db *gorm.DB) {
		projectStorage := NewProjectStorage(db)

		projects := []models.Project{
			{
				Name:           "project_1",
				Administrators: []string{"user@example.com"},
				Readers:        []string{"user-2@example.com"},
			},
			{
				Name:           "project_2",
				Administrators: []string{"user@example.com"},
				Readers:        []string{"user-2@example.com"},
			},
			{
				Name:           "project_3",
				Administrators: []string{"user@example.com"},
				Readers:        []string{"user-2@example.com"},
			},
			{
				Name:           "my-project",
				Administrators: []string{"user@example.com"},
				Readers:        []string{"user-2@example.com"},
			},
		}

		for _, p := range projects {
			projectStorage.Save(&p)
		}

		res, err := projectStorage.ListProjects("")
		assert.NoError(t, err)
		assert.Len(t, res, len(projects))

		res, err = projectStorage.ListProjects("my-project")
		assert.NoError(t, err)
		assert.Len(t, res, 1)
		assert.Equal(t, "my-project", res[0].Name)

		res, err = projectStorage.ListProjects("unknown-project")
		assert.NoError(t, err)
		assert.Len(t, res, 0)
	})
}
