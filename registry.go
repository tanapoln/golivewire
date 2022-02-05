package golivewire

import (
	"context"
	"errors"
	"github.com/rs/xid"
	"strings"
)

var (
	componentRegistry  = factoryRegistry{}
	ErrNotComponent    = errors.New("object must embeded golivewire.BaseComponent")
	ErrNoNameDefined   = errors.New("component must have name defined, cannot be empty")
	ErrNotRenderer     = errors.New("object must implement golivewire.Renderer")
	ErrCreateComponent = errors.New("cannot create component, component is not valid")
)

type factoryRegistry map[string]factory

func (r factoryRegistry) register(name string, fn ComponentFactoryFunc) error {
	comp := fn()
	if _, ok := comp.(Renderer); !ok {
		return ErrNotRenderer
	}

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

func (f factory) valid() bool {
	return f.fn != nil
}

func (f factory) createInstance(ctx context.Context) (Component, error) {
	if !f.valid() {
		return nil, ErrCreateComponent
	}

	comp := f.fn()

	manager := managerFromCtx(ctx)
	base := comp.getBaseComponent()
	base.id = xid.New().String()
	base.name = f.name
	base.ctx = ctx
	base.component = comp
	base.manager = manager
	return comp, nil
}

// RegisterFactory register component factory. It's not thread-safe.
func RegisterFactory(componentName string, fn ComponentFactoryFunc) {
	err := componentRegistry.register(componentName, fn)
	if err != nil {
		panic(err)
	}
}
