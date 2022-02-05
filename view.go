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
	err = v.init()
	if err != nil {
		return nil, err
	}

	return v, nil
}

type htmlView struct {
	doc       *html.Node
	firstNode *html.Node
	nodes     []*html.Node
}

func (v *htmlView) RenderSafe() (template.HTML, error) {
	buf := &bytes.Buffer{}

	for _, n := range v.nodes {
		err := html.Render(buf, n)
		if err != nil {
			return "", err
		}
	}

	return template.HTML(buf.String()), nil
}

func (v *htmlView) AddWireTag(key string, data string) {
	v.firstNode.Attr = append(v.firstNode.Attr, html.Attribute{
		Key: "wire:" + key,
		Val: data,
	})
}

func (v *htmlView) init() error {
	nodes, err := extractNodesFromDoc(v.doc)
	if err != nil {
		return err
	}

	v.nodes = nodes
	if len(nodes) > 0 {
		v.firstNode = nodes[0]
	}
	return nil
}

func extractNodesFromDoc(node *html.Node) ([]*html.Node, error) {
	if node.Type != html.DocumentNode {
		return nil, ErrInvalidHTMLContent
	}
	body := node.FirstChild.FirstChild.NextSibling

	if body.Type != html.ElementNode && body.Data != "body" {
		return nil, ErrInvalidHTMLContent
	}

	nodes := getAllHTMLChildNodeFirstMatch(body, func(n *html.Node) bool {
		return n.Type == html.ElementNode
	})
	return nodes, nil
}

func getAllHTMLChildNodeFirstMatch(node *html.Node, pred func(n *html.Node) bool) []*html.Node {
	nodes := []*html.Node{}
	n := node.FirstChild
	for ; n != nil; n = n.NextSibling {
		if pred(n) {
			nodes = append(nodes, n)
		}
	}
	return nodes
}
