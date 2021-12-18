package component

import (
	"fmt"
	"strings"
	"time"

	"github.com/tanapoln/golivewire"
)

func init() {
	golivewire.RegisterFactory(func() interface{} {
		return &OutputComponent{
			BaseComponent: golivewire.BaseComponent{
				Name: "output",
			},
		}
	})
}

type OutputComponent struct {
	golivewire.BaseComponent
	Output string `json:"output"`
}

func (o *OutputComponent) SetOutputTo(str ...string) error {
	o.Output = strings.Join(str, "")
	return nil
}

func (o *OutputComponent) SetOutputToFoo() error {
	o.Output = "foo"
	return nil
}

func (o *OutputComponent) AppendToOutput(str ...string) error {
	o.Output += strings.Join(str, "")
	return nil
}

func (o *OutputComponent) Pause() {
	time.Sleep(1 * time.Second)
}

func (o *OutputComponent) Render() (string, error) {
	template := `
	<div>
		<button type="button" wire:click="SetOutputToFoo" dusk="foo">Foo</button>
		<button type="button" wire:click="SetOutputTo('bar', 'bell')" dusk="bar">Bar</button>
		<button type="button" wire:click="SetOutputTo('a', &quot;b&quot; , 'c','d' ,'e', ''.concat('f'))" dusk="ball">Ball</button>
		<button type="button" wire:click="SetOutputToFoo()" dusk="bowl">Bowl</button>
		<button type="button" wire:click="@if (1) SetOutputToFoo() @else SetOutputToFoo() @endif" dusk="baw">Baw</button>
		<button type="button" wire:click="SetOutputTo('baz')" dusk="baz.outer"><button type="button" wire:click="$refresh" dusk="baz.inner">Inner</button> Outer</button>
		<input type="text" wire:blur="AppendToOutput('bop')" dusk="bop.input"><button type="button" wire:mousedown="AppendToOutput('bop')" dusk="bop.button">Blur &</button>
		<input type="text" wire:keydown="AppendToOutput('bob')" wire:keydown.enter="AppendToOutput('bob')" dusk="bob">
		<input type="text" wire:keydown.enter="SetOutputTo('lob')" dusk="lob">
		<input type="text" wire:keydown.shift.enter="SetOutputTo('law')" dusk="law">
		<input type="text" wire:keydown.space="SetOutputTo('spa')" dusk="spa">
		<form wire:submit.prevent="Pause">
			<div wire:ignore>
				<input type="text" dusk="blog.input.ignored">
			</div>
	
			<input type="text" dusk="blog.input">
			<button type="submit" dusk="blog.button">Submit</button>
		</form>
		<form wire:submit.prevent="ThrowError">
			<button type="submit" dusk="boo.button">Submit</button>
		</form>
		<input wire:keydown.debounce.75ms="SetOutputTo('bap')" dusk="bap"></button>
		<span dusk="output">%s</span>
	</div>
	`
	return fmt.Sprintf(template, o.Output), nil
}
