package policy

import (
    "os"
    "io/ioutil"
    "path/filepath"
    "encoding/json"
)

// 1. Overview
// Based on openstack oslo.policy, used to check authorization policy
// Read docs on oslo.policy for more details.
//
// 2. Policy files
//
// Policy files are json format files that stores key:value pairs.
// 
// Key is used to identify the service/API/or other entities that this policy
// applies. Example: "identity:get_user" represents the keystone identity service
// get_user API.
//
// Value represents the rule or check (they are roughly the samething in our docs) for
// the service/api identified by Key. Example will be: "role:admin or is_admin:true".
//
// So the overall key:value pair example will be "identity:get_user":"role:admin or {{.is_admin}}:true"
//
// 3. Rules/Checks
// 
// There are several types of rules
//   - logical: and/or/not
//   - role: the creds["roles"] caontains a list of roles that current session possesses, the 
//           role check passes (returns true) only if the role listed behind is in creds["roles"].
//           E.g., "role:admin" passes only creds["roles"] contains an element called "admin".
//   - rule: based on other rules in the policy. E.g.,
//           { 
//              "admin_required":"role:admin or is_admin:true",
//              "identity:get_user":"rule:admin_required or role:operator"
//           }
//    - generic rules, e.g., "{{.is_admin}}:true", the check passes only if creds["is_admin"] is true.
//      Another example: "true:{{.is_admin}}", the check passes only if target["is_admin"] is true.
//      "{{.is_admin}}:{{.is_admin}}", the check passes only if creds["is_admin"] == target["is_admin"]
// 
// 4. Usage
//
// Firstly, create a PolicyEnforcer object given PolicyOpts, the opts should contain the policy file to load
// etc. information.
//
//   enforcer := LoadRules(opts)
//
// Then we can enforce a policy by,
//
//   enforcer.Enforce(key, target, creds)
//
// where key is the Key such as "identity:get_user"
//
// If we already have a Check object, we can also pass check as key,
//
//    enforcer.Enforce(check, target, creds)
//
// Enforce will return true if the policy allows, or it will return false.
//

type PolicyRules struct {
    rules map[string]PolicyCheck
    defaultRule string
}
func (r *PolicyRules) RuleCheck(rule string, target map[string]interface{}, creds map[string]interface{}, enforcer *PolicyEnforcer) bool {
    if ruleObject, ok := r.rules[rule]; ok {
        return ruleObject.Check(target, creds, enforcer)
    } else if r.defaultRule != "" {
        if ruleObject, ok := r.rules[r.defaultRule]; ok {
            return ruleObject.Check(target, creds, enforcer)
        }
    }
    return false
}
func NewFromJson(data []byte) *PolicyRules {
    r := &PolicyRules { rules: make(map[string]PolicyCheck) }
    r.UpdateFromJson(data)
    return r
}
func NewFromMap(m map[string]interface{}) *PolicyRules {
    r := &PolicyRules { rules: make(map[string]PolicyCheck) }
    r.UpdateFromMap(m)
    return r
}
func (r *PolicyRules) UpdateFromJson(data []byte) {
    p := &policyParser{}
    
    var g interface{}
    if err := json.Unmarshal(data, &g); err == nil {
        if m, ok := g.(map[string]interface{}); ok {
            for k, v := range m {
                rule, err := p.parseRule(v)
                if err == nil {
                    r.rules[k] = rule
                }
            }
        }
    }
}
func (r *PolicyRules) UpdateFromMap(m map[string]interface{}) {
    p := &policyParser{}
    for k, v := range m {
        rule, err := p.parseRule(v)
        if err == nil {
            r.rules[k] = rule
        }
    }
}

type PolicyEnforcer struct {
    r *PolicyRules	
}
func (e *PolicyEnforcer) RuleCheck(rule string, target map[string]interface{}, creds map[string]interface{}) bool {
    return e.r.RuleCheck(rule, target, creds, e)
}
func (e *PolicyEnforcer) Enforce(rule interface{},  target map[string]interface{}, creds map[string]interface{}) bool {
    switch r := rule.(type) {
    case string:
        return e.RuleCheck(r, target, creds)
    case PolicyCheck:
        return r.Check(target, creds, e)
    default:
        return false
    }
}
func (e *PolicyEnforcer) LoadRules(o *PolicyOpts) {
    //load from json file if exists
    if o.File != "" {
        e.loadJsonFile(o.File)
    }
    
    //load from all files in json dir if exists
    if o.Dirs != nil {
        for _, dir := range o.Dirs {
            filepath.Walk(dir, e.visit)
        }
    }
    
    //set default rule
    if e.r == nil {
        e.r = &PolicyRules{ rules: make(map[string]PolicyCheck) }
    }
    e.r.defaultRule = o.DefaultRule
}
func (e *PolicyEnforcer) loadJsonFile(path string) {
    data, err := ioutil.ReadFile(path)
    if err != nil {
        //log error?
    } else {
        if e.r == nil || e.r.rules == nil {
            e.r = NewFromJson(data)
        } else {
            e.r.UpdateFromJson(data)
        }
    }
}
func (e *PolicyEnforcer) visit(path string, f os.FileInfo, err error) error {
    if err != nil || f.IsDir() {
        return nil
    } else {
        e.loadJsonFile(path)
        return nil
    }
}