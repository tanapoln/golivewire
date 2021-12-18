package golivewire

import (
	"errors"
	"strings"
)

var (
	componentRegistry = factoryRegistry{}
	ErrNotComponent   = errors.New("object must embeded golivewire.BaseComponent")
	ErrNoNameDefined  = errors.New("component must have name defined, cannot be empty")
	ErrNotRenderer    = errors.New("object must implement golivewire.Renderer")
)

type factoryRegistry map[string]factory

func (r factoryRegistry) register(fn ComponentFactoryFunc) error {
	comp := fn()
	if _, ok := comp.(Renderer); !ok {
		return ErrNotRenderer
	}
	if _, ok := comp.(baseComponentSupport); !ok {
		return ErrNotComponent
	}

	name := comp.(baseComponentSupport).getBaseComponent().Name
	name = strings.TrimSpace(name)
	if name == "" {
		return ErrNoNameDefined
	}

	r[name] = factory{
		name: name,
		fn:   fn,
	}
	return nil
}

type factory struct {
	fn   ComponentFactoryFunc
	name string
}

func (f factory) createInstance() baseComponentSupport {
	comp := f.fn()
	return comp.(baseComponentSupport)
}

// RegisterFactory register component factory. It's not thread-safe.
func RegisterFactory(fn ComponentFactoryFunc) {
	err := componentRegistry.register(fn)
	if err != nil {
		panic(err)
	}
}
