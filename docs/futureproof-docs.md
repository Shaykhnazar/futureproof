# 🌍 FutureProof — Platform Architecture & Implementation Guide

> Global Career Intelligence Platform — Post-AI World Edition

---

## 📁 Document Index

| File | Purpose |
|------|---------|
| `01-stack-decisions.md` | Tech stack rationale & tradeoffs |
| `02-architecture.md` | System architecture & data flow |
| `03-backend-go.md` | Go backend structure & API design |
| `04-frontend-react.md` | React + Three.js frontend guide |
| `05-database.md` | Database schema & data strategy |
| `06-ai-integration.md` | AI/LLM layer design |
| `07-data-sources.md` | Real-world data pipelines |
| `08-deployment.md` | Infrastructure & DevOps |
| `09-roadmap.md` | Phased feature roadmap |

---

---

# 01 — Stack Decisions & Rationale

## My Recommendation: Stick with Go + React ✅

Go is an excellent choice for this platform. Here's a detailed comparison so you can decide with confidence.

---

## Backend: Go vs Alternatives

### ✅ Go (Recommended)

```
Strengths:
- Extremely fast HTTP servers (net/http, Fiber, Echo)
- Native concurrency with goroutines → ideal for scraping + AI calls in parallel
- Single binary deployment → simple Docker images
- Strong typing catches bugs at compile time
- Excellent for WebSocket support (real-time globe updates)
- Low memory footprint vs Node/Python

Weaknesses:
- More verbose than Python for ML tasks
- Smaller AI/data ecosystem vs Python
- No built-in ORM as mature as Django/SQLAlchemy
```

**Verdict for FutureProof:** Go handles the API layer, data aggregation, caching, and WebSockets perfectly. AI/ML tasks should be delegated to external services (OpenAI/Anthropic API), so Go's lack of native ML is irrelevant.

---

### Python FastAPI (Strong Alternative)

```
Strengths:
- Best AI/ML ecosystem (scikit-learn, pandas, transformers)
- Fastest prototyping for data pipelines
- Native integration with LangChain, LlamaIndex
- Async support with uvicorn

Weaknesses:
- 3-5x slower than Go under load
- GIL limits true parallelism
- Higher memory usage
- Dependency hell (pip/conda)
```

**Verdict:** Use Python ONLY if you plan to run local ML models. For API-first AI (calling Anthropic/OpenAI), Go is superior.

---

### Node.js (Not Recommended for this use case)

```
Weaknesses for this platform:
- Single-threaded (though async)
- Weaker typing even with TypeScript
- Not ideal for CPU-bound data aggregation
- npm ecosystem instability
```

---

## 🏆 Final Stack Decision

| Layer | Technology | Reason |
|-------|-----------|--------|
| **Backend API** | Go + Fiber v2 | Fast, concurrent, typed |
| **AI Layer** | Anthropic Claude API | Career analysis, pivots |
| **Data Scraping** | Go goroutines + colly | Parallel job market crawlers |
| **Frontend** | React 18 + Vite | Fast dev experience |
| **3D Globe** | Three.js | Already proven, great ecosystem |
| **UI Components** | Tailwind CSS + shadcn/ui | Consistent design system |
| **Primary DB** | PostgreSQL | Relational career/city data |
| **Cache** | Redis | Globe data, AI responses, sessions |
| **Search** | Elasticsearch | Full-text job/skill search |
| **Auth** | JWT + OAuth2 (Google/GitHub) | Standard, easy to implement |
| **Real-time** | WebSockets (Go + gorilla/websocket) | Live globe updates |
| **Deployment** | Docker + Railway / Fly.io | Simple, affordable |
| **CDN** | Cloudflare | Globe assets, static files |

---

---

# 02 — System Architecture

## High-Level Overview

