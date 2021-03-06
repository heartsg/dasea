//Based on httprouter with modification for context based handler
//Replace httprouter's Param with golang's context library


// Package router is a trie based high performance HTTP request router.
//
// A trivial example is:
//
//  package main
//
//  import (
//      "fmt"
//      "dasea/router"
//      "net/http"
//      "log"
//		"golang.org/x/net/context
//  )
//
//  func Index(ctx context.Context, w http.ResponseWriter, r *http.Request) {
//      fmt.Fprint(w, "Welcome!\n")
//  }
//
//  func Hello(ctx context.Context, w http.ResponseWriter, r *http.Request) {
//      fmt.Fprintf(w, "hello, %s!\n", router.PathParam(ctx, "name"))
//  }
//
//  func main() {
//      router := router.New()
//      router.GET("/", Index)
//      router.GET("/hello/:name", Hello)
//
//      log.Fatal(http.ListenAndServe(":8080", router))
//  }
//
// The router matches incoming requests by the request method and the path.
// If a handle is registered for this path and method, the router delegates the
// request to that function.
// For the methods GET, POST, PUT, PATCH and DELETE shortcut functions exist to
// register handles, for all other methods router.Handle can be used.
//
// The registered path, against which the router matches incoming requests, can
// contain two types of parameters:
//  Syntax    Type
//  :name     named parameter
//  *name     catch-all parameter
//
// Named parameters are dynamic path segments. They match anything until the
// next '/' or the path end:
//  Path: /blog/:category/:post
//
//  Requests:
//   /blog/go/request-routers            match: category="go", post="request-routers"
//   /blog/go/request-routers/           no match, but the router would redirect
//   /blog/go/                           no match
//   /blog/go/request-routers/comments   no match
//
// Catch-all parameters match anything until the path end, including the
// directory index (the '/' before the catch-all). Since they match anything
// until the end, catch-all parameters must always be the final path element.
//  Path: /files/*filepath
//
//  Requests:
//   /files/                             match: filepath="/"
//   /files/LICENSE                      match: filepath="/LICENSE"
//   /files/templates/article.html       match: filepath="/templates/article.html"
//   /files                              no match, but the router would redirect
//
// The value of parameters is saved as a slice of the Param struct, consisting
// each of a key and a value. The slice is passed to the Handle func as a third
// parameter.
// There are two ways to retrieve the value of a parameter:
//  // by the name of the parameter
//  user := PathParam(ctx, "user") // defined by :user or *user
//
//  // by the index of the parameter. This way you can also get the name (key)
//  //not implemented yet
package router

import (
	"net/http"
	"golang.org/x/net/context"
)

// Router is a http.Handler which can be used to dispatch requests to different
// handler functions via configurable routes
type Router struct {
	trees map[string]*node

	// Enables automatic redirection if the current route can't be matched but a
	// handler for the path with (without) the trailing slash exists.
	// For example if /foo/ is requested but a route only exists for /foo, the
	// client is redirected to /foo with http status code 301 for GET requests
	// and 307 for all other request methods.
	RedirectTrailingSlash bool

	// If enabled, the router tries to fix the current request path, if no
	// handle is registered for it.
	// First superfluous path elements like ../ or // are removed.
	// Afterwards the router does a case-insensitive lookup of the cleaned path.
	// If a handle can be found for this route, the router makes a redirection
	// to the corrected path with status code 301 for GET requests and 307 for
	// all other request methods.
	// For example /FOO and /..//Foo could be redirected to /foo.
	// RedirectTrailingSlash is independent of this option.
	RedirectFixedPath bool

	// If enabled, the router checks if another method is allowed for the
	// current route, if the current request can not be routed.
	// If this is the case, the request is answered with 'Method Not Allowed'
	// and HTTP status code 405.
	// If no other Method is allowed, the request is delegated to the NotFound
	// handler.
	HandleMethodNotAllowed bool

	// Configurable http.Handler which is called when no matching route is
	// found. If it is not set, http.NotFound is used.
	NotFound http.Handler

	// Configurable http.Handler which is called when a request
	// cannot be routed and HandleMethodNotAllowed is true.
	// If it is not set, http.Error with http.StatusMethodNotAllowed is used.
	MethodNotAllowed http.Handler

	// Function to handle panics recovered from http handlers.
	// It should be used to generate a error page and return the http error code
	// 500 (Internal Server Error).
	// The handler can be used to keep your server from crashing because of
	// unrecovered panics.
	PanicHandler func(context.Context, http.ResponseWriter, *http.Request)
}

