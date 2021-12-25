package golivewire

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type customHandlerFunc func(w http.ResponseWriter, r *http.Request, params httprouter.Params) (interface{}, error)
type H map[string]interface{}

func wrapHandlerFunc(fn customHandlerFunc) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		i, err := fn(w, r, params)
		if err != nil {
			if v, ok := err.(HTTPError); ok {
				w.WriteHeader(v.HTTPStatusCode())
				_ = writeJsonBody(w, H{
					"error": v.Error(),
				})
			} else {
				w.WriteHeader(500)
				_ = writeJsonBody(w, H{
					"error": v.Error(),
				})
			}

			return
		}

		w.WriteHeader(200)
		_ = writeJsonBody(w, i)
	}
}

func writeJsonBody(w http.ResponseWriter, body interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func bindJSONRequest(r *http.Request, o interface{}) error {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, o)
	if err != nil {
		return err
	}

	return nil
}
