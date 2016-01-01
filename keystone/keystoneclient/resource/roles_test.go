package resource

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"io/ioutil"
	"encoding/json"
	
	"github.com/heartsg/dasea/requests"
	"github.com/heartsg/dasea/keystone/keystoneclient"
	"github.com/heartsg/dasea/keystone/keystoneclient/client"
	"github.com/heartsg/dasea/keystone/keystoneclient/types"
)

func TestRoleCreate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != requests.POST {
			t.Errorf("Expected method %q; got %q", requests.POST, r.Method)
		}
		if r.Header == nil {
			t.Errorf("Expected non-nil request Header")
		}
		if r.Header.Get("User-Agent") != "Dasea Keystone client" {
			t.Errorf("Expected 'User-Agent' == %q; got %q", "Dasea Keystone client", r.Header.Get("User-Agent"))
		}
		if r.Header.Get("X-Auth-Token") != "12345" {
			t.Errorf("Expected 'X-Auth-Token' == %q; got %q", "12345", r.Header.Get("X-Auth-Token"))
		}
		
		body, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		tmp := `{
    "role": {
        "name": "developer"
    }
}`
		tmpStruct := &types.RoleResponse{}
		_ = json.Unmarshal([]byte(tmp), tmpStruct)
		tmpbytes, _ := json.Marshal(tmpStruct)
		if string(body) != string(tmpbytes) {
			t.Errorf("Expected auth request body == %q; got %q", string(tmpbytes), string(body))
		}
		
		w.Write([]byte(`{
    "role": {
        "id": "1e443fa8cee3482a8a2b6954dd5c8f12",
        "links": {
            "self": "http://localhost:5000/v3/roles/1e443fa8cee3482a8a2b6954dd5c8f12"
        },
        "name": "developer"
    }
}`))
	}))	

	defer ts.Close()
	
	auth := &client.AuthAccess {
		Access: &client.AccessInfo {
			Token: "12345",
		},
	}
	
	session := client.NewSession(&keystoneclient.Opts{ AuthUrl: ts.URL }, auth)
	
	roleResource := &Role {
		Session: session,
	}
	role, err := roleResource.Create(&types.Role{
		Name: "developer",
	})
	if err != nil {
		t.Error(err)
	}
	if role.Id != "1e443fa8cee3482a8a2b6954dd5c8f12" {
		t.Errorf("Expect role.Id == %q; got %q", "161718", role.Id)
	}
}

func TestRoleList(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != requests.GET {
			t.Errorf("Expected method %q; got %q", requests.GET, r.Method)
		}
		if r.Header == nil {
			t.Errorf("Expected non-nil request Header")
		}
		if r.Header.Get("User-Agent") != "Dasea Keystone client" {
			t.Errorf("Expected 'User-Agent' == %q; got %q", "Dasea Keystone client", r.Header.Get("User-Agent"))
		}
		if r.Header.Get("X-Auth-Token") != "12345" {
			t.Errorf("Expected 'X-Auth-Token' == %q; got %q", "12345", r.Header.Get("X-Auth-Token"))
		}
		
		//test query parameter
		q := r.URL.Query()
		if q.Get("name") != "test" {
			t.Errorf("Expected 'name' == %q; got %q", "test", q.Get("name"))
		}
		
		w.Write([]byte(`{
    "links": {
        "next": null,
        "previous": null,
        "self": "http://localhost:5000/v3/roles"
    },
    "roles": [
        {
            "id": "5318e65d75574c17bf5339d3df33a5a3",
            "links": {
                "self": "http://localhost:5000/v3/roles/5318e65d75574c17bf5339d3df33a5a3"
            },
            "name": "admin"
        },
        {
            "id": "642bcfc75c384fd181adf34d9b2df897",
            "links": {
                "self": "http://localhost:5000/v3/roles/642bcfc75c384fd181adf34d9b2df897"
            },
            "name": "anotherrole"
        }
	]
}`))
	}))	
	defer ts.Close()
	
	auth := &client.AuthAccess {
		Access: &client.AccessInfo {
			Token: "12345",
		},
	}
	
	session := client.NewSession(&keystoneclient.Opts{ AuthUrl: ts.URL }, auth)
	
	roleResource := &Role {
		Session: session,
	}
	roles, err := roleResource.List("test")
	if err != nil {
		t.Error(err)
	}
	if roles.Links["self"] != "http://localhost:5000/v3/roles" {
		t.Errorf("Expect roles.Links[self] == %q; got %q", "http://localhost:5000/v3/roles", roles.Links["self"])
	}
	if roles.Roles[0].Id != "5318e65d75574c17bf5339d3df33a5a3" {
		t.Errorf("Expect roles.Roles[0].Id == %q; got %q", "5318e65d75574c17bf5339d3df33a5a3", roles.Roles[0].Id)
	}
	if roles.Roles[1].Links["self"] != "http://localhost:5000/v3/roles/642bcfc75c384fd181adf34d9b2df897" {
		t.Errorf("Expect roles.Roles[1].Links[self] == %q; got %q", "http://localhost:5000/v3/roles/642bcfc75c384fd181adf34d9b2df897", roles.Roles[1].Links["self"])
	}
}

