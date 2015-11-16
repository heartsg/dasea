package policy

import (
	"fmt"
    "strings"
    "bytes"
    "strconv"
    "text/template"
)


type PolicyCheck interface {
	//The interface is same as oslo.policy's Check.__call__
	//target: used specifically for GenericCheck's template
	//creds: credential map that contains authorizaiton information
	//enforcer: used specifically for RuleCheck
	Check(target map[string]interface{}, creds map[string]interface{}, enforcer *PolicyEnforcer ) bool
}

type FalseCheck struct {}
func (c* FalseCheck) Check(target map[string]interface{}, creds map[string]interface{}, enforcer *PolicyEnforcer ) bool {
	//ignore all inputs
	return false
}
func (c* FalseCheck) String() string {
	return "!"
}


type TrueCheck struct {}
func (c* TrueCheck) Check(target map[string]interface{}, creds map[string]interface{}, enforcer *PolicyEnforcer ) bool {
	//ignore all inputs
	return true
}
func (c* TrueCheck) String() string {
	return "@"
}

type NotCheck struct {
	rule PolicyCheck
}
func (c *NotCheck) Check(target map[string]interface{}, creds map[string]interface{}, enforcer *PolicyEnforcer ) bool {
	return !c.rule.Check(target, creds, enforcer)
}
func (c* NotCheck) String() string {
	return fmt.Sprintf("not %s", c.rule)
}

type AndCheck struct {
	rules []PolicyCheck
}
func (c *AndCheck) Check(target map[string]interface{}, creds map[string]interface{}, enforcer *PolicyEnforcer ) bool {
	for _, rule := range c.rules {
		if !rule.Check(target, creds, enforcer) {
			return false
		}
	}
	return true
}
func (c* AndCheck) String() string {
	s := "("
	for i, r := range c.rules {
		s += fmt.Sprintf("%s", r)
		if i < len(c.rules) - 1 {
			s += " and "
		} 
	}
	s += ")"
	return s
}
func (c *AndCheck) AddCheck(rule PolicyCheck) {
	c.rules = append(c.rules, rule)
}

type OrCheck struct {
	rules []PolicyCheck
}
func (c *OrCheck) Check(target map[string]interface{}, creds map[string]interface{}, enforcer *PolicyEnforcer ) bool {
	for _, rule := range c.rules {
		if rule.Check(target, creds, enforcer) {
			return true
		}
	}
	return false
}
func (c* OrCheck) String() string {
	s := "("
	for i, r := range c.rules {
		s += fmt.Sprintf("%s", r)
		if i < len(c.rules) - 1 {
			s += " or "
		} 
	}
	s += ")"
	return s
}
func (c *OrCheck) AddCheck(rule PolicyCheck) {
	c.rules = append(c.rules, rule)
}


type RuleCheck struct {
	kind string
    match string
}
func (c *RuleCheck) Check(target map[string]interface{}, creds map[string]interface{}, enforcer *PolicyEnforcer ) bool {
    if enforcer != nil {
        return enforcer.RuleCheck(c.match, target, creds)
    }
	return true
}
func (c* RuleCheck) String() string {
	return fmt.Sprintf("%s:%s", c.kind, c.match)
}

type RoleCheck struct {
	kind string
    match string
}
func (c *RoleCheck) Check(target map[string]interface{}, creds map[string]interface{}, enforcer *PolicyEnforcer ) bool {
    if rolesInterface, ok := creds["roles"]; ok {
		if 	roles, ok := rolesInterface.([]string); ok {
			match := strings.ToLower(c.match)
			for _, role := range roles {
				if match == strings.ToLower(role) {
					return true
				}
			}
		}
	}
	
	return false
}
func (c* RoleCheck) String() string {
	return fmt.Sprintf("%s:%s", c.kind, c.match)
}

// Matches look like:
//    - {{.tenant}}:{{.tenant_id}}
//    - role:compute:admin
//    - True:{{.user_enabled}}
//    - 'Member':{{.role_name}}
type GenericCheck struct {
	kind string
    match string
}

//Returns literal and string representation if c.kind is a literal, only supports float64, int64 and string literals
func  literalEval(s string) (interface{}, string, bool) {
    if len(s) >= 2 && s[0] == '\'' && s[len(s)-1] == '\'' {
        return s[1:len(s)-1], s[1:len(s)-1], true
    } else if v1, err := strconv.ParseBool(s); err == nil {
        return v1, fmt.Sprintf("%t", v1), true
    } else if v2, err := strconv.ParseInt(s, 10, 64); err == nil {
        return v2, fmt.Sprintf("%d", v2), true
    } else if v3, err := strconv.ParseFloat(s, 64); err == nil {
        return v3, fmt.Sprintf("%g", v3), true
    }
    
    return nil, "", false
}
func (c *GenericCheck) Check(target map[string]interface{}, creds map[string]interface{}, enforcer *PolicyEnforcer ) bool {
	//firstly, check if match is literal
	_, match, ok := literalEval(c.match)
	if !ok {
		//{{.key}} in match is mapped to target["key"], convert to none template
		tmpl, err := template.New("gc").Parse(c.match)
		if err != nil {
			return false
		}
		var matchbuf bytes.Buffer
		err = tmpl.Execute(&matchbuf, target)
		if err != nil {
			return false
		}
		
		match = matchbuf.String()
	}

    
    //check if c.kind is literal
   	_, kind, ok := literalEval(c.kind); 
	
	if !ok {
		//{{.key}} in kind is maped to creds["key"], convert to none template
		tmpl, err := template.New("gc").Parse(c.kind)
		if err != nil {
			return false
		}
		var kindbuf bytes.Buffer
		err = tmpl.Execute(&kindbuf, creds)
		if err != nil {
			return false
		}
		
		kind = kindbuf.String()
	}

    return match == kind
}
func (c* GenericCheck) String() string {
	return fmt.Sprintf("%s:%s", c.kind, c.match)
}

//Parse a single check without compositions
//Generally in forms of "something:somethingelse"
func ParseCheck(rule string) PolicyCheck {
	if rule == "!" {
		return &FalseCheck{}
	} else if rule == "@" {
		return &TrueCheck{}
	}
	
	result := strings.SplitN(rule, ":", 2)
	if len(result) == 1 {
		//Do a warning log later, something wrong, 
		//I don't understand the format
		return &FalseCheck{}
	}
	kind := result[0]
	match := result[1]
	
	//currently we only support role and rule checks, all others
	// will be directed to generic check
	switch kind {
	case "role":
		return &RoleCheck{kind:kind, match:match}
	case "rule":
		return &RuleCheck{kind:kind, match:match}
	default:
		return &GenericCheck{kind:kind, match:match}
	}
}
