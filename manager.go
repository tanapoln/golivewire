package golivewire

import (
	"context"
	"net/http"
)

type managerCtxKey struct{}

func newManagerCtx(ctx context.Context, req *http.Request) context.Context {
	mgr := &livewireManager{}
	newctx := context.WithValue(ctx, managerCtxKey{}, mgr)
	mgr.ctx = newctx
	mgr.req = req

	return newctx
}

func managerFromCtx(ctx context.Context) *livewireManager {
	return ctx.Value(managerCtxKey{}).(*livewireManager)
}

type livewireManager struct {
	req *http.Request
	ctx context.Context
}

func (l *livewireManager) OriginalPath() string {
	return l.req.URL.Path
}

func (l *livewireManager) OriginalMethod() string {
	return l.req.Method
}
