package policy

import (
	"fmt"
)

//Based on openstack oslo.policy, used to check authorization policy


//PolicyParser is used to parse rules into key, object map

type struct policyParser {
}

/*

func (p *policyParser) parseTextRule(rule string) error {
	return nil
}

func (p *policyParser) parseListRule(rule []string) error {
	return nil
}

func (p *policyParser) parseRule(rule interface{}) error {
	switch r := rule.type {
	case string:
		return p.parseTextRule(r)
	case []string:
		return p.parseListRule(r)
	}
	
	return fmt.Errorf("policy.policyParser: rule format not supported")
}*/