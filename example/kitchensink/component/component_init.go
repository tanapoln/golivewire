package component

import (
	"fmt"

	"github.com/tanapoln/golivewire"
)

func init() {
	golivewire.RegisterFactory(func() interface{} {
		return &ComponentInit{
			BaseComponent: golivewire.BaseComponent{
				Name: "tests.browser.init.component",
			},
		}
	})
}

type ComponentInit struct {
	golivewire.BaseComponent

	Output string `json:"output"`
}

func (o *ComponentInit) SetOutputToFoo() {
	o.Output = "foo"
}

func (c *ComponentInit) Render() (string, error) {
	tmpl := `
		<div wire:init="setOutputToFoo">
			<span dusk="output">%s</span>
		</div>
	`
	return fmt.Sprintf(tmpl, c.Output), nil
}
