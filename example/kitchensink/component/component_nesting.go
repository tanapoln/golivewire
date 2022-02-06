package component

import (
	"context"
	"net/url"
	"strconv"

	"github.com/tanapoln/golivewire"
)

func init() {
	golivewire.RegisterFactory("tests.browser.nesting.component", func() golivewire.Component {
		return &ComponentNestingComponent{}
	})

	golivewire.RegisterFactory("tests.browser.nesting.nestedcomponent", func() golivewire.Component {
		return &ComponentNestingNestedComponent{}
	})

	golivewire.RegisterFactory("tests.browser.nesting.rendercontextcomponent", func() golivewire.Component {
		return &ComponentNestingRenderContextComponent{}
	})
}

type ComponentNestingComponent struct {
	golivewire.BaseComponent

	ShowChild bool   `json:"showChild" query:"showChild"`
	Key       string `json:"key"`
}

func (c *ComponentNestingComponent) Querystring() url.Values {
	v := url.Values{}
	v.Set("showChild", strconv.FormatBool(c.ShowChild))
	return v
}

func (c *ComponentNestingComponent) Render(ctx context.Context) (string, error) {
	return RenderTemplate(c, `
<div>
    <button wire:click="$toggle('showChild')" dusk="button.toggleChild"></button>

    <button wire:click="$set('key', 'bar')" dusk="button.changeKey"></button>

    {{if .ShowChild}}
        {{livewire "tests.browser.nesting.nestedcomponent" .Key .}}
    {{end}}
</div>
`)
}

type ComponentNestingNestedComponent struct {
	golivewire.BaseComponent

	Output string `json:"output"`
}

func (c *ComponentNestingNestedComponent) Render(ctx context.Context) (string, error) {
	return RenderTemplate(c, `
<div>
    <button wire:click="$set('output', 'foo')" dusk="button.nested"></button>

    <span dusk="output.nested">{{ .Output }}</span>
</div>
`)
}

type ComponentNestingRenderContextComponent struct {
	golivewire.BaseComponent

	One   string `json:"one"`
	Two   string `json:"two"`
	Three string `json:"three"`
}

func (c *ComponentNestingRenderContextComponent) Render(ctx context.Context) (string, error) {
	return RenderTemplate(c, `
<div>
    <x-blade-component dusk="output.blade-component1" property="one" />
    <x-blade-component dusk="output.blade-component2" property="two" />

    <div>
        {{livewire "tests.browser.nesting.nestedcomponent" .NestedArg .}}
    </div>

    <x-blade-component dusk="output.blade-component3" property="three" />
</div>
`)
}

func (c *ComponentNestingRenderContextComponent) NestedArg() map[string]interface{} {
	return map[string]interface{}{
		"output": "Sub render",
	}
}
