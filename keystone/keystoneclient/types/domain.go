package types

// A domain is a collection of users, groups, and projects. 
// Each group and project is owned by exactly one domain.
//
// Each domain defines a namespace where certain API-visible name attributes exist, 
// which affects whether those names must be globally unique or unique within that domain. 
// In the Identity API, the uniqueness of these attributes is as follows:
//
//    Domain name. Globally unique across all domains.
//    Role name. Globally unique across all domains.
//    User name. Unique within the owning domain.
//    Project name. Unique within the owning domain.
//    Group name. Unique within the owning domain.

type Domain struct {
	Id string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Links Link `json:"links,omitempty"`
	Description string `json:"description,omitempty"`
	Enabled bool `json:"enabled,omitempty"`
}

type DomainRequest struct {
	Domain *Domain `json:"domain"`
}
type DomainResponse struct {
	Domain *Domain `json:"domain"`
}
type DomainsResponse struct {
	Domains []*Domain `json:"domains"`
	Links Link `json:"links"`
}

func NewDomainRequest(d *Domain) *DomainRequest {
	return &DomainRequest {
		Domain: d,
	}
}