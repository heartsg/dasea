package resource


import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/heartsg/dasea/keystone/keystoneclient/types"
	"github.com/heartsg/dasea/keystone/keystoneclient/client"
	"github.com/heartsg/dasea/requests"
)
type Region struct {
	Session *client.Session
}

func (r *Region) Create(d *types.Region) (*types.Region, error) {
	if r.Session == nil  {
		return nil, errors.New("No client.Session to send the request.")
	}
	
	regionRequest := types.NewRegionRequest(d)
	
	_, body, err := r.Session.Request("/regions", 
		requests.POST, 
		nil,  //header
		nil,  // query params
		regionRequest,  // body data
		true)
	if err != nil {
		return nil, err
	}
	
	regionResponse := &types.RegionResponse{}
	err = json.Unmarshal(body, regionResponse)
	if err != nil {
		return nil, err
	}
	
	return regionResponse.Region, nil
}

func (r *Region) List(parentRegionId string) (*types.RegionsResponse, error) {
	if r.Session == nil  {
		return nil, errors.New("No client.Session to send the request.")
	}
	
	queryParams := make(map[string]string)
	if parentRegionId != "" {
		queryParams["parent_region_id"] = parentRegionId
	}
	
	_, body, err := r.Session.Request("/regions", 
		requests.GET, 
		nil,  //header
		queryParams,  // query params
		nil,  // body data
		true)
	if err != nil {
		return nil, err
	}
	
	regionsResponse := &types.RegionsResponse{}
	err = json.Unmarshal(body, regionsResponse)
	if err != nil {
		return nil, err
	}
	
	return regionsResponse, nil	
}

func (r *Region) Delete(id string) error {
	if r.Session == nil  {
		return errors.New("No client.Session to send the request.")
	}
	
	_, _, err := r.Session.Request(fmt.Sprintf("/regions/%s", id), 
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

func (r *Region) Get(id string) (*types.Region, error) {
	if r.Session == nil  {
		return nil, errors.New("No client.Session to send the request.")
	}
	
	_, body, err := r.Session.Request(fmt.Sprintf("/regions/%s", id), 
		requests.GET, 
		nil, //header
		nil,  // query params
		nil,  // body data
		true)
	if err != nil {
		return nil, err
	}

	regionResponse := &types.RegionResponse{}
	err = json.Unmarshal(body, regionResponse)
	if err != nil {
		return nil, err
	}
	
	return regionResponse.Region, nil
}


func (r *Region) Update(id string, d *types.Region) (*types.Region, error) {
	if r.Session == nil  {
		return nil, errors.New("No client.Session to send the request.")
	}
	
	regionRequest := types.NewRegionRequest(d)
	
	_, body, err := r.Session.Request(fmt.Sprintf("/regions/%s", id), 
		requests.PATCH, 
		nil,  //header
		nil,  // query params
		regionRequest,  // body data
		true)
	if err != nil {
		return nil, err
	}
	
	regionResponse := &types.RegionResponse{}
	err = json.Unmarshal(body, regionResponse)
	if err != nil {
		return nil, err
	}
	
	return regionResponse.Region, nil	
}