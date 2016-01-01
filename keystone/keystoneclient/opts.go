package keystoneclient

//
// AuthUrl : Keystone server address for version 3 identity APIs
//           We only support version 3. Therefore the identity endpoint
//           should be AuthUrl/v3/some_resources.
//
type Opts struct {
	// Keystone identity server address
	AuthUrl string 		`default:"127.0.0.1:35357"`
	
	// For identity server connection timeout in seconds (keystone server), 0 means no timeout
	HttpConnectionTimeout int

	// Currently we don't use the following
	HttpRequestMaxRetries int `default:"3"`
	CertFile string
	KeyFile string
	CAFile string
}