# Tasks: 001-crm-core CRM 客戶關係管理系統核心

**Plan Package**: `specs/plans/001-crm-core`  
**Truth Root**: `specs/truth`  
**Core Inputs**: `spec.md`, `research.md`, `plan.md`, `truth-delta.md`, `specs/truth/techstack.md`, `specs/truth/contracts/openapi.yaml`, `specs/truth/data/schema.dbml`, `specs/truth/features/**`

---

## Phase 1: Setup

> 本輪引入新技術（Go 1.22+, Chi Router, Godog, SQLite modernc.org/sqlite）。只做環境配置與冒煙測試，不寫 DSL 語意與業務邏輯。

- [X] T001 初始化 Go 模組、安裝核心依賴與執行 Smoke Test
  - 初始化 `go.mod` (module `crm`)
  - 安裝核心依賴：`github.com/go-chi/chi/v5`, `github.com/cucumber/godog`, `modernc.org/sqlite`
  - 建立基礎目錄結構 (`cmd/server`, `internal/domain`, `internal/service`, `internal/adapter/http`, `internal/adapter/repository`, `tests/bdd`)
  - 驗證 `go test ./...` 正常退出（Smoke Test 通過）

---

## Phase 2: Foundational

> 建立領域實體、資料庫模型與儲存骨架、HTTP 基礎設施與測試運行基底。每則任務嚴格遵循「只做／不做」。

- [X] T002 定義核心領域模型實體與 Repository 介面
  - **只做**: 在 `internal/domain` 建立 `Company`, `Contact`, `Interaction`, `Opportunity`, `SalesStage`, `User`, `UserRole` 結構體與領域約束，宣告各 Repository 介面。
  - **不做**: 不寫具體 SQL 實作、不寫 HTTP 邏輯與不寫商業用例。
- [X] T003 建立 SQLite 資料庫遷移與 Repository 實作
  - **只做**: 在 `internal/adapter/repository` 建立 `schema.sql` 與 SQLite 連線池管理，實作支援 `:memory:` 與檔案資料庫的 Repository，包含外鍵約束開啟與 WAL 模式。
  - **不做**: 不處理 HTTP 請求，不實作上層業務應用服務。
- [X] T004 建立 HTTP 伺服器、Chi 路由器、認證/RBAC 中介軟體與前端靜態託管
  - **只做**: 在 `internal/adapter/http` 建立 HTTP 路由骨架、JSON 序列化與錯誤回應 helper、身分提取中介軟體 (`X-User-Name`, `X-User-Role`)，以及靜態資源託管 (`specs/plans/001-crm-core/ui/**`)。
  - **不做**: 不實作各業務功能之具體 Handler 商業邏輯。
- [X] T005 建立 Godog BDD 測試套件執行器與測試情境生命週期
  - **只做**: 在 `tests/bdd/godog_test.go` 建立 TestMain 與 Scenario 勾點 (BeforeScenario/AfterScenario)，確保每個 Scenario 執行前自動初始化全新的 SQLite `:memory:` 資料庫與測試 HTTP 伺服器。
  - **不做**: 不實作具體的 Step Definitions 匹配邏輯。

---

## Phase 3: Test Alignment & Implementation

> 集中對齊與實作所有 BDD Step Definitions。本階段不實作任何產品碼，只建立測試層。

### DSL 參照
- 根共用 DSL: `specs/truth/features/backend/dsl.md`
- 客戶與公司模組 DSL: `specs/truth/features/backend/customer/dsl.md`
- 聯絡歷程模組 DSL: `specs/truth/features/backend/interaction/dsl.md`
- 銷售機會模組 DSL: `specs/truth/features/backend/opportunity/dsl.md`
- 銷售管線模組 DSL: `specs/truth/features/backend/pipeline/dsl.md`

### Markers
- `[BDD-RED]`: 建立 Step Definitions，執行測試時必須因尚未有產品實作而呈現乾淨的失敗（RED）訊號。

### Shared Must Read
- `specs/truth/contracts/openapi.yaml`
- `specs/truth/data/schema.dbml`
- 各模組對應之 `specs/truth/features/backend/**/dsl.md`

### Boundary
- 嚴格限制在 `tests/bdd/` 下編寫測試步驟程式碼，不更動 `internal/` 產品程式碼。

### Parallel Hint
- T006 ~ T010 為各模組獨立的 Step Definitions，可平行實作。

