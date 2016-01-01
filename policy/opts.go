package policy

//Based on openstack oslo.policy, used to check authorization policy

//Policy options, use config to retrieve the options from file/flag/env etc.

type Opts struct {
	File       string `default:"policy.json"`
	DefaultRule string `default:"default"`
	Dirs		[]string `default:"policy.d"`
}
