package resource

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/heartsg/dasea/keystone/keystoneclient/types"
	"github.com/heartsg/dasea/keystone/keystoneclient/client"
	"github.com/heartsg/dasea/requests"
)
type User struct {
	Session *client.Session
}

func (r *User) Create(u *types.User) (*types.User, error) {
	if r.Session == nil  {
		return nil, errors.New("No client.Session to send the request.")
	}
	
	userRequest := types.NewUserRequest(u)
	
	_, body, err := r.Session.Request("/users", 
		requests.POST, 
		nil,  //header
		nil,  // query params
		userRequest,  // body data
		true)
	if err != nil {
		return nil, err
	}
	
	userResponse := &types.UserResponse{}
	err = json.Unmarshal(body, userResponse)
	if err != nil {
		return nil, err
	}
	
	return userResponse.User, nil
}

func (r *User) List(domainId string, name string, enabled string) (*types.UsersResponse, error) {
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
	if enabled == "true" || enabled == "false" {
		queryParams["enabled"] = enabled
	}
	
	_, body, err := r.Session.Request("/users", 
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

func (r *User) Delete(id string) error {
	if r.Session == nil  {
		return errors.New("No client.Session to send the request.")
	}
	
	_, _, err := r.Session.Request(fmt.Sprintf("/users/%s", id), 
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

func (r *User) Get(id string) (*types.User, error) {
	if r.Session == nil  {
		return nil, errors.New("No client.Session to send the request.")
	}
	
	_, body, err := r.Session.Request(fmt.Sprintf("/users/%s", id), 
		requests.GET, 
		nil, //header
		nil,  // query params
		nil,  // body data
		true)
	if err != nil {
		return nil, err
	}

	userResponse := &types.UserResponse{}
	err = json.Unmarshal(body, userResponse)
	if err != nil {
		return nil, err
	}
	
	return userResponse.User, nil
}


func (r *User) Update(id string, u *types.User) (*types.User, error) {
	if r.Session == nil  {
		return nil, errors.New("No client.Session to send the request.")
	}
	
	userRequest := types.NewUserRequest(u)
	
	_, body, err := r.Session.Request(fmt.Sprintf("/users/%s", id), 
		requests.PATCH, 
		nil,  //header
		nil,  // query params
		userRequest,  // body data
		true)
	if err != nil {
		return nil, err
	}
	
	userResponse := &types.UserResponse{}
	err = json.Unmarshal(body, userResponse)
	if err != nil {
		return nil, err
	}
	
	return userResponse.User, nil	
}


func (r *User) ChangePassword(id string, password string, originalPassword string) error {
	if r.Session == nil  {
		return errors.New("No client.Session to send the request.")
	}
	
	userRequest := types.NewUserRequest(&types.User{
		Password: password,
		OriginalPassword: originalPassword,
	})
	
	_, _, err := r.Session.Request(fmt.Sprintf("/users/%s/password", id), 
		requests.POST, 
		nil,  //header
		nil,  // query params
		userRequest,  // body data
		true)
	if err != nil {
		return err
	}
	
	return nil	
}


func (r *User) ListGroups(id string) (*types.GroupsResponse, error) {
	if r.Session == nil  {
		return nil, errors.New("No client.Session to send the request.")
	}
	
	_, body, err := r.Session.Request(fmt.Sprintf("/users/%s/groups", id), 
		requests.GET, 
		nil,  //header
		nil,  // query params
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


func (r *User) ListProjects(id string) (*types.ProjectsResponse, error) {
	if r.Session == nil  {
		return nil, errors.New("No client.Session to send the request.")
	}
	
	_, body, err := r.Session.Request(fmt.Sprintf("/users/%s/projects", id), 
		requests.GET, 
		nil,  //header
		nil,  // query params
		nil,  // body data
		true)
	if err != nil {
		return nil, err
	}
	
	projectsResponse := &types.ProjectsResponse{}
	err = json.Unmarshal(body, projectsResponse)
	if err != nil {
		return nil, err
	}
	
	return projectsResponse, nil	
}