-- +goose Up
ALTER TABLE feeds ADD last_fetched_at TIMESTAMP;

-- +goose Down
ALTER TABLE feeds DROP COLLUMN last_fetched_at;