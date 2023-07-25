package enforcer

import (
	"context"
	"fmt"
	"sync"

	ory "github.com/ory/keto-client-go"
	"golang.org/x/exp/slices"
	"golang.org/x/sync/errgroup"
)

// Enforcer interface to enforce authorization
type Enforcer interface {
	// IsUserGrantedPermission check whether user has the required permission, both directly and indirectly
	IsUserGrantedPermission(ctx context.Context, user string, permission string) (bool, error)
	// GetUserRoles get all roles directly associated with a user
	GetUserRoles(ctx context.Context, user string) ([]string, error)
	// GetRolePermissions get all permissions directly associated with a role
	GetRolePermissions(ctx context.Context, role string) ([]string, error)
	// GetRoleMembers get all members for a role
	GetRoleMembers(ctx context.Context, role string) ([]string, error)
	// UpdateAuthorization update authorization rules in batches
	UpdateAuthorization(ctx context.Context, updateRequest AuthorizationUpdateRequest) error
}

// CacheConfig holds the configuration for the in-memory cache, if enabled
type CacheConfig struct {
	KeyExpirySeconds            int
	CacheCleanUpIntervalSeconds int
}

// MaxKeyExpirySeconds is the max allowed value for the KeyExpirySeconds.
const MaxKeyExpirySeconds = 600

type enforcer struct {
	cache           *InMemoryCache
	ketoReadClient  *ory.APIClient
	ketoWriteClient *ory.APIClient
}

func newEnforcer(
	ketoRemoteRead string,
	ketoRemoteWrite string,
	cacheConfig *CacheConfig,
) (*enforcer, error) {
	readConfiguration := ory.NewConfiguration()
	readConfiguration.Servers = []ory.ServerConfiguration{
		{
			URL: ketoRemoteRead,
		},
	}
	writeConfiguration := ory.NewConfiguration()
	writeConfiguration.Servers = []ory.ServerConfiguration{
		{
			URL: ketoRemoteWrite,
		},
	}
	enforcer := &enforcer{
		ketoReadClient:  ory.NewAPIClient(readConfiguration),
		ketoWriteClient: ory.NewAPIClient(writeConfiguration),
	}

	if cacheConfig != nil {
		if cacheConfig.KeyExpirySeconds > MaxKeyExpirySeconds {
			return nil, fmt.Errorf("Configured KeyExpirySeconds is larger than the max permitted value of %d",
				MaxKeyExpirySeconds)
		}
		enforcer.cache = newInMemoryCache(cacheConfig.KeyExpirySeconds, cacheConfig.CacheCleanUpIntervalSeconds)
	}
	return enforcer, nil
}

func (e *enforcer) IsUserGrantedPermission(ctx context.Context, user string, permission string) (bool, error) {
	if e.isCacheEnabled() {
		if isAllowed, found := e.cache.LookUpUserPermissions(user, permission); found {
			return *isAllowed, nil
		}
	}
	checkPermissionResult, _, err := e.ketoReadClient.PermissionApi.CheckPermission(ctx).
		Namespace("Permission").
		Object(permission).
		Relation("granted").
		SubjectSetNamespace("Subject").
		SubjectSetObject(user).
		SubjectSetRelation("").
		Execute()
	if err != nil {
		return false, err
	}
	userHasPermission := checkPermissionResult.Allowed
	if e.isCacheEnabled() {
		e.cache.StoreUserPermissions(user, permission, userHasPermission)
	}
	return userHasPermission, nil
}

func (e *enforcer) GetUserRoles(ctx context.Context, user string) ([]string, error) {
	roleRelationships, _, err := e.ketoReadClient.RelationshipApi.GetRelationships(ctx).
		Namespace("Role").
		Relation("member").
		SubjectSetNamespace("Subject").
		SubjectSetRelation("").
		SubjectSetObject(user).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}
	roles := make([]string, 0)
	for _, tuple := range roleRelationships.RelationTuples {
		roles = append(roles, tuple.Object)
	}

	return roles, nil
}

func (e *enforcer) GetRolePermissions(ctx context.Context, role string) ([]string, error) {
	permissionRelationships, _, err := e.ketoReadClient.RelationshipApi.GetRelationships(ctx).
		Namespace("Permission").
		Relation("granted").
		SubjectSetNamespace("Role").
		SubjectSetRelation("member").
		SubjectSetObject(role).Execute()
	if err != nil {
		return nil, err
	}
	permissions := make([]string, 0)
	for _, tuple := range permissionRelationships.RelationTuples {
		permissions = append(permissions, tuple.Object)
	}

	return permissions, nil
}

