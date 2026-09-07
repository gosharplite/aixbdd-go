package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	crmhttp "crm/internal/adapter/http"
	"crm/internal/adapter/repository"
	"crm/internal/domain"
)

func TestAuthMiddlewareAndMeEndpoint(t *testing.T) {
	db, err := repository.InitDB(":memory:")
	if err != nil {
		t.Fatalf("failed to init db: %v", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepo(db)
	router := crmhttp.NewRouter(crmhttp.ServerConfig{
		UserRepo: userRepo,
		UIDir:    "../../specs/plans/001-crm-core/ui",
	})

	// 1. Unauthenticated request to /api/me -> 401
	req := httptest.NewRequest("GET", "/api/me", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}

	// 2. Request with X-User-Name: Carol -> 200 with SalesManager role
	req = httptest.NewRequest("GET", "/api/me", nil)
	req.Header.Set("X-User-Name", "Carol")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestActorPermissions(t *testing.T) {
	managerActor := &domain.CurrentActor{
		User:            &domain.User{ID: 1, Name: "Carol", Role: domain.RoleSalesManager},
		IsAuthenticated: true,
	}
	repAlice := &domain.CurrentActor{
		User:            &domain.User{ID: 2, Name: "Alice", Role: domain.RoleSalesRep},
		IsAuthenticated: true,
	}
	repBob := &domain.CurrentActor{
		User:            &domain.User{ID: 3, Name: "Bob", Role: domain.RoleSalesRep},
		IsAuthenticated: true,
	}

	opp := &domain.Opportunity{
		ID:          10,
		Name:        "Test Opp",
		OwnerUserID: 2,
		OwnerName:   "Alice",
	}

	// Carol (Manager) can advance
	if !managerActor.CanAdvanceOpportunity(opp) {
		t.Errorf("Manager Carol should be able to advance opp")
	}

	// Alice (Owner) can advance
	if !repAlice.CanAdvanceOpportunity(opp) {
		t.Errorf("Owner Alice should be able to advance opp")
	}

	// Bob (Non-owner Rep) cannot advance
	if repBob.CanAdvanceOpportunity(opp) {
		t.Errorf("Non-owner Bob should NOT be able to advance opp")
	}
}
