CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name TEXT NOT NULL,
    email TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'clerk'
        CHECK (role IN ('administrator', 'clerk')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    employee_number TEXT,
    phone_number TEXT
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_lower
    ON users (LOWER(email));

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_employee_number
    ON users (employee_number)
    WHERE employee_number IS NOT NULL;

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
