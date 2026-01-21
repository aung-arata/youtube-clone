# API Documentation

## Base URL
```
http://localhost:8080/api
```

## Endpoints

### Health Check

#### GET /health
Check if the API is running.

**Response:**
```json
{
  "status": "ok"
}
```

**Status Codes:**
- `200 OK` - API is healthy

---

### Videos

#### GET /videos
Retrieve a list of videos with optional search and pagination.

**Query Parameters:**
- `q` (optional): Search query to filter videos by title, description, or channel name
- `page` (optional): Page number for pagination (default: 1, min: 1)
- `limit` (optional): Number of videos per page (default: 20, max: 100)

**Example Request:**
```bash
# Get all videos
curl http://localhost:8080/api/videos

# Search for videos
curl "http://localhost:8080/api/videos?q=react"

# Get videos with pagination
curl "http://localhost:8080/api/videos?page=2&limit=10"

# Search with pagination
curl "http://localhost:8080/api/videos?q=tutorial&page=1&limit=5"
```

**Response:**
```json
[
  {
    "id": 1,
    "title": "Building a YouTube Clone",
    "description": "Learn how to build a YouTube clone...",
    "url": "https://example.com/video.mp4",
    "thumbnail": "https://example.com/thumb.jpg",
    "channel_name": "Code Master",
    "channel_avatar": "https://example.com/avatar.jpg",
    "views": 125000,
    "likes": 0,
    "dislikes": 0,
    "duration": "12:34",
    "uploaded_at": "2024-01-10T10:30:00Z",
    "created_at": "2024-01-10T10:30:00Z",
    "updated_at": "2024-01-10T10:30:00Z"
  }
]
```

**Status Codes:**
- `200 OK` - Videos retrieved successfully
- `500 Internal Server Error` - Database error

---

#### GET /videos/{id}
Retrieve a specific video by ID.

**Path Parameters:**
- `id` (required): Video ID

**Example Request:**
```bash
curl http://localhost:8080/api/videos/1
```

**Response:**
```json
{
  "id": 1,
  "title": "Building a YouTube Clone",
  "description": "Learn how to build a YouTube clone...",
  "url": "https://example.com/video.mp4",
  "thumbnail": "https://example.com/thumb.jpg",
  "channel_name": "Code Master",
  "channel_avatar": "https://example.com/avatar.jpg",
  "views": 125000,
  "likes": 0,
  "dislikes": 0,
  "duration": "12:34",
  "uploaded_at": "2024-01-10T10:30:00Z",
  "created_at": "2024-01-10T10:30:00Z",
  "updated_at": "2024-01-10T10:30:00Z"
}
```

**Status Codes:**
- `200 OK` - Video found
- `400 Bad Request` - Invalid video ID
- `404 Not Found` - Video not found
- `500 Internal Server Error` - Database error

---

#### POST /videos
Create a new video.

**Request Body:**
```json
{
  "title": "My New Video",
  "description": "Video description",
  "url": "https://example.com/video.mp4",
  "thumbnail": "https://example.com/thumb.jpg",
  "channel_name": "My Channel",
  "channel_avatar": "https://example.com/avatar.jpg",
  "duration": "10:30"
}
```

**Required Fields:**
- `title` - Video title (non-empty string)
- `url` - Video URL (non-empty string)
- `channel_name` - Channel name (non-empty string)

**Optional Fields:**
- `description` - Video description
- `thumbnail` - Thumbnail image URL
- `channel_avatar` - Channel avatar URL
- `duration` - Video duration (format: MM:SS or HH:MM:SS)

**Example Request:**
```bash
curl -X POST http://localhost:8080/api/videos \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My Video",
    "description": "Video description",
    "url": "https://example.com/video.mp4",
    "thumbnail": "https://example.com/thumb.jpg",
    "channel_name": "My Channel",
    "channel_avatar": "https://example.com/avatar.jpg",
    "duration": "10:30"
  }'
```

**Response:**
```json
{
  "id": 11,
  "title": "My New Video",
  "description": "Video description",
  "url": "https://example.com/video.mp4",
  "thumbnail": "https://example.com/thumb.jpg",
  "channel_name": "My Channel",
  "channel_avatar": "https://example.com/avatar.jpg",
  "views": 0,
  "likes": 0,
  "dislikes": 0,
  "duration": "10:30",
  "uploaded_at": "2024-01-12T04:30:00Z",
  "created_at": "2024-01-12T04:30:00Z",
  "updated_at": "2024-01-12T04:30:00Z"
}
```

