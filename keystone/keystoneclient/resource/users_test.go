package resource

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"io/ioutil"
	
	"github.com/heartsg/dasea/requests"
	"github.com/heartsg/dasea/keystone/keystoneclient/client"
	"github.com/heartsg/dasea/keystone/keystoneclient/types"
)

func TestUserCreate(t *testing.T) {
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
		// should be
		//{
    	//	"user": {
		//		"default_project_id": "263fd9",
		//		"description": "James Doe user",
		//		"domain_id": "1789d1",
		//		"email": "jdoe@example.com",
		//		"enabled": true,
		//		"name": "James Doe",
		//		"password": "secretsecret"
		//	}
		//}
		// But sequence not really defined. We try to know the sequence first, and
		// then hard code here.
		if string(body) != `{"user":{"name":"James Doe","domain_id":"1789d1","default_project_id":"263fd9","email":"jdoe@example.com","description":"James Doe user","password":"secretsecret","enabled":true}}` {
			t.Errorf("Expected auth request body == %q; got %q", 
				`{"user":{"name":"James Doe","domain_id":"1789d1","default_project_id":"263fd9","email":"jdoe@example.com","description":"James Doe user","password":"secretsecret","enabled":true}}`,
				string(body))
		}
		
		w.Write([]byte(`{
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
}`))
	}))	

	defer ts.Close()
	
	auth := &client.AuthAccess {
		Access: &client.AccessInfo {
			Token: "12345",
		},
	}
	
	session := client.NewSession(ts.URL, auth)
	
	userResource := &User {
		Session: session,
	}
	user, err := userResource.Create(&types.User{
		DefaultProjectId: "263fd9",
		Description: "James Doe user",
		DomainId: "1789d1",
		Email: "jdoe@example.com",
		Enabled: true,
		Name: "James Doe",
		Password: "secretsecret",
	})
	if err != nil {
		t.Error(err)
	}
	if user.Name != "James Doe" {
		t.Errorf("Expect user.Name == %q; got %q", "James Doe", user.Name)
	}
	if user.Links["self"] != "https://identity:35357/v3/users/ff4e51" {
		t.Errorf("Expect user.Links[self] == %q; got %q", "https://identity:35357/v3/users/ff4e51", user.Links["self"])
	}
}



func TestList(t *testing.T) {
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
		if q.Get("domain_id") != "default" {
			t.Errorf("Expected 'domain_id' == %q; got %q", "default", q.Get("domain_id"))
		}
		if q.Get("enabled") != "true" {
			t.Errorf("Expected 'enabled' == %q; got %q", "true", q.Get("enabled"))
		}
		w.Write([]byte(`{
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
		}`))
	}))	
	defer ts.Close()
	
	auth := &client.AuthAccess {
		Access: &client.AccessInfo {
			Token: "12345",
		},
	}
	
	session := client.NewSession(ts.URL, auth)
	
	userResource := &User {
		Session: session,
	}
	users, err := userResource.List("default", "", "true")
	if err != nil {
		t.Error(err)
	}
	if users.Links["self"] != "http://localhost:5000/v3/users" {
		t.Errorf("Expect users.Links[self] == %q; got %q", "http://localhost:5000/v3/users", users.Links["self"])
	}
	if users.Links["next"] != "" {
		t.Errorf("Expect users.Links[next] empty, got %q", users.Links["next"])
	}
	if users.Links["previous"] != "" {
		t.Errorf("Expect users.Links[previous] empty, got %q", users.Links["previous"])
	}
	if users.Users[0].Name != "glance" {
		t.Errorf("Expect users.Users[0].Name == %q; got %q", "glance", users.Users[0].Name)
	}
	if users.Users[1].Links["self"] != "http://localhost:5000/v3/users/4ab84ab39de54f4d96eaff8f2145a7cd" {
		t.Errorf("Expect users.Users[1].Links[self] == %q; got %q", "http://localhost:5000/v3/users/4ab84ab39de54f4d96eaff8f2145a7cd", users.Users[1].Links["self"])
	}
}

