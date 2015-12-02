package keystoneclient

//
// AuthUrl : Keystone server address for version 3 identity APIs
//           We only support version 3. Therefore the identity endpoint
//           should be AuthUrl/v3/some_resources.
//
type KeystoneclientOpts struct {
	AuthUrl string 		`default:"127.0.0.1:35357"`
}