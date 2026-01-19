# Admin Service - Symfony 6.4

Admin Dashboard and CMS Service for YouTube Clone built with **Symfony 6.4 LTS**.

## Features

### Admin Dashboard
- User management with role-based access
- Content moderation queue
- System statistics and monitoring
- Real-time integration with Go microservices

### CMS (Content Management System)
- **Blog Posts**: Full CRUD with Doctrine ORM, automatic slug generation, categories, publish workflow
- **Documentation**: Organized documentation with categories and sorting
- **Help Center**: Searchable help articles with view tracking

### Email Template System
- Template management with Symfony Twig integration
- Variable substitution support
- Email delivery tracking
- HTML and plain text templates

### Batch Reporting
- Analytics reports with data aggregation
- User activity reports
- Video statistics by category
- Report persistence with JSON data storage

## Technology Stack

- **Symfony 6.4 LTS** - Full-stack PHP framework
- **Doctrine ORM** - Database abstraction and ORM
- **PostgreSQL 15** - Database
- **Symfony HTTP Client** - Integration with Go services
- **Twig** - Template engine (ready for admin UI)

## Project Structure (Symfony Standard)

```
admin-service/
├── bin/
│   └── console              # Symfony console commands
├── config/
│   ├── packages/           # Bundle configurations
│   ├── routes.yaml         # Route definitions
│   └── services.yaml       # Service container config
├── migrations/             # Database migrations
├── public/
│   └── index.php           # Application entry point
├── src/
│   ├── Controller/         # Symfony controllers with attributes
│   │   ├── AdminController.php
│   │   ├── CMSController.php
│   │   ├── EmailController.php
│   │   ├── ReportController.php
│   │   └── HealthController.php
│   ├── Entity/             # Doctrine entities
│   │   ├── BlogPost.php
│   │   ├── Documentation.php
│   │   ├── HelpArticle.php
│   │   ├── EmailTemplate.php
│   │   ├── EmailLog.php
│   │   ├── Report.php
│   │   ├── AdminUser.php
│   │   └── ModerationQueue.php
│   ├── Repository/         # Doctrine repositories
│   ├── Service/            # Business logic services
│   │   └── GoServiceClient.php
│   └── Kernel.php          # Application kernel
├── templates/              # Twig templates (for future admin UI)
├── var/                    # Cache and logs
├── composer.json           # Dependencies
└── .env                    # Environment variables
```

## API Endpoints

All endpoints use Symfony's routing system with attributes.

### Health Check
```
GET /health
```

### Admin Dashboard
```
GET  /admin                 # Dashboard statistics
GET  /admin/users           # List users with pagination
GET  /admin/moderation      # Content moderation queue
```

### CMS - Blog Posts
```
GET    /cms/blog            # List blog posts (paginated)
POST   /cms/blog            # Create blog post
GET    /cms/blog/{id}       # Get blog post
PUT    /cms/blog/{id}       # Update blog post
DELETE /cms/blog/{id}       # Delete blog post
```

### CMS - Documentation
```
GET  /cms/docs              # List documentation
POST /cms/docs              # Create documentation
```

### CMS - Help Articles
```
GET  /cms/help              # List help articles
POST /cms/help              # Create help article
```

### Email Templates
```
GET    /email/templates        # List templates
POST   /email/templates        # Create template
GET    /email/templates/{id}   # Get template
PUT    /email/templates/{id}   # Update template
DELETE /email/templates/{id}   # Delete template
POST   /email/send             # Send email using template
```

### Reports
```
GET /reports/analytics      # Generate analytics report
GET /reports/users          # Generate user report
GET /reports/videos         # Generate video report
```

## Database Schema

Managed by Doctrine ORM with migrations.

### Entities

- **BlogPost** - Blog posts with slug, categories, publish status
- **Documentation** - Technical documentation with sorting
- **HelpArticle** - Help center articles with view tracking
- **EmailTemplate** - Email templates with variable support
- **EmailLog** - Email delivery tracking
- **Report** - Generated reports with JSON data
- **AdminUser** - Admin user roles and permissions
- **ModerationQueue** - Content moderation workflow

## Development

### Local Development with Symfony CLI

```bash
cd services/admin-service

# Install dependencies
composer install

# Create database
php bin/console doctrine:database:create

# Run migrations
php bin/console doctrine:migrations:migrate

# Start development server
symfony server:start --port=8085
# or
php -S localhost:8085 -t public
```

### Database Migrations

