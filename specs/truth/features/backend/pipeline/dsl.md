# Backend Pipeline Module DSL

> **模組 DSL (Backend: Pipeline & Aggregations)**  
> **Truth Root**: `specs/truth/features/backend/pipeline/dsl.md`  
> **維護者**: `/dsl-refine`

| 句型 (Gherkin Pattern) | 類型 | 參數 / DataTable | 預設值 | StepDef 實作語意 |
|---|---|---|---|---|
| `^CRM 系統中已存在以下銷售機會資料：$` | Given | DataTable: `[機會名稱, 關聯公司, 負責人, 階段, 預估金額]` | 無 | 批量在資料庫建立公司與所屬銷售機會資料。 |
| `^"([^"]*)" 查看全公司銷售管線看板$` | When | `actor: string` | 無 | 以 `{actor}` (SalesManager) 身分呼叫 `GET /api/pipeline`。 |
| `^"([^"]*)" 查看銷售管線$` | When | `actor: string` | 無 | 以 `{actor}` (SalesRep) 身分呼叫 `GET /api/pipeline`。 |
| `^系統在銷售管線中呈現包含 "([^"]*)" 與 "([^"]*)" 負責的所有銷售機會$` | Then | `owner1, owner2: string` | 無 | 驗證管線資料包含兩位業務負責之項目清單。 |
| `^各階段統計指標如下：$` | Then | DataTable: `[階段, 機會筆數, 階段加總金額]` | 無 | 比對 `GET /api/pipeline` 回傳之 `stages` 陣列中每個 stage 的 `count` 與 `total_amount`。 |
| `^進行中階段（潛在、已聯繫、提案中）總加總預估金額為 (\d+) 元$` | Then | `expectedSum: float` | 無 | 斷言回傳之 `active_total_amount` 等於 `{expectedSum}`。 |
| `^系統僅顯示負責人為 "([^"]*)" 的銷售機會如下：$` | Then | `owner: string`, DataTable: `[機會名稱, 關聯公司, 階段, 預估金額]` | 無 | 驗證回傳的管線商機皆屬於 `{owner}`，且項目與 DataTable 一致。 |
| `^不應顯示負責人為 "([^"]*)" 的銷售機會$` | Then | `unauthorizedOwner: string` | 無 | 斷言回傳項目中無任何 `owner_name` 為 `{unauthorizedOwner}` 的商機。 |
