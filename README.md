# YouTube Clone

A full-stack web-based YouTube clone built with **Microservices Architecture** using React, Tailwind CSS, Golang, and PostgreSQL.

## Architecture

This project follows a **microservices architecture** pattern with:
- **API Gateway** - Single entry point for all client requests
- **Video Service** - Handles video management, views, likes/dislikes
- **User Service** - Manages user profiles and authentication
- **Comment Service** - Handles comment CRUD operations
- **History Service** - Tracks user watch history
- **Separate databases** - Each service has its own PostgreSQL database for data isolation

### Architecture Diagram

```
┌─────────────┐
│   Frontend  │
│   (React)   │
└──────┬──────┘
       │
       ▼
┌─────────────────────────────────────────────────┐
│           API Gateway (Port 8080)               │
│  - Request Routing                              │
│  - CORS Handling                                │
│  - Rate Limiting                                │
│  - Request Logging                              │
└──────┬──────────────┬──────────┬────────────┬───┘
       │              │          │            │
       ▼              ▼          ▼            ▼
┌──────────┐   ┌──────────┐  ┌──────────┐  ┌──────────┐
│  Video   │   │   User   │  │ Comment  │  │ History  │
│ Service  │   │ Service  │  │ Service  │  │ Service  │
│ (8081)   │   │ (8082)   │  │ (8083)   │  │ (8084)   │
└────┬─────┘   └────┬─────┘  └────┬─────┘  └────┬─────┘
     │              │              │              │
     ▼              ▼              ▼              ▼
┌──────────┐   ┌──────────┐  ┌──────────┐  ┌──────────┐
│ Video DB │   │ User DB  │  │Comment DB│  │History DB│
│(Postgres)│   │(Postgres)│  │(Postgres)│  │(Postgres)│
└──────────┘   └──────────┘  └──────────┘  └──────────┘
```

### Benefits of Microservices Architecture

1. **Independent Scaling** - Each service can be scaled independently based on demand
2. **Technology Flexibility** - Different services can use different technologies if needed
3. **Fault Isolation** - Failure in one service doesn't bring down the entire system
4. **Independent Deployment** - Services can be deployed independently without affecting others
5. **Team Autonomy** - Different teams can work on different services independently
6. **Database per Service** - Each service owns its data, ensuring loose coupling

For detailed microservices documentation, see [MICROSERVICES.md](MICROSERVICES.md).

## Features

### Core Features
- 🎨 Modern UI with React and Tailwind CSS
- 🚀 Fast backend API with Golang
- 💾 PostgreSQL database for data persistence
- 📱 Responsive design for all devices
- 🎥 Video listing and playback
- 🔍 Full-text search functionality
- 📊 Video management API with CRUD operations
- 👀 View count tracking
- 📄 Pagination support for efficient data loading
- 💬 Comments system with full CRUD operations
- 👍 Like/dislike functionality for videos
- 🌙 Dark mode support with theme persistence
- 🏷️ Video categories and filtering
- 👤 User profiles with edit functionality
- 📜 Watch history tracking

### Backend Features (Microservices)
- 🏗️ **Microservices Architecture** with independent services
- 🚪 **API Gateway** for routing and middleware
- ✅ Input validation and error handling
- 🔒 Rate limiting middleware (100 requests/minute)
- 📝 Request logging middleware
- 🔎 Search videos by title, description, or channel name
- 📈 View count increment API
- 👍 Like/dislike API endpoints
- 💬 Comment management API (Create, Read, Update, Delete)
- 👤 User profile API (Create, Read, Update)
- 📜 Watch history API with pagination
- 🏷️ Category filtering and management
- 🗄️ **Database per Service** pattern for data isolation
- 🧪 Comprehensive unit tests
- 🐳 Docker support with multi-stage builds

### Frontend Features
- ⚡ Real-time API integration
- 🔄 Loading states and error handling
- 🎯 Dynamic video search
- 📊 View count formatting (K, M)
- ⏱️ Relative time display (e.g., "2 days ago")
- 🎬 Video view tracking on click
- 🌙 Dark mode toggle with localStorage persistence
- 🎨 Dark mode styling across all components
- 🏷️ Category filter with horizontal scroll
- 👤 User profile management component
- 📜 Watch history component and tracking

