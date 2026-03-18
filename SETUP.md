# FutureProof - Setup Guide

Complete setup instructions for the FutureProof career intelligence platform.

## Prerequisites

- **Go 1.22+** - [Download](https://golang.org/dl/)
- **Node.js 20+** - [Download](https://nodejs.org/)
- **Docker & Docker Compose** - [Download](https://www.docker.com/products/docker-desktop/)
- **Anthropic API Key** - [Get one here](https://console.anthropic.com/)

## Quick Start

### 1. Clone the Repository

```bash
cd d:\Projects\Lessons\Go\futureproof
```

### 2. Set Up Environment Variables

```bash
cp .env.example .env
```

Edit `.env` and add your API keys:

```env
# Required
ANTHROPIC_API_KEY=sk-ant-api03-your-key-here

# Optional (can add later)
ADZUNA_APP_ID=your-adzuna-app-id
ADZUNA_API_KEY=your-adzuna-api-key
```

### 3. Start Infrastructure (PostgreSQL + Redis)

```bash
docker-compose up -d postgres redis
```

Wait for services to be ready (about 10-15 seconds):

```bash
docker-compose ps
```

Both services should show "healthy" status.

### 4. Run Backend

```bash
cd backend
go mod download
go run cmd/server/main.go
```

You should see:
```
Server starting on :8080
```

The backend will automatically:
- Connect to PostgreSQL
- Run database migrations
- Seed initial data (20 cities, 28 professions)
- Connect to Redis

### 5. Run Frontend (New Terminal)

```bash
cd frontend
npm install
npm run dev
```

You should see:
```
VITE v5.x.x ready in X ms
Local: http://localhost:5173/
```

### 6. Open the Application

Navigate to [http://localhost:5173](http://localhost:5173)

You should see:
- Interactive 3D globe with 20 cities marked
- Real-time WebSocket connection
- Career Analyzer tool

## API Endpoints

### Public Endpoints

```bash
# Health check
curl http://localhost:8080/api/v1/health

# Get all cities
curl http://localhost:8080/api/v1/cities

# Get all professions
curl http://localhost:8080/api/v1/professions

# Get future professions
curl http://localhost:8080/api/v1/professions/future

# Analyze a career
curl -X POST http://localhost:8080/api/v1/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "profession_slug": "software-engineer",
    "location": "San Francisco",
    "years_exp": 5,
    "current_skills": ["Python", "React", "Docker"]
  }'
```

### Protected Endpoints (Require Authentication)

```bash
# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "name": "Test User",
    "password": "password123"
  }'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'

# Get current user (with token)
curl http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## Database Access

### PostgreSQL

```bash
# Connect to database
docker exec -it futureproof_postgres psql -U futureproof -d futureproof_db

# View tables
\dt

# View cities
SELECT name, country, region FROM cities;

# View professions
SELECT title, category, ai_risk_score FROM professions;

# Exit
\q
```

### Redis

```bash
# Connect to Redis
docker exec -it futureproof_redis redis-cli

# View all keys
KEYS *

# Get a value
GET "cities:all"

# Exit
exit
```

## Testing WebSocket

```javascript
// In browser console
const ws = new WebSocket('ws://localhost:8080/ws/globe');

ws.onopen = () => console.log('Connected');
ws.onmessage = (e) => console.log('Message:', JSON.parse(e.data));

// Trigger an update
fetch('http://localhost:8080/api/v1/trigger-update', { method: 'POST' });
```

## Development Workflow

### Backend Development

```bash
cd backend

# Run with hot reload (install air first: go install github.com/cosmtrek/air@latest)
air

# Run tests
go test ./...

# Build binary
go build -o bin/server cmd/server/main.go

# Run binary
./bin/server
```

### Frontend Development

```bash
cd frontend

# Start dev server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview

# Type check
npm run lint
```

## Troubleshooting

### PostgreSQL connection fails

```bash
# Check if container is running
docker ps | grep postgres

# View logs
docker logs futureproof_postgres

# Restart
docker-compose restart postgres
```

### Redis connection fails

```bash
# Check if container is running
docker ps | grep redis

# View logs
docker logs futureproof_redis

# Restart
docker-compose restart redis
```

### Port already in use

```bash
# Find process using port 8080
# Windows
netstat -ano | findstr :8080

# Linux/Mac
lsof -i :8080

# Kill the process or change PORT in .env
```

### Frontend can't connect to backend

1. Check backend is running on port 8080
2. Check CORS settings in `.env` include `http://localhost:5173`
3. Check browser console for errors
4. Try clearing browser cache

## Production Deployment

### Backend (Railway/Fly.io)

```bash
# Build Docker image
docker build -t futureproof-api ./backend

# Tag for registry
docker tag futureproof-api registry.example.com/futureproof-api

# Push
docker push registry.example.com/futureproof-api
```

### Frontend (Cloudflare Pages)

```bash
cd frontend

# Build
npm run build

# Deploy to Cloudflare Pages
# Connect your GitHub repo to Cloudflare Pages
# Build command: npm run build
# Build output directory: dist
```

## Project Structure

```
futureproof/
├── backend/
│   ├── cmd/server/           # Application entry point
│   ├── internal/
│   │   ├── api/              # HTTP handlers & routes
│   │   ├── services/         # Business logic
│   │   ├── repository/       # Database access
│   │   ├── models/           # Data models
│   │   ├── workers/          # Background jobs
│   │   └── config/           # Configuration
│   ├── pkg/                  # Shared packages
│   ├── migrations/           # SQL migrations
│   ├── go.mod
│   └── Dockerfile
├── frontend/
│   ├── src/
│   │   ├── components/       # React components
│   │   ├── hooks/            # Custom hooks
│   │   ├── store/            # Zustand stores
│   │   ├── api/              # API client
│   │   ├── types/            # TypeScript types
│   │   ├── App.tsx
│   │   └── main.tsx
│   ├── package.json
│   └── vite.config.ts
├── docker-compose.yml
├── .env.example
└── README.md
```

## Next Steps

1. **Customize the data**: Edit `backend/migrations/004_seed_data.sql`
2. **Add real AI integration**: Update `backend/internal/services/ai_service.go`
3. **Implement Adzuna API**: Update `backend/internal/workers/job_scraper.go`
4. **Add authentication UI**: Create login/register pages
5. **Enhance the globe**: Add more interactive features
6. **Deploy to production**: Follow deployment guides above

## Support

For issues and questions:
- Check the main [README.md](README.md)
- Review the [documentation](docs/)
- Open an issue on GitHub

## License

MIT License - see LICENSE file
