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
	Checksum string      `json:"checksum,omitempty"`
	Data     interface{} `json:"data"`
	HTMLHash string      `json:"htmlHash,omitempty"`
}

type messageEffects struct {
	Html  string   `json:"html"`
	Dirty []string `json:"dirty"`
	Path  string   `json:"path,omitempty"`
}

type componentEffects struct {
	Listeners []string `json:"listeners,omitempty"`
	Path      string   `json:"path,omitempty"`
}

type updateAction struct {
	Type    string              `json:"type,omitempty"`
	Payload updateActionPayload `json:"payload,omitempty"`
}

type updateActionPayload struct {
	ID string `json:"id"`

	// For type callMethod
	Method string        `json:"method,omitempty"`
	Params []interface{} `json:"params,omitempty"`

	// For type syncInput
	Name  string      `json:"name,omitempty"`
	Value interface{} `json:"value,omitempty"`
}
