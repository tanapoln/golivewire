package golivewire

import (
	"context"

	"github.com/rs/xid"
)

type BaseComponent struct {
	id        string
	name      string
	listeners []string
	ctx       context.Context
}

func (b *BaseComponent) ID() string {
	if b.id == "" {
		b.id = xid.New().String()
	}
	return b.id
}

func (c *BaseComponent) Name() string {
	return c.name
}

func (b *BaseComponent) WithListeners(listeners ...string) {
	b.listeners = listeners
}

func (b *BaseComponent) getContext() context.Context {
	if b.ctx != nil {
		return b.ctx
	} else {
		return context.Background()
	}
}

func (b *BaseComponent) getBaseComponent() *BaseComponent {
	return b
}

type baseComponentSupport interface {
	getBaseComponent() *BaseComponent
}