func TestRoleDelete (t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != requests.DELETE {
			t.Errorf("Expected method %q; got %q", requests.DELETE, r.Method)
		}
		if r.Header == nil {
			t.Errorf("Expected non-nil request Header")
		}
		if r.Header.Get("User-Agent") != "Dasea Keystone client" {
			t.Errorf("Expected 'User-Agent' == %q; got %q", "Dasea Keystone client", r.Header.Get("User-Agent"))
		}
		if r.Header.Get("X-Auth-Token") != "12345" {
			t.Errorf("Expected 'X-Auth-Token' == %q; got %q", "12345", r.Header.Get("X-Auth-Token"))
		}
		if r.URL.Path != "/roles/xyz" {
			t.Errorf("Expected 'url' == %q; got %q", "/roles/xyz", r.URL.Path)
		}
	}))
	defer ts.Close()
	
	
	auth := &client.AuthAccess {
		Access: &client.AccessInfo {
			Token: "12345",
		},
	}
	
	session := client.NewSession(&keystoneclient.Opts{ AuthUrl: ts.URL }, auth)
	
	roleResource := &Role {
		Session: session,
	}
	err := roleResource.Delete("xyz")
	
	if err != nil {
		t.Error(err)
	}
}

func TestRoleGet (t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != requests.GET {
			t.Errorf("Expected method %q; got %q", requests.GET, r.Method)
		}
		if r.Header == nil {
			t.Errorf("Expected non-nil request Header")
		}
		if r.Header.Get("User-Agent") != "Dasea Keystone client" {
			t.Errorf("Expected 'User-Agent' == %q; got %q", "Dasea Keystone client", r.Header.Get("User-Agent"))
		}
		if r.Header.Get("X-Auth-Token") != "12345" {
			t.Errorf("Expected 'X-Auth-Token' == %q; got %q", "12345", r.Header.Get("X-Auth-Token"))
		}
		if r.URL.Path != "/roles/xyz" {
			t.Errorf("Expected 'url' == %q; got %q", "/roles/xyz", r.URL.Path)
		}
		
		w.Write([]byte(`{
    "role": {
        "id": "1e443fa8cee3482a8a2b6954dd5c8f12",
        "links": {
            "self": "http://localhost:5000/v3/roles/1e443fa8cee3482a8a2b6954dd5c8f12"
        },
        "name": "Developer"
    }
}`))
	}))
	defer ts.Close()
	
	
	auth := &client.AuthAccess {
		Access: &client.AccessInfo {
			Token: "12345",
		},
	}
	
	session := client.NewSession(&keystoneclient.Opts{ AuthUrl: ts.URL }, auth)
	
	roleResource := &Role {
		Session: session,
	}
	role, err := roleResource.Get("xyz")
	
	if err != nil {
		t.Error(err)
	}
	
	if role.Id != "1e443fa8cee3482a8a2b6954dd5c8f12" {
		t.Errorf("Expected id == %q; got %q", "1e443fa8cee3482a8a2b6954dd5c8f12", role.Id)
	}
}

