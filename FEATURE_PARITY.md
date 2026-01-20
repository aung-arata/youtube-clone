# Feature Parity: Monolith vs Microservices

This document describes the feature parity between the monolithic backend and the microservices architecture.

## Core API Features

The following core features are available in **both** architectures with equivalent functionality:

### Video Features
- ✅ List videos with pagination and search
- ✅ Get single video by ID
- ✅ Create new video
- ✅ Increment view count
- ✅ Like/dislike videos
- ✅ Get video categories

**Implementation:**
- **Monolith**: `/backend/internal/handlers/video_handler.go`
- **Microservices**: `/services/video-service/internal/handlers/video_handler.go`

### User Features
- ✅ Get user by ID
- ✅ Create new user
- ✅ Update user profile

**Implementation:**
- **Monolith**: `/backend/internal/handlers/user_handler.go`
- **Microservices**: `/services/user-service/internal/handlers/user_handler.go`

### Comment Features
- ✅ Get comments for a video
- ✅ Get single comment by ID
- ✅ Create new comment
- ✅ Update comment
- ✅ Delete comment

**Implementation:**
- **Monolith**: `/backend/internal/handlers/comment_handler.go`
- **Microservices**: `/services/comment-service/internal/handlers/comment_handler.go`

### Watch History Features
- ✅ Add video to watch history
- ✅ Get user's watch history with video details

**Implementation:**
- **Monolith**: `/backend/internal/handlers/history_handler.go` (uses SQL JOIN)
- **Microservices**: `/services/history-service/internal/handlers/history_handler.go` (calls video-service API)

**Note**: Both implementations return the same data structure - watch history enriched with full video details. The monolith uses a database JOIN while the microservices version makes HTTP calls to the video service.

## Microservices-Only Features

The following features are **only available in the microservices architecture** and are not present in the monolith:

### Admin Service (PHP/Symfony)

Located in `/services/admin-service/`, this is a separate Symfony 8.0 application providing administrative and content management features.

#### CMS (Content Management System)
- ❌ **Blog Posts**: Full CRUD operations, automatic slug generation, categories, publish workflow
- ❌ **Documentation**: Organized documentation with categories and sorting
- ❌ **Help Center**: Searchable help articles with view tracking

#### Email System
- ❌ **Email Templates**: Template management with variable substitution
- ❌ **Email Logs**: Delivery tracking and history

#### Reporting
- ❌ **Analytics Reports**: Data aggregation and analytics
- ❌ **User Activity Reports**: User behavior tracking
- ❌ **Video Statistics**: Category-based video statistics

#### Admin Dashboard
- ❌ **User Management**: Admin user roles and permissions
- ❌ **Content Moderation**: Moderation queue and workflow
- ❌ **System Monitoring**: Dashboard statistics and monitoring

**Why Not in Monolith?**

The admin service is intentionally separate because:
1. It serves a different purpose (administrative/CMS vs core API)
2. It uses a different technology stack (PHP/Symfony vs Go)
3. It can be scaled independently
4. It has its own database (`admin_service_db`)
5. It's designed for administrative users, not end users

### API Gateway

Located in `/services/api-gateway/`, this service is **only needed for microservices**.

Features:
- Request routing to appropriate services
- CORS handling
- Rate limiting (100 requests/minute)
- Request logging

**Why Not in Monolith?**

The monolith backend directly handles all requests and doesn't need a gateway. In microservices, the gateway provides a single entry point and routes requests to the appropriate service.

## Architectural Differences

### Data Access

**Monolith:**
- Single database: `youtube_clone`
- All tables in one database
- Can perform cross-table JOINs efficiently
- Single database connection pool

**Microservices:**
- Four separate databases:
  - `video_service_db` (videos)
  - `user_service_db` (users)
  - `comment_service_db` (comments)
  - `history_service_db` (watch_history)
  - `admin_service_db` (admin features)
- Each service owns its data
- Cross-service data access via HTTP APIs
- Database isolation for better scalability

### Example: Watch History with Video Details

**Monolith Approach:**
```sql
SELECT v.*, wh.watched_at
FROM watch_history wh
JOIN videos v ON wh.video_id = v.id
WHERE wh.user_id = $1
```

**Microservices Approach:**
1. History service queries `watch_history` table
2. For each history entry, makes HTTP GET to `video-service/videos/{id}`
3. Combines results into enriched response

Both return the same data structure to clients.

## Deployment Differences

### Monolith

```bash
docker-compose up -d
# Starts: PostgreSQL + Backend + Frontend
# Port: 8080
```

### Microservices

```bash
docker-compose -f docker-compose.microservices.yml up -d
# Starts: 
# - 5 databases (video, user, comment, history, admin)
# - 5 Go services (video, user, comment, history, gateway)
# - 1 PHP service (admin)
# - Frontend
# Ports: 8080 (gateway), 8081-8085 (services)
```

## Choosing Between Architectures

### Use Monolith When:
- Starting a new project
- Team is small (< 5 developers)
- Traffic is moderate
- Simpler deployment preferred
- Don't need admin/CMS features

### Use Microservices When:
- Scaling specific features independently
- Multiple teams working on different services
- Need admin/CMS functionality
- Want technology flexibility per service
- Have DevOps resources for complex deployment

## API Compatibility

Both architectures expose the same core API endpoints through their respective entry points:

- **Monolith**: `http://localhost:8080/api/*`
- **Microservices**: `http://localhost:8080/api/*` (via API Gateway)

Frontend applications can switch between architectures by simply changing the `VITE_API_URL` environment variable.

## Summary

| Feature | Monolith | Microservices |
|---------|----------|---------------|
| Videos | ✅ | ✅ |
| Users | ✅ | ✅ |
| Comments | ✅ | ✅ |
| Watch History | ✅ | ✅ |
| CMS (Blog, Docs, Help) | ❌ | ✅ |
| Email Templates | ❌ | ✅ |
| Admin Dashboard | ❌ | ✅ |
| Reports & Analytics | ❌ | ✅ |
| Independent Scaling | ❌ | ✅ |
| Simple Deployment | ✅ | ❌ |
| Single Database | ✅ | ❌ |
| Technology Flexibility | ❌ | ✅ |

## Conclusion

The core video platform features have **full parity** between both architectures. The microservices architecture adds **administrative and CMS features** that are intentionally not in the monolith, as they serve different use cases and can evolve independently.
