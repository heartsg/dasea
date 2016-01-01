package keystonemiddleware

import (
	"testing"
	"time"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"github.com/heartsg/dasea/router"
	"golang.org/x/net/context"
	"github.com/heartsg/dasea/keystone/keystoneclient/client"
	"github.com/heartsg/dasea/keystone/keystoneclient"
)

func TestUserToken(t *testing.T) {
	// We will create two servers, one is for authentication (token <-> user info)
	// One is where AuthToken middleware is running
	// Client connects to AuthToken middleware server, which authenticates the client
	// from the authentication server.
	
	// The server will do 3 things:
	// - Auth a service user with token
	// - Get a user token with service token as authtoken, return a normal user info
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/auth/tokens":
			if r.Method == "POST" {
				_, _ = ioutil.ReadAll(r.Body)
				defer r.Body.Close()
			
				w.Header().Set("X-Subject-Token", "servicetoken")
				w.Write([]byte(`{
					"token": {
						"methods": [
							"password"
						],
						"expires_at": "2016-11-06T15:32:17.893769Z",
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
				}`))
			} else if r.Method == "GET" {
				if r.Header.Get("X-Auth-Token") != "servicetoken" {
					t.Errorf("Expected 'X-Auth-Token' == %q; got %q", "servicetoken", r.Header.Get("X-Auth-Token"))
				}
				if r.Header.Get("X-Subject-Token") != "usertoken" {
					t.Errorf("Expected 'X-Subject-Token' == %q; got %q", "usertoken", r.Header.Get("X-Subject-Token"))			
				}
				
				w.Header().Set("X-Auth-Token", "servicetoken")
				w.Header().Set("X-Subject-Token", "usertoken")
				w.Write([]byte(`{
					"token": {
						"methods": [
							"token"
						],
						"expires_at": "2016-11-05T22:00:11.000000Z",
						"extras": {},
						"user": {
							"domain": {
								"id": "default",
								"name": "Default"
							},
							"id": "10a2e6e717a245d9acad3e5f97aeca3d",
							"name": "testuser"
						},
						"audit_ids": [
							"mAjXQhiYRyKwkB4qygdLVg"
						],
						"issued_at": "2015-11-05T21:00:33.819948Z"
					}
				}`))
			} else {
				t.Errorf("Expected method GET or POST; got %q", r.Method)
			}
		default:
			t.Errorf("Invalid path %q", r.URL.Path)
		}
	}))
		
	defer ts.Close()
	
	// Create a middleware and app, test whether the middleware will be able to 
	// auto handle service token retrieval and authentication for us.
	opts := Opts {
		AuthMethod: "password",
		Username: "admin",
		Password: "devstacker",
		UserDomainId: "default",
		Client: keystoneclient.Opts {
			AuthUrl: ts.URL,
		},
		DelayAuthDecision: true,
		TokenCacheTime: time.Duration(300)*time.Second,
	}
	authToken := NewAuthToken(&opts)
   	handler := router.MiddlewareHandlerChain(router.ContextHandlerFunc(appUserToken), authToken)
	ts_service := httptest.NewServer(
		//wrap handler into an http handler, also insert t into ctx
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.Background()
			ctx = context.WithValue(ctx, "test", t)
			handler(ctx, w, r)
		}))
	defer ts_service.Close()
	
	// create a client and request to ts_service
	client := &http.Client{}
	req, _ := http.NewRequest("GET", ts_service.URL, nil)
	req.Header.Add("X-Auth-Token", "usertoken")
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	
	body, _ := ioutil.ReadAll(resp.Body)
	if string(body) != "Test success!" {
		t.Errorf("Response error, should be Test Success!, got: %v", body)
	}
}

