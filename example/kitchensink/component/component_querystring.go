package component

import (
	"context"
	"net/url"

	"github.com/tanapoln/golivewire"
)

func init() {
	golivewire.RegisterFactory("tests.browser.querystring.component", func() golivewire.Component {
		return &ComponentQueryString{
			Foo: "bar",
			Bar: "baz",
			Bob: []string{"foo", "bar"},
			Qux: map[string]interface{}{
				"hyphen":    "quux-quuz",
				"comma":     "quux,quuz",
				"ampersand": "quux&quuz",
				"space":     "quux quuz",
				"array":     []string{"quux", "quuz"},
			},
			ShowNestedComponent: false,
		}
	})

	golivewire.RegisterFactory("tests.browser.querystring.nestedcomponent", func() golivewire.Component {
		return &NestedComponentQueryString{
			Baz: "bop",
		}
	})
}

type ComponentQueryString struct {
	golivewire.BaseComponent

	Foo string                 `json:"foo" query:"foo"`
	Bar string                 `json:"bar" query:"bar"`
	Bob []string               `json:"bob" query:"bob"`
	Qux map[string]interface{} `json:"qux" query:"qux"`

	ShowNestedComponent bool `json:"showNestedComponent"`
}

func (q *ComponentQueryString) Querystring() url.Values {
	v := url.Values{}
	v.Set("foo", q.Foo)
	v.Set("bar", q.Bar)
	for _, s := range q.Bob {
		v.Add("bob", s)
	}
	return v
}

func (q *ComponentQueryString) ModifyBob() {
	q.Bob = []string{"foo", "bar", "baz"}
}

func (c *ComponentQueryString) Render(ctx context.Context) (string, error) {
	return RenderTemplate(c, `
	<div>
		<span dusk="output">{{ .Foo }}</span>
		<span dusk="bar-output">{{ .Bar }}</span>
	
		<span dusk="qux.hyphen">{{ .Qux.hyphen }}</span>
		<span dusk="qux.comma">{{ .Qux.comma }}</span>
		<span dusk="qux.ampersand">{{ .Qux.ampersand }}</span>
		<span dusk="qux.space">{{ .Qux.space }}</span>
		<span dusk="qux.array">{{ json_encode .Qux.array }}</span>
	
		<input wire:model="foo" type="text" dusk="input">
		<input wire:model="bar" type="text" dusk="bar-input">
	
		<button wire:click="$set('showNestedComponent', true)" dusk="show-nested">Show Nested Component</button>
	
		<button wire:click="modifyBob" dusk="bob.modify">Modify Bob (Array Property)</button>
		<span dusk="bob.output">{{ json .Bob }}</span>
	
		{{if .ShowNestedComponent}}
			{{livewire "tests.browser.querystring.nestedcomponent" .}}
		{{end}}
	</div>
	`)
}

type NestedComponentQueryString struct {
	golivewire.BaseComponent

	Baz string `json:"baz" query:"baz"`
}

func (n *NestedComponentQueryString) Querystring() url.Values {
	v := url.Values{}
	v.Set("baz", n.Baz)
	return v
}

func (n *NestedComponentQueryString) Render(ctx context.Context) (string, error) {
	return RenderTemplate(n, `
	<div>
		<span dusk="baz-output">{{ .Baz }}</span>
	
		<input wire:model="baz" type="text" dusk="baz-input">
	</div>
	`)
}
