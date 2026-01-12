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
- `Title is required` - Required field missing
- `URL is required` - Required field missing
- `Channel name is required` - Required field missing
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
