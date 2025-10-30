package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	chiadapter "github.com/awslabs/aws-lambda-go-api-proxy/chi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jessicaandtommymccay/realty-wizard/backend/internal/handlers"
	"github.com/jessicaandtommymccay/realty-wizard/backend/internal/storage"
)

var chiLambda *chiadapter.ChiLambda

// init runs once when Lambda container starts (stays warm)
func init() {
	// Database path - Lambda has /tmp for writable storage
	// For persistent storage, we'll use EFS mounted at /mnt/efs
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		// Check if EFS is mounted
		if _, err := os.Stat("/mnt/efs"); err == nil {
			dbPath = "/mnt/efs/realty-wizard.db"
			os.MkdirAll("/mnt/efs", 0755)
		} else {
			// Fallback to /tmp (ephemeral, resets on cold start)
			dbPath = "/tmp/realty-wizard.db"
			log.Println("WARNING: Using /tmp for database - data will not persist!")
		}
	}

	log.Printf("Initializing database at %s", dbPath)

	// Ensure directory exists
	os.MkdirAll(filepath.Dir(dbPath), 0755)

	// Initialize storage
	store, err := storage.NewSQLiteStorage(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	log.Println("Database initialized successfully")

	// Get JWT secret
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "development-secret-key-change-in-production"
		log.Println("WARNING: Using default JWT_SECRET")
	}

	// Initialize handlers
	h := handlers.NewHandler(store, jwtSecret)

	// Setup router (same as regular API)
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://realty-wizard-1760549770.s3-website.us-east-2.amazonaws.com", "*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false, // Must be false when using wildcard
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
		r.Post("/projects/{id}/documents/generate", h.GenerateDocument)
		r.Get("/documents/{id}/download", h.DownloadDocument)
	})

	// Create Lambda adapter
	chiLambda = chiadapter.New(r)
}

// Handler is the Lambda function handler
func Handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	// Convert V2 request to V1 for chi adapter
	v1req := events.APIGatewayProxyRequest{
		HTTPMethod: req.RequestContext.HTTP.Method,
		Path:       req.RawPath,
		Body:       req.Body,
		Headers:    req.Headers,
	}

	v1resp, err := chiLambda.ProxyWithContext(ctx, v1req)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Body:       `{"error":"Internal server error"}`,
		}, err
	}

	// Convert V1 response to V2
	return events.APIGatewayV2HTTPResponse{
		StatusCode: v1resp.StatusCode,
		Headers:    v1resp.Headers,
		Body:       v1resp.Body,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
