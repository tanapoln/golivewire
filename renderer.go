package golivewire

import (
	"bytes"
	"context"
	"encoding/base32"
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"
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
	rendererPipeline = append(rendererPipeline,
		htmlDecoratorFunc(livewireIdRenderer),
		htmlDecoratorFunc(livewireQuerystring),
	)

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

func livewireQuerystring(node *html.Node, obj interface{}) error {
	comp, ok := obj.(baseComponentSupport)
	if !ok {
		return ErrNotComponent
	}

	qs := qsStoreFromCtx(comp.getBaseComponent().getContext())
	if qs == nil {
		return ErrInvalidContext
	}

	q, err := unbindQuery(comp)
	if err != nil {
		return err
	}

	qs.Merge(q)
	return nil
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

func livewireInitialDataRenderer(node *html.Node, obj interface{}) error {
	var component baseComponentSupport
	var baseComp *BaseComponent
	if v, ok := obj.(baseComponentSupport); !ok {
		return ErrNotComponent
	} else {
		component = v
		baseComp = v.getBaseComponent()
	}

	url, err := renderMessageResponsePath(component, nil)
	if err != nil {
		return err
	}

	htmlHash, _ := crc32HTML(node)

	initData := componentData{
		Fingerprint: fingerprint{
			ID:     baseComp.GetID(),
			Name:   baseComp.name,
			Path:   url.EscapedPath(),
			Method: "GET",
			Locale: "en",
		},
		Effects: componentEffects{
			Listeners: baseComp.Listeners,
			Path:      url.String(),
		},
		ServerMemo: serverMemo{
			Checksum: "",
			Data:     component,
			HTMLHash: htmlHash,
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
	if _, ok := obj.(baseComponentSupport); !ok {
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
		err := decorator.Decorate(nodes[0], obj)
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
		Data: fmt.Sprintf("Livewire Component wire-end:%s", obj.(baseComponentSupport).getBaseComponent().GetID()),
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

func crc32HTML(n *html.Node) (string, error) {
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer func() {
		bufPool.Put(buf)
	}()

	err := html.Render(buf, n)
	if err != nil {
		return "", err
	}

	return crc32Hash(buf.String())
}

func crc32Hash(str string) (string, error) {
	ch := crc32.NewIEEE()
	_, err := ch.Write([]byte(str))
	if err != nil {
		return "", err
	}
	sum := ch.Sum(nil)
	return base32.HexEncoding.EncodeToString(sum), nil
}
