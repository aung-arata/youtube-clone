.PHONY: help frontend-install frontend-dev frontend-build backend-run backend-build db-up db-down clean

help:
	@echo "Available commands:"
	@echo "  make frontend-install  - Install frontend dependencies"
	@echo "  make frontend-dev      - Start frontend development server"
	@echo "  make frontend-build    - Build frontend for production"
	@echo "  make backend-run       - Run backend server"
	@echo "  make backend-build     - Build backend binary"
	@echo "  make db-up             - Start PostgreSQL with Docker"
	@echo "  make db-down           - Stop PostgreSQL"
	@echo "  make clean             - Clean build artifacts"

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

db-up:
	docker-compose up -d

db-down:
	docker-compose down

clean:
	rm -rf frontend/dist
	rm -rf frontend/node_modules
	rm -f backend/server
