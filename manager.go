package golivewire

import (
	"context"
	"net/http"
)

type managerCtxKey struct{}

func newManagerCtx(ctx context.Context, req *http.Request) (context.Context, *livewireManager) {
	mgr := &livewireManager{}
	newctx := context.WithValue(ctx, managerCtxKey{}, mgr)
	mgr.ctx = newctx
	mgr.httpReq = req

	return newctx, mgr
}

func managerFromCtx(ctx context.Context) *livewireManager {
	v := ctx.Value(managerCtxKey{})
	if v == nil {
		return nil
	}
	return v.(*livewireManager)
}

type livewireManager struct {
	httpReq *http.Request
	ctx     context.Context
	req     Request
}

func (l *livewireManager) ProcessRequest() error {
	return bindJSONRequest(l.httpReq, &l.req)
}

func (l *livewireManager) OriginalPath() string {
	return l.httpReq.URL.Path
}

func (l *livewireManager) OriginalMethod() string {
	return l.httpReq.Method
}

func (l *livewireManager) NewComponentInstance(name string) (Component, error) {
	compFactory, ok := componentRegistry[name]
	if !ok {
		return nil, ErrNotFound.Message("component name is not found: " + name)
	}
	return compFactory.createInstance(l.ctx)
}

func (l *livewireManager) GetComponentInstance(name string, id string) (Component, error) {
	compFactory, ok := componentRegistry[name]
	if !ok {
		return nil, ErrNotFound.Message("component name is not found: " + name)
	}
	return compFactory.createInstanceWithID(l.ctx, id)
}

func (l *livewireManager) IsLivewireRequest() bool {
	return l.httpReq.Header.Get("x-livewire") != ""
}