// New returns a new initialized Router.
// Path auto-correction, including trailing slashes, is enabled by default.
func NewRouter() *Router {
	return &Router{
		RedirectTrailingSlash:  true,
		RedirectFixedPath:      true,
		HandleMethodNotAllowed: true,
	}
}

// GET is a shortcut for router.Handle("GET", path, handle)
func (r *Router) GET(path string, handle ContextHandlerFunc) {
	r.ContextHandlerFunc("GET", path, handle)
}

// HEAD is a shortcut for router.Handle("HEAD", path, handle)
func (r *Router) HEAD(path string, handle ContextHandlerFunc) {
	r.ContextHandlerFunc("HEAD", path, handle)
}

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handle)
func (r *Router) OPTIONS(path string, handle ContextHandlerFunc) {
	r.ContextHandlerFunc("OPTIONS", path, handle)
}

// POST is a shortcut for router.Handle("POST", path, handle)
func (r *Router) POST(path string, handle ContextHandlerFunc) {
	r.ContextHandlerFunc("POST", path, handle)
}

// PUT is a shortcut for router.Handle("PUT", path, handle)
func (r *Router) PUT(path string, handle ContextHandlerFunc) {
	r.ContextHandlerFunc("PUT", path, handle)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle)
func (r *Router) PATCH(path string, handle ContextHandlerFunc) {
	r.ContextHandlerFunc("PATCH", path, handle)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle)
func (r *Router) DELETE(path string, handle ContextHandlerFunc) {
	r.ContextHandlerFunc("DELETE", path, handle)
}

// Handle registers a new request handle with the given path and method.
//
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut
// functions can be used.
//
// This function is intended for bulk loading and to allow the usage of less
// frequently used, non-standardized or custom methods (e.g. for internal
// communication with a proxy).
func (r *Router) HandleFunc(method, path string, handle ContextHandlerFunc) {
	r.ContextHandlerFunc(method, path, handle)
}
func (r *Router) Handle(method, path string, handle ContextHandler) {
	r.ContextHandler(method, path, handle)
}
func (r *Router) ContextHandlerFunc(method, path string, handle ContextHandlerFunc) {
	if path[0] != '/' {
		panic("path must begin with '/' in path '" + path + "'")
	}

	if r.trees == nil {
		r.trees = make(map[string]*node)
	}

	root := r.trees[method]
	if root == nil {
		root = new(node)
		r.trees[method] = root
	}
	//Before put handle, we put common middlewares around the handle
	handle = MiddlewareHandlerAfterwareChain(commonMiddlewares, handle, commonAfterwares)
	
	root.addRoute(path, handle)
}

// ContextHandler is an adapter which allows the usage of an router.ContextHandler
// as a request handle.
func (r *Router) ContextHandler(method, path string, handler ContextHandler) {
	r.ContextHandlerFunc(method, path, 
		func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
			handler.ServeHTTPContext(ctx, w, req)
		},
	)
}

// HTTPHandler is an adapter which allows the usage of an http.Handler as a
// request handle.
func (r *Router) HTTPHandler(method, path string, handler http.Handler) {
	r.ContextHandlerFunc(method, path,
		func(_ context.Context, w http.ResponseWriter, req *http.Request) {
			handler.ServeHTTP(w, req)
		},
	)
}

