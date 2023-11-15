-- +goose Up
CREATE TABLE IF NOT EXISTS Users (
    user_id SERIAL PRIMARY KEY,
    login VARCHAR(50) UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    current_balance DECIMAL(10,2),
    withdrawn DECIMAL(10,2)
);

-- +goose Down
DROP TABLE IF EXISTS Users;