func (e *enforcer) GetRoleMembers(ctx context.Context, role string) ([]string, error) {
	expandedRole, _, err := e.ketoReadClient.PermissionApi.ExpandPermissions(ctx).
		Namespace("Role").
		Relation("member").
		Object(role).
		MaxDepth(2).Execute()
	if err != nil {
		return nil, err
	}
	members := make([]string, 0)
	for _, child := range expandedRole.GetChildren() {
		members = append(members, child.Tuple.SubjectSet.Object)
	}

	return members, nil
}

func newRolePermissionPatch(action string, permission string, role string) ory.RelationshipPatch {
	return ory.RelationshipPatch{
		Action: &action,
		RelationTuple: &ory.Relationship{
			Namespace:  "Permission",
			Object:     permission,
			Relation:   "granted",
			SubjectSet: ory.NewSubjectSet("Role", role, "member"),
		},
	}
}

func newRoleMemberPatch(action string, role string, member string) ory.RelationshipPatch {
	return ory.RelationshipPatch{
		Action: &action,
		RelationTuple: &ory.Relationship{
			Namespace:  "Role",
			Object:     role,
			Relation:   "member",
			SubjectSet: ory.NewSubjectSet("Subject", member, ""),
		},
	}
}

func (e *enforcer) UpdateAuthorization(ctx context.Context, updateRequest AuthorizationUpdateRequest) error {
	var existingRolePermissions sync.Map
	var existingRoleMembers sync.Map
	getRelationsWorkersGroup := new(errgroup.Group)
	for role := range updateRequest.RolePermissions {
		updatedRole := role
		getRelationsWorkersGroup.Go(func() error {
			permissions, err := e.GetRolePermissions(ctx, updatedRole)
			if err != nil {
				return err
			}
			existingRolePermissions.Store(updatedRole, permissions)
			return nil
		})
	}
	for role := range updateRequest.RoleMembers {
		updatedRole := role
		getRelationsWorkersGroup.Go(func() error {
			members, err := e.GetRoleMembers(ctx, updatedRole)
			if err != nil {
				return err
			}
			existingRoleMembers.Store(updatedRole, members)
			return nil
		})
	}
	err := getRelationsWorkersGroup.Wait()
	if err != nil {
		return err
	}
	patches := make([]ory.RelationshipPatch, 0)
	existingRolePermissions.Range(func(key, value interface{}) bool {
		role := key.(string)
		permissions := value.([]string)
		for _, permission := range permissions {
			if !slices.Contains(updateRequest.RolePermissions[role], permission) {
				patches = append(patches, newRolePermissionPatch("delete", permission, role))
			}
		}
		return true
	})

	for role, permissions := range updateRequest.RolePermissions {
		for _, permission := range permissions {
			result, found := existingRolePermissions.Load(role)
			existingPermissions := result.([]string)
			if found && !slices.Contains(existingPermissions, permission) {
				patches = append(patches, newRolePermissionPatch("insert", permission, role))
			}
		}
	}

	existingRoleMembers.Range(func(key, value interface{}) bool {
		role := key.(string)
		members := value.([]string)
		for _, member := range members {
			if !slices.Contains(updateRequest.RoleMembers[role], member) {
				patches = append(patches, newRoleMemberPatch("delete", role, member))
			}
		}
		return true
	})

	for role, members := range updateRequest.RoleMembers {
		for _, member := range members {
			result, found := existingRoleMembers.Load(role)
			existingMembers := result.([]string)
			if found && !slices.Contains(existingMembers, member) {
				patches = append(patches, newRoleMemberPatch("insert", role, member))
			}
		}
	}

	_, err = e.ketoWriteClient.RelationshipApi.PatchRelationships(ctx).RelationshipPatch(patches).Execute()
	return err
}

func (e *enforcer) isCacheEnabled() bool {
	return e.cache != nil
}

// NewAuthorizationUpdateRequest create a new AuthorizationUpdateRequest. Multiple operations can be chained together
// using the SetRolePermissions and SetRoleMembers methods. No changes will be made until the AuthorizationUpdateRequest
// object is passed to the Enforcer, in which all the previously chained operations will be executed in batch.
func NewAuthorizationUpdateRequest() AuthorizationUpdateRequest {
	return AuthorizationUpdateRequest{
		RolePermissions: make(map[string][]string),
		RoleMembers:     make(map[string][]string),
	}
}

type AuthorizationUpdateRequest struct {
	RolePermissions map[string][]string
	RoleMembers     map[string][]string
}

// SetRolePermissions set the permissions for a role. If the role already has permissions, they will be replaced.
func (a AuthorizationUpdateRequest) SetRolePermissions(role string,
	permissions []string) AuthorizationUpdateRequest {
	a.RolePermissions[role] = permissions
	return a
}

// SetRoleMembers set the members for a role. If the role already has members, they will be replaced.
func (a AuthorizationUpdateRequest) SetRoleMembers(role string, members []string) AuthorizationUpdateRequest {
	a.RoleMembers[role] = members
	return a
}
