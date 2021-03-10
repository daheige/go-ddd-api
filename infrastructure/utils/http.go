package utils

import (
	"context"
	"net/http"
)

// ContextGet get value from http ctx
func ContextGet(r *http.Request, key interface{}) interface{} {
	return r.Context().Value(key)
}

// ContextSet set value to http ctx
func ContextSet(r *http.Request, key, val interface{}) *http.Request {
	if val == nil {
		return r
	}

	return r.WithContext(context.WithValue(r.Context(), key, val))
}

// GetStringByCtx get string key from ctx
func GetStringByCtx(ctx context.Context, key string) string {
	val := ctx.Value(key)
	if val == nil {
		return ""
	}

	str, ok := val.(string)
	if !ok {
		return ""
	}

	return str
}
