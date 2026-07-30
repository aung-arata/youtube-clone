# Environment Variables

All services are configured via environment variables. In Docker Compose these are set in `docker-compose.yml`, with sensitive values sourced from a root `.env` file. For local dev, export them in your shell or use a `.env` file.

## Root `.env` (required for Docker Compose)

Copy `.env.example` to `.env` and set secret values:

```env
POSTGRES_PASSWORD=your_strong_password_here
JWT_SECRET=your_jwt_secret_key_here
```

These variables are shared across microservice containers.

## API Gateway (port 8080)

```env
PORT=8080
VIDEO_SERVICE_URL=http://video-service:8081
USER_SERVICE_URL=http://user-service:8082
COMMENT_SERVICE_URL=http://comment-service:8083
HISTORY_SERVICE_URL=http://history-service:8084
NOTIFICATION_SERVICE_URL=http://notification-service:8086
```

> For local dev replace service hostnames with `localhost` and the corresponding port.

## Video Service (port 8081)

```env
PORT=8081
DB_HOST=video-db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=$POSTGRES_PASSWORD
DB_NAME=video_service_db
```

Optional — for HTTPS and transcoding:
```env
TLS_ENABLED=true
TLS_CERT_FILE=/path/to/cert.pem
TLS_KEY_FILE=/path/to/key.pem
FFMPEG_PATH=/usr/bin/ffmpeg
TRANSCODING_OUTPUT_DIR=/uploads/transcoded
```

## User Service (port 8082)

```env
PORT=8082
DB_HOST=user-db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=$POSTGRES_PASSWORD
DB_NAME=user_service_db
```

## Comment Service (port 8083)

```env
PORT=8083
DB_HOST=comment-db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=$POSTGRES_PASSWORD
DB_NAME=comment_service_db
```

## History Service (port 8084)

```env
PORT=8084
DB_HOST=history-db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=$POSTGRES_PASSWORD
DB_NAME=history_service_db
VIDEO_SERVICE_URL=http://video-service:8081
```

## Notification Service (port 8086)

```env
PORT=8086
DB_HOST=notification-db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=$POSTGRES_PASSWORD
DB_NAME=notification_service_db
```

## Admin Service (port 8085) — PHP/Symfony

```env
PORT=8085
DATABASE_URL=postgresql://postgres:$POSTGRES_PASSWORD@admin-db:5432/admin_service_db
VIDEO_SERVICE_URL=http://video-service:8081
USER_SERVICE_URL=http://user-service:8082
COMMENT_SERVICE_URL=http://comment-service:8083
HISTORY_SERVICE_URL=http://history-service:8084
```

## Frontend

Create `frontend/.env` (copy from `.env.example`):

```env
VITE_API_URL=http://localhost:8080
```

## Legacy Monolith (Deprecated)

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=$POSTGRES_PASSWORD
DB_NAME=youtube_clone
PORT=8080
```
