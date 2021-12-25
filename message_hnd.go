package golivewire

import (
	"fmt"
	"reflect"
	"strings"
)

var (
	EnableMethodCamelCaseSupport = false

	errorInterface = reflect.TypeOf((*error)(nil)).Elem()
)

func HandleComponentMessage(req *messageRequest, comp baseComponentSupport) error {
	compf := reflect.ValueOf(comp)
	for _, upd := range req.Updates {
		switch upd.Type {
		case "callMethod":
			methodName := coalesceMethodName(upd.Payload.Method)
			args := make([]reflect.Value, 0, len(upd.Payload.Params))
			for _, param := range upd.Payload.Params {
				args = append(args, reflect.ValueOf(param))
			}

			var method reflect.Value
			switch methodName {
			case "$refresh":
				method = compf.MethodByName("Refresh")
			default:
				method = compf.MethodByName(methodName)
			}
			if method.IsValid() {
				returns := method.Call(args)
				if len(returns) > 0 {
					last := returns[len(returns)-1]
					if err, ok := last.Interface().(error); ok {
						if err != nil {
							return err
						}
					}
					// if last.Kind() == reflect.Interface && last.Type().Implements(errorInterface) {
					// 	last.Interface().(error)
					// }
				}
			} else {
				fmt.Printf("Error method is invalid: %s\n", methodName)
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
