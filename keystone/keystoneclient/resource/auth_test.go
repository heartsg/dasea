package resource

import (
	"net/http"
	"net/http/httptest"
	"testing"
	
	"github.com/heartsg/dasea/requests"
	"github.com/heartsg/dasea/keystone/keystoneclient/client"
)

func TestGetAccess(t *testing.T) {
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
		if r.Header.Get("X-Subject-Token") != "67890" {
			t.Errorf("Expected 'X-Subject-Token' == %q; got %q", "67890", r.Header.Get("X-Auth-Token"))			
		}
		
		w.Header().Set("X-Auth-Token", "12345")
		w.Header().Set("X-Subject-Token", "67890")
		w.Write([]byte(`{
    "token": {
        "methods": [
            "token"
        ],
        "expires_at": "2015-11-05T22:00:11.000000Z",
        "extras": {},
        "user": {
            "domain": {
                "id": "default",
                "name": "Default"
            },
            "id": "10a2e6e717a245d9acad3e5f97aeca3d",
            "name": "admin"
        },
        "audit_ids": [
            "mAjXQhiYRyKwkB4qygdLVg"
        ],
        "issued_at": "2015-11-05T21:00:33.819948Z"
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
	
	authResource := &Auth {
		Session: session,
	}
	access, err := authResource.GetAccessFromToken("67890")
	if err != nil {
		t.Error(err)
	}
	if access.Token != "67890" {
		t.Errorf("Expect token = %q; got %q", "67890", access.Token)
	}
}