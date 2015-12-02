package resource


import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/heartsg/dasea/keystone/keystoneclient/types"
	"github.com/heartsg/dasea/keystone/keystoneclient/client"
	"github.com/heartsg/dasea/requests"
)
type Role struct {
	Session *client.Session
}

func (r *Role) Create(d *types.Role) (*types.Role, error) {
	if r.Session == nil  {
		return nil, errors.New("No client.Session to send the request.")
	}
	
	roleRequest := types.NewRoleRequest(d)
	
	_, body, err := r.Session.Request("/roles", 
		requests.POST, 
		nil,  //header
		nil,  // query params
		roleRequest,  // body data
		true)
	if err != nil {
		return nil, err
	}
	
	roleResponse := &types.RoleResponse{}
	err = json.Unmarshal(body, roleResponse)
	if err != nil {
		return nil, err
	}
	
	return roleResponse.Role, nil
}

func (r *Role) List(name string) (*types.RolesResponse, error) {
	if r.Session == nil  {
		return nil, errors.New("No client.Session to send the request.")
	}
	
	queryParams := make(map[string]string)
	if name != "" {
		queryParams["name"] = name
	}
	
	_, body, err := r.Session.Request("/roles", 
		requests.GET, 
		nil,  //header
		queryParams,  // query params
		nil,  // body data
		true)
	if err != nil {
		return nil, err
	}
	
	rolesResponse := &types.RolesResponse{}
	err = json.Unmarshal(body, rolesResponse)
	if err != nil {
		return nil, err
	}
	
	return rolesResponse, nil	
}

