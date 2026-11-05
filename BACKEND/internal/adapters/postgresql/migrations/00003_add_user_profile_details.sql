-- +goose Up
ALTER TABLE users
    ADD COLUMN employee_number TEXT,
    ADD COLUMN phone_number TEXT;

CREATE UNIQUE INDEX idx_users_employee_number
    ON users (employee_number)
    WHERE employee_number IS NOT NULL;

-- +goose Down
DROP INDEX IF EXISTS idx_users_employee_number;

ALTER TABLE users
    DROP COLUMN IF EXISTS phone_number,
    DROP COLUMN IF EXISTS employee_number;
