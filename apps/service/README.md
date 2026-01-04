# Service

Go API service for Weiss built with Fiber, GORM, and PostgreSQL.

## Requirements

- Go 1.25+
- PostgreSQL 16+

## Quickstart

1. Copy `.env.example` to `.env` and configure
2. Install dependencies: `just deps`
3. Run: `just run` or `just dev` (with hot reload)

## Commands

```bash
just run          # Run service
just dev          # Run with hot reload (air)
just build        # Build binary
just test         # Run tests
just deps         # Download dependencies

just docker-up    # Start Docker environment
just docker-down  # Stop Docker environment
just docker-logs # View logs
```

## Migrations

Database migrations are managed using [goose](https://github.com/pressly/goose).

### Usage

```bash
# Create a new migration
just migration-create add_users_table

# Run migrations
just migration-up

# Rollback last migration
just migration-down

# Check migration status
just migration-status
```

## Migration Files

Migration files follow the naming pattern: `YYYYMMDDHHMMSS_name.sql`

- Up migrations: SQL statements to apply changes
- Down migrations: SQL statements to revert changes



## Docker

Development: `docker-compose up`  
Production: `docker build -t weiss-service:latest -f Dockerfile .`