```
┌─────────────────────────────────────────────────────────┐
│                     CLIENT BROWSER                       │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌────────┐  │
│  │ 3D Globe │  │ Career   │  │ Future   │  │ User   │  │
│  │ Three.js │  │ Analyzer │  │ Jobs     │  │ Profile│  │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └───┬────┘  │
│       └──────────────┴──────────────┴────────────┘       │
│                    React 18 + Vite                        │
└─────────────────────────┬───────────────────────────────┘
                          │ HTTPS + WSS
┌─────────────────────────▼───────────────────────────────┐
│                   Go API SERVER (Fiber v2)                │
│  ┌─────────────┐  ┌──────────────┐  ┌────────────────┐  │
│  │ REST API    │  │ WebSocket    │  │ AI Proxy       │  │
│  │ /api/v1/    │  │ /ws/globe    │  │ /api/v1/ai/    │  │
│  └──────┬──────┘  └──────┬───────┘  └───────┬────────┘  │
│         │                │                   │            │
│  ┌──────▼──────┐  ┌──────▼──────┐  ┌────────▼──────┐   │
│  │ Controllers │  │ Hub Manager │  │ AI Service    │   │
│  │ Services    │  │ (Broadcast) │  │ (Claude API)  │   │
│  │ Repositories│  └─────────────┘  └───────────────┘   │
│  └──────┬──────┘                                         │
└─────────┼───────────────────────────────────────────────┘
          │
┌─────────▼──────────────────────────────────────────────┐
│                    DATA LAYER                            │
│  ┌─────────────┐  ┌──────────┐  ┌────────────────────┐ │
│  │ PostgreSQL  │  │  Redis   │  │  Elasticsearch     │ │
│  │ (Primary)   │  │  (Cache) │  │  (Job Search)      │ │
│  └─────────────┘  └──────────┘  └────────────────────┘ │
└────────────────────────────────────────────────────────┘
          │
┌─────────▼──────────────────────────────────────────────┐
│               BACKGROUND WORKERS (Go goroutines)         │
│  ┌─────────────┐  ┌──────────────┐  ┌───────────────┐  │
│  │ Job Market  │  │ LinkedIn/    │  │ World Bank    │  │
│  │ Scraper     │  │ Indeed API   │  │ Data Fetcher  │  │
│  └─────────────┘  └──────────────┘  └───────────────┘  │
└────────────────────────────────────────────────────────┘
```

## Data Flow: Career Analysis Request

```
User selects "Software Engineer"
        │
        ▼
React sends POST /api/v1/careers/analyze
        │
        ▼
Go Controller → Career Service
        │
        ├── PostgreSQL: fetch profession baseline data
        ├── Redis: check cached AI response (TTL: 24h)
        │         └── Cache HIT? Return immediately
        │         └── Cache MISS? ↓
        ├── Anthropic API: generate personalized analysis
        ├── Redis: cache result
        └── Return JSON to client
```

---

---

# 03 — Backend: Go (Fiber v2)

## Project Structure

```
futureproof-api/
├── cmd/
│   └── server/
│       └── main.go               # Entry point
├── internal/
│   ├── api/
│   │   ├── router.go             # Route definitions
│   │   ├── middleware/
│   │   │   ├── auth.go           # JWT middleware
│   │   │   ├── cors.go
│   │   │   ├── ratelimit.go      # Redis-backed rate limiter
│   │   │   └── logger.go
│   │   └── handlers/
│   │       ├── careers.go
│   │       ├── cities.go
│   │       ├── ai.go
│   │       ├── users.go
│   │       └── websocket.go
│   ├── services/
│   │   ├── career_service.go
│   │   ├── city_service.go
│   │   ├── ai_service.go         # Anthropic/OpenAI client
│   │   ├── scraper_service.go
│   │   └── auth_service.go
│   ├── repository/
│   │   ├── career_repo.go        # DB queries
│   │   ├── city_repo.go
│   │   └── user_repo.go
│   ├── models/
│   │   ├── career.go
│   │   ├── city.go
│   │   ├── user.go
│   │   └── analysis.go
│   ├── workers/
│   │   ├── job_scraper.go        # Background goroutine
│   │   ├── data_fetcher.go
│   │   └── scheduler.go
│   └── config/
│       └── config.go             # Env vars, secrets
├── pkg/
│   ├── cache/
│   │   └── redis.go
│   ├── database/
│   │   └── postgres.go
│   └── logger/
│       └── logger.go
├── migrations/
│   ├── 001_create_cities.sql
│   ├── 002_create_professions.sql
│   └── 003_create_users.sql
├── scripts/
│   └── seed.go                   # Seed initial data
├── .env.example
├── Dockerfile
├── docker-compose.yml
└── go.mod
```