### DevOps & Code Quality
- 🏗️ **Microservices Architecture** with service isolation
- 🔄 CI/CD pipeline with GitHub Actions
- 🐳 Full Docker Compose setup for all microservices
- 🔀 Service orchestration with Docker Compose
- 📦 Multi-stage Docker builds for optimization
- 🗄️ Separate PostgreSQL databases for each service
- 🧪 Backend unit tests with sqlmock
- 🔍 Linting configuration (golangci-lint)
- 📋 Environment-based configuration

## Tech Stack

### Frontend
- **React** - UI library
- **Tailwind CSS** - Utility-first CSS framework
- **Vite** - Build tool and development server
- **Nginx** - Web server for production

### Backend - Microservices
- **Golang** - Backend language for all services
- **Gorilla Mux** - HTTP router
- **PostgreSQL** - Database (one per service)
- **lib/pq** - PostgreSQL driver

### Microservices
1. **API Gateway** (Port 8080)
   - Routes requests to appropriate services
   - CORS handling
   - Rate limiting (100 req/min)
   - Request logging
   
2. **Video Service** (Port 8081)
   - Video CRUD operations
   - Search functionality
   - View count tracking
   - Like/dislike management
   - Category management
   
3. **User Service** (Port 8082)
   - User profile management
   - User CRUD operations
   
4. **Comment Service** (Port 8083)
   - Comment CRUD operations
   - Video comment associations
   
5. **History Service** (Port 8084)
   - Watch history tracking
   - History retrieval with pagination

### Infrastructure
- **Docker** - Containerization
- **Docker Compose** - Service orchestration
- **PostgreSQL 15** - Multiple database instances

## Project Structure

```
youtube-clone/
├── frontend/                      # React frontend
│   ├── src/
│   │   ├── components/           # React components
│   │   ├── pages/                # Page components
│   │   ├── assets/               # Static assets
│   │   ├── App.jsx               # Main App component
│   │   ├── main.jsx              # Entry point
│   │   └── index.css             # Global styles
│   ├── Dockerfile                # Frontend Docker configuration
│   └── package.json              # Frontend dependencies
│
├── services/                      # Microservices
│   ├── api-gateway/              # API Gateway service
│   │   ├── cmd/server/           # Gateway entry point
│   │   ├── internal/middleware/  # Logging and rate limiting
│   │   ├── Dockerfile
│   │   └── go.mod
│   │
│   ├── video-service/            # Video microservice
│   │   ├── cmd/server/           # Service entry point
│   │   ├── internal/
│   │   │   ├── handlers/         # Video handlers
│   │   │   ├── models/           # Video models
│   │   │   └── database/         # DB connection
│   │   ├── Dockerfile
│   │   └── go.mod
│   │
│   ├── user-service/             # User microservice
│   │   ├── cmd/server/           # Service entry point
│   │   ├── internal/
│   │   │   ├── handlers/         # User handlers
│   │   │   ├── models/           # User models
│   │   │   └── database/         # DB connection
│   │   ├── Dockerfile
│   │   └── go.mod
│   │
│   ├── comment-service/          # Comment microservice
│   │   ├── cmd/server/           # Service entry point
│   │   ├── internal/
│   │   │   ├── handlers/         # Comment handlers
│   │   │   ├── models/           # Comment models
│   │   │   └── database/         # DB connection
│   │   ├── Dockerfile
│   │   └── go.mod
│   │
│   └── history-service/          # History microservice
│       ├── cmd/server/           # Service entry point
│       ├── internal/
│       │   ├── handlers/         # History handlers
│       │   ├── models/           # History models
│       │   └── database/         # DB connection
│       ├── Dockerfile
│       └── go.mod
│
├── backend/                       # Legacy monolithic backend (deprecated)
│
├── docker-compose.yml             # Original monolithic setup (deprecated)
└── docker-compose.microservices.yml  # Microservices Docker Compose
```

## Getting Started

### Prerequisites

- Docker and Docker Compose (recommended for microservices)
- **OR** for local development:
  - Node.js 18+ and npm
  - Go 1.21+
  - PostgreSQL 15+

### Installation with Microservices (Recommended)

1. **Clone the repository**
   ```bash
   git clone https://github.com/aung-arata/youtube-clone.git
   cd youtube-clone
   ```

