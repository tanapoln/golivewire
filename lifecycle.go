package golivewire

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

func (l *lifecycleManager) boot() error {
	return nil
}

func (l *lifecycleManager) hydrate() error {
	return nil
}

func (l *lifecycleManager) initialHydrate() error {
	return nil
}

func (l *lifecycleManager) month() error {
	//TODO Sharp: bind params to component
	return nil
}

func (l *lifecycleManager) renderToView() error {
	_, err := l.component.getBaseComponent().renderToView()
	if err != nil {
		return err
	}

	return nil
}

func (l *lifecycleManager) dehydrate() error {
	return nil
}

func (l *lifecycleManager) initialDehydrate() error {
	return nil
}

func (l *lifecycleManager) toInitialResponse() error {
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

func (l *lifecycleManager) toSubsequentResponse() error {
	comp := l.component.getBaseComponent()
	view := comp.preRenderView
	html, err := view.RenderSafe()
	if err != nil {
		return err
	}

	l.response.Effects.Html = html
	l.response.ServerMemo.Data = l.component
	return nil
}
