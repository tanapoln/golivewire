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
		_, manager := newManagerCtx(r.Context(), r)
		if err := manager.ProcessRequest(); err != nil {
			return nil, err
		}

		lifecycle, err := newLifecycleFromSubsequentRequest(manager)
		if err != nil {
			return nil, err
		}
		if err := lifecycle.hydrate(); err != nil {
			return nil, err
		}
		if err := lifecycle.renderToView(); err != nil {
			return nil, err
		}
		if err := lifecycle.dehydrate(); err != nil {
			return nil, err
		}
		if err := lifecycle.toSubsequentResponse(); err != nil {
			return nil, err
		}

		return lifecycle.response, nil
	}))

	var hnd http.Handler = router
	if CORSOptions != nil {
		c := cors.New(*CORSOptions)
		hnd = c.Handler(router)
	}

	return hnd
}
