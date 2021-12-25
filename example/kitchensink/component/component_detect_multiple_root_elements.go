package component

import (
	"fmt"

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

func (c *ComponentWithCommentAsFirstElement) Render() (string, error) {
	tmpl := `
	<!-- A comment here -->
	<div>Element</div>
	`
	return tmpl, nil
}

type ComponentWithMultipleRootElements struct {
	golivewire.BaseComponent
}

func (c *ComponentWithMultipleRootElements) Render() (string, error) {
	tmpl := `
	<div>Element 1</div>
	<div>Element 2</div>
	`
	return tmpl, nil
}

type ComponentWithNestedSingleRootElement struct {
	golivewire.BaseComponent
}

func (c *ComponentWithNestedSingleRootElement) Render() (string, error) {
	tmpl := `
	<div>
		Nested: %s
		<span>Dummy Element</span>
	</div>
	`

	t, err := golivewire.LivewireTemplateFunc("tests.browser.detectmultiplerootelements.componentwithsinglerootelement")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(tmpl, t), nil
}

type ComponentWithSingleRootElement struct {
	golivewire.BaseComponent
}

func (c *ComponentWithSingleRootElement) Render() (string, error) {
	tmpl := `
	<div>Only Element</div>
	`
	return tmpl, nil
}
