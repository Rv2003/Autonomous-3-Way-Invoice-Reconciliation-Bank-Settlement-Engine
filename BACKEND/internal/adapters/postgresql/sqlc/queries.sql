-- name: ListPdf :many
SELECT id, original_name, size_bytes, status, uploaded_at
FROM pdfs
ORDER BY uploaded_at DESC;

-- name: GetPdf :one
SELECT * FROM pdfs
WHERE id = $1;

-- name: CreatePdf :one
INSERT INTO pdfs (original_name, stored_path, size_bytes, sha256)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdatePdfStatus :exec
UPDATE pdfs
SET status = $2, processed_at = now()
WHERE id = $1;