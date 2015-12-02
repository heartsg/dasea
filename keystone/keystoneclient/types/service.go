package types

			
type Service struct {
	Description string `json:"description,omitempty"`
	Id string `json:"id,omitempty"`
	Links Link `json:"links,omitempty"`
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

type ServiceRequest struct {
	Service *Service `json:"service"`
}
type ServiceResponse struct {
	Service *Service `json:"service"`
}
type ServicesResponse struct {
	Services []*Service `json:"services"`
	Links Link `json:"links"` 
}

type Endpoint struct {
	RegionId string `json:"region_id,omitempty"`
	Url string `json:"url,omitempty"`
	Region string `json:"region,omitempty"`
	Interface string `json:"interface,omitempty"`
	Id string `json:"id,omitempty"`
	Enabled bool `json:"enabled,omitempty"`
	Links Link `json:"links,omitempty"`
	ServiceId string `json:"service_id,omitempty"`
}

type EndpointRequest struct {
	Endpoint *Endpoint `json:"endpoint"`
}

type EndpointResponse struct {
	Endpoint *Endpoint `json:"endpoint"`
}

type EndpointsResponse struct {
	Endpoints []*Endpoint `json:"endpoints"`
}