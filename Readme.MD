# AccessMesh - Advanced Access Control System

A robust, Go-based access control system designed for microservices architecture, implementing dynamic, fine-grained access control using MongoDB and Casbin. AccessMesh provides comprehensive policy management, rate limiting, and authentication features for modern distributed applications.

## Features

- Dynamic policy management via REST API
- Role-based access control (RBAC) with hierarchical support
- JWT-based authentication with robust middleware
- Rate limiting support for API endpoints
- MongoDB integration for policy and user data storage
- Context-aware authorization (IP, time-based)
- Audit logging for all policy evaluations
- Modular architecture with clean separation of concerns

## Prerequisites

- Go 1.19 or higher
- MongoDB 4.4 or higher
- Docker (optional, for containerized deployment)

## Installation

1. Clone the repository:

```bash
git clone https://github.com/nakulbharti/AccessMesh.git
cd AccessMesh
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
export PORT="8080"
```

4. Run the server:

```bash
go run cmd/server/main.go
```

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - Login and get JWT token

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

## Project Structure

```
├── cmd/
│   └── server/          # Application entry point
├── internal/
│   ├── api/            # API handlers and routes
│   │   ├── handlers/   # Request handlers
│   │   └── middleware/ # Authentication and rate limiting
│   ├── config/         # Configuration management
│   ├── models/         # Data models
│   └── store/          # Database interactions
├── pkg/                # Reusable packages
└── model.conf          # Casbin model configuration
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
docker build -t accessmesh .
docker run -p 8080:8080 accessmesh
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

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
