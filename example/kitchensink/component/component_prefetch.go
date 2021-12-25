package component

import (
	"github.com/tanapoln/golivewire"
)

func init() {
	golivewire.RegisterFactory("tests.browser.prefetch.component", func() interface{} {
		return &ComponentPrefetch{}
	})
}

var (
	prefetchCount int
)

type ComponentPrefetch struct {
	golivewire.BaseComponent
}

func (c *ComponentPrefetch) Render() (string, error) {
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
