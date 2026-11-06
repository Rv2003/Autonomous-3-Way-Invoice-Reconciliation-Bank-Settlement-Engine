# Autonomous 3-Way Invoice Reconciliation & Bank Settlement Engine

This repository contains a Go API, PostgreSQL database, and React/Vite frontend for invoice reconciliation.

## Prerequisites

- Docker with Docker Compose
- Go 1.27.1 or later (for local backend development and creating users)
- Node.js 22.12 or later and npm (for the frontend)
- OpenSSL (to generate secrets)

## Environment setup

The backend Compose project reads environment variables from `BACKEND/.env`. Create it from the committed template:

### WSL / Bash

```bash
cp BACKEND/.env.example BACKEND/.env
openssl rand -hex 32
openssl rand -hex 32
```

Use the first generated value for `JWT_SECRET` and the second for `POSTGRES_PASSWORD` in `BACKEND/.env`. Both should be distinct random values. Keep `.env` private; it is ignored by Git.

### PowerShell

```powershell
Copy-Item BACKEND/.env.example BACKEND/.env
notepad BACKEND\.env
```

Generate two distinct values with an available OpenSSL installation:

```powershell
openssl rand -hex 32
openssl rand -hex 32
```

Paste one value into `JWT_SECRET` and the other into `POSTGRES_PASSWORD` in `BACKEND/.env`. Do not commit or share this file.

| Variable | Required | Purpose |
|---|---|---|
| `JWT_SECRET` | Yes | Signs API login tokens; use a random value at least 32 characters long. |
| `POSTGRES_PASSWORD` | Recommended | Sets the local database password used by PostgreSQL and the API connection. |
| `COOKIE_SECURE` | Production | Set to `true` when the API is served over HTTPS; leave `false` for local HTTP development. |
| `GOOSE_DBSTRING` | For local Go runs/tools | Local PostgreSQL connection string; set it in `BACKEND/.env` using the same password as `POSTGRES_PASSWORD`. Docker Compose supplies the API's value automatically. |
| `SEED_USER_PASSWORD` | Only when creating a test user | Password input consumed by the one-off user-creation command. |
| `GMAIL_CLIENT_ID` | For Gmail authorization/testing | OAuth client ID from your Google Cloud Web application credentials. |
| `GMAIL_CLIENT_SECRET` | For Gmail authorization/testing | OAuth client secret from your Google Cloud **Web application** credentials. Keep this value private. |
| `GMAIL_TOKEN_FILE` | Optional | Path to the OAuth token JSON file; defaults to `tokens/gmail-token.json` when run from `BACKEND`. Docker Compose mounts this folder for the API. |
| `GMAIL_REDIRECT_URL` | Optional | OAuth callback override for local API runs; defaults to `http://localhost:8080/email/gmail/callback`. Docker Compose uses `http://localhost:8081/email/gmail/callback`. Register both URIs in Google Cloud if using both modes. |
| `GMAIL_FRONTEND_URL` | Optional | Where to return after Gmail authorization; defaults to `http://localhost:5173/`. |
| `VITE_API_BASE_URL` | Optional | Frontend API base URL; defaults to `http://localhost:8080`. |

### Configure Gmail OAuth

In Google Cloud, enable the Gmail API, configure the OAuth consent screen, and create an OAuth client with the **Web application** type. Add the exact callback URI `http://localhost:8081/email/gmail/callback` for Docker Compose and/or `http://localhost:8080/email/gmail/callback` for a locally run API to its authorized redirect URIs. Add the client ID and secret to `BACKEND/.env` as `GMAIL_CLIENT_ID` and `GMAIL_CLIENT_SECRET`. If the consent screen is in testing mode, add the Gmail account as a test user. Keep `.env` private.

Set `GOOSE_DBSTRING` in `BACKEND/.env` to connect local Go tools to PostgreSQL. Use the same password as `POSTGRES_PASSWORD` and `host=localhost` when the database runs in Docker.

Admins connect the shared inbox from the dashboard using **Connect Gmail**. Google prompts for read-only access once; on callback, the app verifies the Gmail address, registers it as an active mailbox, saves the token, and returns to the dashboard. For Docker Compose, use port `8081` in `GMAIL_REDIRECT_URL`. For `go run ./cmd`, set it to `http://localhost:8080/email/gmail/callback` and register that exact URI in Google Cloud. The Gmail account must be an approved test user if the consent screen is in testing mode.

To check an already connected mailbox, run this from the `BACKEND` directory:

```powershell
go run ./cmd/gmail-test
```

The dashboard fetches the latest 10 inbox messages for administrator accounts; use **Fetch latest emails** to refresh. The messages endpoint returns sender, subject, snippet, received date, and whether attachments exist; it does not download message bodies or attachments. Docker Compose passes OAuth settings from `BACKEND/.env` to the API and mounts the token directory so the callback can save the token. The token directory is ignored by Git.

## Start the application with Docker

Run these commands from the repository root:

```bash
cd BACKEND
docker compose up -d --build
docker compose ps
docker compose logs -f api
```

