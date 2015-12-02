package resource

import (
	"errors"
	
	"github.com/heartsg/dasea/keystone/keystoneclient/client"
	"github.com/heartsg/dasea/requests"
)


type Auth struct {
	// Every resource object should have a client.Session member
	// client.Session will automatically insert token into "X-Auth-Token" header
	Session *client.Session
}

// Get AccessInfo from a token
func (a *Auth) GetAccessFromToken(token string) (*client.AccessInfo, error) {	
	if a.Session == nil  {
		return nil, errors.New("No client.Session to send the request.")
	}
	
	resp, body, err := a.Session.Request("/auth/tokens", 
		requests.GET, 
		map[string]string{"X-Subject-Token": token}, //header
		nil,  // query params
		nil,  // body data
		true)
	if err != nil {
		return nil, err
	}
	access, err := client.NewAccessFromResponseBody(resp, body)
	if err != nil {
		return nil, err
	}
	
	return access, nil
}

// Validate a token
func (a *Auth) ValidateToken(token string) (bool, error) {
	if a.Session == nil  {
		return false, errors.New("No client.Session to send the request.")
	}
	
	_, _, err := a.Session.Request("/auth/tokens", 
		requests.HEAD, 
		map[string]string{"X-Subject-Token": token}, //header
		nil,  // query params
		nil,  // body data
		true)
	if err != nil {
		return false, err
	}
	
	return true, nil
}


// Revoke a token
func (a *Auth)DeleteToken(token string) error {
	if a.Session == nil  {
		return errors.New("No client.Session to send the request.")
	}
	
	_, _, err := a.Session.Request("/auth/tokens", 
		requests.DELETE, 
		map[string]string{"X-Subject-Token": token}, //header
		nil,  // query params
		nil,  // body data
		true)
	if err != nil {
		return err
	}

	return nil
}