package main

import (
	"encoding/json"
	"net/http"

	"fmt"
)

func errPage(status int, text string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(status)

	errpageTemplate.Execute(w, fmt.Sprintf("%d, %s", status, text))
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func errAPI(status int, text string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	enc := json.NewEncoder(w)
	err := enc.Encode(&ErrorResponse{
		Error: text,
	})

	if err != nil {
		fmt.Fprintf(w, "Error encoding: '%s'", text)
		return
	}
}
