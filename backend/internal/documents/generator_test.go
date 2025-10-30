package documents

import (
	"testing"
	"time"

	"github.com/jessicaandtommymccay/realty-wizard/backend/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestAvailableTemplates(t *testing.T) {
	g := NewGenerator()
	templates := g.AvailableTemplates()

	assert.NotEmpty(t, templates)
	assert.GreaterOrEqual(t, len(templates), 6, "Should have at least 6 templates")

	// Check for key templates
	var hasContract, hasDisclosure, hasLeadPaint bool
	for _, tmpl := range templates {
		switch tmpl.Type {
		case "residential_contract":
			hasContract = true
			assert.Equal(t, "TREC 20-18", tmpl.FormNumber)
		case "sellers_disclosure":
			hasDisclosure = true
			assert.Equal(t, "TREC OP-H", tmpl.FormNumber)
		case "lead_paint":
			hasLeadPaint = true
			assert.Equal(t, "OP-L", tmpl.FormNumber)
		}
	}

	assert.True(t, hasContract, "Should have residential contract template")
	assert.True(t, hasDisclosure, "Should have sellers disclosure template")
	assert.True(t, hasLeadPaint, "Should have lead paint disclosure template")
}

func TestGetApplicableDocuments_WithContract(t *testing.T) {
	g := NewGenerator()

	project := &models.Project{
		ID:              "test-project-1",
		PropertyAddress: "123 Test St",
		SellerNames:     []string{"John Doe"},
	}

	property := &models.Property{
		ProjectID:  "test-project-1",
		YearBuilt:  2015,
		IsHomestead: true,
		HasHOA:     true,
		HOAName:    "Test HOA",
	}

	contract := &models.ContractTerms{
		ProjectID:     "test-project-1",
		BuyerFinancing: true,
		SellerStaysPostClose: false,
	}

	docs := g.GetApplicableDocuments(project, property, contract)

	assert.NotEmpty(t, docs)

	// Count document types
	docTypes := make(map[string]bool)
	for _, doc := range docs {
		docTypes[doc.Type] = true
		assert.Equal(t, "test-project-1", doc.ProjectID)
		assert.Equal(t, "pending", doc.Status)
	}

	// Should include these documents
	assert.True(t, docTypes["sellers_disclosure"], "Should include sellers disclosure")
	assert.True(t, docTypes["residential_contract"], "Should include residential contract")
	assert.True(t, docTypes["third_party_financing"], "Should include financing addendum")
	assert.True(t, docTypes["hoa_addendum"], "Should include HOA addendum")

	// Should NOT include these (conditions not met)
	assert.False(t, docTypes["lead_paint"], "Should not include lead paint (built after 1978)")
	assert.False(t, docTypes["seller_lease"], "Should not include seller lease (seller not staying)")
}

func TestGetApplicableDocuments_Pre1978Property(t *testing.T) {
	g := NewGenerator()

	project := &models.Project{ID: "test-project-1"}
	property := &models.Property{
		ProjectID: "test-project-1",
		YearBuilt: 1970, // Pre-1978
	}
	contract := &models.ContractTerms{ProjectID: "test-project-1"}

	docs := g.GetApplicableDocuments(project, property, contract)

	hasLeadPaint := false
	for _, doc := range docs {
		if doc.Type == "lead_paint" {
			hasLeadPaint = true
		}
	}

	assert.True(t, hasLeadPaint, "Pre-1978 property should require lead paint disclosure")
}

func TestGetApplicableDocuments_WithSellerLease(t *testing.T) {
	g := NewGenerator()

	project := &models.Project{ID: "test-project-1"}
	property := &models.Property{ProjectID: "test-project-1", YearBuilt: 2015}
	contract := &models.ContractTerms{
		ProjectID:            "test-project-1",
		SellerStaysPostClose: true,
		SellerLeaseDays:      30,
	}

	docs := g.GetApplicableDocuments(project, property, contract)

	hasSellerLease := false
	for _, doc := range docs {
		if doc.Type == "seller_lease" {
			hasSellerLease = true
		}
	}

	assert.True(t, hasSellerLease, "Should include seller lease when seller stays post-close")
}

