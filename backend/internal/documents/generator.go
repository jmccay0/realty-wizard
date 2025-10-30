package documents

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/jessicaandtommymccay/realty-wizard/backend/internal/models"
	"github.com/jung-kurt/gofpdf"
)

// Generator handles document generation
type Generator struct{}

// NewGenerator creates a new document generator
func NewGenerator() *Generator {
	return &Generator{}
}

// DocumentTemplate represents a document that can be generated
type DocumentTemplate struct {
	Type       string
	FormNumber string
	Name       string
	Required   bool
	Condition  func(*models.Project, *models.Property, *models.ContractTerms) bool
}

// AvailableTemplates returns all document templates
func (g *Generator) AvailableTemplates() []DocumentTemplate {
	return []DocumentTemplate{
		{
			Type:       "sellers_disclosure",
			FormNumber: "TREC OP-H",
			Name:       "Seller's Disclosure Notice",
			Required:   true,
			Condition:  func(p *models.Project, prop *models.Property, c *models.ContractTerms) bool { return prop != nil },
		},
		{
			Type:       "lead_paint",
			FormNumber: "OP-L",
			Name:       "Lead-Based Paint Disclosure",
			Required:   false,
			Condition:  func(p *models.Project, prop *models.Property, c *models.ContractTerms) bool { return prop != nil && prop.YearBuilt < 1978 },
		},
		{
			Type:       "residential_contract",
			FormNumber: "TREC 20-18",
			Name:       "One to Four Family Residential Contract",
			Required:   true,
			Condition:  func(p *models.Project, prop *models.Property, c *models.ContractTerms) bool { return c != nil },
		},
		{
			Type:       "third_party_financing",
			FormNumber: "TREC 40-9",
			Name:       "Third Party Financing Addendum",
			Required:   false,
			Condition:  func(p *models.Project, prop *models.Property, c *models.ContractTerms) bool { return c != nil && c.BuyerFinancing },
		},
		{
			Type:       "seller_lease",
			FormNumber: "TREC 15-6",
			Name:       "Seller's Temporary Residential Lease",
			Required:   false,
			Condition:  func(p *models.Project, prop *models.Property, c *models.ContractTerms) bool { return c != nil && c.SellerStaysPostClose },
		},
		{
			Type:       "hoa_addendum",
			FormNumber: "TREC 36-10",
			Name:       "Addendum for Property Subject to Mandatory Membership in HOA",
			Required:   false,
			Condition:  func(p *models.Project, prop *models.Property, c *models.ContractTerms) bool { return prop != nil && prop.HasHOA },
		},
	}
}

// GetApplicableDocuments returns documents that apply to this project
func (g *Generator) GetApplicableDocuments(project *models.Project, property *models.Property, contract *models.ContractTerms) []*models.Document {
	templates := g.AvailableTemplates()
	var docs []*models.Document

	for _, tmpl := range templates {
		if tmpl.Condition(project, property, contract) {
			doc := &models.Document{
				ProjectID:  project.ID,
				Type:       tmpl.Type,
				FormNumber: tmpl.FormNumber,
				Status:     "pending",
			}
			docs = append(docs, doc)
		}
	}

	return docs
}

// GeneratePDF generates a PDF for the given document
// For MVP, we'll generate simple text-based PDFs
// In production, this would use proper PDF templates
func (g *Generator) GeneratePDF(
	doc *models.Document,
	project *models.Project,
	property *models.Property,
	disclosure *models.Disclosure,
	contract *models.ContractTerms,
) ([]byte, error) {
	var buf bytes.Buffer

	switch doc.Type {
	case "sellers_disclosure":
		return g.generateSellersDisclosure(&buf, project, property, disclosure)
	case "lead_paint":
		return g.generateLeadPaintDisclosure(&buf, project, property)
	case "residential_contract":
		return g.generateResidentialContract(&buf, project, property, contract)
	case "third_party_financing":
		return g.generateFinancingAddendum(&buf, project, contract)
	case "seller_lease":
		return g.generateSellerLease(&buf, project, contract)
	case "hoa_addendum":
		return g.generateHOAAddendum(&buf, project, property)
	default:
		return nil, fmt.Errorf("unknown document type: %s", doc.Type)
	}
}

