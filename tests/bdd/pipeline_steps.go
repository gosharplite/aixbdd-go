package bdd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"

	"github.com/gosharplite/aixbdd-go/internal/domain"
	"github.com/cucumber/godog"
)

func RegisterPipelineSteps(ctx *godog.ScenarioContext, tc *TestContext) {
	ctx.Step(`^CRM 系統中已存在以下銷售機會資料：$`, func(table *godog.Table) error {
		tableRows := table.Rows
		if len(tableRows) > 0 && tableRows[0].Cells[0].Value == "機會名稱" {
			tableRows = tableRows[1:]
		}

		for _, row := range tableRows {
			oppName := row.Cells[0].Value
			compName := row.Cells[1].Value
			ownerName := row.Cells[2].Value
			stageStr := row.Cells[3].Value
			amtStr := row.Cells[4].Value
			amt, _ := strconv.ParseFloat(amtStr, 64)

			comp, _ := tc.CompanyRepo.GetByName(context.Background(), compName)
			if comp == nil {
				comp = &domain.Company{Name: compName, Industry: "科技軟體"}
				_ = tc.CompanyRepo.Create(context.Background(), comp)
			}

			user, _ := tc.UserRepo.GetByName(context.Background(), ownerName)
			var userID int64 = 2
			if user != nil {
				userID = user.ID
			} else {
				u := &domain.User{Username: ownerName, Name: ownerName, Role: domain.RoleSalesRep}
				_ = tc.UserRepo.Create(context.Background(), u)
				userID = u.ID
			}

			contact := &domain.Contact{CompanyID: comp.ID, Name: "窗口", Email: "contact@test.com"}
			_ = tc.ContactRepo.Create(context.Background(), contact)

			opp := &domain.Opportunity{
				Name:              oppName,
				CompanyID:         comp.ID,
				ContactID:         contact.ID,
				OwnerUserID:       userID,
				OwnerName:         ownerName,
				Stage:             mapStageName(stageStr),
				Amount:            amt,
				ExpectedCloseDate: "2026-06-30",
				StageUpdatedBy:    ownerName,
			}
			if err := tc.OpportunityRepo.Create(context.Background(), opp); err != nil {
				return err
			}
		}
		return nil
	})

	ctx.Step(`^"([^"]*)" 查看全公司銷售管線看板$`, func(actor string) error {
		tc.CurrentActorName = actor
		tc.CurrentActorRole = domain.RoleSalesManager
		tc.IsAuthenticated = true
		req := httptest.NewRequest("GET", "/api/pipeline", nil)
		tc.DoRequest(req)
		return nil
	})

	ctx.Step(`^"([^"]*)" 查看銷售管線$`, func(actor string) error {
		tc.CurrentActorName = actor
		tc.CurrentActorRole = domain.RoleSalesRep
		tc.IsAuthenticated = true
		req := httptest.NewRequest("GET", "/api/pipeline", nil)
		tc.DoRequest(req)
		return nil
	})

	ctx.Step(`^系統在銷售管線中呈現包含 "([^"]*)" 與 "([^"]*)" 負責的所有銷售機會$`, func(owner1, owner2 string) error {
		if tc.LastResponse.Code != http.StatusOK {
			return fmt.Errorf("expected 200 OK, got %d", tc.LastResponse.Code)
		}
		var summary domain.PipelineSummary
		if err := json.Unmarshal(tc.LastResponse.Body.Bytes(), &summary); err != nil {
			return err
		}

		found1 := false
		found2 := false
		for _, st := range summary.Stages {
			for _, o := range st.Opportunities {
				if o.OwnerName == owner1 {
					found1 = true
				}
				if o.OwnerName == owner2 {
					found2 = true
				}
			}
		}
		if !found1 || !found2 {
			return fmt.Errorf("pipeline did not contain opportunities from both %s and %s", owner1, owner2)
		}
		return nil
	})

	ctx.Step(`^各階段統計指標如下：$`, func(table *godog.Table) error {
		var summary domain.PipelineSummary
		if err := json.Unmarshal(tc.LastResponse.Body.Bytes(), &summary); err != nil {
			return err
		}

		stageMap := make(map[domain.SalesStage]domain.StageSummary)
		for _, st := range summary.Stages {
			stageMap[st.Stage] = st
		}

		tableRows := table.Rows
		if len(tableRows) > 0 && tableRows[0].Cells[0].Value == "階段" {
			tableRows = tableRows[1:]
		}

		for _, row := range tableRows {
			stageName := row.Cells[0].Value
			expectedCount, _ := strconv.Atoi(row.Cells[1].Value)
			expectedAmt, _ := strconv.ParseFloat(row.Cells[2].Value, 64)

			stage := mapStageName(stageName)
			actual, ok := stageMap[stage]
			if !ok {
				return fmt.Errorf("stage %s not found in pipeline response", stageName)
			}
			if actual.Count != expectedCount {
				return fmt.Errorf("stage %s count mismatch: expected %d, got %d", stageName, expectedCount, actual.Count)
			}
			if actual.TotalAmount != expectedAmt {
				return fmt.Errorf("stage %s amount mismatch: expected %f, got %f", stageName, expectedAmt, actual.TotalAmount)
			}
		}
		return nil
	})

	ctx.Step(`^進行中階段（潛在、已聯繫、提案中）總加總預估金額為 (\d+) 元$`, func(expectedSum int) error {
		var summary domain.PipelineSummary
		if err := json.Unmarshal(tc.LastResponse.Body.Bytes(), &summary); err != nil {
			return err
		}
		if summary.ActiveTotalAmount != float64(expectedSum) {
			return fmt.Errorf("expected active total amount %d, got %f", expectedSum, summary.ActiveTotalAmount)
		}
		return nil
	})

	ctx.Step(`^系統僅顯示負責人為 "([^"]*)" 的銷售機會如下：$`, func(owner string, table *godog.Table) error {
		if tc.LastResponse.Code != http.StatusOK {
			return fmt.Errorf("expected 200 OK, got %d", tc.LastResponse.Code)
		}
		var summary domain.PipelineSummary
		if err := json.Unmarshal(tc.LastResponse.Body.Bytes(), &summary); err != nil {
			return err
		}

		tableRows := table.Rows
		if len(tableRows) > 0 && tableRows[0].Cells[0].Value == "機會名稱" {
			tableRows = tableRows[1:]
		}

		var allOpps []domain.Opportunity
		for _, st := range summary.Stages {
			for _, o := range st.Opportunities {
				if o.OwnerName != owner {
					return fmt.Errorf("found opportunity with unauthorized owner %s: %s", o.OwnerName, o.Name)
				}
				allOpps = append(allOpps, o)
			}
		}

		if len(allOpps) != len(tableRows) {
			return fmt.Errorf("expected %d opportunities for %s, got %d", len(tableRows), owner, len(allOpps))
		}

		return nil
	})

	ctx.Step(`^不應顯示負責人為 "([^"]*)" 的銷售機會$`, func(unauthorizedOwner string) error {
		var summary domain.PipelineSummary
		if err := json.Unmarshal(tc.LastResponse.Body.Bytes(), &summary); err != nil {
			return err
		}
		for _, st := range summary.Stages {
			for _, o := range st.Opportunities {
				if o.OwnerName == unauthorizedOwner {
					return fmt.Errorf("unauthorized opportunity from %s was displayed: %s", unauthorizedOwner, o.Name)
				}
			}
		}
		return nil
	})
}
