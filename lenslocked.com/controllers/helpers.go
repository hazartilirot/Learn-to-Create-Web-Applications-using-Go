package controllers

import (
	"github.com/gorilla/schema"
	"net/http"
	"net/url"
)

func parseForm(r *http.Request, dst interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	return parseValues(r.PostForm, dst)
}

func parseURLParams(r *http.Request, dst interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	return parseValues(r.Form, dst)
}

func parseValues(values url.Values, dst interface{}) error {
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	if err := decoder.Decode(dst, values); err != nil {
		return err
	}
	return nil
}