func (r *Role) Delete(id string) error {
	if r.Session == nil  {
		return errors.New("No client.Session to send the request.")
	}
	
	_, _, err := r.Session.Request(fmt.Sprintf("/roles/%s", id), 
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

func (r *Role) Get(id string) (*types.Role, error) {
	if r.Session == nil  {
		return nil, errors.New("No client.Session to send the request.")
	}
	
	_, body, err := r.Session.Request(fmt.Sprintf("/roles/%s", id), 
		requests.GET, 
		nil, //header
		nil,  // query params
		nil,  // body data
		true)
	if err != nil {
		return nil, err
	}

	roleResponse := &types.RoleResponse{}
	err = json.Unmarshal(body, roleResponse)
	if err != nil {
		return nil, err
	}
	
	return roleResponse.Role, nil
}


func (r *Role) Update(id string, d *types.Role) (*types.Role, error) {
	if r.Session == nil  {
		return nil, errors.New("No client.Session to send the request.")
	}
	
	roleRequest := types.NewRoleRequest(d)
	
	_, body, err := r.Session.Request(fmt.Sprintf("/roles/%s", id), 
		requests.PATCH, 
		nil,  //header
		nil,  // query params
		roleRequest,  // body data
		true)
	if err != nil {
		return nil, err
	}
	
	roleResponse := &types.RoleResponse{}
	err = json.Unmarshal(body, roleResponse)
	if err != nil {
		return nil, err
	}
	
	return roleResponse.Role, nil	
}



func (r *Role) ListForUserOnDomain(userId string, domainId string) (*types.RolesResponse, error) {
	if r.Session == nil  {
		return nil, errors.New("No client.Session to send the request.")
	}
	
	_, body, err := r.Session.Request(fmt.Sprintf("/domains/%s/users/%s/roles", domainId, userId), 
		requests.GET, 
		nil,  //header
		nil,  // query params
		nil,  // body data
		true)
	if err != nil {
		return nil, err
	}
	
	rolesResponse := &types.RolesResponse{}
	err = json.Unmarshal(body, rolesResponse)
	if err != nil {
		return nil, err
	}
	
	return rolesResponse, nil	
}


func (r *Role) GrantToUserOnDomain(roleId string, userId string, domainId string) error {
	if r.Session == nil  {
		return errors.New("No client.Session to send the request.")
	}
	
	_, _, err := r.Session.Request(fmt.Sprintf("/domains/%s/users/%s/roles/%s", domainId, userId, roleId), 
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

func (r *Role) CheckUserOnDomain(roleId string, userId string, domainId string) (bool, error) {
	if r.Session == nil  {
		return false, errors.New("No client.Session to send the request.")
	}
	
	_, _, err := r.Session.Request(fmt.Sprintf("/domains/%s/users/%s/roles/%s", domainId, userId, roleId), 
		requests.HEAD, 
		nil,  //header
		nil,  // query params
		nil,  // body data
		true)
	if err != nil {
		return false, err
	}
	
	return true, nil	
}

func (r *Role) DeleteFromUserOnDomain(roleId string, userId string, domainId string) error {
	if r.Session == nil  {
		return errors.New("No client.Session to send the request.")
	}
	
	_, _, err := r.Session.Request(fmt.Sprintf("/domains/%s/users/%s/roles/%s", domainId, userId, roleId), 
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



func (r *Role) ListForGroupOnDomain(groupId string, domainId string) (*types.RolesResponse, error) {
	if r.Session == nil  {
		return nil, errors.New("No client.Session to send the request.")
	}
	
	_, body, err := r.Session.Request(fmt.Sprintf("/domains/%s/groups/%s/roles", domainId, groupId), 
		requests.GET, 
		nil,  //header
		nil,  // query params
		nil,  // body data
		true)
	if err != nil {
		return nil, err
	}
	
	rolesResponse := &types.RolesResponse{}
	err = json.Unmarshal(body, rolesResponse)
	if err != nil {
		return nil, err
	}
	
	return rolesResponse, nil	
}


func (r *Role) GrantToGroupOnDomain(roleId string, groupId string, domainId string) error {
	if r.Session == nil  {
		return errors.New("No client.Session to send the request.")
	}
	
	_, _, err := r.Session.Request(fmt.Sprintf("/domains/%s/groups/%s/roles/%s", domainId, groupId, roleId), 
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

func (r *Role) CheckGroupOnDomain(roleId string, groupId string, domainId string) (bool, error) {
	if r.Session == nil  {
		return false, errors.New("No client.Session to send the request.")
	}
	
	_, _, err := r.Session.Request(fmt.Sprintf("/domains/%s/groups/%s/roles/%s", domainId, groupId, roleId), 
		requests.HEAD, 
		nil,  //header
		nil,  // query params
		nil,  // body data
		true)
	if err != nil {
		return false, err
	}
	
	return true, nil	
}

func (r *Role) DeleteFromGroupOnDomain(roleId string, groupId string, domainId string) error {
	if r.Session == nil  {
		return errors.New("No client.Session to send the request.")
	}
	
	_, _, err := r.Session.Request(fmt.Sprintf("/domains/%s/groups/%s/roles/%s", domainId, groupId, roleId), 
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



func (r *Role) ListForUserOnProject(userId string, projectId string) (*types.RolesResponse, error) {
	if r.Session == nil  {
		return nil, errors.New("No client.Session to send the request.")
	}
	
	_, body, err := r.Session.Request(fmt.Sprintf("/projects/%s/users/%s/roles", projectId, userId), 
		requests.GET, 
		nil,  //header
		nil,  // query params
		nil,  // body data
		true)
	if err != nil {
		return nil, err
	}
	
	rolesResponse := &types.RolesResponse{}
	err = json.Unmarshal(body, rolesResponse)
	if err != nil {
		return nil, err
	}
	
	return rolesResponse, nil	
}


func (r *Role) GrantToUserOnProject(roleId string, userId string, projectId string) error {
	if r.Session == nil  {
		return errors.New("No client.Session to send the request.")
	}
	
	_, _, err := r.Session.Request(fmt.Sprintf("/projects/%s/users/%s/roles/%s", projectId, userId, roleId), 
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

func (r *Role) CheckUserOnProject(roleId string, userId string, projectId string) (bool, error) {
	if r.Session == nil  {
		return false, errors.New("No client.Session to send the request.")
	}
	
	_, _, err := r.Session.Request(fmt.Sprintf("/projects/%s/users/%s/roles/%s", projectId, userId, roleId), 
		requests.HEAD, 
		nil,  //header
		nil,  // query params
		nil,  // body data
		true)
	if err != nil {
		return false, err
	}
	
	return true, nil	
}

func (r *Role) DeleteFromUserOnProject(roleId string, userId string, projectId string) error {
	if r.Session == nil  {
		return errors.New("No client.Session to send the request.")
	}
	
	_, _, err := r.Session.Request(fmt.Sprintf("/projects/%s/users/%s/roles/%s", projectId, userId, roleId), 
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



func (r *Role) ListForGroupOnProject(groupId string, projectId string) (*types.RolesResponse, error) {
	if r.Session == nil  {
		return nil, errors.New("No client.Session to send the request.")
	}
	
	_, body, err := r.Session.Request(fmt.Sprintf("/projects/%s/groups/%s/roles", projectId, groupId), 
		requests.GET, 
		nil,  //header
		nil,  // query params
		nil,  // body data
		true)
	if err != nil {
		return nil, err
	}
	
	rolesResponse := &types.RolesResponse{}
	err = json.Unmarshal(body, rolesResponse)
	if err != nil {
		return nil, err
	}
	
	return rolesResponse, nil	
}


func (r *Role) GrantToGroupOnProject(roleId string, groupId string, projectId string) error {
	if r.Session == nil  {
		return errors.New("No client.Session to send the request.")
	}
	
	_, _, err := r.Session.Request(fmt.Sprintf("/projects/%s/groups/%s/roles/%s", projectId, groupId, roleId), 
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

func (r *Role) CheckGroupOnProject(roleId string, groupId string, projectId string) (bool, error) {
	if r.Session == nil  {
		return false, errors.New("No client.Session to send the request.")
	}
	
	_, _, err := r.Session.Request(fmt.Sprintf("/projects/%s/groups/%s/roles/%s", projectId, groupId, roleId), 
		requests.HEAD, 
		nil,  //header
		nil,  // query params
		nil,  // body data
		true)
	if err != nil {
		return false, err
	}
	
	return true, nil	
}

func (r *Role) DeleteFromGroupOnProject(roleId string, groupId string, projectId string) error {
	if r.Session == nil  {
		return errors.New("No client.Session to send the request.")
	}
	
	_, _, err := r.Session.Request(fmt.Sprintf("/projects/%s/groups/%s/roles/%s", projectId, groupId, roleId), 
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


// Supported input
//
// group.id (Optional) 	query 	xsd:string 	
// Filters the response by a group ID. To list all role assignments for a group, specify group.id={group_id}.
//
// role.id (Optional) 	query 	xsd:string 	
// Filters the response by a role ID. To list all role assignments for a role, specify role.id={role_id}.
//
// scope.domain.id (Optional) 	query 	xsd:string 	
// Filters the response by a domain ID. To list all role assignments for a domain, specify scope.domain.id={domain_id}.
//
// scope.project.id (Optional) 	query 	xsd:string 	
// Filters the response by a project ID. To list all role assignments for a project, specify scope.project.id={project_id}.
//
// user.id (Optional) 	query 	xsd:string 	
// Filters the response by a user ID. To list all role assignments for a user, specify user.id={user_id}.
//
// effective (Optional) 	query 	xsd:key 	
// Lists effective assignments at the user, project, and domain level, allowing for the effects of group membership.
// The group role assignment entities themselves are not returned in the collection.
// This represents the effective role assignments that would be included in a scoped token. You can use the other query parameters with the effective parameter.

// include_subtree (Optional) 	query 	xsd:boolean 	
// (Since v3.6) Lists all role assignments within a tree of projects. The following call lists all role assignments for a project and its sub-projects:
// GET /role_assignments?scope.project.id={project_id}?include_subtree=true
// You can specify include_subtree=true only in combination with scope.project.id. If you do not include the project ID, this call returns the Bad Request (400) response code.
// Each role assignment entity in the collection contains a link to the assignment that created the entity. 
func (r *Role) Assignments(params map[string]string) (*types.RoleAssignmentsResponse, error) {
	if r.Session == nil  {
		return nil, errors.New("No client.Session to send the request.")
	}
	_, body, err := r.Session.Request("/role_assignments", 
		requests.GET, 
		nil,  //header
		params,  // query params
		nil,  // body data
		true)
	if err != nil {
		return nil, err
	}
	roleAssignmentsResponse := &types.RoleAssignmentsResponse{}
	err = json.Unmarshal(body, roleAssignmentsResponse)
	if err != nil {
		return nil, err
	}
	
	return roleAssignmentsResponse, nil
}
