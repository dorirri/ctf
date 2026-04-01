# CTF Platform

A Jeopardy-style Capture The Flag platform built with Go + PostgreSQL + React.

## Quick Start

```bash
# copy and edit secrets before first run
cp .env.example .env

docker compose up --build
```

The API will be available at `http://localhost:8080`.

## Environment Variables

| Variable     | Default                                        | Description          |
|--------------|------------------------------------------------|----------------------|
| `PORT`       | `8080`                                         | HTTP listen port     |
| `DB_URL`     | `postgres://ctf:ctf@db:5432/ctfdb?sslmode=disable` | Postgres DSN     |
| `JWT_SECRET` | `changeme`                                     | HS256 signing secret |

Set `JWT_SECRET` to a long random string in production:

```bash
openssl rand -hex 32
```

## Creating the First Admin

After startup, promote a registered user directly in the database:

```bash
docker compose exec db psql -U ctf ctfdb \
  -c "UPDATE users SET role = 'admin' WHERE email = 'you@example.com';"
```

## API Overview

### Auth (public)

| Method | Path                  | Body                              | Returns              |
|--------|-----------------------|-----------------------------------|----------------------|
| POST   | `/api/auth/register`  | `{username, email, password}`     | `{token, user}`      |
| POST   | `/api/auth/login`     | `{email, password}`               | `{token, user}`      |

All other endpoints require `Authorization: Bearer <token>`.

### Challenges (authenticated)

| Method | Path                    | Description              |
|--------|-------------------------|--------------------------|
| GET    | `/api/challenges`       | List visible challenges  |
| GET    | `/api/challenges/{id}`  | Get challenge by ID      |

### Submissions (authenticated)

| Method | Path           | Body                         | Description     |
|--------|----------------|------------------------------|-----------------|
| POST   | `/api/submit`  | `{challenge_id, flag}`       | Submit a flag   |

### Scoreboard (authenticated)

| Method | Path              | Description       |
|--------|-------------------|-------------------|
| GET    | `/api/scoreboard` | Leaderboard       |

### Admin (admin role required)

| Method | Path                          | Description            |
|--------|-------------------------------|------------------------|
| POST   | `/api/admin/challenges`       | Create challenge       |
| PUT    | `/api/admin/challenges/{id}`  | Update challenge       |
| DELETE | `/api/admin/challenges/{id}`  | Delete challenge       |
| GET    | `/api/admin/users`            | List all users         |
| PUT    | `/api/admin/users/{id}/disable` | Disable a user       |

> Challenge/scoreboard/submission handlers currently return **501 Not Implemented** — only auth is wired end-to-end.

## Project Structure

```
ctf-platform/
├── backend/
│   ├── main.go              # router wiring + migration runner
│   ├── config/              # env var loading
│   ├── db/                  # connection pool
│   ├── migrations/          # SQL schema
│   ├── middleware/          # JWT auth + admin guard
│   ├── handlers/            # HTTP handlers
│   ├── models/              # DB structs
│   └── utils/               # bcrypt helpers
├── frontend/                # React app (scaffold only)
└── docker-compose.yml
```
