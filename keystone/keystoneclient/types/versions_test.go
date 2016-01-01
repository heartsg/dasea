package types

import (
	"testing"
	"encoding/json"
    "time"
    "github.com/heartsg/dasea/keystone/keystoneclient/util"
    "github.com/heartsg/dasea/testutil"
)

var vs = `{
    "versions": {
        "values": [
            {
                "id": "v3.4",
                "links": [
                    {
                        "href": "http://localhost:5000/v3/",
                        "rel": "self"
                    }
                ],
                "media-types": [
                    {
                        "base": "application/json",
                        "type": "application/vnd.openstack.identity-v3+json"
                    }
                ],
                "status": "stable",
                "updated": "2015-03-30T00:00:00Z"
            },
            {
                "id": "v2.0",
                "links": [
                    {
                        "href": "http://localhost:5000/v2.0/",
                        "rel": "self"
                    },
                    {
                        "href": "http://docs.openstack.org/",
                        "rel": "describedby",
                        "type": "text/html"
                    }
                ],
                "media-types": [
                    {
                        "base": "application/json",
                        "type": "application/vnd.openstack.identity-v2.0+json"
                    }
                ],
                "status": "stable",
                "updated": "2014-04-17T00:00:00Z"
            }
        ]
    }
}`

var structVs = &Versions {
    Versions: &VersionValues {
        Values:[]*VersionValue {
            &VersionValue {
                Id: "v3.4",
                Links: []Link {
                    Link {
                        "href": "http://localhost:5000/v3/",
                        "rel": "self",
                    },
                },
                MediaTypes: []MediaType{
                    MediaType {
                        "base": "application/json",
                        "type": "application/vnd.openstack.identity-v3+json",
                    },
                },
                Status: "stable",
                Updated: &util.Iso8601DateTime{time.Date(2015, time.March, 30, 0, 0, 0, 0, time.UTC)},
            },
            &VersionValue {
                Id: "v2.0",
                Links: []Link{
                    Link {
                        "href": "http://localhost:5000/v2.0/",
                        "rel": "self",
                    },
                    Link {
                        "href": "http://docs.openstack.org/",
                        "rel": "describedby",
                        "type": "text/html",                        
                    },
                },
                MediaTypes: []MediaType{
                    MediaType {
                        "base": "application/json",
                        "type": "application/vnd.openstack.identity-v2.0+json",
                    },
                },
                Status: "stable",
                Updated: &util.Iso8601DateTime{time.Date(2014, time.April, 17, 0, 0, 0, 0, time.UTC)},                
            },
        },
    },
}

func TestVersions(t *testing.T) {
	versions := &Versions{}
	err := json.Unmarshal([]byte(vs), versions)
	if err  != nil {
		t.Error("Failed to unmarshal version: %s", err)
	}
	
	testutil.Equals(t, versions, structVs)
}



var v = `{
    "version": {
        "id": "v3.4",
        "links": [
            {
                "href": "http://localhost:5000/v3/",
                "rel": "self"
            }
        ],
        "media-types": [
            {
                "base": "application/json",
                "type": "application/vnd.openstack.identity-v3+json"
            }
        ],
        "status": "stable",
        "updated": "2015-03-30T00:00:00Z"
    }
}`
var structV = &Version {
    &VersionValue {
        Id: "v3.4",
        Links: []Link {
            Link {
                "href": "http://localhost:5000/v3/",
                "rel": "self",
            },
        },
        MediaTypes: []MediaType{
            MediaType {
                "base": "application/json",
                "type": "application/vnd.openstack.identity-v3+json",
            },
        },
        Status: "stable",
        Updated: &util.Iso8601DateTime{time.Date(2015, time.March, 30, 0, 0, 0, 0, time.UTC)},
    },
}

func TestVersion(t *testing.T) {
	version := &Version{}
	err := json.Unmarshal([]byte(v), version)
	if err  != nil {
		t.Error("Failed to unmarshal version: %s", err)
	}
	
	testutil.Equals(t, version, structV)
}