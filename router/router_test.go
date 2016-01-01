//Based on httprouter with modification for context based handler
//Replace httprouter's Param with golang's context library

package router

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"golang.org/x/net/context"
)

type mockResponseWriter struct{}

func (m *mockResponseWriter) Header() (h http.Header) {
	return http.Header{}
}

func (m *mockResponseWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (m *mockResponseWriter) WriteString(s string) (n int, err error) {
	return len(s), nil
}

func (m *mockResponseWriter) WriteHeader(int) {}

func TestParams(t *testing.T) {
	ps := pathParams{
		pathParam{"param1", "value1"},
		pathParam{"param2", "value2"},
		pathParam{"param3", "value3"},
	}
	for i := range ps {
		if val := ps.byName(ps[i].Key); val != ps[i].Value {
			t.Errorf("Wrong value for %s: Got %s; Want %s", ps[i].Key, val, ps[i].Value)
		}
	}
	if val := ps.byName("noKey"); val != "" {
		t.Errorf("Expected empty string for not found key; got: %s", val)
	}
}

func TestRouter(t *testing.T) {
	router := NewRouter()

	routed := false
	router.ContextHandlerFunc("GET", "/user/:name", func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		routed = true
		want := pathParams{pathParam{"name", "gopher"}}
		if !reflect.DeepEqual(ctx.Value(PathParamsKey).(pathParams), want) {
			t.Fatalf("wrong wildcard values: want %v, got %v", want, ctx.Value(PathParamsKey).(pathParams))
		}
	})

	w := new(mockResponseWriter)

	req, _ := http.NewRequest("GET", "/user/gopher", nil)
	router.ServeHTTP(w, req)

	if !routed {
		t.Fatal("routing failed")
	}
}

type handlerStruct struct {
	handeled *bool
}

func (h handlerStruct) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	*h.handeled = true
}

func TestRouterAPI(t *testing.T) {
	var get, head, options, post, put, patch, delete, handler, handlerFunc bool

	httpHandler := handlerStruct{&handler}

	router := NewRouter()
	router.GET("/GET", func(_ context.Context, w http.ResponseWriter, r *http.Request) {
		get = true
	})
	router.HEAD("/GET", func(_ context.Context, w http.ResponseWriter, r *http.Request) {
		head = true
	})
	router.OPTIONS("/GET", func(_ context.Context, w http.ResponseWriter, r *http.Request) {
		options = true
	})
	router.POST("/POST", func(_ context.Context, w http.ResponseWriter, r *http.Request) {
		post = true
	})
	router.PUT("/PUT", func(_ context.Context, w http.ResponseWriter, r *http.Request) {
		put = true
	})
	router.PATCH("/PATCH", func(_ context.Context, w http.ResponseWriter, r *http.Request) {
		patch = true
	})
	router.DELETE("/DELETE", func(_ context.Context, w http.ResponseWriter, r *http.Request) {
		delete = true
	})
	router.HTTPHandler("GET", "/Handler", httpHandler)
	router.HTTPHandlerFunc("GET", "/HandlerFunc", func(w http.ResponseWriter, r *http.Request) {
		handlerFunc = true
	})

	w := new(mockResponseWriter)

	r, _ := http.NewRequest("GET", "/GET", nil)
	router.ServeHTTP(w, r)
	if !get {
		t.Error("routing GET failed")
	}

	r, _ = http.NewRequest("HEAD", "/GET", nil)
	router.ServeHTTP(w, r)
	if !head {
		t.Error("routing HEAD failed")
	}

	r, _ = http.NewRequest("OPTIONS", "/GET", nil)
	router.ServeHTTP(w, r)
	if !options {
		t.Error("routing OPTIONS failed")
	}

	r, _ = http.NewRequest("POST", "/POST", nil)
	router.ServeHTTP(w, r)
	if !post {
		t.Error("routing POST failed")
	}

	r, _ = http.NewRequest("PUT", "/PUT", nil)
	router.ServeHTTP(w, r)
	if !put {
		t.Error("routing PUT failed")
	}

	r, _ = http.NewRequest("PATCH", "/PATCH", nil)
	router.ServeHTTP(w, r)
	if !patch {
		t.Error("routing PATCH failed")
	}

	r, _ = http.NewRequest("DELETE", "/DELETE", nil)
	router.ServeHTTP(w, r)
	if !delete {
		t.Error("routing DELETE failed")
	}

	r, _ = http.NewRequest("GET", "/Handler", nil)
	router.ServeHTTP(w, r)
	if !handler {
		t.Error("routing Handler failed")
	}

	r, _ = http.NewRequest("GET", "/HandlerFunc", nil)
	router.ServeHTTP(w, r)
	if !handlerFunc {
		t.Error("routing HandlerFunc failed")
	}
}

