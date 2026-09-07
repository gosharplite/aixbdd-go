# 技術研究與架構決策：001-crm-core

**Plan Package**: `specs/plans/001-crm-core`  
**建立日期**: 2026-03-30  
**負責角色**: RD (`/technical-research`)  
**狀態**: 已確認 (Approved)

---

## 1. 需求與研究背景

依據 `specs/plans/001-crm-core/spec.md` 及 4 份 PM 驗收 Feature 規格（客戶與公司管理、聯絡歷程、銷售機會推進、銷售管線綜覽），本專案目標為從零建立高可靠度的 CRM 核心系統，涵蓋後端 RESTful API 與現代化網頁介面。

本次技術研究聚焦於以下三大核心決策及相關支援技術：
1. **系統端點與架構劃分**
2. **BDD 自動化測試技術堆疊**
3. **資料持久化與測試驗證策略**
4. **身份驗證、權限控制（RBAC）與資料隔離機制**

---

## 2. 核心技術決策

### 決策 1：系統端點與前後端架構劃分

- **決策結果**: 採用 **Go 後端 RESTful API 服務 + 輕量靜態前端資源託管 (HTML5 / ES Modules / Tailwind CSS)** 單一應用架構。
- **架構特點**:
  - 後端以 Go 原生高效 HTTP 服務（採用 `net/http` 搭配輕量且標準相容的 `go-chi/chi/v5` 路由器）提供完整的 RESTful JSON API。
  - 前端直接整合並深化 PM 已交付且驗收通過的 4 份高擬真 HTML 原型（`ui/*.html`），透過標準 ES Modules 與 `fetch` API 進行前後端非同步資料交換。
  - 後端靜態檔案伺服器直接託管前端靜態檔案（亦支援 `embed.FS` 嵌入打包），產生單一可執行檔，部署無繁瑣 Node.js 運行時依賴。
- **替代方案比較**:
  - *方案 B (Go API + 獨立 Vite/React SPA)*: 需要額外維護前端 Node.js 建置管線、TypeScript 型別重複定義，且需將現有 4 頁完整 HTML 打碎重寫為 React 元件，增加前期迭代成本。**[未採用]**
  - *方案 C (Go SSR / html/template + HTMX)*: 雖然免除 API 序列化，但與 `spec.md`（FR-017 規定必須具備獨立 RESTful API）及未來行動端或第三方整合契約不符。**[未採用]**

### 決策 2：BDD 自動化測試技術堆疊

- **決策結果**: 採用 **雙層 BDD 驗證策略：後端 Go Godog (Cucumber for Go) + 前端/全端 Playwright BDD**。
- **架構特點**:
  - **後端 BDD**: 使用 `github.com/cucumber/godog`，直接將 PM 的 Gherkin 驗收規格及 RD 拆解的後端特徵檔（`specs/truth/features/backend/*.feature`）綁定到 Go Step Definitions，直接呼叫 HTTP Handler 或 Domain Service 驗證狀態移轉、資料完整性與錯誤碼（HTTP 401/403/404/422）。
  - **前端/E2E BDD**: 使用 Playwright 驅動無頭瀏覽器，自動化驗證使用者真實互動歷程（例如從管線 Kanban 拖拉/檢視、角色身分切換 Alice vs Carol、新增聯絡歷程並即時反映於時間軸）。
- **替代方案比較**:
  - *方案 B (純 Playwright-BDD E2E)*: 缺乏後端單元與契約層級的快速回饋（回饋循環較慢），不利於 TDD Red-Green-Refactor 極速推進。**[未採用]**
  - *方案 C (純 Go Table-Driven Test，無 Gherkin 引擎)*: 無法直接以 PM 可閱讀之 Gherkin 檔案作為測試真相，失去 AIBDD 流程之溯源性。**[未採用]**

### 決策 3：資料持久化與測試驗證策略

- **決策結果**: 採用 **SQLite (純 Go 實作驅動 `modernc.org/sqlite` 或 `mattn/go-sqlite3`) + Clean Architecture Repository Pattern**。
- **架構特點**:
  - 生產/日常開發運行支援本地檔案資料庫（如 `crm.db`），零外部資料庫依賴，開箱即用。
  - 自動化測試（Godog 與單元測試）支援 SQLite In-Memory 模式 (`file::memory:?cache=shared`)，每個測試情境（Scenario）皆能秒級啟動、獨立重設並具備完全資料隔離。
  - 核心領域模型與儲存實作透過 Go Interface 隔離，未來若需要無縫遷移至 PostgreSQL / MySQL 僅需抽換 Repository 實作。
- **替代方案比較**:
  - *方案 B (強制依賴 Docker PostgreSQL)*: 增加本機測試與 CI 的啟動門檻與網路延遲，無法做到開箱即跑。**[未採用]**
  - *方案 C (僅使用純記憶體 Slice/Map Mock)*: 無法驗證真實 SQL 交易、Foreign Key 限制（如公司防刪除規則）與關聯查詢邏輯。**[未採用]**

### 決策 4：身份驗證、權限控制 (RBAC) 與資料可視性

- **決策結果**: 採用 **Bearer Token / 輕量 Session Context + 宣告式 RBAC 中介軟體 (Middleware)**。
- **架構特點**:
  - 支援預設測試帳號切換（如 `Alice` [業務人員 Sales Rep]、`Bob` [其他業務人員]、`Carol` [業務主管 Sales Manager]）。
  - 後端 Context 注入 `CurrentActor`（包含 `UserID`, `Name`, `Role`）。
  - 嚴格落實權限矩陣：
    1. **公司與聯絡人資料**: 訪客/未登入一律拒絕 (401/403)。全團隊成員可瀏覽基本清冊，但僅負責業務與主管有修改權限。
    2. **銷售機會階段推進與金額修改**: 僅負責業務本人或業務主管具備推進權限，其餘業務嘗試修改回傳 403 Forbidden。
    3. **銷售管線綜覽 (Pipeline)**: 業務主管（Carol）能查看全公司所有人員機會與整體加總金額；一般業務（Alice）僅能查看自己名下負責之機會。

---

## 3. 風險評估與緩解措施

| 風險項目 | 影響程度 | 緩解措施 |
|---|---|---|
| **SQLite 併發寫入鎖定 (busy/lock)** | 中 | 啟用 WAL 模式 (`PRAGMA journal_mode=WAL;`)，後端連線池配置合理之 busy_timeout，足夠應付 CRM 核心場景與高併發測試。 |
| **Gherkin 步驟與 API 實作細節耦合** | 中 | 嚴格依循 AIBDD DSL 規範，Gherkin 維持業務與網域語意，後端 Step Definitions 負責對齊 REST API 請求/回應與 DB 狀態。 |
| **前端靜態檔案快取問題** | 低 | 開發模式關閉快取，發布模式透過版本雜湊或 ETag 管理。 |

---

## 4. 下一步委派

技術決策已完整收斂，後續將依 AIBDD SOP 推進：
1. 建立系統層級唯一的技術真相：`specs/truth/techstack.md`。
2. 啟動 `/system-analysis` 規劃系統 Wave 與介面邊界，並委派 `/api-plan` 與 `/data-plan`。
