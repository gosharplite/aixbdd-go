package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"crm/internal/domain"
)

type OpportunityService struct {
	oppRepo     domain.OpportunityRepository
	companyRepo domain.CompanyRepository
	contactRepo domain.ContactRepository
	userRepo    domain.UserRepository
}

func NewOpportunityService(
	oppRepo domain.OpportunityRepository,
	companyRepo domain.CompanyRepository,
	contactRepo domain.ContactRepository,
	userRepo domain.UserRepository,
) *OpportunityService {
	return &OpportunityService{
		oppRepo:     oppRepo,
		companyRepo: companyRepo,
		contactRepo: contactRepo,
		userRepo:    userRepo,
	}
}

type CreateOpportunityRequest struct {
	Name              string  `json:"name"`
	CompanyID         int64   `json:"company_id"`
	ContactID         int64   `json:"contact_id"`
	Stage             string  `json:"stage,omitempty"`
	Amount            float64 `json:"amount"`
	ExpectedCloseDate string  `json:"expected_close_date"`
}

func (s *OpportunityService) CreateOpportunity(ctx context.Context, actor *domain.CurrentActor, req CreateOpportunityRequest) (*domain.Opportunity, error) {
	if actor == nil || !actor.IsAuthenticated {
		return nil, domain.ErrUnauthorized
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, fmt.Errorf("%w: 機會名稱為必填", domain.ErrValidation)
	}
	if req.Amount < 0 {
		return nil, fmt.Errorf("%w: 預估金額不可為負數", domain.ErrValidation)
	}
	if req.CompanyID <= 0 {
		return nil, fmt.Errorf("%w: 必須指定關聯公司", domain.ErrValidation)
	}
	comp, err := s.companyRepo.GetByID(ctx, req.CompanyID)
	if err != nil {
		return nil, fmt.Errorf("company not found: %w", err)
	}

	contactName := ""
	if req.ContactID > 0 {
		cont, err := s.contactRepo.GetByID(ctx, req.ContactID)
		if err == nil && cont != nil {
			contactName = cont.Name
		}
	}

	stage := domain.StageProspecting
	if req.Stage != "" {
		stage = mapStage(req.Stage)
	}

	expectedDate := strings.TrimSpace(req.ExpectedCloseDate)
	if expectedDate == "" {
		expectedDate = time.Now().AddDate(0, 3, 0).Format("2006-01-02")
	}

	opp := &domain.Opportunity{
		Name:              name,
		CompanyID:         req.CompanyID,
		CompanyName:       comp.Name,
		ContactID:         req.ContactID,
		ContactName:       contactName,
		OwnerUserID:       actor.User.ID,
		OwnerName:         actor.User.Name,
		Stage:             stage,
		Amount:            req.Amount,
		ExpectedCloseDate: expectedDate,
		StageUpdatedBy:    actor.User.Name,
	}

	if err := s.oppRepo.Create(ctx, opp); err != nil {
		return nil, fmt.Errorf("create opportunity: %w", err)
	}
	return opp, nil
}

func (s *OpportunityService) GetOpportunity(ctx context.Context, actor *domain.CurrentActor, id int64) (*domain.Opportunity, error) {
	if actor == nil || !actor.IsAuthenticated {
		return nil, domain.ErrUnauthorized
	}
	return s.oppRepo.GetByID(ctx, id)
}

type UpdateStageRequest struct {
	Stage  string   `json:"stage"`
	Amount *float64 `json:"amount,omitempty"`
}

func (s *OpportunityService) UpdateStage(ctx context.Context, actor *domain.CurrentActor, id int64, req UpdateStageRequest) (*domain.Opportunity, error) {
	if actor == nil || !actor.IsAuthenticated {
		return nil, domain.ErrUnauthorized
	}
	opp, err := s.oppRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// RBAC Enforcement: Only owner or SalesManager can advance stage or update amount
	if !actor.CanAdvanceOpportunity(opp) {
		return nil, domain.ErrForbidden
	}

	newStage := mapStage(req.Stage)
	if err := s.oppRepo.UpdateStage(ctx, id, newStage, req.Amount, actor.User.ID, actor.User.Name); err != nil {
		return nil, fmt.Errorf("update opportunity stage: %w", err)
	}

	return s.oppRepo.GetByID(ctx, id)
}

func (s *OpportunityService) ListOpportunities(ctx context.Context, actor *domain.CurrentActor, stage *domain.SalesStage, ownerID *int64) ([]domain.Opportunity, error) {
	if actor == nil || !actor.IsAuthenticated {
		return nil, domain.ErrUnauthorized
	}
	// Role data isolation
	var targetOwner *int64 = ownerID
	if actor.User.Role != domain.RoleSalesManager {
		// SalesRep can only see their own opportunities
		targetOwner = &actor.User.ID
	}
	return s.oppRepo.List(ctx, targetOwner, stage)
}

func mapStage(st string) domain.SalesStage {
	switch st {
	case "潛在", "Prospecting":
		return domain.StageProspecting
	case "已聯繫", "Contacted":
		return domain.StageContacted
	case "提案中", "Proposal":
		return domain.StageProposal
	case "已成交", "ClosedWon":
		return domain.StageClosedWon
	case "已失敗", "ClosedLost":
		return domain.StageClosedLost
	default:
		return domain.SalesStage(st)
	}
}
