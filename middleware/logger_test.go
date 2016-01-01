package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"golang.org/x/net/context"
	"github.com/heartsg/dasea/router"
)

var testApp = router.ContextHandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
})

func TestLogger(t *testing.T) {
	chained := router.MiddlewareHandlerAfterwareChain([]router.Middleware{ router.MiddlewareFunc(StartLogger) },
		testApp, []router.Middleware{ router.MiddlewareFunc(EndLogger) })

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	chained.ServeHTTPContext(context.Background(), w, r)	
}