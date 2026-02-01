package docs

import (
	"encoding/json"
	"net/http"
)

// SwaggerInfo holds the API documentation information
type SwaggerInfo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
	BasePath    string `json:"basePath"`
}

// OpenAPISpec represents the OpenAPI 3.0 specification
type OpenAPISpec struct {
	OpenAPI    string              `json:"openapi"`
	Info       SwaggerInfo         `json:"info"`
	Servers    []Server            `json:"servers"`
	Paths      map[string]PathItem `json:"paths"`
	Components Components          `json:"components"`
	Tags       []Tag               `json:"tags"`
}

// Server represents a server in the OpenAPI spec
type Server struct {
	URL         string `json:"url"`
	Description string `json:"description"`
}

// Tag represents a tag in the OpenAPI spec
type Tag struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// PathItem represents a path item in the OpenAPI spec
type PathItem map[string]Operation

// Operation represents an operation in the OpenAPI spec
type Operation struct {
	Tags        []string              `json:"tags,omitempty"`
	Summary     string                `json:"summary"`
	Description string                `json:"description,omitempty"`
	OperationID string                `json:"operationId,omitempty"`
	Parameters  []Parameter           `json:"parameters,omitempty"`
	RequestBody *RequestBody          `json:"requestBody,omitempty"`
	Responses   map[string]Response   `json:"responses"`
	Security    []map[string][]string `json:"security,omitempty"`
}

// Parameter represents a parameter in the OpenAPI spec
type Parameter struct {
	Name        string  `json:"name"`
	In          string  `json:"in"`
	Description string  `json:"description,omitempty"`
	Required    bool    `json:"required,omitempty"`
	Schema      *Schema `json:"schema,omitempty"`
}

// RequestBody represents a request body in the OpenAPI spec
type RequestBody struct {
	Description string               `json:"description,omitempty"`
	Required    bool                 `json:"required,omitempty"`
	Content     map[string]MediaType `json:"content"`
}

// Response represents a response in the OpenAPI spec
type Response struct {
	Description string               `json:"description"`
	Content     map[string]MediaType `json:"content,omitempty"`
}

// MediaType represents a media type in the OpenAPI spec
type MediaType struct {
	Schema *Schema `json:"schema,omitempty"`
}

// Schema represents a schema in the OpenAPI spec
type Schema struct {
	Type        string             `json:"type,omitempty"`
	Format      string             `json:"format,omitempty"`
	Properties  map[string]*Schema `json:"properties,omitempty"`
	Items       *Schema            `json:"items,omitempty"`
	Ref         string             `json:"$ref,omitempty"`
	Description string             `json:"description,omitempty"`
	Example     interface{}        `json:"example,omitempty"`
}

// Components represents components in the OpenAPI spec
type Components struct {
	Schemas         map[string]*Schema        `json:"schemas,omitempty"`
	SecuritySchemes map[string]SecurityScheme `json:"securitySchemes,omitempty"`
}

// SecurityScheme represents a security scheme in the OpenAPI spec
type SecurityScheme struct {
	Type         string `json:"type"`
	Scheme       string `json:"scheme,omitempty"`
	BearerFormat string `json:"bearerFormat,omitempty"`
	Description  string `json:"description,omitempty"`
}

// GetOpenAPISpec returns the OpenAPI specification for the YouTube Clone API
func GetOpenAPISpec() *OpenAPISpec {
	return &OpenAPISpec{
		OpenAPI: "3.0.3",
		Info: SwaggerInfo{
			Title:       "YouTube Clone API",
			Description: "A comprehensive API for a YouTube clone application with video management, user authentication, comments, notifications, and more.",
			Version:     "1.0.0",
			BasePath:    "/api",
		},
		Servers: []Server{
			{URL: "http://localhost:8080", Description: "Development server"},
			{URL: "https://api.yourdomain.com", Description: "Production server"},
		},
		Tags: []Tag{
			{Name: "Authentication", Description: "User authentication operations"},
			{Name: "Videos", Description: "Video management operations"},
			{Name: "Comments", Description: "Comment management operations"},
			{Name: "Users", Description: "User profile operations"},
			{Name: "History", Description: "Watch history operations"},
			{Name: "Notifications", Description: "Notification management operations"},
			{Name: "Playlists", Description: "Playlist management operations"},
			{Name: "Subscriptions", Description: "Channel subscription operations"},
			{Name: "Plans", Description: "Subscription plan operations"},
		},
		Paths:      buildPaths(),
		Components: Components{
			Schemas:         buildSchemas(),
			SecuritySchemes: buildSecuritySchemes(),
		},
	}
}

