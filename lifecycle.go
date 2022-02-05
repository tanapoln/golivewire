package golivewire

type lifecycleManager struct {
	component Component
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
	return nil
}

func (l *lifecycleManager) renderToView() error {
	return nil
}

func (l *lifecycleManager) dehydrate() error {
	return nil
}

func (l *lifecycleManager) initialDehydrate() error {
	return nil
}

func (l *lifecycleManager) toInitialResponse() error {
	return nil
}

func (l *lifecycleManager) toSubsequentResponse() error {
	return nil
}
