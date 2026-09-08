# Project Status & Handoff

> Current development status, completed artifacts, and next steps for the AIBDD workflow.

---

## Current Plan Package: `001-crm-core`

**Status**: **DELIVERED & COMMITTED** (`0038ace`)  
The initial iteration `001-crm-core` (CRM Core Customer Relationship Management System) has completed all PM and RD phases of the AIBDD workflow—from specification and acceptance testing to system design, BDD automation, implementation, and verification.

---

### Completed Artifacts

| Phase | Artifact | Path | Status | Owner |
|---|---|---|---|---|
| **PM** | **Specification** | `specs/plans/001-crm-core/spec.md` | Completed | PM (`/axb-specify`) |
| **PM** | **Requirements Checklist** | `specs/plans/001-crm-core/checklists/requirements.md` | Completed (100% Ready) | PM (`/axb-specify`) |
| **PM** | **Acceptance Features** | `specs/plans/001-crm-core/features/acceptance/*.feature` | 4 Feature files verified | PM (`/axb-spec-by-example`) |
| **PM** | **UI Control Plan** | `specs/plans/001-crm-core/ui/ui-plan.md` | Completed | PM (`/axb-ui-plan`) |
| **PM** | **Interactive Prototypes** | `specs/plans/001-crm-core/ui/*.html` (`index.html` entry) | 4 Screens verified | PM (`/axb-ui-plan`) |
| **RD** | **Technical Research** | `specs/plans/001-crm-core/research.md` | Completed | RD (`/axb-technical-research`) |
| **RD** | **Tech Stack Truth** | `specs/truth/techstack.md` | Completed (Canonical Truth) | RD (`/axb-technical-research`) |
| **RD** | **System Analysis Plan** | `specs/plans/001-crm-core/plan.md` | Completed | RD (`/axb-system-analysis`) |
| **RD** | **Data Model Truth** | `specs/truth/data/schema.dbml` | Completed (Canonical Truth) | RD (`/axb-data-plan`) |
| **RD** | **API Contract Truth** | `specs/truth/contracts/openapi.yaml` | Completed (Canonical Truth) | RD (`/axb-api-plan`) |
| **RD** | **Executable Features & DSL** | `specs/truth/features/**` | Completed (Canonical Truth) | RD (`/axb-dsl-refine`) |
| **Shared** | **Truth Delta Ledger** | `specs/plans/001-crm-core/truth-delta.md` | Completed (All Owners ADD) | Shared (`/axb-truth-delta`) |
| **RD** | **BDD Development Tasks** | `specs/plans/001-crm-core/tasks.md` | Completed (21/21 Tasks Done) | RD (`/axb-tasks`) |
| **RD** | **Full Implementation & Verification** | `cmd/server/main.go`, `internal/**`, `tests/bdd/**` | 100% Tests Passing (12/12 Scenarios, 65/65 Steps Green) | RD (`/axb-implement`) |

---

## Acceptance Features & Automated BDD Coverage

All 4 PM acceptance feature files have been mapped 1:1 to executable BDD specifications and automated tests:

1. **`specs/plans/001-crm-core/features/acceptance/客戶與公司資料管理.feature`**
   * Company creation, editing, and contact management.
   * Format validation (email regex) and mandatory fields.
   * Unauthenticated access rejected with HTTP 401 Unauthorized.
   * **Verification**: ✅ 100% Passing (`tests/bdd/customer_steps.go`)

2. **`specs/plans/001-crm-core/features/acceptance/記錄聯絡歷程與下一步行動.feature`**
   * Append-only interaction logging (Call, Meeting, Email).
   * Reverse-chronological timeline retrieval (`ORDER BY interaction_time DESC`).
   * Latest next action and date synchronization on customer profiles.
   * **Verification**: ✅ 100% Passing (`tests/bdd/interaction_steps.go`)

3. **`specs/plans/001-crm-core/features/acceptance/銷售機會建立與階段推進.feature`**
   * Deal value, expected close date, and standard 5-stage workflow.
   * Stage advancement strictly restricted to opportunity owner (*Alice*) or manager (*Carol*).
   * Unauthorized modification attempts by non-owner rep (*Bob*) rejected with HTTP 403 Forbidden.
   * Differentiation of terminal closed states (*Closed Won* vs. *Closed Lost*).
   * **Verification**: ✅ 100% Passing (`tests/bdd/opportunity_steps.go`)

4. **`specs/plans/001-crm-core/features/acceptance/銷售管線綜覽與業績預測.feature`**
   * Pipeline Kanban board summary with stage count and total amount aggregations.
   * Active pipeline sum (Prospecting + Contacted + Proposal) calculation.
   * Role-based data isolation (Sales Manager sees full team; Sales Rep restricted to own deals).
   * **Verification**: ✅ 100% Passing (`tests/bdd/pipeline_steps.go`)

---

## How to Run & Verify

### 1. Run All Automated BDD & Unit Tests
```bash
# Run all tests with verbose output
go test -v ./...

# Run all tests with Go race detector enabled
go test -race ./...
```

### 2. Start the Production / Local Server
```bash
# Start server on default port 8080 (auto-serves UI and API)
go run cmd/server/main.go
```
* **Web UI Entry Point**: `http://localhost:8080/`
* **API Endpoints**: `http://localhost:8080/api/pipeline`, `http://localhost:8080/api/me`, `http://localhost:8080/api/companies`, etc.

---

## Workflow Status & Handoff

* **Current Iteration (`001-crm-core`)**: Completed, verified, and committed.
* **Next Steps**:
  * For new requirements or next iteration features, start a new plan package using `/axb-specify` (e.g., `specs/plans/002-<feature-slug>/`).
  * System truth specifications in `specs/truth/` remain the authoritative source of truth for all subsequent iterations.
