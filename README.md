# Ergo API Gateway

**Zero-Trust Security for your APIs, made simple.**

A lightweight API Gateway written in Go. It adds enterprise-grade security features—like JWT Authentication, RBAC, Threat Detection, and Rate Limiting—to any application without requiring code changes.

![License](https://img.shields.io/badge/license-Apache-2.0-blue.svg)
![Go Version](https://img.shields.io/badge/go-1.25+-00ADD8.svg)

## Features

*   **Zero-Trust Security**: Validates JWT tokens on every request.
*   **RBAC Policy Engine**: Fine-grained access control using Open Policy Agent (OPA) embedded directly in the binary.
*   **Threat Detection**: Automatically blocks suspicious activity (e.g., "Impossible Travel" login attempts).
*   **Rate Limiting**: Built-in sliding window limiter (Redis-backed or In-Memory fallback).
*   **Observability**: Prometheus metrics endpoint (`/metrics`) for real-time monitoring.
*   **Circuit Breaking**: Fails fast when your backend is down to prevent cascading failures.
*   **Single Binary**: Distributed as a single executable with no external dependencies required.

## Installation

### From Source
```bash
# Clone the repository
git clone https://github.com/cronenberg64/ergo-api.git
cd ergo-api

# Build the binary
go build -o ergo cmd/ergo/*.go

# Move to your path (optional)
sudo mv ergo /usr/local/bin/
```

### Docker
```bash
docker compose up --build -d
```

## Usage

### 1. Initialize Configuration
Run `init` to generate a default `.env` file in your current directory.
```bash
ergo init
```

### 2. Start the Gateway
```bash
ergo start
```
By default, Ergo listens on port `8080` and forwards traffic to `http://localhost:3000`.

### Configuration (.env)
| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | The port Ergo listens on. |
| `BACKEND_URL` | `http://localhost:3000` | The URL of your upstream service. |
| `JWT_SECRET` | *(Required)* | Secret key for validating JWT signatures. |
| `REDIS_ADDR` | `""` | Redis address (e.g., `localhost:6379`). If empty, uses In-Memory limiter. |
| `POLICY_PATH` | `""` | Path to custom OPA policy file. If empty, uses embedded RBAC policy. |

## Development

### Running Locally
```bash
# Run without building
go run cmd/ergo/*.go start
```

### Testing
```bash
go test ./...
```

## License
Apache-2.0