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
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		email TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		name TEXT NOT NULL,
		phone TEXT NOT NULL,
		default_role TEXT NOT NULL,
		created_at DATETIME NOT NULL
	);

	CREATE TABLE IF NOT EXISTS refresh_tokens (
		token TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		expires_at DATETIME NOT NULL,
		revoked BOOLEAN DEFAULT 0,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);

	CREATE TABLE IF NOT EXISTS project_participants (
		project_id TEXT NOT NULL,
		user_id TEXT NOT NULL,
		role TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		PRIMARY KEY (project_id, user_id),
		FOREIGN KEY (project_id) REFERENCES projects(id),
		FOREIGN KEY (user_id) REFERENCES users(id)
	);

	CREATE TABLE IF NOT EXISTS projects (
		id TEXT PRIMARY KEY,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		property_address TEXT NOT NULL,
		user_role TEXT NOT NULL,
		seller_names TEXT NOT NULL,
		seller_email TEXT,
		seller_phone TEXT,
		buyer_names TEXT NOT NULL,
		buyer_email TEXT,
		buyer_phone TEXT,
		has_agent BOOLEAN,
		agent_name TEXT,
		target_list_date DATETIME,
		status TEXT NOT NULL,
		owner_user_id TEXT,
		FOREIGN KEY (owner_user_id) REFERENCES users(id)
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

	CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user ON refresh_tokens(user_id);
	CREATE INDEX IF NOT EXISTS idx_project_participants_user ON project_participants(user_id);
	CREATE INDEX IF NOT EXISTS idx_projects_owner ON projects(owner_user_id);
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
	// Ensure seller_names and buyer_names are at least empty arrays, not nil
	if p.SellerNames == nil {
		p.SellerNames = []string{}
	}
	if p.BuyerNames == nil {
		p.BuyerNames = []string{}
	}

	sellerNamesJSON, err := json.Marshal(p.SellerNames)
	if err != nil {
		return fmt.Errorf("failed to marshal seller names: %w", err)
	}

	buyerNamesJSON, err := json.Marshal(p.BuyerNames)
	if err != nil {
		return fmt.Errorf("failed to marshal buyer names: %w", err)
	}

	_, err = s.db.Exec(`
		INSERT INTO projects (id, created_at, updated_at, property_address, user_role, seller_names,
			seller_email, seller_phone, buyer_names, buyer_email, buyer_phone, has_agent, agent_name, target_list_date, status, owner_user_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		p.ID, p.CreatedAt, p.UpdatedAt, p.PropertyAddress, p.UserRole, sellerNamesJSON,
		p.SellerEmail, p.SellerPhone, buyerNamesJSON, p.BuyerEmail, p.BuyerPhone, p.HasAgent, p.AgentName, p.TargetListDate, p.Status, p.OwnerUserID)
	if err != nil {
		return fmt.Errorf("failed to insert project: %w", err)
	}
	return nil
}

// GetProject retrieves a project by ID
func (s *SQLiteStorage) GetProject(id string) (*models.Project, error) {
	var p models.Project
	var sellerNamesJSON, buyerNamesJSON string
	var targetListDate sql.NullTime
	var ownerUserID sql.NullString

	err := s.db.QueryRow(`
		SELECT id, created_at, updated_at, property_address, user_role, seller_names,
			seller_email, seller_phone, buyer_names, buyer_email, buyer_phone, has_agent, agent_name, target_list_date, status, owner_user_id
		FROM projects WHERE id = ?`, id).Scan(
		&p.ID, &p.CreatedAt, &p.UpdatedAt, &p.PropertyAddress, &p.UserRole, &sellerNamesJSON,
		&p.SellerEmail, &p.SellerPhone, &buyerNamesJSON, &p.BuyerEmail, &p.BuyerPhone, &p.HasAgent, &p.AgentName, &targetListDate, &p.Status, &ownerUserID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("project not found")
		}
		return nil, err
	}

	json.Unmarshal([]byte(sellerNamesJSON), &p.SellerNames)
	json.Unmarshal([]byte(buyerNamesJSON), &p.BuyerNames)
	if targetListDate.Valid {
		p.TargetListDate = &targetListDate.Time
	}
	if ownerUserID.Valid {
		p.OwnerUserID = &ownerUserID.String
	}

	return &p, nil
}

// UpdateProject updates an existing project
func (s *SQLiteStorage) UpdateProject(p *models.Project) error {
	p.UpdatedAt = time.Now()
	sellerNamesJSON, err := json.Marshal(p.SellerNames)
	if err != nil {
		return fmt.Errorf("failed to marshal seller names: %w", err)
	}
	buyerNamesJSON, err := json.Marshal(p.BuyerNames)
	if err != nil {
		return fmt.Errorf("failed to marshal buyer names: %w", err)
	}
	_, err = s.db.Exec(`
		UPDATE projects SET updated_at = ?, property_address = ?, user_role = ?, seller_names = ?,
			seller_email = ?, seller_phone = ?, buyer_names = ?, buyer_email = ?, buyer_phone = ?,
			has_agent = ?, agent_name = ?, target_list_date = ?, status = ?, owner_user_id = ?
		WHERE id = ?`,
		p.UpdatedAt, p.PropertyAddress, p.UserRole, sellerNamesJSON, p.SellerEmail, p.SellerPhone,
		buyerNamesJSON, p.BuyerEmail, p.BuyerPhone, p.HasAgent, p.AgentName, p.TargetListDate, p.Status, p.OwnerUserID, p.ID)
	return err
}

// ListProjects returns all projects
func (s *SQLiteStorage) ListProjects() ([]*models.Project, error) {
	rows, err := s.db.Query(`
		SELECT id, created_at, updated_at, property_address, user_role, seller_names,
			seller_email, seller_phone, buyer_names, buyer_email, buyer_phone, has_agent, agent_name, target_list_date, status, owner_user_id
		FROM projects ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*models.Project
	for rows.Next() {
		var p models.Project
		var sellerNamesJSON, buyerNamesJSON string
		var targetListDate sql.NullTime
		var ownerUserID sql.NullString

		err := rows.Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt, &p.PropertyAddress, &p.UserRole, &sellerNamesJSON,
			&p.SellerEmail, &p.SellerPhone, &buyerNamesJSON, &p.BuyerEmail, &p.BuyerPhone, &p.HasAgent, &p.AgentName, &targetListDate, &p.Status, &ownerUserID)
		if err != nil {
			return nil, err
		}

		json.Unmarshal([]byte(sellerNamesJSON), &p.SellerNames)
		json.Unmarshal([]byte(buyerNamesJSON), &p.BuyerNames)
		if targetListDate.Valid {
			p.TargetListDate = &targetListDate.Time
		}
		if ownerUserID.Valid {
			p.OwnerUserID = &ownerUserID.String
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

// CreateUser inserts a new user
func (s *SQLiteStorage) CreateUser(user *models.User) error {
	_, err := s.db.Exec(`
		INSERT INTO users (id, email, password_hash, name, phone, default_role, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		user.ID, user.Email, user.PasswordHash, user.Name, user.Phone, user.DefaultRole, user.CreatedAt)
	return err
}

// GetUser retrieves a user by ID
func (s *SQLiteStorage) GetUser(id string) (*models.User, error) {
	var user models.User
	err := s.db.QueryRow(`
		SELECT id, email, password_hash, name, phone, default_role, created_at
		FROM users WHERE id = ?`, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Phone, &user.DefaultRole, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// GetUserByEmail retrieves a user by email
func (s *SQLiteStorage) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := s.db.QueryRow(`
		SELECT id, email, password_hash, name, phone, default_role, created_at
		FROM users WHERE email = ?`, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Phone, &user.DefaultRole, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// CreateRefreshToken stores a new refresh token
func (s *SQLiteStorage) CreateRefreshToken(token *models.RefreshToken) error {
	_, err := s.db.Exec(`
		INSERT INTO refresh_tokens (token, user_id, expires_at, revoked)
		VALUES (?, ?, ?, ?)`,
		token.Token, token.UserID, token.ExpiresAt, token.Revoked)
	return err
}

// GetRefreshToken retrieves a refresh token
func (s *SQLiteStorage) GetRefreshToken(token string) (*models.RefreshToken, error) {
	var rt models.RefreshToken
	err := s.db.QueryRow(`
		SELECT token, user_id, expires_at, revoked
		FROM refresh_tokens WHERE token = ?`, token).Scan(
		&rt.Token, &rt.UserID, &rt.ExpiresAt, &rt.Revoked)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("refresh token not found")
		}
		return nil, err
	}

	return &rt, nil
}

// RevokeRefreshToken marks a refresh token as revoked
func (s *SQLiteStorage) RevokeRefreshToken(token string) error {
	_, err := s.db.Exec(`UPDATE refresh_tokens SET revoked = 1 WHERE token = ?`, token)
	return err
}

// DeleteExpiredRefreshTokens removes expired tokens
func (s *SQLiteStorage) DeleteExpiredRefreshTokens() error {
	_, err := s.db.Exec(`DELETE FROM refresh_tokens WHERE expires_at < ?`, time.Now())
	return err
}

// CreateProjectParticipant adds a participant to a project
func (s *SQLiteStorage) CreateProjectParticipant(participant *models.ProjectParticipant) error {
	_, err := s.db.Exec(`
		INSERT INTO project_participants (project_id, user_id, role, created_at)
		VALUES (?, ?, ?, ?)`,
		participant.ProjectID, participant.UserID, participant.Role, participant.CreatedAt)
	return err
}

// GetProjectParticipant retrieves a specific participant
func (s *SQLiteStorage) GetProjectParticipant(projectID, userID string) (*models.ProjectParticipant, error) {
	var participant models.ProjectParticipant
	err := s.db.QueryRow(`
		SELECT project_id, user_id, role, created_at
		FROM project_participants WHERE project_id = ? AND user_id = ?`, projectID, userID).Scan(
		&participant.ProjectID, &participant.UserID, &participant.Role, &participant.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("participant not found")
		}
		return nil, err
	}

	return &participant, nil
}

// ListProjectParticipants retrieves all participants for a project
func (s *SQLiteStorage) ListProjectParticipants(projectID string) ([]*models.ProjectParticipant, error) {
	rows, err := s.db.Query(`
		SELECT project_id, user_id, role, created_at
		FROM project_participants WHERE project_id = ?`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []*models.ProjectParticipant
	for rows.Next() {
		var participant models.ProjectParticipant
		err := rows.Scan(&participant.ProjectID, &participant.UserID, &participant.Role, &participant.CreatedAt)
		if err != nil {
			return nil, err
		}
		participants = append(participants, &participant)
	}

	return participants, nil
}

// ListUserProjects retrieves all projects where user is a participant
func (s *SQLiteStorage) ListUserProjects(userID string) ([]*models.Project, error) {
	rows, err := s.db.Query(`
		SELECT p.id, p.created_at, p.updated_at, p.property_address, p.user_role, p.seller_names,
			p.seller_email, p.seller_phone, p.buyer_names, p.buyer_email, p.buyer_phone,
			p.has_agent, p.agent_name, p.target_list_date, p.status, p.owner_user_id
		FROM projects p
		INNER JOIN project_participants pp ON p.id = pp.project_id
		WHERE pp.user_id = ?
		ORDER BY p.created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*models.Project
	for rows.Next() {
		var p models.Project
		var sellerNamesJSON, buyerNamesJSON string
		var targetListDate sql.NullTime
		var ownerUserID sql.NullString

		err := rows.Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt, &p.PropertyAddress, &p.UserRole, &sellerNamesJSON,
			&p.SellerEmail, &p.SellerPhone, &buyerNamesJSON, &p.BuyerEmail, &p.BuyerPhone,
			&p.HasAgent, &p.AgentName, &targetListDate, &p.Status, &ownerUserID)
		if err != nil {
			return nil, err
		}

		json.Unmarshal([]byte(sellerNamesJSON), &p.SellerNames)
		json.Unmarshal([]byte(buyerNamesJSON), &p.BuyerNames)
		if targetListDate.Valid {
			p.TargetListDate = &targetListDate.Time
		}
		if ownerUserID.Valid {
			p.OwnerUserID = &ownerUserID.String
		}

		projects = append(projects, &p)
	}

	return projects, nil
}

// Close closes the database connection
func (s *SQLiteStorage) Close() error {
	return s.db.Close()
}
