package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jessicaandtommymccay/realty-wizard/backend/internal/auth"
	"github.com/jessicaandtommymccay/realty-wizard/backend/internal/documents"
	"github.com/jessicaandtommymccay/realty-wizard/backend/internal/models"
	"github.com/jessicaandtommymccay/realty-wizard/backend/internal/rules"
	"github.com/jessicaandtommymccay/realty-wizard/backend/internal/storage"
)

// Handler holds dependencies for HTTP handlers
type Handler struct {
	storage        storage.Storage
	deadlineEngine *rules.DeadlineEngine
	docGenerator   *documents.Generator
	jwtManager     *auth.JWTManager
}

// NewHandler creates a new handler with dependencies
func NewHandler(store storage.Storage, jwtSecret string) *Handler {
	return &Handler{
		storage:        store,
		deadlineEngine: rules.NewDeadlineEngine(),
		docGenerator:   documents.NewGenerator(),
		jwtManager:     auth.NewJWTManager(jwtSecret),
	}
}

// respondJSON sends a JSON response
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// respondError sends an error response
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

// CreateProject handles POST /api/projects
func (h *Handler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var req models.Project
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Set defaults
	req.ID = uuid.New().String()
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()
	if req.Status == "" {
		req.Status = "setup"
	}

	// Get user from context (if auth is enabled)
	var userID string
	if user, ok := r.Context().Value("user").(*models.User); ok {
		userID = user.ID
		req.OwnerUserID = &userID
	}

	if err := h.storage.CreateProject(&req); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create project")
		return
	}

	// If auth is enabled, add user as owner participant
	if userID != "" {
		participant := &models.ProjectParticipant{
			ProjectID: req.ID,
			UserID:    userID,
			Role:      "owner",
			CreatedAt: time.Now(),
		}
		if err := h.storage.CreateProjectParticipant(participant); err != nil {
			// Log error but don't fail the request
			// The project was created successfully
		}
	}

	respondJSON(w, http.StatusCreated, req)
}

// GetProject handles GET /api/projects/:id
func (h *Handler) GetProject(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	project, err := h.storage.GetProject(id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Project not found")
		return
	}

	respondJSON(w, http.StatusOK, project)
}

// ListProjects handles GET /api/projects
func (h *Handler) ListProjects(w http.ResponseWriter, r *http.Request) {
	// If auth is enabled, filter by user's projects
	var projects []*models.Project
	var err error

	if user, ok := r.Context().Value("user").(*models.User); ok {
		// Auth enabled - get only user's projects
		projects, err = h.storage.ListUserProjects(user.ID)
	} else {
		// Auth disabled - get all projects
		projects, err = h.storage.ListProjects()
	}

	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to list projects")
		return
	}

	if projects == nil {
		projects = []*models.Project{}
	}

	respondJSON(w, http.StatusOK, projects)
}

// UpdateProject handles PUT /api/projects/:id
func (h *Handler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req models.Project
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	req.ID = id
	req.UpdatedAt = time.Now()

	if err := h.storage.UpdateProject(&req); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update project")
		return
	}

	respondJSON(w, http.StatusOK, req)
}

// CreateProperty handles POST /api/projects/:id/property
func (h *Handler) CreateProperty(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")

	var req models.Property
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	req.ID = uuid.New().String()
	req.ProjectID = projectID

	if err := h.storage.CreateProperty(&req); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create property")
		return
	}

	respondJSON(w, http.StatusCreated, req)
}

// GetProperty handles GET /api/projects/:id/property
func (h *Handler) GetProperty(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")

	property, err := h.storage.GetProperty(projectID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Property not found")
		return
	}

	respondJSON(w, http.StatusOK, property)
}

// UpdateProperty handles PUT /api/projects/:id/property
func (h *Handler) UpdateProperty(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")

	var req models.Property
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Ensure project ID matches
	req.ProjectID = projectID

	if err := h.storage.UpdateProperty(&req); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update property")
		return
	}

	respondJSON(w, http.StatusOK, req)
}

// CreateDisclosure handles POST /api/projects/:id/disclosure
func (h *Handler) CreateDisclosure(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")

	var req models.Disclosure
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	req.ID = uuid.New().String()
	req.ProjectID = projectID
	req.UpdatedAt = time.Now()

	// Auto-set lead based paint if property built before 1978
	property, err := h.storage.GetProperty(projectID)
	if err == nil && property.YearBuilt > 0 && property.YearBuilt < 1978 {
		req.LeadBasedPaint = true
	}

	if err := h.storage.CreateDisclosure(&req); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create disclosure")
		return
	}

	respondJSON(w, http.StatusCreated, req)
}

// GetDisclosure handles GET /api/projects/:id/disclosure
func (h *Handler) GetDisclosure(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")

	disclosure, err := h.storage.GetDisclosure(projectID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Disclosure not found")
		return
	}

	respondJSON(w, http.StatusOK, disclosure)
}

// UpdateDisclosure handles PUT /api/projects/:id/disclosure
func (h *Handler) UpdateDisclosure(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")

	var req models.Disclosure
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	req.ProjectID = projectID
	req.UpdatedAt = time.Now()

	if err := h.storage.UpdateDisclosure(&req); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update disclosure")
		return
	}

	respondJSON(w, http.StatusOK, req)
}

