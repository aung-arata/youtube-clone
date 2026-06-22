# YouTube Clone

A full-stack YouTube clone using a **hybrid microservices architecture**: React + Tailwind CSS frontend, Go microservices, PHP admin service, and PostgreSQL per service.

## Architecture

```
┌─────────────┐
│   Frontend  │
│   (React)   │
└──────┬──────┘
       │
       ▼
┌────────────────────────────────────────────────────────────────────────┐
│                    API Gateway (Port 8080) [Go]                        │
│  - Request Routing  - CORS  - Rate Limiting  - Request Logging         │
└──────┬──────────┬──────────┬────────────┬──────────────┬──────────────┘
       │          │          │            │              │
       ▼          ▼          ▼            ▼              ▼
┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐
│  Video   │ │   User   │ │ Comment  │ │ History  │ │Notif.    │ │ Admin    │
│ Service  │ │ Service  │ │ Service  │ │ Service  │ │ Service  │ │ Service  │
│ (8081)   │ │ (8082)   │ │ (8083)   │ │ (8084)   │ │ (8086)   │ │ (8085)  │
│  [Go]    │ │  [Go]    │ │  [Go]    │ │  [Go]    │ │  [Go]    │ │ [PHP]   │
└────┬─────┘ └────┬─────┘ └────┬─────┘ └────┬─────┘ └────┬─────┘ └────┬───┘
     ▼            ▼            ▼            ▼            ▼            ▼
  video_db     user_db    comment_db   history_db   notif_db     admin_db
```

See [MICROSERVICES.md](MICROSERVICES.md) for details on each service.

## Quick Start

```bash
git clone https://github.com/aung-arata/youtube-clone.git
cd youtube-clone
cp .env.example .env   # set POSTGRES_PASSWORD
docker compose up -d
```

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

## Features

- Video upload with drag-and-drop, thumbnail preview, and real-time progress
- Video playback with view count (deduplicated per session — no inflating on reload)
- Comments with nested replies (one level, YouTube-style `@mention` prefill)
- User auth (JWT), channel profiles, watch history
- Notification system, playlist management
- Admin panel (PHP/Symfony)

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Frontend | React, Tailwind CSS, Vite, Nginx |
| Backend (services) | Go 1.21+, Gorilla Mux |
| Backend (admin) | PHP 8.2+, Symfony |
| Database | PostgreSQL 15 (one per service) |
| Infrastructure | Docker, Docker Compose |

## Documentation

| Document | Description |
|----------|-------------|
| [GETTING_STARTED.md](GETTING_STARTED.md) | Full setup guide — Docker and local dev |
| [MICROSERVICES.md](MICROSERVICES.md) | Service responsibilities, ports, health checks |
| [API.md](API.md) | Complete API reference for the gateway |
| [DATABASE.md](DATABASE.md) | Database schemas for all services |
| [ENVIRONMENT.md](ENVIRONMENT.md) | Environment variables for all services |
| [DEVELOPMENT.md](DEVELOPMENT.md) | Per-service dev commands, testing, linting |
| [CONTRIBUTING.md](CONTRIBUTING.md) | Contribution guidelines and coding standards |
| [ROADMAP.md](ROADMAP.md) | Project phases and completion status |
| [ARCHITECTURE_COMPARISON.md](ARCHITECTURE_COMPARISON.md) | Monolith vs microservices trade-offs |
| [FEATURE_PARITY.md](FEATURE_PARITY.md) | Feature comparison between architectures |
| [PHP_GO_INTEGRATION_ANALYSIS.md](PHP_GO_INTEGRATION_ANALYSIS.md) | Analysis behind the hybrid tech choice |

## License

ISC License


