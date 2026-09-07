package http

import (
	"net/http"
	"os"
	"path/filepath"

	"crm/internal/domain"
	"crm/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type ServerConfig struct {
	UserRepo        domain.UserRepository
	CompanyRepo     domain.CompanyRepository
	ContactRepo     domain.ContactRepository
	InteractionRepo domain.InteractionRepository
	OpportunityRepo domain.OpportunityRepository
	UIDir           string
}

func NewRouter(cfg ServerConfig) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(AuthMiddleware(cfg.UserRepo))

	// Services
	customerSvc := service.NewCustomerService(cfg.CompanyRepo, cfg.ContactRepo, cfg.InteractionRepo)
	customerHandler := NewCustomerHandler(customerSvc)

	interactionSvc := service.NewInteractionService(cfg.InteractionRepo, cfg.ContactRepo)
	interactionHandler := NewInteractionHandler(interactionSvc)

	opportunitySvc := service.NewOpportunityService(cfg.OpportunityRepo, cfg.CompanyRepo, cfg.ContactRepo, cfg.UserRepo)
	opportunityHandler := NewOpportunityHandler(opportunitySvc)

	pipelineSvc := service.NewPipelineService(cfg.OpportunityRepo)
	pipelineHandler := NewPipelineHandler(pipelineSvc)

	// API Routes
	r.Route("/api", func(api chi.Router) {
		api.Get("/me", func(w http.ResponseWriter, r *http.Request) {
			actor := GetActor(r)
			if !actor.IsAuthenticated {
				RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "未登入")
				return
			}
			RespondJSON(w, http.StatusOK, actor.User)
		})

		// Companies & Contacts (auth handled inside handlers to return 401 when unauthenticated)
		api.Post("/companies", customerHandler.CreateCompany)
		api.Get("/companies", customerHandler.ListCompanies)
		api.Get("/companies/{id}", customerHandler.GetCompany)
		api.Post("/contacts", customerHandler.CreateContact)
		api.Get("/contacts", customerHandler.ListContacts)
		api.Get("/contacts/{id}", customerHandler.GetContact)
		api.Put("/contacts/{id}", customerHandler.UpdateContact)

		// Interactions
		api.Post("/contacts/{id}/interactions", interactionHandler.CreateInteraction)
		api.Get("/contacts/{id}/interactions", interactionHandler.ListInteractions)

		// Opportunities
		api.Post("/opportunities", opportunityHandler.CreateOpportunity)
		api.Get("/opportunities", opportunityHandler.ListOpportunities)
		api.Get("/opportunities/{id}", opportunityHandler.GetOpportunity)
		api.Patch("/opportunities/{id}/stage", opportunityHandler.UpdateStage)

		// Pipeline
		api.Get("/pipeline", pipelineHandler.GetPipeline)

		// Placeholder groups for subsequent Feature phases
		api.Group(func(protected chi.Router) {
			protected.Use(RequireAuth)

			// Interactions will be hooked up in Feature phase 5
			// Opportunities will be hooked up in Feature phase 6
			// Pipeline will be hooked up in Feature phase 7
		})
	})

	// Static UI file serving
	uiPath := cfg.UIDir
	if uiPath == "" {
		uiPath = "specs/plans/001-crm-core/ui"
	}
	if _, err := os.Stat(uiPath); err == nil {
		fileServer := http.FileServer(http.Dir(uiPath))
		r.Handle("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if file exists in uiPath
			targetPath := filepath.Join(uiPath, filepath.Clean(r.URL.Path))
			if _, err := os.Stat(targetPath); err != nil {
				// Fallback to index.html for SPA/root
				http.ServeFile(w, r, filepath.Join(uiPath, "index.html"))
				return
			}
			fileServer.ServeHTTP(w, r)
		}))
	}

	return r
}
