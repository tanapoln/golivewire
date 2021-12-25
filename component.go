package golivewire

import (
	"context"

	"github.com/rs/xid"
)

type BaseComponent struct {
	id        string
	Name      string   `json:"-"`
	Listeners []string `json:"-"`
	ctx       context.Context
}

func (b *BaseComponent) GetID() string {
	if b.id == "" {
		b.id = xid.New().String()
	}
	return b.id
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
