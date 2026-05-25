# Getting Started

## Prerequisites

- **Docker + Docker Compose** — recommended for running all services
- **Or** for local development:
  - Go 1.21+
  - Node.js 18+ and npm
  - PHP 8.2+ with Composer (admin service only)
  - PostgreSQL 15+

## Docker Compose (Recommended)

The fastest way to run the full stack:

```bash
git clone https://github.com/aung-arata/youtube-clone.git
cd youtube-clone

# Build and start all services
docker compose up -d

# View logs (all services)
docker compose logs -f

# View logs (specific service)
docker compose logs -f video-service

# Rebuild after code changes
docker compose up -d --build

# Stop everything
docker compose down
```

**Running services:**

| Service | URL |
|---------|-----|
| Frontend | http://localhost:80 |
| API Gateway | http://localhost:8080 |
| Video Service | http://localhost:8081 |
| User Service | http://localhost:8082 |
| Comment Service | http://localhost:8083 |
| History Service | http://localhost:8084 |
| Admin Service | http://localhost:8085 |
| Notification Service | http://localhost:8086 |

## Local Development Setup

For running individual services without Docker.

### 1. Create PostgreSQL databases

```bash
psql -U postgres -c "CREATE DATABASE video_service_db;"
psql -U postgres -c "CREATE DATABASE user_service_db;"
psql -U postgres -c "CREATE DATABASE comment_service_db;"
psql -U postgres -c "CREATE DATABASE history_service_db;"
psql -U postgres -c "CREATE DATABASE admin_service_db;"
psql -U postgres -c "CREATE DATABASE notification_service_db;"
```

### 2. Run backend services (each in a separate terminal)

**Video Service (port 8081):**
```bash
cd services/video-service
export DB_HOST=localhost DB_PORT=5432 DB_USER=postgres DB_PASSWORD=postgres DB_NAME=video_service_db PORT=8081
go mod download && go run cmd/server/main.go
```

**User Service (port 8082):**
```bash
cd services/user-service
export DB_HOST=localhost DB_PORT=5432 DB_USER=postgres DB_PASSWORD=postgres DB_NAME=user_service_db PORT=8082
go mod download && go run cmd/server/main.go
```

**Comment Service (port 8083):**
```bash
cd services/comment-service
export DB_HOST=localhost DB_PORT=5432 DB_USER=postgres DB_PASSWORD=postgres DB_NAME=comment_service_db PORT=8083
go mod download && go run cmd/server/main.go
```

**History Service (port 8084):**
```bash
cd services/history-service
export DB_HOST=localhost DB_PORT=5432 DB_USER=postgres DB_PASSWORD=postgres DB_NAME=history_service_db VIDEO_SERVICE_URL=http://localhost:8081 PORT=8084
go mod download && go run cmd/server/main.go
```

**Notification Service (port 8086):**
```bash
cd services/notification-service
export DB_HOST=localhost DB_PORT=5432 DB_USER=postgres DB_PASSWORD=postgres DB_NAME=notification_service_db PORT=8086
go mod download && go run cmd/server/main.go
```

**API Gateway (port 8080):**
```bash
cd services/api-gateway
export VIDEO_SERVICE_URL=http://localhost:8081 USER_SERVICE_URL=http://localhost:8082 COMMENT_SERVICE_URL=http://localhost:8083 HISTORY_SERVICE_URL=http://localhost:8084 NOTIFICATION_SERVICE_URL=http://localhost:8086 PORT=8080
go mod download && go run cmd/server/main.go
```

**Admin Service (port 8085):**
```bash
cd services/admin-service
composer install
export DATABASE_URL="postgresql://postgres:postgres@localhost:5432/admin_service_db" PORT=8085
php -S 0.0.0.0:8085 -t public/
```

### 3. Run the frontend

```bash
cd frontend
npm install
cp .env.example .env   # set VITE_API_URL=http://localhost:8080
npm run dev
```

Frontend starts on `http://localhost:5173` (Vite default).

## Production Build

**Frontend:**
```bash
cd frontend
npm run build
# Output: frontend/dist/
```

**Go services:**
```bash
# Example for video-service
cd services/video-service
go build -o video-service cmd/server/main.go
./video-service
```