**Status Codes:**
- `201 Created` - Video created successfully
- `400 Bad Request` - Invalid request body or missing required fields
- `500 Internal Server Error` - Database error

---

#### POST /videos/{id}/views
Increment the view count for a video.

**Path Parameters:**
- `id` (required): Video ID

**Example Request:**
```bash
curl -X POST http://localhost:8080/api/videos/1/views
```

**Response:**
```json
{
  "views": 125001
}
```

**Status Codes:**
- `200 OK` - Views incremented successfully
- `400 Bad Request` - Invalid video ID
- `404 Not Found` - Video not found
- `500 Internal Server Error` - Database error

---

#### POST /videos/{id}/like
Increment the like count for a video.

**Path Parameters:**
- `id` (required): Video ID

**Example Request:**
```bash
curl -X POST http://localhost:8080/api/videos/1/like
```

**Response:**
```json
{
  "likes": 1
}
```

**Status Codes:**
- `200 OK` - Likes incremented successfully
- `400 Bad Request` - Invalid video ID
- `404 Not Found` - Video not found
- `500 Internal Server Error` - Database error

---

#### POST /videos/{id}/dislike
Increment the dislike count for a video.

**Path Parameters:**
- `id` (required): Video ID

**Example Request:**
```bash
curl -X POST http://localhost:8080/api/videos/1/dislike
```

**Response:**
```json
{
  "dislikes": 1
}
```

**Status Codes:**
- `200 OK` - Dislikes incremented successfully
- `400 Bad Request` - Invalid video ID
- `404 Not Found` - Video not found
- `500 Internal Server Error` - Database error

---

### Comments

#### GET /videos/{videoId}/comments
Retrieve all comments for a specific video.

**Path Parameters:**
- `videoId` (required): Video ID

**Example Request:**
```bash
curl http://localhost:8080/api/videos/1/comments
```

**Response:**
```json
[
  {
    "id": 1,
    "video_id": 1,
    "user_id": 1,
    "content": "Great video!",
    "created_at": "2024-01-12T10:30:00Z",
    "updated_at": "2024-01-12T10:30:00Z"
  },
  {
    "id": 2,
    "video_id": 1,
    "user_id": 2,
    "content": "Thanks for sharing!",
    "created_at": "2024-01-12T11:00:00Z",
    "updated_at": "2024-01-12T11:00:00Z"
  }
]
```

**Status Codes:**
- `200 OK` - Comments retrieved successfully
- `400 Bad Request` - Invalid video ID
- `500 Internal Server Error` - Database error

---

#### POST /videos/{videoId}/comments
Create a new comment on a video.

**Path Parameters:**
- `videoId` (required): Video ID

**Request Body:**
```json
{
  "user_id": 1,
  "content": "Great video!"
}
```

**Required Fields:**
- `user_id` - User ID (positive integer)
- `content` - Comment content (non-empty string)

**Example Request:**
```bash
curl -X POST http://localhost:8080/api/videos/1/comments \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "content": "Great video!"
  }'
```

**Response:**
```json
{
  "id": 3,
  "video_id": 1,
  "user_id": 1,
  "content": "Great video!",
  "created_at": "2024-01-12T12:00:00Z",
  "updated_at": "2024-01-12T12:00:00Z"
}
```

**Status Codes:**
- `201 Created` - Comment created successfully
- `400 Bad Request` - Invalid request body or missing required fields
- `500 Internal Server Error` - Database error

---

#### GET /comments/{id}
Retrieve a specific comment by ID.

**Path Parameters:**
- `id` (required): Comment ID

**Example Request:**
```bash
curl http://localhost:8080/api/comments/1
```

**Response:**
```json
{
  "id": 1,
  "video_id": 1,
  "user_id": 1,
  "content": "Great video!",
  "created_at": "2024-01-12T10:30:00Z",
  "updated_at": "2024-01-12T10:30:00Z"
}
```

**Status Codes:**
- `200 OK` - Comment found
- `400 Bad Request` - Invalid comment ID
- `404 Not Found` - Comment not found
- `500 Internal Server Error` - Database error

---

#### PUT /comments/{id}
Update an existing comment.

**Path Parameters:**
- `id` (required): Comment ID

**Request Body:**
```json
{
  "content": "Updated comment!"
}
```

