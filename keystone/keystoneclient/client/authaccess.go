package client

import (
	"errors"
)

type AuthAccess struct {
	Access *AccessInfo
}

func (a *AuthAccess) GetAccess(session *Session) (*AccessInfo, error) {
	if a.Access.WillExpireSoon() {
		return nil, errors.New("Token expired")
	}
	return a.Access, nil
}