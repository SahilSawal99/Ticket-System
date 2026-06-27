# Ticket System API
A lightweight Go backend for managing tickets with JWT authentication and an in-memory data store.

## Features
- User registration and login
- JWT-based authentication
- Ticket creation, listing, retrieval, and status updates
- Simple in-memory persistence for quick local development
- Docker-ready build for `linux/amd64`

## Endpoints
- `GET /health` — service health check
- `POST /auth/register` — register a new user
- `POST /auth/login` — login and receive a JWT
- `POST /tickets` — create a new ticket (requires Bearer token)
- `GET /tickets` — list tickets for the logged-in user (requires Bearer token)
- `GET /tickets/{id}` — get a single ticket by ID (requires Bearer token)
- `PATCH /tickets/{id}/status` — update a ticket status (requires Bearer token)

## Local Development
1. Install Go 1.25.
2. From the project root:

```bash
go mod tidy
go run main.go
```

3.Open your browser or API client at `http://localhost:8080`.

## Deployment
The application is deployed at:

- `https://ticket-system-etst.onrender.com`

## Docker
Build the Docker image for `linux/amd64` & tag it:

```bash
docker build --platform linux/amd64 -t ticket-system:latest .
```

Run the container locally:

```bash
docker run -d -p 8080:8080 --name ticket-service ticket-system:latest
```

## Important Point
- This service uses an in-memory data store. All data is lost when the app stops.
- No database or persistent storage is configured.
- The JWT signing key is currently hardcoded in `internal/auth/auth.go`


## Project Structure
- `main.go` — application startup and route registration
- `internal/auth` — JWT generation and request authentication
- `internal/controller` — HTTP handlers and request/response logic
- `internal/service` — business logic and validation
- `internal/repository` — in-memory data store operations
- `internal/model` — shared data model definition