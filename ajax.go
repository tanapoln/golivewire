package golivewire

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/gorilla/mux"
)

func NewAjaxHandler() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/livewire/message/{componentName}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(404)
			return
		}

		vars := mux.Vars(r)
		name, ok := vars["componentName"]
		if !ok {
			w.WriteHeader(400)
			return
		}

		factory, ok := componentRegistry[name]
		if !ok {
			w.WriteHeader(404)
			return
		}

		comp := factory.createInstance()
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(400)
			return
		}

		req := &MessageRequest{
			ServerMemo: ServerMemo{
				Data: comp,
			},
		}

		err = json.Unmarshal(body, req)
		if err != nil {
			w.WriteHeader(400)
			return
		}

		comp.getBaseComponent().id = req.Fingerprint.ID

		compf := reflect.ValueOf(comp)
		for _, upd := range req.Updates {
			switch upd.Type {
			case "callMethod":
				methodName := upd.Payload.Method
				args := make([]reflect.Value, 0, len(upd.Payload.Params))
				for _, param := range upd.Payload.Params {
					args = append(args, reflect.ValueOf(param))
				}

				switch methodName {
				case "$refresh":
					method := compf.MethodByName("Refresh")
					if method.IsValid() {
						method.Call(args)
					}
				default:
					method := compf.MethodByName(methodName)
					if method.IsValid() {
						method.Call(args)
					} else {
						fmt.Printf("Error method is invalid: %s\n", methodName)
					}
				}
			}
		}

		html, err := renderWithDecorators(comp, rendererPipeline...)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		resp := MessageResponse{
			Effects: MessageEffects{
				Html:  html,
				Dirty: []string{},
			},
			ServerMemo: ServerMemo{
				Data: comp,
			},
		}
		respBody, err := json.Marshal(resp)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(200)
		_, _ = w.Write(respBody)
	})

	return router
}
