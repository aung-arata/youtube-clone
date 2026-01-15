# YouTube Clone

A full-stack web-based YouTube clone built with React, Tailwind CSS, Golang, and PostgreSQL.

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

### Backend Features
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
- 🔄 CI/CD pipeline with GitHub Actions
- 🐳 Full Docker Compose setup for all services
- 📦 Multi-stage Docker builds for optimization
- 🧪 Backend unit tests with sqlmock
- 🔍 Linting configuration (golangci-lint)
- 📋 Environment-based configuration

## Tech Stack

### Frontend
- **React** - UI library
- **Tailwind CSS** - Utility-first CSS framework
- **Vite** - Build tool and development server

### Backend
- **Golang** - Backend language
- **Gorilla Mux** - HTTP router
- **PostgreSQL** - Database
- **lib/pq** - PostgreSQL driver

## Project Structure

```
youtube-clone/
├── frontend/              # React frontend
│   ├── src/
│   │   ├── components/   # React components
│   │   ├── pages/        # Page components
│   │   ├── assets/       # Static assets
│   │   ├── App.jsx       # Main App component
│   │   ├── main.jsx      # Entry point
│   │   └── index.css     # Global styles
│   ├── index.html        # HTML template
│   ├── package.json      # Frontend dependencies
│   ├── vite.config.js    # Vite configuration
│   └── tailwind.config.js # Tailwind configuration
│
├── backend/              # Golang backend
│   ├── cmd/
│   │   └── server/       # Server entry point
│   ├── internal/
│   │   ├── handlers/     # HTTP handlers
│   │   ├── models/       # Data models
│   │   └── database/     # Database connection
│   ├── go.mod            # Go dependencies
│   └── .env.example      # Environment variables example
│
└── docker-compose.yml    # Docker setup for PostgreSQL
```

## Getting Started

### Prerequisites

- Node.js 18+ and npm
- Go 1.21+
- PostgreSQL 15+ (or use Docker)
- Docker and Docker Compose (optional, for database)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/aung-arata/youtube-clone.git
   cd youtube-clone
   ```

2. **Set up the database**

   Option A: Using Docker (recommended)
   ```bash
   docker-compose up -d postgres
   ```

   Option B: Using local PostgreSQL
   - Install PostgreSQL
   - Create a database named `youtube_clone`
   ```sql
   CREATE DATABASE youtube_clone;
   ```

3. **Set up the backend**
   ```bash
   cd backend
   
   # Copy environment file
   cp .env.example .env
   
   # Install dependencies
   go mod download
   
   # Run the server
   go run cmd/server/main.go
   ```

   The backend will start on `http://localhost:8080`

4. **Set up the frontend**
   ```bash
   cd frontend
   
   # Copy environment file (optional)
   cp .env.example .env
   
   # Install dependencies
   npm install
   
   # Start development server
   npm run dev
   ```

   The frontend will start on `http://localhost:3000`

### Using Docker Compose (Recommended)

Run the entire stack with a single command:

```bash
# Build and start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop all services
docker-compose down

# Rebuild and restart
docker-compose up -d --build
```

Services:
- Frontend: http://localhost:80
- Backend API: http://localhost:8080
- PostgreSQL: localhost:5432

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

### Videos

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

### Comments

- `GET /api/videos/{videoId}/comments` - Get all comments for a video
- `POST /api/videos/{videoId}/comments` - Create a new comment on a video
- `GET /api/comments/{id}` - Get a specific comment
- `PUT /api/comments/{id}` - Update a comment
- `DELETE /api/comments/{id}` - Delete a comment

### Users

- `POST /api/users` - Create a new user
- `GET /api/users/{id}` - Get user profile
- `PUT /api/users/{id}` - Update user profile

### Watch History

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

### Videos Table
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
```

### Users Table
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

### Comments Table
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

### Watch History Table
```sql
CREATE TABLE watch_history (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    video_id INTEGER REFERENCES videos(id) ON DELETE CASCADE,
    watched_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, video_id)
);
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
```bash
cd backend
go run cmd/server/main.go  # Run server
go build ./...             # Build all packages
go test ./...              # Run tests
go test -v -race -coverprofile=coverage.out ./...  # Run tests with coverage
```

### Testing

**Backend Tests:**
```bash
cd backend
go test -v ./...                    # Run all tests
go test -v ./internal/handlers/...  # Run specific package tests
go test -race ./...                 # Run tests with race detection
go test -cover ./...                # Run tests with coverage
```

**Linting:**
```bash
cd backend
golangci-lint run  # Run Go linter
```

### Code Quality

The project includes:
- Backend unit tests with `go-sqlmock` for database mocking
- CI/CD pipeline with GitHub Actions
- Automatic linting and testing on pull requests
- Code coverage reporting
- Docker builds for all services

## Environment Variables

### Backend (.env)
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