func buildSecuritySchemes() map[string]SecurityScheme {
	return map[string]SecurityScheme{
		"bearerAuth": {
			Type:         "http",
			Scheme:       "bearer",
			BearerFormat: "JWT",
			Description:  "JWT Authorization header using the Bearer scheme",
		},
	}
}

func buildSchemas() map[string]*Schema {
	return map[string]*Schema{
		"Video": {
			Type: "object",
			Properties: map[string]*Schema{
				"id":             {Type: "integer", Description: "Video ID"},
				"title":          {Type: "string", Description: "Video title"},
				"description":    {Type: "string", Description: "Video description"},
				"url":            {Type: "string", Description: "Video URL"},
				"thumbnail":      {Type: "string", Description: "Thumbnail URL"},
				"channel_name":   {Type: "string", Description: "Channel name"},
				"channel_avatar": {Type: "string", Description: "Channel avatar URL"},
				"views":          {Type: "integer", Description: "View count"},
				"likes":          {Type: "integer", Description: "Like count"},
				"dislikes":       {Type: "integer", Description: "Dislike count"},
				"category":       {Type: "string", Description: "Video category"},
				"duration":       {Type: "string", Description: "Video duration"},
				"uploaded_at":    {Type: "string", Format: "date-time", Description: "Upload timestamp"},
			},
		},
		"User": {
			Type: "object",
			Properties: map[string]*Schema{
				"id":         {Type: "integer", Description: "User ID"},
				"username":   {Type: "string", Description: "Username"},
				"email":      {Type: "string", Format: "email", Description: "Email address"},
				"avatar":     {Type: "string", Description: "Avatar URL"},
				"role":       {Type: "string", Description: "User role"},
				"created_at": {Type: "string", Format: "date-time", Description: "Account creation timestamp"},
			},
		},
		"Comment": {
			Type: "object",
			Properties: map[string]*Schema{
				"id":         {Type: "integer", Description: "Comment ID"},
				"video_id":   {Type: "integer", Description: "Video ID"},
				"user_id":    {Type: "integer", Description: "User ID"},
				"content":    {Type: "string", Description: "Comment content"},
				"created_at": {Type: "string", Format: "date-time", Description: "Creation timestamp"},
				"updated_at": {Type: "string", Format: "date-time", Description: "Last update timestamp"},
			},
		},
		"Notification": {
			Type: "object",
			Properties: map[string]*Schema{
				"id":         {Type: "integer", Description: "Notification ID"},
				"user_id":    {Type: "integer", Description: "User ID"},
				"type":       {Type: "string", Description: "Notification type"},
				"title":      {Type: "string", Description: "Notification title"},
				"message":    {Type: "string", Description: "Notification message"},
				"link":       {Type: "string", Description: "Related link"},
				"is_read":    {Type: "boolean", Description: "Read status"},
				"created_at": {Type: "string", Format: "date-time", Description: "Creation timestamp"},
			},
		},
		"Playlist": {
			Type: "object",
			Properties: map[string]*Schema{
				"id":          {Type: "integer", Description: "Playlist ID"},
				"user_id":     {Type: "integer", Description: "Owner user ID"},
				"name":        {Type: "string", Description: "Playlist name"},
				"description": {Type: "string", Description: "Playlist description"},
				"created_at":  {Type: "string", Format: "date-time", Description: "Creation timestamp"},
			},
		},
		"Plan": {
			Type: "object",
			Properties: map[string]*Schema{
				"id":                    {Type: "integer", Description: "Plan ID"},
				"name":                  {Type: "string", Description: "Plan name"},
				"price":                 {Type: "number", Format: "float", Description: "Price"},
				"max_video_quality":     {Type: "string", Description: "Maximum video quality"},
				"max_uploads_per_month": {Type: "integer", Description: "Maximum uploads per month"},
				"ads_free":              {Type: "boolean", Description: "Ad-free experience"},
			},
		},
		"AuthResponse": {
			Type: "object",
			Properties: map[string]*Schema{
				"token": {Type: "string", Description: "JWT token"},
				"user":  {Ref: "#/components/schemas/User"},
			},
		},
		"Error": {
			Type: "object",
			Properties: map[string]*Schema{
				"error":   {Type: "string", Description: "Error message"},
				"code":    {Type: "integer", Description: "Error code"},
				"details": {Type: "string", Description: "Error details"},
			},
		},
	}
}

