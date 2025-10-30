package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jessicaandtommymccay/realty-wizard/backend/internal/models"
	_ "github.com/lib/pq"
)

// PostgresStorage implements Storage using PostgreSQL
type PostgresStorage struct {
	db *sql.DB
}

// NewPostgresStorage creates a new PostgreSQL storage instance
func NewPostgresStorage(connectionString string) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	storage := &PostgresStorage{db: db}
	if err := storage.migrate(); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return storage, nil
}

// migrate creates the database schema
func (s *PostgresStorage) migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		email TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		name TEXT NOT NULL,
		phone TEXT NOT NULL,
		default_role TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL
	);

	CREATE TABLE IF NOT EXISTS refresh_tokens (
		token TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		expires_at TIMESTAMP NOT NULL,
		revoked BOOLEAN DEFAULT false,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);

	CREATE TABLE IF NOT EXISTS projects (
		id TEXT PRIMARY KEY,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL,
		property_address TEXT NOT NULL,
		user_role TEXT NOT NULL,
		seller_names JSONB NOT NULL,
		seller_email TEXT,
		seller_phone TEXT,
		buyer_names JSONB NOT NULL,
		buyer_email TEXT,
		buyer_phone TEXT,
		has_agent BOOLEAN,
		agent_name TEXT,
		target_list_date TIMESTAMP,
		status TEXT NOT NULL,
		owner_user_id TEXT,
		FOREIGN KEY (owner_user_id) REFERENCES users(id)
	);

	CREATE TABLE IF NOT EXISTS project_participants (
		project_id TEXT NOT NULL,
		user_id TEXT NOT NULL,
		role TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL,
		PRIMARY KEY (project_id, user_id),
		FOREIGN KEY (project_id) REFERENCES projects(id),
		FOREIGN KEY (user_id) REFERENCES users(id)
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
		hoa_fee DOUBLE PRECISION,
		hoa_frequency TEXT,
		has_survey BOOLEAN,
		survey_date TIMESTAMP,
		survey_conditions TEXT,
		FOREIGN KEY (project_id) REFERENCES projects(id)
	);

	CREATE TABLE IF NOT EXISTS disclosures (
		id TEXT PRIMARY KEY,
		project_id TEXT NOT NULL,
		updated_at TIMESTAMP NOT NULL,
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
		effective_date TIMESTAMP NOT NULL,
		sales_price DOUBLE PRECISION NOT NULL,
		closing_date TIMESTAMP NOT NULL,
		option_fee DOUBLE PRECISION,
		option_period_days INTEGER,
		earnest_money DOUBLE PRECISION,
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
		due_date TIMESTAMP NOT NULL,
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
		generated_at TIMESTAMP,
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

// User operations
func (s *PostgresStorage) CreateUser(user *models.User) error {
	_, err := s.db.Exec(`
		INSERT INTO users (id, email, password_hash, name, phone, default_role, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		user.ID, user.Email, user.PasswordHash, user.Name, user.Phone, user.DefaultRole, user.CreatedAt)
	return err
}

func (s *PostgresStorage) GetUser(id string) (*models.User, error) {
	var user models.User
	err := s.db.QueryRow(`
		SELECT id, email, password_hash, name, phone, default_role, created_at
		FROM users WHERE id = $1`, id).
		Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Phone, &user.DefaultRole, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *PostgresStorage) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := s.db.QueryRow(`
		SELECT id, email, password_hash, name, phone, default_role, created_at
		FROM users WHERE email = $1`, email).
		Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Phone, &user.DefaultRole, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Refresh token operations
func (s *PostgresStorage) CreateRefreshToken(token *models.RefreshToken) error {
	_, err := s.db.Exec(`
		INSERT INTO refresh_tokens (token, user_id, expires_at, revoked)
		VALUES ($1, $2, $3, $4)`,
		token.Token, token.UserID, token.ExpiresAt, token.Revoked)
	return err
}

func (s *PostgresStorage) GetRefreshToken(token string) (*models.RefreshToken, error) {
	var rt models.RefreshToken
	err := s.db.QueryRow(`
		SELECT token, user_id, expires_at, revoked
		FROM refresh_tokens WHERE token = $1`, token).
		Scan(&rt.Token, &rt.UserID, &rt.ExpiresAt, &rt.Revoked)
	if err != nil {
		return nil, err
	}
	return &rt, nil
}

func (s *PostgresStorage) RevokeRefreshToken(token string) error {
	_, err := s.db.Exec(`UPDATE refresh_tokens SET revoked = true WHERE token = $1`, token)
	return err
}

func (s *PostgresStorage) DeleteExpiredRefreshTokens() error {
	_, err := s.db.Exec(`DELETE FROM refresh_tokens WHERE expires_at < $1`, time.Now())
	return err
}

// Project participant operations
func (s *PostgresStorage) CreateProjectParticipant(participant *models.ProjectParticipant) error {
	_, err := s.db.Exec(`
		INSERT INTO project_participants (project_id, user_id, role, created_at)
		VALUES ($1, $2, $3, $4)`,
		participant.ProjectID, participant.UserID, participant.Role, participant.CreatedAt)
	return err
}

func (s *PostgresStorage) GetProjectParticipant(projectID, userID string) (*models.ProjectParticipant, error) {
	var pp models.ProjectParticipant
	err := s.db.QueryRow(`
		SELECT project_id, user_id, role, created_at
		FROM project_participants WHERE project_id = $1 AND user_id = $2`, projectID, userID).
		Scan(&pp.ProjectID, &pp.UserID, &pp.Role, &pp.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &pp, nil
}

func (s *PostgresStorage) ListProjectParticipants(projectID string) ([]*models.ProjectParticipant, error) {
	rows, err := s.db.Query(`
		SELECT project_id, user_id, role, created_at
		FROM project_participants WHERE project_id = $1`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []*models.ProjectParticipant
	for rows.Next() {
		var pp models.ProjectParticipant
		if err := rows.Scan(&pp.ProjectID, &pp.UserID, &pp.Role, &pp.CreatedAt); err != nil {
			return nil, err
		}
		participants = append(participants, &pp)
	}
	return participants, nil
}

func (s *PostgresStorage) ListUserProjects(userID string) ([]*models.Project, error) {
	rows, err := s.db.Query(`
		SELECT p.id, p.created_at, p.updated_at, p.property_address, p.user_role, p.seller_names,
			p.seller_email, p.seller_phone, p.buyer_names, p.buyer_email, p.buyer_phone,
			p.has_agent, p.agent_name, p.target_list_date, p.status, p.owner_user_id
		FROM projects p
		INNER JOIN project_participants pp ON p.id = pp.project_id
		WHERE pp.user_id = $1
		ORDER BY p.updated_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*models.Project
	for rows.Next() {
		var p models.Project
		var sellerNamesJSON, buyerNamesJSON []byte
		var targetListDate sql.NullTime
		var ownerUserID sql.NullString

		if err := rows.Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt, &p.PropertyAddress, &p.UserRole, &sellerNamesJSON,
			&p.SellerEmail, &p.SellerPhone, &buyerNamesJSON, &p.BuyerEmail, &p.BuyerPhone,
			&p.HasAgent, &p.AgentName, &targetListDate, &p.Status, &ownerUserID); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(sellerNamesJSON, &p.SellerNames); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(buyerNamesJSON, &p.BuyerNames); err != nil {
			return nil, err
		}
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

// Project operations
func (s *PostgresStorage) CreateProject(p *models.Project) error {
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
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`,
		p.ID, p.CreatedAt, p.UpdatedAt, p.PropertyAddress, p.UserRole, sellerNamesJSON,
		p.SellerEmail, p.SellerPhone, buyerNamesJSON, p.BuyerEmail, p.BuyerPhone, p.HasAgent, p.AgentName, p.TargetListDate, p.Status, p.OwnerUserID)
	return err
}

func (s *PostgresStorage) GetProject(id string) (*models.Project, error) {
	var p models.Project
	var sellerNamesJSON, buyerNamesJSON []byte
	var targetListDate sql.NullTime
	var ownerUserID sql.NullString

	err := s.db.QueryRow(`
		SELECT id, created_at, updated_at, property_address, user_role, seller_names,
			seller_email, seller_phone, buyer_names, buyer_email, buyer_phone, has_agent, agent_name, target_list_date, status, owner_user_id
		FROM projects WHERE id = $1`, id).
		Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt, &p.PropertyAddress, &p.UserRole, &sellerNamesJSON,
			&p.SellerEmail, &p.SellerPhone, &buyerNamesJSON, &p.BuyerEmail, &p.BuyerPhone, &p.HasAgent, &p.AgentName, &targetListDate, &p.Status, &ownerUserID)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(sellerNamesJSON, &p.SellerNames); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(buyerNamesJSON, &p.BuyerNames); err != nil {
		return nil, err
	}
	if targetListDate.Valid {
		p.TargetListDate = &targetListDate.Time
	}
	if ownerUserID.Valid {
		p.OwnerUserID = &ownerUserID.String
	}

	return &p, nil
}

func (s *PostgresStorage) UpdateProject(p *models.Project) error {
	sellerNamesJSON, err := json.Marshal(p.SellerNames)
	if err != nil {
		return fmt.Errorf("failed to marshal seller names: %w", err)
	}
	buyerNamesJSON, err := json.Marshal(p.BuyerNames)
	if err != nil {
		return fmt.Errorf("failed to marshal buyer names: %w", err)
	}

	_, err = s.db.Exec(`
		UPDATE projects SET updated_at = $1, property_address = $2, user_role = $3, seller_names = $4,
			seller_email = $5, seller_phone = $6, buyer_names = $7, buyer_email = $8, buyer_phone = $9,
			has_agent = $10, agent_name = $11, target_list_date = $12, status = $13, owner_user_id = $14
		WHERE id = $15`,
		p.UpdatedAt, p.PropertyAddress, p.UserRole, sellerNamesJSON, p.SellerEmail, p.SellerPhone,
		buyerNamesJSON, p.BuyerEmail, p.BuyerPhone, p.HasAgent, p.AgentName, p.TargetListDate, p.Status, p.OwnerUserID, p.ID)
	return err
}

func (s *PostgresStorage) ListProjects() ([]*models.Project, error) {
	rows, err := s.db.Query(`
		SELECT id, created_at, updated_at, property_address, user_role, seller_names,
			seller_email, seller_phone, buyer_names, buyer_email, buyer_phone, has_agent, agent_name, target_list_date, status, owner_user_id
		FROM projects ORDER BY updated_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*models.Project
	for rows.Next() {
		var p models.Project
		var sellerNamesJSON, buyerNamesJSON []byte
		var targetListDate sql.NullTime
		var ownerUserID sql.NullString

		if err := rows.Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt, &p.PropertyAddress, &p.UserRole, &sellerNamesJSON,
			&p.SellerEmail, &p.SellerPhone, &buyerNamesJSON, &p.BuyerEmail, &p.BuyerPhone, &p.HasAgent, &p.AgentName,
			&targetListDate, &p.Status, &ownerUserID); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(sellerNamesJSON, &p.SellerNames); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(buyerNamesJSON, &p.BuyerNames); err != nil {
			return nil, err
		}
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

// Property operations
func (s *PostgresStorage) CreateProperty(property *models.Property) error {
	_, err := s.db.Exec(`
		INSERT INTO properties (id, project_id, year_built, legal_description, tax_id,
			is_homestead, has_hoa, hoa_name, hoa_fee, hoa_frequency, has_survey, survey_date, survey_conditions)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
		property.ID, property.ProjectID, property.YearBuilt, property.LegalDescription, property.TaxID,
		property.IsHomestead, property.HasHOA, property.HOAName, property.HOAFee, property.HOAFrequency,
		property.HasSurvey, property.SurveyDate, property.SurveyConditions)
	return err
}

func (s *PostgresStorage) GetProperty(projectID string) (*models.Property, error) {
	var prop models.Property
	var yearBuilt sql.NullInt64
	var surveyDate sql.NullTime

	err := s.db.QueryRow(`
		SELECT id, project_id, year_built, legal_description, tax_id,
			is_homestead, has_hoa, hoa_name, hoa_fee, hoa_frequency, has_survey, survey_date, survey_conditions
		FROM properties WHERE project_id = $1`, projectID).
		Scan(&prop.ID, &prop.ProjectID, &yearBuilt, &prop.LegalDescription, &prop.TaxID,
			&prop.IsHomestead, &prop.HasHOA, &prop.HOAName, &prop.HOAFee, &prop.HOAFrequency,
			&prop.HasSurvey, &surveyDate, &prop.SurveyConditions)
	if err != nil {
		return nil, err
	}

	if yearBuilt.Valid {
		prop.YearBuilt = int(yearBuilt.Int64)
	}
	if surveyDate.Valid {
		prop.SurveyDate = &surveyDate.Time
	}

	return &prop, nil
}

func (s *PostgresStorage) UpdateProperty(property *models.Property) error {
	_, err := s.db.Exec(`
		UPDATE properties SET year_built = $1, legal_description = $2, tax_id = $3,
			is_homestead = $4, has_hoa = $5, hoa_name = $6, hoa_fee = $7, hoa_frequency = $8,
			has_survey = $9, survey_date = $10, survey_conditions = $11
		WHERE id = $12`,
		property.YearBuilt, property.LegalDescription, property.TaxID, property.IsHomestead,
		property.HasHOA, property.HOAName, property.HOAFee, property.HOAFrequency,
		property.HasSurvey, property.SurveyDate, property.SurveyConditions, property.ID)
	return err
}

// Disclosure operations
func (s *PostgresStorage) CreateDisclosure(disclosure *models.Disclosure) error {
	_, err := s.db.Exec(`
		INSERT INTO disclosures (id, project_id, updated_at, flooding_history, flooding_details,
			in_flood_plain, flood_insurance_required, foundation_issues, foundation_details,
			roof_age, roof_issues, roof_details, plumbing_issues, plumbing_details,
			electrical_issues, electrical_details, hvac_age, hvac_issues, hvac_details,
			lead_based_paint, insurance_claims_history, claims_details, other_material_defects)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23)`,
		disclosure.ID, disclosure.ProjectID, disclosure.UpdatedAt,
		disclosure.FloodingHistory, disclosure.FloodingDetails, disclosure.InFloodPlain,
		disclosure.FloodInsuranceRequired, disclosure.FoundationIssues, disclosure.FoundationDetails,
		disclosure.RoofAge, disclosure.RoofIssues, disclosure.RoofDetails,
		disclosure.PlumbingIssues, disclosure.PlumbingDetails, disclosure.ElectricalIssues,
		disclosure.ElectricalDetails, disclosure.HVACAge, disclosure.HVACIssues,
		disclosure.HVACDetails, disclosure.LeadBasedPaint, disclosure.InsuranceClaimsHistory,
		disclosure.ClaimsDetails, disclosure.OtherMaterialDefects)
	return err
}

func (s *PostgresStorage) GetDisclosure(projectID string) (*models.Disclosure, error) {
	var d models.Disclosure
	var roofAge, hvacAge sql.NullInt64

	err := s.db.QueryRow(`
		SELECT id, project_id, updated_at, flooding_history, flooding_details,
			in_flood_plain, flood_insurance_required, foundation_issues, foundation_details,
			roof_age, roof_issues, roof_details, plumbing_issues, plumbing_details,
			electrical_issues, electrical_details, hvac_age, hvac_issues, hvac_details,
			lead_based_paint, insurance_claims_history, claims_details, other_material_defects
		FROM disclosures WHERE project_id = $1`, projectID).
		Scan(&d.ID, &d.ProjectID, &d.UpdatedAt,
			&d.FloodingHistory, &d.FloodingDetails, &d.InFloodPlain,
			&d.FloodInsuranceRequired, &d.FoundationIssues, &d.FoundationDetails,
			&roofAge, &d.RoofIssues, &d.RoofDetails,
			&d.PlumbingIssues, &d.PlumbingDetails, &d.ElectricalIssues,
			&d.ElectricalDetails, &hvacAge, &d.HVACIssues,
			&d.HVACDetails, &d.LeadBasedPaint, &d.InsuranceClaimsHistory,
			&d.ClaimsDetails, &d.OtherMaterialDefects)
	if err != nil {
		return nil, err
	}

	if roofAge.Valid {
		d.RoofAge = int(roofAge.Int64)
	}
	if hvacAge.Valid {
		d.HVACAge = int(hvacAge.Int64)
	}

	return &d, nil
}

func (s *PostgresStorage) UpdateDisclosure(disclosure *models.Disclosure) error {
	_, err := s.db.Exec(`
		UPDATE disclosures SET updated_at = $1, flooding_history = $2, flooding_details = $3,
			in_flood_plain = $4, flood_insurance_required = $5, foundation_issues = $6, foundation_details = $7,
			roof_age = $8, roof_issues = $9, roof_details = $10, plumbing_issues = $11, plumbing_details = $12,
			electrical_issues = $13, electrical_details = $14, hvac_age = $15, hvac_issues = $16, hvac_details = $17,
			lead_based_paint = $18, insurance_claims_history = $19, claims_details = $20, other_material_defects = $21
		WHERE id = $22`,
		disclosure.UpdatedAt, disclosure.FloodingHistory, disclosure.FloodingDetails,
		disclosure.InFloodPlain, disclosure.FloodInsuranceRequired, disclosure.FoundationIssues,
		disclosure.FoundationDetails, disclosure.RoofAge, disclosure.RoofIssues, disclosure.RoofDetails,
		disclosure.PlumbingIssues, disclosure.PlumbingDetails, disclosure.ElectricalIssues,
		disclosure.ElectricalDetails, disclosure.HVACAge, disclosure.HVACIssues, disclosure.HVACDetails,
		disclosure.LeadBasedPaint, disclosure.InsuranceClaimsHistory, disclosure.ClaimsDetails,
		disclosure.OtherMaterialDefects, disclosure.ID)
	return err
}

// Contract operations
func (s *PostgresStorage) CreateContract(contract *models.ContractTerms) error {
	// Marshal arrays to JSON
	if contract.ExcludedItems == nil {
		contract.ExcludedItems = []string{}
	}
	if contract.IncludedPersonalItems == nil {
		contract.IncludedPersonalItems = []string{}
	}

	excludedItemsJSON, err := json.Marshal(contract.ExcludedItems)
	if err != nil {
		return err
	}
	includedItemsJSON, err := json.Marshal(contract.IncludedPersonalItems)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`
		INSERT INTO contracts (id, project_id, effective_date, sales_price, closing_date,
			option_fee, option_period_days, earnest_money, title_commitment_days,
			buyer_financing, seller_stays_post_close, seller_lease_days,
			excluded_items, included_personal_items)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`,
		contract.ID, contract.ProjectID, contract.EffectiveDate, contract.SalesPrice,
		contract.ClosingDate, contract.OptionFee, contract.OptionPeriodDays,
		contract.EarnestMoney, contract.TitleCommitmentDays, contract.BuyerFinancing,
		contract.SellerStaysPostClose, contract.SellerLeaseDays,
		string(excludedItemsJSON), string(includedItemsJSON))
	return err
}

func (s *PostgresStorage) GetContract(projectID string) (*models.ContractTerms, error) {
	var c models.ContractTerms
	var optionFee, earnestMoney sql.NullFloat64
	var optionPeriodDays, titleCommitmentDays, sellerLeaseDays sql.NullInt64
	var excludedItemsJSON, includedItemsJSON string

	err := s.db.QueryRow(`
		SELECT id, project_id, effective_date, sales_price, closing_date,
			option_fee, option_period_days, earnest_money, title_commitment_days,
			buyer_financing, seller_stays_post_close, seller_lease_days,
			excluded_items, included_personal_items
		FROM contracts WHERE project_id = $1`, projectID).
		Scan(&c.ID, &c.ProjectID, &c.EffectiveDate, &c.SalesPrice, &c.ClosingDate,
			&optionFee, &optionPeriodDays, &earnestMoney, &titleCommitmentDays,
			&c.BuyerFinancing, &c.SellerStaysPostClose, &sellerLeaseDays,
			&excludedItemsJSON, &includedItemsJSON)
	if err != nil {
		return nil, err
	}

	if optionFee.Valid {
		c.OptionFee = optionFee.Float64
	}
	if optionPeriodDays.Valid {
		c.OptionPeriodDays = int(optionPeriodDays.Int64)
	}
	if earnestMoney.Valid {
		c.EarnestMoney = earnestMoney.Float64
	}
	if titleCommitmentDays.Valid {
		c.TitleCommitmentDays = int(titleCommitmentDays.Int64)
	}
	if sellerLeaseDays.Valid {
		c.SellerLeaseDays = int(sellerLeaseDays.Int64)
	}

	// Unmarshal JSON arrays
	if err := json.Unmarshal([]byte(excludedItemsJSON), &c.ExcludedItems); err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(includedItemsJSON), &c.IncludedPersonalItems); err != nil {
		return nil, err
	}

	return &c, nil
}

func (s *PostgresStorage) UpdateContract(contract *models.ContractTerms) error {
	// Marshal arrays to JSON
	if contract.ExcludedItems == nil {
		contract.ExcludedItems = []string{}
	}
	if contract.IncludedPersonalItems == nil {
		contract.IncludedPersonalItems = []string{}
	}

	excludedItemsJSON, err := json.Marshal(contract.ExcludedItems)
	if err != nil {
		return err
	}
	includedItemsJSON, err := json.Marshal(contract.IncludedPersonalItems)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`
		UPDATE contracts SET effective_date = $1, sales_price = $2, closing_date = $3,
			option_fee = $4, option_period_days = $5, earnest_money = $6, title_commitment_days = $7,
			buyer_financing = $8, seller_stays_post_close = $9, seller_lease_days = $10,
			excluded_items = $11, included_personal_items = $12
		WHERE id = $13`,
		contract.EffectiveDate, contract.SalesPrice, contract.ClosingDate,
		contract.OptionFee, contract.OptionPeriodDays, contract.EarnestMoney,
		contract.TitleCommitmentDays, contract.BuyerFinancing, contract.SellerStaysPostClose,
		contract.SellerLeaseDays, string(excludedItemsJSON), string(includedItemsJSON), contract.ID)
	return err
}

// Deadline operations
func (s *PostgresStorage) CreateDeadline(deadline *models.Deadline) error {
	_, err := s.db.Exec(`
		INSERT INTO deadlines (id, project_id, type, description, due_date,
			is_business_day, priority, completed, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		deadline.ID, deadline.ProjectID, deadline.Type, deadline.Description,
		deadline.DueDate, deadline.IsBusinessDay, deadline.Priority,
		deadline.Completed, deadline.Notes)
	return err
}

func (s *PostgresStorage) GetDeadline(id string) (*models.Deadline, error) {
	var d models.Deadline
	err := s.db.QueryRow(`
		SELECT id, project_id, type, description, due_date,
			is_business_day, priority, completed, notes
		FROM deadlines WHERE id = $1`, id).
		Scan(&d.ID, &d.ProjectID, &d.Type, &d.Description, &d.DueDate,
			&d.IsBusinessDay, &d.Priority, &d.Completed, &d.Notes)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (s *PostgresStorage) ListDeadlines(projectID string) ([]*models.Deadline, error) {
	rows, err := s.db.Query(`
		SELECT id, project_id, type, description, due_date,
			is_business_day, priority, completed, notes
		FROM deadlines WHERE project_id = $1 ORDER BY due_date`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deadlines []*models.Deadline
	for rows.Next() {
		var d models.Deadline
		if err := rows.Scan(&d.ID, &d.ProjectID, &d.Type, &d.Description, &d.DueDate,
			&d.IsBusinessDay, &d.Priority, &d.Completed, &d.Notes); err != nil {
			return nil, err
		}
		deadlines = append(deadlines, &d)
	}
	return deadlines, nil
}

func (s *PostgresStorage) UpdateDeadline(deadline *models.Deadline) error {
	_, err := s.db.Exec(`
		UPDATE deadlines SET type = $1, description = $2, due_date = $3,
			is_business_day = $4, priority = $5, completed = $6, notes = $7
		WHERE id = $8`,
		deadline.Type, deadline.Description, deadline.DueDate,
		deadline.IsBusinessDay, deadline.Priority, deadline.Completed,
		deadline.Notes, deadline.ID)
	return err
}

func (s *PostgresStorage) DeleteDeadline(id string) error {
	_, err := s.db.Exec(`DELETE FROM deadlines WHERE id = $1`, id)
	return err
}

// Document operations
func (s *PostgresStorage) CreateDocument(doc *models.Document) error {
	_, err := s.db.Exec(`
		INSERT INTO documents (id, project_id, type, form_number, status, generated_at, file_path)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		doc.ID, doc.ProjectID, doc.Type, doc.FormNumber, doc.Status, doc.GeneratedAt, doc.FilePath)
	return err
}

func (s *PostgresStorage) GetDocument(id string) (*models.Document, error) {
	var doc models.Document
	var generatedAt sql.NullTime

	err := s.db.QueryRow(`
		SELECT id, project_id, type, form_number, status, generated_at, file_path
		FROM documents WHERE id = $1`, id).
		Scan(&doc.ID, &doc.ProjectID, &doc.Type, &doc.FormNumber, &doc.Status, &generatedAt, &doc.FilePath)
	if err != nil {
		return nil, err
	}

	if generatedAt.Valid {
		doc.GeneratedAt = &generatedAt.Time
	}

	return &doc, nil
}

func (s *PostgresStorage) ListDocuments(projectID string) ([]*models.Document, error) {
	rows, err := s.db.Query(`
		SELECT id, project_id, type, form_number, status, generated_at, file_path
		FROM documents WHERE project_id = $1`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var documents []*models.Document
	for rows.Next() {
		var doc models.Document
		var generatedAt sql.NullTime

		if err := rows.Scan(&doc.ID, &doc.ProjectID, &doc.Type, &doc.FormNumber, &doc.Status, &generatedAt, &doc.FilePath); err != nil {
			return nil, err
		}

		if generatedAt.Valid {
			doc.GeneratedAt = &generatedAt.Time
		}

		documents = append(documents, &doc)
	}
	return documents, nil
}

func (s *PostgresStorage) UpdateDocument(doc *models.Document) error {
	_, err := s.db.Exec(`
		UPDATE documents SET type = $1, form_number = $2, status = $3,
			generated_at = $4, file_path = $5
		WHERE id = $6`,
		doc.Type, doc.FormNumber, doc.Status, doc.GeneratedAt, doc.FilePath, doc.ID)
	return err
}

func (s *PostgresStorage) Close() error {
	return s.db.Close()
}
