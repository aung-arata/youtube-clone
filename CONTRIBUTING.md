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

```bash
git clone https://github.com/YOUR_USERNAME/youtube-clone.git
cd youtube-clone
docker compose -f docker-compose.microservices.yml up -d
```

See [GETTING_STARTED.md](GETTING_STARTED.md) for local (non-Docker) setup.

## Project Structure

```
youtube-clone/
├── frontend/            # React + Tailwind CSS
│   ├── src/
│   │   ├── components/  # React components
│   │   ├── App.jsx
│   │   └── main.jsx
│   ├── Dockerfile
│   └── nginx.conf
│
├── services/
│   ├── api-gateway/           # Go — routing, CORS, rate limiting
│   ├── video-service/         # Go — video CRUD, search, analytics
│   ├── user-service/          # Go — user profiles
│   ├── comment-service/       # Go — comment CRUD
│   ├── history-service/       # Go — watch history
│   ├── notification-service/  # Go — notifications, WebSocket
│   └── admin-service/         # PHP/Symfony — admin dashboard, CMS
│
├── backend/                   # Deprecated monolith (reference only)
├── docker-compose.microservices.yml
├── docker-compose.yml         # Deprecated monolith compose
└── .github/workflows/         # CI/CD pipelines
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

See [DEVELOPMENT.md](DEVELOPMENT.md) for full test and lint commands.

### Backend (Go)

Run tests inside the service directory:
```bash
cd services/video-service
go test ./...
go test -cover ./...
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
   cd services/video-service && go test ./...
   cd frontend && npm test
   ```

2. **Format your code:**
   ```bash
   # Go (inside service directory)
   gofmt -w .
   golangci-lint run
   
   # Frontend
   cd frontend && npm run lint
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
lsof -i :8080  # find process
kill -9 <PID>
```

### Database Connection Issues
- Ensure the correct service database is running: `docker compose -f docker-compose.microservices.yml ps`
- Check env vars in your shell or compose file
- See [ENVIRONMENT.md](ENVIRONMENT.md) for all variables

### Frontend Not Connecting to Backend
- Verify the API Gateway is running on port 8080
- Check `VITE_API_URL` in `frontend/.env`
- Check CORS configuration in `services/api-gateway`

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
