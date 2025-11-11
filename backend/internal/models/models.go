package models

import "time"

// Project represents a seller's transaction workflow
type Project struct {
	ID                string    `json:"id"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	PropertyAddress   string    `json:"property_address"`
	SellerNames       []string  `json:"seller_names"`
	SellerEmail       string    `json:"seller_email"`
	SellerPhone       string    `json:"seller_phone"`
	HasAgent          bool      `json:"has_agent"`
	AgentName         string    `json:"agent_name,omitempty"`
	TitleCompany      string    `json:"title_company,omitempty"`
	TargetListDate    *time.Time `json:"target_list_date,omitempty"`
	Status            string    `json:"status"` // "setup", "listing_prep", "under_contract", "closing"
}

// Property holds property-specific details
type Property struct {
	ID               string    `json:"id"`
	ProjectID        string    `json:"project_id"`
	YearBuilt        int       `json:"year_built"`
	LegalDescription string    `json:"legal_description"`
	TaxID            string    `json:"tax_id"`
	IsHomestead      bool      `json:"is_homestead"`
	HasHOA           bool      `json:"has_hoa"`
	HOAName          string    `json:"hoa_name,omitempty"`
	HOAFee           float64   `json:"hoa_fee,omitempty"`
	HOAFrequency     string    `json:"hoa_frequency,omitempty"` // "monthly", "quarterly", "annual"
	HasSurvey        bool      `json:"has_survey"`
	SurveyDate       *time.Time `json:"survey_date,omitempty"`
	SurveyConditions string    `json:"survey_conditions,omitempty"`
}

// Disclosure represents seller's disclosure responses
type Disclosure struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"project_id"`
	UpdatedAt time.Time `json:"updated_at"`

	// Flooding & Water
	FloodingHistory        bool   `json:"flooding_history"`
	FloodingDetails        string `json:"flooding_details,omitempty"`
	InFloodPlain           bool   `json:"in_flood_plain"`
	FloodInsuranceRequired bool   `json:"flood_insurance_required"`

	// Structure & Systems
	FoundationIssues       bool   `json:"foundation_issues"`
	FoundationDetails      string `json:"foundation_details,omitempty"`
	RoofAge                int    `json:"roof_age"`
	RoofIssues             bool   `json:"roof_issues"`
	RoofDetails            string `json:"roof_details,omitempty"`
	PlumbingIssues         bool   `json:"plumbing_issues"`
	PlumbingDetails        string `json:"plumbing_details,omitempty"`
	ElectricalIssues       bool   `json:"electrical_issues"`
	ElectricalDetails      string `json:"electrical_details,omitempty"`
	HVACAge                int    `json:"hvac_age"`
	HVACIssues             bool   `json:"hvac_issues"`
	HVACDetails            string `json:"hvac_details,omitempty"`

	// Other Disclosures
	LeadBasedPaint         bool   `json:"lead_based_paint"` // auto-set if year < 1978
	InsuranceClaimsHistory bool   `json:"insurance_claims_history"`
	ClaimsDetails          string `json:"claims_details,omitempty"`
	OtherMaterialDefects   string `json:"other_material_defects,omitempty"`
}

// ContractTerms holds the agreed contract details
type ContractTerms struct {
	ID                    string    `json:"id"`
	ProjectID             string    `json:"project_id"`
	EffectiveDate         time.Time `json:"effective_date"`
	SalesPrice            float64   `json:"sales_price"`
	ClosingDate           time.Time `json:"closing_date"`
	OptionFee             float64   `json:"option_fee"`
	OptionPeriodDays      int       `json:"option_period_days"`
	EarnestMoney          float64   `json:"earnest_money"`
	TitleCommitmentDays   int       `json:"title_commitment_days"` // default 20
	BuyerFinancing        bool      `json:"buyer_financing"`
	SellerStaysPostClose  bool      `json:"seller_stays_post_close"`
	SellerLeaseDays       int       `json:"seller_lease_days,omitempty"`
	ExcludedItems         []string  `json:"excluded_items,omitempty"`
	IncludedPersonalItems []string  `json:"included_personal_items,omitempty"`
}

