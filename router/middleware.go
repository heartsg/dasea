//Based on kami & alice
package router

import (
	"net/http"
	"fmt"
	"golang.org/x/net/context"
)

// Middleware is a function that takes the current request context and returns a new request context.
// You can use middleware to build your context before your handler handles a request.
// As a special case, middleware that returns nil will halt middleware and handler execution.
//
// ToDo: Make LogHandler a special middleware/afterware and let it still run even others return nil.
//
type Middleware interface {
	ServeHTTPContext(context.Context, http.ResponseWriter, *http.Request) context.Context
}
type MiddlewareFunc func(context.Context, http.ResponseWriter, *http.Request) context.Context
func (m MiddlewareFunc) ServeHTTPContext(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	return m(ctx, w, r)
}

// MiddlewareType represents types that can convert to Middleware.
// The following concrete types are accepted:
// 	- Middleware
// 	- func(context.Context, http.ResponseWriter, *http.Request) context.Context
// 	- func(http.Handler) http.Handler [will run sequentially, not in a chain]
// 	- func(http.ContextHandler) http.ContextHandler [will run sequentially, not in a chain]
type MiddlewareType interface{}

// MiddlewareWrap turns standard http middleware into context based Middleware if needed.
// These standard http middleware has forms of func(http.Handler) http.Handler.
// One example of standard http middleware would be authentication.
// func Auth(next http.Handler) http.Handler {
//   fn := func(w http.ResponseWriter, r *http.Request) {
//     token := globalContext.Get("token")
//     if user, ok := validate(token); ok {
//       globalContext.Set("user", user)
//       next(w, r)
//     } else {
//       w.Write("auth failed")
//     }
//   }
//   return http.HandlerFunc(fn)
// }
// The above standard middleware generally makes use of a global context (protected by mutex).
// 
// To convert standard middleware into context based middleware, we have to know whether next(w, r)
// is called, so we pass in a dummyHandler which will be true if next(w, r) is called.
//
// Note that such conersion is not perfect.
//
// The middleware conversion is best for those standard middlewares that do not have
// dependencies on global context. Because there's no way for us to know which global context
// they use (we will consider to support some well-known context such as gorilla's context?), 
// if the standard middleware output some value into global context, such as the user above, 
// we cannot put the user into ctx and pass to next middleware. The remaining middlewares or
// applications that needs the "user" variable has to retrieve it from the global context.
// 
func MiddlewareWrap(mw MiddlewareType) Middleware {
	switch x := mw.(type) {
	case Middleware:
		return x
	case func(context.Context, http.ResponseWriter, *http.Request) context.Context:
		return MiddlewareFunc(x)
	case func(ContextHandler) ContextHandler:
		return MiddlewareFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
			var dh dummyHandler
			x(&dh).ServeHTTPContext(ctx, w, r)
			if !dh {
				return nil
			}
			return ctx
		})
	case func(http.Handler) http.Handler:
		return MiddlewareFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
			var dh dummyHandler
			x(&dh).ServeHTTP(w, r)
			if !dh {
				return nil
			}
			return ctx
		})
	}
	panic(fmt.Errorf("unsupported MiddlewareType: %T", mw))
}

// Convert a series of middlewares into one middleware
// Once an intermediate middleware returns nil, we shall safely assume that
// the middleware already takes care of response, and propergating along the chain is not
// necessary after that middleware.
// One example would be auth failed, we can return immediately without caring about the remaining
// middleware
func MiddlewareChain(middlewares ... Middleware) Middleware {
	return MiddlewareFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
		for _, mw := range middlewares {
			ctx = mw.ServeHTTPContext(ctx, w, r)
			if ctx == nil {
				return nil
			}
		}
		return ctx
	})
}
// Difference between afterware chain and middleware chain is that,
// the chain will not be broken if one afterware returns nil
func AfterwareChain(middlewares ... Middleware) Middleware {
	return MiddlewareFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
		for _, mw := range middlewares {
			ctx = mw.ServeHTTPContext(ctx, w, r)
		}
		return ctx
	})
}
func MiddlewareHandlerChain(handler ContextHandler, middlewares... Middleware) ContextHandlerFunc {
	return ContextHandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		ctx = MiddlewareChain(middlewares...).ServeHTTPContext(ctx, w, r)
		if ctx == nil {
			return
		}
		handler.ServeHTTPContext(ctx, w, r)
	})
}
//This is the only way to add post middlewares. There is no way to use variadic variables, have to use slice
func MiddlewareHandlerAfterwareChain(middlewares[] Middleware, handler ContextHandler, afterwares[] Middleware) ContextHandlerFunc {
	return ContextHandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		ctx = MiddlewareChain(middlewares...).ServeHTTPContext(ctx, w, r)
		if ctx == nil {
			return
		}
		handler.ServeHTTPContext(ctx, w, r)
		AfterwareChain(afterwares...).ServeHTTPContext(ctx, w, r)
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
var commonAfterwares []Middleware

func RegisterCommonMiddleware(middleware... Middleware) {
	commonMiddlewares = append(commonMiddlewares, middleware...)
}
func RegisterCommonAfterware(middleware...Middleware) {
	commonAfterwares = append(commonAfterwares, middleware...)
}