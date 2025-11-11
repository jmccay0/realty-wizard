package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jessicaandtommymccay/realty-wizard/backend/internal/models"
	_ "modernc.org/sqlite"
)

// SQLiteStorage implements Storage using SQLite
type SQLiteStorage struct {
	db *sql.DB
}

// NewSQLiteStorage creates a new SQLite storage instance
func NewSQLiteStorage(dbPath string) (*SQLiteStorage, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	storage := &SQLiteStorage{db: db}
	if err := storage.migrate(); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return storage, nil
}

// migrate creates the database schema
func (s *SQLiteStorage) migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS projects (
		id TEXT PRIMARY KEY,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		property_address TEXT NOT NULL,
		seller_names TEXT NOT NULL,
		seller_email TEXT,
		seller_phone TEXT,
		has_agent BOOLEAN,
		agent_name TEXT,
		title_company TEXT,
		target_list_date DATETIME,
		status TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS properties (
		id TEXT PRIMARY KEY,
		project_id TEXT NOT NULL,
		year_built INTEGER,
		legal_description TEXT,
		tax_id TEXT,
		is_homestead BOOLEAN,
		has_hoa BOOLEAN,
		hoa_name TEXT,
		hoa_fee REAL,
		hoa_frequency TEXT,
		has_survey BOOLEAN,
		survey_date DATETIME,
		survey_conditions TEXT,
		FOREIGN KEY (project_id) REFERENCES projects(id)
	);

	CREATE TABLE IF NOT EXISTS disclosures (
		id TEXT PRIMARY KEY,
		project_id TEXT NOT NULL,
		updated_at DATETIME NOT NULL,
		flooding_history BOOLEAN,
		flooding_details TEXT,
		in_flood_plain BOOLEAN,
		flood_insurance_required BOOLEAN,
		foundation_issues BOOLEAN,
		foundation_details TEXT,
		roof_age INTEGER,
		roof_issues BOOLEAN,
		roof_details TEXT,
		plumbing_issues BOOLEAN,
		plumbing_details TEXT,
		electrical_issues BOOLEAN,
		electrical_details TEXT,
		hvac_age INTEGER,
		hvac_issues BOOLEAN,
		hvac_details TEXT,
		lead_based_paint BOOLEAN,
		insurance_claims_history BOOLEAN,
		claims_details TEXT,
		other_material_defects TEXT,
		FOREIGN KEY (project_id) REFERENCES projects(id)
	);

	CREATE TABLE IF NOT EXISTS contracts (
		id TEXT PRIMARY KEY,
		project_id TEXT NOT NULL,
		effective_date DATETIME NOT NULL,
		sales_price REAL NOT NULL,
		closing_date DATETIME NOT NULL,
		option_fee REAL,
		option_period_days INTEGER,
		earnest_money REAL,
		title_commitment_days INTEGER,
		buyer_financing BOOLEAN,
		seller_stays_post_close BOOLEAN,
		seller_lease_days INTEGER,
		excluded_items TEXT,
		included_personal_items TEXT,
		FOREIGN KEY (project_id) REFERENCES projects(id)
	);

	CREATE TABLE IF NOT EXISTS deadlines (
		id TEXT PRIMARY KEY,
		project_id TEXT NOT NULL,
		type TEXT NOT NULL,
		description TEXT NOT NULL,
		due_date DATETIME NOT NULL,
		is_business_day BOOLEAN,
		priority TEXT,
		completed BOOLEAN,
		notes TEXT,
		FOREIGN KEY (project_id) REFERENCES projects(id)
	);

	CREATE TABLE IF NOT EXISTS documents (
		id TEXT PRIMARY KEY,
		project_id TEXT NOT NULL,
		type TEXT NOT NULL,
		form_number TEXT,
		status TEXT NOT NULL,
		generated_at DATETIME,
		file_path TEXT,
		FOREIGN KEY (project_id) REFERENCES projects(id)
	);

	CREATE INDEX IF NOT EXISTS idx_properties_project ON properties(project_id);
	CREATE INDEX IF NOT EXISTS idx_disclosures_project ON disclosures(project_id);
	CREATE INDEX IF NOT EXISTS idx_contracts_project ON contracts(project_id);
	CREATE INDEX IF NOT EXISTS idx_deadlines_project ON deadlines(project_id);
	CREATE INDEX IF NOT EXISTS idx_documents_project ON documents(project_id);

	-- Service Marketplace Tables
	CREATE TABLE IF NOT EXISTS services (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		category TEXT NOT NULL,
		description TEXT NOT NULL,
		typical_timeline TEXT,
		estimated_cost_min REAL,
		estimated_cost_max REAL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);

	CREATE TABLE IF NOT EXISTS providers (
		id TEXT PRIMARY KEY,
		service_id TEXT NOT NULL,
		business_name TEXT NOT NULL,
		contact_name TEXT,
		email TEXT NOT NULL,
		phone TEXT NOT NULL,
		address TEXT,
		city TEXT,
		state TEXT,
		zip TEXT,
		bio TEXT,
		years_experience INTEGER,
		license_number TEXT,
		insurance_verified BOOLEAN DEFAULT FALSE,
		availability_status TEXT DEFAULT 'available',
		rating_average REAL DEFAULT 0,
		rating_count INTEGER DEFAULT 0,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY (service_id) REFERENCES services(id)
	);

	CREATE TABLE IF NOT EXISTS service_requests (
		id TEXT PRIMARY KEY,
		project_id TEXT,
		user_email TEXT NOT NULL,
		user_name TEXT NOT NULL,
		user_phone TEXT,
		service_id TEXT NOT NULL,
		provider_id TEXT,
		property_address TEXT,
		requested_date DATETIME,
		preferred_time TEXT,
		status TEXT NOT NULL DEFAULT 'pending',
		notes TEXT,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY (project_id) REFERENCES projects(id),
		FOREIGN KEY (service_id) REFERENCES services(id),
		FOREIGN KEY (provider_id) REFERENCES providers(id)
	);

	CREATE TABLE IF NOT EXISTS provider_reviews (
		id TEXT PRIMARY KEY,
		provider_id TEXT NOT NULL,
		service_request_id TEXT NOT NULL,
		user_email TEXT NOT NULL,
		rating INTEGER NOT NULL,
		review_text TEXT,
		created_at DATETIME NOT NULL,
		FOREIGN KEY (provider_id) REFERENCES providers(id),
		FOREIGN KEY (service_request_id) REFERENCES service_requests(id)
	);

	CREATE INDEX IF NOT EXISTS idx_providers_service ON providers(service_id);
	CREATE INDEX IF NOT EXISTS idx_service_requests_project ON service_requests(project_id);
	CREATE INDEX IF NOT EXISTS idx_service_requests_service ON service_requests(service_id);
	CREATE INDEX IF NOT EXISTS idx_service_requests_provider ON service_requests(provider_id);
	CREATE INDEX IF NOT EXISTS idx_provider_reviews_provider ON provider_reviews(provider_id);
	`

	_, err := s.db.Exec(schema)
	return err
}

// CreateProject inserts a new project
func (s *SQLiteStorage) CreateProject(p *models.Project) error {
	sellerNamesJSON, _ := json.Marshal(p.SellerNames)
	_, err := s.db.Exec(`
		INSERT INTO projects (id, created_at, updated_at, property_address, seller_names,
			seller_email, seller_phone, has_agent, agent_name, title_company, target_list_date, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		p.ID, p.CreatedAt, p.UpdatedAt, p.PropertyAddress, sellerNamesJSON,
		p.SellerEmail, p.SellerPhone, p.HasAgent, p.AgentName, p.TitleCompany, p.TargetListDate, p.Status)
	return err
}

