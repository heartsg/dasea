package policy

//Based on openstack oslo.policy, used to check authorization policy

//Policy options, use config to retrieve the options from file/flag/env etc.

type PolicyOpts struct {
	PolicyFile       string `default:"policy.json"`
	PolicyDefaultRule string `default:"default"`
	PolicyDirs		[]string `default:"policy.d"`
}
