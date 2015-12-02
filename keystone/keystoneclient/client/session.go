package client

import (
	"time"
	"fmt"
	"crypto/tls"
	"net/http"
	"errors"
	
	"github.com/heartsg/dasea/requests"
)
//
// Session will automatically deal with a lot of repeated stuff
// for a http.Client request.
//
// Includes:
//  - auto `X-Auth-Token` insertion for authentication purpose
//  - tls and certificates
//  - Auto authentication & re-authenticate (auto insert AccessInfo after re-auth)
//  - dealing with response errors
//
// Note, the auto-authentication also relies on session (without the
// `X-Auth-Token` insertion)
//
// Note, for future version. We roughly follow the same idea as python-
// golangclient's session implementation. The python version has
// native implementation of retries/redirects. Currently we let redirects
// handled by requests.Requests, and do not attempt retries if failed.
// We may add in these functions later on.
//


// requests: used to send requests. Requests caters httpClient.Do(request), and
//           also deals with redirect, timeout, certificates/https and etc.
// Auth: used to authenticate and get access information
type Session struct {
	requests *requests.Requests
	Auth Auth
	
	// keystone server url, normally it's the same value
	// as in KeystoneclientOpts.AuthUrl.
	// e.g.: 127.0.0.1:35357/v3
	BaseUrl string
	
	// session related data
	originalIp string
	userAgent string
}

//
// Python-keystoneclient's session has the following attributes during init.
//  - Auth plugin (either password or token)
//  - requests.session (deal with cookies, cert, tls etc.)
//  - original_ip of user
//  - verify (whether client verify server certificates)
//  - client certificates
//  - connection timeout
//  - user_agent
//
// Make these parameters as individually configured.
//
func NewSession(url string, auth Auth) *Session {
	s := &Session {
		requests: requests.New(),
		Auth: auth,
		BaseUrl: url,
		userAgent: "Dasea Keystone client",
	}
	s.requests.SetResponseBodyValid(false)
	return s
}
func (s *Session) Timeout(timeout time.Duration) *Session {
	s.requests.Timeout(timeout)
	return s
}
func (s *Session) TLSClientConfig(config *tls.Config) *Session {
	s.requests.TLSClientConfig(config)
	return s
}
func (s *Session) MaxRedirect(maxRedirect int) *Session {
	s.requests.Redirect(maxRedirect, true)
	return s
}
//Headers will be cleared for each request, so we will save them
func (s *Session) OriginalIp(originalIp string) *Session {
	s.originalIp = originalIp
	return s
}
func (s *Session) UserAgent(userAgent string) *Session {
	s.userAgent = userAgent
	return s
}

//
// Python-keystoneclient's session request allows overwriting session's own
// parameters such as timeout, originalIp in request.
//
// However, since we already allowed these settings to be set individually,
// there's no need to provide them here too.
//
// Python's implementation also provides endpoint-filters etc to automatically
// check which endpoint the request uses. We don't implement this.
//
// url : the partial url. Examples: /auth/tokens, /users/{user_id} etc.
// method: "GET", "PUT", "POST" and etc.
// queryParams: will be put as ?key1=value1&key2=value2
// data: will be converted to json and put in body
// requireAuthentication: if true, will insert X-Auth-Token (by session.Auth)
func (s *Session) Request(url string, 
	method string, 
	header map[string]string, 
	queryParams map[string]string, 
	data interface{}, 
	requireAuthentication bool) (*http.Response, []byte, error) {
	
	//check authentication, we use same session for auth
	var token string
	
	if requireAuthentication {
		if s.Auth != nil {
			access, err := s.Auth.GetAccess(s)
			if err != nil {
				return nil, nil, err
			}
			token = access.Token
		} else {
			return nil, nil, errors.New("No auth plugin to retrieve auth-token")
		}
	}
	
	s.requests.Request(s.BaseUrl + url, method)
	if token != "" {
		s.requests.Set("X-Auth-Token", token)
	}
	if s.userAgent != "" {
		s.requests.Set("User-Agent", s.userAgent)
	}
	if s.originalIp != "" {
		s.requests.Set("Forwarded", fmt.Sprintf("for=%s;by=%s", s.originalIp, s.userAgent))
	}
	
	if header != nil {
		for k, v := range header {
			s.requests.Set(k, v)
		}
	}
	if queryParams != nil {
		for k, v := range queryParams {
			s.requests.Param(k, v)
		}
	}
	if data != nil {
		//Try to optimize by not calling Send(data)
		s.requests.SendRawStruct(data)
	}
	
	resp, body, errs := s.requests.EndBytes()
	
	if errs != nil {
		return nil, nil, fmt.Errorf("%s", errs)
	}
	
	//check resp errors
	//Bad Request (400)
	//Unauthorized (401)
	//Forbidden (403)
	//Not Found (404)
	//Method Not Allowed (405)
	//conflict (409)
	//Request Entity Too Large (413)
	//Unsupported Media Type (415)
	//Service Unavailable (503)
	
	if resp.StatusCode >= 400 {
		return resp, nil, fmt.Errorf("Response error: (%d) %s", resp.StatusCode, resp.Status)
	}
	
	return resp, body, nil
}