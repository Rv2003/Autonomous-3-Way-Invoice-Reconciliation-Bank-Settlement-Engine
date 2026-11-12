# Invoice Recon · Settlement Engine (frontend)

React + Material UI implementation of the Figma design "3-way bank reconciliation system".

For complete project setup, environment variables, Docker, database migrations, and test-user instructions, see the [root README](../../../README.md).

```
npm ci
npm run dev
```

Routes: `/login` · `/register` · `/` dashboard · `/documents` · `/review` · `/review/:id` (3-way match detail) · `/ledger` · `/users`. The dashboard and workspace routes require a verified login.

Login posts credentials to `POST /auth/login`, then fetches the signed-in user from
`GET /auth/me` using the HttpOnly auth cookie. Tokens and user data are not stored in
local or session storage. The auth context exposes the user, loading state, login, and
logout functions to components. Registration is currently a presentation-only page;
create users with the backend `cmd/create-user` command.

The Documents page loads PDF records from `http://localhost:8080/pdf` and uploads files
to `http://localhost:8080/pdf/submit` as multipart form data in the `file` field.
Both manual upload screens support multiple PDFs and send them sequentially through
an upload queue. Set `VITE_API_BASE_URL` to override the backend base URL. API requests
that use authentication must include browser credentials; the backend must allow
credentialed requests from the frontend origin. For a hosted frontend, set
`CORS_ALLOWED_ORIGINS` in `BACKEND/.env` to the exact origin (for example,
`https://recon.example.com`, without a trailing slash) and restart the API.
