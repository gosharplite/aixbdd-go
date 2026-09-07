# 系統分析計畫：001-crm-core

**Plan Package**: `specs/plans/001-crm-core`  
**建立日期**: 2026-03-30  
**負責角色**: RD (`/system-analysis`)  
**狀態**: 已確認 (Approved)

---

## 1. 系統介面與邊界盤點

依據 `spec.md`、`research.md` 及 PM 驗收 Feature，本輪迭代涉及三大系統介面：

| 介面名稱 | 介面類型 | 負責 Planner | 目標 Truth Artifact | 範圍與重點說明 |
|---|---|---|---|---|
| **資料模型介面 (Data Schema)** | 資料庫模型 (DBML) | `/data-plan` | `specs/truth/data/schema.dbml` | 公司、客戶聯絡人、聯絡歷程、銷售機會、使用者與權限之實體、欄位型別、外鍵關聯與約束。 |
| **後端契約介面 (Backend API)** | RESTful API (OpenAPI 3.0) | `/api-plan` | `specs/truth/contracts/openapi.yaml` | 客戶/公司 CRUD、聯絡歷程記錄與倒序查詢、銷售機會建立與階段推進 (RBAC 檢查)、銷售管線統計指標。 |
| **前端使用者介面 (Web UI)** | 靜態 Web 原型審查 | (PM 已交付，RD 複核) | `specs/plans/001-crm-core/ui/**` | 複核 4 份 HTML 原型之表單欄位、狀態機過渡、多角色視圖（Carol vs Alice）與後端 API 契約的對接可行性。 |

---

## 2. 分析 Wave 規劃與依賴排序

為確保資料完整性與契約一致性，採分波依序（Wave-based）委派推進：

### Wave 1: 資料模型設計 (`/data-plan`)
- **委派目標**: 產出完整系統唯一資料模型真相 `specs/truth/data/schema.dbml`。
- **分析重點**:
  1. `companies` 表：包含名稱、產業、統編、地址，設定名稱非空與刪除保護（若有客戶或商機則限制外鍵刪除 `RESTRICT`）。
  2. `contacts` 表：關聯至 `companies.id`，包含姓名、Email（驗證格式）、電話、職稱。
  3. `interactions` 表：關聯至 `contacts.id`，記錄方式 (Email/電話/會議)、時間戳、內容摘要、下一步行動內容與預計日期。歷史紀錄禁止任意修改刪除。
  4. `opportunities` 表：關聯至 `companies.id` 與 `contacts.id`，包含機會名稱、負責人 (`owner_user_id`)、階段（Enum: Prospecting, Contacted, Proposal, ClosedWon, ClosedLost）、預估金額（Decimal）、預計成交日。
  5. `users` 表：帳號名稱、角色 (`SalesManager`, `SalesRep`)。

### Wave 2: 後端 API 契約設計 (`/api-plan`)
- **依賴**: Wave 1 資料模型。
- **委派目標**: 產出完整系統唯一 API 契約真相 `specs/truth/contracts/openapi.yaml`。
- **分析重點**:
  1. 認證與身分：`GET /api/me`、測試角色切換。
  2. 客戶與公司 API：`GET/POST /api/companies`, `GET/PUT /api/companies/{id}`, `GET/POST /api/contacts`, `GET/PUT /api/contacts/{id}`。
  3. 聯絡歷程 API：`POST /api/contacts/{id}/interactions`, `GET /api/contacts/{id}/interactions` (依 `interaction_time` 倒序排序)。
  4. 銷售機會 API：`GET/POST /api/opportunities`, `GET /api/opportunities/{id}`, `PATCH /api/opportunities/{id}/stage` (執行 RBAC 嚴格驗證：僅負責人或主管可變更，非負責人回傳 403 Forbidden)。
  5. 銷售管線 API：`GET /api/pipeline` (主管回傳全員分階段指標與清單；一般業務僅回傳自己負責之機會)。
  6. 錯誤模型：標準結構化錯誤（RFC 7807 衍生），統一錯誤碼（`UNAUTHORIZED`, `FORBIDDEN`, `VALIDATION_FAILED`, `NOT_FOUND`）。

### Wave 3: 前端原型對接複核 (Frontend Alignment Review)
- **依賴**: Wave 2 API 契約。
- **分析重點**:
  1. 複核 `ui/index.html` (管線看板) 之階段欄位、金額合計、卡片拖曳與角色切換是否與 `GET /api/pipeline` 一致。
  2. 複核 `ui/customers.html` 與 `ui/customer-detail.html` 之表單欄位與時間軸渲染是否能完整映射 API 回傳資料。
  3. 複核 `ui/opportunity-detail.html` 之 5 階段步進器及非負責人禁止變更之 UI 提示是否與 API 403 錯誤行為對齊。

---

## 3. 委派清單與交付產物

- **Wave 1**: 委派 `/data-plan` -> 產出 `specs/truth/data/schema.dbml` 並更新 `truth-delta.md`。
- **Wave 2**: 委派 `/api-plan` -> 產出 `specs/truth/contracts/openapi.yaml` 並更新 `truth-delta.md`。
- **後續步驟**: 完成後進入 `/dsl-refine`，將 PM 驗收 Gherkin 拆解為前後端可執行測試規格與 DSL。
