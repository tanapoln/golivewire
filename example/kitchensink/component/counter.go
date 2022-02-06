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
		}
	})
}

type Counter struct {
	golivewire.BaseComponent

	Banned bool `json:"banned"`
	Count  int  `json:"count" query:"count"`
}

func (c *Counter) Render(ctx context.Context) (string, error) {
	time.Sleep(100 * time.Millisecond)
	return RenderTemplate(c, `
	<div>
		Count: {{.Count}}
		{{if not .Banned}}
		<button wire:click="Incr">incr</button>
		{{else}}
		<button>Banned, disabled</button>
		{{end}}
	</div>`)
}

func (c *Counter) Incr() error {
	c.Count++
	c.Banned = true
	return nil
}
