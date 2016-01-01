package keystonemiddleware

import (
	"net/http"
	"errors"
	"strings"
	
	"github.com/heartsg/dasea/keystone/keystoneclient/types"
	"github.com/heartsg/dasea/keystone/keystoneclient/client"
	"github.com/heartsg/dasea/keystone/keystoneclient/resource"
	"github.com/heartsg/dasea/keystone/keystonemiddleware/cache"
	"github.com/heartsg/dasea/router"
	"golang.org/x/net/context"
)

const (
	UserAccessInfoKey = "UserAccessInfo"
)

//
// AuthToken middleware
//
type AuthToken struct {
	session *client.Session
	cache *cache.TokenCache
	revocations []string
	delayAuthDecision bool
}

func NewAuthToken(opts *Opts) *AuthToken {
	// Init session, session's auth shall be a auto-fetch password
	// or token authentication for a service user (have super previlege)
	// We will use this service user for:
	// 1. validate a normal user token (sent in via request's "X-Auth-Token")
	// 2. validate another service token (sent in via request's "X-Service-Token")
	// 3. Retrieve revocation list from identity server (not implemented)
	params := &types.AuthRequestParams {
		Username : opts.Username,
		UserId : opts.UserId,
		Password : opts.Password,
		Token : opts.Token,
		DomainName: opts.UserDomainName,
		DomainId: opts.UserDomainId,
		ProjectName: opts.ProjectName,
		ProjectId: opts.ProjectId,
		Scope: true,
	}
	var auth client.Auth
	if opts.AuthMethod == "password" {
		auth = &client.AuthPassword { Params: params }
	} else {
		auth = &client.AuthToken { Params: params }
	}
	s := client.NewSession(&opts.Client, auth)
	
	t := cache.NewTokenCache(opts.MemcacheServers, opts.TokenCacheTime)
	
	return &AuthToken { session: s, cache: t, delayAuthDecision: opts.DelayAuthDecision }
}

// The sequence :
// 1. Get X_AUTH_TOKEN and X_SERVICE_TOKEN from request
// 2. Get token info (access) from cache server
// 3. If token valid, put related data into context
// 4. If not, fetch token from identity server
// 5. If token valid, put related data into context
// 6. Or put invalidated token data into context
// 7. If auth delay enabled, return context
// 8. Or return nil
func (a *AuthToken) ServeHTTPContext(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	userToken := r.Header.Get("X-Auth-Token")
	serviceToken := r.Header.Get("X-Service-Token")
	
	var userValid bool
	var serviceValid bool
	
	if userToken != "" {
		access, err := a.fetchTokenInfo(userToken)
		if err == nil && !access.Expired() {
			ctx = router.SetMiddlewareParam(ctx, UserAccessInfoKey, access)
			r.Header.Set("X-Identity-Status", "Confirmed")
			if access.TokenInfo != nil {
				if access.TokenInfo.User != nil && access.TokenInfo.User.Domain != nil {
					r.Header.Set("X-Domain-Id", access.TokenInfo.User.Domain.Id)
					r.Header.Set("X-Domain-Name", access.TokenInfo.User.Domain.Name)
					r.Header.Set("X-User-Domain-Id", access.TokenInfo.User.Domain.Id)
					r.Header.Set("X-User-Domain-Name", access.TokenInfo.User.Domain.Name)
				}
				if access.TokenInfo.Project != nil {
					r.Header.Set("X-Project-Id", access.TokenInfo.Project.Id)
					r.Header.Set("X-Project-Name", access.TokenInfo.Project.Name)
				}
				if access.TokenInfo.Project != nil && access.TokenInfo.Project.Domain != nil {
					r.Header.Set("X-Project-Domain-Id", access.TokenInfo.Project.Domain.Id)
					r.Header.Set("X-Project-Domain-Name", access.TokenInfo.Project.Domain.Name)
				}
				if access.TokenInfo.User != nil {
					r.Header.Set("X-User-Id", access.TokenInfo.User.Id)
					r.Header.Set("X-User-Name", access.TokenInfo.User.Name)
				}
				if access.TokenInfo.Roles != nil {
					roleNames := make([]string, len(access.TokenInfo.Roles))
					for i, r := range access.TokenInfo.Roles {
						roleNames[i] = r.Name
					}
					r.Header.Set("X-Roles", strings.Join(roleNames, ",") )
				}
			}
			userValid = true
		}
	}
	
	if !userValid {
		r.Header.Set("X-Identity-Status", "Invalid")
	}
	
	if serviceToken != "" {
		access, err := a.fetchTokenInfo(serviceToken)
		if err == nil && !access.Expired() {
			r.Header.Set("X-Service-Identity-Status", "Confirmed")
			if access.TokenInfo != nil {
				if access.TokenInfo.User != nil && access.TokenInfo.User.Domain != nil {
					r.Header.Set("X-Service-Domain-Id", access.TokenInfo.User.Domain.Id)
					r.Header.Set("X-Service-Domain-Name", access.TokenInfo.User.Domain.Name)
					r.Header.Set("X-Service-User-Domain-Id", access.TokenInfo.User.Domain.Id)
					r.Header.Set("X-Service-User-Domain-Name", access.TokenInfo.User.Domain.Name)
				}
				if access.TokenInfo.Project != nil {
					r.Header.Set("X-Service-Project-Id", access.TokenInfo.Project.Id)
					r.Header.Set("X-Service-Project-Name", access.TokenInfo.Project.Name)
				}
				if access.TokenInfo.Project != nil && access.TokenInfo.Project.Domain != nil {
					r.Header.Set("X-Service-Project-Domain-Id", access.TokenInfo.Project.Domain.Id)
					r.Header.Set("X-Service-Project-Domain-Name", access.TokenInfo.Project.Domain.Name)
				}
				if access.TokenInfo.User != nil {
					r.Header.Set("X-Service-User-Id", access.TokenInfo.User.Id)
					r.Header.Set("X-Service-User-Name", access.TokenInfo.User.Name)
				}
				if access.TokenInfo.Roles != nil {
					roleNames := make([]string, len(access.TokenInfo.Roles))
					for i, r := range access.TokenInfo.Roles {
						roleNames[i] = r.Name
					}
					r.Header.Set("X-Service-Roles", strings.Join(roleNames, ",") )
				}
			}
			serviceValid = true
		}
	} else {
		// if service token not present, make service true here
		serviceValid = true
	}
	if !(userValid && serviceValid) {
		if a.delayAuthDecision {
			return ctx
		} else {
			// reject due to invalid tokens (either user or service or both)
			w.Header().Set("WWW-Authenticate", a.session.BaseUrl)
			if userValid {
				w.WriteHeader(http.StatusForbidden)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
			}
			//return nil context will stop the middleware chain
			return nil
		}
	}
	return ctx
}

func (a *AuthToken) fetchTokenInfo(token string) (*client.AccessInfo, error) {
	info, err := a.cache.Get(token)
	if err == nil {
		return client.NewAccess(token, info), nil
	}
	//if got err, try to retrive from itentity server
	auth := &resource.Auth { Session: a.session }
	accessInfo, err := auth.GetAccessFromToken(token)
	
	if err == nil {
		a.cache.Set(token, accessInfo.TokenInfo)
		return accessInfo, nil
	}
	
	return nil, errors.New("Token info cannot be fetched")
}