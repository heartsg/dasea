
// Logger is a middleware that logs the start and end of each request, along
// with some useful data about what was requested, what the response status was,
// and how long it took to return. When standard output is a TTY, Logger will
// print in color, otherwise it will print in black and white.

package middleware

import (
	"bytes"
	"log"
	"net/http"
	"time"
	"dasea/router"
	"golang.org/x/net/context"
)

func StartLogger(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	var buf bytes.Buffer
	buf.WriteString("Started ")
	cW(&buf, bMagenta, "%s ", r.Method)
	cW(&buf, nBlue, "%q ", r.URL.String())
	buf.WriteString("from ")
	buf.WriteString(r.RemoteAddr)
	startTime := time.Now()
	ctx = router.SetMiddlewareParam(ctx, "start_time", startTime)
	log.Print(buf.String())
	return ctx
}

func EndLogger(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	endTime := time.Now()
	var buf bytes.Buffer
	
	buf.WriteString("Returning ")
	buf.WriteString(" in ")
	
	startTime, ok := router.MiddlewareParam(ctx, "start_time").(time.Time)
	if (ok) {
		dt := endTime.Sub(startTime)
		if dt < 500*time.Millisecond {
			cW(&buf, nGreen, "%s", dt)
		} else if dt < 5*time.Second {
			cW(&buf, nYellow, "%s", dt)
		} else {
			cW(&buf, nRed, "%s", dt)
		}
	} else {
		buf.WriteString("unknown time")
	}

	log.Print(buf.String())
	return ctx
}
