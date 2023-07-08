package util

import (
	"encoding/json"
	"net/http"
)

func EncodeJson(r interface{}) ([]byte, error) {
	return json.Marshal(r)
}

func DecodeJson(w *http.Response, r interface{}) error {
	return json.NewDecoder(w.Body).Decode(&r)
}
