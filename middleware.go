package golivewire

import (
	"context"
	"net/http"
)

type httpRequestContext struct{}

func LivewireMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, WithRequestContext(r))
	})
}

func WithRequestContext(r *http.Request) *http.Request {
	ctx := context.WithValue(r.Context(), httpRequestContext{}, r)
	newReq := r.WithContext(ctx)
	return newReq
}

func httpRequestFromContext(ctx context.Context) *http.Request {
	val := ctx.Value(httpRequestContext{})
	if val == nil {
		return nil
	}
	if req, ok := val.(*http.Request); ok {
		return req
	} else {
		return nil
	}
}