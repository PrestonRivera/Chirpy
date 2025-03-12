# Chirpy

Chirpy is a social media platform project inspired by Twitter/X, built as a REST API to learn API fundamentals.
With Docker integration for easy setup and deployment.

## Current Features

- User authentication (signup and login)
- Create and delete posts ("chirps")
- View all chirps

## Running with Docker

The application can be easily run using Docker compose:

1. Clone the repository

```git clone https://github.com/PrestonRivera/Chirpy.git```

2. Navigate to the root of the repo

```cd Chirpy```

3. Start the application with Docker Compose

```docker-compose up```

## The server will be available at http://localhost:8080

### API Highlights

- Authentication using JWT tokens
- Create and manage chirps
- User account management

### Technologies Used

- Go (backend)
- PostgreSQL (database)
- HTML/CSS/JavaScript (frontend)
- Docker (containerization)

### Development Status

This project is in active developement. The following features are planned for future implementation:

- User profiles
- Post comments
- Like/unlike functionality
- User following
- Image uploads
- and MORE!