// Deadline represents a calculated deadline with reminders
type Deadline struct {
	ID            string    `json:"id"`
	ProjectID     string    `json:"project_id"`
	Type          string    `json:"type"` // "option_fee", "earnest_money", "title_commitment", "inspection_end", etc.
	Description   string    `json:"description"`
	DueDate       time.Time `json:"due_date"`
	IsBusinessDay bool      `json:"is_business_day"` // whether calculated using business days
	Priority      string    `json:"priority"` // "critical", "high", "normal"
	Completed     bool      `json:"completed"`
	Notes         string    `json:"notes,omitempty"`
}

// Document represents a generated or required document
type Document struct {
	ID          string    `json:"id"`
	ProjectID   string    `json:"project_id"`
	Type        string    `json:"type"` // "sellers_disclosure", "op_l", "hoa_addendum", etc.
	FormNumber  string    `json:"form_number"` // e.g., "TREC 20-18", "OP-L"
	Status      string    `json:"status"` // "pending", "draft", "ready", "delivered"
	GeneratedAt *time.Time `json:"generated_at,omitempty"`
	FilePath    string    `json:"file_path,omitempty"`
}

// Service represents a type of service offered in the marketplace
type Service struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Category         string    `json:"category"` // "inspections", "contracts_title_legal", "repairs_trades"
	Description      string    `json:"description"`
	TypicalTimeline  string    `json:"typical_timeline,omitempty"`
	EstimatedCostMin float64   `json:"estimated_cost_min,omitempty"`
	EstimatedCostMax float64   `json:"estimated_cost_max,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// Provider represents a business or individual offering a service
type Provider struct {
	ID                 string    `json:"id"`
	ServiceID          string    `json:"service_id"`
	BusinessName       string    `json:"business_name"`
	ContactName        string    `json:"contact_name,omitempty"`
	Email              string    `json:"email"`
	Phone              string    `json:"phone"`
	Address            string    `json:"address,omitempty"`
	City               string    `json:"city,omitempty"`
	State              string    `json:"state,omitempty"`
	Zip                string    `json:"zip,omitempty"`
	Bio                string    `json:"bio,omitempty"`
	YearsExperience    int       `json:"years_experience,omitempty"`
	LicenseNumber      string    `json:"license_number,omitempty"`
	InsuranceVerified  bool      `json:"insurance_verified"`
	AvailabilityStatus string    `json:"availability_status"` // "available", "limited", "unavailable"
	RatingAverage      float64   `json:"rating_average"`
	RatingCount        int       `json:"rating_count"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// ServiceRequest represents a user's request for a service
type ServiceRequest struct {
	ID              string    `json:"id"`
	ProjectID       string    `json:"project_id,omitempty"`
	UserEmail       string    `json:"user_email"`
	UserName        string    `json:"user_name"`
	UserPhone       string    `json:"user_phone,omitempty"`
	ServiceID       string    `json:"service_id"`
	ProviderID      string    `json:"provider_id,omitempty"`
	PropertyAddress string    `json:"property_address,omitempty"`
	RequestedDate   *time.Time `json:"requested_date,omitempty"`
	PreferredTime   string    `json:"preferred_time,omitempty"`
	Status          string    `json:"status"` // "pending", "matched", "scheduled", "completed", "cancelled"
	Notes           string    `json:"notes,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ProviderReview represents a user review of a provider
type ProviderReview struct {
	ID               string    `json:"id"`
	ProviderID       string    `json:"provider_id"`
	ServiceRequestID string    `json:"service_request_id"`
	UserEmail        string    `json:"user_email"`
	Rating           int       `json:"rating"` // 1-5 stars
	ReviewText       string    `json:"review_text,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
}
