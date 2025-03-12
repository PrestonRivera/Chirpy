# Chirpy

Chirpy is a social media platform project inspired by Twitter/X, built as a REST API to learn API fundamentals.
With Docker integration for easy setup and deployment.

## Current Features

- User authentication (signup and login)
- Create and delete posts ("chirps")
- View all chirps

## Technologies Used

- Go (backend)
- PostgreSQL (database)
- HTML/CSS/JavaScript (frontend)
- Docker (containerization)

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)

## Running with Docker

The application can be easily run using Docker compose:

1. Clone the repository

```git clone https://github.com/PrestonRivera/Chirpy.git```

2. Navigate to the root of the repo

```cd Chirpy```

3. Start the application with Docker Compose

```docker-compose up```

## The server will be available at http://localhost:8080

## Database Initialization

The PostgreSQL database is automatically initialized with the necessary schema using:
- The init.sql script (for initial database setup)
- Goose migrations (for schema versioning)

No manual database setup is required.

### API Highlights

- Authentication using JWT tokens
- Create and manage chirps
- User account management

### Development Status

This project is in active developement. The following features are planned for future implementation:

- User profiles
- Post comments
- Like/unlike functionality
- User following
- Image uploads
- and MORE!

## API Reference
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