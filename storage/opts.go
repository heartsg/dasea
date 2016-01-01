package storage

import (
	"github.com/heartsg/dasea/keystone/keystonemiddleware"
)
type DBOpts struct {
    Type string
    Hosts []string
    User string
    Password string
}
type Opts struct {
    // We currently uses sql-like relational database for meta data
	MetaDB DBOpts
    // We currently choose to use cassandra for real-time data
    DataDB DBOpts
	KeystoneMiddleware keystonemiddleware.Opts
}

