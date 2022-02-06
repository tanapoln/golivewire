package component

import (
	"context"

	"github.com/tanapoln/golivewire"
)

func init() {
	golivewire.RegisterFactory("tests.browser.dirty.component", func() golivewire.Component {
		return &ComponentDirty{}
	})
}

type ComponentDirty struct {
	golivewire.BaseComponent

	Foo string `json:"foo"`
	Bar string `json:"bar"`
	Baz string `json:"baz"`
	Bob string `json:"bob"`
}

func (c *ComponentDirty) Render(ctx context.Context) (string, error) {
	return RenderTemplate(c, `
<div>
    <input wire:model.lazy="foo" wire:dirty.class="foo-dirty" dusk="foo">
    <input wire:model.lazy="bar" wire:dirty.class.remove="bar-dirty" class="bar-dirty" dusk="bar">
    <span wire:dirty.class="baz-dirty" wire:target="baz" dusk="baz.target"><input wire:model.lazy="baz" dusk="baz.input"></span>
    <span wire:dirty wire:target="bob" dusk="bob.target">Dirty Indicator</span><input wire:model.lazy="bob" dusk="bob.input">

    <button type="button" dusk="dummy"></button>
</div>
`)
}
