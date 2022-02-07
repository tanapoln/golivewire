package golivewire

import (
	"context"
	"html/template"
	"net/url"
)

type Renderer interface {
	Render(ctx context.Context) (string, error)
}

type OnBoot interface {
	Boot() error
}

type ComponentValidator interface {
	Validate(ctx context.Context) error
}

type Querystringer interface {
	Querystring() url.Values
}

type LifecycleHook interface {
	Execute(ctx context.Context, component Component, request *Request, response *Response) error
}

type LifecycleHookFunc func(ctx context.Context, component Component, request *Request, response *Response) error

func (f LifecycleHookFunc) Execute(ctx context.Context, component Component, request *Request, response *Response) error {
	return f(ctx, component, request, response)
}

type ComponentFactoryFunc func() Component

type Request struct {
	Fingerprint fingerprint    `json:"fingerprint,omitempty"`
	ServerMemo  serverMemo     `json:"serverMemo,omitempty"`
	Updates     []updateAction `json:"updates,omitempty"`
}

type Response struct {
	Fingerprint fingerprint `json:"fingerprint,omitempty"`
	Effects     effects     `json:"effects,omitempty"`
	ServerMemo  serverMemo  `json:"serverMemo,omitempty"`
}

type fingerprint struct {
	ID     string `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Locale string `json:"locale,omitempty"`
	Path   string `json:"path,omitempty"`
	Method string `json:"method,omitempty"`
}

type serverMemo struct {
	Data     interface{}      `json:"data"`
	HtmlHash string           `json:"htmlHash,omitempty"`
	Checksum string           `json:"checksum,omitempty"`
	Children []childComponent `json:"children,omitempty"`
	DataMeta dataMeta         `json:"dataMeta,omitempty"`
}

type childComponent struct {
	ID  string `json:"id"`
	Tag string `json:"tag"`
}

type dataMeta struct {
	Date        map[string]string `json:"date,omitempty"`
	Collections interface{}       `json:"collections,omitempty"`
	Wirables    []interface{}     `json:"wirables,omitempty"`
	Stringables []interface{}     `json:"stringables,omitempty"`
}

type effects struct {
	Html       template.HTML          `json:"html"`
	Dirty      []string               `json:"dirty"`
	HtmlHash   string                 `json:"htmlHash,omitempty"`
	Returns    map[string]interface{} `json:"returns,omitempty"`
	Path       string                 `json:"path,omitempty"`
	Listeners  []string               `json:"listeners,omitempty"`
	Emits      []interface{}          `json:"emits,omitempty"`
	Dispatches []interface{}          `json:"dispatches,omitempty"`
	Download   interface{}            `json:"download,omitempty"`
	Redirect   interface{}            `json:"redirect,omitempty"`
	ForStack   interface{}            `json:"forStack,omitempty"`
}

type updateAction struct {
	// Possible: callMethod, syncInput, fireEvent
	Type    string              `json:"type,omitempty"`
	Payload updateActionPayload `json:"payload,omitempty"`
}

type updateActionPayload struct {
	// for callMethod and fireEvent
	ID string `json:"id"`
	// for callMethod and fireEvent
	Params []interface{} `json:"params"`

	// for callMethod. common method: $sync, $set, $toggle, $refresh. see: \Livewire\ComponentConcerns\HandlesActions::callMethod
	Method string `json:"method"`

	// for syncInput
	Name  string      `json:"name"`
	Value interface{} `json:"value"`

	// for fireEvent
	Event string `json:"event"`
}

// ErrorBag hold error message. it is not threadsafe.
type ErrorBag struct {
	data map[string][]string
}

func (e *ErrorBag) AddError(key string, message string) {
	if e.data == nil {
		e.data = make(map[string][]string)
	}
	e.data[key] = append(e.data[key], message)
}

func (e *ErrorBag) Keys() []string {
	if e.data == nil {
		return []string{}
	}

	keys := make([]string, 0, len(e.data))
	for k := range e.data {
		keys = append(keys, k)
	}
	return keys
}

func (e *ErrorBag) Messages() map[string][]string {
	if e.data == nil {
		return map[string][]string{}
	}

	return e.data
}

func (e *ErrorBag) Get(key string) []string {
	return e.Messages()[key]
}

func (e *ErrorBag) GetFirst(key string) string {
	arr := e.Get(key)
	if len(arr) > 0 {
		return arr[0]
	} else {
		return ""
	}
}

func (e *ErrorBag) Any() bool {
	return len(e.data) > 0
}
