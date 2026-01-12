.PHONY: help frontend-install frontend-dev frontend-build backend-run backend-build backend-test backend-lint db-up db-down db-seed clean setup docker-up docker-down docker-logs

help:
	@echo "Available commands:"
	@echo "  make setup             - Initial project setup"
	@echo "  make frontend-install  - Install frontend dependencies"
	@echo "  make frontend-dev      - Start frontend development server"
	@echo "  make frontend-build    - Build frontend for production"
	@echo "  make backend-run       - Run backend server"
	@echo "  make backend-build     - Build backend binary"
	@echo "  make backend-test      - Run backend tests"
	@echo "  make backend-lint      - Run backend linter"
	@echo "  make db-up             - Start PostgreSQL with Docker"
	@echo "  make db-down           - Stop PostgreSQL"
	@echo "  make db-seed           - Seed database with sample data"
	@echo "  make docker-up         - Start all services with Docker Compose"
	@echo "  make docker-down       - Stop all Docker services"
	@echo "  make docker-logs       - View Docker logs"
	@echo "  make clean             - Clean build artifacts"
	@echo "  make test              - Run all tests"

setup:
	@echo "Setting up project..."
	@chmod +x dev-setup.sh
	@./dev-setup.sh

frontend-install:
	cd frontend && npm install

frontend-dev:
	cd frontend && npm run dev

frontend-build:
	cd frontend && npm run build

backend-run:
	cd backend && go run cmd/server/main.go

backend-build:
	cd backend && go build -o server cmd/server/main.go

backend-test:
	cd backend && go test -v -race -cover ./...

backend-lint:
	cd backend && golangci-lint run

db-up:
	docker-compose up -d postgres

db-down:
	docker-compose down

db-seed:
	docker-compose exec -T postgres psql -U postgres -d youtube_clone < backend/seed.sql

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

test: backend-test
	@echo "All tests completed"

clean:
	rm -rf frontend/dist
	rm -rf frontend/node_modules
	rm -f backend/server
	docker-compose down -v
