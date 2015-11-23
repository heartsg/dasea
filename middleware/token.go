// To check whether token is valid (authentication)
//

package middleware

import (
	"net/http"
	"dasea/router"
	"golang.org/x/net/context"
)

func Token(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	return ctx
}
