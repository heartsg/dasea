package cache

import (
	"time"
	"encoding/json"
	"github.com/heartsg/dasea/keystone/keystoneclient/types"
)
type TokenCache struct {
	cache Cache
	defaultExpiration int
}

func NewTokenCache(servers[] string, expiration time.Duration) *TokenCache {
	var c Cache
	if servers == nil || len(servers) == 0 {
		c = NewSimplecache()
	} else {
		c = NewMemcache(servers...)
	}
	e := int(expiration/time.Second)
	return &TokenCache { cache: c, defaultExpiration : e}
}

// token as key, and tokeninfo (marshalled to string in json) as data
func (c *TokenCache) Set(token string, t *types.Token) error {
	value, err := json.Marshal(t)
	if err != nil {
		return err
	}
	return c.cache.Set(token, value, c.defaultExpiration)
}
func (c *TokenCache) Get(token string) (*types.Token, error)  {
	data, err := c.cache.Get(token)
	if err != nil {
		return nil, err
	}
	var t types.Token
	err = json.Unmarshal(data, &t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

