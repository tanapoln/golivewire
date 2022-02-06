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
	mgr.hooks = make(map[EventName][]LifecycleHook)

	mgr.boot()

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

	hooks map[EventName][]LifecycleHook
}

func (l *livewireManager) boot() {
	l.HookRegister(EventComponentHydrateInitial, LifecycleHookFunc(hookQueryParamHydration))

	hookUrlQuerySupport := newHookUrlQuerySupport()
	l.HookRegister(EventComponentDehydrateInitial, LifecycleHookFunc(hookUrlQuerySupport.dehydrateInitial))
	l.HookRegister(EventComponentDehydrateSubsequent, LifecycleHookFunc(hookUrlQuerySupport.dehydrateSubsequent))
}

func (l *livewireManager) HookRegister(name EventName, fn LifecycleHook) {
	l.hooks[name] = append(l.hooks[name], fn)
}

func (l *livewireManager) HookDispatch(name EventName, lm *lifecycleManager) error {
	comp := lm.component
	base := comp.getBaseComponent()
	for _, fn := range l.hooks[name] {
		if err := fn.Execute(base.ctx, comp, &lm.request, &lm.response); err != nil {
			return err
		}
	}
	return nil
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

func (l *livewireManager) OriginalBaseURL() string {
	return "http://localhost:8081"
}

func (l *livewireManager) Queryparams() map[string]interface{} {
	result := map[string]interface{}{}
	for key, val := range l.httpReq.URL.Query() {
		switch len(val) {
		case 0:
			result[key] = ""
		case 1:
			result[key] = val[0]
		default:
			result[key] = val
		}
	}
	return result
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
