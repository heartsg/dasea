// Based on Package alice

package router

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"golang.org/x/net/context"
)

// A constructor for middleware
// that writes its own "tag" into the RW and does nothing else.
// Useful in checking if a chain is behaving in the right order.
func tagMiddleware(tag string) Middleware {
	return MiddlewareFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
		w.Write([]byte(tag))
		return ctx
	})
}

// Not recommended (https://golang.org/pkg/reflect/#Value.Pointer),
// but the best we can do.
func funcsEqual(f1, f2 interface{}) bool {
	val1 := reflect.ValueOf(f1)
	val2 := reflect.ValueOf(f2)
	return val1.Pointer() == val2.Pointer()
}

var testApp = ContextHandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("app\n"))
})


func TestMiddlewareHandlerChainCorrectly(t *testing.T) {
	chain := MiddlewareChain(tagMiddleware("t1\n"), tagMiddleware("t2\n"))
	newChain := MiddlewareChain(tagMiddleware("t3\n"), tagMiddleware("t4\n"))

	chained := MiddlewareHandlerChain(testApp, chain, newChain)

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	chained.ServeHTTPContext(context.Background(), w, r)


	if (w.Body.String() != "t1\nt2\nt3\nt4\napp\n") {
		t.Fatalf("wrong string: want %s, got %s", "t1\nt2\nt3\nt4\napp\n",  w.Body.String())		
	}
}

func TestPostMiddlewareHandlerChainCorrectly(t *testing.T) {
	chain := MiddlewareChain(tagMiddleware("t1\n"), tagMiddleware("t2\n"))
	newChain := MiddlewareChain(tagMiddleware("t3\n"), tagMiddleware("t4\n"))

	chained := MiddlewareHandlerAfterwareChain([]Middleware{chain}, testApp, []Middleware{newChain})

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	chained.ServeHTTPContext(context.Background(), w, r)


	if (w.Body.String() != "t1\nt2\napp\nt3\nt4\n") {
		t.Fatalf("wrong string: want %s, got %s", "t1\nt2\napp\nt3\nt4\n",  w.Body.String())		
	}	
}

func contextTagMiddleware(tag string, t *testing.T) Middleware {
	return MiddlewareFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
		ctx = SetMiddlewareParam(ctx, tag, tag)
		if (tag == "t1\n") {
			t1, ok := MiddlewareParam(ctx, tag).(string)
			if !ok || t1 != "t1\n" {
				t.Fatal("context t1 test failed: want t1, got nil" )
			}
		}
		if (tag == "t2\n") {
			t1, ok := MiddlewareParam(ctx, "t1\n").(string)
			if !ok || t1 != "t1\n" {
				t.Fatal("context t2 test failed: want t1, got nil" )
			}
			t2, ok := MiddlewareParam(ctx, tag).(string)
			if !ok || t2 != "t2\n" {
				t.Fatal("context t2 test failed: want t2, got nil" )
			}
		}
		if (tag == "t3\n") {
			t1, ok := MiddlewareParam(ctx, "t1\n").(string)
			if !ok || t1 != "t1\n" {
				t.Fatal("context t3 test failed: want t1, got nil" )
			}
			t2, ok := MiddlewareParam(ctx, "t2\n").(string)
			if !ok || t2 != "t2\n" {
				t.Fatal("context t3 test failed: want t2, got nil" )
			}
			t3, ok := MiddlewareParam(ctx, tag).(string)
			if !ok || t3 != "t3\n" {
				t.Fatal("context t3 test failed: want t3, got nil" )
			}
		}
		if (tag == "t4\n") {
			t1, ok := MiddlewareParam(ctx, "t1\n").(string)
			if !ok || t1 != "t1\n" {
				t.Fatal("context t4 test failed: want t1, got nil" )
			}
			t2, ok := MiddlewareParam(ctx, "t2\n").(string)
			if !ok || t2 != "t2\n" {
				t.Fatal("context t4 test failed: want t2, got nil" )
			}
			t3, ok := MiddlewareParam(ctx, "t3\n").(string)
			if !ok || t3 != "t3\n" {
				t.Fatal("context t4 test failed: want t3, got nil" )
			}
			t4, ok := MiddlewareParam(ctx, tag).(string)
			if !ok || t4 != "t4\n" {
				t.Fatal("context t4 test failed: want t4, got nil" )
			}
		}
		return ctx
	})
}
func TestContextCorrectly(t *testing.T) {
	chain := MiddlewareChain(contextTagMiddleware("t1\n", t), contextTagMiddleware("t2\n", t))
	newChain := MiddlewareChain(contextTagMiddleware("t3\n", t), contextTagMiddleware("t4\n", t))
	
	chained := MiddlewareHandlerAfterwareChain([]Middleware{chain}, testApp, []Middleware{newChain})
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	chained.ServeHTTPContext(context.Background(), w, r)
}