package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	crmhttp "github.com/gosharplite/aixbdd-go/internal/adapter/http"
	"github.com/gosharplite/aixbdd-go/internal/adapter/repository"
)

func main() {
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "crm.db"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	uiDir := os.Getenv("UI_DIR")
	if uiDir == "" {
		uiDir = "specs/plans/001-crm-core/ui"
	}

	db, err := repository.InitDB(dbPath)
	if err != nil {
		log.Fatalf("failed to init database: %v", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepo(db)
	companyRepo := repository.NewCompanyRepo(db)
	contactRepo := repository.NewContactRepo(db)
	interactionRepo := repository.NewInteractionRepo(db)
	oppRepo := repository.NewOpportunityRepo(db)

	router := crmhttp.NewRouter(crmhttp.ServerConfig{
		UserRepo:        userRepo,
		CompanyRepo:     companyRepo,
		ContactRepo:     contactRepo,
		InteractionRepo: interactionRepo,
		OpportunityRepo: oppRepo,
		UIDir:           uiDir,
	})

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("CRM Server listening on http://localhost:%s", port)
		log.Printf("Serving UI from %s", uiDir)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down CRM server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server forced shutdown: %v", err)
	}
	fmt.Println("CRM server gracefully stopped")
}
