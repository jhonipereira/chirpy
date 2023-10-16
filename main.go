package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
	const filepathRoot = "."
	const port = "8080"
	var apiConfig apiConfig

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	fsHandler := apiConfig.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	r.Handle("/app", fsHandler)
	r.Handle("/app/*", fsHandler)

	// Create a subrouter for /api
	apiRouter := chi.NewRouter()
	apiRouter.Get("/reset", apiConfig.handleReset)
	apiRouter.Get("/healthz", handleServerHealth)
	apiRouter.Post("/validate_chirp", handleValidateChirp)

	r.Mount("/api", apiRouter)

	// Create a subrouter for /admin
	adminRouter := chi.NewRouter()
	adminRouter.Get("/metrics", apiConfig.handleMetrics)
	r.Mount("/admin", adminRouter)

	corsMux := middlewareCors(r)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Server running on port: %s", port)
	log.Fatal(srv.ListenAndServe())
}