The Compose stack starts the API and PostgreSQL. The API listens on container port `8080` and is exposed on host port `8081`; PostgreSQL is available on port `5432`. Database migrations are embedded in the API and run automatically before it begins serving requests.

Check the API health endpoint at <http://localhost:8081/health>. Stop the services without deleting database data with:

```bash
docker compose down
```

Avoid `docker compose down -v` unless you intend to delete the database volume.

### Automatic API rebuilds during development

The API service watches backend files and rebuilds its image when they change. From `BACKEND`, run this in a terminal and leave it open while developing:

```bash
docker compose up --build --watch
```

Edits under `BACKEND` trigger an API rebuild and restart; `.git`, `uploads`, and Markdown files are ignored. Press `Ctrl+C` to stop the watch session. For a one-time detached start without watching, use `docker compose up -d --build` instead.

## Start the frontend

In a second terminal, from the repository root:

```bash
cd "FRONTEND/New folder/recon-ui"
npm ci
npm run dev
```

Open <http://localhost:5173>. The frontend defaults to `http://localhost:8080`, which is the local Go API port. When using Docker Compose, create `FRONTEND/New folder/recon-ui/.env.local` with:

```env
VITE_API_BASE_URL=http://localhost:8081
```

Restart Vite after changing frontend environment variables. The API allows local loopback origins and credentialed requests. Use `http://localhost:5173` during development. For HTTPS deployments, set `COOKIE_SECURE=true` so the authentication cookie is only sent over HTTPS.

## Create a test login

Accounts are not created by the frontend registration screen. Use the backend CLI to create a test account; it hashes the password with bcrypt. Start the database first (the full Docker setup above is fine), then run these commands from WSL/Bash, starting at the repository root:

```bash
cd BACKEND
set -a
source .env
set +a
export GOOSE_DBSTRING="host=localhost user=invoice_app password=${POSTGRES_PASSWORD} dbname=invoice_reconciliation sslmode=disable"
read -r -p "Full name: " USER_NAME
read -r -p "Email: " USER_EMAIL
read -r -s -p "Test account password: " SEED_USER_PASSWORD
printf '\n'
export SEED_USER_PASSWORD
go run ./cmd/create-user -name "$USER_NAME" -email "$USER_EMAIL"
unset SEED_USER_PASSWORD GOOSE_DBSTRING JWT_SECRET POSTGRES_PASSWORD
```

Use the email and password entered above to sign in at <http://localhost:5173/login>. The login flow calls `POST /auth/login`, which returns the signed JWT and user as JSON and sets an HttpOnly authentication cookie, then fetches the current user's profile from `GET /auth/me`. The verified user is shared with components through the React auth context; the JWT is not stored in browser local or session storage. Unchecking “Keep me signed in” creates a session cookie/token; checking it makes the cookie/token valid for 30 days.

## Core API routes

| Method | Path | Access | Purpose |
|---|---|---|---|
| `GET` | `/health` | Public | Health check |
| `POST` | `/auth/login` | Public | Authenticate with email and password |
| `GET` | `/auth/me` | HttpOnly auth cookie | Return the authenticated user's profile |
| `POST` | `/auth/logout` | Public; clears auth cookie | Sign out |
| `GET` | `/email/gmail/connect` | Administrator | Start shared Gmail OAuth connection |
| `GET` | `/email/gmail/callback` | Administrator | Complete OAuth, register the mailbox, and save its token |
| `GET` | `/email/messages` | Administrator | Fetch recent messages from the registered Gmail inbox |
| `GET` | `/pdf` | Public | List PDF records |
| `GET` | `/pdf/{id}/file` | Public | View an uploaded PDF inline |
| `POST` | `/pdf/submit` | Public | Upload a PDF |

Authentication uses an HttpOnly, SameSite=Lax cookie. The login JSON response includes the JWT as `token` alongside the `user` object. JWTs expire after 30 minutes by default or 30 days when “Keep me signed in” is selected. JavaScript cannot read the cookie; browser requests send it with credentials.

## Local development and checks

Frontend:

```bash
cd "FRONTEND/New folder/recon-ui"
npm run build
```

Backend:

```bash
cd BACKEND
go test ./...
go build ./cmd
```

To run the Go API directly instead of the API container, keep PostgreSQL running, load `BACKEND/.env` into the shell, set `GOOSE_DBSTRING` to use `host=localhost` and the same `POSTGRES_PASSWORD`, then run `go run ./cmd`. Do not run both API versions on port `8080` simultaneously.

## Troubleshooting

- **Compose says `JWT_SECRET` is missing:** create `BACKEND/.env` from `BACKEND/.env.example` and set `JWT_SECRET`.
- **Login returns `401`:** check that the user exists, the email is correct, and the password matches the one used when creating the account.
- **Login reports a network or CORS error:** verify the API is running on port `8080` for `go run ./cmd` or `8081` for Docker Compose, open the frontend at `http://localhost:5173`, and restart/rebuild the API after backend CORS changes.
- **A database table is missing:** check `docker compose logs api` for migration errors. Migrations run at API startup and do not require deleting the database volume.