---

## go.mod Dependencies

```go
module github.com/yourname/futureproof-api

go 1.22

require (
    github.com/gofiber/fiber/v2 v2.52.0
    github.com/gofiber/websocket/v2 v2.2.1
    github.com/golang-jwt/jwt/v5 v5.2.0
    github.com/jackc/pgx/v5 v5.5.0         // PostgreSQL driver
    github.com/redis/go-redis/v9 v9.4.0
    github.com/elastic/go-elasticsearch/v8 v8.11.0
    github.com/gocolly/colly/v2 v2.1.0     // Web scraper
    github.com/robfig/cron/v3 v3.0.1       // Task scheduler
    github.com/joho/godotenv v1.5.1
    go.uber.org/zap v1.27.0                // Structured logging
    github.com/anthropics/anthropic-sdk-go v0.2.0
)
```

---

## REST API Endpoints

```
AUTH
  POST   /api/v1/auth/register
  POST   /api/v1/auth/login
  POST   /api/v1/auth/refresh
  GET    /api/v1/auth/google/callback

CITIES (Globe Data)
  GET    /api/v1/cities                    # All cities with scores
  GET    /api/v1/cities/:id                # City detail
  GET    /api/v1/cities/:id/professions    # Top professions for city
  GET    /api/v1/cities/top?limit=10       # Top opportunity hubs
  GET    /api/v1/cities/search?q=berlin    # Search cities

CAREERS
  GET    /api/v1/careers                   # All tracked professions
  GET    /api/v1/careers/:slug             # Profession detail + AI risk
  POST   /api/v1/careers/analyze           # Analyze risk + get pivots
  GET    /api/v1/careers/future            # Top future jobs list
  GET    /api/v1/careers/transitions/:from # Career transition paths

AI (Claude-powered)
  POST   /api/v1/ai/career-roadmap         # Personalized roadmap
  POST   /api/v1/ai/skill-gap              # Skill gap analysis
  POST   /api/v1/ai/relocation-advice      # Where to move for career
  POST   /api/v1/ai/chat                   # Free-form career coach chat

USERS (Authenticated)
  GET    /api/v1/users/me
  PUT    /api/v1/users/me
  GET    /api/v1/users/me/saved-careers
  POST   /api/v1/users/me/saved-careers
  GET    /api/v1/users/me/roadmap

REAL-TIME
  WS     /ws/globe                         # Live opportunity score updates
```

---

## Key Go Code Snippets

### main.go

```go
package main

import (
    "log"
    "github.com/gofiber/fiber/v2"
    "github.com/yourname/futureproof-api/internal/api"
    "github.com/yourname/futureproof-api/internal/config"
    "github.com/yourname/futureproof-api/pkg/database"
    "github.com/yourname/futureproof-api/pkg/cache"
)

func main() {
    cfg := config.Load()

    db := database.NewPostgres(cfg.DatabaseURL)
    defer db.Close()

    rdb := cache.NewRedis(cfg.RedisURL)
    defer rdb.Close()

    app := fiber.New(fiber.Config{
        AppName: "FutureProof API v1",
        ErrorHandler: api.ErrorHandler,
    })

    api.SetupRoutes(app, db, rdb, cfg)

    log.Fatal(app.Listen(":" + cfg.Port))
}
```

### Career Analysis Handler

