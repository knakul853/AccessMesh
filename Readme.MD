# Policy-Based Access Control (PBAC) System

A Go-based PBAC system for microservices implementing dynamic, fine-grained access control with MongoDB and Casbin.

## Features

- Dynamic policy management via REST API
- Role-based access control with hierarchical support
- JWT-based authentication
- MongoDB integration for policy storage
- Context-aware authorization (IP, time-based)
- Audit logging for all policy evaluations

## Prerequisites

- Go 1.19 or higher
- MongoDB 4.4 or higher
- Docker (optional)

## Installation

1. Clone the repository:

```bash
git clone https://github.com/knakul853/accessmesh
cd pbac
```

2. Install dependencies:

```bash
go mod tidy
```

3. Set up environment variables:

```bash
export MONGO_URI="mongodb://localhost:27017"
export JWT_SECRET="your-secret-key"
export ENV="development"
```

4. Run the server:

```bash
go run cmd/server/main.go
```

## API Endpoints

### Policies

- `POST /api/v1/policies` - Create a new policy
- `GET /api/v1/policies` - List all policies
- `GET /api/v1/policies/{id}` - Get a specific policy
- `PUT /api/v1/policies/{id}` - Update a policy
- `DELETE /api/v1/policies/{id}` - Delete a policy

### Example Policy

```json
{
  "role": "manager",
  "resource": "/api/v1/orders",
  "action": "read",
  "conditions": {
    "ip_range": ["10.0.0.0/16"],
    "time_range": ["08:00-20:00"]
  }
}
```

## Docker Deployment

```dockerfile
FROM golang:1.19-alpine

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main cmd/server/main.go

EXPOSE 8080
CMD ["./main"]
```

Build and run:

```bash
docker build -t pbac-service .
docker run -p 8080:8080 pbac-service
```

## Architecture

```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│   Client    │────▶│  API Server  │────▶│   MongoDB   │
└─────────────┘     └──────────────┘     └─────────────┘
                           │
                    ┌──────┴───────┐
                    │    Casbin    │
                    └──────────────┘
```
