package bdd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"crm/internal/domain"
	"github.com/cucumber/godog"
)

func RegisterCustomerSteps(ctx *godog.ScenarioContext, tc *TestContext) {
	ctx.Step(`^系統中已存在公司 "([^"]*)"$`, func(companyName string) error {
		c, err := tc.CompanyRepo.GetByName(context.Background(), companyName)
		if err == nil && c != nil {
			tc.LastCreatedCompanyID = c.ID
			return nil
		}
		newComp := &domain.Company{
			Name:     companyName,
			Industry: "軟體資訊",
		}
		if err := tc.CompanyRepo.Create(context.Background(), newComp); err != nil {
			return err
		}
		tc.LastCreatedCompanyID = newComp.ID
		return nil
	})

	ctx.Step(`^公司 "([^"]*)" 底下已存在客戶聯絡人 "([^"]*)"，職稱為 "([^"]*)"，電話為 "([^"]*)"$`, func(compName, contactName, title, phone string) error {
		comp, err := tc.CompanyRepo.GetByName(context.Background(), compName)
		if err != nil {
			return fmt.Errorf("company not found: %s", compName)
		}
		contact := &domain.Contact{
			CompanyID: comp.ID,
			Name:      contactName,
			Email:     strings.ToLower(contactName) + "@test.com",
			Phone:     phone,
			Title:     title,
		}
		if err := tc.ContactRepo.Create(context.Background(), contact); err != nil {
			return err
		}
		tc.LastCreatedContactID = contact.ID
		return nil
	})

	ctx.Step(`^"([^"]*)" 建立公司資料如下：$`, func(actor string, table *godog.Table) error {
		tc.CurrentActorName = actor
		payload := make(map[string]interface{})
		for _, row := range table.Rows {
			if len(row.Cells) < 2 {
				continue
			}
			key := row.Cells[0].Value
			val := row.Cells[1].Value
			switch key {
			case "公司名稱":
				payload["name"] = val
			case "產業類別":
				payload["industry"] = val
			case "聯絡地址":
				payload["address"] = val
			case "統一編號":
				payload["tax_id"] = val
			}
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/companies", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		tc.DoRequest(req)
		return nil
	})

	ctx.Step(`^"([^"]*)" 為公司 "([^"]*)" 新增客戶聯絡人如下：$`, func(actor, compName string, table *godog.Table) error {
		tc.CurrentActorName = actor
		comp, err := tc.CompanyRepo.GetByName(context.Background(), compName)
		if err != nil {
			return fmt.Errorf("company not found: %s", compName)
		}

		payload := make(map[string]interface{})
		payload["company_id"] = comp.ID
		for _, row := range table.Rows {
			if len(row.Cells) < 2 {
				continue
			}
			key := row.Cells[0].Value
			val := row.Cells[1].Value
			switch key {
			case "姓名":
				payload["name"] = val
			case "電子郵件":
				payload["email"] = val
			case "聯絡電話":
				payload["phone"] = val
			case "職稱":
				payload["title"] = val
			}
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/contacts", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		tc.DoRequest(req)
		return nil
	})

	ctx.Step(`^"([^"]*)" 將客戶 "([^"]*)" 的職稱更新為 "([^"]*)"，電話更新為 "([^"]*)"$`, func(actor, contactName, title, phone string) error {
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

		payload := map[string]interface{}{
			"title": title,
			"phone": phone,
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("PUT", fmt.Sprintf("/api/contacts/%d", target.ID), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		tc.DoRequest(req)
		return nil
	})

	ctx.Step(`^該訪客嘗試提交建立公司 "([^"]*)" 的請求$`, func(compName string) error {
		tc.IsAuthenticated = false
		tc.CurrentActorName = ""
		payload := map[string]interface{}{
			"name": compName,
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/companies", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		tc.DoRequest(req)
		return nil
	})

	ctx.Step(`^系統建立公司 "([^"]*)" 成功$`, func(expectedName string) error {
		if tc.LastResponse.Code != http.StatusCreated {
			return fmt.Errorf("expected 201 Created, got %d: %s", tc.LastResponse.Code, tc.LastResponse.Body.String())
		}
		var comp domain.Company
		if err := json.Unmarshal(tc.LastResponse.Body.Bytes(), &comp); err != nil {
			return err
		}
		if comp.Name != expectedName {
			return fmt.Errorf("expected company name %s, got %s", expectedName, comp.Name)
		}
		tc.LastCreatedCompanyID = comp.ID
		return nil
	})

	ctx.Step(`^系統建立客戶聯絡人 "([^"]*)" 成功$`, func(expectedName string) error {
		if tc.LastResponse.Code != http.StatusCreated {
			return fmt.Errorf("expected 201 Created, got %d: %s", tc.LastResponse.Code, tc.LastResponse.Body.String())
		}
		var contact domain.Contact
		if err := json.Unmarshal(tc.LastResponse.Body.Bytes(), &contact); err != nil {
			return err
		}
		if contact.Name != expectedName {
			return fmt.Errorf("expected contact name %s, got %s", expectedName, contact.Name)
		}
		tc.LastCreatedContactID = contact.ID
		return nil
	})

	ctx.Step(`^"([^"]*)" 的公司關聯為 "([^"]*)"$`, func(contactName, compName string) error {
		req := httptest.NewRequest("GET", fmt.Sprintf("/api/contacts/%d", tc.LastCreatedContactID), nil)
		tc.DoRequest(req)
		if tc.LastResponse.Code != http.StatusOK {
			return fmt.Errorf("failed to get contact: %d", tc.LastResponse.Code)
		}
		var contact domain.Contact
		if err := json.Unmarshal(tc.LastResponse.Body.Bytes(), &contact); err != nil {
			return err
		}
		if contact.CompanyName != compName {
			return fmt.Errorf("expected company %s, got %s", compName, contact.CompanyName)
		}
		return nil
	})

	ctx.Step(`^客戶聯絡人 "([^"]*)" 的最新資訊如下：$`, func(contactName string, table *godog.Table) error {
		contacts, err := tc.ContactRepo.List(context.Background(), nil)
		if err != nil {
			return err
		}
		var targetID int64
		for _, c := range contacts {
			if c.Name == contactName {
				targetID = c.ID
				break
			}
		}
		req := httptest.NewRequest("GET", fmt.Sprintf("/api/contacts/%d", targetID), nil)
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
			case "職稱":
				if contact.Title != val {
					return fmt.Errorf("expected title %s, got %s", val, contact.Title)
				}
			case "聯絡電話":
				if contact.Phone != val {
					return fmt.Errorf("expected phone %s, got %s", val, contact.Phone)
				}
			}
		}
		return nil
	})

	ctx.Step(`^系統中不應存在公司 "([^"]*)"$`, func(compName string) error {
		_, err := tc.CompanyRepo.GetByName(context.Background(), compName)
		if err == nil {
			return fmt.Errorf("company %s should not exist", compName)
		}
		return nil
	})
}