```go
// internal/api/handlers/careers.go
package handlers

type AnalyzeRequest struct {
    Profession string `json:"profession" validate:"required"`
    Location   string `json:"location"`
    YearsExp   int    `json:"years_experience"`
    Skills     []string `json:"skills"`
}

func (h *CareerHandler) Analyze(c *fiber.Ctx) error {
    var req AnalyzeRequest
    if err := c.BodyParser(&req); err != nil {
        return fiber.NewError(400, "invalid request body")
    }

    // Check Redis cache first
    cacheKey := fmt.Sprintf("analysis:%s:%s", req.Profession, req.Location)
    if cached, err := h.cache.Get(ctx, cacheKey).Result(); err == nil {
        return c.JSON(json.RawMessage(cached))
    }

    // Fetch base data from DB
    profession, err := h.careerSvc.GetProfession(req.Profession)
    if err != nil {
        return fiber.NewError(404, "profession not found")
    }

    // Call AI service for personalized analysis
    analysis, err := h.aiSvc.AnalyzeCareer(req, profession)
    if err != nil {
        return fiber.NewError(500, "AI analysis failed")
    }

    // Cache for 24 hours
    h.cache.Set(ctx, cacheKey, analysis, 24*time.Hour)

    return c.JSON(analysis)
}
```

### AI Service (Anthropic)

```go
// internal/services/ai_service.go
package services

import (
    "github.com/anthropics/anthropic-sdk-go"
)

type AIService struct {
    client *anthropic.Client
}

func (s *AIService) AnalyzeCareer(req AnalyzeRequest, prof Profession) (*CareerAnalysis, error) {
    prompt := fmt.Sprintf(`
You are a world-class career strategist in 2025.

Profession: %s
Current Skills: %v
Years of Experience: %d
Location: %s
Industry: %s

Provide a JSON response with:
1. ai_replacement_risk: 0-100 score with detailed reasoning
2. resilient_skills: which of their skills AI cannot replace
3. at_risk_tasks: specific tasks likely to be automated
4. top_5_pivots: best career transitions with match_score and skill_gap
5. timeline: when to expect major disruption
6. immediate_actions: 3 things to do in the next 6 months

Be specific, data-driven, and brutally honest.
`, req.Profession, req.Skills, req.YearsExp, req.Location, prof.Industry)

    msg, err := s.client.Messages.New(context.TODO(), anthropic.MessageNewParams{
        Model:     anthropic.F(anthropic.ModelClaude3_5SonnetLatest),
        MaxTokens: anthropic.F(int64(2048)),
        Messages: anthropic.F([]anthropic.MessageParam{
            anthropic.UserMessageParam(
                anthropic.NewTextBlock(prompt),
            ),
        }),
    })
    if err != nil {
        return nil, err
    }

    // Parse JSON from response
    var analysis CareerAnalysis
    json.Unmarshal([]byte(msg.Content[0].Text), &analysis)
    return &analysis, nil
}
```

### WebSocket Globe Hub

```go
// internal/api/handlers/websocket.go
package handlers

type GlobeHub struct {
    clients   map[*websocket.Conn]bool
    broadcast chan GlobeUpdate
    mu        sync.RWMutex
}

type GlobeUpdate struct {
    Type    string      `json:"type"`
    Payload interface{} `json:"payload"`
}

func (h *GlobeHub) Run() {
    ticker := time.NewTicker(30 * time.Second)
    for {
        select {
        case <-ticker.C:
            // Fetch fresh city scores and broadcast
            update := GlobeUpdate{
                Type:    "scores_update",
                Payload: h.cityService.GetAllScores(),
            }
            h.broadcastToAll(update)
        case msg := <-h.broadcast:
            h.broadcastToAll(msg)
        }
    }
}
```

---

---

# 04 — Frontend: React + Three.js

## Project Structure

