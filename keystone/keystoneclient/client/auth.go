package client

import (
)

//
// Session is used to communicate with keystone identity servers
// to retrieve access tokens.
//

type Auth interface {
	GetAccess(session *Session) (*AccessInfo, error)
}