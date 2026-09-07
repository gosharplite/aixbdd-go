# Backend Customer Module DSL

> **模組 DSL (Backend: Customer & Company)**  
> **Truth Root**: `specs/truth/features/backend/customer/dsl.md`  
> **維護者**: `/dsl-refine`

| 句型 (Gherkin Pattern) | 類型 | 參數 / DataTable | 預設值 | StepDef 實作語意 |
|---|---|---|---|---|
| `^系統中已存在公司 "([^"]*)"$` | Given | `companyName: string` | 無 | 透過 Repository 或 API 確保建立公司名稱為 `{companyName}`，記錄其 `company_id`。 |
| `^公司 "([^"]*)" 底下已存在客戶聯絡人 "([^"]*)"，職稱為 "([^"]*)"，電話為 "([^"]*)"$` | Given | `companyName, contactName, title, phone: string` | 無 | 在指定公司下透過 Repository 建立聯絡人，記錄其 `contact_id`。預設 Email 為 `{contactName}@test.com`。 |
| `^"([^"]*)" 建立公司資料如下：$` | When | `actor: string`, DataTable: `[欄位, 內容]` | 無 | 將 DataTable 解析為 `CreateCompanyRequest`，呼叫 `POST /api/companies`。 |
| `^"([^"]*)" 為公司 "([^"]*)" 新增客戶聯絡人如下：$` | When | `actor, companyName: string`, DataTable: `[欄位, 內容]` | 無 | 將 DataTable 解析為 `CreateContactRequest`，呼叫 `POST /api/contacts`。 |
| `^"([^"]*)" 將客戶 "([^"]*)" 的職稱更新為 "([^"]*)"，電話更新為 "([^"]*)"$` | When | `actor, contactName, title, phone: string` | 無 | 呼叫 `PUT /api/contacts/{id}` 更新指定聯絡人之職稱與電話。 |
| `^該訪客嘗試提交建立公司 "([^"]*)" 的請求$` | When | `companyName: string` | 無 | 未帶認證 Header 呼叫 `POST /api/companies`，Body 含 `{"name": "{companyName}"}`。 |
| `^系統建立公司 "([^"]*)" 成功$` | Then | `companyName: string` | 無 | 斷言 HTTP 回應碼為 201 Created，且回傳之公司名稱為 `{companyName}`，ID 不為空。 |
| `^系統建立客戶聯絡人 "([^"]*)" 成功$` | Then | `contactName: string` | 無 | 斷言 HTTP 回應碼為 201 Created，且回傳之聯絡人姓名為 `{contactName}`。 |
| `^"([^"]*)" 的公司關聯為 "([^"]*)"$` | Then | `contactName, companyName: string` | 無 | 呼叫 `GET /api/contacts/{id}` 驗證其所屬 `company_name` 等於 `{companyName}`。 |
| `^客戶聯絡人 "([^"]*)" 的最新資訊如下：$` | Then | `contactName: string`, DataTable: `[項目, 值]` | 無 | 呼叫 `GET /api/contacts/{id}` 驗證各屬性（職稱、聯絡電話）皆與 DataTable 相符。 |
| `^系統中不應存在公司 "([^"]*)"$` | Then | `companyName: string` | 無 | 查詢資料庫確認公司 `{companyName}` 筆數為 0。 |
