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

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo -e "${RED}❌ Docker Compose is not installed. Please install Docker Compose first.${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Docker and Docker Compose found${NC}"

# Create backend .env file if it doesn't exist
if [ ! -f backend/.env ]; then
    echo -e "${YELLOW}📝 Creating backend .env file...${NC}"
    cp backend/.env.example backend/.env
    echo -e "${GREEN}✅ Backend .env created${NC}"
else
    echo -e "${GREEN}✅ Backend .env already exists${NC}"
fi

# Create frontend .env file if it doesn't exist
if [ ! -f frontend/.env ]; then
    echo -e "${YELLOW}📝 Creating frontend .env file...${NC}"
    cp frontend/.env.example frontend/.env
    echo -e "${GREEN}✅ Frontend .env created${NC}"
else
    echo -e "${GREEN}✅ Frontend .env already exists${NC}"
fi

# Start PostgreSQL
echo -e "${YELLOW}🐘 Starting PostgreSQL...${NC}"
docker-compose up -d postgres

# Wait for PostgreSQL to be ready
echo -e "${YELLOW}⏳ Waiting for PostgreSQL to be ready...${NC}"
sleep 5

# Check if PostgreSQL is healthy
if docker-compose ps postgres | grep -q "healthy"; then
    echo -e "${GREEN}✅ PostgreSQL is ready${NC}"
else
    echo -e "${YELLOW}⏳ Still waiting for PostgreSQL...${NC}"
    sleep 5
fi

# Seed the database (optional)
read -p "Do you want to seed the database with sample data? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}🌱 Seeding database...${NC}"
    # Wait for migrations to complete
    sleep 2
    # Check if videos table exists
    if docker-compose exec -T postgres psql -U postgres -d youtube_clone -c "\dt videos" | grep -q "videos"; then
        docker-compose exec -T postgres psql -U postgres -d youtube_clone < backend/seed.sql
        echo -e "${GREEN}✅ Database seeded${NC}"
    else
        echo -e "${RED}❌ Videos table not found. Please ensure migrations have run.${NC}"
    fi
fi

# Ask user what they want to run
echo ""
echo "What would you like to do?"
echo "1) Start full stack with Docker Compose"
echo "2) Start backend only (for development)"
echo "3) Start frontend only (for development)"
echo "4) Exit"
read -p "Enter your choice (1-4): " choice

case $choice in
    1)
        echo -e "${YELLOW}🚀 Starting full stack...${NC}"
        docker-compose up -d
        echo ""
        echo -e "${GREEN}✅ Full stack is running!${NC}"
        echo -e "Frontend: ${GREEN}http://localhost:80${NC}"
        echo -e "Backend API: ${GREEN}http://localhost:8080${NC}"
        echo -e "PostgreSQL: ${GREEN}localhost:5432${NC}"
        echo ""
        echo "To view logs: docker-compose logs -f"
        echo "To stop: docker-compose down"
        ;;
    2)
        echo -e "${YELLOW}🔧 Backend development mode${NC}"
        echo "PostgreSQL is already running on localhost:5432"
        echo "Run: cd backend && go run cmd/server/main.go"
        ;;
    3)
        echo -e "${YELLOW}💻 Frontend development mode${NC}"
        echo "Make sure the backend is running on http://localhost:8080"
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
