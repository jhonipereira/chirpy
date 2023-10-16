package main

import "net/http"

func (api *apiConfig) handleReset(w http.ResponseWriter, r *http.Request) {
	api.fileserverHits = 0
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
