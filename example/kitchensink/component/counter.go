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
	golivewire.RegisterFactory("counter", func() golivewire.Component {
		c := &Counter{}
		c.WithListeners("test")
		return c
	})
}

type Counter struct {
	golivewire.BaseComponent

	Count int `json:"count" query:"count"`
}

func (c *Counter) Validate(ctx context.Context) error {
	if c.Count == 7 {
		return golivewire.ErrBadRequest.Message("error 400")
	}
	return nil
}

func (c *Counter) Render(ctx context.Context) (string, error) {
	time.Sleep(100 * time.Millisecond)
	return RenderTemplate(c, `
	<div>
		Count: {{.Count}}
		<button wire:click="Incr">incr</button>
	</div>`)
}

func (c *Counter) Refresh() error {
	return nil
}

func (c *Counter) Incr() error {
	c.Count++
	return nil
}
