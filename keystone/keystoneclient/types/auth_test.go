package types

import (
	"testing"
	"encoding/json"
	"time"
	"github.com/heartsg/dasea/keystone/keystoneclient/util"
	"github.com/heartsg/dasea/testutil"
)


func TestAuthRequest(t *testing.T) {
	//request 1 : unscoped password authentication
	request1Raw := `{"auth":{"identity":{"methods":["password"],"password":{"user":{"id":"423f19a4ac1e4f48bbb4180756e6eb6c","password":"devstacker"}}}}}`
	request1Struct := &AuthRequest{
		Auth: &AuthRequestAuth {
			Identity: &AuthRequestIdentity {
				Methods: []string { "password" },
				Password: &AuthRequestPassword {
					User: &User {
						Id: "423f19a4ac1e4f48bbb4180756e6eb6c",
						Password: "devstacker",
					},
				},
			},
		},
	}
	
	request1Unmarshal := &AuthRequest{}
	err := json.Unmarshal([]byte(request1Raw), request1Unmarshal)
	
	testutil.IsNil(t, err)
	testutil.Equals(t, request1Struct, request1Unmarshal)
	
	request1Marshal, err := json.Marshal(request1Struct)
	testutil.IsNil(t, err)
	testutil.Equals(t, request1Raw, string(request1Marshal))
	
	//request 2 : another unscoped password authentication
	request2Raw := `{"auth":{"identity":{"methods":["password"],"password":{"user":{"name":"admin","domain":{"id":"default"},"password":"devstacker"}}}}}`
	request2Struct := &AuthRequest{
		Auth: &AuthRequestAuth {
			Identity: &AuthRequestIdentity {
				Methods: []string { "password" },
				Password: &AuthRequestPassword {
					User: &User {
						Name: "admin",
						Domain: &Domain {
							Id: "default",	
						},
						Password: "devstacker",
					},
				},
			},
		},
	}
	
	request2Unmarshal := &AuthRequest{}
	err = json.Unmarshal([]byte(request2Raw), request2Unmarshal)
	
	testutil.IsNil(t, err)
	testutil.Equals(t, request2Struct, request2Unmarshal)
	
	request2Marshal, err := json.Marshal(request2Struct)
	testutil.IsNil(t, err)
	testutil.Equals(t, request2Raw, string(request2Marshal))
	
	//request 3: explicit unscoped password authentication

	request3Raw := `{"auth":{"identity":{"methods":["password"],"password":{"user":{"id":"ee4dfb6e5540447cb3741905149d9b6e","password":"devstacker"}}},"scope":"unscoped"}}`
	request3Struct := &AuthRequest{
		Auth: &AuthRequestAuth {
			Identity: &AuthRequestIdentity {
				Methods: []string { "password" },
				Password: &AuthRequestPassword {
					User: &User {
						Id: "ee4dfb6e5540447cb3741905149d9b6e",
						Password: "devstacker",
					},
				},
			},
			Scope: "unscoped",
		},
	}
	
	request3Unmarshal := &AuthRequest{}
	err = json.Unmarshal([]byte(request3Raw), request3Unmarshal)
	
	testutil.IsNil(t, err)
	testutil.Equals(t, request3Struct, request3Unmarshal)
	
	request3Marshal, err := json.Marshal(request3Struct)
	testutil.IsNil(t, err)
	testutil.Equals(t, request3Raw, string(request3Marshal))
	
	//request 4: scoped password authentication

	request4Raw := `{"auth":{"identity":{"methods":["password"],"password":{"user":{"id":"ee4dfb6e5540447cb3741905149d9b6e","password":"devstacker"}}},"scope":{"project":{"id":"a6944d763bf64ee6a275f1263fae0352"}}}}`
	request4Struct1 := &AuthRequest{
		Auth: &AuthRequestAuth {
			Identity: &AuthRequestIdentity {
				Methods: []string { "password" },
				Password: &AuthRequestPassword {
					User: &User {
						Id: "ee4dfb6e5540447cb3741905149d9b6e",
						Password: "devstacker",
					},
				},
			},
			Scope: map[string]interface{} {
				"project": map[string]interface{} {
					"id": "a6944d763bf64ee6a275f1263fae0352",
				},
			},
		},
	}
	
	request4Struct2 := &AuthRequest{
		Auth: &AuthRequestAuth {
			Identity: &AuthRequestIdentity {
				Methods: []string { "password" },
				Password: &AuthRequestPassword {
					User: &User {
						Id: "ee4dfb6e5540447cb3741905149d9b6e",
						Password: "devstacker",
					},
				},
			},
			Scope: &Scope {
				Project: &Project {
					Id: "a6944d763bf64ee6a275f1263fae0352",
				},
			},
		},
	}
	
	request4Unmarshal := &AuthRequest{}
	err = json.Unmarshal([]byte(request4Raw), request4Unmarshal)
	
	testutil.IsNil(t, err)
	testutil.Equals(t, request4Struct1, request4Unmarshal)
	
	request4Marshal, err := json.Marshal(request4Struct1)
	testutil.IsNil(t, err)
	testutil.Equals(t, request4Raw, string(request4Marshal))
	
	request4Marshal, err = json.Marshal(request4Struct2)
	testutil.IsNil(t, err)
	testutil.Equals(t, request4Raw, string(request4Marshal))
	
	//request 5: unscoped token authorization

	request5Raw := `{"auth":{"identity":{"methods":["token"],"token":{"id":"'$OS_TOKEN'"}}}}`
	request5Struct := &AuthRequest{
		Auth: &AuthRequestAuth {
			Identity: &AuthRequestIdentity {
				Methods: []string { "token" },
				Token: &AuthRequestToken {
					Id: "'$OS_TOKEN'",
				},
			},
		},
	}
	
	request5Unmarshal := &AuthRequest{}
	err = json.Unmarshal([]byte(request5Raw), request5Unmarshal)
	
	testutil.IsNil(t, err)
	testutil.Equals(t, request5Struct, request5Unmarshal)
	
	request5Marshal, err := json.Marshal(request5Struct)
	testutil.IsNil(t, err)
	testutil.Equals(t, request5Raw, string(request5Marshal))
	
	//request 6: scoped token authorization
	request6Raw := `{"auth":{"identity":{"methods":["token"],"token":{"id":"'$OS_TOKEN'"}},"scope":{"project":{"id":"5b50efd009b540559104ee3c03bbb2b7"}}}}`
	request6Struct := &AuthRequest{
		Auth: &AuthRequestAuth {
			Identity: &AuthRequestIdentity {
				Methods: []string { "token" },
				Token: &AuthRequestToken {
					Id: "'$OS_TOKEN'",
				},
			},
			Scope: map[string]interface{} {
				"project": map[string]interface{} {
					"id": "5b50efd009b540559104ee3c03bbb2b7",
				},
			},
		},
	}
	
	request6Unmarshal := &AuthRequest{}
	err = json.Unmarshal([]byte(request6Raw), request6Unmarshal)
	
	testutil.IsNil(t, err)
	testutil.Equals(t, request6Struct, request6Unmarshal)
	
	request6Marshal, err := json.Marshal(request6Struct)
	testutil.IsNil(t, err)
	testutil.Equals(t, request6Raw, string(request6Marshal))
}


