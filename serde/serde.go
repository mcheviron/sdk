package serde

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Encode encodes the given value v as JSON and writes it to the http.ResponseWriter w.
// The status parameter is used to set the HTTP status code.
// The generic parameter T can be elided and inferred from the type of the value v.
func Encode[T any](w http.ResponseWriter, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

// Decode decodes the JSON data from the given HTTP request into a value of type T.
// The generic parameter T needs to be passed while invoking the function.
func Decode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}
	return v, nil
}
