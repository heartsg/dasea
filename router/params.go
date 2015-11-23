//Package router is in charge of routing and middleware supports

package router

import (
	"golang.org/x/net/context"
)

type key int

const (
	PathParamsKey key = iota
	PanicKey
	MiddlewareParamsKey
)


// pathParam is a single URL parameter, consisting of a key and a value.
type pathParam struct {
	Key   string
	Value string
}

// pathParams is a Param-slice, as returned by the router.
// The slice is ordered, the first URL parameter is also the first slice value.
// It is therefore safe to read values by the index.
type pathParams []pathParam


// ByName returns the value of the first Param which key matches the given name.
// If no matching Param is found, an empty string is returned.
func (ps pathParams) byName(name string) string {
	for i := range ps {
		if ps[i].Key == name {
			return ps[i].Value
		}
	}
	return ""
}

// PathParam returns a request URL parameter, or a blank string if it doesn't exist.
// For example, with the path /v2/papers/:page
// use router.PathParam(ctx, "page") to access the :page variable.
func PathParam(ctx context.Context, name string) string {
	params, ok := ctx.Value(PathParamsKey).(pathParams)
	if !ok {
		return ""
	}
	return params.byName(name)
}

func PanicException(ctx context.Context) interface{} {
	return ctx.Value(PanicKey)
}

func ContextWithPathParams(params pathParams) context.Context {
	return context.WithValue(context.Background(), PathParamsKey, params)
}

func ContextWithPanicException(e interface{}) context.Context {
	return context.WithValue(context.Background(), PanicKey, e)
}

func mergePathParams(ctx context.Context, params pathParams) context.Context {
	current, _ := ctx.Value(PathParamsKey).(pathParams)
	current = append(current, params...)
	return context.WithValue(ctx, PathParamsKey, current)
}



type middlewareParams map[string]interface{}

func MiddlewareParam(ctx context.Context, key string) interface{} {
	params, ok := ctx.Value(MiddlewareParamsKey).(middlewareParams)
	if !ok {
		return nil
	}
	return params[key]
}
func SetMiddlewareParam(ctx context.Context, key string, value interface{}) context.Context {
	params, ok := ctx.Value(MiddlewareParamsKey).(middlewareParams)
	if !ok {
		params := make(middlewareParams)
		params[key] = value
		return context.WithValue(ctx, MiddlewareParamsKey, params)
	} else {
		params[key] = value
		return ctx
	}
}

