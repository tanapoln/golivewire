package golivewire

import (
	"context"
	"fmt"
	"net/url"

	"github.com/tanapoln/golivewire/lib/mapstructure"
)

type EventName int

const (
	EventUnknown EventName = iota
	EventComponentRendering
	EventComponentRendered
	EventViewRendered
	EventComponentBoot
	EventFailedValidation
	EventComponentMount
	EventComponentBooted
	EventMounted
	EventComponentUpdating
	EventComponentUpdated
	EventActionReturned
	EventComponentHydrate
	EventComponentHydrateInitial
	EventComponentHydrateSubsequent
	EventComponentDehydrate
	EventComponentDehydrateInitial
	EventComponentDehydrateSubsequent
	EventPropertyHydrate
	EventPropertyDehydrate
)

func init() {
	registerBooter(func(manager *livewireManager) {
		manager.HookRegister(EventComponentHydrateInitial, LifecycleHookFunc(hookQueryParamHydration))
	})

	registerBooter(func(manager *livewireManager) {
		hookUrlQuerySupport := newHookUrlQuerySupport()
		manager.HookRegister(EventComponentDehydrateInitial, LifecycleHookFunc(hookUrlQuerySupport.dehydrateInitial))
		manager.HookRegister(EventComponentDehydrateSubsequent, LifecycleHookFunc(hookUrlQuerySupport.dehydrateSubsequent))
	})

	registerBooter(func(manager *livewireManager) {
		manager.HookRegister(EventComponentDehydrate, LifecycleHookFunc(func(ctx context.Context, component Component, request *Request, response *Response) error {
			for _, child := range component.getBaseComponent().children {
				response.ServerMemo.Children = append(response.ServerMemo.Children, childComponent{
					ID:  child.getBaseComponent().ID(),
					Tag: child.getBaseComponent().preRenderView.firstNode.Data,
				})
			}

			return nil
		}))

		manager.HookRegister(EventComponentBooted, LifecycleHookFunc(func(ctx context.Context, component Component, request *Request, response *Response) error {
			if c, ok := component.(ComponentValidator); ok {
				if err := c.Validate(ctx); err != nil {
					return err
				} else {
					return nil
				}
			} else {
				return nil
			}
		}))
	})
}

func hookQueryParamHydration(ctx context.Context, component Component, request *Request, response *Response) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName:          "query",
		WeaklyTypedInput: true,
		Result:           component,
	})
	if err != nil {
		return err
	}

	manager := managerFromCtx(ctx)
	query := manager.Queryparams()
	if err := decoder.Decode(query); err != nil {
		return err
	}
	return nil
}

func newHookUrlQuerySupport() *hookUrlQuerySupport {
	return &hookUrlQuerySupport{
		qs: url.Values{},
	}
}

type hookUrlQuerySupport struct {
	qs url.Values
}

func (h *hookUrlQuerySupport) replaceQuery(u url.Values) {
	for key, vals := range u {
		h.qs[key] = vals
	}
}

func (h *hookUrlQuerySupport) dehydrateInitial(ctx context.Context, component Component, request *Request, response *Response) error {
	if q, ok := component.(Querystringer); ok {
		fmt.Printf("[DEBUG] initial dehydrate component:%v\n", component.getBaseComponent().Name())

		manager := managerFromCtx(ctx)
		if manager.IsLivewireRequest() {
			return h.dehydrateSubsequent(ctx, component, request, response)
		}

		var existingURL *url.URL
		if response.Effects.Path != "" {
			u, err := url.Parse(response.Effects.Path)
			if err != nil {
				return err
			}
			existingURL = u
		} else {
			u, err := url.Parse(request.Fingerprint.Path)
			if err != nil {
				return err
			}
			existingURL = u
		}

		h.replaceQuery(existingURL.Query())
		h.replaceQuery(q.Querystring())
		existingURL.RawQuery = h.qs.Encode()
		response.Effects.Path = manager.OriginalBaseURL() + existingURL.String()
	}

	return nil
}

func (h *hookUrlQuerySupport) dehydrateSubsequent(ctx context.Context, component Component, request *Request, response *Response) error {
	if q, ok := component.(Querystringer); ok {
		fmt.Printf("[DEBUG] subsequent dehydrate component:%v\n", component.getBaseComponent().Name())

		manager := managerFromCtx(ctx)
		var existingURL *url.URL
		if response.Effects.Path != "" {
			u, err := url.Parse(response.Effects.Path)
			if err != nil {
				return err
			}
			existingURL = u
		} else {
			referer := manager.httpReq.Header.Get("referer")
			u, err := url.Parse(referer)
			if err != nil {
				return err
			}
			existingURL = u
		}

		h.replaceQuery(existingURL.Query())
		h.replaceQuery(q.Querystring())
		existingURL.RawQuery = h.qs.Encode()
		response.Effects.Path = existingURL.String()
	}

	return nil
}
