package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jessicaandtommymccay/realty-wizard/backend/internal/models"
	"github.com/jessicaandtommymccay/realty-wizard/backend/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestHandler creates a handler with in-memory storage for testing
func setupTestHandler(t *testing.T) *Handler {
	store, err := storage.NewSQLiteStorage(":memory:")
	require.NoError(t, err)
	return NewHandler(store)
}

// createTestProject creates a test project and returns its ID
func createTestProject(t *testing.T, h *Handler) string {
	project := &models.Project{
		ID:              "test-project-1",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		PropertyAddress: "123 Test St, Austin, TX 78701",
		SellerNames:     []string{"John Doe", "Jane Doe"},
		SellerEmail:     "seller@example.com",
		SellerPhone:     "512-555-0100",
		HasAgent:        true,
		AgentName:       "Agent Smith",
		Status:          "setup",
	}

	err := h.storage.CreateProject(project)
	require.NoError(t, err)
	return project.ID
}

// createTestProperty creates a test property for a project
func createTestProperty(t *testing.T, h *Handler, projectID string) {
	surveyDate := time.Date(2020, 1, 15, 0, 0, 0, 0, time.UTC)
	property := &models.Property{
		ID:               "test-property-1",
		ProjectID:        projectID,
		YearBuilt:        2015,
		LegalDescription: "Lot 42, Block 7, Test Subdivision",
		TaxID:            "12345-67890",
		IsHomestead:      true,
		HasHOA:           true,
		HOAName:          "Test HOA",
		HOAFee:           150.00,
		HOAFrequency:     "monthly",
		HasSurvey:        true,
		SurveyDate:       &surveyDate,
	}

	err := h.storage.CreateProperty(property)
	require.NoError(t, err)
}

