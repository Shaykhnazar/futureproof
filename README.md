# FutureProof

**Global Career Intelligence Platform for the AI Era**

FutureProof helps professionals navigate AI-driven career disruption through interactive 3D visualization, AI-powered risk analysis, and real-time job market intelligence.

## Features

- **Interactive 3D Globe**: Visualize global opportunity scores across 20+ major cities
- **AI Risk Analysis**: Claude-powered career automation risk assessment
- **Career Pivot Recommendations**: Smart suggestions for future-proof transitions
- **Real-time Job Market Data**: Live updates via WebSockets
- **Future Jobs Database**: Emerging professions (2025-2035)
- **Personalized Career Roadmaps**: Tailored skill development paths

## Tech Stack

### Backend
- **Go 1.22** with Fiber v2 web framework
- **PostgreSQL 16** for primary data storage
- **Redis 7** for caching and sessions
- **Anthropic Claude API** for AI analysis
- **WebSockets** for real-time updates

### Frontend
- **React 18** with Vite build tool
- **Three.js** + @react-three/fiber for 3D globe
- **Tailwind CSS** + shadcn/ui components
- **Zustand** for state management
- **TypeScript** for type safety

## Quick Start

### Prerequisites
- Go 1.22+
- Node.js 20+
- Docker & Docker Compose
- Anthropic API key

### 1. Clone and Setup

```bash
git clone <repository-url>
cd futureproof
cp .env.example .env
```

### 2. Add API Keys to .env

```bash
# Required
ANTHROPIC_API_KEY=sk-ant-api03-your-key-here

# Optional (for enhanced features)
ADZUNA_APP_ID=your-adzuna-app-id
ADZUNA_API_KEY=your-adzuna-api-key
```

### 3. Start Infrastructure

```bash
docker-compose up -d postgres redis
```

Wait for services to be healthy:
```bash
docker-compose ps
```

### 4. Run Backend

```bash
cd backend
go mod download
go run cmd/server/main.go
```

Backend will be available at `http://localhost:8080`

### 5. Run Frontend

```bash
cd frontend
npm install
npm run dev
```

Frontend will be available at `http://localhost:5173`

## Project Structure

```
futureproof/
├── backend/
│   ├── cmd/server/           # Application entry point
│   ├── internal/
│   │   ├── api/              # HTTP handlers, middleware, router
│   │   ├── services/         # Business logic layer
│   │   ├── repository/       # Database access layer
│   │   ├── models/           # Domain models
│   │   ├── workers/          # Background jobs
│   │   └── config/           # Configuration management
│   ├── pkg/                  # Reusable packages
│   ├── migrations/           # SQL migration files
│   └── go.mod
├── frontend/
│   ├── src/
│   │   ├── components/       # React components
│   │   ├── hooks/            # Custom React hooks
│   │   ├── store/            # Zustand stores
│   │   ├── api/              # API client
│   │   └── types/            # TypeScript types
│   └── package.json
├── docker-compose.yml
└── README.md
```

## API Endpoints

### Cities
- `GET /api/v1/cities` - List all cities with opportunity scores
- `GET /api/v1/cities/:id` - Get city details
- `GET /api/v1/cities/:id/jobs` - Get jobs for a city

### Careers
- `GET /api/v1/professions` - List all professions
- `POST /api/v1/analyze` - Analyze career automation risk
- `GET /api/v1/professions/:id/pivots` - Get pivot recommendations

### AI
- `POST /api/v1/ai/chat` - Chat with AI career coach (SSE stream)

### Users (Protected)
- `POST /api/v1/auth/register` - Create account
- `POST /api/v1/auth/login` - Login
- `GET /api/v1/users/me` - Get user profile
- `POST /api/v1/users/careers/save` - Save career analysis

### WebSocket
- `WS /ws/globe` - Real-time globe updates

## Database Migrations

Migrations are automatically run on container startup. To run manually:

```bash
cd backend
psql -h localhost -U futureproof -d futureproof_db -f migrations/001_create_cities.sql
psql -h localhost -U futureproof -d futureproof_db -f migrations/002_create_professions.sql
psql -h localhost -U futureproof -d futureproof_db -f migrations/003_create_users.sql
psql -h localhost -U futureproof -d futureproof_db -f migrations/004_seed_data.sql
```

## Development

### Run Tests
```bash
# Backend
cd backend
go test ./...

# Frontend
cd frontend
npm run test
```

### Build for Production
```bash
# Backend
cd backend
go build -o bin/server cmd/server/main.go

# Frontend
cd frontend
npm run build
```

## Deployment

See [docs/deployment.md](docs/deployment.md) for detailed deployment instructions for:
- Railway (backend)
- Cloudflare Pages (frontend)
- Fly.io alternative

## Architecture

```
Client (React + Three.js)
         ↓
   API Gateway (Fiber)
         ↓
    ┌────┴────┐
Services    Workers
    ├────────┤
Repositories
    ├────────┤
PostgreSQL  Redis
```

## License

MIT License - see LICENSE file for details

## Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

## Support

For issues and questions:
- GitHub Issues: [repository-url]/issues
- Documentation: [docs/](docs/)
