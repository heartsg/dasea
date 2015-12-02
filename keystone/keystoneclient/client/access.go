package client

import (
	"net/http"
	"encoding/json"
	"fmt"
	"time"
	
	"github.com/heartsg/dasea/keystone/keystoneclient/types"
	"github.com/heartsg/dasea/keystone/keystoneclient/util"
)

const (
	//30 seconds 
	StaleDuration int = 30
)

// If more keystone versions are to be supported, we could use
// AccessInfo as interface, and AccessInfoV3, AccessInfoV4 etc.
// as implementation.
//
// However, we currently decide to only support keystone identity
// API v3. We leave the generalization to future updates.

type AccessInfo struct {
	// New token is returned from identity service via
	// response.Header.Get("X-Subject-Token"). One can also
	// manually input token values if it is known in advance.
	Token string
	// Either user provided, or returned from identity service
	// response. The format of this stirng is as in 
	// http://developer.openstack.org/api-ref-identity-v3.html.
	// An example would be (JSON format):
	//   {
	//		"token": {
	//			"methods": [
	//				"password"
	//			],
	//			"expires_at": "2015-11-06T15:32:17.893769Z",
	//			"extras": {},
	//			"user": {
	//				"domain": {
	//					"id": "default",
	//					"name": "Default"
	//				},
	//				"id": "423f19a4ac1e4f48bbb4180756e6eb6c",
	//				"name": "admin"
	//			},
	//			"audit_ids": [
	//				"ZzZwkUflQfygX7pdYDBCQQ"
	//			],
	//			"issued_at": "2015-11-06T14:32:17.893797Z"
	//		}
	//	}
	TokenInfo *types.Token
}

// Manually create an AccessInfo
func NewAccess(token string, tokenInfo *types.Token) *AccessInfo {
	return &AccessInfo{
		Token: token,
		TokenInfo: tokenInfo,
	}
}


// The interface is different, an error has to be returned because
// it is unknown whether the response contains the required infomation
func NewAccessFromResponse(resp *http.Response) (*AccessInfo, error) {
	token := resp.Header.Get("X-Subject-Token")
	
	if token == "" {
		return nil, fmt.Errorf("Unrecognized auth response: X-Subject-Token header does not exist.")
	}
	

	f := &types.AuthResponse{}
	err := json.NewDecoder(resp.Body).Decode(f)
	if (err != nil) {
		return nil, fmt.Errorf("Unrecognized auth response: %s", err)
	}
	
	return NewAccess(token, (*types.Token)(f.Token)), nil
}

func NewAccessFromResponseBody(resp *http.Response, body []byte) (*AccessInfo, error) {
	token := resp.Header.Get("X-Subject-Token")
	
	if token == "" {
		return nil, fmt.Errorf("Unrecognized auth response: X-Subject-Token header does not exist.")
	}
	

	f := &types.AuthResponse{}
	err := json.Unmarshal(body, f)
	if (err != nil) {
		return nil, fmt.Errorf("Unrecognized auth response: %s", err)
	}
	
	return NewAccess(token, (*types.Token)(f.Token)), nil
}

func (a *AccessInfo) Expires() *util.Iso8601DateTime {
	if a.TokenInfo == nil {
		return nil
	}
	return a.TokenInfo.ExpiresAt
}

func (a *AccessInfo) WillExpireSoon() bool {
	expires := a.Expires()
	if expires == nil {
		return false
	}
	t := expires.Add(time.Duration(-StaleDuration) * time.Second)
	return t.Before(time.Now())
}



