package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func (api *apiConfig) handleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	html, err := template.ParseFiles("layout.html")
	if err != nil {
		fmt.Errorf(err.Error())
	}
	html.Execute(w, api.fileserverHits)
}

func (api *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.fileserverHits += 1
		next.ServeHTTP(w, r)
	})
}