func TestAuthResponse(t *testing.T) {
	//response 1 : unscoped response
	response1Raw := `{"token":{"methods":["password"],"expires_at":"2015-11-06T15:32:17.893769Z","extras":{},"user":{"domain":{"id":"default","name":"Default"},"id":"423f19a4ac1e4f48bbb4180756e6eb6c","name":"admin"},"audit_ids":["ZzZwkUflQfygX7pdYDBCQQ"],"issued_at":"2015-11-06T14:32:17.893797Z"}}`
	response1Struct := &AuthResponse{
		Token: &AuthResponseToken {
			Methods: []string { "password" },
			ExpiresAt: &util.Iso8601DateTime{ time.Date(2015, time.November, 06, 15, 32, 17, 893769000, time.UTC) },
			Extras: make(map[string]interface{}),
			User: &User {
				Id: "423f19a4ac1e4f48bbb4180756e6eb6c",
				Name: "admin",
				Domain: &Domain {
					Id: "default",
					Name: "Default",
				},
			},
			AuditIds: []string { "ZzZwkUflQfygX7pdYDBCQQ" },
			IssuedAt: &util.Iso8601DateTime{ time.Date(2015, time.November, 06, 14, 32, 17, 893797000, time.UTC) },
		},
	}
	
	response1Unmarshal := &AuthResponse{}
	err := json.Unmarshal([]byte(response1Raw), response1Unmarshal)
	testutil.IsNil(t, err)
	testutil.Equals(t, response1Struct, response1Unmarshal)
	
	//response 2: password scoped response
	response2Raw := `{
		"token": {
			"methods": [
				"password"
			],
			"roles": [
				{
					"id": "51cc68287d524c759f47c811e6463340",
					"name": "admin"
				}
			],
			"expires_at": "2015-11-07T02:58:43.578887Z",
			"project": {
				"domain": {
					"id": "default",
					"name": "Default"
				},
				"id": "a6944d763bf64ee6a275f1263fae0352",
				"name": "admin"
			},
			"catalog": [
				{
					"endpoints": [
						{
							"region_id": "RegionOne",
							"url": "http://23.253.248.171:5000/v2.0",
							"region": "RegionOne",
							"interface": "public",
							"id": "068d1b359ee84b438266cb736d81de97"
						},
						{
							"region_id": "RegionOne",
							"url": "http://23.253.248.171:35357/v2.0",
							"region": "RegionOne",
							"interface": "admin",
							"id": "8bfc846841ab441ca38471be6d164ced"
						},
						{
							"region_id": "RegionOne",
							"url": "http://23.253.248.171:5000/v2.0",
							"region": "RegionOne",
							"interface": "internal",
							"id": "beb6d358c3654b4bada04d4663b640b9"
						}
					],
					"type": "identity",
					"id": "050726f278654128aba89757ae25950c",
					"name": "keystone"
				},
				{
					"endpoints": [
						{
							"region_id": "RegionOne",
							"url": "http://23.253.248.171:8774/v2/a6944d763bf64ee6a275f1263fae0352",
							"region": "RegionOne",
							"interface": "admin",
							"id": "ae36c0dbb0634e1dbf711f9fc2359975"
						},
						{
							"region_id": "RegionOne",
							"url": "http://23.253.248.171:8774/v2/a6944d763bf64ee6a275f1263fae0352",
							"region": "RegionOne",
							"interface": "internal",
							"id": "d286b51530144d90a4de52d214d3ad1e"
						},
						{
							"region_id": "RegionOne",
							"url": "http://23.253.248.171:8774/v2/a6944d763bf64ee6a275f1263fae0352",
							"region": "RegionOne",
							"interface": "public",
							"id": "d6e681dd4aab4ae5a0937ed60bb4ae33"
						}
					],
					"type": "compute_legacy",
					"id": "1c4bfbabe3b346b1bbe27a4b3258964f",
					"name": "nova_legacy"
				}
			],
			"extras": {},
			"user": {
				"domain": {
					"id": "default",
					"name": "Default"
				},
				"id": "ee4dfb6e5540447cb3741905149d9b6e",
				"name": "admin"
			},
			"audit_ids": [
				"3T2dc1CGQxyJsHdDu1xkcw"
			],
			"issued_at": "2015-11-07T01:58:43.578929Z"
		}
	}`
	response2Struct := &AuthResponse{
		Token: &AuthResponseToken {
			Methods: []string { "password" },
			Roles: []*Role {
				&Role {
					Id: "51cc68287d524c759f47c811e6463340",
					Name: "admin",
				},
			},
			ExpiresAt: &util.Iso8601DateTime{ time.Date(2015, time.November, 07, 02, 58, 43, 578887000, time.UTC) },
			Project: &Project {
				Domain: &Domain {
					Id: "default",
					Name: "Default",
				},
				Id: "a6944d763bf64ee6a275f1263fae0352",
				Name: "admin",
			},
			Catalog: []*AuthResponseCatalog {
				&AuthResponseCatalog {
					Endpoints: []*Endpoint {
						&Endpoint {
							RegionId: "RegionOne",
							Url: "http://23.253.248.171:5000/v2.0",
							Interface: "public",
							Region: "RegionOne",
							Id: "068d1b359ee84b438266cb736d81de97",
						},
						&Endpoint {
							RegionId: "RegionOne",
							Url: "http://23.253.248.171:35357/v2.0",
							Interface: "admin",
							Region: "RegionOne",
							Id: "8bfc846841ab441ca38471be6d164ced",
						},
						&Endpoint {
							RegionId: "RegionOne",
							Url: "http://23.253.248.171:5000/v2.0",
							Interface: "internal",
							Region: "RegionOne",
							Id: "beb6d358c3654b4bada04d4663b640b9",
						},
					},
					Type: "identity",
					Id: "050726f278654128aba89757ae25950c",
					Name: "keystone",
				},
				&AuthResponseCatalog {
					Endpoints: []*Endpoint {
						&Endpoint {
							RegionId: "RegionOne",
							Url: "http://23.253.248.171:8774/v2/a6944d763bf64ee6a275f1263fae0352",
							Interface: "admin",
							Region: "RegionOne",
							Id: "ae36c0dbb0634e1dbf711f9fc2359975",
						},
						&Endpoint {
							RegionId: "RegionOne",
							Url: "http://23.253.248.171:8774/v2/a6944d763bf64ee6a275f1263fae0352",
							Interface: "internal",
							Region: "RegionOne",
							Id: "d286b51530144d90a4de52d214d3ad1e",
						},
						&Endpoint {
							RegionId: "RegionOne",
							Url: "http://23.253.248.171:8774/v2/a6944d763bf64ee6a275f1263fae0352",
							Interface: "public",
							Region: "RegionOne",
							Id: "d6e681dd4aab4ae5a0937ed60bb4ae33",
						},
					},
					Type: "compute_legacy",
					Id: "1c4bfbabe3b346b1bbe27a4b3258964f",
					Name: "nova_legacy",
				},
			},
			Extras: make(map[string]interface{}),
			User: &User {
				Id: "ee4dfb6e5540447cb3741905149d9b6e",
				Name: "admin",
				Domain: &Domain {
					Id: "default",
					Name: "Default",
				},
			},
			AuditIds: []string { "3T2dc1CGQxyJsHdDu1xkcw" },
			IssuedAt: &util.Iso8601DateTime{ time.Date(2015, time.November, 07, 01, 58, 43, 578929000, time.UTC) },
		},
	}

	response2Unmarshal := &AuthResponse{}
	err = json.Unmarshal([]byte(response2Raw), response2Unmarshal)
	testutil.IsNil(t, err)
	testutil.Equals(t, response2Struct, response2Unmarshal)
}

