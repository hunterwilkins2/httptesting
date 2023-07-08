// Package util Utility library for httptesting
package util

import (
	"encoding/json"
	"net/http"
)

// EncodeJSON helper function for encoding a struct to JSON
func EncodeJSON(r interface{}) ([]byte, error) {
	return json.Marshal(r)
}

// DecodeJSON helper function for decoding a JSON response body into a struct
func DecodeJSON(w *http.Response, r interface{}) error {
	return json.NewDecoder(w.Body).Decode(&r)
}
