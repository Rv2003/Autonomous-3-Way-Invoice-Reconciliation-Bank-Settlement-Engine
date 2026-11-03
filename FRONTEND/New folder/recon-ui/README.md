# Invoice Recon · Settlement Engine (frontend)

React + Material UI implementation of the Figma design "3-way bank reconciliation system".

```
npm install
npm run dev
```

Routes: `/login` · `/register` · `/` dashboard · `/documents` · `/review` · `/review/:id` (3-way match detail) · `/ledger` · `/users`.

The login and registration pages are static Material UI designs; their forms do not
connect to an authentication service.

The Documents page loads PDF records from `http://localhost:8080/pdf` and uploads files
to `http://localhost:8080/pdf/submit` as multipart form data in the `file` field.
Both manual upload screens support multiple PDFs and send them sequentially through
an upload queue. Set `VITE_API_BASE_URL` to override the backend base URL. The backend
permits browser requests from loopback origins for local development.
