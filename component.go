package golivewire

import (
	"context"
	"github.com/rs/xid"
)

type BaseComponent struct {
	id               string
	name             string
	listeners        []string
	ctx              context.Context
	component        interface{}
	manager          *livewireManager
	preRenderView    *htmlView
	shouldSkipRender bool
}

func (c *BaseComponent) ID() string {
	if c.id == "" {
		c.id = xid.New().String()
	}
	return c.id
}

func (c *BaseComponent) Name() string {
	return c.name
}

func (c *BaseComponent) WithListeners(listeners ...string) {
	c.listeners = listeners
}

func (c *BaseComponent) getBaseComponent() *BaseComponent {
	return c
}

func (c *BaseComponent) skipRender() {
	c.shouldSkipRender = true
}

func (c *BaseComponent) renderToView() (*htmlView, error) {
	if c.shouldSkipRender {
		return nil, nil
	}

	renderer, ok := c.component.(Renderer)
	if ok {
		raw, err := renderer.Render(c.ctx)
		if err != nil {
			return nil, err
		}
		view, err := newHTMLView(raw)
		if err != nil {
			return nil, err
		}
		c.preRenderView = view
		return view, nil
	}

	return nil, ErrNotRenderer
}

type Component interface {
	getBaseComponent() *BaseComponent
}
