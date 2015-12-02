package types

type User struct {
	Id string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	DomainId string `json:"domain_id,omitempty"`
	DefaultProjectId string `json:"default_project_id,omitempty"`
	Email string `json:"email,omitempty"`
	Description string `json:"description,omitempty"`
	Domain *Domain `json:"domain,omitempty"`
	Password string `json:"password,omitempty"`
	OriginalPassword string `json:"original_password,omitempty"`
	Enabled bool `json:"enabled,omitempty"`
	Links Link `json:"links,omitempty"`
}


type UsersResponse struct {
	Users []*User `json:"users"`
	Links Link `json:"links"`	
}
type UserRequest struct {
	User *User `json:"user"`
}
type UserResponse struct {
	User *User `json:"user"`
}

func NewUserRequest(u *User) *UserRequest {
	return &UserRequest {
		User: u,
	}
}

// List users
// Get /v3/users
// Request: domain_id, name, enabled (all in query string)
// Response: UsersResponse


// Create a user
// Post /v3/users
// Request: UserRequest
// Response: UserResponse

// Show user details
// Get /v3/users/{user_id}
// Response: UserResponse

// Update user details
// Patch /v3/users/{user_id}
// Request: UserRequest
// Response: UserResponse

// Delete /v3/users/{user_id}

// Change password
// Post /v3/users/{user_id}/password
// Request: UserRequest (with password & original_password)

// Get user groups
// Get /v3/users/{user_id}/groups
// Response: GroupsResponse

// Get user projects
// Get /v3/user/{user_id}/projects
// Response: ProjectsResponse