```
futureproof-web/
├── public/
│   └── textures/
│       ├── earth_day.jpg         # 4K Earth texture
│       ├── earth_night.jpg       # City lights texture
│       ├── earth_clouds.jpg      # Cloud overlay
│       └── earth_normal.jpg      # Bump map
├── src/
│   ├── main.tsx
│   ├── App.tsx
│   ├── components/
│   │   ├── globe/
│   │   │   ├── Globe.tsx         # Three.js scene
│   │   │   ├── GlobeMarkers.tsx  # City dots + popups
│   │   │   ├── GlobeArcs.tsx     # Migration arcs
│   │   │   └── useGlobe.ts       # Globe logic hook
│   │   ├── career/
│   │   │   ├── CareerAnalyzer.tsx
│   │   │   ├── RiskMeter.tsx
│   │   │   ├── PivotCard.tsx
│   │   │   └── SkillGapChart.tsx
│   │   ├── city/
│   │   │   ├── CityPanel.tsx
│   │   │   ├── CitySearch.tsx
│   │   │   └── CityCompare.tsx
│   │   ├── ai/
│   │   │   ├── CareerCoach.tsx   # AI chat interface
│   │   │   └── RoadmapView.tsx
│   │   └── ui/
│   │       ├── Button.tsx
│   │       ├── Card.tsx
│   │       ├── Modal.tsx
│   │       └── ProgressBar.tsx
│   ├── hooks/
│   │   ├── useWebSocket.ts       # Live globe updates
│   │   ├── useCareerAnalysis.ts
│   │   └── useAuth.ts
│   ├── store/
│   │   ├── globeStore.ts         # Zustand stores
│   │   ├── careerStore.ts
│   │   └── userStore.ts
│   ├── api/
│   │   ├── client.ts             # Axios instance
│   │   ├── careers.ts
│   │   ├── cities.ts
│   │   └── ai.ts
│   ├── types/
│   │   ├── career.ts
│   │   ├── city.ts
│   │   └── user.ts
│   └── styles/
│       └── globals.css
├── package.json
├── vite.config.ts
├── tailwind.config.ts
└── tsconfig.json
```

---

## package.json Key Dependencies

```json
{
  "dependencies": {
    "react": "^18.3.0",
    "react-dom": "^18.3.0",
    "three": "^0.162.0",
    "@react-three/fiber": "^8.16.0",
    "@react-three/drei": "^9.99.0",
    "zustand": "^4.5.0",
    "axios": "^1.6.0",
    "framer-motion": "^11.0.0",
    "recharts": "^2.12.0",
    "tailwindcss": "^3.4.0",
    "@radix-ui/react-dialog": "^1.0.0",
    "react-query": "^5.0.0",
    "react-router-dom": "^6.22.0"
  }
}
```

> **Upgrade tip:** Use `@react-three/fiber` + `@react-three/drei` instead of raw Three.js. It gives you declarative 3D, built-in helpers like `<OrbitControls>`, `<Stars>`, and `<Html>` popups, significantly reducing boilerplate.

---

## Globe Component with Real Textures

```tsx
// src/components/globe/Globe.tsx
import { Canvas } from "@react-three/fiber"
import { OrbitControls, Stars, useTexture } from "@react-three/drei"
import { Suspense } from "react"
import GlobeMarkers from "./GlobeMarkers"
import GlobeArcs from "./GlobeArcs"

export function Globe({ cities, onCityClick }) {
  return (
    <Canvas camera={{ position: [0, 0, 2.8], fov: 42 }}>
      <ambientLight intensity={0.4} color="#1a3355" />
      <directionalLight position={[5, 3, 4]} intensity={2} color="#4499ff" />
      <Stars radius={100} depth={50} count={6000} factor={4} fade />
      <Suspense fallback={null}>
        <EarthSphere />
        <GlobeMarkers cities={cities} onCityClick={onCityClick} />
        <GlobeArcs />
      </Suspense>
      <OrbitControls
        enableZoom={true}
        minDistance={1.5}
        maxDistance={5}
        autoRotate
        autoRotateSpeed={0.3}
        enablePan={false}
      />
    </Canvas>
  )
}

function EarthSphere() {
  const [dayMap, normalMap, cloudsMap] = useTexture([
    "/textures/earth_day.jpg",
    "/textures/earth_normal.jpg",
    "/textures/earth_clouds.jpg",
  ])
  return (
    <group>
      <mesh>
        <sphereGeometry args={[1, 64, 64]} />
        <meshPhongMaterial map={dayMap} normalMap={normalMap} />
      </mesh>
      {/* Animated cloud layer */}
      <mesh rotation={[0, 0.3, 0]}>
        <sphereGeometry args={[1.003, 32, 32]} />
        <meshPhongMaterial map={cloudsMap} transparent opacity={0.25} />
      </mesh>
      {/* Atmosphere glow */}
      <mesh>
        <sphereGeometry args={[1.08, 32, 32]} />
        <meshPhongMaterial color="#0055bb" transparent opacity={0.06} side={THREE.BackSide} />
      </mesh>
    </group>
  )
}
```

