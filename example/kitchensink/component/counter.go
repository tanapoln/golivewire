package component

import (
	"fmt"

	"github.com/tanapoln/golivewire"
)

var (
	count int
)

func init() {
	golivewire.RegisterFactory(func() interface{} {
		return &Counter{
			BaseComponent: golivewire.BaseComponent{
				Name:      "counter",
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

func (c *Counter) Render() (string, error) {
	template := `
	<div wire:poll>
		Count: %d 
		<button wire:click="Incr">incr</button>
	</div>
	`
	return fmt.Sprintf(template, c.Count), nil
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
