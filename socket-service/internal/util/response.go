package util

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func WriteJson(w http.ResponseWriter, code int, payload any) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Println("Error when Marshal json response:", payload)
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func WriteMessage(w http.ResponseWriter, code int, message string) {
	data, _ := json.Marshal(struct {
		Message string `json:"message"`
	}{
		Message: message,
	})
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func WriteError(w http.ResponseWriter, code int, errs ...string) {
	data, _ := json.Marshal(struct {
		Error string `json:"error"`
	}{
		Error: strings.Join(errs, ": "),
	})
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}
