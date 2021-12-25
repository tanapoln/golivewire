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

	ctx := context.Background()
	for _, arg := range args {
		if v, ok := arg.(context.Context); ok {
			ctx = v
			break
		}
		if v, ok := arg.(baseComponentSupport); ok {
			ctx = v.getBaseComponent().getContext()
			break
		}
	}

	factory, ok := componentRegistry[componentName]
	if !ok {
		return "", ErrComponentNotFound
	}

	comp, err := factory.createInstance(ctx)
	if err != nil {
		return "", err
	}
	
	raw, err := InitialRender(ctx, comp.(Renderer))
	if err != nil {
		return "", err
	}
	return template.HTML(raw), nil
}
