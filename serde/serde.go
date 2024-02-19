package serde

import (
	"encoding/json"
	"fmt"
	"io"
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

// Decode decodes JSON data from the given io.Reader into a value of type T.
// It returns the decoded value and an error if decoding fails.
func Decode[T any](r io.Reader) (T, error) {
	var v T
	if err := json.NewDecoder(r).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}
	return v, nil
}
