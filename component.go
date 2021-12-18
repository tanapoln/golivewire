package golivewire

import "github.com/rs/xid"

type BaseComponent struct {
	id        string
	Name      string   `json:"-"`
	Listeners []string `json:"-"`
}

func (b *BaseComponent) GetID() string {
	if b.id == "" {
		b.id = xid.New().String()
	}
	return b.id
}

func (b *BaseComponent) getBaseComponent() *BaseComponent {
	return b
}

type baseComponentSupport interface {
	getBaseComponent() *BaseComponent
}
