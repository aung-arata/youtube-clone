# Contributing to YouTube Clone

Thank you for your interest in contributing to this YouTube Clone project! This document provides guidelines and instructions for contributing.

## Table of Contents
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Project Structure](#project-structure)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/youtube-clone.git
   cd youtube-clone
   ```
3. **Create a branch** for your changes:
   ```bash
   git checkout -b feature/your-feature-name
   ```

## Development Setup

### Quick Setup (Recommended)
Use the development setup script:
```bash
./dev-setup.sh
```

### Manual Setup

1. **Start PostgreSQL:**
   ```bash
   docker-compose up -d postgres
   ```

2. **Set up Backend:**
   ```bash
   cd backend
   cp .env.example .env
   go mod download
   go run cmd/server/main.go
   ```

3. **Set up Frontend:**
   ```bash
   cd frontend
   cp .env.example .env
   npm install
   npm run dev
   ```

### Seed Database (Optional)
```bash
docker-compose exec postgres psql -U postgres -d youtube_clone < backend/seed.sql
```

## Project Structure

```
youtube-clone/
├── backend/              # Go backend
│   ├── cmd/
│   │   └── server/      # Main application entry point
│   ├── internal/
│   │   ├── handlers/    # HTTP request handlers
│   │   ├── models/      # Data models
│   │   ├── database/    # Database connection and migrations
│   │   └── middleware/  # HTTP middleware
│   ├── Dockerfile
│   └── seed.sql         # Sample data
│
├── frontend/            # React frontend
│   ├── src/
│   │   ├── components/  # React components
│   │   ├── App.jsx      # Main app component
│   │   └── main.jsx     # Entry point
│   ├── Dockerfile
│   └── nginx.conf
│
├── .github/
│   └── workflows/       # CI/CD pipelines
└── docker-compose.yml   # Docker services configuration
```

## Coding Standards

### Backend (Go)

- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` to format your code
- Run `golangci-lint run` before committing
- Write tests for all new functionality
- Keep handlers thin; move business logic to separate packages
- Use meaningful variable names
- Add comments for exported functions and types

Example:
```go
// GetVideo retrieves a single video by ID
func (h *VideoHandler) GetVideo(w http.ResponseWriter, r *http.Request) {
    // Implementation
}
```

### Frontend (React)

- Use functional components with hooks
- Follow React best practices
- Use meaningful component and variable names
- Keep components small and focused
- Use PropTypes or TypeScript for type checking (if added)
- Format code consistently

Example:
```jsx
function VideoCard({ video }) {
  // Component logic
  return (
    // JSX
  )
}
```

### Database

- Use migrations for schema changes
- Never modify existing migrations
- Add appropriate indexes for queries
- Use meaningful column names
- Add foreign key constraints where appropriate

## Testing

### Backend Tests

Run all tests:
```bash
cd backend
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

Run tests with race detection:
```bash
go test -race ./...
```

### Writing Tests

- Write table-driven tests when appropriate
- Use `go-sqlmock` for database testing
- Test both success and error cases
- Aim for high code coverage

Example:
```go
func TestGetVideo_Success(t *testing.T) {
    // Setup
    // Execute
    // Assert
}
```

### Frontend Tests

(Tests to be added)

```bash
cd frontend
npm test
```

## Submitting Changes

1. **Ensure all tests pass:**
   ```bash
   cd backend && go test ./...
   ```

2. **Format your code:**
   ```bash
   cd backend && gofmt -w .
   ```

3. **Commit your changes:**
   ```bash
   git add .
   git commit -m "Add feature: brief description"
   ```
   
   Follow [Conventional Commits](https://www.conventionalcommits.org/):
   - `feat:` for new features
   - `fix:` for bug fixes
   - `docs:` for documentation changes
   - `test:` for test additions or changes
   - `refactor:` for code refactoring
   - `chore:` for maintenance tasks

4. **Push to your fork:**
   ```bash
   git push origin feature/your-feature-name
   ```

5. **Create a Pull Request** on GitHub with:
   - Clear description of changes
   - Reference to any related issues
   - Screenshots for UI changes
   - Test results if applicable

## Pull Request Checklist

- [ ] Code follows project style guidelines
- [ ] Tests added for new functionality
- [ ] All tests pass
- [ ] Documentation updated (if needed)
- [ ] No unnecessary dependencies added
- [ ] Commit messages follow conventional commits
- [ ] PR description clearly explains the changes

## Common Issues

### Port Already in Use
If you get "port already in use" errors:
```bash
# Find process using port 8080
lsof -i :8080
# Kill the process
kill -9 <PID>
```

### Database Connection Issues
- Ensure PostgreSQL is running: `docker-compose ps`
- Check credentials in `.env` file
- Verify database exists: `docker-compose exec postgres psql -U postgres -l`

### Frontend Not Connecting to Backend
- Verify backend is running on port 8080
- Check `VITE_API_URL` in frontend `.env`
- Check CORS configuration in backend

## Getting Help

- Open an issue for bugs or feature requests
- Use GitHub Discussions for questions
- Check existing issues before creating new ones

## Code of Conduct

- Be respectful and inclusive
- Provide constructive feedback
- Focus on the code, not the person
- Help others learn and grow

Thank you for contributing! 🎉
