-- +goose Up
CREATE TABLE IF NOT EXISTS Users (
    user_id SERIAL PRIMARY KEY,
    login VARCHAR(50) UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    current_balance NUMERIC NOT NULL DEFAULT 0,
    withdrawn NUMERIC NOT NULL DEFAULT 0
);

-- +goose Down
DROP TABLE IF EXISTS Users;