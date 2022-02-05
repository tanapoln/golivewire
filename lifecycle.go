package golivewire

func newLifecycleFromInitialComponent(comp Component) *lifecycleManager {
	l := &lifecycleManager{}
	l.component = comp

	base := comp.getBaseComponent()
	l.response.Fingerprint.ID = base.ID()
	l.response.Fingerprint.Name = base.Name()
	l.response.Fingerprint.Path = base.manager.OriginalPath()
	l.response.Fingerprint.Method = base.manager.OriginalMethod()

	return l
}

type lifecycleManager struct {
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
	return nil
}
