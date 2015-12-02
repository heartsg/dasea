package types

type Project struct {
	Id string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Domain *Domain `json:"domain,omitempty"`
	Description string `json:"description,omitempty"`
	DomainId string `json:"domain_id,omitempty"`
	Enabled bool `json:"enabled,omitempty"`
	ParentId string `json:"parent_id,omitempty"`
	IsDomain bool `json:"is_domain,omitempty"`
	Extra map[string]interface{} `json:"extra,omitempty"`
	Links Link `json:"links,omitempty"`
}

type ProjectRequest struct {
	Project *Project `json:"project"`
}

type ProjectResponse struct {
	Project *Project `json:"project"`
}


type ProjectsResponse struct {
	Projects []*Project `json:"projects"`
	Links Link `json:"links"`
}

func NewProjectRequest(p *Project) *ProjectRequest {
	return &ProjectRequest {
		Project: p,
	}
}

// List projects
// Get /v3/projects
// Request: domain_id, parent_id, name, enabled (in query string)
// Response: ProjectsResponse

// Post /v3/projects
// Get /v3/projects/{project_id}
// Patch /v3/projets/{project_id}
// Delete /v3/projects/{project_id}


// Patch /v3/projects/{project_id}/cascade
// Delete /v3/projects/{project_id}/cascade