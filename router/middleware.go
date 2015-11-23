//Based on kami & alice
package router

import (
	"net/http"
	"fmt"
	"golang.org/x/net/context"
)

// Middleware is a function that takes the current request context and returns a new request context.
// You can use middleware to build your context before your handler handles a request.
// As a special case, middleware that returns nil will halt middleware and handler execution (LogHandler will still run).
type Middleware func(context.Context, http.ResponseWriter, *http.Request) context.Context

// MiddlewareType represents types that can convert to Middleware.
// The following concrete types are accepted:
// 	- Middleware
// 	- func(context.Context, http.ResponseWriter, *http.Request) context.Context
// 	- func(http.Handler) http.Handler [will run sequentially, not in a chain]
// 	- func(http.ContextHandler) http.ContextHandler [will run sequentially, not in a chain]
type MiddlewareType interface{}

// MiddlewareWrap turns standard http middleware into context based Middleware if needed.
func MiddlewareWrap(mw MiddlewareType) Middleware {
	switch x := mw.(type) {
	case Middleware:
		return x
	case func(context.Context, http.ResponseWriter, *http.Request) context.Context:
		return Middleware(x)
	case func(ContextHandler) ContextHandler:
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
			var dh dummyHandler
			x(&dh).ServeHTTPContext(ctx, w, r)
			if !dh {
				return nil
			}
			return ctx
		}
	case func(http.Handler) http.Handler:
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
			var dh dummyHandler
			x(&dh).ServeHTTP(w, r)
			if !dh {
				return nil
			}
			return ctx
		}
	}
	panic(fmt.Errorf("unsupported MiddlewareType: %T", mw))
}

// Convert a series of middlewares and a handler into a handler
func MiddlewareChain(middlewares ... Middleware) Middleware {
	return Middleware(func(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
		for _, mw := range middlewares {
			ctx = mw(ctx, w, r)
		}
		return ctx
	})
}
func MiddlewareChainWithHandler(handler ContextHandlerFunc, middlewares... Middleware) ContextHandlerFunc {
	return ContextHandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		ctx = MiddlewareChain(middlewares...)(ctx, w, r)
		handler(ctx, w, r)
	})
}
//This is the only way to add post middlewares. There is no way to use variadic variables, have to use slice
func MiddlewareChainWithPostMiddleware(middlewares[] Middleware, handler ContextHandlerFunc, postMiddlewares[] Middleware) ContextHandlerFunc {
	return ContextHandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		ctx = MiddlewareChain(middlewares...)(ctx, w, r)
		handler(ctx, w, r)
		MiddlewareChain(postMiddlewares...)(ctx, w, r)
	})
}

// dummyHandler is used to keep track of whether the next middleware was called or not.
type dummyHandler bool

func (dh *dummyHandler) ServeHTTP(_ http.ResponseWriter, _ *http.Request) {
	*dh = true
}

func (dh *dummyHandler) ServeHTTPContext(_ context.Context, _ http.ResponseWriter, _ *http.Request) {
	*dh = true
}

var commonMiddlewares []Middleware
var commonPostMiddlewares []Middleware

func RegisterCommonMiddleware(middleware... Middleware) {
	commonMiddlewares = append(commonMiddlewares, middleware...)
}
func RegisterCommonPostMiddleware(middleware...Middleware) {
	commonPostMiddlewares = append(commonPostMiddlewares, middleware...)
}