2. **Run with Docker Compose**
   ```bash
   # Build and start all microservices
   docker-compose -f docker-compose.microservices.yml up -d
   
   # View logs
   docker-compose -f docker-compose.microservices.yml logs -f
   
   # Stop all services
   docker-compose -f docker-compose.microservices.yml down
   
   # Rebuild and restart
   docker-compose -f docker-compose.microservices.yml up -d --build
   ```

   **Services will be available at:**
   - Frontend: http://localhost:80
   - API Gateway: http://localhost:8080
   - Video Service: http://localhost:8081
   - User Service: http://localhost:8082
   - Comment Service: http://localhost:8083
   - History Service: http://localhost:8084

### Installation for Local Development

If you want to run services individually for development:

1. **Clone the repository**
   ```bash
   git clone https://github.com/aung-arata/youtube-clone.git
   cd youtube-clone
   ```

2. **Start PostgreSQL databases**
   ```bash
   # You'll need to create 4 databases:
   CREATE DATABASE video_service_db;
   CREATE DATABASE user_service_db;
   CREATE DATABASE comment_service_db;
   CREATE DATABASE history_service_db;
   ```

3. **Run each microservice**
   
   Terminal 1 - Video Service:
   ```bash
   cd services/video-service
   export DB_HOST=localhost
   export DB_NAME=video_service_db
   go mod download
   go run cmd/server/main.go
   ```
   
   Terminal 2 - User Service:
   ```bash
   cd services/user-service
   export DB_HOST=localhost
   export DB_NAME=user_service_db
   go mod download
   go run cmd/server/main.go
   ```
   
   Terminal 3 - Comment Service:
   ```bash
   cd services/comment-service
   export DB_HOST=localhost
   export DB_NAME=comment_service_db
   go mod download
   go run cmd/server/main.go
   ```
   
   Terminal 4 - History Service:
   ```bash
   cd services/history-service
   export DB_HOST=localhost
   export DB_NAME=history_service_db
   go mod download
   go run cmd/server/main.go
   ```
   
   Terminal 5 - API Gateway:
   ```bash
   cd services/api-gateway
   export VIDEO_SERVICE_URL=http://localhost:8081
   export USER_SERVICE_URL=http://localhost:8082
   export COMMENT_SERVICE_URL=http://localhost:8083
   export HISTORY_SERVICE_URL=http://localhost:8084
   go mod download
   go run cmd/server/main.go
   ```

4. **Set up the frontend**
   ```bash
   cd frontend
   
   # Install dependencies
   npm install
   
   # Start development server
   npm run dev
   ```

   The frontend will start on `http://localhost:3000`

### Monolithic Setup (Legacy - Deprecated)

The original monolithic backend is still available in the `backend/` directory for backward compatibility:

```bash
# Using Docker Compose (monolithic)
docker-compose up -d

# Services:
# - Frontend: http://localhost:80
# - Backend API: http://localhost:8080
# - PostgreSQL: localhost:5432
```

**Note:** The microservices architecture is the recommended approach.

### Building for Production

**Frontend:**
```bash
cd frontend
npm run build
```

The production build will be in the `frontend/dist` directory.

**Backend:**
```bash
cd backend
go build -o server cmd/server/main.go
./server
```

## API Endpoints

All API requests go through the **API Gateway** at `http://localhost:8080/api`

The gateway routes requests to the appropriate microservice:

### Videos (Video Service)

- `GET /api/videos` - Get all videos
  - Query Parameters:
    - `q` (optional): Search query for title, description, or channel name
    - `category` (optional): Filter by category
    - `page` (optional): Page number (default: 1)
    - `limit` (optional): Items per page (default: 20, max: 100)
- `GET /api/videos/categories` - Get all unique video categories
- `GET /api/videos/{id}` - Get a specific video
- `POST /api/videos` - Create a new video
- `POST /api/videos/{id}/views` - Increment view count
- `POST /api/videos/{id}/like` - Increment like count
- `POST /api/videos/{id}/dislike` - Increment dislike count

### Comments (Comment Service)

- `GET /api/videos/{videoId}/comments` - Get all comments for a video
- `POST /api/videos/{videoId}/comments` - Create a new comment on a video
- `GET /api/comments/{id}` - Get a specific comment
- `PUT /api/comments/{id}` - Update a comment
- `DELETE /api/comments/{id}` - Delete a comment

