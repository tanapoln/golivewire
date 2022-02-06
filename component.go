package golivewire

import (
	"context"
)

type BaseComponent struct {
	id               string
	name             string
	listeners        []string
	ctx              context.Context
	component        interface{}
	preRenderView    *htmlView
	shouldSkipRender bool
	children         []Component
}

func (c *BaseComponent) ID() string {
	return c.id
}

func (c *BaseComponent) Name() string {
	return c.name
}

func (c *BaseComponent) WithListeners(listeners ...string) {
	c.listeners = listeners
}

func (c *BaseComponent) Context() context.Context {
	return c.ctx
}

func (c *BaseComponent) getBaseComponent() *BaseComponent {
	return c
}

func (c *BaseComponent) skipRender() {
	c.shouldSkipRender = true
}

func (c *BaseComponent) addChild(comp Component) {
	c.children = append(c.children, comp)
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
