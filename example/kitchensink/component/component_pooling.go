package component

import (
	"context"

	"github.com/tanapoln/golivewire"
)

func init() {
	golivewire.RegisterFactory("tests.browser.polling.component", func() interface{} {
		return &ComponentPooling{}
	})
}

type ComponentPooling struct {
	golivewire.BaseComponent

	Enabled bool
	Count   int
}

func (c *ComponentPooling) Render(ctx context.Context) (string, error) {
	c.Count++

	return RenderTemplate(c, `
	<div {{if .Enabled}} wire:poll.500ms {{end}}>
		<button wire:click="$refresh" dusk="refresh">count++</button>
		<button wire:click="$set('enabled', true)" dusk="enable">enable</button>
		<button wire:click="$set('enabled', false)" dusk="disable">disable</button>
	
		<span dusk="output">{{.Count}}</span>
	</div>`)
}
