# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

See README.md for project overview, setup steps, and API endpoint listing.

## Development Commands

```bash
# Backend (run from backend/)
go mod tidy                              # Sync dependencies and go.sum
go run cmd/server/main.go               # Dev server (port 8080)
go test ./...                           # All tests
go test ./internal/services/...         # Single package tests
air                                     # Hot reload (go install github.com/cosmtrek/air@latest)

# Frontend (run from frontend/)
npm run lint    # TypeScript type checking
```

## Gotchas

- `.env` must be at `backend/.env` (not project root) — the server loads it relative to its working directory. Copy with `cp .env.example backend/.env`.
- `ANTHROPIC_API_KEY` is required; the server exits immediately if missing.
- After cloning or adding dependencies, run `go mod tidy` (not just `go mod download`) to keep `go.sum` in sync.

## Architecture

Layered: **handlers → services → repositories**. Services own all caching logic.

```
backend/
├── cmd/server/main.go       # Entry point: wires config, DB, Redis, services, Fiber
├── internal/
│   ├── config/              # Env var loading; all tunables in one struct
│   ├── models/              # Domain structs (City, Profession, User, AIAnalysis)
│   ├── repository/          # Raw pgx/v5 queries — no ORM
│   ├── services/            # Business logic + Redis caching
│   ├── api/
│   │   ├── router.go        # Route registration
│   │   ├── handlers/        # One file per domain (cities, careers, ai, users, websocket)
│   │   └── middleware/      # JWT auth, CORS, Zap request logger, Redis rate limiter
│   └── workers/             # Background jobs (stubs only, not implemented)
└── pkg/
    ├── database/            # pgxpool (min 5 / max 25 connections)
    ├── cache/               # Redis client wrapper
    └── logger/              # Zap initialization
```

**Redis cache TTLs**: 1 hour for city/profession/career-transition data; 24 hours for AI analysis results (keyed by request hash).

**Auth**: JWT HMAC-SHA256, 24h access / 7-day refresh tokens, bcrypt passwords.

## Implementation Status

- **AI Service** (`internal/services/ai_service.go`): returns mock data — real Anthropic API call is TODO.
- **Workers** (`internal/workers/`): job scraper and data fetcher are empty stubs.
- **Seed data** (`migrations/004_seed_data.sql`): tables exist but data rows need populating.
