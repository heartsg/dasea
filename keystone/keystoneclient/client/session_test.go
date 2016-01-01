package client

import (
	"net/http/httptest"
	"net/http"
	"testing"
	"io/ioutil"
	
	"github.com/heartsg/dasea/requests"
	"github.com/heartsg/dasea/keystone/keystoneclient"
	"github.com/heartsg/dasea/keystone/keystoneclient/types"
)

func TestSession(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check method is GET before going to check other features
		if r.Method != requests.GET {
			t.Errorf("Expected method %q; got %q", requests.GET, r.Method)
		}
		if r.Header == nil {
			t.Errorf("Expected non-nil request Header")
		}
		if r.Header.Get("User-Agent") != "Dasea Keystone client" {
			t.Errorf("Expected 'User-Agent' == %q; got %q", "Dasea Keystone client", r.Header.Get("User-Agent"))
		}
		if r.Header.Get("Forwarded") != "for=127.0.0.1;by=Dasea Keystone client" {
			t.Errorf("Expected 'Forwarded' == %q; got %q", "for=127.0.0.1;by=Dasea Keystone client", r.Header.Get("Forwarded"))
		}
	}))	

	defer ts.Close()
	
	session := NewSession(&keystoneclient.Opts{ AuthUrl: ts.URL }, nil)
	_, _, err := session.OriginalIp("127.0.0.1").Request("/", requests.GET, nil, nil, nil, false)
	if err != nil {
		t.Error(err)
	}

}

func TestSessionAuth(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check method is GET before going to check other features
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
		if r.Header.Get("Forwarded") != "for=127.0.0.1;by=Dasea Keystone client" {
			t.Errorf("Expected 'Forwarded' == %q; got %q", "for=127.0.0.1;by=Dasea Keystone client", r.Header.Get("Forwarded"))
		}
	}))	

	defer ts.Close()
	
	auth := &AuthAccess {
		Access: &AccessInfo {
			Token: "12345",
		},
	}
	
	session := NewSession(&keystoneclient.Opts{ AuthUrl: ts.URL }, auth)
	_, _, err := session.OriginalIp("127.0.0.1").Request("/", requests.GET, nil, nil, nil, true)
	if err != nil {
		t.Error(err)
	}

}

func TestSessionAuthPassword(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header == nil {
			t.Errorf("Expected non-nil request Header")
		}
		if r.Header.Get("User-Agent") != "Dasea Keystone client" {
			t.Errorf("Expected 'User-Agent' == %q; got %q", "Dasea Keystone client", r.Header.Get("User-Agent"))
		}
		
		if r.Header.Get("Forwarded") != "for=127.0.0.1;by=Dasea Keystone client" {
			t.Errorf("Expected 'Forwarded' == %q; got %q", "for=127.0.0.1;by=Dasea Keystone client", r.Header.Get("Forwarded"))
		}
		switch r.URL.Path {
		case "/":
			if r.Method != requests.GET {
				t.Errorf("Expected method %q; got %q", requests.GET, r.Method)
			}
			if r.Header.Get("X-Auth-Token") != "12345" {
				t.Errorf("Expected 'X-Auth-Token' == %q; got %q", "12345", r.Header.Get("X-Auth-Token"))
			}
		case "/auth/tokens":
			if r.Method != requests.POST {
				t.Errorf("Expected method %q; got %q", requests.POST, r.Method)
			}
			body, _ := ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			
			if string(body) != `{"auth":{"identity":{"methods":["password"],"password":{"user":{"name":"admin","domain":{"id":"default"},"password":"devstacker"}}}}}` {
				t.Errorf("Expected auth request body == %q; got %q", 
					`{"auth":{"identity":{"methods":["password"],"password":{"user":{"name":"admin","domain":{"id":"default"},"password":"devstacker"}}}}}`,
					string(body))
			}
			w.Header().Set("X-Subject-Token", "12345")
			w.Write([]byte(`{
    "token": {
        "methods": [
            "password"
        ],
        "expires_at": "2015-11-06T15:32:17.893769Z",
        "extras": {},
        "user": {
            "domain": {
                "id": "default",
                "name": "Default"
            },
            "id": "423f19a4ac1e4f48bbb4180756e6eb6c",
            "name": "admin"
        },
        "audit_ids": [
            "ZzZwkUflQfygX7pdYDBCQQ"
        ],
        "issued_at": "2015-11-06T14:32:17.893797Z"
    }
}`))
		}
	}))	

	defer ts.Close()
	
	auth := &AuthPassword {
		Params: &types.AuthRequestParams {
			Username: "admin",
			DomainId: "default",
			Password: "devstacker",
		},
	}
	
	session := NewSession(&keystoneclient.Opts{ AuthUrl: ts.URL }, auth)
	_, _, err := session.OriginalIp("127.0.0.1").Request("/", requests.GET, nil, nil, nil, true)
	if err != nil {
		t.Error(err)
	}
}

func TestSessionAuthToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header == nil {
			t.Errorf("Expected non-nil request Header")
		}
		if r.Header.Get("User-Agent") != "Dasea Keystone client" {
			t.Errorf("Expected 'User-Agent' == %q; got %q", "Dasea Keystone client", r.Header.Get("User-Agent"))
		}
		
		if r.Header.Get("Forwarded") != "for=127.0.0.1;by=Dasea Keystone client" {
			t.Errorf("Expected 'Forwarded' == %q; got %q", "for=127.0.0.1;by=Dasea Keystone client", r.Header.Get("Forwarded"))
		}
		switch r.URL.Path {
		case "/":
			if r.Method != requests.GET {
				t.Errorf("Expected method %q; got %q", requests.GET, r.Method)
			}
			if r.Header.Get("X-Auth-Token") != "12345" {
				t.Errorf("Expected 'X-Auth-Token' == %q; got %q", "12345", r.Header.Get("X-Auth-Token"))
			}
		case "/auth/tokens":
			if r.Method != requests.POST {
				t.Errorf("Expected method %q; got %q", requests.POST, r.Method)
			}
			body, _ := ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			
			if string(body) != `{"auth":{"identity":{"methods":["token"],"token":{"id":"'$OS_TOKEN'"}}}}` {
				t.Errorf("Expected auth request body == %q; got %q", 
					`{"auth":{"identity":{"methods":["token"],"token":{"id":"'$OS_TOKEN'"}}}}`,
					string(body))
			}
			w.Header().Set("X-Subject-Token", "12345")
			w.Write([]byte(`{
    "token": {
        "methods": [
            "token"
        ],
        "expires_at": "2015-11-06T15:32:17.893769Z",
        "extras": {},
        "user": {
            "domain": {
                "id": "default",
                "name": "Default"
            },
            "id": "423f19a4ac1e4f48bbb4180756e6eb6c",
            "name": "admin"
        },
        "audit_ids": [
            "ZzZwkUflQfygX7pdYDBCQQ"
        ],
        "issued_at": "2015-11-06T14:32:17.893797Z"
    }
}`))
		}
	}))	

	defer ts.Close()
	
	auth := &AuthToken {
		Params: &types.AuthRequestParams {
			Token: "'$OS_TOKEN'",
		},
	}
	
	session := NewSession(&keystoneclient.Opts{ AuthUrl: ts.URL }, auth)
	_, _, err := session.OriginalIp("127.0.0.1").Request("/", requests.GET, nil, nil, nil, true)
	if err != nil {
		t.Error(err)
	}	
}