---

---

# 05 — Database Schema

## PostgreSQL Tables

```sql
-- Cities & Regions
CREATE TABLE cities (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(100) NOT NULL,
    country     VARCHAR(100) NOT NULL,
    region      VARCHAR(100),
    lat         DECIMAL(9,6) NOT NULL,
    lng         DECIMAL(9,6) NOT NULL,
    population  INTEGER,
    timezone    VARCHAR(50),
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    updated_at  TIMESTAMPTZ DEFAULT NOW()
);

-- Opportunity Scores (versioned over time)
CREATE TABLE city_scores (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    city_id         UUID REFERENCES cities(id),
    score           INTEGER CHECK (score BETWEEN 0 AND 100),
    job_growth_pct  DECIMAL(5,2),
    remote_score    INTEGER,
    ai_investment   INTEGER,     -- AI startup investment score
    talent_demand   INTEGER,
    cost_of_living  INTEGER,     -- Lower = better for workers
    snapshot_date   DATE NOT NULL,
    source          VARCHAR(50), -- 'manual', 'linkedin_api', 'world_bank'
    created_at      TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_city_scores_city_date ON city_scores(city_id, snapshot_date DESC);

-- Professions
CREATE TABLE professions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug            VARCHAR(100) UNIQUE NOT NULL,
    title           VARCHAR(200) NOT NULL,
    category        VARCHAR(50),
    ai_risk_score   INTEGER CHECK (ai_risk_score BETWEEN 0 AND 100),
    avg_salary_usd  INTEGER,
    description     TEXT,
    is_future_job   BOOLEAN DEFAULT false,
    demand_index    INTEGER,
    growth_pct      DECIMAL(5,2),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

-- Skills
CREATE TABLE skills (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(100) UNIQUE NOT NULL,
    category    VARCHAR(50),     -- 'technical', 'soft', 'domain'
    is_ai_proof BOOLEAN DEFAULT false
);

-- Profession ↔ Skills
CREATE TABLE profession_skills (
    profession_id   UUID REFERENCES professions(id),
    skill_id        UUID REFERENCES skills(id),
    importance      INTEGER CHECK (importance BETWEEN 1 AND 5),
    is_at_risk      BOOLEAN DEFAULT false,
    PRIMARY KEY (profession_id, skill_id)
);

-- Career Transitions
CREATE TABLE career_transitions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    from_profession UUID REFERENCES professions(id),
    to_profession   UUID REFERENCES professions(id),
    match_score     INTEGER,
    transition_reason TEXT,
    avg_reskill_months INTEGER,
    UNIQUE (from_profession, to_profession)
);

-- City Top Jobs
CREATE TABLE city_top_professions (
    city_id         UUID REFERENCES cities(id),
    profession_id   UUID REFERENCES professions(id),
    rank            INTEGER,
    local_demand    INTEGER,
    snapshot_date   DATE,
    PRIMARY KEY (city_id, profession_id, snapshot_date)
);

-- Users
CREATE TABLE users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email           VARCHAR(255) UNIQUE NOT NULL,
    name            VARCHAR(100),
    avatar_url      TEXT,
    auth_provider   VARCHAR(20),   -- 'email', 'google', 'github'
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- User Profiles
CREATE TABLE user_profiles (
    user_id         UUID REFERENCES users(id) PRIMARY KEY,
    current_job_id  UUID REFERENCES professions(id),
    city_id         UUID REFERENCES cities(id),
    years_exp       INTEGER,
    education       VARCHAR(50),
    target_job_id   UUID REFERENCES professions(id),
    skills          TEXT[],         -- Array of skill names
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

-- Saved Careers
CREATE TABLE saved_careers (
    user_id         UUID REFERENCES users(id),
    profession_id   UUID REFERENCES professions(id),
    notes           TEXT,
    saved_at        TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (user_id, profession_id)
);

-- AI Analysis Cache (stored for audit + replay)
CREATE TABLE ai_analyses (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID REFERENCES users(id),
    profession_slug VARCHAR(100),
    location        VARCHAR(100),
    request_hash    VARCHAR(64) UNIQUE,  -- SHA256 of inputs
    result          JSONB NOT NULL,
    model_used      VARCHAR(50),
    tokens_used     INTEGER,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_ai_analyses_hash ON ai_analyses(request_hash);
```

