package resource


import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/heartsg/dasea/keystone/keystoneclient/types"
	"github.com/heartsg/dasea/keystone/keystoneclient/client"
	"github.com/heartsg/dasea/requests"
)
type Credential struct {
	Session *client.Session
}

func (r *Credential) Create(c *types.Credential) (*types.Credential, error) {
	if r.Session == nil  {
		return nil, errors.New("No client.Session to send the request.")
	}
	
	credentialRequest := types.NewCredentialRequest(c)
	
	_, body, err := r.Session.Request("/credentials", 
		requests.POST, 
		nil,  //header
		nil,  // query params
		credentialRequest,  // body data
		true)
	if err != nil {
		return nil, err
	}
	
	credentialResponse := &types.CredentialResponse{}
	err = json.Unmarshal(body, credentialResponse)
	if err != nil {
		return nil, err
	}
	
	return credentialResponse.Credential, nil
}

func (r *Credential) List(userId string) (*types.CredentialsResponse, error) {
	if r.Session == nil  {
		return nil, errors.New("No client.Session to send the request.")
	}
	
	queryParams := make(map[string]string)
	if userId != "" {
		queryParams["user_id"] = userId
	}
	
	_, body, err := r.Session.Request("/credentials", 
		requests.GET, 
		nil,  //header
		queryParams,  // query params
		nil,  // body data
		true)
	if err != nil {
		return nil, err
	}
	
	credentialsResponse := &types.CredentialsResponse{}
	err = json.Unmarshal(body, credentialsResponse)
	if err != nil {
		return nil, err
	}
	
	return credentialsResponse, nil	
}

func (r *Credential) Delete(id string) error {
	if r.Session == nil  {
		return errors.New("No client.Session to send the request.")
	}
	
	_, _, err := r.Session.Request(fmt.Sprintf("/credentials/%s", id), 
		requests.DELETE, 
		nil, //header
		nil,  // query params
		nil,  // body data
		true)
	if err != nil {
		return err
	}

	return nil
}

func (r *Credential) Get(id string) (*types.Credential, error) {
	if r.Session == nil  {
		return nil, errors.New("No client.Session to send the request.")
	}
	
	_, body, err := r.Session.Request(fmt.Sprintf("/credentials/%s", id), 
		requests.GET, 
		nil, //header
		nil,  // query params
		nil,  // body data
		true)
	if err != nil {
		return nil, err
	}

	credentialResponse := &types.CredentialResponse{}
	err = json.Unmarshal(body, credentialResponse)
	if err != nil {
		return nil, err
	}
	
	return credentialResponse.Credential, nil
}


func (r *Credential) Update(id string, c *types.Credential) (*types.Credential, error) {
	if r.Session == nil  {
		return nil, errors.New("No client.Session to send the request.")
	}
	
	credentialRequest := types.NewCredentialRequest(c)
	
	_, body, err := r.Session.Request(fmt.Sprintf("/credentials/%s", id), 
		requests.PATCH, 
		nil,  //header
		nil,  // query params
		credentialRequest,  // body data
		true)
	if err != nil {
		return nil, err
	}
	
	credentialResponse := &types.CredentialResponse{}
	err = json.Unmarshal(body, credentialResponse)
	if err != nil {
		return nil, err
	}
	
	return credentialResponse.Credential, nil	
}