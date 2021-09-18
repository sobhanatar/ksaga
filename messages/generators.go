package messages

import (
	"encoding/json"
	"net/http"
)

const (
	Key   = "Content-Type"
	Value = "application/json"
)

func GenerateSuccessMessage(w *http.ResponseWriter, m map[string]string) (resp []byte) {
	resp, _ = json.Marshal(m)
	(*w).Header().Add(Key, Value)

	return
}

func GenerateRollbackSuccessMessage(w *http.ResponseWriter, m map[string]string) (resp []byte) {
	resp, _ = json.Marshal(m)
	(*w).Header().Add(Key, Value)

	return
}

func GenerateRollbackFailMessage(w *http.ResponseWriter, m map[string]string) (resp []byte) {
	resp, _ = json.Marshal(m)
	(*w).Header().Add(Key, Value)

	return
}
