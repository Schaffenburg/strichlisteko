package main

import (
	"net/http"
)

func AddPrefix(prefix string, h http.Handler) http.Handler {
	if prefix == "" {
		return h
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = prefix + r.URL.Path

		h.ServeHTTP(w, r)
	})

}
