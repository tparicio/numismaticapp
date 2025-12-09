# AGENTS.md - Master Context for NumismaticApp

## Project Context
NumismaticApp is a web application designed to manage a numismatic coin collection. It is deployed on a QNAP NAS using Docker. The core workflow involves uploading coin images (front/back), processing them (cropping/rotation), analyzing them with Google Gemini AI to extract metadata, and storing the data in PostgreSQL for visualization in a Vue.js frontend.

## Tech Stack Rules
- **Backend:** Golang (Latest Stable). Framework: Fiber v2/v3.
- **Database:** PostgreSQL 16.
- **ORM/Queries:** **ALWAYS use SQLC**. **NEVER use GORM**.
- **Driver:** pgx/v5.
- **Architecture:** Strict Domain-Driven Design (DDD) / Hexagonal Architecture.
- **Image Processing:** libvips (via `h2non/bimg`).
- **AI:** Google Gemini API (Vision model).
- **Frontend:** Vue 3 (Composition API) + Vite.
- **Styling:** Tailwind CSS + DaisyUI.
- **Deployment:** Docker & Docker Compose (Multi-stage builds).

## Directory Map
- `/cmd/api/`: Application entrypoint (`main.go`).
- `/internal/domain/`: Pure domain logic, entities, and repository interfaces. No external dependencies.
- `/internal/application/`: Use cases and business logic orchestration. Depends on domain.
- `/internal/infrastructure/`: Implementation of interfaces (Postgres, Gemini, FileSystem, Vips).
  - `/internal/infrastructure/db/`: **Generated** SQLC code. DO NOT EDIT MANUALLY.
- `/internal/api/`: Fiber handlers, Router, Middleware. Adapters for the HTTP layer.
- `/web/`: Vue.js frontend application.
- `/deployment/`: Docker and deployment scripts.
- `/migrations/`: SQL migration scripts.

## Code Style & Conventions
- **Golang:**
  - Follow standard Go conventions (`gofmt`).
  - Error handling: Wrap errors with context (e.g., `fmt.Errorf("failed to process image: %w", err)`).
  - Variable names: CamelCase.
  - Interfaces defined in `domain`, implemented in `infrastructure`.
- **Vue.js:**
  - Use **Composition API** (`<script setup>`).
  - Use TypeScript if possible, otherwise strict JavaScript.
  - Components should be small and focused.
  - Use DaisyUI classes for UI components.

## Maintenance Notes
- **SQLC Generation:** Run `sqlc generate` after modifying `schema.sql` or queries.
- **Migrations:** Add new SQL files to `/migrations/` with sequential numbering (e.g., `002_add_index.up.sql`).
- **Gemini Prompting:** When modifying the AI prompt, ensure the `vertical_correction_angle` requirement remains prominent.
