package policy

import (
    "os"
    "io/ioutil"
    "path/filepath"
    "encoding/json"
)

//Based on openstack oslo.policy, used to check authorization policy

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
    if o.PolicyFile != "" {
        e.loadJsonFile(o.PolicyFile)
    }
    
    //load from all files in json dir if exists
    if o.PolicyDirs != nil {
        for _, dir := range o.PolicyDirs {
            filepath.Walk(dir, e.visit)
        }
    }
    
    //set default rule
    if e.r == nil {
        e.r = &PolicyRules{ rules: make(map[string]PolicyCheck) }
    }
    e.r.defaultRule = o.PolicyDefaultRule
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