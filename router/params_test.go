package router

import (
	"testing"
	"golang.org/x/net/context"	
)

func TestMiddlewareParams(t *testing.T) {
	ctx := context.Background()
	ctx = SetMiddlewareParam(ctx, "key1", "value1")
	SetMiddlewareParam(ctx, "key2", "value2")
	value1 := MiddlewareParam(ctx, "key1").(string)
	value2 := MiddlewareParam(ctx, "key2").(string)
	
	if (value1 != "value1") {
		t.Fatalf("middleware params: want value1, got %s", value1)
	}
	if (value2 != "value2") {
		t.Fatalf("middleware params: want value2, got %s", value2)
	}
}