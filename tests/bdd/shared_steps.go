package bdd

import (
	"fmt"
	"net/http"

	"crm/internal/domain"
	"github.com/cucumber/godog"
)

func RegisterSharedSteps(ctx *godog.ScenarioContext, tc *TestContext) {
	ctx.Step(`^業務人員 "([^"]*)" 已登入 CRM 系統且具備業務權限$`, func(name string) error {
		tc.CurrentActorName = name
		tc.CurrentActorRole = domain.RoleSalesRep
		tc.IsAuthenticated = true
		return nil
	})

	ctx.Step(`^"([^"]*)" 為業務主管$`, func(name string) error {
		tc.CurrentActorName = name
		tc.CurrentActorRole = domain.RoleSalesManager
		tc.IsAuthenticated = true
		return nil
	})

	ctx.Step(`^"([^"]*)" 是一般業務人員$`, func(name string) error {
		tc.CurrentActorName = name
		tc.CurrentActorRole = domain.RoleSalesRep
		tc.IsAuthenticated = true
		return nil
	})

	ctx.Step(`^訪客尚未通過 CRM 身份驗證$`, func() error {
		tc.CurrentActorName = ""
		tc.IsAuthenticated = false
		tc.Headers = make(map[string]string)
		return nil
	})

	ctx.Step(`^系統拒絕該操作，回應狀態碼為 (\d+)$`, func(expectedStatus int) error {
		if tc.LastResponse.Code != expectedStatus {
			return fmt.Errorf("expected status %d, got %d. Body: %s", expectedStatus, tc.LastResponse.Code, tc.LastResponse.Body.String())
		}
		return nil
	})

	ctx.Step(`^系統拒絕這次操作$`, func() error {
		if tc.LastResponse.Code != http.StatusForbidden {
			return fmt.Errorf("expected status 403 Forbidden, got %d. Body: %s", tc.LastResponse.Code, tc.LastResponse.Body.String())
		}
		return nil
	})
}
