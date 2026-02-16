# 🏛 Root Cursor Rules – Legislative System Monorepo

## 🎯 Project Philosophy

This monorepo contains multiple applications (backend services,
frontend, shared libraries).

All code must follow:

-   Clean Architecture
-   Domain-first modeling
-   Strict Test-Driven Development (TDD)
-   Modular design
-   Explicit boundaries
-   OOP principles implemented idiomatically (interfaces + composition)

No shortcuts. No speculative architecture.

------------------------------------------------------------------------

# 🧱 1️⃣ Architecture (Non-Negotiable)

Every backend service must follow this structure:

    /internal
        /domain
        /application
        /infrastructure
        /api

## Layer Responsibilities

### Domain

-   Pure business logic
-   No framework imports
-   No DB logic
-   No HTTP logic
-   Defines repository interfaces
-   Protects invariants
-   Owns state transitions

### Application

-   Orchestrates domain
-   Uses domain interfaces
-   No direct DB access
-   No framework-specific types
-   No business rules (only coordination)

### Infrastructure

-   Implements domain interfaces
-   Contains DB, messaging, external systems
-   Converts persistence models ↔ domain entities

### API

-   HTTP handlers
-   Converts DTO ↔ application commands
-   No business logic
-   No persistence logic

Violation of layering is forbidden.

------------------------------------------------------------------------

# 🧠 2️⃣ Domain Modeling Rules

1.  Domain must be framework-agnostic.
2.  Entities must protect their invariants.
3.  All state transitions must be explicit methods.
4.  No public mutable fields.
5.  IDs must be stable identifiers (UUID recommended).
6.  Display numbering (e.g., § 4) must not be identity.
7.  Domain must not depend on infrastructure.

Example:

``` go
func (a *Act) StartVoting() error
func (a *Act) Accept() error
```

Direct status mutation from outside is forbidden.

------------------------------------------------------------------------

# 🧩 3️⃣ Interface & OOP Discipline

Even though Go is not classical OOP:

-   Use composition over inheritance.
-   Define interfaces where behavior is required.
-   Interfaces live in the domain layer.
-   Infrastructure implements domain interfaces.
-   Avoid fat structs and god services.

Each struct must have a single responsibility.

------------------------------------------------------------------------

# 🧪 4️⃣ STRICT TDD RULES (Mandatory)

Every feature must follow this exact cycle:

## Step 1 – Write ONE failing test

-   The test must describe behavior.
-   The test must fail.
-   No implementation before failing test.

## Step 2 – Minimal implementation

-   Write the smallest code necessary to pass.
-   No extra functionality.
-   No premature abstractions.

## Step 3 – Refactor

-   Improve naming.
-   Remove duplication.
-   Extract clear responsibilities.
-   Keep tests green.

## Step 4 – Repeat

------------------------------------------------------------------------

## 🚫 Absolutely Forbidden

-   Writing multiple failing tests at once.
-   Implementing multiple behaviors in one step.
-   Writing implementation before a failing test.
-   Refactoring while tests are red.
-   Adding future-proof abstractions without immediate need.

------------------------------------------------------------------------

# 🧪 5️⃣ Testing Rules

## Domain Tests

-   No mocks.
-   Pure logic tests.
-   Behavior-focused.

## Application Tests

-   Mock repositories.
-   Test orchestration only.

## Infrastructure Tests

-   Optional integration tests.
-   Must not leak into domain tests.

## API Tests

-   Only after domain logic is stable.

------------------------------------------------------------------------

# 📦 6️⃣ Modular Feature Development Flow

When implementing a feature:

1.  Start in the domain layer.
2.  Add a failing test.
3.  Implement minimal domain logic.
4.  Refactor.
5.  Add application orchestration test.
6.  Implement minimal application service.
7.  Only then implement infrastructure.
8.  API last.

Never start from HTTP.

------------------------------------------------------------------------

# 🧼 7️⃣ Refactoring Discipline

Refactor only when:

-   Tests pass.
-   Duplication appears.
-   Behavior becomes unclear.
-   A clear abstraction emerges.

Do NOT: - Generalize for future jurisdictions yet. - Add patterns “just
in case.” - Introduce generics without clear duplication.

------------------------------------------------------------------------

# 🧨 8️⃣ Anti-Patterns Prohibited Globally

-   Business logic inside handlers
-   Domain importing database packages
-   Returning DB models to API
-   Mutating entity state from outside
-   Large multi-responsibility services
-   Skipping tests for “simple” logic
-   Speculative microservices

------------------------------------------------------------------------

# 📌 9️⃣ MVP Constraint

Focus on correctness before scale.

No: - Multi-tenant complexity - Cross-service orchestration - Event
buses - Advanced abstraction layers

Build a correct legislative engine first.

------------------------------------------------------------------------

# 🧭 10️⃣ Code Generation Constraint for Cursor

When generating code:

-   Never generate large multi-feature implementations.
-   Always start with a failing test.
-   Generate only one small behavior per response.
-   Wait for confirmation before continuing to the next step.
-   Prefer explicitness over cleverness.
-   Avoid magic.

------------------------------------------------------------------------

# 🧠 11️⃣ Monorepo Consistency Rule

App-level `.cursor/rules.md` files:

-   May extend these rules.
-   May add technology-specific constraints.
-   Must NOT violate architecture or TDD rules defined here.

This file is the constitutional layer of the monorepo.

------------------------------------------------------------------------

# 🔥 12️⃣ Core Principle

Build a correct domain model.

Everything else (DB, API, frontend) is an adapter.
