package golivewire

import (
	"bytes"
	"golang.org/x/net/html"
	"html/template"
	"strings"
)

func newHTMLView(raw string) (*htmlView, error) {
	wrapped := "<html><body>" + raw + "</body></html>"
	doc, err := html.Parse(strings.NewReader(wrapped))
	if err != nil {
		return nil, err
	}

	v := &htmlView{
		doc: doc,
	}

	return v, nil
}

type htmlView struct {
	doc *html.Node
}

func (v *htmlView) RenderSafe() (template.HTML, error) {
	//TODO Sharp: fix this & optimize
	buf := &bytes.Buffer{}
	err := html.Render(buf, v.doc)
	if err != nil {
		return "", err
	}
	return template.HTML(buf.String()), nil
}
