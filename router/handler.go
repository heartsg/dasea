//Package router is in charge of routing and middleware supports

package router

import (
	"fmt"
	"net/http"
	"golang.org/x/net/context"
)

// Dasea internally uses golang.org/x/net/context for context management. To work
// with other frameworks, we uses wrappers from other handler types to ContextHandler.
//
// Inspired by and based on github.com/guregu/kami
//

// HandlerType is the type of Handlers and types converted to ContextHandler.
// Currently only the following concrete types are accepted (will support more later on):
// 	- types that implement http.Handler
// 	- types that implement ContextHandler
// 	- func(http.ResponseWriter, *http.Request)
// 	- func(context.Context, http.ResponseWriter, *http.Request)
type HandlerType interface{}

// ContextHandler is like http.Handler but supports context.
type ContextHandler interface {
	ServeHTTPContext(context.Context, http.ResponseWriter, *http.Request)
}

// HandlerFunc is like http.HandlerFunc with context.
type ContextHandlerFunc func(context.Context, http.ResponseWriter, *http.Request)

// ContextHandlerFunc implements ContextHandler interface
func (h ContextHandlerFunc) ServeHTTPContext(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	h(ctx, w, r)
}

// Make ContextHandlerFunc compatible to net/http by implementing http.Handler interface
func (h ContextHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h(context.Background(), w, r)
}

// Turn a HandlerType into a ContextHandler
func ContextHandlerWrap(h HandlerType) ContextHandler {
	switch x := h.(type) {
	case ContextHandler:
		return x
	case func(context.Context, http.ResponseWriter, *http.Request):
		return ContextHandlerFunc(x)
	case http.Handler:
		return ContextHandlerFunc(func(_ context.Context, w http.ResponseWriter, r *http.Request) {
			x.ServeHTTP(w, r)
		})
	case func(http.ResponseWriter, *http.Request):
		return ContextHandlerFunc(func(_ context.Context, w http.ResponseWriter, r *http.Request) {
			x(w, r)
		})
	}
	panic(fmt.Errorf("unsupported HandlerType: %T", h))
}

// Turn a HandlerType into a http.Handler
func HTTPHandlerWrap(h HandlerType) http.Handler {
	switch x := h.(type) {
	case http.Handler:
		return x
	case func(http.ResponseWriter, *http.Request):
		return http.HandlerFunc(x)
	case ContextHandler:
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			x.ServeHTTPContext(context.Background(), w, r)
		})
	case func(context.Context, http.ResponseWriter, *http.Request):
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			x(context.Background(), w, r)
		})
	}
	panic(fmt.Errorf("unsupported HandlerType: %T", h))
}