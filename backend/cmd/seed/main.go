package main

import (
	"log"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/jessicaandtommymccay/realty-wizard/backend/internal/auth"
	"github.com/jessicaandtommymccay/realty-wizard/backend/internal/models"
	"github.com/jessicaandtommymccay/realty-wizard/backend/internal/storage"
)

func main() {
	log.Println("Starting database seed...")

	// Get database path
	dataDir := "data"
	dbPath := filepath.Join(dataDir, "realty-wizard.db")

	// Initialize storage
	store, err := storage.NewSQLiteStorage(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
	defer store.Close()

	// Check if admin user already exists
	existingAdmin, _ := store.GetUserByEmail("admin@realwiz.local")
	if existingAdmin != nil {
		log.Println("Admin user already exists. Skipping seed.")
		return
	}

	// Create admin user
	log.Println("Creating admin user...")
	passwordHash, err := auth.HashPassword("admin123")
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	adminUser := &models.User{
		ID:           uuid.New().String(),
		Email:        "admin@realwiz.local",
		PasswordHash: passwordHash,
		Name:         "Admin User",
		Phone:        "555-0000",
		DefaultRole:  "admin",
		CreatedAt:    time.Now(),
	}

	if err := store.CreateUser(adminUser); err != nil {
		log.Fatalf("Failed to create admin user: %v", err)
	}

	log.Printf("✅ Admin user created: %s", adminUser.Email)

	// Assign all existing projects to admin and create participants
	projects, err := store.ListProjects()
	if err != nil {
		log.Fatalf("Failed to list projects: %v", err)
	}

	if len(projects) > 0 {
		log.Printf("Found %d existing projects. Assigning to admin user...", len(projects))

		for _, project := range projects {
			// Update project owner
			project.OwnerUserID = &adminUser.ID
			if err := store.UpdateProject(project); err != nil {
				log.Printf("Warning: Failed to update project %s: %v", project.ID, err)
				continue
			}

			// Create participant record
			participant := &models.ProjectParticipant{
				ProjectID: project.ID,
				UserID:    adminUser.ID,
				Role:      "owner",
				CreatedAt: time.Now(),
			}

			if err := store.CreateProjectParticipant(participant); err != nil {
				log.Printf("Warning: Failed to create participant for project %s: %v", project.ID, err)
			}
		}

		log.Printf("✅ Assigned %d projects to admin user", len(projects))
	}

	log.Println("✅ Database seed completed successfully!")
	log.Println("\nAdmin credentials:")
	log.Println("  Email: admin@realwiz.local")
	log.Println("  Password: admin123")
	log.Println("\n⚠️  IMPORTANT: Change the admin password after first login!")
}
