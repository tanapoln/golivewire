package golivewire

import (
	"encoding/json"

	"github.com/tanapoln/golivewire/lib/mapstructure"
	"golang.org/x/net/html"
)

func newLifecycleFromSubsequentRequest(manager *livewireManager) (*lifecycleManager, error) {
	l := &lifecycleManager{}
	l.manager = manager
	l.request = manager.req

	comp, err := manager.GetComponentInstance(l.request.Fingerprint.Name, l.request.Fingerprint.ID)
	if err != nil {
		return nil, err
	}
	l.component = comp

	return l, nil
}

func newLifecycleFromInitialComponent(manager *livewireManager, componentName string) (*lifecycleManager, error) {
	comp, err := manager.NewComponentInstance(componentName)
	if err != nil {
		return nil, err
	}

	l := &lifecycleManager{}
	l.manager = manager
	l.component = comp

	base := comp.getBaseComponent()
	l.request.Fingerprint.ID = base.ID()
	l.request.Fingerprint.Name = base.Name()
	l.request.Fingerprint.Path = manager.OriginalPath()
	l.request.Fingerprint.Method = manager.OriginalMethod()

	return l, nil
}

type lifecycleManager struct {
	manager   *livewireManager
	request   Request
	component Component
	response  Response
}

func (l *lifecycleManager) Boot() error {
	if err := l.manager.HookDispatch(EventComponentBoot, l); err != nil {
		return err
	}

	if v, ok := l.component.(OnBoot); ok {
		return v.Boot()
	}

	if err := l.manager.HookDispatch(EventComponentBooted, l); err != nil {
		return err
	}

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

	return l.manager.HookDispatch(EventComponentHydrate, l)
}

func (l *lifecycleManager) bindDataToComponent() error {
	if data, ok := l.request.ServerMemo.Data.(map[string]interface{}); ok {
		decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			TagName:          "json",
			Result:           l.component,
			WeaklyTypedInput: true,
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
	if err := l.manager.HookDispatch(EventComponentHydrateInitial, l); err != nil {
		return err
	}

	return nil
}

func (l *lifecycleManager) Mount(param map[string]interface{}) error {
	if err := l.manager.HookDispatch(EventComponentMount, l); err != nil {
		return err
	}

	if param != nil {
		decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			TagName:          "json",
			Result:           l.component,
			WeaklyTypedInput: true,
		})
		if err != nil {
			return err
		}
		if err := decoder.Decode(param); err != nil {
			return err
		}
	}

	if err := l.manager.HookDispatch(EventComponentBooted, l); err != nil {
		return err
	}
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
	if err := l.manager.HookDispatch(EventComponentDehydrate, l); err != nil {
		return err
	}
	if err := l.manager.HookDispatch(EventComponentDehydrateSubsequent, l); err != nil {
		return err
	}
	return nil
}

func (l *lifecycleManager) InitialDehydrate() error {
	l.copyRequestToResponse()
	l.response.ServerMemo.Data = l.component
	l.component.getBaseComponent().preRenderView.AppendHtmlNode(&html.Node{
		Type: html.CommentNode,
		Data: "Livewire Component wire-end:" + l.component.getBaseComponent().ID(),
	})

	if err := l.manager.HookDispatch(EventComponentDehydrate, l); err != nil {
		return err
	}
	if err := l.manager.HookDispatch(EventComponentDehydrateInitial, l); err != nil {
		return err
	}

	return nil
}

func (l *lifecycleManager) ToInitialResponse() error {
	comp := l.component.getBaseComponent()
	view := comp.preRenderView

	initialData, err := json.Marshal(l.response)
	if err != nil {
		return err
	}
	view.AddWireTag("initial-data", string(initialData))

	view.AddWireTag("id", comp.id)

	html, err := view.RenderSafe()
	if err != nil {
		return err
	}

	l.response.Effects.Html = html

	if err := l.manager.HookDispatch(EventMounted, l); err != nil {
		return err
	}
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
		case "syncInput":
			err := hnd.OnSyncInput(upd)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
