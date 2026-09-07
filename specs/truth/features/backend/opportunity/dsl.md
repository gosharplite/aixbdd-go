# Backend Opportunity Module DSL

> **模組 DSL (Backend: Opportunity & Stages)**  
> **Truth Root**: `specs/truth/features/backend/opportunity/dsl.md`  
> **維護者**: `/dsl-refine`

| 句型 (Gherkin Pattern) | 類型 | 參數 / DataTable | 預設值 | StepDef 實作語意 |
|---|---|---|---|---|
| `^CRM 系統定義標準銷售階段如下：$` | Given | DataTable: `[順序, 階段]` | 無 | 驗證系統支援之銷售階段清單（潛在、已聯繫、提案中、已成交、已失敗）。 |
| `^系統中存在公司 "([^"]*)" 與客戶聯絡人 "([^"]*)"$` | Given | `companyName, contactName: string` | 無 | 確保資料庫中存在該公司與聯絡人。 |
| `^"([^"]*)" 負責 "([^"]*)" 的銷售機會 "([^"]*)"$` | Given | `ownerName, companyName, oppName: string` | 無 | 建立由 `{ownerName}` 負責之商機 `{oppName}`，關聯至 `{companyName}`。 |
| `^該銷售機會目前階段為 "([^"]*)"$` | Given | `stageName: string` | 無 | 將商機階段設為指定階段。 |
| `^該銷售機會預估金額為 (\d+) 元$` | Given | `amount: float` | 無 | 將商機預估金額設為指定數值。 |
| `^"([^"]*)" 負責的銷售機會 "([^"]*)" 處於 "([^"]*)" 階段$` | Given | `ownerName, oppName, stageName: string` | 無 | 建立或設定該商機負責人與階段。 |
| `^"([^"]*)" 為一般業務人員，不是該銷售機會負責人$` | Given | `actor: string` | 無 | 設定請求身分為 `{actor}` (`X-User-Role: SalesRep`)，且確認非該商機負責人。 |
| `^"([^"]*)" 不是業務主管$` | Given | `actor: string` | 無 | 確認請求身分不具備 `SalesManager` 角色。 |
| `^"([^"]*)" 建立銷售機會如下：$` | When | `actor: string`, DataTable: `[欄位, 內容]` | 無 | 解析 DataTable 為 `CreateOpportunityRequest`，呼叫 `POST /api/opportunities`。 |
| `^"([^"]*)" 將銷售機會推進到 "([^"]*)"$` | When | `actor, stageName: string` | 無 | 以 `{actor}` 身分呼叫 `PATCH /api/opportunities/{id}/stage`，Body 為 `{"stage": "{stageName}"}`。 |
| `^"([^"]*)" 嘗試將該銷售機會推進到 "([^"]*)"$` | When | `actor, stageName: string` | 無 | 以 `{actor}` 身分呼叫 `PATCH /api/opportunities/{id}/stage`，預期觸發 RBAC 攔截。 |
| `^系統建立銷售機會 "([^"]*)" 成功$` | Then | `oppName: string` | 無 | 斷言 HTTP 回應碼為 201 Created，且名稱為 `{oppName}`。 |
| `^銷售機會負責人為 "([^"]*)"$` | Then | `ownerName: string` | 無 | 斷言回傳商機資料之 `owner_name` 為 `{ownerName}`。 |
| `^該銷售機會目前階段為 "([^"]*)"$` | Then | `stageName: string` | 無 | 斷言回傳商機資料之 `stage` 等於 `{stageName}`。 |
| `^銷售機會 "([^"]*)" 目前階段為 "([^"]*)"$` | Then | `oppName, stageName: string` | 無 | 查詢 `GET /api/opportunities/{id}` 驗證其 `stage` 等於 `{stageName}`。 |
| `^系統記錄階段推進時間與操作人為 "([^"]*)"$` | Then | `actor: string` | 無 | 驗證 `stage_updated_by` 為 `{actor}` 且 `stage_updated_at` 有紀錄。 |
| `^銷售機會維持原本階段 "([^"]*)"$` | Then | `stageName: string` | 無 | 查詢資料庫驗證該商機階段仍為 `{stageName}` 未被變更。 |
| `^該銷售機會歸類為已結束結案$` | Then | 無 | 無 | 驗證該商機階段屬於 ClosedWon 或 ClosedLost 終態。 |
