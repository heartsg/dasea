package types

//
// Common types used by all keystone response/requests
//
type Link map[string]string
type MediaType map[string]string
type Scope struct {
	Project *Project `json:"project,omitempty"`
	Domain *Domain `json:"domain,omitempty"`
}
type Token AuthResponseToken


