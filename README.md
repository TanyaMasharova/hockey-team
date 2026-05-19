# hockey-team

Backend REST API for a hockey ticket purchasing application.

## Stack

- **Go 1.25** · **Gin** — HTTP framework
- **PostgreSQL 17** — database
- **golang-migrate** — database migrations
- **logrus** — structured logging
- **godotenv** — `.env` file loading

Architecture: `cmd` → `internal/api` → `internal/service` → `internal/repository` → `pkg/database`.

## Local Development

### Prerequisites

- Go 1.25+
- PostgreSQL 17
- [migrate CLI](https://github.com/golang-migrate/migrate) for running migrations manually

### Environment Variables

Copy `.env.example` to `.env`:

| Variable | Description | Default |
|---|---|---|
| `DB_HOST` | PostgreSQL host | `localhost` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_USER` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | — |
| `DB_NAME` | Database name | `hockey_db` |
| `DB_SSLMODE` | SSL mode | `disable` |
| `HTTP_PORT` | HTTP server port | `8080` |
| `GIN_MODE` | Gin mode (`debug` / `release`) | `debug` |
| `CORS_ALLOWED_ORIGINS` | Allowed CORS origins, comma-separated | `http://localhost:3000` |

### Running

```bash
cp .env.example .env
# Fill in .env

# Apply migrations
migrate -path ./migrations -database "postgres://user:pass@localhost:5432/hockey_db?sslmode=disable" up

# Start the server
go run ./cmd/api
```

To run the full stack (backend + frontend + database + proxy) use [hockey-infra](https://github.com/TanyaMasharova/hockey-infra).

## API

All routes are available under the `/api` prefix.

### Authentication

| Method | Path | Description |
|---|---|---|
| `POST` | `/api/register` | Register with phone, email and password |
| `POST` | `/api/login` | Sign in with phone and password |

### Matches

| Method | Path | Description |
|---|---|---|
| `GET` | `/api/matches` | List all matches with opponent info |
| `GET` | `/api/matches/:id` | Get match by ID |
| `GET` | `/api/matchesStats` | Win/loss statistics by finish type (regular, overtime, penalty) |

### User

| Method | Path | Description |
|---|---|---|
| `GET` | `/api/user/:id` | Get user profile |
| `PATCH` | `/api/user/:user_id/profile/field` | Update a single profile field (`full_name`, `phone`, `email`, `birth_date`) |

### Tickets

| Method | Path | Description |
|---|---|---|
| `POST` | `/api/tickets` | Purchase a ticket (reserves the seat, generates QR hash) |
| `GET` | `/api/tickets/user/:user_id` | List user tickets with match and seat details |

### Stadium

| Method | Path | Description |
|---|---|---|
| `GET` | `/api/stadium/sectors` | List stadium sectors with price coefficients |
| `GET` | `/api/stadium/sectors/:sectorId/seats` | List seats in a sector with availability |

### Admin

| Method | Path | Description |
|---|---|---|
| `GET` | `/api/admin/stats-summary` | Sales summary statistics |
| `GET` | `/api/admin/stats` | Full statistics |

## Docker

The image contains the server binary, `migrate` CLI and the `migrations/` directory, making it self-contained for deployment.

```bash
docker build -t hockey-team .
docker run -p 8080:8080 --env-file .env hockey-team
```

## Migrations

Migration files are located in `migrations/`. In production they run automatically via a dedicated container on each deploy.

Apply manually:

```bash
migrate -path ./migrations -database "postgres://..." up
```

Roll back the last migration:

```bash
migrate -path ./migrations -database "postgres://..." down 1
```

## Release

1. Merge changes into `master` via PR (CI checks formatting, `go vet`, build and Docker image).
2. Push a tag: `git tag v1.2.3 && git push origin v1.2.3` — GitHub Actions builds the image and publishes it to GHCR.
3. Create a GitHub Release for the tag — this triggers an automatic deploy to the server.