func TestCreateContract_Success(t *testing.T) {
	h := setupTestHandler(t)
	projectID := createTestProject(t, h)
	createTestProperty(t, h, projectID)

	effectiveDate := time.Now()
	closingDate := effectiveDate.AddDate(0, 0, 45)

	contractReq := map[string]interface{}{
		"effective_date":        effectiveDate.Format(time.RFC3339),
		"sales_price":           450000.00,
		"closing_date":          closingDate.Format(time.RFC3339),
		"option_fee":            500.00,
		"option_period_days":    10,
		"earnest_money":         5000.00,
		"title_commitment_days": 20,
		"buyer_financing":       true,
	}

	body, _ := json.Marshal(contractReq)
	req := httptest.NewRequest("POST", "/api/projects/"+projectID+"/contract", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Setup chi URLParam
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", projectID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()
	h.CreateContract(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.ContractTerms
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	assert.NotEmpty(t, response.ID)
	assert.Equal(t, projectID, response.ProjectID)
	assert.Equal(t, 450000.00, response.SalesPrice)
	assert.Equal(t, 500.00, response.OptionFee)
	assert.Equal(t, 10, response.OptionPeriodDays)
	assert.Equal(t, 5000.00, response.EarnestMoney)
	assert.Equal(t, 20, response.TitleCommitmentDays)
	assert.True(t, response.BuyerFinancing)
}

func TestCreateContract_GeneratesDeadlines(t *testing.T) {
	h := setupTestHandler(t)
	projectID := createTestProject(t, h)
	createTestProperty(t, h, projectID)

	effectiveDate := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC) // Friday
	closingDate := time.Date(2024, 4, 15, 0, 0, 0, 0, time.UTC)

	contractReq := map[string]interface{}{
		"effective_date":        effectiveDate.Format(time.RFC3339),
		"sales_price":           450000.00,
		"closing_date":          closingDate.Format(time.RFC3339),
		"option_fee":            500.00,
		"option_period_days":    10,
		"earnest_money":         5000.00,
		"title_commitment_days": 20,
		"buyer_financing":       true,
	}

	body, _ := json.Marshal(contractReq)
	req := httptest.NewRequest("POST", "/api/projects/"+projectID+"/contract", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", projectID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()
	h.CreateContract(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// Check that deadlines were created
	deadlines, err := h.storage.ListDeadlines(projectID)
	require.NoError(t, err)
	assert.NotEmpty(t, deadlines, "Should generate deadlines")

	// Count deadline types
	deadlineTypes := make(map[string]int)
	for _, d := range deadlines {
		deadlineTypes[d.Type]++
	}

	// Verify critical deadlines exist
	assert.Equal(t, 1, deadlineTypes["option_fee"], "Should have option fee deadline")
	assert.Equal(t, 1, deadlineTypes["earnest_money"], "Should have earnest money deadline")
	assert.Equal(t, 1, deadlineTypes["inspection_period"], "Should have inspection period deadline")
	assert.Equal(t, 1, deadlineTypes["title_commitment"], "Should have title commitment deadline")
	assert.Equal(t, 1, deadlineTypes["closing"], "Should have closing deadline")
	assert.Equal(t, 1, deadlineTypes["loan_approval"], "Should have loan approval deadline for financed deal")
}

func TestCreateContract_UpdatesProjectStatus(t *testing.T) {
	h := setupTestHandler(t)
	projectID := createTestProject(t, h)

	effectiveDate := time.Now()
	closingDate := effectiveDate.AddDate(0, 0, 45)

	contractReq := map[string]interface{}{
		"effective_date":        effectiveDate.Format(time.RFC3339),
		"sales_price":           450000.00,
		"closing_date":          closingDate.Format(time.RFC3339),
		"option_fee":            500.00,
		"option_period_days":    10,
		"earnest_money":         5000.00,
		"title_commitment_days": 20,
		"buyer_financing":       false,
	}

	body, _ := json.Marshal(contractReq)
	req := httptest.NewRequest("POST", "/api/projects/"+projectID+"/contract", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", projectID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()
	h.CreateContract(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// Verify project status was updated
	project, err := h.storage.GetProject(projectID)
	require.NoError(t, err)
	assert.Equal(t, "under_contract", project.Status)
}

func TestCreateContract_WithSellerLease(t *testing.T) {
	h := setupTestHandler(t)
	projectID := createTestProject(t, h)

	effectiveDate := time.Now()
	closingDate := effectiveDate.AddDate(0, 0, 45)

	contractReq := map[string]interface{}{
		"effective_date":          effectiveDate.Format(time.RFC3339),
		"sales_price":             450000.00,
		"closing_date":            closingDate.Format(time.RFC3339),
		"option_fee":              500.00,
		"option_period_days":      10,
		"earnest_money":           5000.00,
		"title_commitment_days":   20,
		"buyer_financing":         false,
		"seller_stays_post_close": true,
		"seller_lease_days":       30,
	}

	body, _ := json.Marshal(contractReq)
	req := httptest.NewRequest("POST", "/api/projects/"+projectID+"/contract", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", projectID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()
	h.CreateContract(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.ContractTerms
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	assert.True(t, response.SellerStaysPostClose)
	assert.Equal(t, 30, response.SellerLeaseDays)

	// Check for lease_end deadline
	deadlines, err := h.storage.ListDeadlines(projectID)
	require.NoError(t, err)

	var hasLeaseEnd bool
	for _, d := range deadlines {
		if d.Type == "lease_end" {
			hasLeaseEnd = true
			break
		}
	}
	assert.True(t, hasLeaseEnd, "Should have lease end deadline when seller stays post-close")
}

func TestCreateContract_WithExcludedAndIncludedItems(t *testing.T) {
	h := setupTestHandler(t)
	projectID := createTestProject(t, h)

	effectiveDate := time.Now()
	closingDate := effectiveDate.AddDate(0, 0, 45)

	contractReq := map[string]interface{}{
		"effective_date":            effectiveDate.Format(time.RFC3339),
		"sales_price":               450000.00,
		"closing_date":              closingDate.Format(time.RFC3339),
		"option_fee":                500.00,
		"option_period_days":        10,
		"earnest_money":             5000.00,
		"title_commitment_days":     20,
		"buyer_financing":           false,
		"excluded_items":            []string{"chandelier", "custom curtains"},
		"included_personal_items":   []string{"refrigerator", "washer", "dryer"},
	}

	body, _ := json.Marshal(contractReq)
	req := httptest.NewRequest("POST", "/api/projects/"+projectID+"/contract", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", projectID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()
	h.CreateContract(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.ContractTerms
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	assert.ElementsMatch(t, []string{"chandelier", "custom curtains"}, response.ExcludedItems)
	assert.ElementsMatch(t, []string{"refrigerator", "washer", "dryer"}, response.IncludedPersonalItems)
}

func TestCreateContract_DefaultsTitleCommitmentDays(t *testing.T) {
	h := setupTestHandler(t)
	projectID := createTestProject(t, h)

	effectiveDate := time.Now()
	closingDate := effectiveDate.AddDate(0, 0, 45)

	contractReq := map[string]interface{}{
		"effective_date":     effectiveDate.Format(time.RFC3339),
		"sales_price":        450000.00,
		"closing_date":       closingDate.Format(time.RFC3339),
		"option_fee":         500.00,
		"option_period_days": 10,
		"earnest_money":      5000.00,
		"buyer_financing":    false,
		// Omit title_commitment_days to test default
	}

	body, _ := json.Marshal(contractReq)
	req := httptest.NewRequest("POST", "/api/projects/"+projectID+"/contract", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", projectID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()
	h.CreateContract(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.ContractTerms
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, 20, response.TitleCommitmentDays, "Should default to 20 days")
}

func TestCreateContract_InvalidJSON(t *testing.T) {
	h := setupTestHandler(t)
	projectID := createTestProject(t, h)

	req := httptest.NewRequest("POST", "/api/projects/"+projectID+"/contract", bytes.NewBufferString("{invalid json"))
	req.Header.Set("Content-Type", "application/json")

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", projectID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()
	h.CreateContract(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	json.NewDecoder(w.Body).Decode(&response)
	assert.Equal(t, "Invalid request body", response["error"])
}

func TestGetContract_Success(t *testing.T) {
	h := setupTestHandler(t)
	projectID := createTestProject(t, h)

	// First create a contract
	effectiveDate := time.Now()
	contract := &models.ContractTerms{
		ID:                  "test-contract-1",
		ProjectID:           projectID,
		EffectiveDate:       effectiveDate,
		SalesPrice:          450000.00,
		ClosingDate:         effectiveDate.AddDate(0, 0, 45),
		OptionFee:           500.00,
		OptionPeriodDays:    10,
		EarnestMoney:        5000.00,
		TitleCommitmentDays: 20,
		BuyerFinancing:      true,
	}
	err := h.storage.CreateContract(contract)
	require.NoError(t, err)

	// Now retrieve it
	req := httptest.NewRequest("GET", "/api/projects/"+projectID+"/contract", nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", projectID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()
	h.GetContract(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.ContractTerms
	err = json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, contract.ID, response.ID)
	assert.Equal(t, contract.SalesPrice, response.SalesPrice)
	assert.Equal(t, contract.OptionFee, response.OptionFee)
}

func TestGetContract_NotFound(t *testing.T) {
	h := setupTestHandler(t)
	projectID := "nonexistent-project"

	req := httptest.NewRequest("GET", "/api/projects/"+projectID+"/contract", nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", projectID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()
	h.GetContract(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]string
	json.NewDecoder(w.Body).Decode(&response)
	assert.Equal(t, "Contract not found", response["error"])
}