func TestGeneratePDF_SellersDisclosure(t *testing.T) {
	g := NewGenerator()

	project := &models.Project{
		ID:              "test-project-1",
		PropertyAddress: "123 Test St, Austin, TX 78701",
		SellerNames:     []string{"John Doe", "Jane Doe"},
	}

	property := &models.Property{
		ProjectID:   "test-project-1",
		YearBuilt:   2015,
		IsHomestead: true,
	}

	disclosure := &models.Disclosure{
		ProjectID:       "test-project-1",
		FloodingHistory: false,
		InFloodPlain:    false,
		FoundationIssues: false,
		RoofAge:         5,
		RoofIssues:      false,
	}

	doc := &models.Document{
		ID:         "doc-1",
		ProjectID:  "test-project-1",
		Type:       "sellers_disclosure",
		FormNumber: "TREC OP-H",
	}

	pdfBytes, err := g.GeneratePDF(doc, project, property, disclosure, nil)

	assert.NoError(t, err)
	assert.NotEmpty(t, pdfBytes)
	// Check for PDF header to verify it's a valid PDF
	assert.True(t, len(pdfBytes) > 4)
	assert.Equal(t, "%PDF", string(pdfBytes[:4]))
	// Check PDF contains EOF marker
	assert.Contains(t, string(pdfBytes), "%%EOF")
}

func TestGeneratePDF_ResidentialContract(t *testing.T) {
	g := NewGenerator()

	effectiveDate := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	closingDate := time.Date(2024, 4, 15, 0, 0, 0, 0, time.UTC)

	project := &models.Project{
		ID:              "test-project-1",
		PropertyAddress: "123 Test St, Austin, TX 78701",
		SellerNames:     []string{"John Doe"},
	}

	property := &models.Property{
		ProjectID:        "test-project-1",
		YearBuilt:        2015,
		LegalDescription: "Lot 42, Block 7",
		TaxID:            "12345-67890",
	}

	contract := &models.ContractTerms{
		ProjectID:           "test-project-1",
		EffectiveDate:       effectiveDate,
		SalesPrice:          450000,
		ClosingDate:         closingDate,
		OptionFee:           500,
		OptionPeriodDays:    10,
		EarnestMoney:        5000,
		TitleCommitmentDays: 20,
		ExcludedItems:       []string{"chandelier", "curtains"},
		IncludedPersonalItems: []string{"refrigerator", "washer"},
	}

	doc := &models.Document{
		ID:         "doc-1",
		ProjectID:  "test-project-1",
		Type:       "residential_contract",
		FormNumber: "TREC 20-18",
	}

	pdfBytes, err := g.GeneratePDF(doc, project, property, nil, contract)

	assert.NoError(t, err)
	assert.NotEmpty(t, pdfBytes)
	// Check for PDF header to verify it's a valid PDF
	assert.True(t, len(pdfBytes) > 4)
	assert.Equal(t, "%PDF", string(pdfBytes[:4]))
	// Check PDF contains EOF marker
	assert.Contains(t, string(pdfBytes), "%%EOF")
}

func TestGeneratePDF_LeadPaintDisclosure(t *testing.T) {
	g := NewGenerator()

	project := &models.Project{
		ID:              "test-project-1",
		PropertyAddress: "123 Old St, Austin, TX 78701",
	}

	property := &models.Property{
		ProjectID: "test-project-1",
		YearBuilt: 1975, // Pre-1978
	}

	doc := &models.Document{
		ID:         "doc-1",
		ProjectID:  "test-project-1",
		Type:       "lead_paint",
		FormNumber: "OP-L",
	}

	pdfBytes, err := g.GeneratePDF(doc, project, property, nil, nil)

	assert.NoError(t, err)
	assert.NotEmpty(t, pdfBytes)
	// Check for PDF header to verify it's a valid PDF
	assert.True(t, len(pdfBytes) > 4)
	assert.Equal(t, "%PDF", string(pdfBytes[:4]))
	// Check PDF contains EOF marker
	assert.Contains(t, string(pdfBytes), "%%EOF")
}

func TestGeneratePDF_HOAAddendum(t *testing.T) {
	g := NewGenerator()

	project := &models.Project{
		ID:              "test-project-1",
		PropertyAddress: "123 Test St, Austin, TX 78701",
	}

	property := &models.Property{
		ProjectID:    "test-project-1",
		YearBuilt:    2015,
		HasHOA:       true,
		HOAName:      "Test HOA",
		HOAFee:       150.00,
		HOAFrequency: "monthly",
	}

	doc := &models.Document{
		ID:         "doc-1",
		ProjectID:  "test-project-1",
		Type:       "hoa_addendum",
		FormNumber: "TREC 36-10",
	}

	pdfBytes, err := g.GeneratePDF(doc, project, property, nil, nil)

	assert.NoError(t, err)
	assert.NotEmpty(t, pdfBytes)
	// Check for PDF header to verify it's a valid PDF
	assert.True(t, len(pdfBytes) > 4)
	assert.Equal(t, "%PDF", string(pdfBytes[:4]))
	// Check PDF contains EOF marker
	assert.Contains(t, string(pdfBytes), "%%EOF")
}

func TestGeneratePDF_UnknownType(t *testing.T) {
	g := NewGenerator()

	doc := &models.Document{
		Type: "unknown_type",
	}

	_, err := g.GeneratePDF(doc, nil, nil, nil, nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown document type")
}
