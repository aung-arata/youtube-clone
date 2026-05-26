.PHONY: help frontend-install frontend-dev frontend-build gateway-run db-up db-down clean setup docker-up docker-down docker-logs

help:
	@echo "Available commands:"
	@echo "  make setup             - Initial project setup"
	@echo "  make frontend-install  - Install frontend dependencies"
	@echo "  make frontend-dev      - Start frontend development server"
	@echo "  make frontend-build    - Build frontend for production"
	@echo "  make gateway-run       - Run api-gateway locally (port 8080)"
	@echo "  make db-up             - Start all databases with Docker Compose"
	@echo "  make db-down           - Stop all Docker services"
	@echo "  make docker-up         - Start all services with Docker Compose"
	@echo "  make docker-down       - Stop all Docker services"
	@echo "  make docker-logs       - View Docker logs"
	@echo "  make clean             - Clean build artifacts"

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

gateway-run:
	cd services/api-gateway && go run cmd/server/main.go

db-up:
	docker compose up -d video-db user-db comment-db history-db admin-db notification-db

db-down:
	docker compose down

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f

clean:
	rm -rf frontend/dist
	rm -rf frontend/node_modules
	docker compose down -v
