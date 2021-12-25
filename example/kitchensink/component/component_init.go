package component

import (
	"context"
	"fmt"

	"github.com/tanapoln/golivewire"
)

func init() {
	golivewire.RegisterFactory("tests.browser.init.component", func() interface{} {
		return &ComponentInit{}
	})
}

type ComponentInit struct {
	golivewire.BaseComponent

	Output string `json:"output"`
}

func (o *ComponentInit) SetOutputToFoo() {
	o.Output = "foo"
}

func (c *ComponentInit) Render(ctx context.Context) (string, error) {
	tmpl := `
		<div wire:init="setOutputToFoo">
			<span dusk="output">%s</span>
		</div>
	`
	return fmt.Sprintf(tmpl, c.Output), nil
}
