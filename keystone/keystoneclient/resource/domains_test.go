package resource

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"io/ioutil"
	"encoding/json"
	
	"github.com/heartsg/dasea/requests"
	"github.com/heartsg/dasea/keystone/keystoneclient/client"
	"github.com/heartsg/dasea/keystone/keystoneclient/types"
)

func TestDomainCreate(t *testing.T) {
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
    "domain": {
        "description": "Domain description",
        "enabled": true,
        "name": "myDomain"
    }
}`
		tmpStruct := &types.DomainResponse{}
		_ = json.Unmarshal([]byte(tmp), tmpStruct)
		tmpbytes, _ := json.Marshal(tmpStruct)
		if string(body) != string(tmpbytes) {
			t.Errorf("Expected auth request body == %q; got %q", string(tmpbytes), string(body))
		}
		
		w.Write([]byte(`{
    "domain": {
        "description": "Domain description",
        "enabled": true,
        "id": "161718",
        "links": {
            "self": "http://identity:35357/v3/domains/161718"
        },
        "name": "myDomain"
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
	
	domainResource := &Domain {
		Session: session,
	}
	domain, err := domainResource.Create(&types.Domain{
		Description: "Domain description",
		Enabled: true,
		Name: "myDomain",
	})
	if err != nil {
		t.Error(err)
	}
	if domain.Id != "161718" {
		t.Errorf("Expect domain.Id == %q; got %q", "161718", domain.Id)
	}
}

func TestDomainList(t *testing.T) {
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
		if q.Get("enabled") != "true" {
			t.Errorf("Expected 'enabled' == %q; got %q", "true", q.Get("enabled"))
		}
		
		w.Write([]byte(`{
    "domains": [
        {
            "description": "Used for swift functional testing",
            "enabled": true,
            "id": "5a75994a383c449184053ff7270c4e91",
            "links": {
                "self": "http://localhost:5000/v3/domains/5a75994a383c449184053ff7270c4e91"
            },
            "name": "swift_test"
        },
        {
            "description": "Owns users and tenants (i.e. projects) available on Identity API v2.",
            "enabled": true,
            "id": "default",
            "links": {
                "self": "http://localhost:5000/v3/domains/default"
            },
            "name": "Default"
        }
    ],
    "links": {
        "next": null,
        "previous": null,
        "self": "http://localhost:5000/v3/domains"
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
	
	domainResource := &Domain {
		Session: session,
	}
	domains, err := domainResource.List("", "true")
	if err != nil {
		t.Error(err)
	}
	if domains.Links["self"] != "http://localhost:5000/v3/domains" {
		t.Errorf("Expect domains.Links[self] == %q; got %q", "http://localhost:5000/v3/domains", domains.Links["self"])
	}
	if domains.Domains[0].Id != "5a75994a383c449184053ff7270c4e91" {
		t.Errorf("Expect domains.Domains[0].Id == %q; got %q", "5a75994a383c449184053ff7270c4e91", domains.Domains[0].Id)
	}
	if domains.Domains[1].Links["self"] != "http://localhost:5000/v3/domains/default" {
		t.Errorf("Expect domains.Domains[1].Links[self] == %q; got %q", "http://localhost:5000/v3/domains/default", domains.Domains[1].Links["self"])
	}
}

func TestDomainDelete (t *testing.T) {
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
		if r.URL.Path != "/domains/xyz" {
			t.Errorf("Expected 'url' == %q; got %q", "/domains/xyz", r.URL.Path)
		}
	}))
	defer ts.Close()
	
	
	auth := &client.AuthAccess {
		Access: &client.AccessInfo {
			Token: "12345",
		},
	}
	
	session := client.NewSession(ts.URL, auth)
	
	domainResource := &Domain {
		Session: session,
	}
	err := domainResource.Delete("xyz")
	
	if err != nil {
		t.Error(err)
	}
}

func TestDomainGet (t *testing.T) {
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
		if r.URL.Path != "/domains/xyz" {
			t.Errorf("Expected 'url' == %q; got %q", "/domains/xyz", r.URL.Path)
		}
		
		w.Write([]byte(`{
    "domain": {
        "description": "Owns users and tenants (i.e. projects) available on Identity API v2.",
        "enabled": true,
        "id": "default",
        "links": {
            "self": "http://localhost:5000/v3/domains/default"
        },
        "name": "Default"
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
	
	domainResource := &Domain {
		Session: session,
	}
	domain, err := domainResource.Get("xyz")
	
	if err != nil {
		t.Error(err)
	}
	
	if domain.Id != "default" {
		t.Errorf("Expected id == %q; got %q", "default", domain.Id)
	}
}

func TestDomainUpdate(t *testing.T) {
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
		if r.URL.Path != "/domains/xyz" {
			t.Errorf("Expected 'url' == %q; got %q", "/domains/xyz", r.URL.Path)
		}
		
		
		body, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		tmp := `{
    "domain": {
        "description": "Owns users and projects on Identity API v2."
    }
}`
		tmpStruct := &types.DomainResponse{}
		_ = json.Unmarshal([]byte(tmp), tmpStruct)
		tmpbytes, _ := json.Marshal(tmpStruct)

		if string(body) != string(tmpbytes) {
			t.Errorf("Expected auth request body == %q; got %q", string(tmpbytes), string(body))
		}
		
		w.Write([]byte(`{
    "domain": {
        "links": {
            "self": "http://localhost:5000/v3/domains/default"
        },
        "enabled": true,
        "description": "Owns users and projects on Identity API v2.",
        "name": "Default",
        "id": "default"
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
	
	domainResource := &Domain {
		Session: session,
	}

	domain, err := domainResource.Update("xyz", &types.Domain{
		Description: "Owns users and projects on Identity API v2.",
	})
	if err != nil {
		t.Error(err)
	}
	if domain.Id != "default" {
		t.Errorf("Expect domain.Id == %q; got %q", "default", domain.Id)
	}
}
