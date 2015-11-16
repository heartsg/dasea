package policy

import (
	"testing"
	
	"fmt"
)


func TestFalseCheck(t *testing.T) {
	c := &FalseCheck{}
	if fmt.Sprintf("%s", c) != "!" {
		t.Errorf("FalseCheck is %s, want !", c)
	}
	
	if c.Check(nil, nil, nil) != false {
		t.Errorf("FalseCheck check fails")
	}
}



func TestTrueCheck(t *testing.T) {
	c := &TrueCheck{}
	if fmt.Sprintf("%s", c) != "@" {
		t.Errorf("TrueCheck is %s, want @", c)
	}
	
	if c.Check(nil, nil, nil) != true {
		t.Errorf("TrueCheck check fails")
	}
}


func TestNotCheck(t *testing.T) {
	c1 := &NotCheck{ rule:&FalseCheck{} }
	if fmt.Sprintf("%s", c1) != "not !" {
		t.Errorf("NotCheck is %s, want not !", c1)
	}
	
	if c1.Check(nil, nil, nil) != true {
		t.Errorf("NotCheck check fails")
	}

	c2 := &NotCheck{ rule:&TrueCheck{} }
	if fmt.Sprintf("%s", c2) != "not @" {
		t.Errorf("NotCheck is %s, want not !", c2)
	}
	
	if c2.Check(nil, nil, nil) != false {
		t.Errorf("NotCheck check fails")
	}
}


func TestAndCheck(t *testing.T) {
	c1 := &AndCheck{ rules:[]PolicyCheck{ &FalseCheck{}, &TrueCheck{} } }
	if fmt.Sprintf("%s", c1) != "(! and @)" {
		t.Errorf("AndCheck is %s, want (! and @)", c1)
	}
	
	if c1.Check(nil, nil, nil) != false {
		t.Errorf("AndCheck check fails")
	}

	c2 := &AndCheck{ rules:[]PolicyCheck{ &NotCheck{rule:&FalseCheck{}}, &TrueCheck{} } }
	if fmt.Sprintf("%s", c2) != "(not ! and @)" {
		t.Errorf("AndCheck is %s, want (not ! and @)", c2)
	}
	
	if c2.Check(nil, nil, nil) != true {
		t.Errorf("AndCheck check fails")
	}
}


func TestOrCheck(t *testing.T) {
	c1 := &OrCheck{ rules:[]PolicyCheck{ &FalseCheck{}, &TrueCheck{} } }
	if fmt.Sprintf("%s", c1) != "(! or @)" {
		t.Errorf("OrCheck is %s, want (! or @)", c1)
	}
	
	if c1.Check(nil, nil, nil) != true {
		t.Errorf("OrCheck check fails")
	}

	c2 := &OrCheck{ rules:[]PolicyCheck{ &NotCheck{rule:&TrueCheck{}}, &FalseCheck{} } }
	if fmt.Sprintf("%s", c2) != "(not @ or !)" {
		t.Errorf("OrCheck is %s, want (not @ or !)", c2)
	}
	
	if c2.Check(nil, nil, nil) != false {
		t.Errorf("OrCheck check fails")
	}
}


func TestRuleCheck(t *testing.T) {
	e := &PolicyEnforcer {
		r: &PolicyRules {
			rules:map[string]PolicyCheck {
				"true": &TrueCheck{},
				"false": &FalseCheck{},
			},
		},
	}
	c1 := &RuleCheck{kind:"truerule", match:"true"}
	if fmt.Sprintf("%s", c1) != "truerule:true" {
		t.Errorf("RuleCheck is %s, want truerule:true", c1)
	}
	
	if c1.Check(nil, nil, e) != true {
		t.Errorf("RuleCheck check fails")
	}
	
	c2 := &RuleCheck{kind:"falserule", match:"false"}
	if fmt.Sprintf("%s", c2) != "falserule:false" {
		t.Errorf("RuleCheck is %s, want falserule:false", c1)
	}
	
	if c2.Check(nil, nil, e) != false {
		t.Errorf("RuleCheck check fails")
	}
	
	c3 := &RuleCheck{kind:"nullrule", match:"notexist"}
	
	if c3.Check(nil, nil, e) != false {
		t.Errorf("RuleCheck check fails")
	}
}

func TestRoleCheck(t *testing.T) {
	creds := map[string]interface{} {
		"roles":[]string{"admin","operator"},
	}
	
	c1 := &RoleCheck{kind:"role", match:"admin"}
	if fmt.Sprintf("%s", c1) != "role:admin" {
		t.Errorf("RoleCheck is %s, want role:admin", c1)
	}
	
	if c1.Check(nil, creds, nil) != true {
		t.Errorf("RoleCheck check fails")
	}
	
	c2 := &RoleCheck{kind:"role", match:"operator"}
	if fmt.Sprintf("%s", c2) != "role:operator" {
		t.Errorf("RoleCheck is %s, want role:operator", c2)
	}
	
	if c2.Check(nil, creds, nil) != true {
		t.Errorf("RoleCheck check fails")
	}
	
	c3 := &RoleCheck{kind:"role", match:"guest"}
	if fmt.Sprintf("%s", c3) != "role:guest" {
		t.Errorf("RoleCheck is %s, want role:guest", c3)
	}
	
	if c3.Check(nil, creds, nil) != false {
		t.Errorf("RoleCheck check fails")
	}
}



