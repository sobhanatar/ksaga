package messages

import (
	"encoding/json"
	"net/http"
)

const (
	Key   = "Content-Type"
	Value = "application/json"
)

func GenerateMessage(w *http.ResponseWriter, m map[string]interface{}) (resp []byte) {
	resp, _ = json.Marshal(m)
	(*w).Header().Add(Key, Value)

	return
}
