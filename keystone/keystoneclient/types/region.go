package types

type Region struct {
	Id string `json:"id,omitempty"`
	Description string `json:"description,omitempty"`
	ParentRegionId string `json:"parent_region_id,omitempty"`
	Links Link `json:"links,omitempty"`
}

type RegionRequest struct {
	Region *Region `json:"region"`
}

type RegionResponse struct {
	Region *Region `json:"region"`
}

type RegionsResponse struct {
	Regions []*Region `json:"regions"`
	Links Link `json:"links"`
}

func NewRegionRequest(r *Region) *RegionRequest {
	return &RegionRequest {
		Region: r,
	}
}