### Tasks
- [X] T006 `[BDD-RED]` `[P]` 實作後端根共用 DSL Step Definitions
- [X] T007 `[BDD-RED]` `[P]` 實作客戶與公司模組 DSL Step Definitions
- [X] T008 `[BDD-RED]` `[P]` 實作聯絡歷程模組 DSL Step Definitions
- [X] T009 `[BDD-RED]` `[P]` 實作銷售機會模組 DSL Step Definitions
- [X] T010 `[BDD-RED]` `[P]` 實作銷售管線模組 DSL Step Definitions
- [X] T011 審查全體 Step Definitions 並確認測試皆產出有效失敗（RED）訊號
  - 執行 `go test -v ./tests/bdd/...`，確認所有 4 份 backend feature 均被正確識別且因未實作後端服務而產生明確的失敗 (RED)。

---

## Phase 4: Feature Phase - 客戶與公司資料管理 (Customer & Company)

### Shared Must Read
- `specs/truth/features/backend/customer/customer.feature`
- `specs/truth/features/backend/customer/dsl.md`
- `specs/truth/contracts/openapi.yaml`
- `specs/truth/data/schema.dbml`

### Boundary
- 僅實作客戶與公司相關之 Repository 查詢、業務用例與 HTTP Handlers (`/api/companies`, `/api/contacts`)。不更動聯絡歷程或銷售機會。

### Test Scope
- `specs/truth/features/backend/customer/customer.feature`

### Tasks
- [X] T012 `[BDD-GREEN]` 實作公司與客戶業務服務與 REST API Handlers
- [X] T013 `[BDD-REFACTOR]` 重構客戶與公司模組程式碼

---

## Phase 5: Feature Phase - 記錄聯絡歷程與下一步行動 (Interaction & Next Action)

### Shared Must Read
- `specs/truth/features/backend/interaction/interaction.feature`
- `specs/truth/features/backend/interaction/dsl.md`
- `specs/truth/contracts/openapi.yaml`
- `specs/truth/data/schema.dbml`

### Boundary
- 僅實作客戶聯絡歷程與下一步行動之業務邏輯、Append-only 資料寫入與倒序查詢。不更動銷售機會。

### Test Scope
- `specs/truth/features/backend/interaction/interaction.feature`

### Tasks
- [X] T014 `[BDD-GREEN]` 實作聯絡歷程記錄與倒序查詢業務服務與 API
- [X] T015 `[BDD-REFACTOR]` 重構聯絡歷程模組

---

## Phase 6: Feature Phase - 銷售機會建立與階段推進 (Opportunity & Stage Progression)

### Shared Must Read
- `specs/truth/features/backend/opportunity/opportunity.feature`
- `specs/truth/features/backend/opportunity/dsl.md`
- `specs/truth/contracts/openapi.yaml`
- `specs/truth/data/schema.dbml`

### Boundary
- 僅實作銷售機會建立、階段推進與嚴格 RBAC 權限校驗。不包含管線統計指標。

### Test Scope
- `specs/truth/features/backend/opportunity/opportunity.feature`

### Tasks
- [X] T016 `[BDD-GREEN]` 實作銷售機會建立、階段推進與 RBAC 權限防護
- [X] T017 `[BDD-REFACTOR]` 重構銷售機會模組

---

## Phase 7: Feature Phase - 銷售管線綜覽與業績預測 (Pipeline & Forecasting)

### Shared Must Read
- `specs/truth/features/backend/pipeline/pipeline.feature`
- `specs/truth/features/backend/pipeline/dsl.md`
- `specs/truth/contracts/openapi.yaml`
- `specs/truth/data/schema.dbml`

### Boundary
- 僅實作銷售管線多維度統計聚合與角色資料隔離。

### Test Scope
- `specs/truth/features/backend/pipeline/pipeline.feature`

### Tasks
- [X] T018 `[BDD-GREEN]` 實作銷售管線統計聚合與角色資料隔離
- [X] T019 `[BDD-REFACTOR]` 重構管線聚合計算邏輯

---

## Phase 8: Full Integration & Static UI Verification

### Shared Must Read
- `specs/plans/001-crm-core/ui/**`
- `specs/truth/contracts/openapi.yaml`

### Boundary
- 確保靜態 Web 原型在後端服務啟動時可無縫操作與呼叫 API。

### Tasks
- [X] T020 前端原型與後端 RESTful API 整合對接
- [X] T021 全系統迴歸測試與品質驗收
