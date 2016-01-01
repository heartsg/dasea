package keystonemiddleware

import (
    "time"
	"github.com/heartsg/dasea/keystone/keystoneclient"
)

type Opts struct {
	AuthMethod string `default:"password"`
	Username string // generally it is a service account username/password pair
    UserId string
	Password string
	UserDomainId string // In service project domain
    UserDomainName string
	ProjectName string  // In project service
    ProjectId string
	ProjectDomainId string
    ProjectDomainName string
	Token string // Normally won't use it
    
    // Currently we don't use RegionName
    // RegionName string
	
	//Client.AuthUrl contains the URL for authentication
	Client keystoneclient.Opts
    
    // If delay decision, even failed to auth, still return non-nil ctx
	// so that next middleware can still be executed. The final decision
	// can be made by the end handler or middlewares after me
	DelayAuthDecision bool 	`default:"false"`
    
    // The following settings are regarded to token caching mechanisms.
    // Each service that uses AuthToken middleware will try to cache the user
    // tokens so that retrieving token information requently from keystone
    // identity service is not necessary.
    
    // Optionally specify a list of memcached server(s) to
    // use for caching. If left undefined, tokens will instead be
    // cached in-process.
	MemcacheServers []string
    // In order to prevent excessive effort spent validating
    // tokens, the middleware caches previously-seen tokens for a
    // configurable duration (in seconds).
    TokenCacheTime time.Duration `default:"300s"`
    // Determines the frequency at which the list of revoked
    // tokens is retrieved from the Identity service (in seconds). A
    // high number of revocation events combined with a low cache
    // duration may significantly reduce performance.
    RevocationCacheTime time.Duration `default:"10s"`
    
    // (Optional) If defined, indicate whether token data
    // should be authenticated or authenticated and encrypted.
    // Acceptable values are MAC or ENCRYPT.  If MAC, token data is
    // authenticated (with HMAC) in the cache. If ENCRYPT, token
    // data is encrypted and authenticated in the cache. If the
    // value is not one of these options or empty, auth_token will
    // raise an exception on initialization.
    // MemcacheSecurityStrategy string
    // (Optional, mandatory if memcache_security_strategy is
    // defined) This string is used for key derivation.
    // MemcacheSecretKey string
    // Hash algorithms to use for hashing PKI tokens. This may
    // be a single algorithm or multiple. The algorithms are those
    // supported by Python standard hashlib.new(). The hashes will
    // be tried in the order given, so put the preferred one first
    // for performance. The result of the first hash will be stored
    // in the cache. This will typically be set to multiple values
    // only while migrating from a less secure algorithm to a more
    // secure one. Once all the old tokens are expired this option
    // should be set to a single value for better performance.
    // HashAlgorithms []string `default:"md5"`
}