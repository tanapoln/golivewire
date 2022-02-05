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

	component, err := manager.NewComponentInstance(componentName)
	if err != nil {
		return "", err
	}

	lifecycle := newLifecycleFromInitialComponent(component)
	if err := lifecycle.initialHydrate(); err != nil {
		return "", err
	}
	if err := lifecycle.month(); err != nil {
		return "", err
	}
	if err := lifecycle.renderToView(); err != nil {
		return "", err
	}
	if err := lifecycle.initialDehydrate(); err != nil {
		return "", err
	}
	if err := lifecycle.toInitialResponse(); err != nil {
		return "", err
	}

	return lifecycle.response.Effects.Html, nil
}
