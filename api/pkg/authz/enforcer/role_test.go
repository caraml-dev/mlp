package enforcer

import (
	"testing"

	"github.com/caraml-dev/mlp/api/models"
)

func TestParseRole(t *testing.T) {
	type args struct {
		role            string
		templateContext map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"parse role with project id",
			args{
				role:            MLPProjectReaderRole,
				templateContext: map[string]string{"ProjectId": "1"},
			},
			"mlp.projects.1.reader",
			false,
		},
		{
			"parse plain string role without template context",
			args{
				role:            MLPAdminRole,
				templateContext: nil,
			},
			"mlp.administrator",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRole(tt.args.role, tt.args.templateContext)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseRole() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseProjectRole(t *testing.T) {
	type args struct {
		role    string
		project *models.Project
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"parse role with project id",
			args{
				role: MLPProjectReaderRole,
				project: &models.Project{
					ID: 1,
				},
			},
			"mlp.projects.1.reader",
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseProjectRole(tt.args.role, tt.args.project)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseProjectRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseProjectRole() = %v, want %v", got, tt.want)
			}
		})
	}
}
