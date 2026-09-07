package service

import (
	"context"
	"fmt"
	"strings"

	"crm/internal/domain"
)

type CustomerService struct {
	companyRepo     domain.CompanyRepository
	contactRepo     domain.ContactRepository
	interactionRepo domain.InteractionRepository
}

func NewCustomerService(
	companyRepo domain.CompanyRepository,
	contactRepo domain.ContactRepository,
	interactionRepo domain.InteractionRepository,
) *CustomerService {
	return &CustomerService{
		companyRepo:     companyRepo,
		contactRepo:     contactRepo,
		interactionRepo: interactionRepo,
	}
}

type CreateCompanyRequest struct {
	Name     string `json:"name"`
	TaxID    string `json:"tax_id,omitempty"`
	Industry string `json:"industry,omitempty"`
	Address  string `json:"address,omitempty"`
}

func (s *CustomerService) CreateCompany(ctx context.Context, actor *domain.CurrentActor, req CreateCompanyRequest) (*domain.Company, error) {
	if actor == nil || !actor.IsAuthenticated {
		return nil, domain.ErrUnauthorized
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, fmt.Errorf("%w: 公司名稱為必填", domain.ErrValidation)
	}

	comp := &domain.Company{
		Name:     name,
		TaxID:    strings.TrimSpace(req.TaxID),
		Industry: strings.TrimSpace(req.Industry),
		Address:  strings.TrimSpace(req.Address),
	}
	if err := s.companyRepo.Create(ctx, comp); err != nil {
		return nil, fmt.Errorf("create company: %w", err)
	}
	return comp, nil
}

func (s *CustomerService) GetCompany(ctx context.Context, actor *domain.CurrentActor, id int64) (*domain.Company, []domain.Contact, error) {
	if actor == nil || !actor.IsAuthenticated {
		return nil, nil, domain.ErrUnauthorized
	}
	comp, err := s.companyRepo.GetByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	contacts, err := s.contactRepo.List(ctx, &id)
	if err != nil {
		return nil, nil, err
	}
	return comp, contacts, nil
}

func (s *CustomerService) ListCompanies(ctx context.Context, actor *domain.CurrentActor) ([]domain.Company, error) {
	if actor == nil || !actor.IsAuthenticated {
		return nil, domain.ErrUnauthorized
	}
	return s.companyRepo.List(ctx)
}

func (s *CustomerService) UpdateCompany(ctx context.Context, actor *domain.CurrentActor, id int64, req CreateCompanyRequest) (*domain.Company, error) {
	if actor == nil || !actor.IsAuthenticated {
		return nil, domain.ErrUnauthorized
	}
	comp, err := s.companyRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		comp.Name = strings.TrimSpace(req.Name)
	}
	if req.TaxID != "" {
		comp.TaxID = strings.TrimSpace(req.TaxID)
	}
	if req.Industry != "" {
		comp.Industry = strings.TrimSpace(req.Industry)
	}
	if req.Address != "" {
		comp.Address = strings.TrimSpace(req.Address)
	}
	if err := s.companyRepo.Update(ctx, comp); err != nil {
		return nil, err
	}
	return comp, nil
}

type CreateContactRequest struct {
	CompanyID int64  `json:"company_id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone,omitempty"`
	Title     string `json:"title,omitempty"`
}

func (s *CustomerService) CreateContact(ctx context.Context, actor *domain.CurrentActor, req CreateContactRequest) (*domain.Contact, error) {
	if actor == nil || !actor.IsAuthenticated {
		return nil, domain.ErrUnauthorized
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, fmt.Errorf("%w: 聯絡人姓名為必填", domain.ErrValidation)
	}
	email := strings.TrimSpace(req.Email)
	if email == "" || !domain.IsValidEmail(email) {
		return nil, fmt.Errorf("%w: 請輸入有效的電子郵件格式", domain.ErrValidation)
	}
	if req.CompanyID <= 0 {
		return nil, fmt.Errorf("%w: 必須指定所屬公司", domain.ErrValidation)
	}
	comp, err := s.companyRepo.GetByID(ctx, req.CompanyID)
	if err != nil {
		return nil, fmt.Errorf("company not found: %w", err)
	}

	contact := &domain.Contact{
		CompanyID:   req.CompanyID,
		CompanyName: comp.Name,
		Name:        name,
		Email:       email,
		Phone:       strings.TrimSpace(req.Phone),
		Title:       strings.TrimSpace(req.Title),
	}
	if err := s.contactRepo.Create(ctx, contact); err != nil {
		return nil, fmt.Errorf("create contact: %w", err)
	}
	return contact, nil
}

func (s *CustomerService) GetContact(ctx context.Context, actor *domain.CurrentActor, id int64) (*domain.Contact, error) {
	if actor == nil || !actor.IsAuthenticated {
		return nil, domain.ErrUnauthorized
	}
	contact, err := s.contactRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if s.interactionRepo != nil {
		action, date, _ := s.interactionRepo.GetLatestNextAction(ctx, id)
		contact.LatestNextAction = action
		contact.LatestNextActionDate = date
	}
	return contact, nil
}

func (s *CustomerService) ListContacts(ctx context.Context, actor *domain.CurrentActor, companyID *int64) ([]domain.Contact, error) {
	if actor == nil || !actor.IsAuthenticated {
		return nil, domain.ErrUnauthorized
	}
	return s.contactRepo.List(ctx, companyID)
}

type UpdateContactRequest struct {
	Name  string `json:"name,omitempty"`
	Phone string `json:"phone,omitempty"`
	Title string `json:"title,omitempty"`
}

func (s *CustomerService) UpdateContact(ctx context.Context, actor *domain.CurrentActor, id int64, req UpdateContactRequest) (*domain.Contact, error) {
	if actor == nil || !actor.IsAuthenticated {
		return nil, domain.ErrUnauthorized
	}
	contact, err := s.contactRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		contact.Name = strings.TrimSpace(req.Name)
	}
	if req.Phone != "" {
		contact.Phone = strings.TrimSpace(req.Phone)
	}
	if req.Title != "" {
		contact.Title = strings.TrimSpace(req.Title)
	}
	if err := s.contactRepo.Update(ctx, contact); err != nil {
		return nil, err
	}
	return contact, nil
}
