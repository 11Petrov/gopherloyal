-- +goose Up
CREATE TABLE IF NOT EXISTS Orders (
    order_id SERIAL PRIMARY KEY,
    user_id INT REFERENCES Users(user_id),
    order_number VARCHAR(20) UNIQUE,
    uploaded_at TIMESTAMP,
    status VARCHAR(20),
    accrual INT
);

-- +goose Down
DROP TABLE IF EXISTS Orders;