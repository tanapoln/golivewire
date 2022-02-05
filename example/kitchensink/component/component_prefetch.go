package component

import (
	"context"

	"github.com/tanapoln/golivewire"
)

func init() {
	golivewire.RegisterFactory("tests.browser.prefetch.component", func() golivewire.Component {
		return &ComponentPrefetch{}
	})
}

var (
	prefetchCount int
)

type ComponentPrefetch struct {
	golivewire.BaseComponent
}

func (c *ComponentPrefetch) Render(ctx context.Context) (string, error) {
	prefetchCount++

	return RenderTemplate(c, `
	<div>
		<button wire:click.prefetch="$refresh" dusk="button">inc</button>	
		<span dusk="count">{{ .Count }}</span>
	</div>
	`)
}

func (c *ComponentPrefetch) Count() int {
	return prefetchCount
}
