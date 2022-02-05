package golivewire

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

var (
	CORSOptions *cors.Options
)

func NewAjaxHandler() http.Handler {
	router := httprouter.New()

	router.POST("/livewire/message/:componentName", wrapHandlerFunc(func(w http.ResponseWriter, r *http.Request, params httprouter.Params) (interface{}, error) {
		name := params.ByName("componentName")
		compFactory, ok := componentRegistry[name]
		if !ok {
			return nil, ErrNotFound.Message("component name is not found: " + name)
		}

		comp, err := compFactory.createInstance(r.Context())
		if err != nil {
			return nil, err
		}

		req := &Request{
			ServerMemo: serverMemo{
				Data: comp,
			},
		}
		if err := bindJSONRequest(r, req); err != nil {
			return nil, ErrBadRequest.Err(err)
		}
		comp.getBaseComponent().id = req.Fingerprint.ID

		err = HandleComponentMessage(req, comp)
		if err != nil {
			return nil, err
		}

		html, err := renderWithDecorators(r.Context(), comp, rendererPipeline...)
		if err != nil {
			return nil, err
		}

		resp := Response{
			Effects: Effects{
				Html:  html,
				Dirty: []string{},
			},
			ServerMemo: serverMemo{
				Data: comp,
			},
		}
		return resp, nil
	}))

	var hnd http.Handler = router
	if CORSOptions != nil {
		c := cors.New(*CORSOptions)
		hnd = c.Handler(router)
	}

	return hnd
}
