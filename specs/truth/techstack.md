# System Tech Stack Truth

> **系統唯一技術堆疊規格書 (System Tech Stack Truth)**  
> **維護者**: `/technical-research`  
> **最後更新**: 2026-03-30 (Plan Package: `001-crm-core`)

---

## 1. 核心語言與執行環境

| 維度 | 技術選型 | 版本 / 規格 | 採用理由與規範 |
|---|---|---|---|
| **後端語言** | Go | 1.22+ | 高效能、靜態強型別、豐富的標準庫、簡潔的並行模型與部署單一二進位檔。 |
| **模組管理** | Go Modules | `go 1.22` | 專案標準依賴管理，保證建置可重現性。 |
| **程式碼規範** | Idiomatic Go | `gofmt`, `golangci-lint` | 遵從標準 Go 慣例（接受 interface、回傳 struct、顯式錯誤處理、禁止 global 狀態）。 |

---

## 2. 系統架構風格

- **架構模型**: 六角架構 / 清潔架構 (Hexagonal / Clean Architecture / Ports and Adapters)
  - `internal/domain`: 核心業務實體與領域規則（Company, Contact, Interaction, Opportunity, SalesStage, User/Role），不依賴任何外部框架或資料庫。
  - `internal/service` (或 `internal/usecase`): 業務流程用例（客戶管理、歷程追蹤、商機階段推進、管線聚合），定義 Repository 與通知介面。
  - `internal/adapter/http`: RESTful API Handlers、中介軟體（Auth、RBAC、Logger）、靜態前端檔案託管。
  - `internal/adapter/repository`: 資料庫存取實作（SQLite SQL 查詢、事務處理）。
  - `cmd/crm-server`: 應用程式組裝與進入點（Dependency Injection、優雅關機）。

---

## 3. 前後端與 HTTP 服務

| 元件 | 技術選型 | 說明 |
|---|---|---|
| **HTTP 路由** | `go-chi/chi/v5` + `net/http` | 100% 相容 Go `net/http` 標準介面，提供優秀的中介軟體鏈、參數路由與高效比對。 |
| **API 規範** | RESTful JSON API | 遵循 OpenAPI 3.x 契約規範，統一端點風格與狀態碼（200, 201, 400, 401, 403, 404, 422, 500）。 |
| **錯誤模型** | 結構化 JSON 錯誤 | 包含 `error.code`, `error.message`, `error.details`，防範敏感技術細節外露。 |
| **前端架構** | HTML5 + ES Modules + Tailwind CSS | 直接整合 PM 原型頁面，無需複雜打包工具，後端透過 `http.FileServer` 或 Go `embed.FS` 統一提供存取。 |

---

## 4. 資料持久化與儲存

| 元件 | 技術選型 | 說明 |
|---|---|---|
| **資料庫引擎** | SQLite | 支援純本機檔案（如 `crm.db`）與純記憶體模式（`:memory:`）。 |
| **資料庫驅動** | `modernc.org/sqlite` (或 `mattn/go-sqlite3`) | 零 CGO 依賴或標準 CGO 支援，支援交易 (Transactions) 與 Foreign Key 限制檢查。 |
| **Schema 管理** | 原生 SQL 遷移腳本 (Migration) | 定義於 `internal/adapter/repository/schema.sql`，具備冪等性建立能力。 |
| **資料庫設定** | 啟用 WAL 模式 | `PRAGMA journal_mode=WAL; PRAGMA foreign_keys=ON;` 保證併發安全與關聯完整性。 |

---

## 5. 自動化 BDD 與測試技術堆疊

| 測試層級 | 工具 / 框架 | 規範與責任邊界 |
|---|---|---|
| **後端 BDD 規格測試** | `github.com/cucumber/godog` (v0.14+) | 以 Gherkin 特徵檔直接驅動後端 API 與領域層驗證，覆蓋率 100% 綁定 PM 驗收規則。每個 Scenario 使用獨立 SQLite `:memory:` 實例。 |
| **後端單元與整合測試** | Go 原生 `testing` + `net/http/httptest` | 針對領域實體規則、複雜計算（金額匯總、權限過濾）與 HTTP Handlers 進行單元驗證。 |
| **前端 / 端到端 (E2E)** | Playwright (`@playwright/test`) | 驅動無頭瀏覽器，自動化驗證使用者介面旅程（Kanban 看板渲染、表單送出、角色切換）。 |

---

## 6. 認證、授權 (RBAC) 與資料可視性

| 角色 | 英文識別代碼 | 權限範圍 |
|---|---|---|
| **業務主管** | `SalesManager` (如 Carol) | - 檢視與維護所有客戶及公司資料<br>- 檢視全團隊所有銷售機會與管線看板<br>- 查看全管線加總業績統計<br>- 可推進任意銷售機會階段或重新開啟結案商機 |
| **一般業務** | `SalesRep` (如 Alice, Bob) | - 檢視全團隊基本客戶清冊，僅維護自己負責之客戶<br>- 建立與推進自己負責之銷售機會 (禁止修改他人機會)<br>- 銷售管線僅限檢視自己名下負責之商機 |
| **認證傳遞** | Header 傳遞機制 | 支援 `Authorization: Bearer <token>` 及開發測試用之 `X-User-Role: SalesManager` / `X-User-Name: Alice` 角色切換。 |
