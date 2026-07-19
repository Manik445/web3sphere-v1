# Web3Sphere Backend

Enterprise-grade Go backend for the Web3Sphere platform. Built as a modular monolith, ready to scale and extract into microservices.

## Tech Stack
- **Go 1.24+**
- **Gin Web Framework**
- **PostgreSQL** & **GORM**
- **Redis**
- **RabbitMQ** & **Kafka**
- **Zap Logger**

## Project Structure
Adheres to Clean Architecture principles:
- `cmd/server` - Application entry point
- `configs/` - Configuration management
- `internal/` - Business modules (Auth, Users, Projects, Escrow, etc.)
- `pkg/` - Core infrastructure (Database, Redis, MQ, Mailer, Logger)
- `migrations/` - DB Schema and Seed data

## Local Development Workflow

### Prerequisites
- Go 1.24+
- Docker & Docker Compose
- `make`
- [golang-migrate](https://github.com/golang-migrate/migrate)

### 1. Environment Setup
```bash
cp .env.example .env
```
Fill out the required values in `.env` (the defaults work out-of-the-box for local docker compose).

### 2. Start Infrastructure
Start Postgres, Redis, RabbitMQ, Kafka, Mailhog via Docker Compose:
```bash
make docker-up
```

### 3. Database Migration & Seeding
```bash
make migrate-up
make seed
```

### 4. Run the Server
```bash
make run
```
Or for live-reloading (requires `air`):
```bash
go install github.com/air-verse/air@latest
make dev
```

## Available Services (Local)
- **API Server**: http://localhost:8080
- **Postgres**: localhost:5432
- **Redis**: localhost:6379
- **RabbitMQ Admin**: http://localhost:15672 (guest/guest)
- **Mailhog (Emails)**: http://localhost:8025
- **PgAdmin**: http://localhost:5050 (admin@web3sphere.io / admin)
