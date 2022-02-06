package golivewire

import (
	"context"
	"github.com/tanapoln/golivewire/lib/mapstructure"
)

var (
	hooksRegistry = map[EventName][]LifecycleHook{}
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
	HookRegister(EventComponentHydrateInitial, LifecycleHookFunc(HookQueryParamHydration))
}

func HookRegister(name EventName, fn LifecycleHook) {
	hooksRegistry[name] = append(hooksRegistry[name], fn)
}

func hookDispatch(name EventName, lm *lifecycleManager) error {
	comp := lm.component
	base := comp.getBaseComponent()
	for _, fn := range hooksRegistry[name] {
		if err := fn.Execute(base.ctx, comp, &lm.response); err != nil {
			return err
		}
	}
	return nil
}

func HookQueryParamHydration(ctx context.Context, component Component, response *Response) error {
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