### Users (User Service)

- `POST /api/users` - Create a new user
- `GET /api/users/{id}` - Get user profile
- `PUT /api/users/{id}` - Update user profile

### Watch History (History Service)

- `POST /api/users/{userId}/history` - Add video to watch history
- `GET /api/users/{userId}/history` - Get user's watch history
  - Query Parameters:
    - `page` (optional): Page number (default: 1)
    - `limit` (optional): Items per page (default: 20, max: 100)

### System

- `GET /api/health` - Health check endpoint

### Example API Usage

**Get all videos:**
```bash
curl http://localhost:8080/api/videos
```

**Search videos:**
```bash
curl "http://localhost:8080/api/videos?q=react&page=1&limit=10"
```

**Get a specific video:**
```bash
curl http://localhost:8080/api/videos/1
```

**Create a new video:**
```bash
curl -X POST http://localhost:8080/api/videos \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My Video",
    "description": "Video description",
    "url": "https://example.com/video.mp4",
    "thumbnail": "https://example.com/thumb.jpg",
    "channel_name": "My Channel",
    "channel_avatar": "https://example.com/avatar.jpg",
    "duration": "10:30"
  }'
```

**Increment video views:**
```bash
curl -X POST http://localhost:8080/api/videos/1/views
```

**Like a video:**
```bash
curl -X POST http://localhost:8080/api/videos/1/like
```

**Dislike a video:**
```bash
curl -X POST http://localhost:8080/api/videos/1/dislike
```

**Get comments for a video:**
```bash
curl http://localhost:8080/api/videos/1/comments
```

**Create a comment:**
```bash
curl -X POST http://localhost:8080/api/videos/1/comments \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "content": "Great video!"
  }'
```

**Update a comment:**
```bash
curl -X PUT http://localhost:8080/api/comments/1 \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Updated comment!"
  }'
```

**Delete a comment:**
```bash
curl -X DELETE http://localhost:8080/api/comments/1
```

**Filter videos by category:**
```bash
curl "http://localhost:8080/api/videos?category=Education"
```

**Get all categories:**
```bash
curl http://localhost:8080/api/videos/categories
```

**Create a user:**
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "avatar": "https://example.com/avatar.jpg"
  }'
```

**Get user profile:**
```bash
curl http://localhost:8080/api/users/1
```

**Add video to watch history:**
```bash
curl -X POST http://localhost:8080/api/users/1/history \
  -H "Content-Type: application/json" \
  -d '{
    "video_id": 5
  }'