func buildPaths() map[string]PathItem {
	return map[string]PathItem{
		"/api/auth/signup": {
			"post": Operation{
				Tags:    []string{"Authentication"},
				Summary: "Register a new user",
				RequestBody: &RequestBody{
					Required: true,
					Content: map[string]MediaType{
						"application/json": {
							Schema: &Schema{
								Type: "object",
								Properties: map[string]*Schema{
									"username": {Type: "string"},
									"email":    {Type: "string", Format: "email"},
									"password": {Type: "string", Format: "password"},
									"avatar":   {Type: "string"},
								},
							},
						},
					},
				},
				Responses: map[string]Response{
					"201": {Description: "User created successfully", Content: map[string]MediaType{"application/json": {Schema: &Schema{Ref: "#/components/schemas/AuthResponse"}}}},
					"400": {Description: "Bad request"},
					"409": {Description: "User already exists"},
				},
			},
		},
		"/api/auth/login": {
			"post": Operation{
				Tags:    []string{"Authentication"},
				Summary: "Login user",
				RequestBody: &RequestBody{
					Required: true,
					Content: map[string]MediaType{
						"application/json": {
							Schema: &Schema{
								Type: "object",
								Properties: map[string]*Schema{
									"email":    {Type: "string", Format: "email"},
									"password": {Type: "string", Format: "password"},
								},
							},
						},
					},
				},
				Responses: map[string]Response{
					"200": {Description: "Login successful", Content: map[string]MediaType{"application/json": {Schema: &Schema{Ref: "#/components/schemas/AuthResponse"}}}},
					"401": {Description: "Invalid credentials"},
				},
			},
		},
		"/api/auth/me": {
			"get": Operation{
				Tags:     []string{"Authentication"},
				Summary:  "Get current user",
				Security: []map[string][]string{{"bearerAuth": {}}},
				Responses: map[string]Response{
					"200": {Description: "Current user", Content: map[string]MediaType{"application/json": {Schema: &Schema{Ref: "#/components/schemas/User"}}}},
					"401": {Description: "Unauthorized"},
				},
			},
		},
		"/api/videos": {
			"get": Operation{
				Tags:    []string{"Videos"},
				Summary: "Get all videos",
				Parameters: []Parameter{
					{Name: "q", In: "query", Description: "Search query", Schema: &Schema{Type: "string"}},
					{Name: "category", In: "query", Description: "Filter by category", Schema: &Schema{Type: "string"}},
					{Name: "page", In: "query", Description: "Page number", Schema: &Schema{Type: "integer"}},
					{Name: "limit", In: "query", Description: "Items per page", Schema: &Schema{Type: "integer"}},
					{Name: "sort_by", In: "query", Description: "Sort field (views, likes, date, title)", Schema: &Schema{Type: "string"}},
					{Name: "order", In: "query", Description: "Sort order (asc, desc)", Schema: &Schema{Type: "string"}},
				},
				Responses: map[string]Response{
					"200": {Description: "List of videos", Content: map[string]MediaType{"application/json": {Schema: &Schema{Type: "array", Items: &Schema{Ref: "#/components/schemas/Video"}}}}},
				},
			},
			"post": Operation{
				Tags:    []string{"Videos"},
				Summary: "Create a new video",
				RequestBody: &RequestBody{
					Required: true,
					Content: map[string]MediaType{
						"application/json": {
							Schema: &Schema{Ref: "#/components/schemas/Video"},
						},
					},
				},
				Responses: map[string]Response{
					"201": {Description: "Video created", Content: map[string]MediaType{"application/json": {Schema: &Schema{Ref: "#/components/schemas/Video"}}}},
					"400": {Description: "Bad request"},
				},
			},
		},
		"/api/videos/{id}": {
			"get": Operation{
				Tags:    []string{"Videos"},
				Summary: "Get a specific video",
				Parameters: []Parameter{
					{Name: "id", In: "path", Required: true, Description: "Video ID", Schema: &Schema{Type: "integer"}},
				},
				Responses: map[string]Response{
					"200": {Description: "Video details", Content: map[string]MediaType{"application/json": {Schema: &Schema{Ref: "#/components/schemas/Video"}}}},
					"404": {Description: "Video not found"},
				},
			},
		},
		"/api/videos/{id}/views": {
			"post": Operation{
				Tags:    []string{"Videos"},
				Summary: "Increment video view count",
				Parameters: []Parameter{
					{Name: "id", In: "path", Required: true, Description: "Video ID", Schema: &Schema{Type: "integer"}},
				},
				Responses: map[string]Response{
					"200": {Description: "View count incremented"},
					"404": {Description: "Video not found"},
				},
			},
		},
		"/api/videos/{id}/like": {
			"post": Operation{
				Tags:    []string{"Videos"},
				Summary: "Like a video",
				Parameters: []Parameter{
					{Name: "id", In: "path", Required: true, Description: "Video ID", Schema: &Schema{Type: "integer"}},
				},
				Responses: map[string]Response{
					"200": {Description: "Video liked"},
					"404": {Description: "Video not found"},
				},
			},
		},
		"/api/videos/{videoId}/comments": {
			"get": Operation{
				Tags:    []string{"Comments"},
				Summary: "Get comments for a video",
				Parameters: []Parameter{
					{Name: "videoId", In: "path", Required: true, Description: "Video ID", Schema: &Schema{Type: "integer"}},
				},
				Responses: map[string]Response{
					"200": {Description: "List of comments", Content: map[string]MediaType{"application/json": {Schema: &Schema{Type: "array", Items: &Schema{Ref: "#/components/schemas/Comment"}}}}},
				},
			},
			"post": Operation{
				Tags:    []string{"Comments"},
				Summary: "Create a comment",
				Parameters: []Parameter{
					{Name: "videoId", In: "path", Required: true, Description: "Video ID", Schema: &Schema{Type: "integer"}},
				},
				RequestBody: &RequestBody{
					Required: true,
					Content: map[string]MediaType{
						"application/json": {
							Schema: &Schema{
								Type: "object",
								Properties: map[string]*Schema{
									"user_id": {Type: "integer"},
									"content": {Type: "string"},
								},
							},
						},
					},
				},
				Responses: map[string]Response{
					"201": {Description: "Comment created", Content: map[string]MediaType{"application/json": {Schema: &Schema{Ref: "#/components/schemas/Comment"}}}},
					"400": {Description: "Bad request"},
				},
			},
		},
		"/api/users/{userId}/notifications": {
			"get": Operation{
				Tags:    []string{"Notifications"},
				Summary: "Get user notifications",
				Parameters: []Parameter{
					{Name: "userId", In: "path", Required: true, Description: "User ID", Schema: &Schema{Type: "integer"}},
					{Name: "unread", In: "query", Description: "Filter unread only", Schema: &Schema{Type: "boolean"}},
					{Name: "limit", In: "query", Description: "Limit results", Schema: &Schema{Type: "integer"}},
				},
				Responses: map[string]Response{
					"200": {Description: "List of notifications", Content: map[string]MediaType{"application/json": {Schema: &Schema{Type: "array", Items: &Schema{Ref: "#/components/schemas/Notification"}}}}},
				},
			},
		},
		"/api/notifications": {
			"post": Operation{
				Tags:    []string{"Notifications"},
				Summary: "Create a notification",
				RequestBody: &RequestBody{
					Required: true,
					Content: map[string]MediaType{
						"application/json": {
							Schema: &Schema{Ref: "#/components/schemas/Notification"},
						},
					},
				},
				Responses: map[string]Response{
					"201": {Description: "Notification created"},
					"400": {Description: "Bad request"},
				},
			},
		},
		"/api/health": {
			"get": Operation{
				Tags:    []string{"System"},
				Summary: "Health check",
				Responses: map[string]Response{
					"200": {Description: "Service is healthy"},
				},
			},
		},
	}
}

// SwaggerUIHandler serves the Swagger UI HTML page
func SwaggerUIHandler(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>YouTube Clone API Documentation</title>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script>
        window.onload = function() {
            SwaggerUIBundle({
                url: "/api/docs/openapi.json",
                dom_id: '#swagger-ui',
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIBundle.SwaggerUIStandalonePreset
                ],
                layout: "StandaloneLayout"
            });
        };
    </script>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

// OpenAPISpecHandler serves the OpenAPI JSON specification
func OpenAPISpecHandler(w http.ResponseWriter, r *http.Request) {
	spec := GetOpenAPISpec()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(spec)
}
