package types


import (
	"errors"
	"github.com/heartsg/dasea/keystone/keystoneclient/util"
)


// Defines all requests & response structures for keystone Identity API V3.
// Follows the docs in http://developer.openstack.org/api-ref-identity-v3.html.
// Also to facilitate Json marshalling and unmarshalling.


// For auth and token related request/response.

//
// Post /v3/auth/tokens
// ?nocatelog might be appended, such as Post /v3/auth/tokens?nocatelog
//
// Functions: Password or token authentication with unscoped/scoped authroization
//
// Request: struct Auth
// Request rules:
//  1. If methods = ["password"], the AuthRequest.Auth.Identity.Password must be present
//     Password.User.Passowrd must also be present, it is password authentication.
//
//     If Auth.Scope is a string "unscoped" or omitted (no default project for the user), 
//     it is unscoped authentication. The difference is that, if it is omitted, 
//     and the user has a default project, the returned token might be project-scoped.
//
//     If Auth.Scope is a map[string]map[string]string, with conent equals to 
//			{
//				"project": {
//					"id":"xxxxxx"  //or "name":"yyyyyy"
//				}
//			}
//			or {
//				"domain": {
//					"id":"xxxxxx" //or "name":"yyyyyy"
//				}
//			}
//     It will be domain scoped or project scoped. Note that, we can generate this json
//     either from map[string]map[string]string, or from Project or Domain struct with
//     only Id or name field non-empty.
//
// 2. If methods = ["token"], the AuthRequest.Auth.Identity.Token must be present.
//
//    It can be scoped or unscoped the same as password authentication.
//
//

type AuthRequestPassword struct {
	User *User `json:"user"`
}
type AuthRequestToken struct {
	Id string `json:"id"`
}
type AuthRequestIdentity struct {
	Methods []string `json:"methods"`
	Password *AuthRequestPassword `json:"password,omitempty"`
	Token *AuthRequestToken `json:"token,omitempty"`
}
type AuthRequestAuth struct {
	Identity *AuthRequestIdentity `json:"identity"`
	// Use interface{} because scope can be either Scope or
	// a string "unscoped"
	Scope interface{} `json:"scope,omitempty"`
}
type AuthRequest struct {
	Auth *AuthRequestAuth `json:"auth"`
}

func NewAuthRequest() *AuthRequest {
	// Other fields can be empty (so omitted from request)
	// We only create those fields that are not empty during request
	return &AuthRequest{
		Auth: &AuthRequestAuth{
			Identity: &AuthRequestIdentity {
				Methods: make([]string, 1),
			},
		},
	}
}

func NewAuthRequestPassword() *AuthRequest {
	request := NewAuthRequest()
	request.Auth.Identity.Methods[0] = "password"
	request.Auth.Identity.Password = &AuthRequestPassword{
		User: &User {},
	}
	return request
}

func NewAuthRequestToken() *AuthRequest {
	request := NewAuthRequest()
	request.Auth.Identity.Methods[0] = "token"
	request.Auth.Identity.Token = &AuthRequestToken{}
	return request
}
type AuthRequestParams struct {
	Username string
	UserId string
	Password string
	Token string
	DomainName string
	DomainId string
	ProjectName string
	ProjectId string
	Scope bool
	ExplicitUnscope bool
}

func NewAuthRequestFromParams(p *AuthRequestParams) (*AuthRequest, error) {
	// Check input validity
	if (p.Password != "" && p.UserId == "" && p.Username == "") {
		return nil, errors.New("For password authentication, user Id and username cannot be both empty")
	}
	if (p.Password == "" && p.Token == "") {
		return nil, errors.New("Password and token cannot be both empty")
	}
	
	var auth *AuthRequest
	if (p.Password != "") {
		auth = NewAuthRequestPassword()
		auth.Auth.Identity.Password.User.Id = p.UserId
		auth.Auth.Identity.Password.User.Name = p.Username
		auth.Auth.Identity.Password.User.Password = p.Password
		if p.DomainId != "" || p.DomainName != "" {
			auth.Auth.Identity.Password.User.Domain = &Domain {
				Id: p.DomainId,
				Name: p.DomainName,
			}
		}
	} else {
		auth = NewAuthRequestToken()
		auth.Auth.Identity.Token.Id = p.Token
	}
	
	if p.Scope {
		if p.ProjectId != "" || p.ProjectName != "" {
			scope := &Scope {
				Project: &Project {
					Id: p.ProjectId,
					Name: p.ProjectName,
				},
			}
			
			if p.DomainId != "" || p.DomainName != "" {
				scope.Project.Domain = &Domain {
					Id: p.DomainId,
					Name: p.DomainName,
				}
			}
			
			auth.Auth.Scope = scope
		} else {
			auth.Auth.Scope = &Scope {
				Domain: &Domain {
					Id: p.DomainId,
					Name: p.DomainName,
				},
			}
		}
	} else if p.ExplicitUnscope {
		auth.Auth.Scope = "unscoped"
	} 
	
	return auth, nil
}

// Response: struct Token
//
// Note that, the token Id is returned in header rather than body (Json data). The token
// is returned as "X-Subject-Token":"id"
//
// For token authentication, an extral "X-Auth-Token" is also provided in response
// header.
//

// We don't really use catalog, it can be safely ignored so far
type AuthResponseCatalog struct {
	Endpoints []*Endpoint `json:"endpoints"`
	Type string `json:"type"`
	Id string `json:"id"`
	Name string `json:"name"`
}

type AuthResponseToken struct {
	Methods []string `json:"methods"`
	ExpiresAt *util.Iso8601DateTime `json:"expires_at"`
	Extras map[string]interface{} `json:"extras"`
	Roles []*Role `json:"roles,omitempty"`
	Project *Project `json:"project,omitempty"`
	User *User `json:"user,omitempty"`
	Catalog []*AuthResponseCatalog `json:"catalog,omitempty"`
	AuditIds []string `json:"audit_ids"`
	IssuedAt *util.Iso8601DateTime `json:"issued_at"`
}
type AuthResponse struct {
	Token *AuthResponseToken `json:"token"`
}



//
// Get /v3/auth/tokens
//
// Functions: Validate and returns information for this token
//
// Request: Pass own token (service token?) in X-Auth-Token
//          Pass token to be validated in X-Subject-Token.
//          Nothing in body.
//
// Response: X-Auth-Token & X-Subject-Token
//           body is same as Post /v3/auth/tokens which is AuthResponse
//           



//
// Head /v3/auth/tokens
//
// Functions: Validate and returns information for this token, but no body
//            returns. All others similar to Get
//


//
// Delete /v3/auth/tokens
//
// Request: X-Auth-Token & X-Subject-Token
// No response.
//        
