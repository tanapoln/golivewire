package component

import (
	"bytes"
	"encoding/json"
	"html/template"

	"github.com/tanapoln/golivewire"
)

func jsonEncode(obj interface{}) (string, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func RenderTemplate(obj interface{}, tmpl string) (string, error) {
	t := template.New("root")

	t = t.Funcs(map[string]interface{}{
		"livewire":    golivewire.LivewireTemplateFunc,
		"json_encode": jsonEncode,
		"json":        jsonEncode,
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