**Required Fields:**
- `content` - Comment content (non-empty string)

**Example Request:**
```bash
curl -X PUT http://localhost:8080/api/comments/1 \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Updated comment!"
  }'
```

**Response:**
```json
{
  "id": 1,
  "video_id": 1,
  "user_id": 1,
  "content": "Updated comment!",
  "created_at": "2024-01-12T10:30:00Z",
  "updated_at": "2024-01-12T12:30:00Z"
}
```

**Status Codes:**
- `200 OK` - Comment updated successfully
- `400 Bad Request` - Invalid request body or missing required fields
- `404 Not Found` - Comment not found
- `500 Internal Server Error` - Database error

---

#### DELETE /comments/{id}
Delete a comment.

**Path Parameters:**
- `id` (required): Comment ID

**Example Request:**
```bash
curl -X DELETE http://localhost:8080/api/comments/1
```

**Response:**
```
(Empty response body)
```

**Status Codes:**
- `204 No Content` - Comment deleted successfully
- `400 Bad Request` - Invalid comment ID
- `404 Not Found` - Comment not found
- `500 Internal Server Error` - Database error

---

## Rate Limiting

The API implements rate limiting to prevent abuse:
- **Limit:** 100 requests per minute per IP address
- **Response when exceeded:**
  ```
  HTTP/1.1 429 Too Many Requests
  Rate limit exceeded
  ```

## CORS

The API supports Cross-Origin Resource Sharing (CORS):
- **Allowed Origins:** `*` (all origins)
- **Allowed Methods:** GET, POST, PUT, DELETE, OPTIONS
- **Allowed Headers:** Content-Type, Authorization

## Error Responses

All error responses follow this format:

```
HTTP/1.1 <status_code>
<error_message>
```

Common error messages:
- `Invalid video ID` - Provided ID is not a valid number
- `Video not found` - No video exists with the given ID
- `Invalid comment ID` - Provided ID is not a valid number
- `Comment not found` - No comment exists with the given ID
- `Title is required` - Required field missing
- `URL is required` - Required field missing
- `Channel name is required` - Required field missing
- `Content is required` - Required field missing (for comments)
- `User ID is required` - Required field missing (for comments)
- `Rate limit exceeded` - Too many requests

## Best Practices

1. **Pagination:** Always use pagination for listing endpoints to improve performance
2. **Search:** Use specific search terms for better results
3. **Error Handling:** Always handle error responses appropriately
4. **Rate Limiting:** Implement exponential backoff when rate limit is exceeded
5. **Content-Type:** Always set `Content-Type: application/json` for POST requests

## Examples with Different Languages

### JavaScript (fetch)
```javascript
// Get all videos
const response = await fetch('http://localhost:8080/api/videos')
const videos = await response.json()

// Create a video
const response = await fetch('http://localhost:8080/api/videos', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    title: 'My Video',
    url: 'https://example.com/video.mp4',
    channel_name: 'My Channel',
  }),
})
const newVideo = await response.json()
```

### Python (requests)
```python
import requests

# Get all videos
response = requests.get('http://localhost:8080/api/videos')
videos = response.json()

# Create a video
video_data = {
    'title': 'My Video',
    'url': 'https://example.com/video.mp4',
    'channel_name': 'My Channel',
}
response = requests.post(
    'http://localhost:8080/api/videos',
    json=video_data
)
new_video = response.json()
```

### Go
```go
// Get all videos
resp, err := http.Get("http://localhost:8080/api/videos")
if err != nil {
    log.Fatal(err)
}
defer resp.Body.Close()

var videos []Video
json.NewDecoder(resp.Body).Decode(&videos)

// Create a video
videoData := Video{
    Title:       "My Video",
    URL:         "https://example.com/video.mp4",
    ChannelName: "My Channel",
}
body, _ := json.Marshal(videoData)
resp, err = http.Post(
    "http://localhost:8080/api/videos",
    "application/json",
    bytes.NewBuffer(body),
)
```

---

## Subscriptions

### POST /users/{userId}/subscriptions
Subscribe to a channel.

**Path Parameters:**
- `userId` (required): User ID

**Request Body:**
```json
{
  "channel_name": "Code Master"
}
```

**Example Request:**
```bash
curl -X POST http://localhost:8080/api/users/1/subscriptions \
  -H "Content-Type: application/json" \
  -d '{"channel_name":"Code Master"}'
```

