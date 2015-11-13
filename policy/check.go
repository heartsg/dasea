package policy

import (
	"fmt"
)

//Based on openstack oslo.policy, used to check authorization policy

type PolicyCheck interface {
	//The interface is same as oslo.policy's Check.__call__
	//target: used specifically for GenericCheck's template
	//creds: credential map to authorize
	//enforcer: used specifically for RuleCheck
	func Check(target *interface{}, creds map[string]interface{}, enforcer *PolicyEnforcer ) bool
}

type FalseCheck struct {}
func (c* FalseCheck) Check(target *interface{}, creds map[string]interface{}, enforcer *PolicyEnforcer ) bool {
	//ignore all inputs
	return false
}
func (c* FalseCheck) String() string {
	return "!"
}


type TrueCheck struct {}
func (c* TrueCheck) Check(target *interface{}, creds map[string]interface{}, enforcer *PolicyEnforcer ) bool {
	//ignore all inputs
	return true
}
func (c* TrueCheck) String() string {
	return "@"
}

type NotCheck struct {
	rule Check
}
func (c *NotCheck) Check(target *interface{}, creds map[string]interface{}, enforcer *PolicyEnforcer ) bool {
	return !rule.Check(target, creds, enforcer)
}
func (c* NotCheck) String() string {
	return fmt.Sprintf("not %s", c.rule)
}

type AndCheck struct {
	rules []Check
}
func (c *AndCheck) Check(target *interface{}, creds map[string]interface{}, enforcer *PolicyEnforcer ) bool {
	for _, rule := range rules {
		if !rule.Check(target, creds, enforcer) {
			return false
		}
	}
	return true
}
func (c* AndCheck) String() string {
	return fmt.Sprintf("and %s", c.rules)
}
func (c *AndCheck) AddCheck(rule Check) {
	append(c.rules, rule)
}