---

## Redis Cache Keys

```
cities:all                      → Full city list (TTL: 1h)
city:{id}                       → City detail (TTL: 30m)
city:{id}:professions           → Top jobs per city (TTL: 1h)
career:{slug}                   → Profession data (TTL: 6h)
analysis:{hash}                 → AI analysis result (TTL: 24h)
user:{id}:session               → JWT session (TTL: 7d)
globe:scores                    → All city scores for globe (TTL: 30m)
leaderboard:cities              → Top 10 cities sorted set
```

---

---

# 06 — AI Integration Layer

## Architecture

The AI layer handles 4 use cases:

```
1. Career Risk Analysis      → One-shot, cached 24h
2. Career Pivot Suggestions  → One-shot, cached 24h
3. Personalized Roadmap      → One-shot per user profile
4. Career Coach Chat         → Multi-turn streaming
```

---

## Career Coach Chat (Streaming)

```go
// Streaming Anthropic response via SSE
func (h *AIHandler) CareerCoachChat(c *fiber.Ctx) error {
    c.Set("Content-Type", "text/event-stream")
    c.Set("Cache-Control", "no-cache")
    c.Set("X-Accel-Buffering", "no")

    var req ChatRequest
    c.BodyParser(&req)

    stream := h.aiSvc.StreamChat(req.Messages, req.UserProfile)

    c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
        for chunk := range stream {
            fmt.Fprintf(w, "data: %s\n\n", chunk)
            w.Flush()
        }
        fmt.Fprintf(w, "data: [DONE]\n\n")
        w.Flush()
    })
    return nil
}
```

## System Prompts

```go
const CAREER_COACH_SYSTEM = `
You are FutureProof AI — the world's most advanced career intelligence system for the post-AI era.

Your role:
- Help people navigate AI-driven career disruption with honesty and empathy
- Give data-backed, actionable advice (not vague platitudes)
- Be direct about automation risks without causing unnecessary panic
- Suggest realistic transition paths based on transferable skills
- Consider the user's location, experience level, and personal context

User Profile: {{.Profile}}
Current Date: {{.Date}}

Always structure advice around: (1) immediate actions, (2) 6-month plan, (3) 2-year vision.
`
```

---

---

# 07 — Real Data Sources

## Free / Affordable APIs

| Source | Data | Cost |
|--------|------|------|
| **World Bank API** | GDP, employment by country | Free |
| **LinkedIn Jobs API** | Job posting volumes | Requires partnership |
| **Bureau of Labor Statistics** | US occupation projections | Free |
| **OECD.Stat** | International labor data | Free |
| **Glassdoor API** | Salary data | Freemium |
| **GitHub Jobs / Adzuna API** | Tech job demand | Free tier |
| **McKinsey / WEF Reports** | AI displacement forecasts | Scraped PDFs |

## Go Scraper Skeleton

