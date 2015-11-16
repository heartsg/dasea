package policy

import (
	"fmt"
	"testing"
)

func TestLoadPolicyRules(t *testing.T) {
	//test load from map
	test1 := map[string]interface{} {
		"storage_put": "{{.is_admin}}:True or role:admin",
		"storage_get": "role:operator",
	}
	r1 := NewFromMap(test1)
	if storage_put, ok := r1.rules["storage_put"]; !ok {
		t.Errorf("Load policy rules error")
	} else {
		if cast, ok := storage_put.(*OrCheck); !ok {
			t.Errorf("Load policy rules error")
		} else {
			if fmt.Sprintf("%s", cast) != "({{.is_admin}}:True or role:admin)" {
				t.Errorf("Load policy rules error")
			}
		}
	}
	
	//test load from json
	test2 := `{"storage_put": "{{.is_admin}}:True or role:admin", "storage_get": "role:operator"}`
	r2 := NewFromJson([]byte(test2))
	if storage_put, ok := r2.rules["storage_put"]; !ok {
		t.Errorf("Load policy rules error")
	} else {
		if cast, ok := storage_put.(*OrCheck); !ok {
			t.Errorf("Load policy rules error")
		} else {
			if fmt.Sprintf("%s", cast) != "({{.is_admin}}:True or role:admin)" {
				t.Errorf("Load policy rules error")
			}
		}
	}	
}

func TestPolicyEnforcerLoad(t *testing.T) {
	policyOpts := &PolicyOpts{
		PolicyFile:"testdata/policy.json",
		PolicyDirs:[]string{
			"testdata/policy.d1",
			"testdata/policy.d2",
		},
		PolicyDefaultRule:"default",
	}
	
	//start loading
	e := &PolicyEnforcer{}
	e.LoadRules(policyOpts)
	
	//check
	if fmt.Sprintf("%s", e.r.rules["storage.put"]) != "role:admin4" {
		t.Errorf("policy enforcer load failed, storage.put is not role:admin4, got %s", e.r.rules["storage.put"])
	}
	if fmt.Sprintf("%s", e.r.rules["storage.get"]) != "role:admin3" {
		t.Errorf("policy enforcer load failed, storage.get is not role:admin3, got %s", e.r.rules["storage.get"])
	}
	if fmt.Sprintf("%s", e.r.rules["compute"]) != "role:admin2" {
		t.Errorf("policy enforcer load failed, compute is not role:admin2, got %s", e.r.rules["compute"])
	}
	if fmt.Sprintf("%s", e.r.rules["testouter"]) != "role:testouter" {
		t.Errorf("policy enforcer load failed, testouter is not role:testouter, got %s", e.r.rules["testouter"])
	}
	if fmt.Sprintf("%s", e.r.rules["testinner1"]) != "role:testinner1" {
		t.Errorf("policy enforcer load failed, testinner1 is not role:testinner1, got %s", e.r.rules["testinner1"])
	}
	if fmt.Sprintf("%s", e.r.rules["testinner2"]) != "role:testinner2" {
		t.Errorf("policy enforcer load failed, testinner2 is not role:testinner2, got %s", e.r.rules["testinner2"])
	}
	if fmt.Sprintf("%s", e.r.rules["testinner3"]) != "role:testinner3" {
		t.Errorf("policy enforcer load failed, testinner3 is not role:testinner3, got %s", e.r.rules["testinner3"])
	}
	if fmt.Sprintf("%s", e.r.rules["testinner4"]) != "role:testinner4" {
		t.Errorf("policy enforcer load failed, testinner4 is not role:testinner4, got %s", e.r.rules["testinner4"])
	}
}

func TestPolicyEnforcerEnforce(t *testing.T) {
	policyOpts := &PolicyOpts{
		PolicyFile:"testdata/policy.json",
		PolicyDirs:[]string{
			"testdata/policy.d1",
			"testdata/policy.d2",
		},
		PolicyDefaultRule:"default",
	}
	
	//start loading
	e := &PolicyEnforcer{}
	e.LoadRules(policyOpts)

	if e.Enforce(&TrueCheck{}, nil, nil) != true {
		t.Errorf("policy enforcer enforce error")
	}
	
	creds := map[string]interface{} {
		"roles":[]string{"admin3", "testinner2", "testouter"},
		"is_enabled":true,
	}
	target := map[string]interface{} {
		"is_disabled":false,
	}
	if e.Enforce("compute", nil, creds) != false {
		t.Errorf("policy enforcer enforce error")
	}
	if e.Enforce("storage.put", nil, creds) != false {
		t.Errorf("policy enforcer enforce error")
	}
	if e.Enforce("storage.get", nil, creds) != true {
		t.Errorf("policy enforcer enforce error")
	}
	if e.Enforce("testouter", nil, creds) != true {
		t.Errorf("policy enforcer enforce error")
	}
	if e.Enforce("testinner1", nil, creds) != false {
		t.Errorf("policy enforcer enforce error")
	}
	if e.Enforce("testinner2", nil, creds) != true {
		t.Errorf("policy enforcer enforce error")
	}
	if e.Enforce("testinner3", nil, creds) != false {
		t.Errorf("policy enforcer enforce error")
	}
	if e.Enforce("testinner4", nil, creds) != false {
		t.Errorf("policy enforcer enforce error")
	}
	if e.Enforce("testenable", nil, creds) != true {
		t.Errorf("policy enforcer enforce error")
	}
	if e.Enforce("testdisable", target, nil) != true {
		t.Errorf("policy enforcer enforce error")
	}
}