func TestUserDelete (t *testing.T) {
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
		if r.URL.Path != "/users/xyz" {
			t.Errorf("Expected 'url' == %q; got %q", "/users/xyz", r.URL.Path)
		}
	}))
	defer ts.Close()
	
	
	auth := &client.AuthAccess {
		Access: &client.AccessInfo {
			Token: "12345",
		},
	}
	
	session := client.NewSession(ts.URL, auth)
	
	userResource := &User {
		Session: session,
	}
	err := userResource.Delete("xyz")
	
	if err != nil {
		t.Error(err)
	}
}

func TestUserGet (t *testing.T) {
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
		if r.URL.Path != "/users/xyz" {
			t.Errorf("Expected 'url' == %q; got %q", "/users/xyz", r.URL.Path)
		}
		
		w.Write([]byte(`{
    "user": {
        "default_project_id": "263fd9",
        "description": "John Smith's user",
        "domain_id": "1789d1",
        "email": "jsmith@example.com",
        "enabled": true,
        "id": "9fe1d3",
        "links": {
            "self": "https://identity:35357/v3/users/9fe1d3"
        },
        "name": "jsmith"
    }
}`))
	}))
	defer ts.Close()
	
	
	auth := &client.AuthAccess {
		Access: &client.AccessInfo {
			Token: "12345",
		},
	}
	
	session := client.NewSession(ts.URL, auth)
	
	userResource := &User {
		Session: session,
	}
	user, err := userResource.Get("xyz")
	
	if err != nil {
		t.Error(err)
	}
	
	if user.Id != "9fe1d3" {
		t.Errorf("Expected id == %q; got %q", "9fe1d3", user.Id)
	}
}

func TestUserUpdate(t *testing.T) {
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
		if r.URL.Path != "/users/xyz" {
			t.Errorf("Expected 'url' == %q; got %q", "/users/xyz", r.URL.Path)
		}
		
		
		body, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		// should be
		//{
    	//	"user": {
		//		"default_project_id": "263fd9",
		//		"description": "James Doe user",
		//		"email": "jdoe@example.com",
		//		"enabled": true,
		//	}
		//}
		// But sequence not really defined. We try to know the sequence first, and
		// then hard code here.
		if string(body) != `{"user":{"default_project_id":"263fd9","email":"jdoe@example.com","description":"James Doe user","enabled":true}}` {
			t.Errorf("Expected auth request body == %q; got %q", 
				`{"user":"default_project_id":"263fd9","email":"jdoe@example.com","description":"James Doe user","enabled":true}}`,
				string(body))
		}

		
		w.Write([]byte(`{
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
}`))
	}))	

	defer ts.Close()
	
	auth := &client.AuthAccess {
		Access: &client.AccessInfo {
			Token: "12345",
		},
	}
	
	session := client.NewSession(ts.URL, auth)
	
	userResource := &User {
		Session: session,
	}
	user, err := userResource.Update("xyz", &types.User{
		DefaultProjectId: "263fd9",
		Description: "James Doe user",
		Email: "jdoe@example.com",
		Enabled: true,
	})
	if err != nil {
		t.Error(err)
	}
	if user.Name != "James Doe" {
		t.Errorf("Expect user.Name == %q; got %q", "James Doe", user.Name)
	}
	if user.Links["self"] != "https://identity:35357/v3/users/ff4e51" {
		t.Errorf("Expect user.Links[self] == %q; got %q", "https://identity:35357/v3/users/ff4e51", user.Links["self"])
	}
}


func TestUserChangePassword(t *testing.T) {
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
		if r.URL.Path != "/users/xyz/password" {
			t.Errorf("Expected 'url' == %q; got %q", "/users/xyz/password", r.URL.Path)
		}
		
		
		body, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if string(body) != `{"user":{"password":"123","original_password":"456"}}` {
			t.Errorf("Expected auth request body == %q; got %q", 
				`{"user":{"password":"123","original_password":"456"}}`,
				string(body))
		}
	}))	

	defer ts.Close()
	
	auth := &client.AuthAccess {
		Access: &client.AccessInfo {
			Token: "12345",
		},
	}
	
	session := client.NewSession(ts.URL, auth)
	
	userResource := &User {
		Session: session,
	}
	err := userResource.ChangePassword("xyz", "123", "456")
	if err != nil {
		t.Error(err)
	}
}