// HTTPHandlerFunc is an adapter which allows the usage of an http.HandlerFunc as a
// request handle.
func (r *Router) HTTPHandlerFunc(method, path string, handler http.HandlerFunc) {
	r.HTTPHandler(method, path, handler)
}

// ServeFiles serves files from the given file system root.
// The path must end with "/*filepath", files are then served from the local
// path /defined/root/dir/*filepath.
// For example if root is "/etc" and *filepath is "passwd", the local file
// "/etc/passwd" would be served.
// Internally a http.FileServer is used, therefore http.NotFound is used instead
// of the Router's NotFound handler.
// To use the operating system's file system implementation,
// use http.Dir:
//     router.ServeFiles("/src/*filepath", http.Dir("/var/www"))
func (r *Router) ServeFiles(path string, root http.FileSystem) {
	if len(path) < 10 || path[len(path)-10:] != "/*filepath" {
		panic("path must end with /*filepath in path '" + path + "'")
	}

	fileServer := http.FileServer(root)

	r.GET(path, func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		req.URL.Path = PathParam(ctx, "filepath")
		fileServer.ServeHTTP(w, req)
	})
}

func (r *Router) recv(w http.ResponseWriter, req *http.Request) {
	if rcv := recover(); rcv != nil {
		r.PanicHandler(ContextWithPanicException(rcv), w, req)
	}
}

// Lookup allows the manual lookup of a method + path combo.
// This is e.g. useful to build a framework around this router.
// If the path was found, it returns the handle function and the path parameter
// values. Otherwise the third return value indicates whether a redirection to
// the same path with an extra / without the trailing slash should be performed.
func (r *Router) Lookup(method, path string) (ContextHandlerFunc, context.Context, bool) {
	if root := r.trees[method]; root != nil {
		handle, ps, tsr := root.getValue(path)
		return handle, ContextWithPathParams(ps), tsr
	}
	return nil, nil, false
}

// ServeHTTP makes the router implement the http.Handler interface.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if r.PanicHandler != nil {
		defer r.recv(w, req)
	}

	if root := r.trees[req.Method]; root != nil {
		path := req.URL.Path

		if handle, ps, tsr := root.getValue(path); handle != nil {
			handle(ContextWithPathParams(ps), w, req)
			return
		} else if req.Method != "CONNECT" && path != "/" {
			code := 301 // Permanent redirect, request with GET method
			if req.Method != "GET" {
				// Temporary redirect, request with same method
				// As of Go 1.3, Go does not support status code 308.
				code = 307
			}

			if tsr && r.RedirectTrailingSlash {
				if len(path) > 1 && path[len(path)-1] == '/' {
					req.URL.Path = path[:len(path)-1]
				} else {
					req.URL.Path = path + "/"
				}
				http.Redirect(w, req, req.URL.String(), code)
				return
			}

			// Try to fix the request path
			if r.RedirectFixedPath {
				fixedPath, found := root.findCaseInsensitivePath(
					CleanPath(path),
					r.RedirectTrailingSlash,
				)
				if found {
					req.URL.Path = string(fixedPath)
					http.Redirect(w, req, req.URL.String(), code)
					return
				}
			}
		}
	}

	// Handle 405
	if r.HandleMethodNotAllowed {
		for method := range r.trees {
			// Skip the requested method - we already tried this one
			if method == req.Method {
				continue
			}

			handle, _, _ := r.trees[method].getValue(req.URL.Path)
			if handle != nil {
				if r.MethodNotAllowed != nil {
					r.MethodNotAllowed.ServeHTTP(w, req)
				} else {
					http.Error(w,
						http.StatusText(http.StatusMethodNotAllowed),
						http.StatusMethodNotAllowed,
					)
				}
				return
			}
		}
	}

	// Handle 404
	if r.NotFound != nil {
		r.NotFound.ServeHTTP(w, req)
	} else {
		http.NotFound(w, req)
	}
}
