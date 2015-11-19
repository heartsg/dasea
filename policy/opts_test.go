package policy

import (
	"testing"
	
	"github.com/heartsg/dasea/config"
)


func TestPolicyOpts(t *testing.T) {
	loader := config.New()
	policyOpts := &PolicyOpts{}
	loader.Load(policyOpts)
	if policyOpts.File != "policy.json" {
		t.Error("PolicyOpts file error")
	}
	if policyOpts.DefaultRule != "default" {
		t.Error("PolicyOpts default_rule error")
	}
	if policyOpts.Dirs[0] != "policy.d" {
		t.Error("PolicyOpts dirs error")
	}
}
