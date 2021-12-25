package component

import (
	"bytes"
	"html/template"

	"github.com/tanapoln/golivewire"
)

func RenderTemplate(obj interface{}, tmpl string) (string, error) {
	t := template.New("root")
	t = t.Funcs(map[string]interface{}{
		"livewire": golivewire.LivewireTemplateFunc,
	})
	parse, err := t.Parse(tmpl)
	if err != nil {
		return "", err
	}
	buf := &bytes.Buffer{}
	err = parse.Execute(buf, obj)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
