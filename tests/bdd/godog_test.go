package bdd

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	crmhttp "crm/internal/adapter/http"
	"crm/internal/adapter/repository"
	"crm/internal/domain"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
)

type TestContext struct {
	DB                   *sql.DB
	UserRepo             domain.UserRepository
	CompanyRepo          domain.CompanyRepository
	ContactRepo          domain.ContactRepository
	InteractionRepo      domain.InteractionRepository
	OpportunityRepo      domain.OpportunityRepository
	Router               http.Handler
	LastResponse         *httptest.ResponseRecorder
	CurrentActorName     string
	CurrentActorRole     domain.UserRole
	IsAuthenticated      bool
	Headers              map[string]string
	LastCreatedCompanyID int64
	LastCreatedContactID int64
	LastCreatedOppID     int64
	LastOppCurrentStage  string
}

func (tc *TestContext) Reset() error {
	if tc.DB != nil {
		_ = tc.DB.Close()
	}

	db, err := repository.InitDB(":memory:")
	if err != nil {
		return err
	}
	tc.DB = db
	tc.UserRepo = repository.NewUserRepo(db)
	tc.CompanyRepo = repository.NewCompanyRepo(db)
	tc.ContactRepo = repository.NewContactRepo(db)
	tc.InteractionRepo = repository.NewInteractionRepo(db)
	tc.OpportunityRepo = repository.NewOpportunityRepo(db)

	tc.Router = crmhttp.NewRouter(crmhttp.ServerConfig{
		UserRepo:        tc.UserRepo,
		CompanyRepo:     tc.CompanyRepo,
		ContactRepo:     tc.ContactRepo,
		InteractionRepo: tc.InteractionRepo,
		OpportunityRepo: tc.OpportunityRepo,
	})

	tc.LastResponse = httptest.NewRecorder()
	tc.CurrentActorName = ""
	tc.CurrentActorRole = domain.RoleSalesRep
	tc.IsAuthenticated = false
	tc.Headers = make(map[string]string)
	tc.LastCreatedCompanyID = 0
	tc.LastCreatedContactID = 0
	tc.LastCreatedOppID = 0
	tc.LastOppCurrentStage = ""

	return nil
}

func (tc *TestContext) DoRequest(req *http.Request) {
	if tc.IsAuthenticated {
		if tc.CurrentActorName != "" {
			req.Header.Set("X-User-Name", tc.CurrentActorName)
		}
		if tc.CurrentActorRole != "" {
			req.Header.Set("X-User-Role", string(tc.CurrentActorRole))
		}
	}
	for k, v := range tc.Headers {
		req.Header.Set(k, v)
	}

	tc.LastResponse = httptest.NewRecorder()
	tc.Router.ServeHTTP(tc.LastResponse, req)
}

var globalTestContext = &TestContext{}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		err := globalTestContext.Reset()
		return ctx, err
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		if globalTestContext.DB != nil {
			_ = globalTestContext.DB.Close()
		}
		return ctx, nil
	})

	// Step definition registrations will be linked here in Phase 3
	RegisterSharedSteps(ctx, globalTestContext)
	RegisterCustomerSteps(ctx, globalTestContext)
	RegisterInteractionSteps(ctx, globalTestContext)
	RegisterOpportunitySteps(ctx, globalTestContext)
	RegisterPipelineSteps(ctx, globalTestContext)
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../../specs/truth/features/backend"},
			Output:   colors.Colored(os.Stdout),
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
