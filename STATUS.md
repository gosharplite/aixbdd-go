# Project Status & Handoff

> Current development status, completed artifacts, and next steps for the AIBDD workflow.

---

## Current Plan Package: `001-crm-core`

The initial iteration `001-crm-core` has completed the PM Definition, Acceptance Criteria, and Prototyping phases.

### Completed Artifacts

| Artifact | Path | Status | Owner |
|---|---|---|---|
| **Specification** | `specs/plans/001-crm-core/spec.md` | Completed (Draft) | PM (`/specify`) |
| **Requirements Checklist** | `specs/plans/001-crm-core/checklists/requirements.md` | Completed (Ready) | PM (`/specify`) |
| **Truth Delta Ledger** | `specs/plans/001-crm-core/truth-delta.md` | Initialized | Shared (`/truth-delta`) |
| **Acceptance Features** | `specs/plans/001-crm-core/features/acceptance/*.feature` | 4 Feature files verified | PM (`/spec-by-example`) |
| **UI Control Plan** | `specs/plans/001-crm-core/ui/ui-plan.md` | Completed | PM (`/ui-plan`) |
| **Interactive Prototypes** | `specs/plans/001-crm-core/ui/*.html` (`index.html` entry) | 4 Screens verified | PM (`/ui-plan`) |

---

## Acceptance Features Overview

1. **`specs/plans/001-crm-core/features/acceptance/客戶與公司資料管理.feature`**
   * Company and contact person profile creation and editing.
   * RBAC enforcement (unauthenticated or unauthorized modifications rejected).

2. **`specs/plans/001-crm-core/features/acceptance/記錄聯絡歷程與下一步行動.feature`**
   * Interaction logging (call, meeting, email) and next action assignment.
   * Chronological timeline retrieval (reverse chronological order).

3. **`specs/plans/001-crm-core/features/acceptance/銷售機會建立與階段推進.feature`**
   * Opportunity creation with deal value, expected close date, and standard stages.
   * Stage advancement restricted to assigned owner (*Alice*) or manager (*Carol*).
   * Rejection of unauthorized stage alterations by non-owner rep (*Bob*).
   * Differentiation of terminal closed states (*Closed Won* vs. *Closed Lost*).

4. **`specs/plans/001-crm-core/features/acceptance/銷售管線綜覽與業績預測.feature`**
   * Pipeline Kanban board overview for sales managers with stage aggregations and sums.
   * Role-based data isolation restricting general sales reps to viewing their assigned deals.

---

## Interactive HTML Prototypes Overview

* **Entry Point**: `specs/plans/001-crm-core/ui/index.html`
* **Screens**:
  1. `ui/index.html`: Pipeline Kanban board with 5 stages (*Prospecting*, *Contacted*, *Proposal*, *Closed Won*, *Closed Lost*), aggregate metrics, and role switcher (*Carol [Manager]* vs. *Alice [Rep]*).
  2. `ui/customers.html`: Organizations and contacts directory with search and creation modals.
  3. `ui/customer-detail.html`: Customer profile, prominent next action card, and live interaction logging timeline.
  4. `ui/opportunity-detail.html`: 5-step visual stage stepper with live role simulation testing happy-path progression and RBAC rejection.

---

## Next Steps for RD

1. **Review Upstream Artifacts**:
   * Review `specs/plans/001-crm-core/spec.md` and `features/acceptance/*.feature`.
   * Review interactive prototype flow in `specs/plans/001-crm-core/ui/index.html`.
2. **Execute `/technical-research`**:
   * Analyze architecture options, produce `specs/plans/001-crm-core/research.md`, and update `specs/truth/techstack.md`.
3. **Execute `/system-analysis`**:
   * Generate `specs/plans/001-crm-core/plan.md` and delegate to `/api-plan` and `/data-plan`.
4. **Execute `/dsl-refine`**:
   * Decompose acceptance journeys into executable frontend/backend Gherkin and DSL under `specs/truth/features/`.
5. **Execute `/tasks`**:
   * Break down into automated test alignment and implementation tasks in `specs/plans/001-crm-core/tasks.md`.
6. **Execute `/implement`**:
   * Deliver implementation via One-Shot execution with full automated test coverage.
