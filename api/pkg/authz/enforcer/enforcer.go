package enforcer

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/ory/keto-client-go/client"
	"github.com/ory/keto-client-go/client/engines"
	"github.com/ory/keto-client-go/models"

	"github.com/gojek/mlp/pkg/authz/enforcer/types"
)

const (
	// ActionCreate action to create a resource
	ActionCreate = "actions:create"
	// ActionRead action to read a resource
	ActionRead = "actions:read"
	// ActionUpdate action to update a resource
	ActionUpdate = "actions:update"
	// ActionDelete action to delete a resource
	ActionDelete = "actions:delete"
	// ActionAll all action
	ActionAll = "actions:**"
)

// Flavor flavor type
type Flavor string

const (
	// FlavorExact keto flavor using "exact" semantics
	FlavorExact Flavor = "exact"
	// FlavorGlob keto flavor using "glob" pattern matching
	FlavorGlob Flavor = "glob"
	// FlavorRegex keto flavor using "regex" pattern matching
	FlavorRegex Flavor = "regex"
)

// Enforcer thin client providing interface for authorizing users
type Enforcer interface {
	// Enforce check whether user is authorized to do certain action against a resource
	Enforce(user string, resource string, action string) (*bool, error)
	// FilterAuthorizedResource filter and return list of authorized resource for certain user
	FilterAuthorizedResource(user string, resources []string, action string) ([]string, error)
	// GetRole get role with name
	GetRole(roleName string) (*types.Role, error)
	// GetPolicy get policy with name
	GetPolicy(policyName string) (*types.Policy, error)
	// UpsertRole create or update a role containing member as specified by users argument
	UpsertRole(roleName string, users []string) (*types.Role, error)
	// UpsertPolicy create or update a policy to allow subjects do actions against the specified resources
	UpsertPolicy(policyName string, roles []string, users []string, resources []string, actions []string) (*types.Policy, error)
}

type enforcer struct {
	ketoClient *engines.Client
	product    string
	flavor     Flavor
	timeout    time.Duration
}

func newEnforcer(hostURL string, productName string, flavor Flavor, timeout time.Duration) (Enforcer, error) {
	u, err := url.ParseRequestURI(hostURL)
	if err != nil {
		return nil, err
	}
	client := client.NewHTTPClientWithConfig(nil, &client.TransportConfig{
		Host:     u.Host,
		BasePath: u.Path,
		Schemes:  []string{u.Scheme},
	})
	return &enforcer{
		ketoClient: client.Engines,
		product:    productName,
		flavor:     flavor,
		timeout:    timeout,
	}, nil
}

// Enforce check whether user is authorized to do action against a resource
func (e *enforcer) Enforce(user string, resource string, action string) (*bool, error) {
	user = e.formatUser(user)
	resource = e.formatResource(resource)

	return e.isAllowed(user, resource, action)
}

// GetRole get role with name
func (e *enforcer) GetRole(roleName string) (*types.Role, error) {
	fmtRole := e.formatRole(roleName)

	params := &engines.GetOryAccessControlPolicyRoleParams{
		Flavor: string(e.flavor),
		ID:     fmtRole,
	}
	res, err := e.ketoClient.GetOryAccessControlPolicyRole(params.WithTimeout(e.timeout))
	if err != nil {
		return nil, err
	}
	return &types.Role{
		ID:      res.GetPayload().ID,
		Members: res.GetPayload().Members,
	}, nil
}

// GetPolicy get policy with name
func (e *enforcer) GetPolicy(policyName string) (*types.Policy, error) {
	fmtPolicy := e.formatPolicy(policyName)
	params := &engines.GetOryAccessControlPolicyParams{
		Flavor: string(e.flavor),
		ID:     fmtPolicy,
	}
	res, err := e.ketoClient.GetOryAccessControlPolicy(params.WithTimeout(e.timeout))
	if err != nil {
		return nil, err
	}
	payload := res.GetPayload()
	return &types.Policy{
		ID:        payload.ID,
		Actions:   payload.Actions,
		Resources: payload.Resources,
		Subjects:  payload.Subjects,
	}, nil
}

// FilterAuthorizedResource filter and return list of authorized resource for certain user
func (e *enforcer) FilterAuthorizedResource(user string, resources []string, action string) ([]string, error) {
	user = e.formatUser(user)

	allowedResources := make([]string, 0, 0)
	for _, resource := range resources {
		resource = e.formatResource(resource)

		allowed, err := e.isAllowed(user, resource, action)
		if err != nil {
			return nil, err
		}

		if *allowed {
			allowedResources = append(allowedResources, e.stripResourcePrefix(resource))
		}
	}

	return allowedResources, nil
}

