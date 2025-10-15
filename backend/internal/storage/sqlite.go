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

// Close closes the database connection
func (s *SQLiteStorage) Close() error {
	return s.db.Close()
}