func TestRouterRoot(t *testing.T) {
	router := NewRouter()
	recv := catchPanic(func() {
		router.GET("noSlashRoot", nil)
	})
	if recv == nil {
		t.Fatal("registering path not beginning with '/' did not panic")
	}
}

func TestRouterChaining(t *testing.T) {
	router1 := NewRouter()
	router2 := NewRouter()
	router1.NotFound = router2

	fooHit := false
	router1.POST("/foo", func(_ context.Context, w http.ResponseWriter, r *http.Request) {
		fooHit = true
		w.WriteHeader(http.StatusOK)
	})

	barHit := false
	router2.POST("/bar", func(_ context.Context, w http.ResponseWriter, r *http.Request) {
		barHit = true
		w.WriteHeader(http.StatusOK)
	})

	r, _ := http.NewRequest("POST", "/foo", nil)
	w := httptest.NewRecorder()
	router1.ServeHTTP(w, r)
	if !(w.Code == http.StatusOK && fooHit) {
		t.Errorf("Regular routing failed with router chaining.")
		t.FailNow()
	}

	r, _ = http.NewRequest("POST", "/bar", nil)
	w = httptest.NewRecorder()
	router1.ServeHTTP(w, r)
	if !(w.Code == http.StatusOK && barHit) {
		t.Errorf("Chained routing failed with router chaining.")
		t.FailNow()
	}

	r, _ = http.NewRequest("POST", "/qax", nil)
	w = httptest.NewRecorder()
	router1.ServeHTTP(w, r)
	if !(w.Code == http.StatusNotFound) {
		t.Errorf("NotFound behavior failed with router chaining.")
		t.FailNow()
	}
}

func TestRouterNotAllowed(t *testing.T) {
	handlerFunc := func(_ context.Context, _ http.ResponseWriter, _ *http.Request) {}

	router := NewRouter()
	router.POST("/path", handlerFunc)

	// Test not allowed
	r, _ := http.NewRequest("GET", "/path", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if !(w.Code == http.StatusMethodNotAllowed) {
		t.Errorf("NotAllowed handling failed: Code=%d, Header=%v", w.Code, w.Header())
	}

	w = httptest.NewRecorder()
	responseText := "custom method"
	router.MethodNotAllowed = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusTeapot)
		w.Write([]byte(responseText))
	})
	router.ServeHTTP(w, r)
	if got := w.Body.String(); !(got == responseText) {
		t.Errorf("unexpected response got %q want %q", got, responseText)
	}
	if w.Code != http.StatusTeapot {
		t.Errorf("unexpected response code %d want %d", w.Code, http.StatusTeapot)
	}
}

