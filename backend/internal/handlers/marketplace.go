package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jessicaandtommymccay/realty-wizard/backend/internal/models"
)

// ListServices returns all available services
func (h *Handler) ListServices(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")

	var services []*models.Service
	var err error

	if category != "" {
		services, err = h.storage.ListServicesByCategory(category)
	} else {
		services, err = h.storage.ListServices()
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services)
}

// GetService returns a single service by ID
func (h *Handler) GetService(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	service, err := h.storage.GetService(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(service)
}

// ListProviders returns all providers for a service
func (h *Handler) ListProviders(w http.ResponseWriter, r *http.Request) {
	serviceID := chi.URLParam(r, "serviceId")

	providers, err := h.storage.ListProviders(serviceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(providers)
}

// GetProvider returns a single provider by ID
func (h *Handler) GetProvider(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	provider, err := h.storage.GetProvider(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(provider)
}

// CreateServiceRequest creates a new service request
func (h *Handler) CreateServiceRequest(w http.ResponseWriter, r *http.Request) {
	var request models.ServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set defaults
	request.ID = uuid.New().String()
	request.CreatedAt = time.Now()
	request.UpdatedAt = time.Now()
	if request.Status == "" {
		request.Status = "pending"
	}

	if err := h.storage.CreateServiceRequest(&request); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(request)
}

// GetServiceRequest returns a single service request by ID
func (h *Handler) GetServiceRequest(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	request, err := h.storage.GetServiceRequest(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(request)
}

// ListServiceRequests returns all service requests for a user
func (h *Handler) ListServiceRequests(w http.ResponseWriter, r *http.Request) {
	userEmail := r.URL.Query().Get("user_email")
	if userEmail == "" {
		http.Error(w, "user_email query parameter is required", http.StatusBadRequest)
		return
	}

	requests, err := h.storage.ListServiceRequests(userEmail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requests)
}

// UpdateServiceRequest updates an existing service request
func (h *Handler) UpdateServiceRequest(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var updates models.ServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get existing request
	existing, err := h.storage.GetServiceRequest(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Update fields
	existing.ProviderID = updates.ProviderID
	existing.PropertyAddress = updates.PropertyAddress
	existing.RequestedDate = updates.RequestedDate
	existing.PreferredTime = updates.PreferredTime
	existing.Status = updates.Status
	existing.Notes = updates.Notes
	existing.UpdatedAt = time.Now()

	if err := h.storage.UpdateServiceRequest(existing); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existing)
}

// ListProviderReviews returns all reviews for a provider
func (h *Handler) ListProviderReviews(w http.ResponseWriter, r *http.Request) {
	providerID := chi.URLParam(r, "providerId")

	reviews, err := h.storage.ListProviderReviews(providerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reviews)
}

// CreateProviderReview creates a new provider review
func (h *Handler) CreateProviderReview(w http.ResponseWriter, r *http.Request) {
	var review models.ProviderReview
	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set defaults
	review.ID = uuid.New().String()
	review.CreatedAt = time.Now()

	// Validate rating
	if review.Rating < 1 || review.Rating > 5 {
		http.Error(w, "Rating must be between 1 and 5", http.StatusBadRequest)
		return
	}

	if err := h.storage.CreateProviderReview(&review); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update provider's rating average
	// Get all reviews for this provider
	reviews, err := h.storage.ListProviderReviews(review.ProviderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Calculate new average
	var total float64
	for _, r := range reviews {
		total += float64(r.Rating)
	}
	average := total / float64(len(reviews))

	// Get provider and update rating
	provider, err := h.storage.GetProvider(review.ProviderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	provider.RatingAverage = average
	provider.RatingCount = len(reviews)
	provider.UpdatedAt = time.Now()

	if err := h.storage.UpdateProvider(provider); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(review)
}
