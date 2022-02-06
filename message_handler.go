package golivewire

import (
	"reflect"
	"strings"
)

var (
	EnableMethodCamelCaseSupport = false

	errorInterface = reflect.TypeOf((*error)(nil)).Elem()
)

type messageHandler struct {
	comp baseComponentSupport
	ref  reflect.Value
	req  *messageRequest
}

func (h *messageHandler) OnCallMethod(upd updateAction) error {
	var method reflect.Value
	methodName := coalesceMethodName(upd.Payload.Method)

	switch methodName {
	case "$refresh":
		method = h.ref.MethodByName("Refresh")
		if !method.IsValid() {
			return nil
		}
	case "$set":
		params := upd.Payload.Params
		if len(params) != 2 {
			return ErrBadRequest.Message("invalid number of $set parameters, expect 2 parameters")
		}
		if field, ok := params[0].(string); !ok {
			return ErrBadRequest.Message("invalid $set parameters, expect first param to be string")
		} else {
			field = coalesceMethodName(field)
			return h.doSetField(field, params[1])
		}
	default:
		method = h.ref.MethodByName(methodName)
		if !method.IsValid() {
			return ErrBadRequest.Message("method is not valid or not exist: " + methodName)
		}
	}

	return h.doCallMethod(method, upd.Payload.Params...)
}

func (h *messageHandler) doCallMethod(method reflect.Value, args ...interface{}) error {
	argsRef := make([]reflect.Value, 0, len(args))
	for _, param := range args {
		argsRef = append(argsRef, reflect.ValueOf(param))
	}

	returns := method.Call(argsRef)
	if len(returns) > 0 {
		last := returns[len(returns)-1]
		if err, ok := last.Interface().(error); ok {
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (h *messageHandler) doSetField(field string, val interface{}) error {
	fieldRef := h.ref.Elem().FieldByName(field)
	if !fieldRef.IsValid() {
		return ErrBadRequest.Message("invalid field name: " + field)
	}
	if !fieldRef.CanSet() {
		return ErrBadRequest.Message("invalid field cannot be set: " + field)
	}

	fieldRef.Set(reflect.ValueOf(val))
	return nil
}

func (h *messageHandler) OnSyncInput(upd updateAction) error {
	fieldName := upd.Payload.Name
	val := upd.Payload.Value

	fieldName = coalesceMethodName(fieldName)
	return h.doSetField(fieldName, val)
}

func newMessageHandler(req *messageRequest, comp baseComponentSupport) *messageHandler {
	return &messageHandler{
		comp: comp,
		ref:  reflect.ValueOf(comp),
		req:  req,
	}
}

func HandleComponentMessage(req *messageRequest, comp baseComponentSupport) error {
	hnd := newMessageHandler(req, comp)
	for _, upd := range req.Updates {
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

func coalesceMethodName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return name
	}

	if strings.HasPrefix(name, "$") {
		return name
	}

	if EnableMethodCamelCaseSupport {
		name = strings.ToUpper(name[0:1]) + name[1:]
	}

	return name
}