// GetProject retrieves a project by ID
func (s *SQLiteStorage) GetProject(id string) (*models.Project, error) {
	var p models.Project
	var sellerNamesJSON string
	var targetListDate sql.NullTime

	err := s.db.QueryRow(`
		SELECT id, created_at, updated_at, property_address, seller_names,
			seller_email, seller_phone, has_agent, agent_name, title_company, target_list_date, status
		FROM projects WHERE id = ?`, id).Scan(
		&p.ID, &p.CreatedAt, &p.UpdatedAt, &p.PropertyAddress, &sellerNamesJSON,
		&p.SellerEmail, &p.SellerPhone, &p.HasAgent, &p.AgentName, &p.TitleCompany, &targetListDate, &p.Status)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("project not found")
		}
		return nil, err
	}

	json.Unmarshal([]byte(sellerNamesJSON), &p.SellerNames)
	if targetListDate.Valid {
		p.TargetListDate = &targetListDate.Time
	}

	return &p, nil
}

// UpdateProject updates an existing project
func (s *SQLiteStorage) UpdateProject(p *models.Project) error {
	p.UpdatedAt = time.Now()
	sellerNamesJSON, _ := json.Marshal(p.SellerNames)
	_, err := s.db.Exec(`
		UPDATE projects SET updated_at = ?, property_address = ?, seller_names = ?,
			seller_email = ?, seller_phone = ?, has_agent = ?, agent_name = ?,
			title_company = ?, target_list_date = ?, status = ?
		WHERE id = ?`,
		p.UpdatedAt, p.PropertyAddress, sellerNamesJSON, p.SellerEmail, p.SellerPhone,
		p.HasAgent, p.AgentName, p.TitleCompany, p.TargetListDate, p.Status, p.ID)
	return err
}

