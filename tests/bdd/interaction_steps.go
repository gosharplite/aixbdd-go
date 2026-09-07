package bdd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"crm/internal/domain"
	"github.com/cucumber/godog"
)

func RegisterInteractionSteps(ctx *godog.ScenarioContext, tc *TestContext) {
	ctx.Step(`^系統中已存在客戶 "([^"]*)"，任職於 "([^"]*)"$`, func(contactName, compName string) error {
		comp, err := tc.CompanyRepo.GetByName(context.Background(), compName)
		if err != nil {
			comp = &domain.Company{Name: compName, Industry: "軟體資訊"}
			if err := tc.CompanyRepo.Create(context.Background(), comp); err != nil {
				return err
			}
		}
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

	ctx.Step(`^客戶 "([^"]*)" 已有以下歷史聯絡紀錄：$`, func(contactName string, table *godog.Table) error {
		contacts, err := tc.ContactRepo.List(context.Background(), nil)
		if err != nil {
			return err
		}
		var target *domain.Contact
		for _, c := range contacts {
			if c.Name == contactName {
				target = &c
				break
			}
		}
		if target == nil {
			return fmt.Errorf("contact %s not found", contactName)
		}

		user, _ := tc.UserRepo.GetByName(context.Background(), "Alice")
		userID := int64(2)
		if user != nil {
			userID = user.ID
		}

		for i, row := range table.Rows {
			if i == 0 && (row.Cells[0].Value == "時間" || row.Cells[0].Value == "順序") {
				continue
			}
			if len(row.Cells) < 3 {
				continue
			}
			timeStr := row.Cells[0].Value
			typeStr := row.Cells[1].Value
			summary := row.Cells[2].Value

			t, err := time.Parse("2006-01-02 15:04", timeStr)
			if err != nil {
				t = time.Now()
			}

			inter := &domain.Interaction{
				ContactID:       target.ID,
				UserID:          userID,
				Type:            domain.InteractionType(typeStr),
				InteractionTime: t,
				Summary:         summary,
			}
			if err := tc.InteractionRepo.Create(context.Background(), inter); err != nil {
				return err
			}
		}
		return nil
	})

	ctx.Step(`^"([^"]*)" 記錄與客戶 "([^"]*)" 的聯絡內容如下：$`, func(actor, contactName string, table *godog.Table) error {
		tc.CurrentActorName = actor
		contacts, err := tc.ContactRepo.List(context.Background(), nil)
		if err != nil {
			return err
		}
		var target *domain.Contact
		for _, c := range contacts {
			if c.Name == contactName {
				target = &c
				break
			}
		}
		if target == nil {
			return fmt.Errorf("contact %s not found", contactName)
		}

		payload := make(map[string]interface{})
		for _, row := range table.Rows {
			if len(row.Cells) < 2 {
				continue
			}
			key := row.Cells[0].Value
			val := row.Cells[1].Value
			switch key {
			case "聯絡方式":
				payload["type"] = val
			case "聯絡時間":
				payload["interaction_time"] = val
			case "內容摘要":
				payload["summary"] = val
			case "下一步行動":
				payload["next_action"] = val
			case "預計執行日期":
				payload["next_action_date"] = val
			}
		}

		body, _ := json.Marshal(payload)
		url := fmt.Sprintf("/api/contacts/%d/interactions", target.ID)
		req := httptest.NewRequest("POST", url, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		tc.DoRequest(req)
		return nil
	})

	ctx.Step(`^業務人員 "([^"]*)" 查詢客戶 "([^"]*)" 的歷程資訊$`, func(actor, contactName string) error {
		tc.CurrentActorName = actor
		contacts, err := tc.ContactRepo.List(context.Background(), nil)
		if err != nil {
			return err
		}
		var target *domain.Contact
		for _, c := range contacts {
			if c.Name == contactName {
				target = &c
				break
			}
		}
		if target == nil {
			return fmt.Errorf("contact %s not found", contactName)
		}

		url := fmt.Sprintf("/api/contacts/%d/interactions", target.ID)
		req := httptest.NewRequest("GET", url, nil)
		tc.DoRequest(req)
		return nil
	})

	ctx.Step(`^系統成功儲存該筆聯絡紀錄$`, func() error {
		if tc.LastResponse.Code != http.StatusCreated {
			return fmt.Errorf("expected 201 Created, got %d: %s", tc.LastResponse.Code, tc.LastResponse.Body.String())
		}
		return nil
	})

	ctx.Step(`^客戶 "([^"]*)" 的最新下一步行動摘要如下：$`, func(contactName string, table *godog.Table) error {
		contacts, err := tc.ContactRepo.List(context.Background(), nil)
		if err != nil {
			return err
		}
		var target *domain.Contact
		for _, c := range contacts {
			if c.Name == contactName {
				target = &c
				break
			}
		}
		if target == nil {
			return fmt.Errorf("contact %s not found", contactName)
		}

		req := httptest.NewRequest("GET", fmt.Sprintf("/api/contacts/%d", target.ID), nil)
		tc.DoRequest(req)
		if tc.LastResponse.Code != http.StatusOK {
			return fmt.Errorf("failed to get contact: %d", tc.LastResponse.Code)
		}

		var contact domain.Contact
		if err := json.Unmarshal(tc.LastResponse.Body.Bytes(), &contact); err != nil {
			return err
		}

		for _, row := range table.Rows {
			if len(row.Cells) < 2 {
				continue
			}
			item := row.Cells[0].Value
			val := row.Cells[1].Value
			switch item {
			case "待辦行動":
				if contact.LatestNextAction != val {
					return fmt.Errorf("expected next action %q, got %q", val, contact.LatestNextAction)
				}
			case "預計執行日期":
				if contact.LatestNextActionDate != val {
					return fmt.Errorf("expected next action date %q, got %q", val, contact.LatestNextActionDate)
				}
			}
		}
		return nil
	})

	ctx.Step(`^系統依時間由新到舊列出聯絡紀錄：$`, func(table *godog.Table) error {
		if tc.LastResponse.Code != http.StatusOK {
			return fmt.Errorf("expected 200 OK, got %d", tc.LastResponse.Code)
		}
		var list []domain.Interaction
		if err := json.Unmarshal(tc.LastResponse.Body.Bytes(), &list); err != nil {
			return err
		}

		tableRows := table.Rows
		if len(tableRows) > 0 && tableRows[0].Cells[0].Value == "順序" {
			tableRows = tableRows[1:]
		}

		if len(list) != len(tableRows) {
			return fmt.Errorf("expected %d records, got %d", len(tableRows), len(list))
		}

		for idx, row := range tableRows {
			expectedType := row.Cells[2].Value
			expectedSummary := row.Cells[3].Value
			actual := list[idx]
			if string(actual.Type) != expectedType {
				return fmt.Errorf("row %d expected type %s, got %s", idx, expectedType, actual.Type)
			}
			if actual.Summary != expectedSummary {
				return fmt.Errorf("row %d expected summary %s, got %s", idx, expectedSummary, actual.Summary)
			}
		}
		return nil
	})
}
