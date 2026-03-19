# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

FutureProof is a Global Career Intelligence Platform — a Go backend + React/TypeScript frontend — that analyzes AI-driven career disruption risk, renders an interactive 3D globe, and delivers personalized career pivot recommendations via Claude API.

## Commands

### Backend (Go)

```bash
cd backend
go mod download                            # Install dependencies
go run cmd/server/main.go                 # Run dev server (port 8080)
go build -o bin/server cmd/server/main.go # Build binary
go test ./...                             # Run all tests
go test ./internal/services/...           # Run a specific package's tests
air                                       # Hot reload (requires: go install github.com/cosmtrek/air@latest)
```

### Frontend (React/TypeScript)

```bash
cd frontend
npm install
npm run dev     # Dev server (port 5173)
npm run build
npm run lint    # TypeScript type checking
```

### Infrastructure

```bash
docker-compose up -d postgres redis   # Start only DB + cache
docker-compose up -d                  # Start all services (postgres, redis, api)
```

Migrations run automatically when PostgreSQL starts via Docker Compose.

## Architecture

```
futureproof/
├── backend/
│   ├── cmd/server/main.go     # Entry point: wires everything together
│   ├── internal/
│   │   ├── config/            # Reads env vars; ANTHROPIC_API_KEY is required
│   │   ├── models/            # Domain structs (City, Profession, User, AIAnalysis)
│   │   ├── repository/        # PostgreSQL queries via pgx/v5
│   │   ├── services/          # Business logic; AI analysis currently mocked
│   │   ├── api/               # Fiber routes, handlers, middleware
│   │   └── workers/           # Background jobs (stubs, not implemented)
│   ├── pkg/
│   │   ├── database/          # pgxpool connection (min 5, max 25)
│   │   ├── cache/             # Redis client wrapper
│   │   └── logger/            # Zap structured logger
│   └── migrations/            # 4 SQL files (schema + seed)
└── frontend/
    └── src/
        ├── components/        # React + Three.js 3D globe
        ├── hooks/             # Custom hooks
        ├── store/             # Zustand state management
        └── api/               # API client
```

### Request Flow

```
HTTP/WS Request
  → Fiber Middleware (CORS, rate-limit via Redis, logger, JWT auth)
  → Handler (internal/api/handlers/)
  → Service (internal/services/)          ← business logic + Redis caching
  → Repository (internal/repository/)     ← PostgreSQL queries
```

### Key Design Decisions

- **Layered architecture**: handlers → services → repositories. Services own caching logic.
- **Redis cache TTLs**: city/profession/career-transition data caches for 1 hour; AI analysis results cache for 24 hours (keyed by request hash).
- **Authentication**: JWT (HMAC-SHA256), 24h access tokens, 7-day refresh tokens, bcrypt password hashing.
- **Web framework**: [Fiber v2](https://github.com/gofiber/fiber) (not `net/http`).
- **Database driver**: `pgx/v5` directly (no ORM).

### API Routes

| Method | Path | Auth |
|--------|------|------|
| GET | `/api/v1/health` | — |
| GET | `/api/v1/cities`, `/cities/:id`, `/cities/region/:region` | — |
| GET | `/api/v1/professions`, `/professions/future`, `/professions/:slug`, `/professions/:slug/pivots` | — |
| POST | `/api/v1/analyze` | — |
| POST | `/api/v1/ai/chat` | — |
| WS | `/ws/globe` | — |
| POST | `/api/v1/auth/register`, `/auth/login` | — |
| GET/PUT | `/api/v1/users/me`, `/users/profile` | JWT |
| POST/GET | `/api/v1/careers/save`, `/careers/saved` | JWT |

## Environment Variables

Copy `.env.example` to `.env`. Required:

- `ANTHROPIC_API_KEY` — must be set; service panics without it
- `ANTHROPIC_MODEL` — defaults to `claude-3-5-sonnet-20241022`
- `JWT_SECRET` — change from default in production

## Implementation Status

- **AI Service** (`internal/services/ai_service.go`): currently returns mock responses. Real Anthropic Claude API integration is TODO.
- **Workers** (`internal/workers/`): job scraper and data fetcher are stubs, not implemented.
- **Seed data** (`migrations/004_seed_data.sql`): schema exists but actual data rows need to be populated.
