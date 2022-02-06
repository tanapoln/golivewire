package golivewire

import (
	"errors"
	"net/http"
	"net/url"

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
		name := params.ByName("componentName")
		compFactory, ok := componentRegistry[name]
		if !ok {
			return nil, ErrNotFound.Message("component name is not found: " + name)
		}

		comp, err := compFactory.createInstance(r.Context())
		if err != nil {
			return nil, err
		}

		req := &messageRequest{
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

		path, err := renderMessageResponsePath(comp, &req.Fingerprint)
		if err != nil {
			return nil, err
		}
		htmlHash, _ := crc32Hash(html)
		resp := messageResponse{
			Effects: messageEffects{
				Html:  html,
				Dirty: []string{},
				Path:  path.String(),
			},
			ServerMemo: serverMemo{
				Checksum: "",
				Data:     comp,
				HTMLHash: htmlHash,
			},
		}
		return resp, nil
	})))

	var hnd http.Handler = router
	if CORSOptions != nil {
		c := cors.New(*CORSOptions)
		hnd = c.Handler(router)
	}

	return hnd
}

func renderMessageResponsePath(comp baseComponentSupport, fingerprint *fingerprint) (*url.URL, error) {
	httpReq := httpRequestFromContext(comp.getBaseComponent().getContext())
	if httpReq == nil {
		return nil, errors.New("invalid component context, expect http request to be exist")
	}
	oriURL := originalURL(httpReq)

	u := url.URL{}
	u.Scheme = oriURL.Scheme
	u.Host = oriURL.Host

	qs := qsStoreFromCtx(comp.getBaseComponent().getContext())
	if qs == nil {
		return nil, ErrInvalidContext
	}
	u.RawQuery = qs.Encode()

	if fingerprint != nil {
		parsed, err := url.Parse(fingerprint.Path)
		if err != nil {
			return nil, err
		}
		u.Path = parsed.Path

	} else {
		u.Path = oriURL.Path
	}

	return &u, nil
}
func originalURL(req *http.Request) *url.URL {
	u := *req.URL
	u.Host = req.Host

	if req.TLS == nil {
		u.Scheme = "http"
	} else {
		u.Scheme = "https"
	}

	if v := req.Header.Get("x-forwarded-host"); v != "" {
		u.Host = v
	}
	if v := req.Header.Get("x-forwarded-proto"); v != "" {
		u.Scheme = v
	}

	return &u
}
