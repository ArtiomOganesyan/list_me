package server

import (
	"encoding/json"
	"net/http"
)

func errorResponse(w http.ResponseWriter, status int, err error) {
	jsonResponse(w, status, map[string]string{"error": err.Error()})
}

func jsonResponse(w http.ResponseWriter, status int, data any) {
	response, err := json.Marshal(data)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

func statusResponse(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}

func panicRecover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				errorResponse(w, http.StatusInternalServerError, err.(error))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
