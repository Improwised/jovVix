# Getting Started with jovVix

This guide walks you through setting up jovVix for local development.

## Prerequisites

| Tool | Version | Install |
|------|---------|---------|
| Docker | Latest | [docker.com](https://docs.docker.com/get-docker/) |
| Docker Compose | v2+ | Included with Docker Desktop |
| Node.js | v18.0+ | [nodejs.org](https://nodejs.org/) |
| Go | v1.21+ | [golang.org](https://golang.org/dl/) |
| Git | Latest | [git-scm.com](https://git-scm.com/) |

## Quick Start (Docker — Recommended)

```bash
# 1. Clone the repository
git clone https://github.com/Improwised/jovVix.git
cd jovVix

# 2. Start all services
docker-compose up --build

# 3. Open in browser:
#    App:      http://127.0.0.1:5000
#    API:      http://127.0.0.1:3000/api/v1/docs
#    Mailpit:  http://127.0.0.1:8025
```

### Docker Services

| Service | Port | Purpose |
|---------|------|---------|
| `web` | 5000 | Nuxt frontend |
| `api` | 3000 | Go backend API |
| `db` | 5432 | PostgreSQL database |
| `redis` | 6379 | Valkey/Redis cache |
| `kratos` | 4433 | Ory Kratos auth (public) |
| `kratos` | 4434 | Ory Kratos admin API |
| `mailpit` | 8025 | Email testing UI |
| `mailpit` | 1025 | SMTP server |
| `minio` | 9000 | Object storage |

## Local Development (Without Docker)

### 1. Start Backend Services

```bash
cd api
cp .env.example .env
docker-compose up -d
```

### 2. Set Up Backend

```bash
cd api
go mod tidy && go mod vendor
go run app.go migrate up
go run app.go api
```

### 3. Set Up Frontend

```bash
cd app
cp .env.example .env
npm install
npm run dev
```

## Essential Variables to Change

| Variable | File | Default | Change To |
|----------|------|---------|-----------|
| `JWT_SECRET` | `api/.env` | `ThisIsKey` | A random secret string |
| `DB_PASSWORD` | `api/.env` | `jovvix` | A strong password |
| `SMTP_HOST` | `api/.env` | `your-ip-address` | Your SMTP server |
| `EMAIL_FROM` | `api/.env` | `example@gmail.com` | Your email address |

## Verifying the Setup

1. Open `http://127.0.0.1:5000` — jovVix homepage
2. Open `http://127.0.0.1:3000/api/v1/docs` — Swagger API docs
3. Open `http://127.0.0.1:8025` — Mailpit email UI
4. Try registering (check Mailpit for verification email)

## Running Tests

```bash
cd app && npm run test    # Frontend tests
cd api && go test ./...   # Backend tests
cd app && npm run lint    # Frontend linting
cd api && golangci-lint run  # Backend linting
```

## Troubleshooting

### Port Already in Use
```bash
lsof -i :5000
kill <PID>
```

### Database Connection Refused
```bash
docker-compose ps db
docker-compose logs db
```

### Migration Errors
```bash
cd api
go run app.go migrate down
go run app.go migrate up
```

### Frontend Build Errors
```bash
cd app
rm -rf node_modules .nuxt .output
npm install
npm run dev
```

## Next Steps

- [Architecture Overview](./architecture.md)
- [API Development Guide](./api-development.md)
- [Frontend Development Guide](./frontend-development.md)
- [Coding Standards](./coding-standards.md)
