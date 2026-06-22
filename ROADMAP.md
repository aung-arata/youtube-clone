# YouTube Clone — Roadmap

## Legend
- ✅ Done
- 🚧 In Progress
- ⬜ Planned

---

## Phase 0 — Critical Bug Fixes ✅
- ✅ Fix DB password hardcoding → env vars
- ✅ Fix duplicate package declarations in test files
- ✅ Remove compiled binaries from git / update `.gitignore`
- ✅ Remove `go.work` — services are independent modules

## Phase 1 — Backend Safety & Hardening ✅
- ✅ Atomic video status promotion (no partial-ready state)
- ✅ Stricter metadata validation on upload
- ✅ Block video player until transcoding complete
- ✅ Auth handler security fixes (user-service)
- ✅ Hub/WebSocket test hardening
- ✅ Unit tests across all Go services

## Phase 2 — Frontend Navigation Wiring ✅
- ✅ Sidebar navigation wired to routes
- ✅ `/trending` — TrendingPage fetches from video-service
- ✅ `/subscriptions` — SubscriptionsPage fetches from user-service
- ✅ `/notifications` — NotificationsPage fetches from notification-service
- ✅ `/history` — WatchHistory component fetches from history-service
- ✅ `/playlists` — PlaylistsPage + PlaylistDetailPage
- ✅ `/profile/:userId` — UserProfile component
- ✅ CI pipeline fix (frontend leaked into Go services matrix)
- ✅ Node.js upgraded to 20 (Vite 7 requirement)

## Phase 3 — Backend API Completeness ✅
- ✅ `GET /api/videos/trending` — video-service + gateway
- ✅ `GET /api/users/:id/subscriptions` — user-service + gateway
- ✅ `POST/DELETE /api/users/:id/subscriptions/:channel` — subscribe/unsubscribe
- ✅ `GET /api/users/:id/notifications` — notification-service + gateway
- ✅ `POST /api/users/:id/notifications/read` — mark notifications read
- ✅ `GET /api/users/:id/history` — history-service + gateway
- ✅ `GET /api/playlists/:id` — video-service + gateway
- ✅ API gateway routing for all new endpoints

## Phase 4 — Search & Discovery ⬜
- ⬜ Search bar in Header wired to backend
- ⬜ `GET /api/videos/search?q=` endpoint
- ⬜ Search results page
- ⬜ Category/tag filtering

## Phase 5 — Auth & User Experience Polish ⬜
- ⬜ Subscribe button on VideoPage wired to API
- ⬜ Like/dislike on VideoPage wired to API
- ⬜ Comment submit wired to API
- ⬜ User avatar upload
- ⬜ Persistent auth (refresh tokens / remember me)

## Phase 6 — Admin & Observability ⬜
- ⬜ PHP admin service integration with frontend
- ⬜ Admin dashboard: user management, video moderation
- ⬜ Structured logging across all Go services
- ⬜ Health check endpoints standardized
- ⬜ Docker Compose production config

## Phase 7 — Production Readiness ⬜
- ⬜ HTTPS / TLS termination
- ⬜ Rate limiting on API gateway
- ⬜ CDN / object storage for video files (S3-compatible)
- ⬜ Database migrations versioned and automated
- ⬜ CI/CD: automated deploy on merge to main
- ⬜ E2E tests (Playwright or Cypress)
