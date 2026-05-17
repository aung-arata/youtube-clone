# Microservices Architecture Guide

## Overview

This YouTube Clone uses a microservices architecture with 7 independent services (6 in Go, 1 in PHP), each with its own PostgreSQL database.

## Services

### 1. API Gateway (Port 8080)
- **Purpose**: Single entry point for all client requests
- **Responsibilities**:
  - Request routing to appropriate services
  - CORS handling
  - Rate limiting (100 requests/minute)
  - Request logging
- **Technology**: Go, Gorilla Mux

### 2. Video Service (Port 8081)
- **Purpose**: Manage all video-related operations
- **Responsibilities**:
  - Video CRUD operations
  - Search functionality
  - View count tracking
  - Like/dislike management
  - Category management
- **Database**: video_service_db
- **Technology**: Go, PostgreSQL

### 3. User Service (Port 8082)
- **Purpose**: Manage user profiles and data
- **Responsibilities**:
  - User CRUD operations
  - User profile management
- **Database**: user_service_db
- **Technology**: Go, PostgreSQL

### 4. Comment Service (Port 8083)
- **Purpose**: Handle all comment operations
- **Responsibilities**:
  - Comment CRUD operations
  - Associate comments with videos
- **Database**: comment_service_db
- **Technology**: Go, PostgreSQL

### 5. History Service (Port 8084)
- **Purpose**: Track user watch history
- **Responsibilities**:
  - Record video watches
  - Retrieve watch history with pagination
- **Database**: history_service_db
- **Technology**: Go, PostgreSQL
- **Note**: Calls Video Service to enrich history entries with video details

### 6. Notification Service (Port 8086)
- **Purpose**: Manage user notifications with real-time delivery
- **Responsibilities**:
  - Create and store notifications
  - Mark notifications as read
  - Unread notification count
  - WebSocket endpoint for real-time push
- **Database**: notification_service_db
- **Technology**: Go, PostgreSQL

### 7. Admin Service (Port 8085)
- **Purpose**: Administrative dashboard and content management
- **Responsibilities**:
  - Admin user management and roles
  - Content moderation queue
  - Blog post and documentation CMS
  - Help center articles
  - Email template management
  - Batch reporting and analytics
  - Integration with all Go services via HTTP
- **Database**: admin_service_db
- **Technology**: PHP 8.2+, Symfony
- **Note**: Schema is auto-applied from `services/admin-service/config/schema.sql` on first start

## Quick Start

### Using Docker Compose (Recommended)

```bash
# Start all services
docker compose -f docker-compose.microservices.yml up -d

# View logs for all services
docker compose -f docker-compose.microservices.yml logs -f

# View logs for a specific service
docker compose -f docker-compose.microservices.yml logs -f video-service

# Stop all services
docker compose -f docker-compose.microservices.yml down

# Rebuild and restart
docker compose -f docker-compose.microservices.yml up -d --build
```

### Service Health Checks

Each service has a health endpoint:

```bash
# API Gateway
curl http://localhost:8080/api/health

# Video Service
curl http://localhost:8081/health

# User Service
curl http://localhost:8082/health

# Comment Service
curl http://localhost:8083/health

# History Service
curl http://localhost:8084/health

# Notification Service
curl http://localhost:8086/health

# Admin Service
curl http://localhost:8085/health
```

## Data Flow Example

When a user views a video:

1. Frontend sends request to API Gateway: `POST /api/videos/1/views`
2. API Gateway routes to Video Service
3. Video Service increments view count in video_service_db
4. Response flows back through API Gateway to Frontend

When tracking watch history:

1. Frontend sends request to API Gateway: `POST /api/users/1/history`
2. API Gateway routes to History Service
3. History Service records watch in history_service_db
4. Response flows back through API Gateway to Frontend

## Communication Patterns

- **Client-to-Services**: All requests go through API Gateway (synchronous HTTP)
- **Service-to-Service**: Currently, services are independent; if needed, they can communicate via HTTP

## Database Schema per Service

Each service owns its database and schema:

- **video_service_db**: videos table
- **user_service_db**: users table
- **comment_service_db**: comments table
- **history_service_db**: watch_history table
- **notification_service_db**: notifications table
- **admin_service_db**: admin/CMS tables (Symfony-managed)

See [DATABASE.md](DATABASE.md) for full schemas.

## Advantages Over Monolithic

1. **Scalability**: Scale only the services that need it (e.g., video service during high traffic)
2. **Deployment**: Deploy services independently without downtime
3. **Fault Tolerance**: One service failure doesn't crash entire application
4. **Technology Freedom**: Each service can use different tech stack if needed
5. **Team Organization**: Different teams can own different services
6. **Database Isolation**: No shared database, reducing coupling

## Migration from Monolithic

The original monolithic backend is still available in the `backend/` directory for reference. To migrate:

1. Data would need to be migrated from `youtube_clone` into the six separate service databases
2. Ensure all client requests go through the API Gateway endpoints

## Troubleshooting

### Service won't start
- Check if the database is ready (healthcheck)
- Check environment variables are set correctly
- View logs: `docker compose -f docker-compose.microservices.yml logs <service-name>`

### Can't connect to service
- Ensure all services are running: `docker compose -f docker-compose.microservices.yml ps`
- Check port mappings in docker-compose.microservices.yml
- Verify network connectivity between containers

### Database connection issues
- Each service needs its own database
- Check DB_NAME environment variable matches the database created
- Verify database is healthy before service starts (depends_on with condition)

## Development

To develop a single service:

```bash
# Start only required databases
docker compose -f docker-compose.microservices.yml up -d video-db

# Run service locally
cd services/video-service
export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=video_service_db
export DB_USER=postgres
export DB_PASSWORD=postgres
go run cmd/server/main.go
```

## Testing

Test individual services:

```bash
# Test video service
cd services/video-service
go test ./...

# Test with coverage
go test -cover ./...
```

## Future Enhancements

- Service-to-service authentication
- Distributed tracing
- Service mesh (e.g., Istio)
- Message queue for async communication
- API versioning
- Centralized configuration management
- Circuit breakers for resilience
