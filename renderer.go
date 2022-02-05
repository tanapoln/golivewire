package golivewire

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

type HTMLDecorator interface {
	Decorate(ctx context.Context, node *html.Node, component interface{}) error
}

var (
	rendererPipeline        []HTMLDecorator
	initialRendererPipeline []HTMLDecorator
	ErrInvalidHTMLContent   = errors.New("invalid HTML content")
	bufPool                 = &sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
		},
	}
)

func init() {
	rendererPipeline = append(rendererPipeline, htmlDecoratorFunc(livewireIdRenderer))

	initialRendererPipeline = append(initialRendererPipeline, rendererPipeline...)
	initialRendererPipeline = append(initialRendererPipeline, htmlDecoratorFunc(livewireInitialDataRenderer))
}

func AddDecorator(decorator HTMLDecorator) {
	rendererPipeline = append(rendererPipeline, decorator)
}

type htmlDecoratorFunc func(ctx context.Context, node *html.Node, component interface{}) error

func (h htmlDecoratorFunc) Decorate(ctx context.Context, node *html.Node, component interface{}) error {
	return h(ctx, node, component)
}

func livewireIdRenderer(ctx context.Context, node *html.Node, component interface{}) error {
	var baseComp *BaseComponent
	if v, ok := component.(Component); !ok {
		return ErrNotComponent
	} else {
		baseComp = v.getBaseComponent()
	}

	node.Attr = append(node.Attr, html.Attribute{
		Key: "wire:id",
		Val: baseComp.ID(),
	})

	return nil
}

func livewireInitialDataRenderer(ctx context.Context, node *html.Node, component interface{}) error {
	var baseComp *BaseComponent
	if v, ok := component.(Component); !ok {
		return ErrNotComponent
	} else {
		baseComp = v.getBaseComponent()
	}

	initData := Response{
		Fingerprint: fingerprint{
			ID:   baseComp.ID(),
			Name: baseComp.name,
		},
		Effects: Effects{
			Listeners: baseComp.listeners,
		},
		ServerMemo: serverMemo{
			Data: component,
		},
	}

	data, err := json.Marshal(initData)
	if err != nil {
		return err
	}
	node.Attr = append(node.Attr, html.Attribute{
		Key: "wire:initial-data",
		Val: string(data),
	})
	return nil
}

func InitialRender(ctx context.Context, obj interface{}) (string, error) {
	return renderWithDecorators(ctx, obj, initialRendererPipeline...)
}

func renderWithDecorators(ctx context.Context, obj interface{}, decorators ...HTMLDecorator) (string, error) {
	if _, ok := obj.(Renderer); !ok {
		return "", ErrNotRenderer
	}
	if _, ok := obj.(Component); !ok {
		return "", ErrNotComponent
	}

	raw, err := obj.(Renderer).Render(ctx)
	if err != nil {
		return "", err
	}
	if err := ctx.Err(); err != nil {
		return "", err
	}

	wrappedHTML := fmt.Sprintf("<html><body>%s</body></html>", raw)
	doc, err := html.Parse(strings.NewReader(wrappedHTML))
	if err != nil {
		return "", err
	}

	nodes, err := extractNodesFromDoc(doc)
	if err != nil {
		return "", err
	}

	for _, decorator := range decorators {
		err := decorator.Decorate(ctx, nodes[0], obj)
		if err != nil {
			return "", err
		}
		if err := ctx.Err(); err != nil {
			return "", err
		}
	}

	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer func() {
		bufPool.Put(buf)
	}()

	for _, node := range nodes {
		err = html.Render(buf, node)
		if err != nil {
			return "", err
		}
		if err := ctx.Err(); err != nil {
			return "", err
		}
	}

	err = html.Render(buf, &html.Node{
		Type: html.CommentNode,
		Data: fmt.Sprintf("Livewire Component wire-end:%s", obj.(Component).getBaseComponent().ID()),
	})
	if err != nil {
		return "", err
	}
	return buf.String(), nil
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
	if len(nodes) == 0 {
		return nil, ErrInvalidHTMLContent
	}
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
