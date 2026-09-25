# Calculator

A calculator web app that keeps a history of the last 100 calculations, persisted in PostgreSQL.

## Project structure

```
Calculator/
├── backend/          Go API (chi router, pgx/pgxpool for Postgres)
├── frontend/         React + Vite + TypeScript SPA
├── migrations/       golang-migrate SQL migrations
├── docker-compose.yml  Postgres 16 service
└── Makefile          up / migrate / test targets
```

## Architecture

- **Backend** (`backend/`, Go module `calculator-backend`): a chi HTTP server.
  - `cmd/server/main.go` — entrypoint; reads `DATABASE_URL`/`PORT` from env, opens a `pgxpool.Pool`, starts the server.
  - `internal/db` — Postgres connection pool helper (pgx).
  - `internal/server` — chi router setup: `GET /healthz`, and `/api` routes mounted to the history handler.
  - `internal/calc` — recursive-descent expression evaluator (`+ - * /`, parentheses, decimals) using `shopspring/decimal` for exact arithmetic.
  - `internal/history` — `History` model, HTTP handlers for `POST /api/calculate` and `GET /api/history`, and a Postgres-backed repository that persists each calculation and prunes anything beyond the most recent 100 rows.
- **Frontend** (`frontend/`): a React + Vite SPA — "Grey Matter", a physical-calculator-styled UI (pitch-black/neon-green theme) with a nav bar + slide-out menu drawer, the calculator itself (expression display, full button set, AC/DEL), a history overlay backed by `GET /api/history`, and a footer. Wired to the backend over `/api/*` via a Vite dev-server proxy to `http://localhost:8080`. Vitest is configured for component testing (none written yet).
- **Database**: a single `history` table (`id`, `expression`, `result`, `created_at`) created by the migration in `migrations/`.
- **Postgres**: runs only in Docker (`docker-compose.yml`); the backend and frontend run locally against it during development.

## Setup

### Prerequisites
- Go (see `backend/go.mod` for the module; any recent Go toolchain works)
- Node.js + npm
- Docker (with `docker compose`)
- `make` — not native on Windows; use Git Bash/WSL, or install via `choco install make`

### 1. Start Postgres
```sh
make up
```
Starts Postgres 16 on `localhost:5432` (db/user/password: `calculator`).

### 2. Run migrations
```sh
make migrate
```
Applies the SQL migrations in `migrations/` via the dockerized `migrate/migrate` image, creating the `history` table.

### 3. Run the backend
```sh
cd backend
cp .env.example .env   # adjust if needed
go run ./cmd/server
```
Serves on `http://localhost:8080` (`GET /healthz` should return `200 OK`).

### 4. Run the frontend
```sh
cd frontend
npm install
npm run dev
```
Serves the Vite dev server (default `http://localhost:5173`) and proxies `/api/*` to `http://localhost:8080`, so the UI talks to your locally running backend with no extra CORS setup.

### 5. Run tests
```sh
make test
```
Runs `go test ./...` in `backend/` and `npm test` in `frontend/`. Note: `frontend/package.json` currently has no `test` script, so the frontend half of `make test` will fail until one is added (or the Makefile is pointed at `npx vitest run` directly) — run `npx vitest run` in `frontend/` directly in the meantime.

## Status

Backend, database, and frontend are all functional end to end: the calculator evaluates expressions (`+ - * /`, parentheses, decimals) and persists each result to Postgres, capped at the last 100 rows; the "Grey Matter" frontend is a physical-calculator-styled SPA wired to both API endpoints, with a scrollable history overlay for re-using past calculations. Remaining gaps: no automated frontend tests yet, and the Navbar's Log In/Sign Up buttons are placeholders (no auth backend exists).