**Response:**
```json
{
  "id": 1,
  "user_id": 1,
  "channel_name": "Code Master",
  "created_at": "2024-01-10T10:30:00Z"
}
```

**Status Codes:**
- `201 Created` - Subscribed successfully
- `400 Bad Request` - Invalid input
- `409 Conflict` - Already subscribed
- `500 Internal Server Error` - Database error

---

### DELETE /users/{userId}/subscriptions/{channelName}
Unsubscribe from a channel.

**Path Parameters:**
- `userId` (required): User ID
- `channelName` (required): Channel name

**Example Request:**
```bash
curl -X DELETE http://localhost:8080/api/users/1/subscriptions/Code%20Master
```

**Status Codes:**
- `204 No Content` - Unsubscribed successfully
- `404 Not Found` - Subscription not found
- `500 Internal Server Error` - Database error

---

### GET /users/{userId}/subscriptions
Get all subscriptions for a user.

**Path Parameters:**
- `userId` (required): User ID

**Example Request:**
```bash
curl http://localhost:8080/api/users/1/subscriptions
```

**Response:**
```json
[
  {
    "id": 1,
    "user_id": 1,
    "channel_name": "Code Master",
    "created_at": "2024-01-10T10:30:00Z"
  },
  {
    "id": 2,
    "user_id": 1,
    "channel_name": "Tech Tutorials",
    "created_at": "2024-01-11T14:20:00Z"
  }
]
```

**Status Codes:**
- `200 OK` - Subscriptions retrieved successfully
- `500 Internal Server Error` - Database error

---

### GET /users/{userId}/subscriptions/{channelName}
Check if user is subscribed to a channel.

**Path Parameters:**
- `userId` (required): User ID
- `channelName` (required): Channel name

**Example Request:**
```bash
curl http://localhost:8080/api/users/1/subscriptions/Code%20Master
```

**Response:**
```json
{
  "subscribed": true
}
```

**Status Codes:**
- `200 OK` - Status retrieved successfully
- `500 Internal Server Error` - Database error

---

## Playlists

### POST /users/{userId}/playlists
Create a new playlist.

**Path Parameters:**
- `userId` (required): User ID

**Request Body:**
```json
{
  "name": "My Favorite Videos",
  "description": "A collection of my favorite tutorials"
}
```

**Example Request:**
```bash
curl -X POST http://localhost:8080/api/users/1/playlists \
  -H "Content-Type: application/json" \
  -d '{"name":"My Favorite Videos","description":"A collection of my favorite tutorials"}'
```

**Response:**
```json
{
  "id": 1,
  "user_id": 1,
  "name": "My Favorite Videos",
  "description": "A collection of my favorite tutorials",
  "created_at": "2024-01-10T10:30:00Z",
  "updated_at": "2024-01-10T10:30:00Z"
}
```

**Status Codes:**
- `201 Created` - Playlist created successfully
- `400 Bad Request` - Invalid input
- `500 Internal Server Error` - Database error

---

### GET /users/{userId}/playlists
Get all playlists for a user.

**Path Parameters:**
- `userId` (required): User ID

**Example Request:**
```bash
curl http://localhost:8080/api/users/1/playlists
```

**Response:**
```json
[
  {
    "id": 1,
    "user_id": 1,
    "name": "My Favorite Videos",
    "description": "A collection of my favorite tutorials",
    "created_at": "2024-01-10T10:30:00Z",
    "updated_at": "2024-01-10T10:30:00Z"
  }
]
```

**Status Codes:**
- `200 OK` - Playlists retrieved successfully
- `500 Internal Server Error` - Database error

---

### GET /playlists/{id}
Get a specific playlist with its videos.

**Path Parameters:**
- `id` (required): Playlist ID

**Example Request:**
```bash
curl http://localhost:8080/api/playlists/1
```

**Response:**
```json
{
  "id": 1,
  "user_id": 1,
  "name": "My Favorite Videos",
  "description": "A collection of my favorite tutorials",
  "created_at": "2024-01-10T10:30:00Z",
  "updated_at": "2024-01-10T10:30:00Z",
  "videos": [
    {
      "id": 5,
      "title": "React Tutorial",
      "description": "Learn React basics",
      "url": "https://example.com/video.mp4",
      "thumbnail": "https://example.com/thumb.jpg",
      "channel_name": "Code Master",
      "channel_avatar": "https://example.com/avatar.jpg",
      "views": 10000,
      "likes": 500,
      "dislikes": 10,
      "category": "Education",
      "duration": "15:30",
      "uploaded_at": "2024-01-10T10:30:00Z",
      "created_at": "2024-01-10T10:30:00Z",
      "updated_at": "2024-01-10T10:30:00Z"
    }
  ]
}
```

