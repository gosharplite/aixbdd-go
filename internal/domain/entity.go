package domain

import (
	"context"
	"errors"
	"regexp"
	"time"
)

// UserRole represents RBAC role in CRM
type UserRole string

const (
	RoleSalesManager UserRole = "SalesManager"
	RoleSalesRep     UserRole = "SalesRep"
)

// SalesStage represents standard 5 sales stages
type SalesStage string

const (
	StageProspecting SalesStage = "Prospecting" // 1. 潛在
	StageContacted   SalesStage = "Contacted"   // 2. 已聯繫
	StageProposal    SalesStage = "Proposal"    // 3. 提案中
	StageClosedWon   SalesStage = "ClosedWon"   // 4. 已成交
	StageClosedLost  SalesStage = "ClosedLost"  // 5. 已失敗
)

// IsClosedStage checks if stage is in terminal closed state
func IsClosedStage(stage SalesStage) bool {
	return stage == StageClosedWon || stage == StageClosedLost
}

// InteractionType represents contact channel
type InteractionType string

const (
	InteractionCall    InteractionType = "Call"
	InteractionMeeting InteractionType = "Meeting"
	InteractionEmail   InteractionType = "Email"
)

var (
	ErrNotFound             = errors.New("resource not found")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrForbidden            = errors.New("forbidden: insufficient permissions")
	ErrValidation           = errors.New("validation failed")
	ErrDeleteRestricted     = errors.New("cannot delete: related records exist")
	ErrStageAlreadyTerminal = errors.New("opportunity is already closed")
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func IsValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// User entity
type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	Role      UserRole  `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// Company entity
type Company struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	TaxID     string    `json:"tax_id,omitempty"`
	Industry  string    `json:"industry,omitempty"`
	Address   string    `json:"address,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Contact entity
type Contact struct {
	ID                   int64     `json:"id"`
	CompanyID            int64     `json:"company_id"`
	CompanyName          string    `json:"company_name,omitempty"`
	Name                 string    `json:"name"`
	Email                string    `json:"email"`
	Phone                string    `json:"phone,omitempty"`
	Title                string    `json:"title,omitempty"`
	LatestNextAction     string    `json:"latest_next_action,omitempty"`
	LatestNextActionDate string    `json:"latest_next_action_date,omitempty"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// Interaction entity
type Interaction struct {
	ID              int64           `json:"id"`
	ContactID       int64           `json:"contact_id"`
	UserID          int64           `json:"user_id"`
	UserName        string          `json:"user_name,omitempty"`
	Type            InteractionType `json:"type"`
	InteractionTime time.Time       `json:"interaction_time"`
	Summary         string          `json:"summary"`
	NextAction      string          `json:"next_action,omitempty"`
	NextActionDate  string          `json:"next_action_date,omitempty"` // YYYY-MM-DD
	CreatedAt       time.Time       `json:"created_at"`
}

// Opportunity entity
type Opportunity struct {
	ID                int64      `json:"id"`
	Name              string     `json:"name"`
	CompanyID         int64      `json:"company_id"`
	CompanyName       string     `json:"company_name,omitempty"`
	ContactID         int64      `json:"contact_id"`
	ContactName       string     `json:"contact_name,omitempty"`
	OwnerUserID       int64      `json:"owner_user_id"`
	OwnerName         string     `json:"owner_name,omitempty"`
	Stage             SalesStage `json:"stage"`
	Amount            float64    `json:"amount"`
	ExpectedCloseDate string     `json:"expected_close_date"` // YYYY-MM-DD
	StageUpdatedAt    time.Time  `json:"stage_updated_at"`
	StageUpdatedBy    string     `json:"stage_updated_by,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// PipelineSummary provides grouped stages overview
type PipelineSummary struct {
	ViewerRole            UserRole       `json:"viewer_role"`
	Stages                []StageSummary `json:"stages"`
	ActiveTotalAmount     float64        `json:"active_total_amount"`
	TotalOpportunityCount int            `json:"total_opportunity_count"`
}

// StageSummary represents metric for single sales stage
type StageSummary struct {
	Stage         SalesStage    `json:"stage"`
	Count         int           `json:"count"`
	TotalAmount   float64       `json:"total_amount"`
	Opportunities []Opportunity `json:"opportunities"`
}

// CurrentActor provides context about user executing the request
type CurrentActor struct {
	User            *User
	IsAuthenticated bool
}

func (a *CurrentActor) CanManageCompany() bool {
	return a.IsAuthenticated
}

func (a *CurrentActor) CanAdvanceOpportunity(opp *Opportunity) bool {
	if !a.IsAuthenticated || opp == nil {
		return false
	}
	if a.User.Role == RoleSalesManager {
		return true
	}
	return a.User.ID == opp.OwnerUserID || a.User.Name == opp.OwnerName
}

// Repository Interfaces
type UserRepository interface {
	GetByID(ctx context.Context, id int64) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetByName(ctx context.Context, name string) (*User, error)
	Create(ctx context.Context, user *User) error
}

type CompanyRepository interface {
	Create(ctx context.Context, c *Company) error
	GetByID(ctx context.Context, id int64) (*Company, error)
	GetByName(ctx context.Context, name string) (*Company, error)
	List(ctx context.Context) ([]Company, error)
	Update(ctx context.Context, c *Company) error
	HasActiveRelations(ctx context.Context, companyID int64) (bool, error)
}

type ContactRepository interface {
	Create(ctx context.Context, c *Contact) error
	GetByID(ctx context.Context, id int64) (*Contact, error)
	List(ctx context.Context, companyID *int64) ([]Contact, error)
	Update(ctx context.Context, c *Contact) error
}

type InteractionRepository interface {
	Create(ctx context.Context, i *Interaction) error
	ListByContact(ctx context.Context, contactID int64) ([]Interaction, error)
	GetLatestNextAction(ctx context.Context, contactID int64) (action string, date string, err error)
}

type OpportunityRepository interface {
	Create(ctx context.Context, o *Opportunity) error
	GetByID(ctx context.Context, id int64) (*Opportunity, error)
	List(ctx context.Context, ownerUserID *int64, stage *SalesStage) ([]Opportunity, error)
	UpdateStage(ctx context.Context, id int64, stage SalesStage, amount *float64, updatedByID int64, updatedByName string) error
	Update(ctx context.Context, o *Opportunity) error
}
