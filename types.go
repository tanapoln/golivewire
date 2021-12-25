package golivewire

import "context"

type Renderer interface {
	Render(ctx context.Context) (string, error)
}

type ComponentFactoryFunc func() interface{}

type messageRequest struct {
	Fingerprint fingerprint    `json:"fingerprint,omitempty"`
	ServerMemo  serverMemo     `json:"serverMemo,omitempty"`
	Updates     []updateAction `json:"updates,omitempty"`
}

type messageResponse struct {
	Effects    messageEffects `json:"effects,omitempty"`
	ServerMemo serverMemo     `json:"serverMemo,omitempty"`
}

type componentData struct {
	Fingerprint fingerprint      `json:"fingerprint,omitempty"`
	Effects     componentEffects `json:"effects,omitempty"`
	ServerMemo  serverMemo       `json:"serverMemo,omitempty"`
}

type fingerprint struct {
	ID     string `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Locale string `json:"locale,omitempty"`
	Path   string `json:"path,omitempty"`
	Method string `json:"method,omitempty"`
}

type serverMemo struct {
	Data interface{} `json:"data"`
}

type messageEffects struct {
	Html  string   `json:"html"`
	Dirty []string `json:"dirty"`
}

type componentEffects struct {
	Listeners []string `json:"listeners,omitempty"`
}

type updateAction struct {
	Type    string              `json:"type,omitempty"`
	Payload updateActionPayload `json:"payload,omitempty"`
}

type updateActionPayload struct {
	ID     string        `json:"id"`
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
}
