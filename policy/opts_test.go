package policy

import (
	"testing"
	
	"github.com/heartsg/dasea/config"
)


func TestPolicyOpts(t *testing.T) {
	loader := config.New()
	policyOpts := &PolicyOpts{}
	loader.Load(policyOpts)
	if policyOpts.PolicyFile != "policy.json" {
		t.Error("PolicyOpts policy_file error")
	}
	if policyOpts.PolicyDefaultRule != "default" {
		t.Error("PolicyOpts policy_default_rule error")
	}
	if policyOpts.PolicyDirs[0] != "policy.d" {
		t.Error("PolicyOpts policy_dirs error")
	}
}
