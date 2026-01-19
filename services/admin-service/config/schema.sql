-- Admin Service Database Schema

-- Admin Users Table
CREATE TABLE IF NOT EXISTS admin_users (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL UNIQUE,
    role VARCHAR(50) DEFAULT 'user',
    status VARCHAR(20) DEFAULT 'active',
    last_login TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Moderation Queue Table
CREATE TABLE IF NOT EXISTS moderation_queue (
    id SERIAL PRIMARY KEY,
    content_type VARCHAR(50) NOT NULL,
    content_id INTEGER NOT NULL,
    reporter_id INTEGER,
    reason TEXT,
    status VARCHAR(20) DEFAULT 'pending',
    reviewed_by INTEGER,
    reviewed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Blog Posts Table
CREATE TABLE IF NOT EXISTS blog_posts (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    content TEXT NOT NULL,
    excerpt TEXT,
    author_id INTEGER NOT NULL,
    category VARCHAR(100) DEFAULT 'general',
    status VARCHAR(20) DEFAULT 'draft',
    published_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_blog_posts_slug ON blog_posts(slug);
CREATE INDEX IF NOT EXISTS idx_blog_posts_status ON blog_posts(status);
CREATE INDEX IF NOT EXISTS idx_blog_posts_published ON blog_posts(published_at DESC);

-- Documentation Table
CREATE TABLE IF NOT EXISTS documentation (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    content TEXT NOT NULL,
    category VARCHAR(100) DEFAULT 'general',
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_documentation_category ON documentation(category, sort_order);

-- Help Articles Table
CREATE TABLE IF NOT EXISTS help_articles (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    content TEXT NOT NULL,
    category VARCHAR(100) DEFAULT 'general',
    views INTEGER DEFAULT 0,
    helpful_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_help_articles_category ON help_articles(category);

-- Email Templates Table
CREATE TABLE IF NOT EXISTS email_templates (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    subject VARCHAR(500) NOT NULL,
    html_content TEXT NOT NULL,
    text_content TEXT,
    category VARCHAR(100) DEFAULT 'general',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Email Log Table
CREATE TABLE IF NOT EXISTS email_log (
    id SERIAL PRIMARY KEY,
    to_email VARCHAR(255) NOT NULL,
    template_id INTEGER REFERENCES email_templates(id),
    subject VARCHAR(500),
    status VARCHAR(20) DEFAULT 'pending',
    sent_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_email_log_status ON email_log(status);
CREATE INDEX IF NOT EXISTS idx_email_log_sent ON email_log(sent_at DESC);

-- Reports Table
CREATE TABLE IF NOT EXISTS reports (
    id SERIAL PRIMARY KEY,
    type VARCHAR(50) NOT NULL,
    period VARCHAR(20),
    data JSONB,
    generated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_reports_type ON reports(type, generated_at DESC);

-- Insert sample data for testing
INSERT INTO admin_users (user_id, role, status) VALUES
    (1, 'admin', 'active'),
    (2, 'moderator', 'active'),
    (3, 'user', 'active')
ON CONFLICT (user_id) DO NOTHING;

INSERT INTO email_templates (name, subject, html_content, text_content, category) VALUES
    ('welcome', 'Welcome to YouTube Clone!', 
     '<h1>Welcome {{username}}!</h1><p>Thank you for joining our platform.</p>', 
     'Welcome {{username}}! Thank you for joining our platform.',
     'user'),
    ('video_upload_success', 'Your video has been uploaded',
     '<h1>Hi {{username}}</h1><p>Your video "{{video_title}}" has been successfully uploaded.</p>',
     'Hi {{username}}, Your video "{{video_title}}" has been successfully uploaded.',
     'video')
ON CONFLICT (name) DO NOTHING;

INSERT INTO blog_posts (title, slug, content, excerpt, author_id, category, status, published_at) VALUES
    ('Welcome to YouTube Clone', 'welcome-to-youtube-clone',
     'We are excited to launch our YouTube clone platform. This platform offers all the features you love about video sharing.',
     'We are excited to launch our YouTube clone platform.',
     1, 'announcement', 'published', CURRENT_TIMESTAMP)
ON CONFLICT (slug) DO NOTHING;

INSERT INTO documentation (title, slug, content, category, sort_order) VALUES
    ('Getting Started', 'getting-started',
     'Learn how to use YouTube Clone. Upload videos, create playlists, and engage with the community.',
     'basics', 1),
    ('Video Upload Guide', 'video-upload-guide',
     'Step-by-step guide on uploading videos to our platform.',
     'basics', 2),
    ('API Documentation', 'api-documentation',
     'Complete API documentation for developers.',
     'developer', 1)
ON CONFLICT (slug) DO NOTHING;

INSERT INTO help_articles (title, slug, content, category) VALUES
    ('How to upload a video', 'how-to-upload-video',
     'Click the upload button, select your video file, add title and description, then publish.',
     'video'),
    ('How to create a playlist', 'how-to-create-playlist',
     'Go to your library, click create playlist, give it a name and start adding videos.',
     'playlist'),
    ('Account Settings', 'account-settings',
     'Manage your account settings including profile, privacy, and notifications.',
     'account')
ON CONFLICT (slug) DO NOTHING;
