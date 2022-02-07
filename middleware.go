package golivewire

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
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

func ajaxMiddleware(handle httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		r = WithRequestContext(r)
		handle(w, r, p)
	}
}

type AjaxMiddlewareFunc func(original httprouter.Handle) httprouter.Handle
