package component

import (
	"context"

	"github.com/tanapoln/golivewire"
)

func init() {
	golivewire.RegisterFactory("tests.browser.detectmultiplerootelements.componentwithcommentasfirstelement", func() interface{} {
		return &ComponentWithCommentAsFirstElement{}
	})

	golivewire.RegisterFactory("tests.browser.detectmultiplerootelements.componentwithmultiplerootelements", func() interface{} {
		return &ComponentWithMultipleRootElements{}
	})

	golivewire.RegisterFactory("tests.browser.detectmultiplerootelements.componentwithnestedsinglerootelement", func() interface{} {
		return &ComponentWithNestedSingleRootElement{}
	})

	golivewire.RegisterFactory("tests.browser.detectmultiplerootelements.componentwithsinglerootelement", func() interface{} {
		return &ComponentWithSingleRootElement{}
	})
}

type ComponentWithCommentAsFirstElement struct {
	golivewire.BaseComponent
}

func (c *ComponentWithCommentAsFirstElement) Render(ctx context.Context) (string, error) {
	tmpl := `
	<!-- A comment here -->
	<div>Element</div>
	`
	return tmpl, nil
}

type ComponentWithMultipleRootElements struct {
	golivewire.BaseComponent
}

func (c *ComponentWithMultipleRootElements) Render(ctx context.Context) (string, error) {
	tmpl := `
	<div>Element 1</div>
	<div>Element 2</div>
	`
	return tmpl, nil
}

type ComponentWithNestedSingleRootElement struct {
	golivewire.BaseComponent
}

func (c *ComponentWithNestedSingleRootElement) Render(ctx context.Context) (string, error) {
	return RenderTemplate(c, `
	<div>
		Nested: {{livewire "tests.browser.detectmultiplerootelements.componentwithsinglerootelement"}}
		<span>Dummy Element</span>
	</div>
	`)
}

type ComponentWithSingleRootElement struct {
	golivewire.BaseComponent
}

func (c *ComponentWithSingleRootElement) Render(ctx context.Context) (string, error) {
	tmpl := `
	<div>Only Element</div>
	`
	return tmpl, nil
}
