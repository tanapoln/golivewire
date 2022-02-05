package golivewire

import (
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/net/html"
)

func newLifecycleFromSubsequentRequest(manager *livewireManager) (*lifecycleManager, error) {
	l := &lifecycleManager{}
	l.request = manager.req

	comp, err := manager.GetComponentInstance(l.request.Fingerprint.Name, l.request.Fingerprint.ID)
	if err != nil {
		return nil, err
	}
	l.component = comp

	return l, nil
}

func newLifecycleFromInitialComponent(comp Component) *lifecycleManager {
	l := &lifecycleManager{}
	l.component = comp

	base := comp.getBaseComponent()
	l.request.Fingerprint.ID = base.ID()
	l.request.Fingerprint.Name = base.Name()
	l.request.Fingerprint.Path = base.manager.OriginalPath()
	l.request.Fingerprint.Method = base.manager.OriginalMethod()

	return l
}

type lifecycleManager struct {
	request   Request
	component Component
	response  Response
}

func (l *lifecycleManager) Boot() error {
	return nil
}

func (l *lifecycleManager) Hydrate() error {
	err := l.bindDataToComponent()
	if err != nil {
		return err
	}

	if err := l.handleMessage(); err != nil {
		return err
	}

	return nil
}

func (l *lifecycleManager) bindDataToComponent() error {
	if data, ok := l.request.ServerMemo.Data.(map[string]interface{}); ok {
		decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			TagName: "json",
			Result:  l.component,
		})
		if err != nil {
			return err
		}
		if err := decoder.Decode(data); err != nil {
			return err
		}
	}
	return nil
}

func (l *lifecycleManager) InitialHydrate() error {
	return nil
}

func (l *lifecycleManager) Month() error {
	//TODO Sharp: bind params to component
	return nil
}

func (l *lifecycleManager) RenderToView() error {
	_, err := l.component.getBaseComponent().renderToView()
	if err != nil {
		return err
	}

	return nil
}

func (l *lifecycleManager) Dehydrate() error {
	l.copyRequestToResponse()
	l.response.ServerMemo.Data = l.component
	return nil
}

func (l *lifecycleManager) InitialDehydrate() error {
	l.copyRequestToResponse()
	l.response.ServerMemo.Data = l.component
	l.component.getBaseComponent().preRenderView.AppendHtmlNode(&html.Node{
		Type: html.CommentNode,
		Data: "Livewire Component wire-end:" + l.component.getBaseComponent().ID(),
	})

	return nil
}

func (l *lifecycleManager) ToInitialResponse() error {
	comp := l.component.getBaseComponent()
	view := comp.preRenderView

	view.AddWireTag("id", comp.id)

	initialData, err := json.Marshal(l.response)
	if err != nil {
		return err
	}
	view.AddWireTag("initial-data", string(initialData))

	html, err := view.RenderSafe()
	if err != nil {
		return err
	}

	l.response.Effects.Html = html
	return nil
}

func (l *lifecycleManager) ToSubsequentResponse() error {
	comp := l.component.getBaseComponent()
	view := comp.preRenderView

	view.AddWireTag("id", comp.id)

	html, err := view.RenderSafe()
	if err != nil {
		return err
	}

	l.response.Effects.Html = html
	return nil
}

func (l *lifecycleManager) copyRequestToResponse() {
	l.response.Fingerprint = l.request.Fingerprint
	l.response.ServerMemo = l.request.ServerMemo
	l.response.Effects.Dirty = []string{}
}

func (l *lifecycleManager) handleMessage() error {
	hnd := newMessageHandler(&l.request, l.component)
	for _, upd := range l.request.Updates {
		switch upd.Type {
		case "callMethod":
			err := hnd.OnCallMethod(upd)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