// CreateContract handles POST /api/projects/:id/contract
func (h *Handler) CreateContract(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")

	var req models.ContractTerms
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	req.ID = uuid.New().String()
	req.ProjectID = projectID

	// Set default title commitment days if not specified
	if req.TitleCommitmentDays == 0 {
		req.TitleCommitmentDays = 20
	}

	if err := h.storage.CreateContract(&req); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create contract")
		return
	}

	// Auto-generate deadlines when contract is created
	property, _ := h.storage.GetProperty(projectID)
	deadlines := h.deadlineEngine.CalculateDeadlines(&req, property)
	for _, deadline := range deadlines {
		h.storage.CreateDeadline(deadline)
	}

	// Update project status to "under_contract"
	project, _ := h.storage.GetProject(projectID)
	if project != nil {
		project.Status = "under_contract"
		h.storage.UpdateProject(project)
	}

	respondJSON(w, http.StatusCreated, req)
}

// GetContract handles GET /api/projects/:id/contract
func (h *Handler) GetContract(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")

	contract, err := h.storage.GetContract(projectID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Contract not found")
		return
	}

	respondJSON(w, http.StatusOK, contract)
}

// ListDeadlines handles GET /api/projects/:id/deadlines
func (h *Handler) ListDeadlines(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")

	deadlines, err := h.storage.ListDeadlines(projectID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to list deadlines")
		return
	}

	if deadlines == nil {
		deadlines = []*models.Deadline{}
	}

	respondJSON(w, http.StatusOK, deadlines)
}

// UpdateDeadline handles PUT /api/deadlines/:id
func (h *Handler) UpdateDeadline(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req models.Deadline
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	req.ID = id

	if err := h.storage.UpdateDeadline(&req); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update deadline")
		return
	}

	respondJSON(w, http.StatusOK, req)
}

// ListDocuments handles GET /api/projects/:id/documents
func (h *Handler) ListDocuments(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")

	documents, err := h.storage.ListDocuments(projectID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to list documents")
		return
	}

	if documents == nil {
		documents = []*models.Document{}
	}

	respondJSON(w, http.StatusOK, documents)
}

// GenerateDocument handles POST /api/projects/:id/documents/generate
func (h *Handler) GenerateDocument(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")

	var req struct {
		Type string `json:"type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Load project data
	project, err := h.storage.GetProject(projectID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Project not found")
		return
	}

	property, _ := h.storage.GetProperty(projectID)
	disclosure, _ := h.storage.GetDisclosure(projectID)
	contract, _ := h.storage.GetContract(projectID)

	// Check if document type is applicable
	templates := h.docGenerator.AvailableTemplates()
	var template *documents.DocumentTemplate
	for _, t := range templates {
		if t.Type == req.Type {
			template = &t
			break
		}
	}

	if template == nil {
		respondError(w, http.StatusBadRequest, "Unknown document type")
		return
	}

	if !template.Condition(project, property, contract) {
		respondError(w, http.StatusBadRequest, "Document not applicable to this project")
		return
	}

	// Create document record
	doc := &models.Document{
		ID:         uuid.New().String(),
		ProjectID:  projectID,
		Type:       req.Type,
		FormNumber: template.FormNumber,
		Status:     "draft",
	}

	// Generate PDF
	pdfBytes, err := h.docGenerator.GeneratePDF(doc, project, property, disclosure, contract)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate document")
		return
	}

	// In production, save PDF to file storage (S3, etc.)
	// For now, we'll just mark it as generated
	now := time.Now()
	doc.GeneratedAt = &now
	doc.FilePath = "/documents/" + doc.ID + ".pdf"
	doc.Status = "ready"

	// Save document record
	if err := h.storage.CreateDocument(doc); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to save document")
		return
	}

	// Return document info with PDF data
	response := map[string]interface{}{
		"document": doc,
		"pdf":      pdfBytes,
	}

	respondJSON(w, http.StatusCreated, response)
}

// DownloadDocument handles GET /api/documents/:id/download
func (h *Handler) DownloadDocument(w http.ResponseWriter, r *http.Request) {
	docID := chi.URLParam(r, "id")

	doc, err := h.storage.GetDocument(docID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Document not found")
		return
	}

	// Load project data
	project, _ := h.storage.GetProject(doc.ProjectID)
	property, _ := h.storage.GetProperty(doc.ProjectID)
	disclosure, _ := h.storage.GetDisclosure(doc.ProjectID)
	contract, _ := h.storage.GetContract(doc.ProjectID)

	// Regenerate PDF (in production, load from file storage)
	pdfBytes, err := h.docGenerator.GeneratePDF(doc, project, property, disclosure, contract)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate PDF")
		return
	}

	// Set headers for PDF download
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename="+doc.Type+".pdf")
	w.Write(pdfBytes)
}

// GetProjectSummary handles GET /api/projects/:id/summary
// Returns a comprehensive view of project, property, disclosure, contract, deadlines
func (h *Handler) GetProjectSummary(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")

	project, err := h.storage.GetProject(projectID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Project not found")
		return
	}

	property, _ := h.storage.GetProperty(projectID)
	disclosure, _ := h.storage.GetDisclosure(projectID)
	contract, _ := h.storage.GetContract(projectID)
	deadlines, _ := h.storage.ListDeadlines(projectID)
	documents, _ := h.storage.ListDocuments(projectID)

	summary := map[string]interface{}{
		"project":    project,
		"property":   property,
		"disclosure": disclosure,
		"contract":   contract,
		"deadlines":  deadlines,
		"documents":  documents,
	}

	respondJSON(w, http.StatusOK, summary)
}
