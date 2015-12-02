package types

type Role struct {
	Id string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`	
	Links Link `json:"links,omitempty"`
}


type RoleRequest struct {
	Role *Role `json:"role"`
}

type RoleResponse struct {
	Role *Role `json:"role"`
}

type RolesResponse struct {
	Roles []*Role `json:"roles"`
	Links Link `json:"links"`
}

func NewRoleRequest(r *Role) *RoleRequest {
	return &RoleRequest {
		Role: r,
	}
}

// List roles for user on domain
// Get /v3/domains/​{domain_id}​/users/​{user_id}​/roles
// Response: RolesResponse

// Grant role to user on domain
// Put /v3/domains/​{domain_id}​/users/​{user_id}​/roles/​{role_id}​

// Check whether user has a role on domain
// Head /v3/domains/​{domain_id}​/users/​{user_id}​/roles/​{role_id}​

// Revoke
// Delete /v3/domains/​{domain_id}​/users/​{user_id}​/roles/​{role_id}​

// For domain & groups
// /v3/domains/​{domain_id}​/groups/​{group_id}​/roles

// For users & projects
// /v3/projects/​{project_id}​/users/​{user_id}​/roles

// For group & project
// /v3/projects/​{project_id}​/groups/​{group_id}​/roles/​{role_id}​



// /v3/role_assignments
type RoleAssignments struct {
	Links Link `json:"links,omitempty"`
	Role *Role `json:"role,omitempty"`
	User *User `json:"user,omitempty"`
	Group *Group `json:"group,omitempty"`
	Scope *Scope `json:"scope,omitempty"`
}
type RoleAssignmentsResponse struct {
	RoleAssignments []*RoleAssignments `json:"role_assignments"`
	Links Link `json:"links"`
}