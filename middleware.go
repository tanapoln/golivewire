package golivewire

import (
	"net/http"
)

type httpRequestContext struct{}

func LivewireMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, WithRequestContext(r))
	})
}

func WithRequestContext(r *http.Request) *http.Request {
	ctx, _ := newManagerCtx(r.Context(), r)
	return r.WithContext(ctx)
}
