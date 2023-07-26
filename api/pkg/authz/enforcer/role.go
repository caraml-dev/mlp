package enforcer

import (
	"bytes"
	"text/template"

	"github.com/caraml-dev/mlp/api/models"
)

const (
	MLPAdminRole          = "mlp.administrator"
	MLPProjectsReaderRole = "mlp.projects.reader"
	MLPProjectReaderRole  = "mlp.projects.{{ .ProjectId }}.reader"
	MLPProjectAdminRole   = "mlp.projects.{{ .ProjectId }}.administrator"
)

func ParseRole(role string, templateContext map[string]string) (string, error) {
	roleParser, err := template.New("role").Parse(role)
	if err != nil {
		return "", err
	}
	var parseResultBytes bytes.Buffer
	err = roleParser.Execute(&parseResultBytes, templateContext)
	if err != nil {
		return "", err
	}
	return parseResultBytes.String(), nil
}

func ParseProjectRole(roleTemplateString string, project *models.Project) (string, error) {
	parsedRole, err := ParseRole(roleTemplateString, map[string]string{"ProjectId": project.ID.String()})
	if err != nil {
		return "", err
	}
	return parsedRole, nil
}

func ParseProjectRoles(roleTemplateStrings []string, project *models.Project) ([]string, error) {
	roles := make([]string, len(roleTemplateStrings))
	for i, roleTemplateString := range roleTemplateStrings {
		parsedRole, err := ParseProjectRole(roleTemplateString, project)
		roles[i] = parsedRole
		if err != nil {
			return nil, err
		}
	}
	return roles, nil
}
