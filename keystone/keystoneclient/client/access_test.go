package client

import (
	"testing"
    "net/http"
	"net/http/httptest"
    "time"
    "io/ioutil"
    
    "github.com/heartsg/dasea/keystone/keystoneclient/types"
    "github.com/heartsg/dasea/keystone/keystoneclient/util"
)

func TestAccessInfoCreation(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("X-Subject-Token", "12345")
        resRaw := `{
    "token": {
        "methods": [
            "password"
        ],
        "expires_at": "2015-11-06T15:32:17.893769Z",
        "extras": {},
        "user": {
            "domain": {
                "id": "default",
                "name": "Default"
            },
            "id": "423f19a4ac1e4f48bbb4180756e6eb6c",
            "name": "admin"
        },
        "audit_ids": [
            "ZzZwkUflQfygX7pdYDBCQQ"
        ],
        "issued_at": "2015-11-06T14:32:17.893797Z"
    }
}`
    
        w.Write([]byte(resRaw))
    }))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Error(err)
	}
    defer res.Body.Close()
    
    access, err := NewAccessFromResponse(res)
    if err != nil {
		t.Error(err)
	}
    
	if access.Token != "12345" {
        t.Error("token wrong")
    }
    
    if access.TokenInfo.User.Domain.Id != "default" {
        t.Error("token info wrong")
    }
    
    if !access.WillExpireSoon() {
        t.Error("should expire!")
    }
}

func TestAccessInfoCreation2(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("X-Subject-Token", "12345")
        resRaw := `{
    "token": {
        "methods": [
            "password"
        ],
        "expires_at": "2015-11-06T15:32:17.893769Z",
        "extras": {},
        "user": {
            "domain": {
                "id": "default",
                "name": "Default"
            },
            "id": "423f19a4ac1e4f48bbb4180756e6eb6c",
            "name": "admin"
        },
        "audit_ids": [
            "ZzZwkUflQfygX7pdYDBCQQ"
        ],
        "issued_at": "2015-11-06T14:32:17.893797Z"
    }
}`
    
        w.Write([]byte(resRaw))
    }))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Error(err)
	}
    defer res.Body.Close()
    
    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        t.Error(err)
    }
    access, err := NewAccessFromResponseBody(res, body)
    if err != nil {
		t.Error(err)
	}
    
	if access.Token != "12345" {
        t.Error("token wrong")
    }
    
    if access.TokenInfo.User.Domain.Id != "default" {
        t.Error("token info wrong")
    }
    
    if !access.WillExpireSoon() {
        t.Error("should expire!")
    }
}


func TestExpire(t *testing.T) {
    access := NewAccess("12345", &types.Token {
        ExpiresAt: &util.Iso8601DateTime{time.Now().Add(time.Duration(StaleDuration+5)*time.Second)},
    })
    if access.WillExpireSoon() {
        t.Error("should not expire")
    }
    access.TokenInfo.ExpiresAt = &util.Iso8601DateTime{access.TokenInfo.ExpiresAt.Add(time.Duration(-6)*time.Second)}
    if !access.WillExpireSoon() {
        t.Errorf("should expire: now=%s expiresat=%s", time.Now(), access.TokenInfo.ExpiresAt)
    }
}