func (g *Generator) generateSellersDisclosure(buf *bytes.Buffer, project *models.Project, property *models.Property, disclosure *models.Disclosure) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "Letter", "")
	pdf.AddPage()

	// Title
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "SELLER'S DISCLOSURE NOTICE")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, "TREC Form OP-H")
	pdf.Ln(10)

	// Property Info
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 8, "Property Information")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(50, 6, "Property Address:")
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, project.PropertyAddress)
	pdf.Ln(6)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(50, 6, "Seller(s):")
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, strings.Join(project.SellerNames, ", "))
	pdf.Ln(10)

	if disclosure != nil {
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(0, 8, "Property Condition")
		pdf.Ln(8)

		pdf.SetFont("Arial", "", 10)
		pdf.Cell(50, 6, "Year Built:")
		pdf.Cell(0, 6, fmt.Sprintf("%d", property.YearBuilt))
		pdf.Ln(6)

		pdf.Cell(50, 6, "Homestead:")
		pdf.Cell(0, 6, boolToYesNo(property.IsHomestead))
		pdf.Ln(8)

		pdf.SetFont("Arial", "B", 11)
		pdf.Cell(0, 6, "Flooding & Water")
		pdf.Ln(6)
		pdf.SetFont("Arial", "", 10)

		pdf.Cell(50, 6, "Flooding History:")
		pdf.Cell(0, 6, boolToYesNo(disclosure.FloodingHistory))
		pdf.Ln(6)

		if disclosure.FloodingDetails != "" {
			pdf.Cell(50, 6, "Details:")
			pdf.MultiCell(0, 6, disclosure.FloodingDetails, "", "", false)
		}

		pdf.Cell(50, 6, "In Flood Plain:")
		pdf.Cell(0, 6, boolToYesNo(disclosure.InFloodPlain))
		pdf.Ln(8)

		pdf.SetFont("Arial", "B", 11)
		pdf.Cell(0, 6, "Foundation")
		pdf.Ln(6)
		pdf.SetFont("Arial", "", 10)

		pdf.Cell(50, 6, "Foundation Issues:")
		pdf.Cell(0, 6, boolToYesNo(disclosure.FoundationIssues))
		pdf.Ln(6)

		if disclosure.FoundationDetails != "" {
			pdf.Cell(50, 6, "Details:")
			pdf.MultiCell(0, 6, disclosure.FoundationDetails, "", "", false)
		}

		pdf.Ln(4)
		pdf.SetFont("Arial", "B", 11)
		pdf.Cell(0, 6, "Roof")
		pdf.Ln(6)
		pdf.SetFont("Arial", "", 10)

		pdf.Cell(50, 6, "Roof Age:")
		pdf.Cell(0, 6, fmt.Sprintf("%d years", disclosure.RoofAge))
		pdf.Ln(6)

		pdf.Cell(50, 6, "Roof Issues:")
		pdf.Cell(0, 6, boolToYesNo(disclosure.RoofIssues))
		pdf.Ln(6)
	}

	// Footer
	pdf.Ln(10)
	pdf.SetFont("Arial", "I", 9)
	pdf.Cell(0, 6, fmt.Sprintf("Generated: %s", time.Now().Format("January 2, 2006")))
	pdf.Ln(6)
	pdf.SetTextColor(128, 128, 128)
	pdf.MultiCell(0, 5, "Note: This is a simplified document for demonstration. In production, use official TREC forms.", "", "", false)

	err := pdf.Output(buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func boolToYesNo(b bool) string {
	if b {
		return "Yes"
	}
	return "No"
}

func (g *Generator) generateLeadPaintDisclosure(buf *bytes.Buffer, project *models.Project, property *models.Property) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "Letter", "")
	pdf.AddPage()

	// Title
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "LEAD-BASED PAINT DISCLOSURE")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, "TREC Form OP-L")
	pdf.Ln(12)

	// Property Info
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(50, 6, "Property Address:")
	pdf.Cell(0, 6, project.PropertyAddress)
	pdf.Ln(6)

	pdf.Cell(50, 6, "Year Built:")
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(0, 6, fmt.Sprintf("%d (Pre-1978)", property.YearBuilt))
	pdf.Ln(12)

	// Notice
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 8, "NOTICE")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 10)
	pdf.MultiCell(0, 6, "Properties built before 1978 may contain lead-based paint. Lead from paint, paint chips, and dust can pose health hazards if not managed properly.", "", "", false)
	pdf.Ln(4)
	pdf.MultiCell(0, 6, "Buyer has the right to conduct a lead-based paint inspection or risk assessment within the option period.", "", "", false)

	// Footer
	pdf.Ln(10)
	pdf.SetFont("Arial", "I", 9)
	pdf.Cell(0, 6, fmt.Sprintf("Generated: %s", time.Now().Format("January 2, 2006")))

	err := pdf.Output(buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (g *Generator) generateResidentialContract(buf *bytes.Buffer, project *models.Project, property *models.Property, contract *models.ContractTerms) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "Letter", "")
	pdf.AddPage()

	// Title
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "ONE TO FOUR FAMILY RESIDENTIAL CONTRACT")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, "TREC Form 20-18")
	pdf.Ln(12)

	// 1. PARTIES
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 8, "1. PARTIES")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(40, 6, "Seller(s):")
	pdf.Cell(0, 6, strings.Join(project.SellerNames, ", "))
	pdf.Ln(6)

	pdf.Cell(40, 6, "Buyer(s):")
	pdf.Cell(0, 6, "_______________________")
	pdf.Ln(10)

	// 2. PROPERTY
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 8, "2. PROPERTY")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(50, 6, "Address:")
	pdf.Cell(0, 6, project.PropertyAddress)
	pdf.Ln(6)

	pdf.Cell(50, 6, "Legal Description:")
	pdf.MultiCell(0, 6, property.LegalDescription, "", "", false)

	pdf.Cell(50, 6, "Tax ID:")
	pdf.Cell(0, 6, property.TaxID)
	pdf.Ln(10)

	// 3. SALES PRICE
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 8, "3. SALES PRICE")
	pdf.Ln(8)

	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(0, 6, fmt.Sprintf("$%.2f", contract.SalesPrice))
	pdf.Ln(10)

	// 4. KEY DATES
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 8, "4. KEY DATES")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(50, 6, "Effective Date:")
	pdf.Cell(0, 6, contract.EffectiveDate.Format("January 2, 2006"))
	pdf.Ln(6)

	pdf.Cell(50, 6, "Closing Date:")
	pdf.Cell(0, 6, contract.ClosingDate.Format("January 2, 2006"))
	pdf.Ln(10)

	// 5. EARNEST MONEY
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 8, "5. EARNEST MONEY")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, fmt.Sprintf("$%.2f", contract.EarnestMoney))
	pdf.Ln(10)

	// 6. OPTION PERIOD
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 8, "6. OPTION PERIOD")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(50, 6, "Option Fee:")
	pdf.Cell(0, 6, fmt.Sprintf("$%.2f", contract.OptionFee))
	pdf.Ln(6)

	pdf.Cell(50, 6, "Option Period:")
	pdf.Cell(0, 6, fmt.Sprintf("%d days", contract.OptionPeriodDays))
	pdf.Ln(10)

	// 7. TITLE
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 8, "7. TITLE")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, fmt.Sprintf("Title Commitment Due: %d days from effective date", contract.TitleCommitmentDays))
	pdf.Ln(10)

	// 8. EXCLUSIONS (if any)
	if len(contract.ExcludedItems) > 0 {
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(0, 8, "8. EXCLUSIONS")
		pdf.Ln(8)

		pdf.SetFont("Arial", "", 10)
		for _, item := range contract.ExcludedItems {
			pdf.Cell(10, 6, "")
			pdf.Cell(0, 6, "- "+item)
			pdf.Ln(6)
		}
		pdf.Ln(4)
	}

	// 9. INCLUDED PERSONAL PROPERTY (if any)
	if len(contract.IncludedPersonalItems) > 0 {
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(0, 8, "9. INCLUDED PERSONAL PROPERTY")
		pdf.Ln(8)

		pdf.SetFont("Arial", "", 10)
		for _, item := range contract.IncludedPersonalItems {
			pdf.Cell(10, 6, "")
			pdf.Cell(0, 6, "- "+item)
			pdf.Ln(6)
		}
		pdf.Ln(4)
	}

	// Footer
	pdf.Ln(10)
	pdf.SetFont("Arial", "I", 9)
	pdf.Cell(0, 6, fmt.Sprintf("Generated: %s", time.Now().Format("January 2, 2006")))
	pdf.Ln(6)
	pdf.SetTextColor(128, 128, 128)
	pdf.MultiCell(0, 5, "Note: This is a simplified document for demonstration. In production, use official TREC Form 20-18.", "", "", false)

	err := pdf.Output(buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (g *Generator) generateFinancingAddendum(buf *bytes.Buffer, project *models.Project, contract *models.ContractTerms) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "Letter", "")
	pdf.AddPage()

	// Title
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "THIRD PARTY FINANCING ADDENDUM")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, "TREC Form 40-9")
	pdf.Ln(12)

	// Property Info
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(50, 6, "Property:")
	pdf.Cell(0, 6, project.PropertyAddress)
	pdf.Ln(8)

	pdf.Cell(50, 6, "Sales Price:")
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(0, 6, fmt.Sprintf("$%.2f", contract.SalesPrice))
	pdf.Ln(12)

	// Financing terms
	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(0, 8, "Financing Terms")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 10)
	pdf.MultiCell(0, 6, "Buyer shall apply for third party financing within ___ days after the effective date of this contract.", "", "", false)
	pdf.Ln(4)
	pdf.MultiCell(0, 6, "Buyer approval for financing required within ___ days after the effective date.", "", "", false)
	pdf.Ln(6)

	pdf.SetFont("Arial", "I", 9)
	pdf.MultiCell(0, 5, "Note: Specific financing terms and deadlines to be negotiated between parties.", "", "", false)

	// Footer
	pdf.Ln(10)
	pdf.SetFont("Arial", "I", 9)
	pdf.Cell(0, 6, fmt.Sprintf("Generated: %s", time.Now().Format("January 2, 2006")))
	pdf.Ln(6)
	pdf.SetTextColor(128, 128, 128)
	pdf.MultiCell(0, 5, "Note: This is a simplified document for demonstration. In production, use official TREC Form 40-9.", "", "", false)

	err := pdf.Output(buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (g *Generator) generateSellerLease(buf *bytes.Buffer, project *models.Project, contract *models.ContractTerms) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "Letter", "")
	pdf.AddPage()

	// Title
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "SELLER'S TEMPORARY RESIDENTIAL LEASE")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, "TREC Form 15-6")
	pdf.Ln(12)

	// Property Info
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(50, 6, "Property:")
	pdf.Cell(0, 6, project.PropertyAddress)
	pdf.Ln(10)

	// Parties
	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(0, 8, "Parties")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(50, 6, "Landlord (Buyer):")
	pdf.Cell(0, 6, "_______________________")
	pdf.Ln(6)

	pdf.Cell(50, 6, "Tenant (Seller):")
	pdf.Cell(0, 6, strings.Join(project.SellerNames, ", "))
	pdf.Ln(10)

	// Lease Terms
	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(0, 8, "Lease Terms")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(50, 6, "Closing Date:")
	pdf.Cell(0, 6, contract.ClosingDate.Format("January 2, 2006"))
	pdf.Ln(6)

	pdf.Cell(50, 6, "Lease Period:")
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(0, 6, fmt.Sprintf("%d days from closing", contract.SellerLeaseDays))
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 10)
	leaseEndDate := contract.ClosingDate.AddDate(0, 0, contract.SellerLeaseDays)
	pdf.Cell(50, 6, "Lease End Date:")
	pdf.Cell(0, 6, leaseEndDate.Format("January 2, 2006"))
	pdf.Ln(10)

	// Additional Terms
	pdf.SetFont("Arial", "I", 9)
	pdf.MultiCell(0, 5, "Seller agrees to vacate the property by the lease end date. Additional terms and rental amount to be negotiated between parties.", "", "", false)

	// Footer
	pdf.Ln(10)
	pdf.SetFont("Arial", "I", 9)
	pdf.Cell(0, 6, fmt.Sprintf("Generated: %s", time.Now().Format("January 2, 2006")))
	pdf.Ln(6)
	pdf.SetTextColor(128, 128, 128)
	pdf.MultiCell(0, 5, "Note: This is a simplified document for demonstration. In production, use official TREC Form 15-6.", "", "", false)

	err := pdf.Output(buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (g *Generator) generateHOAAddendum(buf *bytes.Buffer, project *models.Project, property *models.Property) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "Letter", "")
	pdf.AddPage()

	// Title
	pdf.SetFont("Arial", "B", 14)
	pdf.MultiCell(0, 8, "ADDENDUM FOR PROPERTY SUBJECT TO MANDATORY MEMBERSHIP IN HOMEOWNERS ASSOCIATION", "", "", false)
	pdf.Ln(6)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, "TREC Form 36-10")
	pdf.Ln(12)

	// Property Info
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(50, 6, "Property:")
	pdf.Cell(0, 6, project.PropertyAddress)
	pdf.Ln(10)

	// HOA Information
	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(0, 8, "Homeowners Association Information")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(50, 6, "HOA Name:")
	pdf.Cell(0, 6, property.HOAName)
	pdf.Ln(6)

	pdf.Cell(50, 6, "HOA Fee:")
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(0, 6, fmt.Sprintf("$%.2f (%s)", property.HOAFee, property.HOAFrequency))
	pdf.Ln(10)

	// Buyer Acknowledgment
	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(0, 8, "Buyer Acknowledgment")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 10)
	pdf.MultiCell(0, 6, "Buyer acknowledges that the property is subject to mandatory membership in the homeowners association and is bound by the HOA's governing documents, including:", "", "", false)
	pdf.Ln(4)

	pdf.Cell(10, 6, "")
	pdf.Cell(0, 6, "- Declaration of Covenants, Conditions, and Restrictions (CC&Rs)")
	pdf.Ln(6)

	pdf.Cell(10, 6, "")
	pdf.Cell(0, 6, "- Bylaws")
	pdf.Ln(6)

	pdf.Cell(10, 6, "")
	pdf.Cell(0, 6, "- Rules and Regulations")
	pdf.Ln(6)

	pdf.Cell(10, 6, "")
	pdf.Cell(0, 6, "- Articles of Incorporation")
	pdf.Ln(10)

	pdf.SetFont("Arial", "I", 9)
	pdf.MultiCell(0, 5, "Buyer should review all HOA documents, including financial statements and meeting minutes, during the option period.", "", "", false)

	// Footer
	pdf.Ln(10)
	pdf.SetFont("Arial", "I", 9)
	pdf.Cell(0, 6, fmt.Sprintf("Generated: %s", time.Now().Format("January 2, 2006")))
	pdf.Ln(6)
	pdf.SetTextColor(128, 128, 128)
	pdf.MultiCell(0, 5, "Note: This is a simplified document for demonstration. In production, use official TREC Form 36-10.", "", "", false)

	err := pdf.Output(buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