**Status Codes:**
- `200 OK` - Playlist retrieved successfully
- `404 Not Found` - Playlist not found
- `500 Internal Server Error` - Database error

---

### PUT /playlists/{id}
Update a playlist.

**Path Parameters:**
- `id` (required): Playlist ID

**Request Body:**
```json
{
  "name": "Updated Playlist Name",
  "description": "Updated description"
}
```

**Example Request:**
```bash
curl -X PUT http://localhost:8080/api/playlists/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Updated Playlist Name","description":"Updated description"}'
```

**Response:**
```json
{
  "id": 1,
  "user_id": 1,
  "name": "Updated Playlist Name",
  "description": "Updated description",
  "created_at": "2024-01-10T10:30:00Z",
  "updated_at": "2024-01-10T12:00:00Z"
}
```

**Status Codes:**
- `200 OK` - Playlist updated successfully
- `400 Bad Request` - Invalid input
- `404 Not Found` - Playlist not found
- `500 Internal Server Error` - Database error

---

### DELETE /playlists/{id}
Delete a playlist.

**Path Parameters:**
- `id` (required): Playlist ID

**Example Request:**
```bash
curl -X DELETE http://localhost:8080/api/playlists/1
```

**Status Codes:**
- `204 No Content` - Playlist deleted successfully
- `404 Not Found` - Playlist not found
- `500 Internal Server Error` - Database error

---

### POST /playlists/{id}/videos
Add a video to a playlist.

**Path Parameters:**
- `id` (required): Playlist ID

**Request Body:**
```json
{
  "video_id": 5
}
```

**Example Request:**
```bash
curl -X POST http://localhost:8080/api/playlists/1/videos \
  -H "Content-Type: application/json" \
  -d '{"video_id":5}'
```

**Response:**
```json
{
  "id": 1,
  "playlist_id": 1,
  "video_id": 5,
  "position": 0,
  "added_at": "2024-01-10T10:30:00Z"
}
```

**Status Codes:**
- `201 Created` - Video added successfully
- `400 Bad Request` - Invalid input
- `409 Conflict` - Video already in playlist
- `500 Internal Server Error` - Database error

---

### DELETE /playlists/{id}/videos/{videoId}
Remove a video from a playlist.

**Path Parameters:**
- `id` (required): Playlist ID
- `videoId` (required): Video ID

**Example Request:**
```bash
curl -X DELETE http://localhost:8080/api/playlists/1/videos/5
```

**Status Codes:**
- `204 No Content` - Video removed successfully
- `404 Not Found` - Video not in playlist
- `500 Internal Server Error` - Database error

---

## Video Recommendations

### GET /videos/{id}/recommendations
Get recommended videos based on a given video.

**Path Parameters:**
- `id` (required): Video ID

**Query Parameters:**
- `limit` (optional): Number of recommendations (default: 10, max: 50)

**Example Request:**
```bash
curl http://localhost:8080/api/videos/1/recommendations

# With custom limit
curl "http://localhost:8080/api/videos/1/recommendations?limit=5"
```

**Response:**
```json
[
  {
    "id": 3,
    "title": "Advanced React Patterns",
    "description": "Learn advanced React patterns",
    "url": "https://example.com/video3.mp4",
    "thumbnail": "https://example.com/thumb3.jpg",
    "channel_name": "Code Master",
    "channel_avatar": "https://example.com/avatar.jpg",
    "views": 8000,
    "likes": 400,
    "dislikes": 5,
    "category": "Education",
    "duration": "20:15",
    "uploaded_at": "2024-01-12T10:30:00Z",
    "created_at": "2024-01-12T10:30:00Z",
    "updated_at": "2024-01-12T10:30:00Z"
  }
]
```

**Algorithm:**
- Returns videos from the same category
- Excludes the current video
- Sorted by views (descending) and upload date

**Status Codes:**
- `200 OK` - Recommendations retrieved successfully
- `404 Not Found` - Video not found
- `500 Internal Server Error` - Database error

