package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jhonipereira/chirpy/internal/database"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
	jwtSecret      string
	polkaKey       string
}

var debugFlag bool

func init() {
	flag.BoolVar(&debugFlag, "debug", false, "Enable debug mode")
	flag.Parse()
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	if debugFlag {
		// Enable debug mode
		log.Println("Debug mode is enabled")
		// Check if the file exists
		if _, err := os.Stat("database.json"); err == nil {
			// Delete the file
			err := os.Remove("database.json")
			if err != nil {
				log.Println("Error deleting the file:", err)
			}
		}
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	polkaKey := os.Getenv("POLKA_KEY")

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	apiConfig := apiConfig{
		fileserverHits: 0,
		DB:             db,
		jwtSecret:      jwtSecret,
		polkaKey:       polkaKey,
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	fsHandler := apiConfig.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	r.Handle("/app", fsHandler)
	r.Handle("/app/*", fsHandler)

	// Create a subrouter for /api
	apiRouter := chi.NewRouter()
	apiRouter.Get("/reset", apiConfig.handleReset)
	apiRouter.Get("/healthz", handleServerHealth)
	apiRouter.Post("/chirps", apiConfig.handlerChirpsCreate)
	apiRouter.Get("/chirps", apiConfig.handlerChirpsRetrieve)
	apiRouter.Get("/chirps/{chirpID}", apiConfig.handlerChirp)
	apiRouter.Delete("/chirps/{chirpID}", apiConfig.handlerChirpsDelete)

	apiRouter.Post("/polka/webhooks", apiConfig.handlerPolkaWebhook)

	apiRouter.Post("/revoke", apiConfig.handlerRevoke)
	apiRouter.Post("/refresh", apiConfig.handlerRefresh)
	apiRouter.Post("/users", apiConfig.handlerUsersCreate)
	apiRouter.Put("/users", apiConfig.handlerUsersUpdate)
	apiRouter.Post("/login", apiConfig.handlerUsersLogin)

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