func TestGenericCheck(t *testing.T) {
	//check boolean constant
	c1 := &GenericCheck{kind:"True", match:"{{.is_enable}}"}
	if fmt.Sprintf("%s", c1) != "True:{{.is_enable}}" {
		t.Errorf("GenericCheck is %s, want True:{{.is_enable}}", c1)
	}
	
	if c1.Check(map[string]interface{} { "is_enable":true }, nil, nil) != true {
		t.Errorf("Generic check fails")
	}
	
	//check int constant
	c2 := &GenericCheck{kind:"100", match:"{{.count}}"}

	if c2.Check(map[string]interface{} { "is_enable":true, "count":100 }, nil, nil) != true {
		t.Errorf("Generic check fails")
	}
	
	if c2.Check(map[string]interface{} { "is_enable":true, "count":101 }, nil, nil) != false {
		t.Errorf("Generic check fails")
	}
	
	if c2.Check(map[string]interface{} { "is_enable":true }, nil, nil) != false {
		t.Errorf("Generic check fails")
	}
	
	//check float constant
	c3 := &GenericCheck{kind:"3.14", match:"{{.pi}}"}

	if c3.Check(map[string]interface{} { "pi":3.14 }, nil, nil) != true {
		t.Errorf("Generic check fails")
	}
	
	//check string literal
	c4 := &GenericCheck{kind:"'member'", match:"{{.whoami}}"}
	if fmt.Sprintf("%s", c4) != "'member':{{.whoami}}" {
		t.Errorf("GenericCheck is %s, want 'member':{{.whoami}}", c4)
	}
	
	if c4.Check(map[string]interface{} { "whoami":"member" }, nil, nil) != true {
		t.Errorf("Generic check fails")
	}
	
	//check generic, which is, check the creds[c.kind] and target[c.match] equals or not 
	c5 := &GenericCheck{kind:"{{.tenant}}", match:"{{.tenant_id}}"}
	
	if c5.Check(map[string]interface{} { "tenant_id":12345 }, map[string]interface{} { "tenant":12345 }, nil) != true {
		t.Errorf("Generic check fails")
	}
	
	//check literal for match
	c6 := &GenericCheck{match:"True", kind:"{{.is_enable}}"}
	
	if c6.Check(nil, map[string]interface{} { "is_enable":true }, nil) != true {
		t.Errorf("Generic check fails")
	}
	
	c7 := &GenericCheck{match:"'member'", kind:"{{.whoami}}"}
	
	if c7.Check(nil, map[string]interface{} { "whoami":"member" }, nil) != true {
		t.Errorf("Generic check fails")
	}
	
	//check both are literal
	c8 := &GenericCheck{match:"true", kind:"T"}
	
	if c8.Check(nil, nil, nil) != true {
		t.Errorf("Generic check fails")
	}
}

func TestParseCheck(t *testing.T) {
	//Test parse false check
	c1 := ParseCheck("!")
	if _, ok := c1.(*FalseCheck); !ok {
		t.Errorf("! parsed to %s, want *FalseCheck", c1)
	}
	
	//Test parse true check
	c2 := ParseCheck("@")
	if _, ok := c2.(*TrueCheck); !ok {
		t.Errorf("@ parsed to %s, want *TrueCheck", c2)
	}
	
	//Test parse unknown check
	c3 := ParseCheck("unknown??@!@")
	if _, ok := c3.(*FalseCheck); !ok {
		t.Errorf("unknown??@!@ parsed to %s, want *FalseCheck", c3)
	}
	
	//Test rule check
	c4 := ParseCheck("rule:myrule")
	if c, ok := c4.(*RuleCheck); !ok {
		t.Errorf("rule:myrule parsed to %s, want *RuleCheck", c4)
	} else if c.kind != "rule" || c.match != "myrule" {
		t.Errorf("kind:match should be rule:myrule, got %s", c)
	}
	
	//Test role check
	c5 := ParseCheck("role:admin")
	if c, ok := c5.(*RoleCheck); !ok {
		t.Errorf("role:admin parsed to %s, want *RoleCheck", c5)
	} else if c.kind != "role" || c.match != "admin" {
		t.Errorf("kind:match should be role:admin, got %s", c)
	}
	
	//Test generic check
	c6 := ParseCheck("True:T")
	if c, ok := c6.(*GenericCheck); !ok {
		t.Errorf("True:T parsed to %s, want *GenericCheck", c6)
	} else if c.kind != "True" || c.match != "T" {
		t.Errorf("kind:match should be True:T, got %s", c)
	}
	
	c7 := ParseCheck("'member':{{.account_type}}")
	if c, ok := c7.(*GenericCheck); !ok {
		t.Errorf("'member':{{.account_type}} parsed to %s, want *GenericCheck", c7)
	} else if c.kind != "'member'" || c.match != "{{.account_type}}" {
		t.Errorf("kind:match should be 'member':{{.account_type}}, got %s", c)
	}
	
	c8 := ParseCheck("{{.tenant}}:{{.tenant_id}}")
	if c, ok := c8.(*GenericCheck); !ok {
		t.Errorf("{{.tenant}}:{{.tenant_id}} parsed to %s, want *GenericCheck", c8)
	} else if c.kind != "{{.tenant}}" || c.match != "{{.tenant_id}}" {
		t.Errorf("kind:match should be {{.tenant}}:{{.tenant_id}}, got %s", c)
	}
}