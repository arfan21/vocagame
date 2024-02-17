-- +goose Up
-- +goose StatementBegin
CREATE TYPE transaction_status AS ENUM ('PROCESSING', 'COMPLETED', 'FAILED');

CREATE TABLE
    IF NOT EXISTS transactions (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        user_id UUID,
        transaction_type_id SMALLINT,
        status transaction_status NOT NULL DEFAULT 'PROCESSING',
        total_amount DECIMAL NOT NULL,
        created_at TIMESTAMP DEFAULT now (),
        updated_at TIMESTAMP DEFAULT now (),
        CONSTRAINT fk_transactions_users FOREIGN KEY (user_id) REFERENCES users (id),
        CONSTRAINT fk_transactions_transaction_types FOREIGN KEY (transaction_type_id) REFERENCES transaction_types (id)
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS transactions;

DROP TYPE IF EXISTS transaction_status;

-- +goose StatementEnd