func TestUserListGroups(t *testing.T) {
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
		if r.URL.Path != "/users/xyz/groups" {
			t.Errorf("Expected 'url' == %q; got %q", "/users/xyz/groups", r.URL.Path)
		}
		
		w.Write([]byte(`{
    "groups": [
        {
            "description": "Developers cleared for work on all general projects",
            "domain_id": "1789d1",
            "id": "ea167b",
            "links": {
                "self": "https://identity:35357/v3/groups/ea167b"
            },
            "name": "Developers"
        },
        {
            "description": "Developers cleared for work on secret projects",
            "domain_id": "1789d1",
            "id": "a62db1",
            "links": {
                "self": "https://identity:35357/v3/groups/a62db1"
            },
            "name": "Secure Developers"
        }
    ],
    "links": {
        "self": "http://identity:35357/v3/users/9fe1d3/groups",
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
	
	session := client.NewSession(ts.URL, auth)
	
	userResource := &User {
		Session: session,
	}
	groups, err := userResource.ListGroups("xyz")
	if err != nil {
		t.Error(err)
	}
	if groups.Links["self"] != "http://identity:35357/v3/users/9fe1d3/groups" {
		t.Errorf("Expect groups.Links[self] == %q; got %q", "http://identity:35357/v3/users/9fe1d3/groups", groups.Links["self"])
	}
	if groups.Groups[0].Name != "Developers" {
		t.Errorf("Expect groups.Groups[0].Name == %q; got %q", "Developers", groups.Groups[0].Name)
	}
	if groups.Groups[1].Links["self"] != "https://identity:35357/v3/groups/a62db1" {
		t.Errorf("Expect groups.Groups[1].Links[self] == %q; got %q", "https://identity:35357/v3/groups/a62db1", groups.Groups[1].Links["self"])
	}
}


func TestUserListProjects(t *testing.T) {
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
		if r.URL.Path != "/users/xyz/projects" {
			t.Errorf("Expected 'url' == %q; got %q", "/users/xyz/projects", r.URL.Path)
		}
		
		w.Write([]byte(`{
    "projects": [
        {
            "description": "description of this project",
            "domain_id": "161718",
            "enabled": true,
            "id": "456788",
            "parent_id": "212223",
            "links": {
                "self": "http://identity:35357/v3/projects/456788"
            },
            "name": "a project name"
        },
        {
            "description": "description of this project",
            "domain_id": "161718",
            "enabled": true,
            "id": "456789",
            "parent_id": "212223",
            "links": {
                "self": "http://identity:35357/v3/projects/456789"
            },
            "name": "another domain"
        }
    ],
    "links": {
        "self": "http://identity:35357/v3/users/313233/projects",
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
	
	session := client.NewSession(ts.URL, auth)
	
	userResource := &User {
		Session: session,
	}
	projects, err := userResource.ListProjects("xyz")
	if err != nil {
		t.Error(err)
	}
	if projects.Links["self"] != "http://identity:35357/v3/users/313233/projects" {
		t.Errorf("Expect projects.Links[self] == %q; got %q", "http://identity:35357/v3/users/313233/projects", projects.Links["self"])
	}
	if projects.Projects[0].Name != "a project name" {
		t.Errorf("Expect projects.Projects[0].Name == %q; got %q", "a project name", projects.Projects[0].Name)
	}
	if projects.Projects[1].Links["self"] != "http://identity:35357/v3/projects/456789" {
		t.Errorf("Expect projects.Projects[1].Links[self] == %q; got %q", "http://identity:35357/v3/projects/456789", projects.Projects[1].Links["self"])
	}
}