```

**Get watch history:**
```bash
curl http://localhost:8080/api/users/1/history
```

For more detailed API documentation, see [API.md](API.md).

## Database Schema

Each microservice has its own PostgreSQL database following the **Database per Service** pattern.

### Video Service Database (`video_service_db`)

**Videos Table:**
```sql
CREATE TABLE videos (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    url VARCHAR(500) NOT NULL,
    thumbnail VARCHAR(500),
    channel_name VARCHAR(100) NOT NULL,
    channel_avatar VARCHAR(500),
    views INTEGER DEFAULT 0,
    likes INTEGER DEFAULT 0,
    dislikes INTEGER DEFAULT 0,
    category VARCHAR(50) DEFAULT 'General',
    duration VARCHAR(20),
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_videos_title ON videos USING gin(to_tsvector('english', title));
CREATE INDEX idx_videos_description ON videos USING gin(to_tsvector('english', description));
CREATE INDEX idx_videos_channel_name ON videos (channel_name);
CREATE INDEX idx_videos_uploaded_at ON videos (uploaded_at DESC);
CREATE INDEX idx_videos_category ON videos (category);
```

### User Service Database (`user_service_db`)

**Users Table:**
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    avatar VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Comment Service Database (`comment_service_db`)

**Comments Table:**
```sql
CREATE TABLE comments (
    id SERIAL PRIMARY KEY,
    video_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_comments_video_id ON comments (video_id);
CREATE INDEX idx_comments_user_id ON comments (user_id);
```

**Note:** In microservices architecture, foreign key constraints to other services' data are removed. Services communicate through APIs.

### History Service Database (`history_service_db`)

**Watch History Table:**
```sql
CREATE TABLE watch_history (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    video_id INTEGER NOT NULL,
    watched_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, video_id)
);

-- Indexes for performance
CREATE INDEX idx_watch_history_user_id ON watch_history (user_id);
CREATE INDEX idx_watch_history_watched_at ON watch_history (watched_at DESC);
```
```sql
CREATE TABLE comments (
    id SERIAL PRIMARY KEY,
    video_id INTEGER REFERENCES videos(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Development

### Frontend Development
```bash
cd frontend
npm run dev     # Start development server
npm run build   # Build for production
npm run preview # Preview production build
```

### Backend Development

Each microservice can be developed independently:

**Video Service:**
```bash
cd services/video-service
go run cmd/server/main.go  # Run server
go test ./...              # Run tests
```

**User Service:**
```bash
cd services/user-service
go run cmd/server/main.go  # Run server
go test ./...              # Run tests
```

**Comment Service:**
```bash
cd services/comment-service
go run cmd/server/main.go  # Run server
go test ./...              # Run tests
```

**History Service:**
```bash
cd services/history-service
go run cmd/server/main.go  # Run server
go test ./...              # Run tests
```

**API Gateway:**
```bash
cd services/api-gateway
go run cmd/server/main.go  # Run server
```

### Testing

Each microservice has its own tests. Run tests for each service:

```bash
# Test Video Service
cd services/video-service
go test -v ./...
go test -race ./...                 # With race detection
go test -cover ./...                # With coverage

# Test other services similarly
cd services/user-service && go test -v ./...
cd services/comment-service && go test -v ./...
cd services/history-service && go test -v ./...
```

**Linting:**
```bash
# Lint each service
cd services/video-service && golangci-lint run
cd services/user-service && golangci-lint run
cd services/comment-service && golangci-lint run
cd services/history-service && golangci-lint run
cd services/api-gateway && golangci-lint run
```

### Code Quality

The project includes:
- Microservices architecture with independent services
- Backend unit tests with `go-sqlmock` for database mocking
- CI/CD pipeline with GitHub Actions
- Automatic linting and testing on pull requests
- Code coverage reporting
- Docker builds for all services
- Database per service pattern

## Environment Variables

### Microservices Configuration

**Video Service:**
```env
DB_HOST=video-db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=video_service_db
PORT=8081
```

**User Service:**
```env
DB_HOST=user-db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=user_service_db
PORT=8082
```

**Comment Service:**
```env
DB_HOST=comment-db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=comment_service_db
PORT=8083
```

**History Service:**
```env
DB_HOST=history-db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=history_service_db
PORT=8084
```

**API Gateway:**
```env
VIDEO_SERVICE_URL=http://video-service:8081
USER_SERVICE_URL=http://user-service:8082
COMMENT_SERVICE_URL=http://comment-service:8083
HISTORY_SERVICE_URL=http://history-service:8084
PORT=8080
```

### Legacy Backend (.env) - Deprecated
```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=youtube_clone

# Server Configuration
PORT=8080
```

### Frontend (.env)
```env
# API Configuration
VITE_API_URL=http://localhost:8080
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the ISC License.

## Acknowledgments

- YouTube for the design inspiration
- React and Tailwind CSS communities
- Go community

## Future Enhancements

### Completed ✅
- [x] Input validation and error handling
- [x] Search functionality implementation
- [x] View count increment
- [x] Pagination support
- [x] Request logging middleware
- [x] Rate limiting middleware
- [x] Backend unit tests
- [x] CI/CD pipeline
- [x] Docker containerization
- [x] Frontend API integration
- [x] Loading states and error handling
- [x] Comments system with full CRUD operations
- [x] Like/dislike functionality for videos
- [x] Dark mode support
- [x] Video categories and filtering
- [x] User profile pages
- [x] Video watch history tracking

### Planned 🚀
- [ ] User authentication and authorization (JWT-based)
- [ ] Video upload functionality with file handling
- [ ] Subscription system for channels
- [ ] Playlist management
- [ ] Video recommendations algorithm
- [ ] Frontend component tests
- [ ] API integration tests
- [ ] HTTPS support and security headers
- [ ] Comprehensive API documentation (Swagger/OpenAPI)
- [ ] Database connection pooling optimization
- [ ] Database migration versioning system


