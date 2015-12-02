package resource


import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/heartsg/dasea/keystone/keystoneclient/types"
	"github.com/heartsg/dasea/keystone/keystoneclient/client"
	"github.com/heartsg/dasea/requests"
)
type Group struct {
	Session *client.Session
}

func (r *Group) Create(d *types.Group) (*types.Group, error) {
	if r.Session == nil  {
		return nil, errors.New("No client.Session to send the request.")
	}
	
	groupRequest := types.NewGroupRequest(d)
	
	_, body, err := r.Session.Request("/groups", 
		requests.POST, 
		nil,  //header
		nil,  // query params
		groupRequest,  // body data
		true)
	if err != nil {
		return nil, err
	}
	
	groupResponse := &types.GroupResponse{}
	err = json.Unmarshal(body, groupResponse)
	if err != nil {
		return nil, err
	}
	
	return groupResponse.Group, nil
}

func (r *Group) List(domainId string, name string) (*types.GroupsResponse, error) {
	if r.Session == nil  {
		return nil, errors.New("No client.Session to send the request.")
	}
	
	queryParams := make(map[string]string)
	if domainId != "" {
		queryParams["domain_id"] = domainId
	}
	if name != "" {
		queryParams["name"] = name
	}
	
	_, body, err := r.Session.Request("/groups", 
		requests.GET, 
		nil,  //header
		queryParams,  // query params
		nil,  // body data
		true)
	if err != nil {
		return nil, err
	}
	
	groupsResponse := &types.GroupsResponse{}
	err = json.Unmarshal(body, groupsResponse)
	if err != nil {
		return nil, err
	}
	
	return groupsResponse, nil	
}

func (r *Group) Delete(id string) error {
	if r.Session == nil  {
		return errors.New("No client.Session to send the request.")
	}
	
	_, _, err := r.Session.Request(fmt.Sprintf("/groups/%s", id), 
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

func (r *Group) Get(id string) (*types.Group, error) {
	if r.Session == nil  {
		return nil, errors.New("No client.Session to send the request.")
	}
	
	_, body, err := r.Session.Request(fmt.Sprintf("/groups/%s", id), 
		requests.GET, 
		nil, //header
		nil,  // query params
		nil,  // body data
		true)
	if err != nil {
		return nil, err
	}

	groupResponse := &types.GroupResponse{}
	err = json.Unmarshal(body, groupResponse)
	if err != nil {
		return nil, err
	}
	
	return groupResponse.Group, nil
}


func (r *Group) Update(id string, d *types.Group) (*types.Group, error) {
	if r.Session == nil  {
		return nil, errors.New("No client.Session to send the request.")
	}
	
	groupRequest := types.NewGroupRequest(d)
	
	_, body, err := r.Session.Request(fmt.Sprintf("/groups/%s", id), 
		requests.PATCH, 
		nil,  //header
		nil,  // query params
		groupRequest,  // body data
		true)
	if err != nil {
		return nil, err
	}
	
	groupResponse := &types.GroupResponse{}
	err = json.Unmarshal(body, groupResponse)
	if err != nil {
		return nil, err
	}
	
	return groupResponse.Group, nil	
}


func (r *Group) ListUsers(id string, domainId string, description string,
	name string, enabled string) (*types.UsersResponse, error) {
	if r.Session == nil  {
		return nil, errors.New("No client.Session to send the request.")
	}
	queryParams := make(map[string]string)
	if domainId != "" {
		queryParams["domain_id"] = domainId
	}
	if name != "" {
		queryParams["name"] = name
	}
	if description != "" {
		queryParams["description"] = description
	}
	if enabled == "true" || enabled == "false" {
		queryParams["enabled"] = enabled
	}
	
	_, body, err := r.Session.Request(fmt.Sprintf("/groups/%s/users", id), 
		requests.GET, 
		nil,  //header
		queryParams,  // query params
		nil,  // body data
		true)
	if err != nil {
		return nil, err
	}
	
	usersResponse := &types.UsersResponse{}
	err = json.Unmarshal(body, usersResponse)
	if err != nil {
		return nil, err
	}
	
	return usersResponse, nil	
}


func (r *Group) AddUser(id string, userId string) error {
	if r.Session == nil  {
		return errors.New("No client.Session to send the request.")
	}	
	
	_, _, err := r.Session.Request(fmt.Sprintf("/groups/%s/users/%s", id, userId), 
		requests.PUT, 
		nil,  //header
		nil,  // query params
		nil,  // body data
		true)
	if err != nil {
		return err
	}
	
	return nil
}

func (r *Group) CheckUser(id string, userId string) (bool, error) {
	if r.Session == nil  {
		return false, errors.New("No client.Session to send the request.")
	}	
	
	_, _, err := r.Session.Request(fmt.Sprintf("/groups/%s/users/%s", id, userId), 
		requests.HEAD, 
		nil,  //header
		nil,  // query params
		nil,  // body data
		true)
	if err != nil {
		return false, err
	}
	
	return true,nil
}

func (r *Group) DeleteUser(id string, userId string) error {
	if r.Session == nil  {
		return errors.New("No client.Session to send the request.")
	}	
	
	_, _, err := r.Session.Request(fmt.Sprintf("/groups/%s/users/%s", id, userId), 
		requests.DELETE, 
		nil,  //header
		nil,  // query params
		nil,  // body data
		true)
	if err != nil {
		return err
	}
	
	return nil
}