```bash
# Create migration
php bin/console make:migration

# Run migrations
php bin/console doctrine:migrations:migrate

# Check schema
php bin/console doctrine:schema:validate
```

### Symfony Console Commands

```bash
# Clear cache
php bin/console cache:clear

# List routes
php bin/console debug:router

# List services
php bin/console debug:container

# Generate entity
php bin/console make:entity
```

## Docker Deployment

```bash
# Build
docker build -t admin-service .

# Run
docker run -p 8085:8085 \
  -e DATABASE_URL="postgresql://postgres:postgres@admin-db:5432/admin_service_db" \
  -e VIDEO_SERVICE_URL="http://video-service:8081" \
  admin-service
```

Or use Docker Compose:

```bash
docker-compose -f docker-compose.microservices.yml up -d admin-service
```

## Environment Variables

```env
APP_ENV=prod
APP_SECRET=your-secret-key-here
DATABASE_URL=postgresql://postgres:postgres@admin-db:5432/admin_service_db
VIDEO_SERVICE_URL=http://video-service:8081
USER_SERVICE_URL=http://user-service:8082
COMMENT_SERVICE_URL=http://comment-service:8083
HISTORY_SERVICE_URL=http://history-service:8084
```

## Symfony Features Used

- **Symfony Attributes** - Modern PHP 8 attributes for routing and configuration
- **Doctrine ORM** - Full entity management with migrations
- **Dependency Injection** - Autowiring and autoconfiguration
- **HTTP Client** - Symfony HTTP Client for Go service integration
- **Logging** - Monolog integration for error tracking
- **Validation** - Symfony Validator (ready to use)
- **Forms** - Symfony Forms (ready for admin UI)
- **Twig** - Template engine (for future admin UI)

## Sample Requests

### Create Blog Post (Symfony Way)

```bash
curl -X POST http://localhost:8085/cms/blog \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Welcome to Symfony Admin",
    "content": "This is a Symfony-powered admin service...",
    "author_id": 1,
    "category": "announcement",
    "status": "published"
  }'
```

### Create Email Template

```bash
curl -X POST http://localhost:8085/email/templates \
  -H "Content-Type: application/json" \
  -d '{
    "name": "welcome_email",
    "subject": "Welcome {{username}}!",
    "html_content": "<h1>Welcome {{username}}</h1><p>Thanks for joining!</p>",
    "category": "user"
  }'
```

### Generate Analytics Report

```bash
curl http://localhost:8085/reports/analytics?period=week
```

## Benefits of Symfony

1. **Mature Framework**: Battle-tested with 15+ years of development
2. **Best Practices**: Follows PSR standards and SOLID principles
3. **ORM Integration**: Doctrine provides powerful database abstraction
4. **Dependency Injection**: Full DI container with autowiring
5. **Debugging**: Symfony Profiler for development
6. **Extensibility**: Large ecosystem of bundles
7. **LTS Support**: Symfony 6.4 supported until November 2027

## Future Enhancements

- Admin web UI with Symfony UX and Stimulus
- API Platform integration for auto-generated API docs
- Symfony Messenger for async processing
- Symfony Mailer for actual email sending
- Symfony Security for authentication/authorization
- EasyAdmin bundle for rapid admin panel development
- Webpack Encore for frontend asset management

## Testing

```bash
# Run tests (when implemented)
php bin/console --env=test doctrine:database:create
php bin/console --env=test doctrine:migrations:migrate
vendor/bin/phpunit
```

## Architecture

```
┌─────────────────┐
│  Admin Frontend │ (Future: Symfony UX + Twig)
└────────┬────────┘
         │
         ▼
┌─────────────────────────────┐
│   Symfony 6.4 Application   │
│   ├─ Controllers (Routing)  │
│   ├─ Services (Business)    │
│   ├─ Entities (ORM)         │
│   └─ Repositories (Data)    │
└────────┬────────────────────┘
         │
         ├─────────────┬──────────────┬──────────────┐
         ▼             ▼              ▼              ▼
    ┌─────────┐  ┌─────────┐    ┌─────────┐   ┌─────────┐
    │ Video   │  │ User    │    │ Comment │   │ History │
    │ Service │  │ Service │    │ Service │   │ Service │
    │ (Go)    │  │ (Go)    │    │ (Go)    │   │ (Go)    │
    └─────────┘  └─────────┘    └─────────┘   └─────────┘
```

## Support

For Symfony-specific documentation, see:
- https://symfony.com/doc/6.4/index.html
- https://symfony.com/doc/current/doctrine.html
- https://symfony.com/doc/current/routing.html
