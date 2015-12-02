package client

//
// Implements similar functions as AuthPlugin in python-keystoneclient
//
// In python's keystoneclient implementation, GetAccess
// will make use of already created session for http request
// calls (to reuse the client and sessions).
//
// However, in our implementation, we currently create a new session
// or requests for http request calls, because the original design seems over
// complicated. We currently keep our design and only change if
// find that the system performance drops.
//
//

import (
)

type Auth interface {
	GetAccess(session *Session) (*AccessInfo, error)
}