# YouTube Clone

A web-based YouTube clone built with React, Tailwind CSS, Golang, and PostgreSQL.

## Features

- 🎨 Modern UI with React and Tailwind CSS
- 🚀 Fast backend API with Golang
- 💾 PostgreSQL database for data persistence
- 📱 Responsive design for all devices
- 🎥 Video listing and playback
- 🔍 Search functionality (UI ready)
- 📊 Video management API

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
   docker-compose up -d
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
   
   # Install dependencies
   npm install
   
   # Start development server
   npm run dev
   ```

   The frontend will start on `http://localhost:3000`

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
- `GET /api/videos/{id}` - Get a specific video
- `POST /api/videos` - Create a new video
- `GET /api/health` - Health check endpoint

### Example API Usage

**Get all videos:**
```bash
curl http://localhost:8080/api/videos
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
```

## Environment Variables

### Backend (.env)
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=youtube_clone
PORT=8080
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

- [ ] User authentication and authorization
- [ ] Video upload functionality
- [ ] Comments system
- [ ] Like/dislike feature
- [ ] Subscription system
- [ ] Search functionality implementation
- [ ] Video recommendations
- [ ] Playlist management
- [ ] Video analytics

