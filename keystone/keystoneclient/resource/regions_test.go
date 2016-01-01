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

func TestRegionCreate(t *testing.T) {
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
    "region": {
        "description": "My subregion",
        "id": "RegionOneSubRegion",
        "parent_region_id": "RegionOne"
    }
}`
		tmpStruct := &types.RegionResponse{}
		_ = json.Unmarshal([]byte(tmp), tmpStruct)
		tmpbytes, _ := json.Marshal(tmpStruct)
		if string(body) != string(tmpbytes) {
			t.Errorf("Expected auth request body == %q; got %q", string(tmpbytes), string(body))
		}
		
		w.Write([]byte(`{
    "region": {
        "parent_region_id": "RegionOne",
        "id": "RegionOneSubRegion",
        "links": {
            "self": "http://localhost:5000/v3/regions/RegionOneSubRegion"
        },
        "description": "My subregion"
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
	
	regionResource := &Region {
		Session: session,
	}
	region, err := regionResource.Create(&types.Region{
		Description: "My subregion",
		Id: "RegionOneSubRegion",
		ParentRegionId: "RegionOne",
	})
	if err != nil {
		t.Error(err)
	}
	if region.Id != "RegionOneSubRegion" {
		t.Errorf("Expect region.Id == %q; got %q", "RegionOneSubRegion", region.Id)
	}
}

func TestRegionList(t *testing.T) {
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
		if q.Get("parent_region_id") != "test" {
			t.Errorf("Expected 'enabled' == %q; got %q", "true", q.Get("parent_region_id"))
		}
		
		w.Write([]byte(`{
    "links": {
        "next": null,
        "previous": null,
        "self": "http://localhost:5000/v3/regions"
    },
    "regions": [
        {
            "description": "",
            "id": "RegionOne",
            "links": {
                "self": "http://localhost:5000/v3/regions/RegionOne"
            },
            "parent_region_id": null
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
	
	regionResource := &Region {
		Session: session,
	}
	regions, err := regionResource.List("test")
	if err != nil {
		t.Error(err)
	}
	if regions.Links["self"] != "http://localhost:5000/v3/regions" {
		t.Errorf("Expect regions.Links[self] == %q; got %q", "http://localhost:5000/v3/regions", regions.Links["self"])
	}
	if regions.Regions[0].Id != "RegionOne" {
		t.Errorf("Expect regions.Regions[0].Id == %q; got %q", "RegionOne", regions.Regions[0].Id)
	}
}

func TestRegionDelete (t *testing.T) {
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
		if r.URL.Path != "/regions/xyz" {
			t.Errorf("Expected 'url' == %q; got %q", "/regions/xyz", r.URL.Path)
		}
	}))
	defer ts.Close()
	
	
	auth := &client.AuthAccess {
		Access: &client.AccessInfo {
			Token: "12345",
		},
	}
	
	session := client.NewSession(&keystoneclient.Opts{ AuthUrl: ts.URL }, auth)
	
	regionResource := &Region {
		Session: session,
	}
	err := regionResource.Delete("xyz")
	
	if err != nil {
		t.Error(err)
	}
}

func TestRegionGet (t *testing.T) {
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
		if r.URL.Path != "/regions/xyz" {
			t.Errorf("Expected 'url' == %q; got %q", "/regions/xyz", r.URL.Path)
		}
		
		w.Write([]byte(`{
    "region": {
        "description": "My subregion 3",
        "id": "RegionThree",
        "links": {
            "self": "http://localhost:5000/v3/regions/RegionThree"
        },
        "parent_region_id": "RegionOne"
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
	
	regionResource := &Region {
		Session: session,
	}
	region, err := regionResource.Get("xyz")
	
	if err != nil {
		t.Error(err)
	}
	
	if region.Id != "RegionThree" {
		t.Errorf("Expected id == %q; got %q", "RegionThree", region.Id)
	}
}

func TestRegionUpdate(t *testing.T) {
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
		if r.URL.Path != "/regions/xyz" {
			t.Errorf("Expected 'url' == %q; got %q", "/regions/xyz", r.URL.Path)
		}
		
		
		body, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		tmp := `{
    "region": {
        "description": "My subregion 3"
    }
}`
		tmpStruct := &types.RegionResponse{}
		_ = json.Unmarshal([]byte(tmp), tmpStruct)
		tmpbytes, _ := json.Marshal(tmpStruct)

		if string(body) != string(tmpbytes) {
			t.Errorf("Expected auth request body == %q; got %q", string(tmpbytes), string(body))
		}
		
		w.Write([]byte(`{
    "region": {
        "parent_region_id": "RegionOne",
        "id": "RegionThree",
        "links": {
            "self": "http://localhost:5000/v3/regions/RegionThree"
        },
        "description": "My subregion 3"
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
	
	regionResource := &Region {
		Session: session,
	}

	region, err := regionResource.Update("xyz", &types.Region{
		Description: "My subregion 3",
	})
	if err != nil {
		t.Error(err)
	}
	if region.Id != "RegionThree" {
		t.Errorf("Expect region.Id == %q; got %q", "RegionThree", region.Id)
	}
}
