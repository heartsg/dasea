package types

type Group struct {
	Id string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Links Link `json:"links,omitempty"`
	Description string `json:"description,omitempty"`
	DomainId string `json:"domain_id,omitempty"`
}

type GroupRequest struct {
	Group *Group `json:"group"`
}
type GroupResponse struct {
	Group *Group `json:"group"`
}
type GroupsResponse struct {
	Groups []*Group `json:"groups"`
	Links Link `json:"links"`
}

func NewGroupRequest(g *Group) *GroupRequest {
	return &GroupRequest {
		Group: g,
	}
}