func TestRoleUpdate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != requests.PATCH {
			t.Errorf("Expected method %q; got %q", requests.PATCH, r.Method)
		}
		if r.Header == nil {
			t.Errorf("Expected non-nil request Header")
		}
		if r.Header.Get("User-Agent") != "Dasea Keystone client" {
			t.Errorf("Expected 'User-Agent' == %q; got %q", "Dasea Keystone client", r.Header.Get("User-Agent"))
		}
		if r.Header.Get("X-Auth-Token") != "12345" {
			t.Errorf("Expected 'X-Auth-Token' == %q; got %q", "12345", r.Header.Get("X-Auth-Token"))
		}
		if r.URL.Path != "/roles/xyz" {
			t.Errorf("Expected 'url' == %q; got %q", "/roles/xyz", r.URL.Path)
		}
		
		
		body, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		tmp := `{
    "role": {
        "name": "Developer"
    }
}`
		tmpStruct := &types.RoleResponse{}
		_ = json.Unmarshal([]byte(tmp), tmpStruct)
		tmpbytes, _ := json.Marshal(tmpStruct)

		if string(body) != string(tmpbytes) {
			t.Errorf("Expected auth request body == %q; got %q", string(tmpbytes), string(body))
		}
		
		w.Write([]byte(`{
    "role": {
        "id": "1e443fa8cee3482a8a2b6954dd5c8f12",
        "links": {
            "self": "http://localhost:5000/v3/roles/1e443fa8cee3482a8a2b6954dd5c8f12"
        },
        "name": "Developer"
    }
}`))
	}))	

	defer ts.Close()
	
	auth := &client.AuthAccess {
		Access: &client.AccessInfo {
			Token: "12345",
		},
	}
	
	session := client.NewSession(&keystoneclient.Opts{ AuthUrl: ts.URL }, auth)
	
	roleResource := &Role {
		Session: session,
	}

	role, err := roleResource.Update("xyz", &types.Role{
		Name: "Developer",
	})
	if err != nil {
		t.Error(err)
	}
	if role.Id != "1e443fa8cee3482a8a2b6954dd5c8f12" {
		t.Errorf("Expect role.Id == %q; got %q", "1e443fa8cee3482a8a2b6954dd5c8f12", role.Id)
	}
}



func TestRoleListForUserOnDomain(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != requests.GET {
			t.Errorf("Expected method %q; got %q", requests.GET, r.Method)
		}
		if r.Header == nil {
			t.Errorf("Expected non-nil request Header")
		}
		if r.Header.Get("User-Agent") != "Dasea Keystone client" {
			t.Errorf("Expected 'User-Agent' == %q; got %q", "Dasea Keystone client", r.Header.Get("User-Agent"))
		}
		if r.Header.Get("X-Auth-Token") != "12345" {
			t.Errorf("Expected 'X-Auth-Token' == %q; got %q", "12345", r.Header.Get("X-Auth-Token"))
		}
		if r.URL.Path != "/domains/default/users/xyz/roles" {
			t.Errorf("Expected 'url' == %q; got %q", "/domains/default/users/xyz/roles" , r.URL.Path)
		}
		
		w.Write([]byte(`{
    "links": {
        "next": null,
        "previous": null,
        "self": "http://localhost:5000/v3/roles"
    },
    "roles": [
        {
            "id": "5318e65d75574c17bf5339d3df33a5a3",
            "links": {
                "self": "http://localhost:5000/v3/roles/5318e65d75574c17bf5339d3df33a5a3"
            },
            "name": "admin"
        },
        {
            "id": "642bcfc75c384fd181adf34d9b2df897",
            "links": {
                "self": "http://localhost:5000/v3/roles/642bcfc75c384fd181adf34d9b2df897"
            },
            "name": "anotherrole"
        }
	]
}`))
	}))	
	defer ts.Close()
	
	auth := &client.AuthAccess {
		Access: &client.AccessInfo {
			Token: "12345",
		},
	}
	
	session := client.NewSession(&keystoneclient.Opts{ AuthUrl: ts.URL }, auth)
	
	roleResource := &Role {
		Session: session,
	}
	roles, err := roleResource.ListForUserOnDomain("xyz", "default")
	if err != nil {
		t.Error(err)
	}
	if roles.Links["self"] != "http://localhost:5000/v3/roles" {
		t.Errorf("Expect roles.Links[self] == %q; got %q", "http://localhost:5000/v3/roles", roles.Links["self"])
	}
	if roles.Roles[0].Id != "5318e65d75574c17bf5339d3df33a5a3" {
		t.Errorf("Expect roles.Roles[0].Id == %q; got %q", "5318e65d75574c17bf5339d3df33a5a3", roles.Roles[0].Id)
	}
	if roles.Roles[1].Links["self"] != "http://localhost:5000/v3/roles/642bcfc75c384fd181adf34d9b2df897" {
		t.Errorf("Expect roles.Roles[1].Links[self] == %q; got %q", "http://localhost:5000/v3/roles/642bcfc75c384fd181adf34d9b2df897", roles.Roles[1].Links["self"])
	}
}