// ListProjects returns all projects
func (s *SQLiteStorage) ListProjects() ([]*models.Project, error) {
	rows, err := s.db.Query(`
		SELECT id, created_at, updated_at, property_address, seller_names,
			seller_email, seller_phone, has_agent, agent_name, title_company, target_list_date, status
		FROM projects ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*models.Project
	for rows.Next() {
		var p models.Project
		var sellerNamesJSON string
		var targetListDate sql.NullTime

		err := rows.Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt, &p.PropertyAddress, &sellerNamesJSON,
			&p.SellerEmail, &p.SellerPhone, &p.HasAgent, &p.AgentName, &p.TitleCompany, &targetListDate, &p.Status)
		if err != nil {
			return nil, err
		}

		json.Unmarshal([]byte(sellerNamesJSON), &p.SellerNames)
		if targetListDate.Valid {
			p.TargetListDate = &targetListDate.Time
		}

		projects = append(projects, &p)
	}

	return projects, nil
}

// CreateProperty inserts a new property
func (s *SQLiteStorage) CreateProperty(p *models.Property) error {
	_, err := s.db.Exec(`
		INSERT INTO properties (id, project_id, year_built, legal_description, tax_id,
			is_homestead, has_hoa, hoa_name, hoa_fee, hoa_frequency, has_survey, survey_date, survey_conditions)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		p.ID, p.ProjectID, p.YearBuilt, p.LegalDescription, p.TaxID,
		p.IsHomestead, p.HasHOA, p.HOAName, p.HOAFee, p.HOAFrequency, p.HasSurvey, p.SurveyDate, p.SurveyConditions)
	return err
}

// GetProperty retrieves a property by project ID
func (s *SQLiteStorage) GetProperty(projectID string) (*models.Property, error) {
	var p models.Property
	var surveyDate sql.NullTime

	err := s.db.QueryRow(`
		SELECT id, project_id, year_built, legal_description, tax_id,
			is_homestead, has_hoa, hoa_name, hoa_fee, hoa_frequency, has_survey, survey_date, survey_conditions
		FROM properties WHERE project_id = ?`, projectID).Scan(
		&p.ID, &p.ProjectID, &p.YearBuilt, &p.LegalDescription, &p.TaxID,
		&p.IsHomestead, &p.HasHOA, &p.HOAName, &p.HOAFee, &p.HOAFrequency, &p.HasSurvey, &surveyDate, &p.SurveyConditions)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("property not found")
		}
		return nil, err
	}

	if surveyDate.Valid {
		p.SurveyDate = &surveyDate.Time
	}

	return &p, nil
}

// UpdateProperty updates an existing property
func (s *SQLiteStorage) UpdateProperty(p *models.Property) error {
	_, err := s.db.Exec(`
		UPDATE properties SET year_built = ?, legal_description = ?, tax_id = ?,
			is_homestead = ?, has_hoa = ?, hoa_name = ?, hoa_fee = ?, hoa_frequency = ?,
			has_survey = ?, survey_date = ?, survey_conditions = ?
		WHERE id = ?`,
		p.YearBuilt, p.LegalDescription, p.TaxID, p.IsHomestead, p.HasHOA, p.HOAName,
		p.HOAFee, p.HOAFrequency, p.HasSurvey, p.SurveyDate, p.SurveyConditions, p.ID)
	return err
}

