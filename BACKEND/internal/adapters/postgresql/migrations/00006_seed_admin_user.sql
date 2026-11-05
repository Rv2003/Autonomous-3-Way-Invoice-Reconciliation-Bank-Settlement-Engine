-- +goose Up
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
DELETE FROM users WHERE lower(email) = 'admintest@gmail.com';