```go
// internal/workers/job_scraper.go
func (s *Scraper) ScrapeAdzuna(city string) ([]JobPosting, error) {
    c := colly.NewCollector(
        colly.AllowedDomains("api.adzuna.com"),
    )
    url := fmt.Sprintf(
        "https://api.adzuna.com/v1/api/jobs/%s/search/1?app_id=%s&app_key=%s&where=%s&results_per_page=50",
        s.countryCode, s.appID, s.appKey, url.QueryEscape(city),
    )
    var postings []JobPosting
    c.OnResponse(func(r *colly.Response) {
        json.Unmarshal(r.Body, &postings)
    })
    c.Visit(url)
    return postings, nil
}
```

---

---

# 08 — Deployment

## docker-compose.yml (Local Dev)

```yaml
version: "3.9"
services:
  api:
    build: ./futureproof-api
    ports: ["8080:8080"]
    environment:
      - DATABASE_URL=postgres://fp:fp@postgres:5432/futureproof
      - REDIS_URL=redis://redis:6379
      - ANTHROPIC_API_KEY=${ANTHROPIC_API_KEY}
      - JWT_SECRET=${JWT_SECRET}
    depends_on: [postgres, redis]

  web:
    build: ./futureproof-web
    ports: ["3000:3000"]
    environment:
      - VITE_API_URL=http://localhost:8080

  postgres:
    image: postgres:16
    environment:
      POSTGRES_DB: futureproof
      POSTGRES_USER: fp
      POSTGRES_PASSWORD: fp
    volumes: [pgdata:/var/lib/postgresql/data]

  redis:
    image: redis:7-alpine
    volumes: [redisdata:/data]

volumes:
  pgdata:
  redisdata:
```

## Production: Railway / Fly.io

```
1. Railway (recommended for solo/startup):
   - Connect GitHub repo → auto-deploy
   - Add PostgreSQL + Redis as plugins
   - Set env vars in dashboard
   - ~$20/month for small scale

2. Fly.io (more control, better for scale):
   - fly launch → generates fly.toml
   - fly postgres create
   - fly redis create
   - fly deploy

3. Cloudflare Pages (frontend):
   - Connect GitHub repo
   - Build command: npm run build
   - Output: dist/
   - Free tier is generous
```

---

---

# 09 — Phased Roadmap

## Phase 1 — MVP (6–8 weeks)
- [ ] Go API with cities + professions endpoints
- [ ] PostgreSQL schema + seed data (40 cities, 50 professions)
- [ ] React globe with real Earth textures
- [ ] Career risk analyzer (Claude API)
- [ ] Basic JWT auth
- [ ] Deploy to Railway + Cloudflare Pages

## Phase 2 — Growth (Weeks 9–16)
- [ ] User profiles + saved careers
- [ ] Personalized career roadmaps (AI-generated)
- [ ] Skill gap calculator with learning resources
- [ ] Real job data integration (Adzuna / BLS)
- [ ] City comparison tool on globe
- [ ] WebSocket live updates

## Phase 3 — Scale (Months 5–8)
- [ ] AI Career Coach (streaming chat)
- [ ] Career transition forum / community
- [ ] Email digests: "New opportunities in your city"
- [ ] Mobile app (React Native, same API)
- [ ] Partner API for universities
- [ ] Elasticsearch for advanced job search

## Phase 4 — Monetization
- [ ] Freemium: 3 analyses/month free
- [ ] Pro ($12/mo): unlimited analyses, full roadmap, coach chat
- [ ] University license: white-label for career centers
- [ ] API access for B2B integrations

---

## Estimated Costs at Scale

| Service | Free Tier | At 10k users/mo |
|---------|-----------|----------------|
| Railway (API + DB) | $5/mo | ~$30/mo |
| Cloudflare Pages (Frontend) | Free | Free |
| Anthropic Claude API | Pay-per-use | ~$50–150/mo |
| Redis (Railway plugin) | Included | Included |
| Adzuna Jobs API | 1000 req/mo | ~$50/mo |
| **Total** | | **~$130–230/mo** |

At 500 Pro subscribers ($12/mo) = **$6,000/mo revenue vs ~$230 costs**.

---

*Last updated: March 2026 — Built with ❤️ for the post-AI generation*
