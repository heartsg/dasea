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

func TestCredentialCreate(t *testing.T) {
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

		if string(body) != `{"credential":{"blob":"{\"access\":\"181920\",\"secret\":\"secretKey\"}","user_id":"bb5476fd12884539b41d5a88f838d773","project_id":"731fc6f265cd486d900f16e84c5cb594","type":"ec2"}}` {
			t.Errorf("Expected auth request body == %q; got %q", 
				`{"credential":{"blob":"{\"access\":\"181920\",\"secret\":\"secretKey\"}","user_id":"bb5476fd12884539b41d5a88f838d773","project_id":"731fc6f265cd486d900f16e84c5cb594","type":"ec2"}}`,
				string(body))
		}
		
		w.Write([]byte(`{
    "credential": {
        "user_id": "bb5476fd12884539b41d5a88f838d773",
        "links": {
            "self": "http://localhost:5000/v3/credentials/3d3367228f9c7665266604462ec60029bcd83ad89614021a80b2eb879c572510"
        },
        "blob": "{\"access\":\"181920\",\"secret\":\"secretKey\"}",
        "project_id": "731fc6f265cd486d900f16e84c5cb594",
        "type": "ec2",
        "id": "3d3367228f9c7665266604462ec60029bcd83ad89614021a80b2eb879c572510"
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
	
	credentialResource := &Credential {
		Session: session,
	}
	blob, _ := json.Marshal(map[string]string {
			"access":"181920",
			"secret":"secretKey",
		})
	credential, err := credentialResource.Create(&types.Credential{
		Blob: string(blob),
		ProjectId: "731fc6f265cd486d900f16e84c5cb594",
		Type: "ec2",
		UserId: "bb5476fd12884539b41d5a88f838d773",
	})
	if err != nil {
		t.Error(err)
	}
	if credential.Id != "3d3367228f9c7665266604462ec60029bcd83ad89614021a80b2eb879c572510" {
		t.Errorf("Expect credential.Id == %q; got %q", "3d3367228f9c7665266604462ec60029bcd83ad89614021a80b2eb879c572510", credential.Id)
	}
}

func TestCredentialList(t *testing.T) {
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
		if q.Get("user_id") != "xyz" {
			t.Errorf("Expected 'user_id' == %q; got %q", "xyz", q.Get("user_id"))
		}
		
		w.Write([]byte(`{
    "credentials": [
        {
            "user_id": "bb5476fd12884539b41d5a88f838d773",
            "links": {
                "self": "http://localhost:5000/v3/credentials/207e9b76935efc03804d3dd6ab52d22e9b22a0711e4ada4ff8b76165a07311d7"
            },
            "blob": "{\"access\": \"a42a27755ce6442596b049bd7dd8a563\", \"secret\": \"71faf1d40bb24c82b479b1c6fbbd9f0c\", \"trust_id\": null}",
            "project_id": "6e01855f345f4c59812999b5e459137d",
            "type": "ec2",
            "id": "207e9b76935efc03804d3dd6ab52d22e9b22a0711e4ada4ff8b76165a07311d7"
        },
        {
            "user_id": "6f556708d04b4ea6bc72d7df2296b71a",
            "links": {
                "self": "http://localhost:5000/v3/credentials/2441494e52ab6d594a34d74586075cb299489bdd1e9389e3ab06467a4f460609"
            },
            "blob": "{\"access\": \"7da79ff0aa364e1396f067e352b9b79a\", \"secret\": \"7a18d68ba8834b799d396f3ff6f1e98c\", \"trust_id\": null}",
            "project_id": "1a1d14690f3c4ec5bf5f321c5fde3c16",
            "type": "ec2",
            "id": "2441494e52ab6d594a34d74586075cb299489bdd1e9389e3ab06467a4f460609"
        }
	],
	"links": {
        "self": "http://localhost:5000/v3/credentials",
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
	
	credentialResource := &Credential {
		Session: session,
	}
	credentials, err := credentialResource.List("xyz")
	if err != nil {
		t.Error(err)
	}
	if credentials.Links["self"] != "http://localhost:5000/v3/credentials" {
		t.Errorf("Expect credentials.Links[self] == %q; got %q", "http://localhost:5000/v3/credentials", credentials.Links["self"])
	}
	if credentials.Credentials[0].Id != "207e9b76935efc03804d3dd6ab52d22e9b22a0711e4ada4ff8b76165a07311d7" {
		t.Errorf("Expect credentials.Credentials[0].Id == %q; got %q", "207e9b76935efc03804d3dd6ab52d22e9b22a0711e4ada4ff8b76165a07311d7", credentials.Credentials[0].Id)
	}
	if credentials.Credentials[1].Links["self"] != "http://localhost:5000/v3/credentials/2441494e52ab6d594a34d74586075cb299489bdd1e9389e3ab06467a4f460609" {
		t.Errorf("Expect credentials.Credentials[1].Links[self] == %q; got %q", "http://localhost:5000/v3/credentials/2441494e52ab6d594a34d74586075cb299489bdd1e9389e3ab06467a4f460609", credentials.Credentials[1].Links["self"])
	}
}

func TestCredentialDelete (t *testing.T) {
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
		if r.URL.Path != "/credentials/xyz" {
			t.Errorf("Expected 'url' == %q; got %q", "/credentials/xyz", r.URL.Path)
		}
	}))
	defer ts.Close()
	
	
	auth := &client.AuthAccess {
		Access: &client.AccessInfo {
			Token: "12345",
		},
	}
	
	session := client.NewSession(&keystoneclient.Opts{ AuthUrl: ts.URL }, auth)
	
	credentialResource := &Credential {
		Session: session,
	}
	err := credentialResource.Delete("xyz")
	
	if err != nil {
		t.Error(err)
	}
}

func TestCredentialGet (t *testing.T) {
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
		if r.URL.Path != "/credentials/xyz" {
			t.Errorf("Expected 'url' == %q; got %q", "/credentials/xyz", r.URL.Path)
		}
		
		w.Write([]byte(`{
    "credential": {
        "user_id": "bb5476fd12884539b41d5a88f838d773",
        "links": {
            "self": "http://localhost:5000/v3/credentials/207e9b76935efc03804d3dd6ab52d22e9b22a0711e4ada4ff8b76165a07311d7"
        },
        "blob": "{\"access\": \"a42a27755ce6442596b049bd7dd8a563\", \"secret\": \"71faf1d40bb24c82b479b1c6fbbd9f0c\", \"trust_id\": null}",
        "project_id": "6e01855f345f4c59812999b5e459137d",
        "type": "ec2",
        "id": "207e9b76935efc03804d3dd6ab52d22e9b22a0711e4ada4ff8b76165a07311d7"
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
	
	credentialResource := &Credential {
		Session: session,
	}
	credential, err := credentialResource.Get("xyz")
	
	if err != nil {
		t.Error(err)
	}
	
	if credential.Id != "207e9b76935efc03804d3dd6ab52d22e9b22a0711e4ada4ff8b76165a07311d7" {
		t.Errorf("Expected id == %q; got %q", "207e9b76935efc03804d3dd6ab52d22e9b22a0711e4ada4ff8b76165a07311d7", credential.Id)
	}
}

func TestCredentialUpdate(t *testing.T) {
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
		if r.URL.Path != "/credentials/xyz" {
			t.Errorf("Expected 'url' == %q; got %q", "/credentials/xyz", r.URL.Path)
		}
		
		
		body, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		if string(body) != `{"credential":{"blob":"{\"access\":\"181920\",\"secret\":\"secretKey\"}","user_id":"bb5476fd12884539b41d5a88f838d773","project_id":"731fc6f265cd486d900f16e84c5cb594","type":"ec2"}}` {
			t.Errorf("Expected auth request body == %q; got %q", 
				`{"credential":{"blob":"{\"access\":\"181920\",\"secret\":\"secretKey\"}","user_id":"bb5476fd12884539b41d5a88f838d773","project_id":"731fc6f265cd486d900f16e84c5cb594","type":"ec2"}}`,
				string(body))
		}
		
		w.Write([]byte(`{
    "credential": {
        "user_id": "bb5476fd12884539b41d5a88f838d773",
        "links": {
            "self": "http://localhost:5000/v3/credentials/3d3367228f9c7665266604462ec60029bcd83ad89614021a80b2eb879c572510"
        },
        "blob": "{\"access\":\"181920\",\"secret\":\"secretKey\"}",
        "project_id": "731fc6f265cd486d900f16e84c5cb594",
        "type": "ec2",
        "id": "3d3367228f9c7665266604462ec60029bcd83ad89614021a80b2eb879c572510"
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
	
	credentialResource := &Credential {
		Session: session,
	}
	blob, _ := json.Marshal(map[string]string {
			"access":"181920",
			"secret":"secretKey",
		})
	credential, err := credentialResource.Update("xyz", &types.Credential{
		Blob: string(blob),
		ProjectId: "731fc6f265cd486d900f16e84c5cb594",
		Type: "ec2",
		UserId: "bb5476fd12884539b41d5a88f838d773",
	})
	if err != nil {
		t.Error(err)
	}
	if credential.Id != "3d3367228f9c7665266604462ec60029bcd83ad89614021a80b2eb879c572510" {
		t.Errorf("Expect credential.Id == %q; got %q", "3d3367228f9c7665266604462ec60029bcd83ad89614021a80b2eb879c572510", credential.Id)
	}
}
