package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gosharplite/aixbdd-go/internal/domain"
)

type InteractionService struct {
	interactionRepo domain.InteractionRepository
	contactRepo     domain.ContactRepository
}

func NewInteractionService(
	interactionRepo domain.InteractionRepository,
	contactRepo domain.ContactRepository,
) *InteractionService {
	return &InteractionService{
		interactionRepo: interactionRepo,
		contactRepo:     contactRepo,
	}
}

type CreateInteractionRequest struct {
	Type            string `json:"type"`
	InteractionTime string `json:"interaction_time"`
	Summary         string `json:"summary"`
	NextAction      string `json:"next_action,omitempty"`
	NextActionDate  string `json:"next_action_date,omitempty"` // YYYY-MM-DD
}

func (s *InteractionService) CreateInteraction(ctx context.Context, actor *domain.CurrentActor, contactID int64, req CreateInteractionRequest) (*domain.Interaction, error) {
	if actor == nil || !actor.IsAuthenticated {
		return nil, domain.ErrUnauthorized
	}
	if contactID <= 0 {
		return nil, fmt.Errorf("%w: 無效的聯絡人識別碼", domain.ErrValidation)
	}
	_, err := s.contactRepo.GetByID(ctx, contactID)
	if err != nil {
		return nil, fmt.Errorf("contact not found: %w", err)
	}

	summary := strings.TrimSpace(req.Summary)
	if summary == "" {
		return nil, fmt.Errorf("%w: 聯絡內容摘要為必填", domain.ErrValidation)
	}

	itype := domain.InteractionType(strings.TrimSpace(req.Type))
	if itype != domain.InteractionCall && itype != domain.InteractionMeeting && itype != domain.InteractionEmail {
		itype = domain.InteractionCall
	}

	t := time.Now()
	if req.InteractionTime != "" {
		parsed, err := time.Parse(time.RFC3339, req.InteractionTime)
		if err == nil {
			t = parsed
		} else {
			parsed, err = time.Parse("2006-01-02 15:04", req.InteractionTime)
			if err == nil {
				t = parsed
			}
		}
	}

	userID := int64(1)
	if actor.User != nil {
		userID = actor.User.ID
	}

	inter := &domain.Interaction{
		ContactID:       contactID,
		UserID:          userID,
		UserName:        actor.User.Name,
		Type:            itype,
		InteractionTime: t,
		Summary:         summary,
		NextAction:      strings.TrimSpace(req.NextAction),
		NextActionDate:  strings.TrimSpace(req.NextActionDate),
	}

	if err := s.interactionRepo.Create(ctx, inter); err != nil {
		return nil, fmt.Errorf("create interaction: %w", err)
	}
	return inter, nil
}

func (s *InteractionService) ListInteractions(ctx context.Context, actor *domain.CurrentActor, contactID int64) ([]domain.Interaction, error) {
	if actor == nil || !actor.IsAuthenticated {
		return nil, domain.ErrUnauthorized
	}
	if contactID <= 0 {
		return nil, fmt.Errorf("%w: 無效的聯絡人識別碼", domain.ErrValidation)
	}
	_, err := s.contactRepo.GetByID(ctx, contactID)
	if err != nil {
		return nil, fmt.Errorf("contact not found: %w", err)
	}
	return s.interactionRepo.ListByContact(ctx, contactID)
}
