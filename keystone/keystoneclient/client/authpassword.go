package client

import (
	"errors"
	
	"github.com/heartsg/dasea/keystone/keystoneclient/types"
	"github.com/heartsg/dasea/requests"
)


type AuthPassword struct {
	Params *types.AuthRequestParams
	Access *AccessInfo
}


func (a *AuthPassword) GetAccess(session *Session) (*AccessInfo, error) {
	if a.Access != nil && ! a.Access.WillExpireSoon() {
		return a.Access, nil
	}
	authRequest, err := types.NewAuthRequestFromParams( a.Params )
	if err != nil {
		return nil, err
	}	
	if session == nil  {
		return nil, errors.New("No client.Session to send the authentication request.")
	}
	resp, body, err := session.Request("/auth/tokens", requests.POST, nil, nil, authRequest, false)
	if err != nil {
		return nil, err
	}
	a.Access, err = NewAccessFromResponseBody(resp, body)
	if err != nil {
		return nil, err
	}
	
	return a.Access, nil
}