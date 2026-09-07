package repository

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"time"

	"crm/internal/domain"
	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schemaSQL string

func InitDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite db: %w", err)
	}

	// PRAGMA configuration
	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		return nil, fmt.Errorf("enable foreign keys: %w", err)
	}

	if _, err := db.Exec(schemaSQL); err != nil {
		return nil, fmt.Errorf("execute schema: %w", err)
	}

	return db, nil
}

// UserRepo implements domain.UserRepository
type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, username, name, role, created_at FROM users WHERE id = ?", id)
	var u domain.User
	var role string
	var createdAt string
	if err := row.Scan(&u.ID, &u.Username, &u.Name, &role, &createdAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	u.Role = domain.UserRole(role)
	u.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	return &u, nil
}

func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, username, name, role, created_at FROM users WHERE username = ? COLLATE NOCASE", username)
	var u domain.User
	var role string
	var createdAt string
	if err := row.Scan(&u.ID, &u.Username, &u.Name, &role, &createdAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	u.Role = domain.UserRole(role)
	u.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	return &u, nil
}

func (r *UserRepo) GetByName(ctx context.Context, name string) (*domain.User, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, username, name, role, created_at FROM users WHERE name = ? COLLATE NOCASE", name)
	var u domain.User
	var role string
	var createdAt string
	if err := row.Scan(&u.ID, &u.Username, &u.Name, &role, &createdAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	u.Role = domain.UserRole(role)
	u.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	return &u, nil
}

func (r *UserRepo) Create(ctx context.Context, user *domain.User) error {
	res, err := r.db.ExecContext(ctx, "INSERT INTO users (username, name, role) VALUES (?, ?, ?)", user.Username, user.Name, string(user.Role))
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = id
	return nil
}

// CompanyRepo implements domain.CompanyRepository
type CompanyRepo struct {
	db *sql.DB
}

func NewCompanyRepo(db *sql.DB) *CompanyRepo {
	return &CompanyRepo{db: db}
}

func (r *CompanyRepo) Create(ctx context.Context, c *domain.Company) error {
	res, err := r.db.ExecContext(ctx, "INSERT INTO companies (name, tax_id, industry, address) VALUES (?, ?, ?, ?)",
		c.Name, c.TaxID, c.Industry, c.Address)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	c.ID = id
	return nil
}

func (r *CompanyRepo) GetByID(ctx context.Context, id int64) (*domain.Company, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, name, tax_id, industry, address, created_at, updated_at FROM companies WHERE id = ?", id)
	var c domain.Company
	var taxID, industry, address sql.NullString
	var createdAt, updatedAt string
	if err := row.Scan(&c.ID, &c.Name, &taxID, &industry, &address, &createdAt, &updatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	c.TaxID = taxID.String
	c.Industry = industry.String
	c.Address = address.String
	return &c, nil
}

func (r *CompanyRepo) GetByName(ctx context.Context, name string) (*domain.Company, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, name, tax_id, industry, address, created_at, updated_at FROM companies WHERE name = ?", name)
	var c domain.Company
	var taxID, industry, address sql.NullString
	var createdAt, updatedAt string
	if err := row.Scan(&c.ID, &c.Name, &taxID, &industry, &address, &createdAt, &updatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	c.TaxID = taxID.String
	c.Industry = industry.String
	c.Address = address.String
	return &c, nil
}

func (r *CompanyRepo) List(ctx context.Context) ([]domain.Company, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, tax_id, industry, address, created_at, updated_at FROM companies ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.Company
	for rows.Next() {
		var c domain.Company
		var taxID, industry, address sql.NullString
		var createdAt, updatedAt string
		if err := rows.Scan(&c.ID, &c.Name, &taxID, &industry, &address, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		c.TaxID = taxID.String
		c.Industry = industry.String
		c.Address = address.String
		result = append(result, c)
	}
	return result, nil
}

func (r *CompanyRepo) Update(ctx context.Context, c *domain.Company) error {
	res, err := r.db.ExecContext(ctx, "UPDATE companies SET name = ?, tax_id = ?, industry = ?, address = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		c.Name, c.TaxID, c.Industry, c.Address, c.ID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *CompanyRepo) HasActiveRelations(ctx context.Context, companyID int64) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM contacts WHERE company_id = ?", companyID).Scan(&count)
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM opportunities WHERE company_id = ?", companyID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ContactRepo implements domain.ContactRepository
type ContactRepo struct {
	db *sql.DB
}

func NewContactRepo(db *sql.DB) *ContactRepo {
	return &ContactRepo{db: db}
}

func (r *ContactRepo) Create(ctx context.Context, c *domain.Contact) error {
	res, err := r.db.ExecContext(ctx, "INSERT INTO contacts (company_id, name, email, phone, title) VALUES (?, ?, ?, ?, ?)",
		c.CompanyID, c.Name, c.Email, c.Phone, c.Title)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	c.ID = id
	return nil
}

func (r *ContactRepo) GetByID(ctx context.Context, id int64) (*domain.Contact, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT c.id, c.company_id, comp.name, c.name, c.email, c.phone, c.title, c.created_at, c.updated_at
		FROM contacts c
		JOIN companies comp ON c.company_id = comp.id
		WHERE c.id = ?`, id)
	var c domain.Contact
	var phone, title sql.NullString
	var createdAt, updatedAt string
	if err := row.Scan(&c.ID, &c.CompanyID, &c.CompanyName, &c.Name, &c.Email, &phone, &title, &createdAt, &updatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	c.Phone = phone.String
	c.Title = title.String
	return &c, nil
}

func (r *ContactRepo) List(ctx context.Context, companyID *int64) ([]domain.Contact, error) {
	query := `
		SELECT c.id, c.company_id, comp.name, c.name, c.email, c.phone, c.title, c.created_at, c.updated_at
		FROM contacts c
		JOIN companies comp ON c.company_id = comp.id`
	var args []interface{}
	if companyID != nil {
		query += " WHERE c.company_id = ?"
		args = append(args, *companyID)
	}
	query += " ORDER BY c.id ASC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.Contact
	for rows.Next() {
		var c domain.Contact
		var phone, title sql.NullString
		var createdAt, updatedAt string
		if err := rows.Scan(&c.ID, &c.CompanyID, &c.CompanyName, &c.Name, &c.Email, &phone, &title, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		c.Phone = phone.String
		c.Title = title.String
		result = append(result, c)
	}
	return result, nil
}

func (r *ContactRepo) Update(ctx context.Context, c *domain.Contact) error {
	res, err := r.db.ExecContext(ctx, "UPDATE contacts SET name = ?, phone = ?, title = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		c.Name, c.Phone, c.Title, c.ID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// InteractionRepo implements domain.InteractionRepository
type InteractionRepo struct {
	db *sql.DB
}

func NewInteractionRepo(db *sql.DB) *InteractionRepo {
	return &InteractionRepo{db: db}
}

func (r *InteractionRepo) Create(ctx context.Context, i *domain.Interaction) error {
	res, err := r.db.ExecContext(ctx, `
		INSERT INTO interactions (contact_id, user_id, type, interaction_time, summary, next_action, next_action_date)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		i.ContactID, i.UserID, string(i.Type), i.InteractionTime.Format(time.RFC3339), i.Summary, i.NextAction, i.NextActionDate)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	i.ID = id
	return nil
}

func (r *InteractionRepo) ListByContact(ctx context.Context, contactID int64) ([]domain.Interaction, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT i.id, i.contact_id, i.user_id, u.name, i.type, i.interaction_time, i.summary, i.next_action, i.next_action_date, i.created_at
		FROM interactions i
		JOIN users u ON i.user_id = u.id
		WHERE i.contact_id = ?
		ORDER BY i.interaction_time DESC`, contactID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.Interaction
	for rows.Next() {
		var i domain.Interaction
		var itype, itime string
		var nextAction, nextDate sql.NullString
		var createdAt string
		if err := rows.Scan(&i.ID, &i.ContactID, &i.UserID, &i.UserName, &itype, &itime, &i.Summary, &nextAction, &nextDate, &createdAt); err != nil {
			return nil, err
		}
		i.Type = domain.InteractionType(itype)
		i.InteractionTime, _ = time.Parse(time.RFC3339, itime)
		if i.InteractionTime.IsZero() {
			i.InteractionTime, _ = time.Parse("2006-01-02 15:04:05", itime)
		}
		i.NextAction = nextAction.String
		i.NextActionDate = nextDate.String
		result = append(result, i)
	}
	return result, nil
}

func (r *InteractionRepo) GetLatestNextAction(ctx context.Context, contactID int64) (string, string, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT next_action, next_action_date 
		FROM interactions 
		WHERE contact_id = ? AND next_action IS NOT NULL AND next_action != ''
		ORDER BY interaction_time DESC LIMIT 1`, contactID)
	var action, date sql.NullString
	if err := row.Scan(&action, &date); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", nil
		}
		return "", "", err
	}
	return action.String, date.String, nil
}

// OpportunityRepo implements domain.OpportunityRepository
type OpportunityRepo struct {
	db *sql.DB
}

func NewOpportunityRepo(db *sql.DB) *OpportunityRepo {
	return &OpportunityRepo{db: db}
}

func (r *OpportunityRepo) Create(ctx context.Context, o *domain.Opportunity) error {
	res, err := r.db.ExecContext(ctx, `
		INSERT INTO opportunities (name, company_id, contact_id, owner_user_id, stage, amount, expected_close_date, stage_updated_at, stage_updated_by)
		VALUES (?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, ?)`,
		o.Name, o.CompanyID, o.ContactID, o.OwnerUserID, string(o.Stage), o.Amount, o.ExpectedCloseDate, o.StageUpdatedBy)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	o.ID = id
	return nil
}

func (r *OpportunityRepo) GetByID(ctx context.Context, id int64) (*domain.Opportunity, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT o.id, o.name, o.company_id, comp.name, o.contact_id, cont.name, o.owner_user_id, u.name,
		       o.stage, o.amount, o.expected_close_date, o.stage_updated_at, o.stage_updated_by, o.created_at, o.updated_at
		FROM opportunities o
		JOIN companies comp ON o.company_id = comp.id
		JOIN contacts cont ON o.contact_id = cont.id
		JOIN users u ON o.owner_user_id = u.id
		WHERE o.id = ?`, id)
	var o domain.Opportunity
	var stage string
	var stageUpdatedAt, createdAt, updatedAt string
	if err := row.Scan(&o.ID, &o.Name, &o.CompanyID, &o.CompanyName, &o.ContactID, &o.ContactName,
		&o.OwnerUserID, &o.OwnerName, &stage, &o.Amount, &o.ExpectedCloseDate,
		&stageUpdatedAt, &o.StageUpdatedBy, &createdAt, &updatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	o.Stage = domain.SalesStage(stage)
	o.StageUpdatedAt, _ = time.Parse(time.RFC3339, stageUpdatedAt)
	return &o, nil
}

func (r *OpportunityRepo) List(ctx context.Context, ownerUserID *int64, stage *domain.SalesStage) ([]domain.Opportunity, error) {
	query := `
		SELECT o.id, o.name, o.company_id, comp.name, o.contact_id, cont.name, o.owner_user_id, u.name,
		       o.stage, o.amount, o.expected_close_date, o.stage_updated_at, o.stage_updated_by, o.created_at, o.updated_at
		FROM opportunities o
		JOIN companies comp ON o.company_id = comp.id
		JOIN contacts cont ON o.contact_id = cont.id
		JOIN users u ON o.owner_user_id = u.id
		WHERE 1=1`
	var args []interface{}
	if ownerUserID != nil {
		query += " AND o.owner_user_id = ?"
		args = append(args, *ownerUserID)
	}
	if stage != nil {
		query += " AND o.stage = ?"
		args = append(args, string(*stage))
	}
	query += " ORDER BY o.id ASC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.Opportunity
	for rows.Next() {
		var o domain.Opportunity
		var st string
		var stageUpdatedAt, createdAt, updatedAt string
		if err := rows.Scan(&o.ID, &o.Name, &o.CompanyID, &o.CompanyName, &o.ContactID, &o.ContactName,
			&o.OwnerUserID, &o.OwnerName, &st, &o.Amount, &o.ExpectedCloseDate,
			&stageUpdatedAt, &o.StageUpdatedBy, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		o.Stage = domain.SalesStage(st)
		o.StageUpdatedAt, _ = time.Parse(time.RFC3339, stageUpdatedAt)
		result = append(result, o)
	}
	return result, nil
}

func (r *OpportunityRepo) UpdateStage(ctx context.Context, id int64, stage domain.SalesStage, amount *float64, updatedByID int64, updatedByName string) error {
	query := `UPDATE opportunities SET stage = ?, stage_updated_at = CURRENT_TIMESTAMP, stage_updated_by = ?, updated_at = CURRENT_TIMESTAMP`
	args := []interface{}{string(stage), updatedByName}
	if amount != nil {
		query += ", amount = ?"
		args = append(args, *amount)
	}
	query += " WHERE id = ?"
	args = append(args, id)

	res, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *OpportunityRepo) Update(ctx context.Context, o *domain.Opportunity) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE opportunities 
		SET name = ?, stage = ?, amount = ?, expected_close_date = ?, updated_at = CURRENT_TIMESTAMP 
		WHERE id = ?`,
		o.Name, string(o.Stage), o.Amount, o.ExpectedCloseDate, o.ID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}
