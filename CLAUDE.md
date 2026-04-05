# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Go REST API backend for managing user subscriptions and payments. Integrates with CloudPayments, Prodamus (payment providers), UniSender (email), and Telegram (admin notifications).

## Common Commands

```bash
make build              # Build binary to ./build/main
make run                # Run with dev env (loads dev/.dev.env)
make run-race           # Run with Go race detector
make install-deps       # go mod tidy
make infra-up           # Start local PostgreSQL via Docker Compose
make infra-down         # Stop local PostgreSQL

go test ./...           # Run all tests
go test ./... -race     # Run tests with race detector
go test -v -run TestName ./internal/pkg/errs/  # Run a single test
```

## Architecture

**Layered structure under `internal/`:**

- `core/domain/` — User entity and business rules (tariffs, subscription statuses)
- `core/payments/` — Payment service: orchestrates user creation, subscription updates, email scheduling, admin notifications
- `core/notify/` — Notification service interface (implemented by Telegram bot or no-op fallback)
- `inra/httpapi/` — Chi router handlers, request/response parsing. Routes registered via `RegIn*` functions
- `inra/pg/` — PostgreSQL repository (sqlx). Uses advisory locks (`pg_advisory_xact_lock`) to prevent race conditions on find-or-create/update operations
- `inra/tgbot/` — Telegram bot (telebot.v4) for admin notifications, restricted to a single admin chat ID
- `inra/emails/` — UniSender email integration for post-payment/post-trial email scheduling
- `inra/cloudpayments/` — CloudPayments REST API client (subscription management)
- `pkg/` — Shared utilities: database connection (`db/`), structured logging (`logger/`), HTTP server with ACME/autocert (`httpx/`), env config loading (`envconf/`)

**Bootstrap flow:** `cmd/api/main.go` → `boot.StartApp()` → `Server.Init()` (config, DB, bot, services) → `Server.Run()` (HTTP + bot via `oklog/run.Group`) → `Server.Shutdown()`

**Payment webhook routes** are secured by URL path hashes (configured via env vars like `HTTP_TRIAL_PAYMENTS_ROUTE_HASH`), not standard auth middleware.

## Key Conventions

- Go 1.24, modules in `go.mod`
- PostgreSQL 16 with `lumos` schema; schema defined in `dev/db/10.init.sql`
- Config loaded from env vars via `caarlos0/env/v11`; dotenv support toggled by `DOTENV_ENABLED`
- HTTP router: `go-chi/chi/v5` with CORS, recoverer, 60s timeout middleware
- Database driver: `pgx/v5` + `jmoiron/sqlx`
- Logging: `log/slog` with `lmittmann/tint` for pretty dev output
- Tests use `stretchr/testify`
- Infrastructure package is named `inra` (not `infra`)
