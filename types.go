package golivewire

type Renderer interface {
	Render() (string, error)
}

type ComponentFactoryFunc func() interface{}

type MessageRequest struct {
	Fingerprint Fingerprint    `json:"fingerprint,omitempty"`
	ServerMemo  ServerMemo     `json:"serverMemo,omitempty"`
	Updates     []UpdateAction `json:"updates,omitempty"`
}

type MessageResponse struct {
	Effects    MessageEffects `json:"effects,omitempty"`
	ServerMemo ServerMemo     `json:"serverMemo,omitempty"`
}

type ComponentData struct {
	Fingerprint Fingerprint      `json:"fingerprint,omitempty"`
	Effects     ComponentEffects `json:"effects,omitempty"`
	ServerMemo  ServerMemo       `json:"serverMemo,omitempty"`
}

type Fingerprint struct {
	ID     string `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Locale string `json:"locale,omitempty"`
	Path   string `json:"path,omitempty"`
	Method string `json:"method,omitempty"`
}

type ServerMemo struct {
	Data interface{} `json:"data"`
}

type MessageEffects struct {
	Html  string   `json:"html"`
	Dirty []string `json:"dirty"`
}

type ComponentEffects struct {
	Listeners []string `json:"listeners,omitempty"`
}

type UpdateAction struct {
	Type    string              `json:"type,omitempty"`
	Payload UpdateActionPayload `json:"payload,omitempty"`
}

type UpdateActionPayload struct {
	ID     string        `json:"id"`
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
}
