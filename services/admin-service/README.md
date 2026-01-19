# Admin Service

PHP-based admin dashboard and CMS service for YouTube Clone.

## Features

### Admin Dashboard
- User management
- Content moderation queue
- System statistics and monitoring
- Integration with all Go microservices

### CMS (Content Management System)
- **Blog Posts**: Create, edit, publish, and manage blog posts
- **Documentation**: Maintain technical documentation
- **Help Center**: Create and manage help articles for users

### Email Template System
- Create and manage HTML email templates
- Variable substitution support (e.g., `{{username}}`, `{{video_title}}`)
- Email sending API
- Email delivery logging and tracking

### Batch Reporting
- Analytics report generation
- User activity reports
- Video statistics reports
- Scheduled batch jobs

## API Endpoints

### Health Check
```
GET /health
```

### Admin Dashboard
```
GET  /admin                    # Dashboard statistics
GET  /admin/users              # List users with admin data
GET  /admin/moderation         # Content moderation queue
```

### CMS - Blog Posts
```
GET    /cms/blog               # List blog posts
POST   /cms/blog               # Create blog post
GET    /cms/blog/{id}          # Get blog post
PUT    /cms/blog/{id}          # Update blog post
DELETE /cms/blog/{id}          # Delete blog post
```

### CMS - Documentation
```
GET    /cms/docs               # List documentation
POST   /cms/docs               # Create documentation
```

### CMS - Help Articles
```
GET    /cms/help               # List help articles
POST   /cms/help               # Create help article
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
GET /reports/analytics         # Generate analytics report
GET /reports/users             # Generate user report
GET /reports/videos            # Generate video report
```

## Database Schema

The admin service uses PostgreSQL with the following tables:

- `admin_users` - User roles and permissions
- `moderation_queue` - Content moderation queue
- `blog_posts` - Blog posts and articles
- `documentation` - Technical documentation
- `help_articles` - Help center articles
- `email_templates` - Email templates
- `email_log` - Email delivery log
- `reports` - Generated reports

## Integration with Go Services

The admin service communicates with Go microservices via HTTP:

- **Video Service** (Port 8081) - Fetch video data, statistics
- **User Service** (Port 8082) - User management
- **Comment Service** (Port 8083) - Comment moderation
- **History Service** (Port 8084) - User activity tracking

## Development

### Local Development

```bash
cd services/admin-service
php -S localhost:8085 -t public
```

### With Docker

```bash
docker-compose -f docker-compose.microservices.yml up admin-service
```

## Technology Stack

- **PHP 8.2+**
- **PostgreSQL 15**
- **Nginx** (web server)
- **PHP-FPM** (FastCGI Process Manager)

## Environment Variables

```env
DATABASE_URL=postgresql://postgres:postgres@admin-db:5432/admin_service_db
PORT=8085
VIDEO_SERVICE_URL=http://video-service:8081
USER_SERVICE_URL=http://user-service:8082
COMMENT_SERVICE_URL=http://comment-service:8083
HISTORY_SERVICE_URL=http://history-service:8084
```

## Sample Requests

### Create Blog Post

```bash
curl -X POST http://localhost:8085/cms/blog \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Welcome to our platform",
    "content": "We are excited to launch...",
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
    "name": "video_notification",
    "subject": "New video from {{channel_name}}",
    "html_content": "<h1>{{video_title}}</h1><p>Watch now!</p>",
    "category": "notification"
  }'
```

### Send Email

```bash
curl -X POST http://localhost:8085/email/send \
  -H "Content-Type: application/json" \
  -d '{
    "to": "user@example.com",
    "template_id": 1,
    "variables": {
      "username": "John",
      "video_title": "My Awesome Video"
    }
  }'
```

### Generate Analytics Report

```bash
curl http://localhost:8085/reports/analytics?period=week
```

## Architecture

```
┌─────────────────┐
│  Admin Frontend │ (Future: React Admin UI)
└────────┬────────┘
         │
         ▼
┌─────────────────────────────┐
│   Admin Service (PHP)       │
│   Port: 8085                │
│   ├─ Admin Dashboard        │
│   ├─ CMS (Blog, Docs, Help) │
│   ├─ Email Templates        │
│   └─ Batch Reports          │
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

## Benefits

1. **Rapid Development**: PHP enables faster development of admin features (2-3x faster than Go)
2. **Rich Ecosystem**: Easy access to CMS libraries, email templates, and admin panels
3. **Isolated**: Failure in admin service doesn't affect user-facing Go services
4. **Flexibility**: Easy to add admin features without rebuilding Go services

## Future Enhancements

- Admin web UI (React-based dashboard)
- Advanced analytics with charts
- Automated email campaigns
- A/B testing for email templates
- Role-based access control (RBAC)
- Audit logging for admin actions
