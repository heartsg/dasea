package resource


import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/heartsg/dasea/keystone/keystoneclient/types"
	"github.com/heartsg/dasea/keystone/keystoneclient/client"
	"github.com/heartsg/dasea/requests"
)
type Project struct {
	Session *client.Session
}

func (r *Project) Create(d *types.Project) (*types.Project, error) {
	if r.Session == nil  {
		return nil, errors.New("No client.Session to send the request.")
	}
	
	projectRequest := types.NewProjectRequest(d)
	
	_, body, err := r.Session.Request("/projects", 
		requests.POST, 
		nil,  //header
		nil,  // query params
		projectRequest,  // body data
		true)
	if err != nil {
		return nil, err
	}
	
	projectResponse := &types.ProjectResponse{}
	err = json.Unmarshal(body, projectResponse)
	if err != nil {
		return nil, err
	}
	
	return projectResponse.Project, nil
}

func (r *Project) List(domainId string, parentId string, name string, enabled string) (*types.ProjectsResponse, error) {
	if r.Session == nil  {
		return nil, errors.New("No client.Session to send the request.")
	}
	
	queryParams := make(map[string]string)
	if domainId != "" {
		queryParams["domain_id"] = domainId
	}
	if parentId != "" {
		queryParams["parent_id"] = parentId
	}
	if name != "" {
		queryParams["name"] = name
	}
	if enabled == "true" || enabled == "false" {
		queryParams["enabled"] = enabled
	}
	
	_, body, err := r.Session.Request("/projects", 
		requests.GET, 
		nil,  //header
		queryParams,  // query params
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

func (r *Project) Delete(id string) error {
	if r.Session == nil  {
		return errors.New("No client.Session to send the request.")
	}
	
	_, _, err := r.Session.Request(fmt.Sprintf("/projects/%s", id), 
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

func (r *Project) Get(id string) (*types.Project, error) {
	if r.Session == nil  {
		return nil, errors.New("No client.Session to send the request.")
	}
	
	_, body, err := r.Session.Request(fmt.Sprintf("/projects/%s", id), 
		requests.GET, 
		nil, //header
		nil,  // query params
		nil,  // body data
		true)
	if err != nil {
		return nil, err
	}

	projectResponse := &types.ProjectResponse{}
	err = json.Unmarshal(body, projectResponse)
	if err != nil {
		return nil, err
	}
	
	return projectResponse.Project, nil
}


func (r *Project) Update(id string, d *types.Project) (*types.Project, error) {
	if r.Session == nil  {
		return nil, errors.New("No client.Session to send the request.")
	}
	
	projectRequest := types.NewProjectRequest(d)
	
	_, body, err := r.Session.Request(fmt.Sprintf("/projects/%s", id), 
		requests.PATCH, 
		nil,  //header
		nil,  // query params
		projectRequest,  // body data
		true)
	if err != nil {
		return nil, err
	}
	
	projectResponse := &types.ProjectResponse{}
	err = json.Unmarshal(body, projectResponse)
	if err != nil {
		return nil, err
	}
	
	return projectResponse.Project, nil	
}