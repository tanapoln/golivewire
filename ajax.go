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

	router.POST("/livewire/message/:componentName", ajaxMiddleware(wrapHandlerFunc(func(w http.ResponseWriter, r *http.Request, params httprouter.Params) (interface{}, error) {
		w.Header().Set("cache-control", "max-age=0, must-revalidate, no-cache, no-store, private")

		manager := managerFromCtx(r.Context())
		if err := manager.ProcessRequest(); err != nil {
			return nil, err
		}

		lifecycle, err := newLifecycleFromSubsequentRequest(manager)
		if err != nil {
			return nil, err
		}
		if err := lifecycle.Boot(); err != nil {
			return nil, err
		}
		if err := lifecycle.Hydrate(); err != nil {
			return nil, err
		}
		if err := lifecycle.RenderToView(); err != nil {
			return nil, err
		}
		if err := lifecycle.Dehydrate(); err != nil {
			return nil, err
		}
		if err := lifecycle.ToSubsequentResponse(); err != nil {
			return nil, err
		}

		return lifecycle.response, nil
	})))

	var hnd http.Handler = router
	if CORSOptions != nil {
		c := cors.New(*CORSOptions)
		hnd = c.Handler(router)
	}

	return hnd
}
