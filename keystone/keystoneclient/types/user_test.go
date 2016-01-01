package types

import (
	"testing"
	"encoding/json"
	"github.com/heartsg/dasea/testutil"
)


func TestUserRequest(t *testing.T) {
	//request 1 : unscoped password authentication
	request1Raw := `{"user":{"name":"James Doe","domain_id":"1789d1","default_project_id":"263fd9","email":"jdoe@example.com","description":"James Doe user","password":"secretsecret","enabled":true}}`
	request1Struct := &UserRequest{
		User: &User{
			DefaultProjectId: "263fd9",
			Description: "James Doe user",
			DomainId: "1789d1",
			Email: "jdoe@example.com",
			Enabled: true,
			Name: "James Doe",
			Password: "secretsecret",
		},
	}
	
	request1Unmarshal := &UserRequest{}
	err := json.Unmarshal([]byte(request1Raw), request1Unmarshal)
	
	testutil.IsNil(t, err)
	testutil.Equals(t, request1Struct, request1Unmarshal)
	
	request1Marshal, err := json.Marshal(request1Struct)
	testutil.IsNil(t, err)
	testutil.Equals(t, request1Raw, string(request1Marshal))
}

func TestUserResponse(t *testing.T) {
	response1Raw := `{
    "user": {
        "default_project_id": "263fd9",
        "description": "James Doe user",
        "domain_id": "1789d1",
        "email": "jdoe@example.com",
        "enabled": true,
        "id": "ff4e51",
        "links": {
            "self": "https://identity:35357/v3/users/ff4e51"
        },
        "name": "James Doe"
    }
}`
	response1Struct := &UserResponse{
		User: &User {
			DefaultProjectId: "263fd9",
			Description: "James Doe user",
			DomainId: "1789d1",
			Email: "jdoe@example.com",
			Enabled: true,
			Name: "James Doe",
			Id: "ff4e51",
			Links: Link {
				"self": "https://identity:35357/v3/users/ff4e51",
			},
		},
	}
	
	response1Unmarshal := &UserResponse{}
	err := json.Unmarshal([]byte(response1Raw), response1Unmarshal)
	testutil.IsNil(t, err)
	testutil.Equals(t, response1Struct, response1Unmarshal)
	
	
	//test link response
	response2Raw := `{
		"links": {
			"next": null,
			"previous": null,
			"self": "http://localhost:5000/v3/users"
		},
		"users": [
			{
				"domain_id": "default",
				"email": null,
				"enabled": true,
				"id": "2844b2a08be147a08ef58317d6471f1f",
				"links": {
					"self": "http://localhost:5000/v3/users/2844b2a08be147a08ef58317d6471f1f"
				},
				"name": "glance"
			},
			{
				"domain_id": "default",
				"email": "test@example.com",
				"enabled": true,
				"id": "4ab84ab39de54f4d96eaff8f2145a7cd",
				"links": {
					"self": "http://localhost:5000/v3/users/4ab84ab39de54f4d96eaff8f2145a7cd"
				},
				"name": "swiftusertest1"
			}
		]
	}`
	response2Struct := &UsersResponse{
		Links: Link {
			"next": "",
			"previous": "",
			"self": "http://localhost:5000/v3/users",
		},
		Users: []*User {
			&User {
				DomainId: "default",
				Email: "",
				Enabled: true,
				Name: "glance",
				Id: "2844b2a08be147a08ef58317d6471f1f",
				Links: Link {
					"self": "http://localhost:5000/v3/users/2844b2a08be147a08ef58317d6471f1f",
				},
			},
			&User {
				DomainId: "default",
				Email: "test@example.com",
				Enabled: true,
				Name: "swiftusertest1",
				Id: "4ab84ab39de54f4d96eaff8f2145a7cd",
				Links: Link {
					"self": "http://localhost:5000/v3/users/4ab84ab39de54f4d96eaff8f2145a7cd",
				},
			},
		},
	}
	
	response2Unmarshal := &UsersResponse{}
	err = json.Unmarshal([]byte(response2Raw), response2Unmarshal)
	testutil.IsNil(t, err)
	testutil.Equals(t, response2Struct, response2Unmarshal)
	
}
	