// CreateDisclosure inserts a new disclosure
func (s *SQLiteStorage) CreateDisclosure(d *models.Disclosure) error {
	_, err := s.db.Exec(`
		INSERT INTO disclosures (id, project_id, updated_at, flooding_history, flooding_details,
			in_flood_plain, flood_insurance_required, foundation_issues, foundation_details,
			roof_age, roof_issues, roof_details, plumbing_issues, plumbing_details,
			electrical_issues, electrical_details, hvac_age, hvac_issues, hvac_details,
			lead_based_paint, insurance_claims_history, claims_details, other_material_defects)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		d.ID, d.ProjectID, d.UpdatedAt, d.FloodingHistory, d.FloodingDetails,
		d.InFloodPlain, d.FloodInsuranceRequired, d.FoundationIssues, d.FoundationDetails,
		d.RoofAge, d.RoofIssues, d.RoofDetails, d.PlumbingIssues, d.PlumbingDetails,
		d.ElectricalIssues, d.ElectricalDetails, d.HVACAge, d.HVACIssues, d.HVACDetails,
		d.LeadBasedPaint, d.InsuranceClaimsHistory, d.ClaimsDetails, d.OtherMaterialDefects)
	return err
}

// GetDisclosure retrieves a disclosure by project ID
func (s *SQLiteStorage) GetDisclosure(projectID string) (*models.Disclosure, error) {
	var d models.Disclosure
	err := s.db.QueryRow(`
		SELECT id, project_id, updated_at, flooding_history, flooding_details,
			in_flood_plain, flood_insurance_required, foundation_issues, foundation_details,
			roof_age, roof_issues, roof_details, plumbing_issues, plumbing_details,
			electrical_issues, electrical_details, hvac_age, hvac_issues, hvac_details,
			lead_based_paint, insurance_claims_history, claims_details, other_material_defects
		FROM disclosures WHERE project_id = ?`, projectID).Scan(
		&d.ID, &d.ProjectID, &d.UpdatedAt, &d.FloodingHistory, &d.FloodingDetails,
		&d.InFloodPlain, &d.FloodInsuranceRequired, &d.FoundationIssues, &d.FoundationDetails,
		&d.RoofAge, &d.RoofIssues, &d.RoofDetails, &d.PlumbingIssues, &d.PlumbingDetails,
		&d.ElectricalIssues, &d.ElectricalDetails, &d.HVACAge, &d.HVACIssues, &d.HVACDetails,
		&d.LeadBasedPaint, &d.InsuranceClaimsHistory, &d.ClaimsDetails, &d.OtherMaterialDefects)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("disclosure not found")
		}
		return nil, err
	}

	return &d, nil
}

// UpdateDisclosure updates an existing disclosure
func (s *SQLiteStorage) UpdateDisclosure(d *models.Disclosure) error {
	d.UpdatedAt = time.Now()
	_, err := s.db.Exec(`
		UPDATE disclosures SET updated_at = ?, flooding_history = ?, flooding_details = ?,
			in_flood_plain = ?, flood_insurance_required = ?, foundation_issues = ?, foundation_details = ?,
			roof_age = ?, roof_issues = ?, roof_details = ?, plumbing_issues = ?, plumbing_details = ?,
			electrical_issues = ?, electrical_details = ?, hvac_age = ?, hvac_issues = ?, hvac_details = ?,
			lead_based_paint = ?, insurance_claims_history = ?, claims_details = ?, other_material_defects = ?
		WHERE id = ?`,
		d.UpdatedAt, d.FloodingHistory, d.FloodingDetails, d.InFloodPlain, d.FloodInsuranceRequired,
		d.FoundationIssues, d.FoundationDetails, d.RoofAge, d.RoofIssues, d.RoofDetails,
		d.PlumbingIssues, d.PlumbingDetails, d.ElectricalIssues, d.ElectricalDetails,
		d.HVACAge, d.HVACIssues, d.HVACDetails, d.LeadBasedPaint, d.InsuranceClaimsHistory,
		d.ClaimsDetails, d.OtherMaterialDefects, d.ID)
	return err
}

// CreateContract inserts a new contract
func (s *SQLiteStorage) CreateContract(c *models.ContractTerms) error {
	excludedJSON, _ := json.Marshal(c.ExcludedItems)
	includedJSON, _ := json.Marshal(c.IncludedPersonalItems)

	_, err := s.db.Exec(`
		INSERT INTO contracts (id, project_id, effective_date, sales_price, closing_date,
			option_fee, option_period_days, earnest_money, title_commitment_days, buyer_financing,
			seller_stays_post_close, seller_lease_days, excluded_items, included_personal_items)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		c.ID, c.ProjectID, c.EffectiveDate, c.SalesPrice, c.ClosingDate,
		c.OptionFee, c.OptionPeriodDays, c.EarnestMoney, c.TitleCommitmentDays, c.BuyerFinancing,
		c.SellerStaysPostClose, c.SellerLeaseDays, excludedJSON, includedJSON)
	return err
}

