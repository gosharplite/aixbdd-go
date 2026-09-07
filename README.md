# CRM Core Platform (CRM 客戶關係管理系統)

> A high-reliability CRM core platform developed using the **AIBDD (AI × BDD)** methodology. Built with an idiomatic Go backend following Clean Architecture, SQLite persistence, and a 100% automated Gherkin BDD test suite powered by Cucumber Godog.

---

## 🌟 Key Features

* **🏢 Company & Customer Contact Management**:
  * Create, update, and manage business organizations and individual contact profiles.
  * Email format validation and company delete protections (`delete: restrict`).
  * Authentication enforcement (unauthorized access rejected with HTTP 401).

* **📞 Interaction History & Next Action Tracking**:
  * Append-only logging for multi-channel touchpoints (*Call*, *Meeting*, *Email*).
  * Reverse-chronological timeline retrieval (`interaction_time DESC`).
  * Real-time next action and scheduled follow-up date synchronization.

* **💼 Opportunity Progression & Role-Based Access Control (RBAC)**:
  * 5 standard sales stages (*Prospecting*, *Contacted*, *Proposal*, *Closed Won*, *Closed Lost*).
  * Strict RBAC permissions: only the assigned owner (*Alice*) or sales manager (*Carol*) can advance stages or adjust amounts.
  * Rejection of unauthorized modifications by other reps (*Bob*) with HTTP 403 Forbidden.
  * Clear terminal closed states (*Closed Won* vs. *Closed Lost*).

* **📊 Sales Pipeline Kanban & Forecasting**:
  * Team-wide pipeline Kanban board with stage counts and aggregated deal sums.
  * Active pipeline sum calculation (*Prospecting + Contacted + Proposal*).
  * Role-based data isolation (Sales Manager sees full team; Sales Rep restricted to personal deals).

---

## 🛠️ Technology Stack

| Component | Technology | Specification / Version |
|---|---|---|
| **Language & Runtime** | Go | 1.22+ |
| **HTTP Routing** | Chi Router | `github.com/go-chi/chi/v5` (`net/http` compatible) |
| **Persistence** | SQLite | `modernc.org/sqlite` (Pure Go, WAL mode, foreign keys enabled) |
| **Architecture** | Clean Architecture | Hexagonal (Domain, Service, HTTP/DB Adapters) |
| **BDD Testing** | Godog | `github.com/cucumber/godog` (Cucumber for Go) |
| **Unit & Integration** | Go Testing | Native `testing`, `httptest`, with `-race` detection |
| **Frontend UI** | Modern Web | HTML5, CSS3, ES Modules, Tailwind CSS |
| **API Contract** | OpenAPI 3.0 | `specs/truth/contracts/openapi.yaml` |
| **Data Schema** | DBML | `specs/truth/data/schema.dbml` |

---

## 📁 Project Directory Layout

```text
.
├── cmd/
│   └── server/
│       └── main.go                  # Application entrypoint & HTTP server
├── internal/
│   ├── domain/                      # Core business entities & repository interfaces
│   ├── service/                     # Business logic use-cases (customer, opp, pipeline, interaction)
│   └── adapter/
│       ├── http/                    # RESTful API handlers, middleware, & static file serving
│       └── repository/              # SQLite repository implementation & schema.sql
├── specs/
│   ├── plans/
│   │   └── 001-crm-core/            # Iteration plan package
│   │       ├── spec.md              # Functional specifications
│   │       ├── research.md          # Technical research & architectural decisions
│   │       ├── plan.md              # System analysis wave plan
│   │       ├── tasks.md             # Executable BDD task decomposition
│   │       ├── truth-delta.md       # Truth modification ledger
│   │       ├── features/acceptance/ # PM acceptance Gherkin feature files
│   │       └── ui/                  # Interactive HTML prototypes & control plan
│   └── truth/                       # Canonical system-wide truth artifacts
│       ├── techstack.md             # Canonical tech stack specification
│       ├── contracts/openapi.yaml   # Canonical OpenAPI 3.0 contract
│       ├── data/schema.dbml         # Canonical DBML database schema
│       └── features/                # Executable Gherkin features & DSL dictionaries
├── tests/
│   └── bdd/                         # Godog test suite & step definitions
├── AIBDD.md                         # AIBDD methodology and workflow SOP
├── SESSION-BOOTSTRAP.md             # AI session orientation guide
├── STATUS.md                        # Living project status and delivery ledger
├── go.mod
└── go.sum
```

---

## 🚀 Getting Started

### Prerequisites

* **Go**: Version 1.22 or higher (`go version`)

### 1. Run Automated BDD & Unit Tests

The test suite runs with isolated in-memory SQLite (`:memory:`), requiring zero external dependencies:

```bash
# Run all tests with verbose output
go test -v ./...

# Run all tests with Go race detection enabled
go test -race ./...
```

### 2. Start the Local Server

```bash
# Start the server (default: port 8080)
go run cmd/server/main.go
```

Once started, open your browser to interact with the platform:

| Page | URL | Description |
|---|---|---|
| **Pipeline Kanban** | [http://localhost:8080/](http://localhost:8080/) | Sales pipeline overview with role switcher (Carol vs. Alice) |
| **Customers Directory** | [http://localhost:8080/customers.html](http://localhost:8080/customers.html) | Company directory and customer contact person management |
| **Customer Detail** | [http://localhost:8080/customer-detail.html](http://localhost:8080/customer-detail.html) | Interaction timeline and next action schedule |
| **Opportunity Progression** | [http://localhost:8080/opportunity-detail.html](http://localhost:8080/opportunity-detail.html) | 5-step visual stage stepper with live RBAC simulation |

---

## 📡 RESTful API Overview

All API endpoints follow OpenAPI 3.0 specification (`specs/truth/contracts/openapi.yaml`):

* `GET /api/me`: Returns current authenticated actor identity and role.
* `GET /api/pipeline`: Returns grouped stages with deal counts, stage sums, and active total.
* `GET /api/companies`, `POST /api/companies`: List and create companies.
* `GET /api/contacts`, `POST /api/contacts`: List and create customer contacts.
* `GET /api/contacts/{id}/interactions`, `POST /api/contacts/{id}/interactions`: Log and retrieve interactions.
* `GET /api/opportunities`, `POST /api/opportunities`: List and create sales opportunities.
* `PATCH /api/opportunities/{id}/stage`: Advance opportunity stage (enforces RBAC permissions).

---

## 📖 Methodology & Governance

* **Workflow SOP**: See [`AIBDD.md`](AIBDD.md) for the complete PM/RD division of responsibilities.
* **Status & Handoff**: See [`STATUS.md`](STATUS.md) for current iteration status and completed artifacts.
* **Session Bootstrap**: See [`SESSION-BOOTSTRAP.md`](SESSION-BOOTSTRAP.md) for pre-loaded context and quick start instructions.
