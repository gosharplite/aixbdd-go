# Backend Root Shared DSL

> **介面根共用 DSL (Backend)**  
> **Truth Root**: `specs/truth/features/backend/dsl.md`  
> **維護者**: `/dsl-refine`

| 句型 (Gherkin Pattern) | 類型 | 參數 / DataTable | 預設值 | StepDef 實作語意 |
|---|---|---|---|---|
| `^業務人員 "([^"]*)" 已登入 CRM 系統且具備業務權限$` | Given | `name: string` | 無 | 設定當前 HTTP 請求 Actor 為該業務人員，Header 附加 `X-User-Name: {name}`, `X-User-Role: SalesRep`。 |
| `^"([^"]*)" 為業務主管$` | Given | `name: string` | 無 | 設定當前 HTTP 請求 Actor 為該主管，Header 附加 `X-User-Name: {name}`, `X-User-Role: SalesManager`。 |
| `^"([^"]*)" 是一般業務人員$` | Given | `name: string` | 無 | 設定當前 HTTP 請求 Actor 為一般業務，Header 附加 `X-User-Name: {name}`, `X-User-Role: SalesRep`。 |
| `^訪客尚未通過 CRM 身份驗證$` | Given | 無 | 無 | 清除當前 HTTP 請求所有認證與身分 Headers（模擬未帶 Token/憑證的訪客）。 |
| `^系統拒絕該操作，回應狀態碼為 (\d+)$` | Then | `statusCode: int` | 401 | 斷言最後一次 HTTP 回應狀態碼等於 `{statusCode}`。 |
| `^系統拒絕這次操作$` | Then | 無 | 403 | 斷言最後一次 HTTP 回應狀態碼為 403 Forbidden。 |
