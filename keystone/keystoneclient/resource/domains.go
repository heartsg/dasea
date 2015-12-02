package resource


import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/heartsg/dasea/keystone/keystoneclient/types"
	"github.com/heartsg/dasea/keystone/keystoneclient/client"
	"github.com/heartsg/dasea/requests"
)
type Domain struct {
	Session *client.Session
}

func (r *Domain) Create(d *types.Domain) (*types.Domain, error) {
	if r.Session == nil  {
		return nil, errors.New("No client.Session to send the request.")
	}
	
	domainRequest := types.NewDomainRequest(d)
	
	_, body, err := r.Session.Request("/domains", 
		requests.POST, 
		nil,  //header
		nil,  // query params
		domainRequest,  // body data
		true)
	if err != nil {
		return nil, err
	}
	
	domainResponse := &types.DomainResponse{}
	err = json.Unmarshal(body, domainResponse)
	if err != nil {
		return nil, err
	}
	
	return domainResponse.Domain, nil
}

func (r *Domain) List(name string, enabled string) (*types.DomainsResponse, error) {
	if r.Session == nil  {
		return nil, errors.New("No client.Session to send the request.")
	}
	
	queryParams := make(map[string]string)
	if name != "" {
		queryParams["name"] = name
	}
	if enabled == "true" || enabled == "false" {
		queryParams["enabled"] = enabled
	}
	
	_, body, err := r.Session.Request("/domains", 
		requests.GET, 
		nil,  //header
		queryParams,  // query params
		nil,  // body data
		true)
	if err != nil {
		return nil, err
	}
	
	domainsResponse := &types.DomainsResponse{}
	err = json.Unmarshal(body, domainsResponse)
	if err != nil {
		return nil, err
	}
	
	return domainsResponse, nil	
}

func (r *Domain) Delete(id string) error {
	if r.Session == nil  {
		return errors.New("No client.Session to send the request.")
	}
	
	_, _, err := r.Session.Request(fmt.Sprintf("/domains/%s", id), 
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

func (r *Domain) Get(id string) (*types.Domain, error) {
	if r.Session == nil  {
		return nil, errors.New("No client.Session to send the request.")
	}
	
	_, body, err := r.Session.Request(fmt.Sprintf("/domains/%s", id), 
		requests.GET, 
		nil, //header
		nil,  // query params
		nil,  // body data
		true)
	if err != nil {
		return nil, err
	}

	domainResponse := &types.DomainResponse{}
	err = json.Unmarshal(body, domainResponse)
	if err != nil {
		return nil, err
	}
	
	return domainResponse.Domain, nil
}


func (r *Domain) Update(id string, d *types.Domain) (*types.Domain, error) {
	if r.Session == nil  {
		return nil, errors.New("No client.Session to send the request.")
	}
	
	domainRequest := types.NewDomainRequest(d)
	
	_, body, err := r.Session.Request(fmt.Sprintf("/domains/%s", id), 
		requests.PATCH, 
		nil,  //header
		nil,  // query params
		domainRequest,  // body data
		true)
	if err != nil {
		return nil, err
	}
	
	domainResponse := &types.DomainResponse{}
	err = json.Unmarshal(body, domainResponse)
	if err != nil {
		return nil, err
	}
	
	return domainResponse.Domain, nil	
}