// GetContract retrieves a contract by project ID
func (s *SQLiteStorage) GetContract(projectID string) (*models.ContractTerms, error) {
	var c models.ContractTerms
	var excludedJSON, includedJSON string

	err := s.db.QueryRow(`
		SELECT id, project_id, effective_date, sales_price, closing_date,
			option_fee, option_period_days, earnest_money, title_commitment_days, buyer_financing,
			seller_stays_post_close, seller_lease_days, excluded_items, included_personal_items
		FROM contracts WHERE project_id = ?`, projectID).Scan(
		&c.ID, &c.ProjectID, &c.EffectiveDate, &c.SalesPrice, &c.ClosingDate,
		&c.OptionFee, &c.OptionPeriodDays, &c.EarnestMoney, &c.TitleCommitmentDays, &c.BuyerFinancing,
		&c.SellerStaysPostClose, &c.SellerLeaseDays, &excludedJSON, &includedJSON)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("contract not found")
		}
		return nil, err
	}

	json.Unmarshal([]byte(excludedJSON), &c.ExcludedItems)
	json.Unmarshal([]byte(includedJSON), &c.IncludedPersonalItems)

	return &c, nil
}

// UpdateContract updates an existing contract
func (s *SQLiteStorage) UpdateContract(c *models.ContractTerms) error {
	excludedJSON, _ := json.Marshal(c.ExcludedItems)
	includedJSON, _ := json.Marshal(c.IncludedPersonalItems)

	_, err := s.db.Exec(`
		UPDATE contracts SET effective_date = ?, sales_price = ?, closing_date = ?,
			option_fee = ?, option_period_days = ?, earnest_money = ?, title_commitment_days = ?,
			buyer_financing = ?, seller_stays_post_close = ?, seller_lease_days = ?,
			excluded_items = ?, included_personal_items = ?
		WHERE id = ?`,
		c.EffectiveDate, c.SalesPrice, c.ClosingDate, c.OptionFee, c.OptionPeriodDays,
		c.EarnestMoney, c.TitleCommitmentDays, c.BuyerFinancing, c.SellerStaysPostClose,
		c.SellerLeaseDays, excludedJSON, includedJSON, c.ID)
	return err
}

// CreateDeadline inserts a new deadline
func (s *SQLiteStorage) CreateDeadline(d *models.Deadline) error {
	_, err := s.db.Exec(`
		INSERT INTO deadlines (id, project_id, type, description, due_date, is_business_day, priority, completed, notes)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		d.ID, d.ProjectID, d.Type, d.Description, d.DueDate, d.IsBusinessDay, d.Priority, d.Completed, d.Notes)
	return err
}

// GetDeadline retrieves a deadline by ID
func (s *SQLiteStorage) GetDeadline(id string) (*models.Deadline, error) {
	var d models.Deadline
	err := s.db.QueryRow(`
		SELECT id, project_id, type, description, due_date, is_business_day, priority, completed, notes
		FROM deadlines WHERE id = ?`, id).Scan(
		&d.ID, &d.ProjectID, &d.Type, &d.Description, &d.DueDate, &d.IsBusinessDay, &d.Priority, &d.Completed, &d.Notes)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("deadline not found")
		}
		return nil, err
	}

	return &d, nil
}

// ListDeadlines retrieves all deadlines for a project
func (s *SQLiteStorage) ListDeadlines(projectID string) ([]*models.Deadline, error) {
	rows, err := s.db.Query(`
		SELECT id, project_id, type, description, due_date, is_business_day, priority, completed, notes
		FROM deadlines WHERE project_id = ? ORDER BY due_date ASC`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deadlines []*models.Deadline
	for rows.Next() {
		var d models.Deadline
		err := rows.Scan(&d.ID, &d.ProjectID, &d.Type, &d.Description, &d.DueDate,
			&d.IsBusinessDay, &d.Priority, &d.Completed, &d.Notes)
		if err != nil {
			return nil, err
		}
		deadlines = append(deadlines, &d)
	}

	return deadlines, nil
}

// UpdateDeadline updates an existing deadline
func (s *SQLiteStorage) UpdateDeadline(d *models.Deadline) error {
	_, err := s.db.Exec(`
		UPDATE deadlines SET type = ?, description = ?, due_date = ?, is_business_day = ?,
			priority = ?, completed = ?, notes = ?
		WHERE id = ?`,
		d.Type, d.Description, d.DueDate, d.IsBusinessDay, d.Priority, d.Completed, d.Notes, d.ID)
	return err
}

// DeleteDeadline removes a deadline
func (s *SQLiteStorage) DeleteDeadline(id string) error {
	_, err := s.db.Exec("DELETE FROM deadlines WHERE id = ?", id)
	return err
}

// CreateDocument inserts a new document
func (s *SQLiteStorage) CreateDocument(doc *models.Document) error {
	_, err := s.db.Exec(`
		INSERT INTO documents (id, project_id, type, form_number, status, generated_at, file_path)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		doc.ID, doc.ProjectID, doc.Type, doc.FormNumber, doc.Status, doc.GeneratedAt, doc.FilePath)
	return err
}

