package repository_test

import (
	"context"
	"testing"
	"time"

	"crm/internal/adapter/repository"
	"crm/internal/domain"
)

func TestSQLiteRepositories(t *testing.T) {
	ctx := context.Background()
	db, err := repository.InitDB(":memory:")
	if err != nil {
		t.Fatalf("failed to init in-memory db: %v", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepo(db)
	companyRepo := repository.NewCompanyRepo(db)
	contactRepo := repository.NewContactRepo(db)
	interactionRepo := repository.NewInteractionRepo(db)
	oppRepo := repository.NewOpportunityRepo(db)

	// 1. Verify seed users
	carol, err := userRepo.GetByUsername(ctx, "carol")
	if err != nil {
		t.Fatalf("failed to get seed user carol: %v", err)
	}
	if carol.Role != domain.RoleSalesManager {
		t.Errorf("expected Carol to be SalesManager, got %s", carol.Role)
	}

	// 2. Company CRUD
	comp := &domain.Company{
		Name:     "水球軟體",
		Industry: "軟體資訊",
		Address:  "台北市信義區",
	}
	if err := companyRepo.Create(ctx, comp); err != nil {
		t.Fatalf("failed to create company: %v", err)
	}
	if comp.ID == 0 {
		t.Errorf("expected company ID > 0")
	}

	// 3. Contact CRUD
	contact := &domain.Contact{
		CompanyID: comp.ID,
		Name:      "王大明",
		Email:     "daming.wang@test.com",
		Phone:     "0912-345-678",
		Title:     "技術長",
	}
	if err := contactRepo.Create(ctx, contact); err != nil {
		t.Fatalf("failed to create contact: %v", err)
	}

	// 4. Interaction CRUD
	inter := &domain.Interaction{
		ContactID:       contact.ID,
		UserID:          carol.ID,
		Type:            domain.InteractionCall,
		InteractionTime: time.Now(),
		Summary:         "需求訪談完成",
		NextAction:      "提供正式提案",
		NextActionDate:  "2026-09-05",
	}
	if err := interactionRepo.Create(ctx, inter); err != nil {
		t.Fatalf("failed to create interaction: %v", err)
	}

	// 5. Opportunity CRUD
	opp := &domain.Opportunity{
		Name:              "企業版合約",
		CompanyID:         comp.ID,
		ContactID:         contact.ID,
		OwnerUserID:       carol.ID,
		Stage:             domain.StageProspecting,
		Amount:            300000,
		ExpectedCloseDate: "2026-06-30",
		StageUpdatedBy:    "Carol",
	}
	if err := oppRepo.Create(ctx, opp); err != nil {
		t.Fatalf("failed to create opportunity: %v", err)
	}
}
