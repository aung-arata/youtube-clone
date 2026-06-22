# Development Guide

## Project Structure

```
youtube-clone/
├── frontend/                   # React + Tailwind CSS
├── services/
│   ├── api-gateway/            # Go — request routing, CORS, rate limiting
│   ├── video-service/          # Go — video CRUD, search, views, likes, analytics
│   ├── user-service/           # Go — user profiles
│   ├── comment-service/        # Go — comment CRUD
│   ├── history-service/        # Go — watch history
│   ├── notification-service/   # Go — notifications, WebSocket
│   └── admin-service/          # PHP (Symfony) — admin dashboard, CMS, reports
├── backend/                    # Deprecated monolith (reference only)
└── docker-compose.yml          # Full microservices stack
```

## Frontend

```bash
cd frontend
npm run dev       # Dev server (http://localhost:5173)
npm run build     # Production build → dist/
npm run preview   # Preview production build
npm run test      # Run component tests (Vitest)
npm run lint      # ESLint
```

## Go Services

Each Go service follows the same structure and commands:

```bash
cd services/<service-name>

go run cmd/server/main.go   # Run
go test ./...               # Test
go test -v ./...            # Verbose test
go test -race ./...         # Race condition detection
go test -cover ./...        # With coverage
golangci-lint run           # Lint
go build -o <name> cmd/server/main.go  # Build binary
```

### Service-specific notes

| Service | Port | DB | Extra env |
|---------|------|----|-----------|
| api-gateway | 8080 | — | `*_SERVICE_URL` vars |
| video-service | 8081 | video_service_db | — |
| user-service | 8082 | user_service_db | — |
| comment-service | 8083 | comment_service_db | — |
| history-service | 8084 | history_service_db | `VIDEO_SERVICE_URL` |
| notification-service | 8086 | notification_service_db | — |

See [ENVIRONMENT.md](ENVIRONMENT.md) for all env var details.

## PHP Admin Service

```bash
cd services/admin-service
composer install
composer test         # Run tests (if configured)
php -S 0.0.0.0:8085 -t public/   # Dev server
```

## Running a single service against Docker databases

```bash
# Start only the database for the service you're working on
docker compose up -d video-db

# Run the service locally
cd services/video-service
export DB_HOST=localhost DB_PORT=5432 DB_USER=postgres DB_PASSWORD=postgres DB_NAME=video_service_db PORT=8081
go run cmd/server/main.go
```

## Testing

### Backend (Go)
```bash
# All services at once (from repo root)
for svc in video-service user-service comment-service history-service notification-service api-gateway; do
  echo "=== $svc ===" && cd services/$svc && go test ./... && cd ../..
done

# Single service with coverage
cd services/video-service
go test -cover ./...
```

### Frontend
```bash
cd frontend
npm run test             # Run all component tests
npm run test -- --watch  # Watch mode
```

### Linting
```bash
# Go — run inside each service directory
golangci-lint run

# Frontend
cd frontend && npm run lint
```

## CI/CD

GitHub Actions runs on every pull request:
- Go tests and linting for all services
- Frontend component tests
- Docker image builds

Workflow files are in `.github/workflows/`.

## Debugging

**Service won't start:**
```bash
# Check logs
docker compose logs <service-name>

# Check all container status
docker compose ps
```

**Database connection errors:**
- Ensure the correct `DB_NAME` env var matches the created database
- Verify the database healthcheck passes before the service starts (`depends_on` with `condition: service_healthy`)

**Service-to-service communication:**
- All services communicate via HTTP through internal Docker network names (e.g., `http://video-service:8081`)
- For local dev, use `http://localhost:<port>`
