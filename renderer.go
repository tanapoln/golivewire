package golivewire

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

type HTMLDecorator interface {
	Decorate(node *html.Node, component interface{}) error
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

type htmlDecoratorFunc func(node *html.Node, component interface{}) error

func (h htmlDecoratorFunc) Decorate(node *html.Node, component interface{}) error {
	return h(node, component)
}

func livewireIdRenderer(node *html.Node, component interface{}) error {
	var baseComp *BaseComponent
	if v, ok := component.(baseComponentSupport); !ok {
		return ErrNotComponent
	} else {
		baseComp = v.getBaseComponent()
	}

	node.Attr = append(node.Attr, html.Attribute{
		Key: "wire:id",
		Val: baseComp.GetID(),
	})

	return nil
}

func livewireInitialDataRenderer(node *html.Node, component interface{}) error {
	var baseComp *BaseComponent
	if v, ok := component.(baseComponentSupport); !ok {
		return ErrNotComponent
	} else {
		baseComp = v.getBaseComponent()
	}

	initData := componentData{
		Fingerprint: fingerprint{
			ID:   baseComp.GetID(),
			Name: baseComp.Name,
		},
		Effects: componentEffects{
			Listeners: baseComp.Listeners,
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

func InitialRender(obj interface{}) (string, error) {
	return renderWithDecorators(obj, initialRendererPipeline...)
}

func renderWithDecorators(obj interface{}, decorators ...HTMLDecorator) (string, error) {
	if _, ok := obj.(Renderer); !ok {
		return "", ErrNotRenderer
	}
	if _, ok := obj.(baseComponentSupport); !ok {
		return "", ErrNotComponent
	}

	raw, err := obj.(Renderer).Render()
	if err != nil {
		return "", err
	}

	doc, err := html.Parse(strings.NewReader(raw))
	if err != nil {
		return "", err
	}

	node, err := extractNodeFromDoc(doc)
	if err != nil {
		return "", err
	}

	for _, decorator := range decorators {
		err := decorator.Decorate(node, obj)
		if err != nil {
			return "", err
		}
	}

	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer func() {
		bufPool.Put(buf)
	}()

	err = html.Render(buf, node)
	if err != nil {
		return "", err
	}
	err = html.Render(buf, &html.Node{
		Type: html.CommentNode,
		Data: fmt.Sprintf("Livewire Component wire-end:%s", obj.(baseComponentSupport).getBaseComponent().GetID()),
	})
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func extractNodeFromDoc(node *html.Node) (*html.Node, error) {
	if node.Type != html.DocumentNode {
		return nil, ErrInvalidHTMLContent
	}
	body := node.FirstChild.FirstChild.NextSibling

	if body.Type != html.ElementNode && body.Data != "body" {
		return nil, ErrInvalidHTMLContent
	}

	first := body.FirstChild
	for n := first.NextSibling; n != nil; n = n.NextSibling {
		if n.Type == html.ElementNode {
			return nil, ErrInvalidHTMLContent
		}
	}
	return first, nil
}
