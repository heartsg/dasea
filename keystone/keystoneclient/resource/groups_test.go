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

func TestGroupCreate(t *testing.T) {
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
    "group": {
        "description": "Contract developers",
        "domain_id": "default",
        "name": "Contract developers"
    }
}`
		tmpStruct := &types.GroupResponse{}
		_ = json.Unmarshal([]byte(tmp), tmpStruct)
		tmpbytes, _ := json.Marshal(tmpStruct)
		if string(body) != string(tmpbytes) {
			t.Errorf("Expected auth request body == %q; got %q", string(tmpbytes), string(body))
		}
		
		w.Write([]byte(`{
    "group": {
        "domain_id": "default",
        "description": "Contract developers",
        "id": "c0d675eac29945ad9dfd08aa1bb75751",
        "links": {
            "self": "http://localhost:5000/v3/groups/c0d675eac29945ad9dfd08aa1bb75751"
        },
        "name": "Contract developers"
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
	
	groupResource := &Group {
		Session: session,
	}
	group, err := groupResource.Create(&types.Group{
		Description: "Contract developers",
		DomainId: "default",
		Name: "Contract developers",
	})
	if err != nil {
		t.Error(err)
	}
	if group.Id != "c0d675eac29945ad9dfd08aa1bb75751" {
		t.Errorf("Expect group.Id == %q; got %q", "c0d675eac29945ad9dfd08aa1bb75751", group.Id)
	}
}

func TestGroupList(t *testing.T) {
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
		
		w.Write([]byte(`{
    "links": {
        "self": "http://localhost:5000/v3/groups",
        "previous": null,
        "next": null
    },
    "groups": [
        {
            "domain_id": "default",
            "description": "non-admin group",
            "id": "96372bbb152f475aa37e9a76a25a029c",
            "links": {
                "self": "http://localhost:5000/v3/groups/96372bbb152f475aa37e9a76a25a029c"
            },
            "name": "nonadmins"
        },
        {
            "domain_id": "default",
            "description": "openstack admin group",
            "id": "9ce0ad4e58a84d7a97b92f7955d10c92",
            "links": {
                "self": "http://localhost:5000/v3/groups/9ce0ad4e58a84d7a97b92f7955d10c92"
            },
            "name": "admins"
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
	
	groupResource := &Group {
		Session: session,
	}
	groups, err := groupResource.List("default", "")
	if err != nil {
		t.Error(err)
	}
	if groups.Links["self"] != "http://localhost:5000/v3/groups" {
		t.Errorf("Expect groups.Links[self] == %q; got %q", "http://localhost:5000/v3/groups", groups.Links["self"])
	}
	if groups.Groups[0].Id != "96372bbb152f475aa37e9a76a25a029c" {
		t.Errorf("Expect groups.Groups[0].Id == %q; got %q", "96372bbb152f475aa37e9a76a25a029c", groups.Groups[0].Id)
	}
	if groups.Groups[1].Links["self"] != "http://localhost:5000/v3/groups/9ce0ad4e58a84d7a97b92f7955d10c92" {
		t.Errorf("Expect groups.Groups[1].Links[self] == %q; got %q", "http://localhost:5000/v3/groups/9ce0ad4e58a84d7a97b92f7955d10c92", groups.Groups[1].Links["self"])
	}
}

func TestGroupDelete (t *testing.T) {
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
		if r.URL.Path != "/groups/xyz" {
			t.Errorf("Expected 'url' == %q; got %q", "/groups/xyz", r.URL.Path)
		}
	}))
	defer ts.Close()
	
	
	auth := &client.AuthAccess {
		Access: &client.AccessInfo {
			Token: "12345",
		},
	}
	
	session := client.NewSession(&keystoneclient.Opts{ AuthUrl: ts.URL }, auth)
	
	groupResource := &Group {
		Session: session,
	}
	err := groupResource.Delete("xyz")
	
	if err != nil {
		t.Error(err)
	}
}

func TestGroupGet (t *testing.T) {
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
		if r.URL.Path != "/groups/xyz" {
			t.Errorf("Expected 'url' == %q; got %q", "/groups/xyz", r.URL.Path)
		}
		
		w.Write([]byte(`{
    "group": {
        "domain_id": "default",
        "description": "Contract developers",
        "id": "c0d675eac29945ad9dfd08aa1bb75751",
        "links": {
            "self": "http://localhost:5000/v3/groups/c0d675eac29945ad9dfd08aa1bb75751"
        },
        "name": "Contract developers"
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
	
	groupResource := &Group {
		Session: session,
	}
	group, err := groupResource.Get("xyz")
	
	if err != nil {
		t.Error(err)
	}
	
	if group.Id != "c0d675eac29945ad9dfd08aa1bb75751" {
		t.Errorf("Expected id == %q; got %q", "c0d675eac29945ad9dfd08aa1bb75751", group.Id)
	}
}

func TestGroupUpdate(t *testing.T) {
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
		if r.URL.Path != "/groups/xyz" {
			t.Errorf("Expected 'url' == %q; got %q", "/groups/xyz", r.URL.Path)
		}
		
		
		body, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		tmp := `{
    "group": {
        "description": "Contract developers 2016",
        "name": "Contract developers 2016"
    }
}`
		tmpStruct := &types.GroupResponse{}
		_ = json.Unmarshal([]byte(tmp), tmpStruct)
		tmpbytes, _ := json.Marshal(tmpStruct)

		if string(body) != string(tmpbytes) {
			t.Errorf("Expected auth request body == %q; got %q", string(tmpbytes), string(body))
		}
		
		w.Write([]byte(`{
    "group": {
        "domain_id": "default",
        "description": "Contract developers 2016",
        "id": "c0d675eac29945ad9dfd08aa1bb75751",
        "links": {
            "self": "http://localhost:5000/v3/groups/c0d675eac29945ad9dfd08aa1bb75751"
        },
        "name": "Contract developers 2016"
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
	
	groupResource := &Group {
		Session: session,
	}

	group, err := groupResource.Update("xyz", &types.Group{
		Description: "Contract developers 2016",
		Name: "Contract developers 2016",
	})
	if err != nil {
		t.Error(err)
	}
	if group.Id != "c0d675eac29945ad9dfd08aa1bb75751" {
		t.Errorf("Expect group.Id == %q; got %q", "c0d675eac29945ad9dfd08aa1bb75751", group.Id)
	}
}


func TestGroupListUsers(t *testing.T) {
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
		if r.URL.Path != "/groups/xyz/users" {
			t.Errorf("Expected 'url' == %q; got %q", "/groups/xyz/users", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("domain_id") != "default" {
			t.Errorf("Expected 'domain_id' == %q; got %q", "default", q.Get("domain_id"))
		}
		w.Write([]byte(`{
    "users": [
        {
            "name": "admin",
            "links": {
                "self": "http://localhost:5000/v3/users/fff603a0829d41e48bc0dd0d72ad61ce"
            },
            "domain_id": "default",
            "enabled": true,
            "email": null,
            "id": "fff603a0829d41e48bc0dd0d72ad61ce"
        }
    ],
    "links": {
        "self": "http://localhost:5000/v3/groups/9ce0ad4e58a84d7a97b92f7955d10c92/users",
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
	
	groupResource := &Group {
		Session: session,
	}
	users, err := groupResource.ListUsers("xyz", "default", "", "", "true")
	if err != nil {
		t.Error(err)
	}
	if users.Links["self"] != "http://localhost:5000/v3/groups/9ce0ad4e58a84d7a97b92f7955d10c92/users" {
		t.Errorf("Expect users.Links[self] == %q; got %q", "http://localhost:5000/v3/groups/9ce0ad4e58a84d7a97b92f7955d10c92/users", users.Links["self"])
	}
}



func TestGroupAddUser(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != requests.PUT {
			t.Errorf("Expected method %q; got %q", requests.PUT, r.Method)
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
		if r.URL.Path != "/groups/xyz/users/abc" {
			t.Errorf("Expected 'url' == %q; got %q", "/groups/xyz/users/abc", r.URL.Path)
		}
	}))	
	defer ts.Close()
	
	auth := &client.AuthAccess {
		Access: &client.AccessInfo {
			Token: "12345",
		},
	}
	
	session := client.NewSession(&keystoneclient.Opts{ AuthUrl: ts.URL }, auth)
	
	groupResource := &Group {
		Session: session,
	}
	err := groupResource.AddUser("xyz", "abc")
	if err != nil {
		t.Error(err)
	}
}

func TestGroupCheckUser(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != requests.HEAD {
			t.Errorf("Expected method %q; got %q", requests.HEAD, r.Method)
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
		if r.URL.Path != "/groups/xyz/users/abc" {
			t.Errorf("Expected 'url' == %q; got %q", "/groups/xyz/users/abc", r.URL.Path)
		}
	}))	
	defer ts.Close()
	
	auth := &client.AuthAccess {
		Access: &client.AccessInfo {
			Token: "12345",
		},
	}
	
	session := client.NewSession(&keystoneclient.Opts{ AuthUrl: ts.URL }, auth)
	
	groupResource := &Group {
		Session: session,
	}
	_, err := groupResource.CheckUser("xyz", "abc")
	if err != nil {
		t.Error(err)
	}
}

func TestGroupDeleteUser(t *testing.T) {
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
		if r.URL.Path != "/groups/xyz/users/abc" {
			t.Errorf("Expected 'url' == %q; got %q", "/groups/xyz/users/abc", r.URL.Path)
		}
	}))	
	defer ts.Close()
	
	auth := &client.AuthAccess {
		Access: &client.AccessInfo {
			Token: "12345",
		},
	}
	
	session := client.NewSession(&keystoneclient.Opts{ AuthUrl: ts.URL }, auth)
	
	groupResource := &Group {
		Session: session,
	}
	err := groupResource.DeleteUser("xyz", "abc")
	if err != nil {
		t.Error(err)
	}
}