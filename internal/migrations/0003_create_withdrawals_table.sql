-- +goose Up
CREATE TABLE IF NOT EXISTS Withdrawals (
    withdrawal_id SERIAL PRIMARY KEY,
    user_id INT REFERENCES Users(user_id),
    order_number VARCHAR(20),
    sum NUMERIC NOT NULL DEFAULT 0,
    processed_at TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS Withdrawals;