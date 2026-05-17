-- Admin Service Database Schema
-- Applied automatically on first container start via docker-entrypoint-initdb.d

CREATE TABLE IF NOT EXISTS admin_users (
    id         SERIAL PRIMARY KEY,
    user_id    INTEGER NOT NULL,
    role       VARCHAR(50)  NOT NULL DEFAULT 'admin',
    status     VARCHAR(20)  NOT NULL DEFAULT 'active',
    last_login TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS blog_posts (
    id           SERIAL PRIMARY KEY,
    title        VARCHAR(255) NOT NULL,
    slug         VARCHAR(255) NOT NULL UNIQUE,
    content      TEXT NOT NULL,
    excerpt      TEXT,
    author_id    INTEGER NOT NULL,
    category     VARCHAR(100) NOT NULL DEFAULT 'general',
    status       VARCHAR(20)  NOT NULL DEFAULT 'draft',
    published_at TIMESTAMP,
    created_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS documentation (
    id         SERIAL PRIMARY KEY,
    title      VARCHAR(255) NOT NULL,
    slug       VARCHAR(255) NOT NULL UNIQUE,
    content    TEXT NOT NULL,
    category   VARCHAR(100) NOT NULL,
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS help_articles (
    id         SERIAL PRIMARY KEY,
    title      VARCHAR(255) NOT NULL,
    slug       VARCHAR(255) NOT NULL UNIQUE,
    content    TEXT NOT NULL,
    category   VARCHAR(100) NOT NULL,
    view_count INTEGER NOT NULL DEFAULT 0,
    is_published BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS email_templates (
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL UNIQUE,
    subject    VARCHAR(255) NOT NULL,
    body       TEXT NOT NULL,
    variables  TEXT,
    category   VARCHAR(100) NOT NULL DEFAULT 'general',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS email_logs (
    id          SERIAL PRIMARY KEY,
    to_email    VARCHAR(255) NOT NULL,
    template_id INTEGER,
    subject     VARCHAR(255),
    status      VARCHAR(20) NOT NULL DEFAULT 'pending',
    sent_at     TIMESTAMP,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS moderation_queue (
    id           SERIAL PRIMARY KEY,
    content_type VARCHAR(50) NOT NULL,
    content_id   INTEGER NOT NULL,
    reporter_id  INTEGER,
    reason       TEXT,
    status       VARCHAR(20) NOT NULL DEFAULT 'pending',
    reviewed_by  INTEGER,
    reviewed_at  TIMESTAMP,
    created_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS reports (
    id           SERIAL PRIMARY KEY,
    type         VARCHAR(100) NOT NULL,
    period       VARCHAR(100),
    data         JSONB NOT NULL DEFAULT '{}',
    generated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
