package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jessicaandtommymccay/realty-wizard/backend/internal/handlers"
	"github.com/jessicaandtommymccay/realty-wizard/backend/internal/storage"
)

func main() {
	// Ensure data directory exists
	dataDir := "data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	// Initialize storage
	dbPath := filepath.Join(dataDir, "realty-wizard.db")
	store, err := storage.NewSQLiteStorage(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
	defer store.Close()

	log.Printf("Database initialized at %s", dbPath)

	// Initialize handlers
	h := handlers.NewHandler(store)

	// Setup router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Routes
	r.Route("/api", func(r chi.Router) {
		// Health check
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})

		// Projects
		r.Get("/projects", h.ListProjects)
		r.Post("/projects", h.CreateProject)
		r.Get("/projects/{id}", h.GetProject)
		r.Put("/projects/{id}", h.UpdateProject)
		r.Get("/projects/{id}/summary", h.GetProjectSummary)

		// Property (per project)
		r.Post("/projects/{id}/property", h.CreateProperty)
		r.Get("/projects/{id}/property", h.GetProperty)
		r.Put("/projects/{id}/property", h.UpdateProperty)

		// Disclosure (per project)
		r.Post("/projects/{id}/disclosure", h.CreateDisclosure)
		r.Get("/projects/{id}/disclosure", h.GetDisclosure)
		r.Put("/projects/{id}/disclosure", h.UpdateDisclosure)

		// Contract (per project)
		r.Post("/projects/{id}/contract", h.CreateContract)
		r.Get("/projects/{id}/contract", h.GetContract)

		// Deadlines
		r.Get("/projects/{id}/deadlines", h.ListDeadlines)
		r.Put("/deadlines/{id}", h.UpdateDeadline)

		// Documents
		r.Get("/projects/{id}/documents", h.ListDocuments)

		// Service Marketplace
		r.Get("/services", h.ListServices)
		r.Get("/services/{id}", h.GetService)
		r.Get("/services/{serviceId}/providers", h.ListProviders)
		r.Get("/providers/{id}", h.GetProvider)
		r.Get("/providers/{providerId}/reviews", h.ListProviderReviews)
		r.Post("/providers/{providerId}/reviews", h.CreateProviderReview)
		r.Post("/service-requests", h.CreateServiceRequest)
		r.Get("/service-requests", h.ListServiceRequests)
		r.Get("/service-requests/{id}", h.GetServiceRequest)
		r.Put("/service-requests/{id}", h.UpdateServiceRequest)
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on :%s", port)
	log.Printf("API available at http://localhost:%s/api", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
