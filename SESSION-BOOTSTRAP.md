# AI Session Bootstrap: AIBDD Workflow

> **READ THIS FIRST in any fresh AI session (PM, RD, Architect, or Subagent).**
> All skills and workflows below are **already pre-loaded into your active prompt context**. You do not need to search the filesystem or look up external paths.

---

## 1. Fast Orientation

* **Project Root**: Current workspace (`.`)
* **Workflow SOP**: See the [AIBDD methodology & skills repo](https://github.com/gosharplite/aixbdd)
* **Living Project Status & Handoff**: See [`STATUS.md`](STATUS.md)
* **Pre-Loaded Skills**: All 24+ skill SOPs, rules, and behavioral guidelines are already in memory. Simply invoke them by name (e.g. `/axb-system-analysis`, `/axb-tasks`, `/axb-implement`).

---

## 2. Directory of Pre-Loaded Skills (Already in Context)

### A. PM & Specification Skills (Problem Space)
* **`axb-specify`**: Creates new plan packages (`specs/plans/NNN-<slug>/`), `spec.md`, `checklists/requirements.md`, and initializes `truth-delta.md`.
* **`axb-clarify`**: Structured user interviews for ambiguities. Formats questions with Context, Question, and Options table.
* **`axb-clarify-over-specs`**: Scans `spec.md` for high-impact gaps, interviews user, and updates spec and checklist.
* **`axb-spec-by-example`**: Produces plan-side acceptance Gherkin features (`specs/plans/NNN-<slug>/features/acceptance/*.feature`).
* **`axb-ui-plan`**: Creates UI control plan (`ui-plan.md`) and interactive HTML prototypes (`ui/*.html`, entry at `ui/index.html`).

### B. RD Architecture & Truth Planning Skills (System Space)
* **`axb-technical-research`**: Analyzes technical architecture options, creates `research.md`, and owns `specs/truth/techstack.md`.
* **`axb-system-analysis`**: Orchestrates system boundaries, dependencies, and waves in `plan.md`; delegates to API, data, and UI planners.
* **`axb-api-plan`**: Truth owner for API contracts (`specs/truth/contracts/**`, OpenAPI 3.x).
* **`axb-data-plan`**: Truth owner for data models (`specs/truth/data/**`, DBML schemas).
* **`axb-dsl-refine`**: Truth owner for interface executable feature files and DSL (`specs/truth/features/**`).
* **`axb-truth-delta`**: Common handoff ledger tracking ADD/MODIFY/DELETE/NOOP changes in `truth-delta.md`.
* **`axb-gherkin-and-dsl`**: Audits and refactors Gherkin syntax, Rule/Example boundaries, and DSL mappings.

### C. RD Implementation & Testing Skills (Solution Space)
* **`axb-tasks`**: Decomposes plans into sequential executable BDD tasks in `specs/plans/NNN-<slug>/tasks.md`.
* **`axb-implement`**: One-Shot task executor following "don't stop until deliver" without skipping steps.
* **`axb-bdd`**: Executes BDD slices via `red`, `green`, and `refactor` entry points against step definitions.
* **`golang-patterns`**: Idiomatic Go design patterns, error handling, concurrency, and architecture best practices.
* **`golang-testing`**: Go testing patterns (table-driven tests, subtests, mocks, fuzzing, benchmarks, TDD).

### D. Domain Modeling & Quality Governance
* **`domain-model-author`** / **`domain-model-context`** / **`domain-model-lint`**: Canonical domain models via `modelith`.
* **`axb-constitution`**: Governs modular AI rules under `.agents/constitution/`.
* **`grilling`**: Relentless adversarial questioning, one question at a time.
* **`tmg-chat-ingroup`** / **`tmg-grill-round`** / **`tmg-issue-to-pr`**: Inter-agent messaging, grill rounds, and PR workflow.

---

## 3. How to Start in a Fresh Session

1. Read [`STATUS.md`](STATUS.md) to check the current iteration and completed artifacts.
2. Follow the next step indicated in [`STATUS.md`](STATUS.md) (e.g. for RD on `001-crm-core`, run `/axb-technical-research` or `/axb-system-analysis`).
3. You do not need to search for skills—call them directly according to the workflow.
