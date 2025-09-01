# test-go — Learning Go by Building Small Projects

This repo contains small, focused Go projects to learn the language, web frameworks, DB access, concurrency, and calling external APIs.

## Projects

- Web-Service-Gin
  - REST API using Gin. CRUD for albums. Environment-driven DB config.
- Web-Service-Chi
  - REST API using Chi + PostgreSQL. Uses sqlc for type-safe queries and pgx (stdlib) driver. Schema and queries are versioned.
- Weather-Api
  - HTTP API that detects client IP (handles proxies) and returns current weather from Open‑Meteo, mapping WMO codes to human-friendly descriptions (embedded JSON).
- Test-Connect-DBMS
  - Minimal examples for connecting to a database with environment variables.
- Go-Routine
  - Concurrency demos: goroutines, scheduling, and GOMAXPROCS.
- Make a Module
  - Basics of modules, packages, tests (greetings) and a hello-world app.

## Prerequisites

- Go 1.21+ installed
- PostgreSQL (for Web-Service-Chi)
- Optional tools:
  - sqlc (for Web-Service-Chi code generation): go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
  - psql (to run schema)
  - curl or HTTP client for testing APIs
  - ngrok (optional, to expose local server)

## How to Run

### 1) Web-Service-Chi (Chi + PostgreSQL + sqlc)

- Copy env:
  - cp Web-Service-Chi/.env-example Web-Service-Chi/.env
  - Edit DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME, DB_SSLMODE
- Create schema (any one approach):
  - psql -h 127.0.0.1 -p 5432 -U postgres -d <DB_NAME> -f "Web-Service-Chi/schema/001_albums.sql"
  - or let the app create the table if coded to do so
- Generate sqlc code:
  - cd Web-Service-Chi
  - sqlc generate
- Run:
  - go run main.go
- Notes:
  - sqlc config: Web-Service-Chi/sqlc.yaml
  - Queries: Web-Service-Chi/queries/\*.sql
  - Generated package: Web-Service-Chi/db
  - Typical endpoints: GET/POST/PUT/DELETE /albums, search, etc.

### 2) Web-Service-Gin (Gin REST API)

- Copy env:
  - cp Web-Service-Gin/.env-example Web-Service-Gin/.env
  - Update DB settings as needed
- Run:
  - cd Web-Service-Gin
  - go run main.go

### 3) Weather-Api

- Purpose: Return weather for the caller’s IP, map WMO codes to readable conditions.
- Run:
  - cd Weather-Api
  - go run main.go
- Endpoints (examples):
  - GET /weather → { request_id, client_ip, weather { description, image } }
  - GET /ipinfo (if implemented)
- Data:
  - Embedded code map: Weather-Api/models/weather_codes.json

### 4) Test-Connect-DBMS

- Copy env:
  - cp Test-Connect-DBMS/main/.env-example Test-Connect-DBMS/main/.env
  - Edit DB credentials
- Run:
  - cd Test-Connect-DBMS/main
  - go run main.go

### 5) Go-Routine (Concurrency demos)

- Run:
  - cd Go-Routine
  - go run main.go
- Try different levels of parallelism:
  - Edit runtime.GOMAXPROCS(n) and observe scheduling/prints.

### 6) Make a Module (Modules, packages, tests)

- greetings package:
  - cd "Make a Module/greetings"
  - go test ./...
- hello-world app:
  - cd "Make a Module/hello-world"
  - go run hello-world.go

## Useful Commands (Windows PowerShell)

- Go tidy:
  - go mod tidy
- Run any project:
  - cd <project-folder>; go run .
- Call API:
  - curl http://localhost:8080/
- psql (example):
  - psql -h 127.0.0.1 -p 5432 -U postgres -d <DB_NAME>

## Notes

- Environment variables are loaded via .env in several projects (using godotenv).
- Web-Service-Chi uses PostgreSQL placeholders ($1, $2, …). MySQL-style (?) won’t work there.
- With proxies/tunnels (e.g., ngrok), client IP is taken from X-Forwarded-For/X-Real-IP (RealIP middleware).
- Prefer sqlc + pgx for type-safe, fast DB access; keep schema/queries in versioned files.
