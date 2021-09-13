package messages

import (
	"encoding/json"
	"net/http"
)

func GenerateSuccessMessage(w *http.ResponseWriter) (resp []byte) {
	resp, _ = json.Marshal(map[string]string{"message": "Transaction completed successfully"})
	(*w).Header().Add("Content-Type", "application/json")

	return
}

func GenerateRollbackFailMessage(w *http.ResponseWriter) (resp []byte) {
	resp, _ = json.Marshal(map[string]string{"message": "Rolling back failed"})
	(*w).Header().Add("Content-Type", "application/json")

	return
}

func GenerateRollbackSuccessMessage(w *http.ResponseWriter) (resp []byte) {
	resp, _ = json.Marshal(map[string]string{"message": "Rolling back completed successfully"})
	(*w).Header().Add("Content-Type", "application/json")

	return
}