func TestAuthRequestCreation(t *testing.T) {
	auth1, err := NewAuthRequestFromParams(&AuthRequestParams{
		UserId: "423f19a4ac1e4f48bbb4180756e6eb6c",
		Password: "devstacker",
	})
	testutil.IsNil(t, err)
	auth1Raw := `{"auth":{"identity":{"methods":["password"],"password":{"user":{"id":"423f19a4ac1e4f48bbb4180756e6eb6c","password":"devstacker"}}}}}`
	auth1Json, err := json.Marshal(auth1)
	testutil.Equals(t, auth1Raw, string(auth1Json))
	
	auth2, err := NewAuthRequestFromParams(&AuthRequestParams{
		Username: "admin",
		Password: "devstacker",
		DomainId: "default",
	})
	testutil.IsNil(t, err)
	auth2Raw := `{"auth":{"identity":{"methods":["password"],"password":{"user":{"name":"admin","domain":{"id":"default"},"password":"devstacker"}}}}}`
	auth2Json, err := json.Marshal(auth2)
	testutil.Equals(t, auth2Raw, string(auth2Json))
	
	auth3, err := NewAuthRequestFromParams(&AuthRequestParams{
		UserId: "ee4dfb6e5540447cb3741905149d9b6e",
		Password: "devstacker",
		ExplicitUnscope: true,
	})
	testutil.IsNil(t, err)
	auth3Raw := `{"auth":{"identity":{"methods":["password"],"password":{"user":{"id":"ee4dfb6e5540447cb3741905149d9b6e","password":"devstacker"}}},"scope":"unscoped"}}`
	auth3Json, err := json.Marshal(auth3)
	testutil.Equals(t, auth3Raw, string(auth3Json))
	
	
	auth4, err := NewAuthRequestFromParams(&AuthRequestParams{
		UserId: "ee4dfb6e5540447cb3741905149d9b6e",
		Password: "devstacker",
		ProjectId: "a6944d763bf64ee6a275f1263fae0352",
		Scope: true,
	})
	testutil.IsNil(t, err)
	auth4Raw := `{"auth":{"identity":{"methods":["password"],"password":{"user":{"id":"ee4dfb6e5540447cb3741905149d9b6e","password":"devstacker"}}},"scope":{"project":{"id":"a6944d763bf64ee6a275f1263fae0352"}}}}`
	auth4Json, err := json.Marshal(auth4)
	testutil.Equals(t, auth4Raw, string(auth4Json))
	
	auth5, err := NewAuthRequestFromParams(&AuthRequestParams{
		Token: "'$OS_TOKEN'",
	})
	testutil.IsNil(t, err)
	auth5Raw := `{"auth":{"identity":{"methods":["token"],"token":{"id":"'$OS_TOKEN'"}}}}`
	auth5Json, err := json.Marshal(auth5)
	testutil.Equals(t, auth5Raw, string(auth5Json))
	
	auth6, err := NewAuthRequestFromParams(&AuthRequestParams{
		Token: "'$OS_TOKEN'",
		ProjectId: "5b50efd009b540559104ee3c03bbb2b7",
		Scope: true,
	})
	testutil.IsNil(t, err)
	auth6Raw := `{"auth":{"identity":{"methods":["token"],"token":{"id":"'$OS_TOKEN'"}},"scope":{"project":{"id":"5b50efd009b540559104ee3c03bbb2b7"}}}}`
	auth6Json, err := json.Marshal(auth6)
	testutil.Equals(t, auth6Raw, string(auth6Json))

}