// GetDocument retrieves a document by ID
func (s *SQLiteStorage) GetDocument(id string) (*models.Document, error) {
	var doc models.Document
	var generatedAt sql.NullTime

	err := s.db.QueryRow(`
		SELECT id, project_id, type, form_number, status, generated_at, file_path
		FROM documents WHERE id = ?`, id).Scan(
		&doc.ID, &doc.ProjectID, &doc.Type, &doc.FormNumber, &doc.Status, &generatedAt, &doc.FilePath)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("document not found")
		}
		return nil, err
	}

	if generatedAt.Valid {
		doc.GeneratedAt = &generatedAt.Time
	}

	return &doc, nil
}

// ListDocuments retrieves all documents for a project
func (s *SQLiteStorage) ListDocuments(projectID string) ([]*models.Document, error) {
	rows, err := s.db.Query(`
		SELECT id, project_id, type, form_number, status, generated_at, file_path
		FROM documents WHERE project_id = ?`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var documents []*models.Document
	for rows.Next() {
		var doc models.Document
		var generatedAt sql.NullTime

		err := rows.Scan(&doc.ID, &doc.ProjectID, &doc.Type, &doc.FormNumber, &doc.Status, &generatedAt, &doc.FilePath)
		if err != nil {
			return nil, err
		}

		if generatedAt.Valid {
			doc.GeneratedAt = &generatedAt.Time
		}

		documents = append(documents, &doc)
	}

	return documents, nil
}

// UpdateDocument updates an existing document
func (s *SQLiteStorage) UpdateDocument(doc *models.Document) error {
	_, err := s.db.Exec(`
		UPDATE documents SET type = ?, form_number = ?, status = ?, generated_at = ?, file_path = ?
		WHERE id = ?`,
		doc.Type, doc.FormNumber, doc.Status, doc.GeneratedAt, doc.FilePath, doc.ID)
	return err
}

// ========== Service Marketplace Methods ==========

// CreateService inserts a new service
func (s *SQLiteStorage) CreateService(service *models.Service) error {
	_, err := s.db.Exec(`
		INSERT INTO services (id, name, category, description, typical_timeline, estimated_cost_min, estimated_cost_max, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		service.ID, service.Name, service.Category, service.Description, service.TypicalTimeline,
		service.EstimatedCostMin, service.EstimatedCostMax, service.CreatedAt, service.UpdatedAt)
	return err
}

// GetService retrieves a service by ID
func (s *SQLiteStorage) GetService(id string) (*models.Service, error) {
	var service models.Service
	err := s.db.QueryRow(`
		SELECT id, name, category, description, typical_timeline, estimated_cost_min, estimated_cost_max, created_at, updated_at
		FROM services WHERE id = ?`, id).Scan(
		&service.ID, &service.Name, &service.Category, &service.Description, &service.TypicalTimeline,
		&service.EstimatedCostMin, &service.EstimatedCostMax, &service.CreatedAt, &service.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &service, nil
}

// ListServices retrieves all services
func (s *SQLiteStorage) ListServices() ([]*models.Service, error) {
	rows, err := s.db.Query(`
		SELECT id, name, category, description, typical_timeline, estimated_cost_min, estimated_cost_max, created_at, updated_at
		FROM services ORDER BY category, name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []*models.Service
	for rows.Next() {
		var service models.Service
		if err := rows.Scan(&service.ID, &service.Name, &service.Category, &service.Description,
			&service.TypicalTimeline, &service.EstimatedCostMin, &service.EstimatedCostMax,
			&service.CreatedAt, &service.UpdatedAt); err != nil {
			return nil, err
		}
		services = append(services, &service)
	}
	return services, nil
}

// ListServicesByCategory retrieves services by category
func (s *SQLiteStorage) ListServicesByCategory(category string) ([]*models.Service, error) {
	rows, err := s.db.Query(`
		SELECT id, name, category, description, typical_timeline, estimated_cost_min, estimated_cost_max, created_at, updated_at
		FROM services WHERE category = ? ORDER BY name`, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []*models.Service
	for rows.Next() {
		var service models.Service
		if err := rows.Scan(&service.ID, &service.Name, &service.Category, &service.Description,
			&service.TypicalTimeline, &service.EstimatedCostMin, &service.EstimatedCostMax,
			&service.CreatedAt, &service.UpdatedAt); err != nil {
			return nil, err
		}
		services = append(services, &service)
	}
	return services, nil
}

// CreateProvider inserts a new provider
func (s *SQLiteStorage) CreateProvider(provider *models.Provider) error {
	_, err := s.db.Exec(`
		INSERT INTO providers (id, service_id, business_name, contact_name, email, phone, address, city, state, zip,
			bio, years_experience, license_number, insurance_verified, availability_status, rating_average, rating_count, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		provider.ID, provider.ServiceID, provider.BusinessName, provider.ContactName, provider.Email, provider.Phone,
		provider.Address, provider.City, provider.State, provider.Zip, provider.Bio, provider.YearsExperience,
		provider.LicenseNumber, provider.InsuranceVerified, provider.AvailabilityStatus, provider.RatingAverage,
		provider.RatingCount, provider.CreatedAt, provider.UpdatedAt)
	return err
}

// GetProvider retrieves a provider by ID
func (s *SQLiteStorage) GetProvider(id string) (*models.Provider, error) {
	var provider models.Provider
	err := s.db.QueryRow(`
		SELECT id, service_id, business_name, contact_name, email, phone, address, city, state, zip,
			bio, years_experience, license_number, insurance_verified, availability_status, rating_average, rating_count, created_at, updated_at
		FROM providers WHERE id = ?`, id).Scan(
		&provider.ID, &provider.ServiceID, &provider.BusinessName, &provider.ContactName, &provider.Email, &provider.Phone,
		&provider.Address, &provider.City, &provider.State, &provider.Zip, &provider.Bio, &provider.YearsExperience,
		&provider.LicenseNumber, &provider.InsuranceVerified, &provider.AvailabilityStatus, &provider.RatingAverage,
		&provider.RatingCount, &provider.CreatedAt, &provider.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &provider, nil
}

// ListProviders retrieves all providers for a service
func (s *SQLiteStorage) ListProviders(serviceID string) ([]*models.Provider, error) {
	rows, err := s.db.Query(`
		SELECT id, service_id, business_name, contact_name, email, phone, address, city, state, zip,
			bio, years_experience, license_number, insurance_verified, availability_status, rating_average, rating_count, created_at, updated_at
		FROM providers WHERE service_id = ? ORDER BY rating_average DESC, business_name`, serviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var providers []*models.Provider
	for rows.Next() {
		var provider models.Provider
		if err := rows.Scan(&provider.ID, &provider.ServiceID, &provider.BusinessName, &provider.ContactName,
			&provider.Email, &provider.Phone, &provider.Address, &provider.City, &provider.State, &provider.Zip,
			&provider.Bio, &provider.YearsExperience, &provider.LicenseNumber, &provider.InsuranceVerified,
			&provider.AvailabilityStatus, &provider.RatingAverage, &provider.RatingCount,
			&provider.CreatedAt, &provider.UpdatedAt); err != nil {
			return nil, err
		}
		providers = append(providers, &provider)
	}
	return providers, nil
}

// UpdateProvider updates an existing provider
func (s *SQLiteStorage) UpdateProvider(provider *models.Provider) error {
	_, err := s.db.Exec(`
		UPDATE providers SET business_name = ?, contact_name = ?, email = ?, phone = ?, address = ?, city = ?,
			state = ?, zip = ?, bio = ?, years_experience = ?, license_number = ?, insurance_verified = ?,
			availability_status = ?, rating_average = ?, rating_count = ?, updated_at = ?
		WHERE id = ?`,
		provider.BusinessName, provider.ContactName, provider.Email, provider.Phone, provider.Address,
		provider.City, provider.State, provider.Zip, provider.Bio, provider.YearsExperience, provider.LicenseNumber,
		provider.InsuranceVerified, provider.AvailabilityStatus, provider.RatingAverage, provider.RatingCount,
		provider.UpdatedAt, provider.ID)
	return err
}

// CreateServiceRequest inserts a new service request
func (s *SQLiteStorage) CreateServiceRequest(request *models.ServiceRequest) error {
	_, err := s.db.Exec(`
		INSERT INTO service_requests (id, project_id, user_email, user_name, user_phone, service_id, provider_id,
			property_address, requested_date, preferred_time, status, notes, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		request.ID, request.ProjectID, request.UserEmail, request.UserName, request.UserPhone, request.ServiceID,
		request.ProviderID, request.PropertyAddress, request.RequestedDate, request.PreferredTime, request.Status,
		request.Notes, request.CreatedAt, request.UpdatedAt)
	return err
}

// GetServiceRequest retrieves a service request by ID
func (s *SQLiteStorage) GetServiceRequest(id string) (*models.ServiceRequest, error) {
	var request models.ServiceRequest
	err := s.db.QueryRow(`
		SELECT id, project_id, user_email, user_name, user_phone, service_id, provider_id,
			property_address, requested_date, preferred_time, status, notes, created_at, updated_at
		FROM service_requests WHERE id = ?`, id).Scan(
		&request.ID, &request.ProjectID, &request.UserEmail, &request.UserName, &request.UserPhone, &request.ServiceID,
		&request.ProviderID, &request.PropertyAddress, &request.RequestedDate, &request.PreferredTime, &request.Status,
		&request.Notes, &request.CreatedAt, &request.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &request, nil
}

// ListServiceRequests retrieves all service requests for a user
func (s *SQLiteStorage) ListServiceRequests(userEmail string) ([]*models.ServiceRequest, error) {
	rows, err := s.db.Query(`
		SELECT id, project_id, user_email, user_name, user_phone, service_id, provider_id,
			property_address, requested_date, preferred_time, status, notes, created_at, updated_at
		FROM service_requests WHERE user_email = ? ORDER BY created_at DESC`, userEmail)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []*models.ServiceRequest
	for rows.Next() {
		var request models.ServiceRequest
		if err := rows.Scan(&request.ID, &request.ProjectID, &request.UserEmail, &request.UserName, &request.UserPhone,
			&request.ServiceID, &request.ProviderID, &request.PropertyAddress, &request.RequestedDate,
			&request.PreferredTime, &request.Status, &request.Notes, &request.CreatedAt, &request.UpdatedAt); err != nil {
			return nil, err
		}
		requests = append(requests, &request)
	}
	return requests, nil
}

// UpdateServiceRequest updates an existing service request
func (s *SQLiteStorage) UpdateServiceRequest(request *models.ServiceRequest) error {
	_, err := s.db.Exec(`
		UPDATE service_requests SET provider_id = ?, property_address = ?, requested_date = ?,
			preferred_time = ?, status = ?, notes = ?, updated_at = ?
		WHERE id = ?`,
		request.ProviderID, request.PropertyAddress, request.RequestedDate, request.PreferredTime,
		request.Status, request.Notes, request.UpdatedAt, request.ID)
	return err
}

// CreateProviderReview inserts a new provider review
func (s *SQLiteStorage) CreateProviderReview(review *models.ProviderReview) error {
	_, err := s.db.Exec(`
		INSERT INTO provider_reviews (id, provider_id, service_request_id, user_email, rating, review_text, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		review.ID, review.ProviderID, review.ServiceRequestID, review.UserEmail, review.Rating, review.ReviewText, review.CreatedAt)
	return err
}

// ListProviderReviews retrieves all reviews for a provider
func (s *SQLiteStorage) ListProviderReviews(providerID string) ([]*models.ProviderReview, error) {
	rows, err := s.db.Query(`
		SELECT id, provider_id, service_request_id, user_email, rating, review_text, created_at
		FROM provider_reviews WHERE provider_id = ? ORDER BY created_at DESC`, providerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []*models.ProviderReview
	for rows.Next() {
		var review models.ProviderReview
		if err := rows.Scan(&review.ID, &review.ProviderID, &review.ServiceRequestID, &review.UserEmail,
			&review.Rating, &review.ReviewText, &review.CreatedAt); err != nil {
			return nil, err
		}
		reviews = append(reviews, &review)
	}
	return reviews, nil
}

// Close closes the database connection
func (s *SQLiteStorage) Close() error {
	return s.db.Close()
}
