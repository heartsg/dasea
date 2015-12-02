package types


import (
	"github.com/heartsg/dasea/keystone/keystoneclient/util"
)


// Defines all requests & response structures for keystone Identity API V3.
// Follows the docs in http://developer.openstack.org/api-ref-identity-v3.html.
// Also to facilitate Json marshalling and unmarshalling.


//
// Get /
// Function: List Versions
// Request: None
// Response: Defined as struct Versions
//
type VersionValue struct {
	Id string `json:"id"`
	Links []Link `json:"links"`
	MediaTypes []MediaType `json:"media-types"`
	Status string `json:"status"`
	Updated *util.Iso8601DateTime `json:"updated"`
}
type VersionValues struct {
	Values []*VersionValue `json:"values"`
}
type Versions struct {
	Versions *VersionValues `json:"versions"`
}


//
// Get /v3
// Function: version 3
// Request: None (Note that request might include in Accept header, but we don't
//          currently support it).
// Response: Defined as struct Version
//           (Note: response may also be Json-home documents, we don't current
//            support it)
type Version struct {
	Version *VersionValue `json:"version"`
}