// UpsertRole create or update a role containing member as specified by users argument
func (e *enforcer) UpsertRole(roleName string, users []string) (*types.Role, error) {
	fmtRoleName := e.formatRole(roleName)
	fmtUser := make([]string, 0, 0)
	for _, user := range users {
		fmtUser = append(fmtUser, e.formatUser(user))
	}

	input := &models.OryAccessControlPolicyRole{
		ID:      fmtRoleName,
		Members: fmtUser,
	}
	params := &engines.UpsertOryAccessControlPolicyRoleParams{
		Body:   input,
		Flavor: string(e.flavor),
	}
	res, err := e.ketoClient.UpsertOryAccessControlPolicyRole(params.WithTimeout(e.timeout))
	if err != nil {
		return nil, err
	}
	return &types.Role{
		ID:      res.GetPayload().ID,
		Members: res.GetPayload().Members,
	}, nil
}

// CreatePolicy create a policy to allow subject do an operation against the specified resource
func (e *enforcer) UpsertPolicy(policyName string, roles []string, users []string, resources []string, actions []string) (*types.Policy, error) {
	fmtPolicy := e.formatPolicy(policyName)

	fmtResources := make([]string, 0, 0)
	for _, res := range resources {
		fmtResources = append(fmtResources, e.formatResource(res))
	}

	fmtRoles := make([]string, 0, 0)
	for _, role := range roles {
		fmtRoles = append(fmtRoles, e.formatRole(role))
	}

	fmtUsers := make([]string, 0, 0)
	for _, user := range users {
		fmtUsers = append(fmtUsers, e.formatUser(user))
	}

	input := &models.OryAccessControlPolicy{
		Actions:   actions,
		Effect:    "allow",
		ID:        fmtPolicy,
		Resources: fmtResources,
		Subjects:  append(fmtRoles, fmtUsers...),
	}
	params := &engines.UpsertOryAccessControlPolicyParams{
		Body:   input,
		Flavor: string(e.flavor),
	}
	res, err := e.ketoClient.UpsertOryAccessControlPolicy(params.WithTimeout(e.timeout))
	if err != nil {
		return nil, err
	}

	payload := res.GetPayload()

	return &types.Policy{
		ID:        payload.ID,
		Subjects:  payload.Subjects,
		Resources: payload.Resources,
		Actions:   payload.Actions,
	}, nil
}

func (e *enforcer) isAllowed(user string, resource string, action string) (*bool, error) {
	input := &models.OryAccessControlPolicyAllowedInput{
		Action:   action,
		Subject:  user,
		Resource: resource,
	}
	params := &engines.DoOryAccessControlPoliciesAllowParams{
		Body:   input,
		Flavor: string(e.flavor),
	}

	res, err := e.ketoClient.DoOryAccessControlPoliciesAllow(params.WithTimeout(e.timeout))
	if err != nil {
		switch d := err.(type) {
		case *engines.DoOryAccessControlPoliciesAllowForbidden:
			return d.GetPayload().Allowed, nil
		default:
			return nil, err
		}
	}

	return res.GetPayload().Allowed, nil
}

func (e *enforcer) formatUser(user string) string {
	match, _ := regexp.MatchString("users:.*", user)
	if match {
		return user
	}
	return fmt.Sprintf("users:%s", user)
}

func (e *enforcer) formatResource(resource string) string {
	match, _ := regexp.MatchString(fmt.Sprintf("resources:%s:.*", e.product), resource)
	if match {
		return resource
	}
	return fmt.Sprintf("resources:%s:%s", e.product, resource)
}

func (e *enforcer) formatRole(role string) string {
	match, _ := regexp.MatchString(fmt.Sprintf("roles:%s:.*", e.product), role)
	if match {
		return role
	}
	return fmt.Sprintf("roles:%s:%s", e.product, role)
}

func (e *enforcer) formatPolicy(policy string) string {
	match, _ := regexp.MatchString(fmt.Sprintf("policies:%s:.*", e.product), policy)
	if match {
		return policy
	}
	return fmt.Sprintf("policies:%s:%s", e.product, policy)
}

func (e *enforcer) stripResourcePrefix(resource string) string {
	return strings.Replace(resource, fmt.Sprintf("resources:%s:", e.product), "", 1)
}
