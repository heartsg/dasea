package types

//
// Credential is the plugin other than password & token?
// Currently I do not know how to authenticate a user via
// credential, because /auth/tokens seems only support
// password or token authenticate.
//
//
// So we currently ignore all credential operations in client
// implementation until in future we may need it.
//
// No test is provided currently for this type
//


type Credential struct {
	Blob string `json:"blob,omitempty"`
	UserId string `json:"user_id,omitempty"`
	ProjectId string `json:"project_id,omitempty"`
	Type string `json:"type,omitempty"`
	Links Link `json:"links,omitempty"`
	Id string `json:"id,omitempty"`
}
type CredentialRequest struct {
	Credential *Credential `json:"credential"`
}
type CredentialResponse struct {
	Credential *Credential `json:"credential"`
}
type CredentialsResponse struct {
	Credentials []*Credential `json:"credentials"`
	Links Link `json:"links"`
}

func NewCredentialRequest(c *Credential) *CredentialRequest {
	return &CredentialRequest {
		Credential: c,
	}
}


// Post /v3/credentials
// Function: create a new credential
// Request: CredentialRequest (without id)
// Response: CredentialResponse (with id)

// Get /v3/credentials
// Function: list all credentials
// Request: user_id in query string ?user_id=xxx
// Response: CredentialsResponse

// Get /v3/credentials/{credential_id}
// Response: CredentialResponse

// Patch /v3/credentials/{credential_id}
// Request: credentialRequest
// Response: credentialResponse

// Delete /v3/credentials/{credential_id}