func TestRoleAssignments(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != requests.GET {
			t.Errorf("Expected method %q; got %q", requests.GET, r.Method)
		}
		if r.Header == nil {
			t.Errorf("Expected non-nil request Header")
		}
		if r.Header.Get("User-Agent") != "Dasea Keystone client" {
			t.Errorf("Expected 'User-Agent' == %q; got %q", "Dasea Keystone client", r.Header.Get("User-Agent"))
		}
		if r.Header.Get("X-Auth-Token") != "12345" {
			t.Errorf("Expected 'X-Auth-Token' == %q; got %q", "12345", r.Header.Get("X-Auth-Token"))
		}
		if r.URL.Path != "/role_assignments" {
			t.Errorf("Expected 'url' == %q; got %q", "/role_assignments" , r.URL.Path)
		}
		
		q := r.URL.Query()
		if q.Get("user.id") != "xyz" {
			t.Errorf("Expected 'user.id' == %q; got %q", "xyz", q.Get("user.id"))
		}
		if q.Get("scope.project.id") != "abc" {
			t.Errorf("Expected 'scope.project.id' == %q; got %q", "abc", q.Get("scope.project.id"))
		}
		
		w.Write([]byte(`{
    "role_assignments": [
        {
            "links": {
                "assignment": "http://identity:35357/v3/domains/161718/users/313233/roles/123456"
            },
            "role": {
                "id": "123456"
            },
            "scope": {
                "domain": {
                    "id": "161718"
                }
            },
            "user": {
                "id": "313233"
            }
        },
        {
            "group": {
                "id": "101112"
            },
            "links": {
                "assignment": "http://identity:35357/v3/projects/456789/groups/101112/roles/123456"
            },
            "role": {
                "id": "123456"
            },
            "scope": {
                "project": {
                    "id": "456789"
                }
            }
        }
    ],
    "links": {
        "self": "http://identity:35357/v3/role_assignments",
        "previous": null,
        "next": null
    }
}`))
	}))	
	defer ts.Close()
	
	auth := &client.AuthAccess {
		Access: &client.AccessInfo {
			Token: "12345",
		},
	}
	
	session := client.NewSession(&keystoneclient.Opts{ AuthUrl: ts.URL }, auth)
	
	roleResource := &Role {
		Session: session,
	}
	roles, err := roleResource.Assignments(map[string]string {
		"user.id":"xyz",
		"scope.project.id":"abc",
		"effective":"true",
	})
	if err != nil {
		t.Error(err)
	}
	if roles.Links["self"] != "http://identity:35357/v3/role_assignments" {
		t.Errorf("Expect roles.Links[self] == %q; got %q", "http://identity:35357/v3/role_assignments", roles.Links["self"])
	}
	if roles.RoleAssignments[0].Role.Id != "123456" {
		t.Errorf("Expect roles.RoleAssignments[0].Role.Id == %q; got %q", "123456", roles.RoleAssignments[0].Role.Id)
	}
	if roles.RoleAssignments[1].Links["assignment"] != "http://identity:35357/v3/projects/456789/groups/101112/roles/123456" {
		t.Errorf("Expect roles.RoleAssignments[1].Links[assignment] == %q; got %q", "http://identity:35357/v3/projects/456789/groups/101112/roles/123456", roles.RoleAssignments[1].Links["assignment"])
	}
}

