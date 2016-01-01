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

func TestProjectCreate(t *testing.T) {
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
    "project": {
        "description": "My new project",
        "domain_id": "default",
        "enabled": true,
        "is_domain": true,
        "name": "myNewProject"
    }
}`
		tmpStruct := &types.ProjectResponse{}
		_ = json.Unmarshal([]byte(tmp), tmpStruct)
		tmpbytes, _ := json.Marshal(tmpStruct)
		if string(body) != string(tmpbytes) {
			t.Errorf("Expected auth request body == %q; got %q", string(tmpbytes), string(body))
		}
		
		w.Write([]byte(`{
    "project": {
        "is_domain": true,
        "description": "My new project",
        "links": {
            "self": "http://localhost:5000/v3/projects/93ebbcc35335488b96ff9cd7d18cbb2e"
        },
        "enabled": true,
        "id": "93ebbcc35335488b96ff9cd7d18cbb2e",
        "parent_id": null,
        "domain_id": "default",
        "name": "myNewProject"
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
	
	projectResource := &Project {
		Session: session,
	}
	project, err := projectResource.Create(&types.Project{
		Description: "My new project",
        DomainId: "default",
        Enabled: true,
        IsDomain: true,
        Name: "myNewProject",
	})
	if err != nil {
		t.Error(err)
	}
	if project.Id != "93ebbcc35335488b96ff9cd7d18cbb2e" {
		t.Errorf("Expect project.Id == %q; got %q", "93ebbcc35335488b96ff9cd7d18cbb2e", project.Id)
	}
}

func TestProjectList(t *testing.T) {
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
		if q.Get("parent_id") != "test" {
			t.Errorf("Expected 'parent_id' == %q; got %q", "test", q.Get("parent_id"))
		}
		if q.Get("name") != "project" {
			t.Errorf("Expected 'name' == %q; got %q", "project", q.Get("name"))
		}
		if q.Get("enabled") != "true" {
			t.Errorf("Expected 'enabled' == %q; got %q", "true", q.Get("enabled"))
		}
		
		w.Write([]byte(`{
    "links": {
        "next": null,
        "previous": null,
        "self": "http://localhost:5000/v3/projects"
    },
    "projects": [
        {
            "description": null,
            "domain_id": "default",
            "enabled": true,
            "id": "0c4e939acacf4376bdcd1129f1a054ad",
            "links": {
                "self": "http://localhost:5000/v3/projects/0c4e939acacf4376bdcd1129f1a054ad"
            },
            "name": "admin",
            "parent_id": null
        },
        {
            "description": null,
            "domain_id": "default",
            "enabled": true,
            "id": "0cbd49cbf76d405d9c86562e1d579bd3",
            "links": {
                "self": "http://localhost:5000/v3/projects/0cbd49cbf76d405d9c86562e1d579bd3"
            },
            "name": "demo",
            "parent_id": null
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
	
	projectResource := &Project {
		Session: session,
	}
	projects, err := projectResource.List("default", "test", "project", "true")
	if err != nil {
		t.Error(err)
	}
	if projects.Links["self"] != "http://localhost:5000/v3/projects" {
		t.Errorf("Expect projects.Links[self] == %q; got %q", "http://localhost:5000/v3/projects", projects.Links["self"])
	}
	if projects.Projects[0].Id != "0c4e939acacf4376bdcd1129f1a054ad" {
		t.Errorf("Expect projects.Projects[0].Id == %q; got %q", "0c4e939acacf4376bdcd1129f1a054ad", projects.Projects[0].Id)
	}
	if projects.Projects[1].Links["self"] != "http://localhost:5000/v3/projects/0cbd49cbf76d405d9c86562e1d579bd3" {
		t.Errorf("Expect projects.Projects[1].Links[self] == %q; got %q", "http://localhost:5000/v3/projects/0cbd49cbf76d405d9c86562e1d579bd3", projects.Projects[1].Links["self"])
	}
}

func TestProjectDelete (t *testing.T) {
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
		if r.URL.Path != "/projects/xyz" {
			t.Errorf("Expected 'url' == %q; got %q", "/projects/xyz", r.URL.Path)
		}
	}))
	defer ts.Close()
	
	
	auth := &client.AuthAccess {
		Access: &client.AccessInfo {
			Token: "12345",
		},
	}
	
	session := client.NewSession(&keystoneclient.Opts{ AuthUrl: ts.URL }, auth)
	
	projectResource := &Project {
		Session: session,
	}
	err := projectResource.Delete("xyz")
	
	if err != nil {
		t.Error(err)
	}
}

func TestProjectGet (t *testing.T) {
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
		if r.URL.Path != "/projects/xyz" {
			t.Errorf("Expected 'url' == %q; got %q", "/projects/xyz", r.URL.Path)
		}
		
		w.Write([]byte(`{
    "project": {
        "description": null,
        "domain_id": "default",
        "enabled": true,
        "id": "0c4e939acacf4376bdcd1129f1a054ad",
        "links": {
            "self": "http://localhost:5000/v3/projects/0c4e939acacf4376bdcd1129f1a054ad"
        },
        "name": "admin",
        "parent_id": null
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
	
	projectResource := &Project {
		Session: session,
	}
	project, err := projectResource.Get("xyz")
	
	if err != nil {
		t.Error(err)
	}
	
	if project.Id != "0c4e939acacf4376bdcd1129f1a054ad" {
		t.Errorf("Expected id == %q; got %q", "default", project.Id)
	}
}

func TestProjectUpdate(t *testing.T) {
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
		if r.URL.Path != "/projects/xyz" {
			t.Errorf("Expected 'url' == %q; got %q", "/projects/xyz", r.URL.Path)
		}
		
		
		body, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		tmp := `{
    "project": {
        "description": "My updated project",
        "domain_id": "default",
        "enabled": true,
        "name": "myUpdatedProject"
    }
}`
		tmpStruct := &types.ProjectResponse{}
		_ = json.Unmarshal([]byte(tmp), tmpStruct)
		tmpbytes, _ := json.Marshal(tmpStruct)

		if string(body) != string(tmpbytes) {
			t.Errorf("Expected auth request body == %q; got %q", string(tmpbytes), string(body))
		}
		
		w.Write([]byte(`{
    "project": {
        "is_domain": true,
        "description": "My updated project",
        "links": {
            "self": "http://localhost:5000/v3/projects/93ebbcc35335488b96ff9cd7d18cbb2e"
        },
        "extra": {
            "is_domain": true
        },
        "enabled": true,
        "id": "93ebbcc35335488b96ff9cd7d18cbb2e",
        "parent_id": null,
        "domain_id": "default",
        "name": "myUpdatedProject"
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
	
	projectResource := &Project {
		Session: session,
	}

	project, err := projectResource.Update("xyz", &types.Project{
        Description: "My updated project",
        DomainId: "default",
        Enabled: true,
        Name: "myUpdatedProject",
	})
	if err != nil {
		t.Error(err)
	}
	if project.Id != "93ebbcc35335488b96ff9cd7d18cbb2e" {
		t.Errorf("Expect project.Id == %q; got %q", "default", project.Id)
	}
	if ! project.IsDomain {
		t.Errorf("Expect is_domain true, got false")
	}
	if v := project.Extra["is_domain"].(bool); !v {
		t.Errorf("Expect extra[is_domain] true, got false")
	}
}
