-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name TEXT NOT NULL,
    email TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'clerk'
        CHECK (role IN ('administrator', 'clerk')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_lower
    ON users (LOWER(email));

INSERT INTO users (full_name, email, password_hash, role)
VALUES (
    'Admin User',
    'admintest@gmail.com',
    '$2a$10$rZaVPvVJg5mPqFvmJtbKv.Y0sQR3gSpPV7GEKcjOP51bSkXp3HBRa',
    'administrator'
)
ON CONFLICT (lower(email)) DO UPDATE
SET role = 'administrator',
    password_hash = EXCLUDED.password_hash;

-- +goose Down
DROP TABLE IF EXISTS users;
