package golivewire

import (
	"errors"
	"html/template"
)

var (
	ErrComponentNotFound = errors.New("components is not found")
)

func LivewireTemplateFunc(componentName string) (template.HTML, error) {
	factory, ok := componentRegistry[componentName]
	if !ok {
		return "", ErrComponentNotFound
	}

	comp := factory.createInstance()
	raw, err := InitialRender(comp.(Renderer))
	if err != nil {
		return "", err
	}
	return template.HTML(raw), nil
}
