package golivewire

import (
	"context"
	"errors"
	"html/template"
)

var (
	ErrComponentNotFound = errors.New("components is not found")
)

func LivewireTemplateFunc(args ...interface{}) (template.HTML, error) {
	if len(args) == 0 {
		return "", errors.New("missing component name for livewire template func")
	}
	var componentName string
	if name, ok := args[0].(string); !ok {
		return "", errors.New("livewire expect component name to be string")
	} else {
		componentName = name
	}

	var ctx context.Context
	for _, arg := range args {
		if v, ok := arg.(context.Context); ok {
			ctx = v
			break
		}
		if v, ok := arg.(Component); ok {
			ctx = v.getBaseComponent().ctx
			break
		}
	}
	if ctx == nil {
		return "", errors.New("no context or component on render")
	}

	manager := managerFromCtx(ctx)
	if manager == nil {
		return "", errors.New("invalid context, no golivewire manager found")
	}

	lifecycle, err := newLifecycleFromInitialComponent(manager, componentName)
	if err != nil {
		return "", err
	}
	if err := lifecycle.InitialHydrate(); err != nil {
		return "", err
	}
	if err := lifecycle.Mount(); err != nil {
		return "", err
	}
	if err := lifecycle.RenderToView(); err != nil {
		return "", err
	}
	if err := lifecycle.InitialDehydrate(); err != nil {
		return "", err
	}
	if err := lifecycle.ToInitialResponse(); err != nil {
		return "", err
	}

	return lifecycle.response.Effects.Html, nil
}
