# Database Schema

Each service owns its own PostgreSQL database — no cross-service foreign keys. Services communicate via HTTP APIs.

## Video Service — `video_service_db`

```sql
CREATE TABLE videos (
    id           SERIAL PRIMARY KEY,
    title        VARCHAR(255) NOT NULL,
    description  TEXT,
    url          VARCHAR(500) NOT NULL,
    thumbnail    VARCHAR(500),
    channel_name VARCHAR(100) NOT NULL,
    channel_avatar VARCHAR(500),
    views        INTEGER DEFAULT 0,
    likes        INTEGER DEFAULT 0,
    dislikes     INTEGER DEFAULT 0,
    category     VARCHAR(50) DEFAULT 'General',
    duration     VARCHAR(20),
    uploaded_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Full-text search indexes
CREATE INDEX idx_videos_title       ON videos USING gin(to_tsvector('english', title));
CREATE INDEX idx_videos_description ON videos USING gin(to_tsvector('english', description));
-- Filter / sort indexes
CREATE INDEX idx_videos_channel_name ON videos (channel_name);
CREATE INDEX idx_videos_uploaded_at  ON videos (uploaded_at DESC);
CREATE INDEX idx_videos_category     ON videos (category);
```

## User Service — `user_service_db`

```sql
CREATE TABLE users (
    id         SERIAL PRIMARY KEY,
    username   VARCHAR(50)  UNIQUE NOT NULL,
    email      VARCHAR(100) UNIQUE NOT NULL,
    password   VARCHAR(255) NOT NULL,  -- bcrypt hash
    avatar     VARCHAR(500),
    role       VARCHAR(20) DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Comment Service — `comment_service_db`

```sql
CREATE TABLE comments (
    id         SERIAL PRIMARY KEY,
    video_id   INTEGER NOT NULL,  -- references video-service via API
    user_id    INTEGER NOT NULL,  -- references user-service via API
    content    TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_comments_video_id ON comments (video_id);
CREATE INDEX idx_comments_user_id  ON comments (user_id);
```

## History Service — `history_service_db`

```sql
CREATE TABLE watch_history (
    id         SERIAL PRIMARY KEY,
    user_id    INTEGER NOT NULL,
    video_id   INTEGER NOT NULL,
    watched_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, video_id)
);

CREATE INDEX idx_watch_history_user_id   ON watch_history (user_id);
CREATE INDEX idx_watch_history_watched_at ON watch_history (watched_at DESC);
```

## Notification Service — `notification_service_db`

```sql
CREATE TABLE notifications (
    id         SERIAL PRIMARY KEY,
    user_id    INTEGER NOT NULL,
    type       VARCHAR(50) NOT NULL,  -- e.g. "subscription", "comment", "like"
    title      VARCHAR(255) NOT NULL,
    message    TEXT NOT NULL,
    link       VARCHAR(500),
    is_read    BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_notifications_user_id   ON notifications (user_id);
CREATE INDEX idx_notifications_is_read   ON notifications (user_id, is_read);
CREATE INDEX idx_notifications_created_at ON notifications (created_at DESC);
```

## Admin Service — `admin_service_db`

The admin service schema is managed by Symfony migrations. The initial schema is in `services/admin-service/config/schema.sql` and is auto-applied when the container starts.

Key tables: `blog_posts`, `documentation`, `help_articles`, `email_templates`, `email_logs`, `admin_users`, `moderation_queue`.

## Migration versioning

Each Go service tracks applied migrations in:

```sql
CREATE TABLE schema_migrations (
    version          INTEGER PRIMARY KEY,
    name             VARCHAR(255) NOT NULL,
    description      TEXT,
    applied_at       TIMESTAMP,
    execution_time_ms INTEGER
);
```

To run or roll back migrations in Go services:

```go
manager := migrations.NewMigrationManager(db)
for _, m := range migrations.GetAllMigrations() {
    manager.Register(m)
}
manager.MigrateUp()    // apply pending
manager.MigrateDown()  // rollback last
manager.Status()       // show current state
```
