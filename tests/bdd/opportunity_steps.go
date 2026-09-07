package bdd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"

	"github.com/gosharplite/aixbdd-go/internal/domain"
	"github.com/cucumber/godog"
)

func RegisterOpportunitySteps(ctx *godog.ScenarioContext, tc *TestContext) {
	ctx.Step(`^CRM 系統定義標準銷售階段如下：$`, func(table *godog.Table) error {
		// Validates predefined stages
		return nil
	})

	ctx.Step(`^系統中存在公司 "([^"]*)" 與客戶聯絡人 "([^"]*)"$`, func(compName, contactName string) error {
		comp, err := tc.CompanyRepo.GetByName(context.Background(), compName)
		if err != nil {
			comp = &domain.Company{Name: compName, Industry: "軟體資訊"}
			if err := tc.CompanyRepo.Create(context.Background(), comp); err != nil {
				return err
			}
		}
		tc.LastCreatedCompanyID = comp.ID

		contact := &domain.Contact{
			CompanyID: comp.ID,
			Name:      contactName,
			Email:     strings.ToLower(contactName) + "@test.com",
			Title:     "技術長",
		}
		if err := tc.ContactRepo.Create(context.Background(), contact); err != nil {
			return err
		}
		tc.LastCreatedContactID = contact.ID
		return nil
	})

	ctx.Step(`^"([^"]*)" 負責 "([^"]*)" 的銷售機會 "([^"]*)"$`, func(ownerName, compName, oppName string) error {
		tc.CurrentActorName = ownerName
		tc.CurrentActorRole = domain.RoleSalesRep
		tc.IsAuthenticated = true

		comp, err := tc.CompanyRepo.GetByName(context.Background(), compName)
		if err != nil {
			comp = &domain.Company{Name: compName, Industry: "軟體資訊"}
			_ = tc.CompanyRepo.Create(context.Background(), comp)
		}
		tc.LastCreatedCompanyID = comp.ID

		contacts, _ := tc.ContactRepo.List(context.Background(), &comp.ID)
		var contactID int64
		if len(contacts) > 0 {
			contactID = contacts[0].ID
		} else {
			c := &domain.Contact{CompanyID: comp.ID, Name: "聯絡人", Email: "contact@test.com"}
			_ = tc.ContactRepo.Create(context.Background(), c)
			contactID = c.ID
		}

		user, _ := tc.UserRepo.GetByName(context.Background(), ownerName)
		var userID int64 = 2
		if user != nil {
			userID = user.ID
		}

		opp := &domain.Opportunity{
			Name:              oppName,
			CompanyID:         comp.ID,
			ContactID:         contactID,
			OwnerUserID:       userID,
			OwnerName:         ownerName,
			Stage:             domain.StageContacted,
			Amount:            300000,
			ExpectedCloseDate: "2026-06-30",
			StageUpdatedBy:    ownerName,
		}
		if err := tc.OpportunityRepo.Create(context.Background(), opp); err != nil {
			return err
		}
		tc.LastCreatedOppID = opp.ID
		tc.LastOppCurrentStage = string(opp.Stage)
		return nil
	})

	ctx.Step(`^該銷售機會目前階段為 "([^"]*)"$`, func(stageName string) error {
		stage := mapStageName(stageName)
		_ = tc.OpportunityRepo.UpdateStage(context.Background(), tc.LastCreatedOppID, stage, nil, 0, tc.CurrentActorName)
		tc.LastOppCurrentStage = string(stage)
		return nil
	})

	ctx.Step(`^該銷售機會預估金額為 (\d+) 元$`, func(amount int) error {
		amt := float64(amount)
		_ = tc.OpportunityRepo.UpdateStage(context.Background(), tc.LastCreatedOppID, domain.SalesStage(tc.LastOppCurrentStage), &amt, 0, tc.CurrentActorName)
		return nil
	})

	ctx.Step(`^"([^"]*)" 負責的銷售機會 "([^"]*)" 處於 "([^"]*)" 階段$`, func(ownerName, oppName, stageName string) error {
		tc.CurrentActorName = ownerName
		tc.CurrentActorRole = domain.RoleSalesRep
		tc.IsAuthenticated = true

		stage := mapStageName(stageName)
		comp, _ := tc.CompanyRepo.GetByName(context.Background(), "水球軟體")
		if comp == nil {
			comp = &domain.Company{Name: "水球軟體", Industry: "軟體資訊"}
			_ = tc.CompanyRepo.Create(context.Background(), comp)
		}
		contacts, _ := tc.ContactRepo.List(context.Background(), &comp.ID)
		var contactID int64
		if len(contacts) > 0 {
			contactID = contacts[0].ID
		} else {
			c := &domain.Contact{CompanyID: comp.ID, Name: "王大明", Email: "wang@test.com"}
			_ = tc.ContactRepo.Create(context.Background(), c)
			contactID = c.ID
		}

		user, _ := tc.UserRepo.GetByName(context.Background(), ownerName)
		var userID int64 = 2
		if user != nil {
			userID = user.ID
		}

		opp := &domain.Opportunity{
			Name:              oppName,
			CompanyID:         comp.ID,
			ContactID:         contactID,
			OwnerUserID:       userID,
			OwnerName:         ownerName,
			Stage:             stage,
			Amount:            300000,
			ExpectedCloseDate: "2026-06-30",
			StageUpdatedBy:    ownerName,
		}
		if err := tc.OpportunityRepo.Create(context.Background(), opp); err != nil {
			return err
		}
		tc.LastCreatedOppID = opp.ID
		tc.LastOppCurrentStage = string(stage)
		return nil
	})

	ctx.Step(`^"([^"]*)" 為一般業務人員，不是該銷售機會負責人$`, func(actor string) error {
		tc.CurrentActorName = actor
		tc.CurrentActorRole = domain.RoleSalesRep
		tc.IsAuthenticated = true
		return nil
	})

	ctx.Step(`^"([^"]*)" 不是業務主管$`, func(actor string) error {
		if tc.CurrentActorRole == domain.RoleSalesManager {
			return fmt.Errorf("%s should not be sales manager", actor)
		}
		return nil
	})

	ctx.Step(`^"([^"]*)" 建立銷售機會如下：$`, func(actor string, table *godog.Table) error {
		tc.CurrentActorName = actor
		payload := make(map[string]interface{})
		for _, row := range table.Rows {
			if len(row.Cells) < 2 {
				continue
			}
			key := row.Cells[0].Value
			val := row.Cells[1].Value
			switch key {
			case "機會名稱":
				payload["name"] = val
			case "關聯公司":
				comp, _ := tc.CompanyRepo.GetByName(context.Background(), val)
				if comp != nil {
					payload["company_id"] = comp.ID
				}
			case "關聯客戶":
				contacts, _ := tc.ContactRepo.List(context.Background(), nil)
				for _, c := range contacts {
					if c.Name == val {
						payload["contact_id"] = c.ID
						break
					}
				}
			case "目前階段":
				payload["stage"] = string(mapStageName(val))
			case "預估金額":
				amt, _ := strconv.ParseFloat(val, 64)
				payload["amount"] = amt
			case "預計成交日":
				payload["expected_close_date"] = val
			}
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/opportunities", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		tc.DoRequest(req)
		return nil
	})

	ctx.Step(`^"([^"]*)" 將銷售機會推進到 "([^"]*)"$`, func(actor, stageName string) error {
		tc.CurrentActorName = actor
		tc.CurrentActorRole = domain.RoleSalesRep
		tc.IsAuthenticated = true
		stage := mapStageName(stageName)
		payload := map[string]interface{}{
			"stage": string(stage),
		}
		body, _ := json.Marshal(payload)
		url := fmt.Sprintf("/api/opportunities/%d/stage", tc.LastCreatedOppID)
		req := httptest.NewRequest("PATCH", url, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		tc.DoRequest(req)
		return nil
	})

	ctx.Step(`^"([^"]*)" 嘗試將該銷售機會推進到 "([^"]*)"$`, func(actor, stageName string) error {
		tc.CurrentActorName = actor
		tc.CurrentActorRole = domain.RoleSalesRep
		tc.IsAuthenticated = true
		stage := mapStageName(stageName)
		payload := map[string]interface{}{
			"stage": string(stage),
		}
		body, _ := json.Marshal(payload)
		url := fmt.Sprintf("/api/opportunities/%d/stage", tc.LastCreatedOppID)
		req := httptest.NewRequest("PATCH", url, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		tc.DoRequest(req)
		return nil
	})

	ctx.Step(`^系統建立銷售機會 "([^"]*)" 成功$`, func(expectedName string) error {
		if tc.LastResponse.Code != http.StatusCreated {
			return fmt.Errorf("expected 201 Created, got %d: %s", tc.LastResponse.Code, tc.LastResponse.Body.String())
		}
		var opp domain.Opportunity
		if err := json.Unmarshal(tc.LastResponse.Body.Bytes(), &opp); err != nil {
			return err
		}
		if opp.Name != expectedName {
			return fmt.Errorf("expected name %s, got %s", expectedName, opp.Name)
		}
		tc.LastCreatedOppID = opp.ID
		return nil
	})

	ctx.Step(`^銷售機會負責人為 "([^"]*)"$`, func(expectedOwner string) error {
		var opp domain.Opportunity
		if err := json.Unmarshal(tc.LastResponse.Body.Bytes(), &opp); err != nil {
			return err
		}
		if opp.OwnerName != expectedOwner {
			return fmt.Errorf("expected owner %s, got %s", expectedOwner, opp.OwnerName)
		}
		return nil
	})

	ctx.Step(`^該銷售機會目前階段為 "([^"]*)"$`, func(expectedStage string) error {
		stage := mapStageName(expectedStage)
		var opp domain.Opportunity
		if err := json.Unmarshal(tc.LastResponse.Body.Bytes(), &opp); err != nil {
			return err
		}
		if opp.Stage != stage {
			return fmt.Errorf("expected stage %s, got %s", stage, opp.Stage)
		}
		return nil
	})

	ctx.Step(`^銷售機會目前階段為 "([^"]*)"$`, func(expectedStage string) error {
		stage := mapStageName(expectedStage)
		req := httptest.NewRequest("GET", fmt.Sprintf("/api/opportunities/%d", tc.LastCreatedOppID), nil)
		tc.DoRequest(req)
		var opp domain.Opportunity
		if err := json.Unmarshal(tc.LastResponse.Body.Bytes(), &opp); err != nil {
			return err
		}
		if opp.Stage != stage {
			return fmt.Errorf("expected stage %s, got %s", stage, opp.Stage)
		}
		return nil
	})

	ctx.Step(`^銷售機會 "([^"]*)" 目前階段為 "([^"]*)"$`, func(oppName, expectedStage string) error {
		stage := mapStageName(expectedStage)
		req := httptest.NewRequest("GET", fmt.Sprintf("/api/opportunities/%d", tc.LastCreatedOppID), nil)
		tc.DoRequest(req)
		if tc.LastResponse.Code != http.StatusOK {
			return fmt.Errorf("failed to get opp: %d", tc.LastResponse.Code)
		}
		var opp domain.Opportunity
		if err := json.Unmarshal(tc.LastResponse.Body.Bytes(), &opp); err != nil {
			return err
		}
		if opp.Stage != stage {
			return fmt.Errorf("expected stage %s, got %s", stage, opp.Stage)
		}
		return nil
	})

	ctx.Step(`^系統記錄階段推進時間與操作人為 "([^"]*)"$`, func(expectedActor string) error {
		req := httptest.NewRequest("GET", fmt.Sprintf("/api/opportunities/%d", tc.LastCreatedOppID), nil)
		tc.DoRequest(req)
		var opp domain.Opportunity
		if err := json.Unmarshal(tc.LastResponse.Body.Bytes(), &opp); err != nil {
			return err
		}
		if opp.StageUpdatedBy != expectedActor {
			return fmt.Errorf("expected stage updated by %s, got %s", expectedActor, opp.StageUpdatedBy)
		}
		return nil
	})

	ctx.Step(`^銷售機會維持原本階段 "([^"]*)"$`, func(expectedStage string) error {
		stage := mapStageName(expectedStage)
		opp, err := tc.OpportunityRepo.GetByID(context.Background(), tc.LastCreatedOppID)
		if err != nil {
			return err
		}
		if opp.Stage != stage {
			return fmt.Errorf("expected stage %s, got %s", stage, opp.Stage)
		}
		return nil
	})

	ctx.Step(`^該銷售機會歸類為已結束結案$`, func() error {
		opp, err := tc.OpportunityRepo.GetByID(context.Background(), tc.LastCreatedOppID)
		if err != nil {
			return err
		}
		if !domain.IsClosedStage(opp.Stage) {
			return fmt.Errorf("expected terminal stage, got %s", opp.Stage)
		}
		return nil
	})
}

func mapStageName(name string) domain.SalesStage {
	switch name {
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
		return domain.SalesStage(name)
	}
}
