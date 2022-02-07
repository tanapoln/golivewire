package golivewire

import (
	"bytes"
	_ "embed"
)

import (
	"context"
	"errors"
	"html/template"
)

var (
	ErrComponentNotFound = errors.New("components is not found")
)

func LivewireTemplateFunc(args ...interface{}) (template.HTML, error) {
	if len(args) == 0 {
		return "", errors.New("missing component name for livewire template func")
	}
	var componentName string
	if name, ok := args[0].(string); !ok {
		return "", errors.New("livewire expect component name to be string")
	} else {
		componentName = name
	}

	var ctx context.Context
	var parentComponent Component
	var componentParams map[string]interface{}
	for _, arg := range args {
		if v, ok := arg.(context.Context); ok {
			ctx = v
		}
		if v, ok := arg.(Component); ok {
			parentComponent = v
			if ctx == nil { // prefer ctx from args, do not overwrite
				ctx = v.getBaseComponent().ctx
			}
		}
		if v, ok := arg.(map[string]interface{}); ok {
			componentParams = v
		}
	}
	if ctx == nil {
		return "", errors.New("no context or component on render")
	}
	if _, ok := ctx.Deadline(); !ok {
		return "", ctx.Err()
	}

	manager := managerFromCtx(ctx)
	if manager == nil {
		return "", errors.New("invalid context, no golivewire manager found")
	}

	lifecycle, err := newLifecycleFromInitialComponent(manager, componentName)
	if err != nil {
		return "", err
	}
	if parentComponent != nil {
		parentComponent.getBaseComponent().addChild(lifecycle.component)
	}

	if err := lifecycle.Boot(); err != nil {
		return "", err
	}
	if err := lifecycle.InitialHydrate(); err != nil {
		return "", err
	}
	if err := lifecycle.Mount(componentParams); err != nil {
		return "", err
	}
	if err := lifecycle.RenderToView(); err != nil {
		return "", err
	}
	if err := lifecycle.InitialDehydrate(); err != nil {
		return "", err
	}
	if err := lifecycle.ToInitialResponse(); err != nil {
		return "", err
	}

	return lifecycle.response.Effects.Html, nil
}

//go:embed static/livewire.init.js.tmpl
var jsInitRaw string

var (
	livewireInitTemplate *template.Template
)

func init() {
	livewireInitTemplate = template.Must(template.New("livewire.init").Parse(jsInitRaw))
}

func LivewireJS(csrfToken string) (template.HTML, error) {
	buf := &bytes.Buffer{}
	err := livewireInitTemplate.Execute(buf, H{
		"Token":       csrfToken,
		"BaseURL":     baseURL,
		"Development": DevelopmentMode,
	})
	if err != nil {
		return "", err
	}

	return template.HTML(buf.String()), nil
}

func TemplateFuncMap() template.FuncMap {
	return template.FuncMap{
		"livewire":   LivewireTemplateFunc,
		"livewireJS": LivewireJS,
	}
}
