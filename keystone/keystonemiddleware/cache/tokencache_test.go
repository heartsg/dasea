package cache

import (
	"time"
	"testing"
	"github.com/heartsg/dasea/keystone/keystoneclient/types"
	"github.com/heartsg/dasea/keystone/keystoneclient/util"
	"github.com/heartsg/dasea/testutil"
)

func TestTokenCache(t *testing.T) {
	c := NewTokenCache(nil, 300)
	t1 := &types.Token {
		Methods: []string { "password" },
		ExpiresAt: &util.Iso8601DateTime{ time.Date(2015, time.November, 06, 15, 32, 17, 893769000, time.UTC) },
		Extras: make(map[string]interface{}),
		User: &types.User {
			Id: "423f19a4ac1e4f48bbb4180756e6eb6c",
			Name: "admin",
			Domain: &types.Domain {
				Id: "default",
				Name: "Default",
			},
		},
	}
	
	err := c.Set("123", t1)
	if err != nil {
		t.Error(err)
	}
	
	t2, err := c.Get("123")
	if err != nil {
		t.Error(err)
	}
	
	testutil.Equals(t, t1, t2)
	
}