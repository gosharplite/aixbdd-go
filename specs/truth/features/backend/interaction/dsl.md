# Backend Interaction Module DSL

> **模組 DSL (Backend: Interaction & History)**  
> **Truth Root**: `specs/truth/features/backend/interaction/dsl.md`  
> **維護者**: `/dsl-refine`

| 句型 (Gherkin Pattern) | 類型 | 參數 / DataTable | 預設值 | StepDef 實作語意 |
|---|---|---|---|---|
| `^系統中已存在客戶 "([^"]*)"，任職於 "([^"]*)"$` | Given | `contactName, companyName: string` | 無 | 確保建立公司與所屬聯絡人，記錄其 `contact_id`。 |
| `^客戶 "([^"]*)" 已有以下歷史聯絡紀錄：$` | Given | `contactName: string`, DataTable: `[時間, 方式, 內容摘要]` | 無 | 依 DataTable 各列依序呼叫 `POST /api/contacts/{id}/interactions` 寫入歷史資料。 |
| `^"([^"]*)" 記錄與客戶 "([^"]*)" 的聯絡內容如下：$` | When | `actor, contactName: string`, DataTable: `[欄位, 內容]` | 無 | 將 DataTable 解析為 `CreateInteractionRequest`，呼叫 `POST /api/contacts/{id}/interactions`。 |
| `^業務人員 "([^"]*)" 查詢客戶 "([^"]*)" 的歷程資訊$` | When | `actor, contactName: string` | 無 | 呼叫 `GET /api/contacts/{id}/interactions`。 |
| `^系統成功儲存該筆聯絡紀錄$` | Then | 無 | 無 | 斷言最後一次 HTTP 回應碼為 201 Created，且資料庫中有此筆記錄。 |
| `^客戶 "([^"]*)" 的最新下一步行動摘要如下：$` | Then | `contactName: string`, DataTable: `[項目, 值]` | 無 | 呼叫 `GET /api/contacts/{id}` 驗證 `latest_next_action` 與 `latest_next_action_date` 符合 DataTable。 |
| `^系統依時間由新到舊列出聯絡紀錄：$` | Then | DataTable: `[順序, 聯絡時間, 方式, 內容摘要]` | 無 | 檢驗 `GET /api/contacts/{id}/interactions` 回傳之陣列順序，第一筆為最新，第二筆次之，且各欄位匹配。 |
