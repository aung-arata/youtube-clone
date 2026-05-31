#!/bin/bash

# YouTube Clone Development Setup Script

set -e

echo "🚀 YouTube Clone Setup Script"
echo "=============================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo -e "${RED}❌ Docker is not installed. Please install Docker first.${NC}"
    exit 1
fi

# Support both Docker Compose v2 plugin (docker compose) and standalone v1 (docker-compose)
if docker compose version &> /dev/null 2>&1; then
    DOCKER_COMPOSE="docker compose"
elif command -v docker-compose &> /dev/null; then
    DOCKER_COMPOSE="docker-compose"
else
    echo -e "${RED}❌ Docker Compose is not installed. Please install Docker Compose first.${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Docker and Docker Compose found${NC}"

# Create root .env file if it doesn't exist
if [ ! -f .env ]; then
    echo -e "${YELLOW}📝 Creating .env file...${NC}"
    cp .env.example .env
    echo -e "${GREEN}✅ .env created (set a strong POSTGRES_PASSWORD before deploying)${NC}"
else
    echo -e "${GREEN}✅ .env already exists${NC}"
fi

# Create frontend .env file if it doesn't exist
if [ ! -f frontend/.env ]; then
    if [ -f frontend/.env.example ]; then
        echo -e "${YELLOW}📝 Creating frontend .env file...${NC}"
        cp frontend/.env.example frontend/.env
        echo -e "${GREEN}✅ Frontend .env created${NC}"
    fi
else
    echo -e "${GREEN}✅ Frontend .env already exists${NC}"
fi

# Start all databases
echo -e "${YELLOW}🐘 Starting databases...${NC}"
$DOCKER_COMPOSE up -d video-db user-db comment-db history-db admin-db notification-db

# Wait for databases to be ready
echo -e "${YELLOW}⏳ Waiting for databases to be ready...${NC}"
sleep 8

echo -e "${GREEN}✅ Databases started${NC}"

# Ask user what they want to run
echo ""
echo "What would you like to do?"
echo "1) Start full stack with Docker Compose"
echo "2) Start databases only (run services locally)"
echo "3) Start frontend only (for development)"
echo "4) Exit"
read -p "Enter your choice (1-4): " choice

case $choice in
    1)
        echo -e "${YELLOW}🚀 Starting full stack...${NC}"
        $DOCKER_COMPOSE up -d
        echo ""
        echo -e "${GREEN}✅ Full stack is running!${NC}"
        echo -e "Frontend:      ${GREEN}http://localhost${NC}"
        echo -e "API Gateway:   ${GREEN}http://localhost:8080${NC}"
        echo -e "User Service:  ${GREEN}http://localhost:8081${NC}"
        echo -e "Video Service: ${GREEN}http://localhost:8082${NC}"
        echo ""
        echo "To view logs: $DOCKER_COMPOSE logs -f"
        echo "To stop:      $DOCKER_COMPOSE down"
        ;;
    2)
        echo -e "${YELLOW}🔧 Databases are running. Start services locally:${NC}"
        echo "  cd services/api-gateway          && go run cmd/server/main.go"
        echo "  cd services/user-service         && go run cmd/server/main.go"
        echo "  cd services/video-service        && go run cmd/server/main.go"
        echo "  cd services/comment-service      && go run cmd/server/main.go"
        echo "  cd services/history-service      && go run cmd/server/main.go"
        echo "  cd services/notification-service && go run cmd/server/main.go"
        ;;
    3)
        echo -e "${YELLOW}💻 Frontend development mode${NC}"
        echo "Make sure the API gateway is running on http://localhost:8080"
        echo "Run: cd frontend && npm install && npm run dev"
        ;;
    4)
        echo "Goodbye!"
        exit 0
        ;;
    *)
        echo -e "${RED}Invalid choice${NC}"
        exit 1
        ;;
esac