func appUserToken(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	// get t from ctx
	t := ctx.Value("test").(*testing.T)
	
	// Test whether the user is authenticated
	if r.Header.Get("X-Identity-Status") != "Confirmed" {
		t.Errorf("At the app, the user token should already be confirmed, got: %v", r.Header.Get("X-Identity-Status"))
	}
	if r.Header.Get("X-User-Id") != "10a2e6e717a245d9acad3e5f97aeca3d" {
		t.Errorf("At the app, the user should be 10a2e6e717a245d9acad3e5f97aeca3d, got: %v", r.Header.Get("X-User-Id"))
	}
	if r.Header.Get("X-User-Name") != "testuser" {
		t.Errorf("At the app, the user should be testuser, got: %v", r.Header.Get("X-User-Name"))
	}
	if r.Header.Get("X-Domain-Id") != "default" {
		t.Errorf("At the app, the user should be default, got: %v", r.Header.Get("X-Domain-Id"))
	}
	
	// Test ctx's Context Parameter with key "UserAccessInfo"
	value := router.MiddlewareParam(ctx, UserAccessInfoKey)
	if value == nil {
		t.Error("ctx should contain user access info")
	} else {
		access, ok := value.(*client.AccessInfo)
		if !ok {
			t.Error("it is not accessinfo, what is it?")
		}
		if access.Token != "usertoken" || access.TokenInfo.User.Domain.Name != "Default" {
			t.Error("ctx's accessinfo contains wrong information")
		}
	}
	
	w.Write([]byte("Test success!"))
}

func TestUserTokenInvalid(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/auth/tokens":
			if r.Method == "POST" {
				_, _ = ioutil.ReadAll(r.Body)
				defer r.Body.Close()
			
				w.Header().Set("X-Subject-Token", "servicetoken")
				w.Write([]byte(`{
					"token": {
						"methods": [
							"password"
						],
						"expires_at": "2016-11-06T15:32:17.893769Z",
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
				}`))
			} else if r.Method == "GET" {
				if r.Header.Get("X-Auth-Token") != "servicetoken" {
					t.Errorf("Expected 'X-Auth-Token' == %q; got %q", "servicetoken", r.Header.Get("X-Auth-Token"))
				}
				if r.Header.Get("X-Subject-Token") != "usertoken" {
					t.Errorf("Expected 'X-Subject-Token' == %q; got %q", "usertoken", r.Header.Get("X-Subject-Token"))			
				}
				//return invalid user!
				w.WriteHeader(http.StatusNotFound)
			} else {
				t.Errorf("Expected method GET or POST; got %q", r.Method)
			}
		default:
			t.Errorf("Invalid path %q", r.URL.Path)
		}
	}))
		
	defer ts.Close()
	
	// Create a middleware and app, test whether the middleware will be able to 
	// auto handle service token retrieval and authentication for us.
	opts := Opts {
		AuthMethod: "password",
		Username: "admin",
		Password: "devstacker",
		UserDomainId: "default",
		Client: keystoneclient.Opts {
			AuthUrl: ts.URL,
		},
		DelayAuthDecision: true,
		TokenCacheTime: time.Duration(300)*time.Second,
	}
	authToken := NewAuthToken(&opts)
   	handler := router.MiddlewareHandlerChain(router.ContextHandlerFunc(appUserTokenInvalid), authToken)
	ts_service := httptest.NewServer(
		//wrap handler into an http handler, also insert t into ctx
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.Background()
			ctx = context.WithValue(ctx, "test", t)
			handler(ctx, w, r)
		}))
	defer ts_service.Close()
	
	// create a client and request to ts_service
	client := &http.Client{}
	req, _ := http.NewRequest("GET", ts_service.URL, nil)
	req.Header.Add("X-Auth-Token", "usertoken")
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	
	body, _ := ioutil.ReadAll(resp.Body)
	//because delay cache decision is true, we will return "test success"
	if string(body) != "Test success!" {
		t.Errorf("Response error, should be Test Success!, got: %v", body)
	}
}

func appUserTokenInvalid(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	// get t from ctx
	t := ctx.Value("test").(*testing.T)
	
	// Test whether the user is authenticated
	if r.Header.Get("X-Identity-Status") != "Invalid" {
		t.Errorf("At the app, the user token should already be invalid, got: %v", r.Header.Get("X-Identity-Status"))
	}
	if r.Header.Get("X-User-Id") != "" {
		t.Errorf("At the app, the user should be empty, got: %v", r.Header.Get("X-User-Id"))
	}
	
	// Test ctx's Context Parameter with key "UserAccessInfo"
	value := router.MiddlewareParam(ctx, UserAccessInfoKey)
	if value != nil {
		t.Error("ctx should not contain user access info")
	}
	
	w.Write([]byte("Test success!"))
}



