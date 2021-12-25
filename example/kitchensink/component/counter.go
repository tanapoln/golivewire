package component

import (
	"context"
	"time"

	"github.com/tanapoln/golivewire"
)

var (
	count int
)

func init() {
	golivewire.RegisterFactory("counter", func() interface{} {
		return &Counter{
			BaseComponent: golivewire.BaseComponent{
				Listeners: []string{"test"},
			},
			Count: count,
		}
	})
}

type Counter struct {
	golivewire.BaseComponent

	Count int `json:"count"`
}

func (c *Counter) Render(ctx context.Context) (string, error) {
	time.Sleep(900 * time.Millisecond)
	return RenderTemplate(c, `
	<div wire:poll>
		Count: {{.Count}}
		<button wire:click="Incr">incr</button>
	</div>`)
}

func (c *Counter) Refresh() error {
	c.Count = count
	return nil
}

func (c *Counter) Incr() error {
	count++
	c.Count = count
	return nil
}
