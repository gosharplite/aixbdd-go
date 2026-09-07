# Truth Delta: 001-crm-core

**Plan Package**: `specs/plans/001-crm-core`
**Truth Root**: `specs/truth`

## /technical-research

| 動作 | Truth 規格 | 改動摘要 | 原因 |
| --- | --- | --- | --- |
| ADD | `specs/truth/techstack.md` | 建立系統唯一技術堆疊真相：Go 1.22+、Hexagonal Architecture、Chi/net/http、SQLite (支援 :memory: 測試)、Godog 後端 BDD、Playwright 前端 E2E、輕量前端靜態整合與宣告式 RBAC。 | 001-crm-core 起始專案確立系統端點架構、資料持久化與雙層 BDD 測試技術選型。 |

## /api-plan

| 動作 | Truth 規格 | 改動摘要 | 原因 |
| --- | --- | --- | --- |
| ADD | `specs/truth/contracts/openapi.yaml` | 建立系統唯一 API 契約真相 (OpenAPI 3.0)，涵蓋 /me、/companies、/contacts、/contacts/{id}/interactions (倒序歷程)、/opportunities (含 /stage 嚴格 RBAC 檢查)、/pipeline (各階段統計與進行中總額彙總) 及標準錯誤模型。 | 001-crm-core 建立 CRM 核心 RESTful API 規格。 |

## /data-plan

| 動作 | Truth 規格 | 改動摘要 | 原因 |
| --- | --- | --- | --- |
| ADD | `specs/truth/data/schema.dbml` | 建立系統資料模型真相，包含 users, companies, contacts, interactions, opportunities 五大實體表，定義 sales_stage、user_role、interaction_type 列舉，以及外鍵防刪約束與 Append-only 歷程規範。 | 001-crm-core 建立 CRM 核心資料實體與關聯規範。 |

## /dsl-refine

| 動作 | Truth 規格 | 改動摘要 | 原因 |
| --- | --- | --- | --- |
| ADD | `specs/truth/features/backend/**` | 建立後端可執行 Gherkin 特徵檔與模組 DSL (customer, interaction, opportunity, pipeline) 及根共用 DSL，全面落實原子化單一 Act 驗收規格。 | 001-crm-core 承接 PM 驗收標準並拆解為可自動化執行之後端 BDD 規格。 |
| ADD | `specs/truth/features/frontend/**` | 建立前端 Playwright BDD 可執行特徵檔與模組 DSL (pipeline, opportunity) 及根共用 DSL，覆蓋管線看板多角色切換與步進器權限攔截。 | 001-crm-core 建立前端自動化端到端驗收規範。 |
