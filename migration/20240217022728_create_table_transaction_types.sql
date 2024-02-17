-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    IF NOT EXISTS transaction_types (
        id SMALLSERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL UNIQUE,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS transaction_types;

-- +goose StatementEnd