func TestUserTokenNodelay(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/auth/tokens":
			if r.Method == "POST" {
				_, _ = ioutil.ReadAll(r.Body)
				defer r.Body.Close()
			
				w.Header().Set("X-Subject-Token", "servicetoken")
				w.Write([]byte(`{
					"token": {
						"methods": [
							"password"
						],
						"expires_at": "2016-11-06T15:32:17.893769Z",
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
				}`))
			} else if r.Method == "GET" {
				if r.Header.Get("X-Auth-Token") != "servicetoken" {
					t.Errorf("Expected 'X-Auth-Token' == %q; got %q", "servicetoken", r.Header.Get("X-Auth-Token"))
				}
				if r.Header.Get("X-Subject-Token") != "usertoken" {
					t.Errorf("Expected 'X-Subject-Token' == %q; got %q", "usertoken", r.Header.Get("X-Subject-Token"))			
				}
				//return invalid user!
				w.WriteHeader(http.StatusNotFound)
			} else {
				t.Errorf("Expected method GET or POST; got %q", r.Method)
			}
		default:
			t.Errorf("Invalid path %q", r.URL.Path)
		}
	}))
		
	defer ts.Close()
	
	// Create a middleware and app, test whether the middleware will be able to 
	// auto handle service token retrieval and authentication for us.
	opts := Opts {
		AuthMethod: "password",
		Username: "admin",
		Password: "devstacker",
		UserDomainId: "default",
		Client: keystoneclient.Opts {
			AuthUrl: ts.URL,
		},
		DelayAuthDecision: false,
		TokenCacheTime: time.Duration(300)*time.Second,
	}
	authToken := NewAuthToken(&opts)
   	handler := router.MiddlewareHandlerChain(router.ContextHandlerFunc(appUserTokenNodelay), authToken)
	ts_service := httptest.NewServer(
		//wrap handler into an http handler, also insert t into ctx
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.Background()
			ctx = context.WithValue(ctx, "test", t)
			handler(ctx, w, r)
		}))
	defer ts_service.Close()
	
	// create a client and request to ts_service
	client := &http.Client{}
	req, _ := http.NewRequest("GET", ts_service.URL, nil)
	req.Header.Add("X-Auth-Token", "usertoken")
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	
	//because delay cache decision is false, we will get "WWW-Authenticate"
	if resp.Header.Get("WWW-Authenticate") != ts.URL {
		t.Errorf("Response error, should be %v, got: %v", ts.URL, resp.Header.Get("WWW-Authenticate"))
	}
}

func appUserTokenNodelay(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	// get t from ctx
	t := ctx.Value("test").(*testing.T)
	
	t.Errorf("No delay, so it should not be here, will return by previous middleware")
}



func TestServiceToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/auth/tokens":
			if r.Method == "POST" {
				_, _ = ioutil.ReadAll(r.Body)
				defer r.Body.Close()
			
				w.Header().Set("X-Subject-Token", "servicetoken")
				w.Write([]byte(`{
					"token": {
						"methods": [
							"password"
						],
						"expires_at": "2016-11-06T15:32:17.893769Z",
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
				}`))
			} else if r.Method == "GET" {
				if r.Header.Get("X-Auth-Token") != "servicetoken" {
					t.Errorf("Expected 'X-Auth-Token' == %q; got %q", "servicetoken", r.Header.Get("X-Auth-Token"))
				}
				if r.Header.Get("X-Subject-Token") != "usertoken" &&
					r.Header.Get("X-Subject-Token") != "myservicetoken" {
					t.Errorf("Expected 'X-Subject-Token' == usertoken or myservicetoken; got %q", r.Header.Get("X-Subject-Token"))			
				}
				
				w.Header().Set("X-Auth-Token", "servicetoken")
				w.Header().Set("X-Subject-Token", r.Header.Get("X-Subject-Token"))
				w.Write([]byte(`{
					"token": {
						"methods": [
							"token"
						],
						"expires_at": "2016-11-05T22:00:11.000000Z",
						"extras": {},
						"user": {
							"domain": {
								"id": "default",
								"name": "Default"
							},
							"id": "10a2e6e717a245d9acad3e5f97aeca3d",
							"name": "testuser"
						},
						"audit_ids": [
							"mAjXQhiYRyKwkB4qygdLVg"
						],
						"issued_at": "2015-11-05T21:00:33.819948Z"
					}
				}`))
			} else {
				t.Errorf("Expected method GET or POST; got %q", r.Method)
			}
		default:
			t.Errorf("Invalid path %q", r.URL.Path)
		}
	}))
		
	defer ts.Close()
	
	// Create a middleware and app, test whether the middleware will be able to 
	// auto handle service token retrieval and authentication for us.
	opts := Opts {
		AuthMethod: "password",
		Username: "admin",
		Password: "devstacker",
		UserDomainId: "default",
		Client: keystoneclient.Opts {
			AuthUrl: ts.URL,
		},
		DelayAuthDecision: true,
		TokenCacheTime: time.Duration(300)*time.Second,
	}
	authToken := NewAuthToken(&opts)
   	handler := router.MiddlewareHandlerChain(router.ContextHandlerFunc(appServiceToken), authToken)
	ts_service := httptest.NewServer(
		//wrap handler into an http handler, also insert t into ctx
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.Background()
			ctx = context.WithValue(ctx, "test", t)
			handler(ctx, w, r)
		}))
	defer ts_service.Close()
	
	// create a client and request to ts_service
	client := &http.Client{}
	req, _ := http.NewRequest("GET", ts_service.URL, nil)
	req.Header.Add("X-Auth-Token", "usertoken")
	req.Header.Add("X-Service-Token", "myservicetoken")
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	
	body, _ := ioutil.ReadAll(resp.Body)
	if string(body) != "Test success!" {
		t.Errorf("Response error, should be Test Success!, got: %v", body)
	}
}

func appServiceToken(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	// get t from ctx
	t := ctx.Value("test").(*testing.T)
	
	// Test whether the user is authenticated
	if r.Header.Get("X-Identity-Status") != "Confirmed" {
		t.Errorf("At the app, the user token should already be confirmed, got: %v", r.Header.Get("X-Identity-Status"))
	}
	if r.Header.Get("X-Service-Identity-Status") != "Confirmed" {
		t.Errorf("At the app, the service token should already be confirmed, got: %v", r.Header.Get("X-Service-Identity-Status"))
	}
	if r.Header.Get("X-User-Id") != "10a2e6e717a245d9acad3e5f97aeca3d" {
		t.Errorf("At the app, the user should be 10a2e6e717a245d9acad3e5f97aeca3d, got: %v", r.Header.Get("X-User-Id"))
	}
	if r.Header.Get("X-Service-User-Id") != "10a2e6e717a245d9acad3e5f97aeca3d" {
		t.Errorf("At the app, the service should be 10a2e6e717a245d9acad3e5f97aeca3d, got: %v", r.Header.Get("X-Service-User-Id"))
	}
	if r.Header.Get("X-User-Name") != "testuser" {
		t.Errorf("At the app, the user should be testuser, got: %v", r.Header.Get("X-User-Name"))
	}
	if r.Header.Get("X-Domain-Id") != "default" {
		t.Errorf("At the app, the user should be default, got: %v", r.Header.Get("X-Domain-Id"))
	}
	
	// Test ctx's Context Parameter with key "UserAccessInfo"
	value := router.MiddlewareParam(ctx, UserAccessInfoKey)
	if value == nil {
		t.Error("ctx should contain user access info")
	} else {
		access, ok := value.(*client.AccessInfo)
		if !ok {
			t.Error("it is not accessinfo, what is it?")
		}
		if access.Token != "usertoken" || access.TokenInfo.User.Domain.Name != "Default" {
			t.Error("ctx's accessinfo contains wrong information")
		}
	}
	
	w.Write([]byte("Test success!"))
}
