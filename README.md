# Chirpy A version of the famous Twitter/X

A REST api I built to learn the fundamentals of an API.

## Getting Started 

Base URL: `http://localhost:8080`

NOTE: This API runs locally on port 8080. Make sure the server is running before making requests. I used AIR to build/run the server and Postman to test my endpoints.

## Running the Server 

1. Must have GO installed on your system
2. Clone the repo
3. Navigate to the project directory
4. Run the server:
``` go run main.go```

## Status Codes 

This API returns standard HTTP status codes:
- 200: Success
- 404: Not found
- 500: Server error

## Authentication 
Many endpoints require JWT authentication via Bearer token. Include the token in the authorization header:
```
Authorization: Bearer <your_jwt_token>
```

Protected endpoints:
- All POST /api/chirps endpoints
- DELETE /api/chirps/{chirpID}
- PUT /api/users

## Query Parameters
### GET /api/chirps
- sort: (optional) Sort chirps by creation date
    - asc: Ascending order
    - desc: Descending order
- author_id: (optional) Filter chirps by author ID

## Endpoints

### Health & Metrics
- GET `/api/healthz` - Check API health status
- GET `/admin/metrics` - Get API metrics

### Chirps
- GET `/api/chirps` - Get all chirps
- GET `/api/chirps/{chirpID}` - Get a specific chirp
- POST `/api/chirps` - Create a new chirp
- DELETE `/api/chirps/{chirpID}` - Delete a specific chirp

### Users
- POST `/api/users` - Create a new user
- PUT `/api/users` - Update user information
- POST `/api/login` - User login
- POST `/api/refresh` - Refresh authentication token
- POST `/api/revoke` - Revoke authentication token

### Admin
- POST `/admin/reset` - Reset users database
- POST `/admin/reset-fileserver-hits` - Reset hit counter
- POST `/polka/webhooks` - Upgrade user status

### Example Requests/Responses

```json
// POST /api/users - Create a user
{
  "password": "04234",
  "email": "preston@example.com"
}

// Response
{
    "id": "f7960cf2-898c-4290-b722-d7c6ff194285",
    "created_at": "2025-01-31T00:51:55.043245Z",
    "updated_at": "2025-01-31T00:51:55.043245Z",
    "email": "preston@example.com",
    "is_chirpy_red": false
}

// POST /api/login - login as that user
{
    "password": "04234",
    "email": "preston@example.com"
}

// Response
{
    "id": "f7960cf2-898c-4290-b722-d7c6ff194285",
    "created_at": "2025-01-31T00:51:55.043245Z",
    "updated_at": "2025-01-31T00:51:55.043245Z",
    "email": "preston@example.com",
    "is_chirpy_red": false,
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHktYWNjZXNzIiwic3ViIjoiZjc5NjBjZjItODk4Yy00MjkwLWI3MjItZDdjNmZmMTk0Mjg1IiwiZXhwIjoxNzM4MzEwMDI4LCJpYXQiOjE3MzgzMDY0Mjh9.18LL-SHykTrHLzyH7SyV8qPf-NuPghUkaCsdZCJ-H_U",
    "refresh_token": "98286e1f29574a6e3f73e351ae52b0ee9f8d17d6aa163994543e58c31cc7a0d9"
}

// POST /api/chirps - Create a Chirp
Headers:
{
  "Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHktYWNjZXNzIiwic3ViIjoiZjc5NjBjZjItODk4Yy00MjkwLWI3MjItZDdjNmZmMTk0Mjg1IiwiZXhwIjoxNzM4MzEwMDI4LCJpYXQiOjE3MzgzMDY0Mjh9.18LL-SHykTrHLzyH7SyV8qPf-NuPghUkaCsdZCJ-H_U"
}
Request Body:
{
  "body": "I am the one who knocks"
}

// Response 
{
    "id": "21365d86-21ae-4a8e-97cd-fdfb5bf3e3ca",
    "created_at": "2025-01-31T01:00:06.20064Z",
    "updated_at": "2025-01-31T01:00:06.20064Z",
    "body": "I am the one who knocks",
    "user_id": "f7960cf2-898c-4290-b722-d7c6ff194285"
}

// GET /api/chirps
// Response
[
    {
        "id": "21365d86-21ae-4a8e-97cd-fdfb5bf3e3ca",
        "created_at": "2025-01-31T01:00:06.20064Z",
        "updated_at": "2025-01-31T01:00:06.20064Z",
        "body": "I am the one who knocks",
        "user_id": "f7960cf2-898c-4290-b722-d7c6ff194285"
    }
]
```