func TestRouterNotFound(t *testing.T) {
	handlerFunc := func(_ context.Context, _ http.ResponseWriter, _ *http.Request) {}

	router := NewRouter()
	router.GET("/path", handlerFunc)
	router.GET("/dir/", handlerFunc)
	router.GET("/", handlerFunc)

	testRoutes := []struct {
		route  string
		code   int
		header string
	}{
		{"/path/", 301, "map[Location:[/path]]"},   // TSR -/
		{"/dir", 301, "map[Location:[/dir/]]"},     // TSR +/
		{"", 301, "map[Location:[/]]"},             // TSR +/
		{"/PATH", 301, "map[Location:[/path]]"},    // Fixed Case
		{"/DIR/", 301, "map[Location:[/dir/]]"},    // Fixed Case
		{"/PATH/", 301, "map[Location:[/path]]"},   // Fixed Case -/
		{"/DIR", 301, "map[Location:[/dir/]]"},     // Fixed Case +/
		{"/../path", 301, "map[Location:[/path]]"}, // CleanPath
		{"/nope", 404, ""},                         // NotFound
	}
	for _, tr := range testRoutes {
		r, _ := http.NewRequest("GET", tr.route, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, r)
		if !(w.Code == tr.code && (w.Code == 404 || fmt.Sprint(w.Header()) == tr.header)) {
			t.Errorf("NotFound handling route %s failed: Code=%d, Header=%v", tr.route, w.Code, w.Header())
		}
	}

	// Test custom not found handler
	var notFound bool
	router.NotFound = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(404)
		notFound = true
	})
	r, _ := http.NewRequest("GET", "/nope", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if !(w.Code == 404 && notFound == true) {
		t.Errorf("Custom NotFound handler failed: Code=%d, Header=%v", w.Code, w.Header())
	}

	// Test other method than GET (want 307 instead of 301)
	router.PATCH("/path", handlerFunc)
	r, _ = http.NewRequest("PATCH", "/path/", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if !(w.Code == 307 && fmt.Sprint(w.Header()) == "map[Location:[/path]]") {
		t.Errorf("Custom NotFound handler failed: Code=%d, Header=%v", w.Code, w.Header())
	}

	// Test special case where no node for the prefix "/" exists
	router = NewRouter()
	router.GET("/a", handlerFunc)
	r, _ = http.NewRequest("GET", "/", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if !(w.Code == 404) {
		t.Errorf("NotFound handling route / failed: Code=%d", w.Code)
	}
}

func TestRouterPanicHandler(t *testing.T) {
	router := NewRouter()
	panicHandled := false

	router.PanicHandler = func(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
		panicHandled = true
	}

	router.ContextHandlerFunc("PUT", "/user/:name", func(_ context.Context, _ http.ResponseWriter, _ *http.Request) {
		panic("oops!")
	})

	w := new(mockResponseWriter)
	req, _ := http.NewRequest("PUT", "/user/gopher", nil)

	defer func() {
		if rcv := recover(); rcv != nil {
			t.Fatal("handling panic failed")
		}
	}()

	router.ServeHTTP(w, req)

	if !panicHandled {
		t.Fatal("simulating failed")
	}
}

func TestRouterLookup(t *testing.T) {
	routed := false
	wantHandle := func(_ context.Context, _ http.ResponseWriter, _ *http.Request) {
		routed = true
	}
	wantParams := pathParams{pathParam{"name", "gopher"}}

	router := NewRouter()

	// try empty router first
	handle, _, tsr := router.Lookup("GET", "/nope")
	if handle != nil {
		t.Fatalf("Got handle for unregistered pattern: %v", handle)
	}
	if tsr {
		t.Error("Got wrong TSR recommendation!")
	}

	// insert route and try again
	router.GET("/user/:name", wantHandle)

	handle, ctx, tsr := router.Lookup("GET", "/user/gopher")
	if handle == nil {
		t.Fatal("Got no handle!")
	} else {
		handle(context.Background(), nil, nil)
		if !routed {
			t.Fatal("Routing failed!")
		}
	}

	if !reflect.DeepEqual(ctx.Value(PathParamsKey).(pathParams), wantParams) {
		t.Fatalf("Wrong parameter values: want %v, got %v", wantParams, ctx.Value(PathParamsKey).(pathParams))
	}

	handle, _, tsr = router.Lookup("GET", "/user/gopher/")
	if handle != nil {
		t.Fatalf("Got handle for unregistered pattern: %v", handle)
	}
	if !tsr {
		t.Error("Got no TSR recommendation!")
	}

	handle, _, tsr = router.Lookup("GET", "/nope")
	if handle != nil {
		t.Fatalf("Got handle for unregistered pattern: %v", handle)
	}
	if tsr {
		t.Error("Got wrong TSR recommendation!")
	}
}

type mockFileSystem struct {
	opened bool
}

func (mfs *mockFileSystem) Open(name string) (http.File, error) {
	mfs.opened = true
	return nil, errors.New("this is just a mock")
}

func TestRouterServeFiles(t *testing.T) {
	router := NewRouter()
	mfs := &mockFileSystem{}

	recv := catchPanic(func() {
		router.ServeFiles("/noFilepath", mfs)
	})
	if recv == nil {
		t.Fatal("registering path not ending with '*filepath' did not panic")
	}

	router.ServeFiles("/*filepath", mfs)
	w := new(mockResponseWriter)
	r, _ := http.NewRequest("GET", "/favicon.ico", nil)
	router.ServeHTTP(w, r)
	if !mfs.opened {
		t.Error("serving file failed")
	}
}

func routerTagMiddleware(tag string, t *testing.T) Middleware {
	return MiddlewareFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
		ctx = SetMiddlewareParam(ctx, tag, tag)
		if (tag == "t4") {
			t1, ok := MiddlewareParam(ctx, "t1").(string)
			if !ok || t1 != "t1" {
				t.Fatal("context t4 test failed: want t1, got nil" )
			}
			t2, ok := MiddlewareParam(ctx, "t2").(string)
			if !ok || t2 != "t2" {
				t.Fatal("context t4 test failed: want t2, got nil" )
			}
			t3, ok := MiddlewareParam(ctx, "t3").(string)
			if !ok || t3 != "t3" {
				t.Fatal("context t4 test failed: want t3, got nil" )
			}
			t4, ok := MiddlewareParam(ctx, tag).(string)
			if !ok || t4 != "t4" {
				t.Fatal("context t4 test failed: want t4, got nil" )
			}
		}
		return ctx
	})
}
func TestCommonMiddleware(t *testing.T) {
	RegisterCommonMiddleware(routerTagMiddleware("t1", t))
	RegisterCommonMiddleware(routerTagMiddleware("t2", t))
	RegisterCommonAfterware(routerTagMiddleware("t3", t))
	RegisterCommonAfterware(routerTagMiddleware("t4", t))
	
	router := NewRouter()
	router.HandleFunc("GET", "/", func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	})

	w := new(mockResponseWriter)

	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)
}
