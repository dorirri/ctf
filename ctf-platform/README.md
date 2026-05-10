# CTF Platform

A Jeopardy-style CTF platform built with Go, PostgreSQL, and React.

## Prerequisites

- Docker 24+
- Docker Compose v2 (comes with Docker Desktop)

## Quickstart

```bash
docker compose up --build
```

Frontend runs at http://localhost:3000, API at http://localhost:8080.

Default development admin account: `admin@ctf.local` / `admin123`

## Configuration

The backend reads:

- `DATABASE_URL` or `DB_URL` - PostgreSQL connection string
- `JWT_SECRET` - token signing secret
- `APP_ENV` - defaults to `development`
- `ADMIN_USERNAME`, `ADMIN_EMAIL`, `ADMIN_PASSWORD` - optional startup admin seed

Outside `development`, the backend refuses to start with the default `JWT_SECRET=changeme`.

## Local development

Backend:

```bash
cd backend
go run .
```

Frontend:

```bash
cd frontend
npm install
npm run dev
```

Useful checks:

```bash
cd backend && go test ./...
cd frontend && npm run build
```

Healthcheck: http://localhost:8080/api/health

## Seeding

Docker Compose creates the development admin from the `ADMIN_*` environment variables.
After signing in as an admin, call `POST /api/admin/seed` to add sample challenges.

