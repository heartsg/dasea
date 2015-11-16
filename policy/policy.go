package policy


//Based on openstack oslo.policy, used to check authorization policy

type PolicyRules struct {
    rules map[string]PolicyCheck
}
func (r *PolicyRules) RuleCheck(rule string, target map[string]interface{}, creds map[string]interface{}, enforcer *PolicyEnforcer) bool {
    if ruleObject, ok := r.rules[rule]; ok {
        return ruleObject.Check(target, creds, enforcer)
    } else {
        return false
    }
}

type PolicyEnforcer struct {
    rules *PolicyRules	
}
func (e *PolicyEnforcer) RuleCheck(rule string, target map[string]interface{}, creds map[string]interface{}) bool {
    return e.rules.RuleCheck(rule, target, creds, e)
}
