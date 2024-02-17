-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    IF NOT EXISTS transaction_details (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        transaction_id UUID,
        product_id UUID,
        qty INT NOT NULL,
        created_at TIMESTAMP DEFAULT now (),
        updated_at TIMESTAMP DEFAULT now (),
        CONSTRAINT fk_transaction_details_transactions FOREIGN KEY (transaction_id) REFERENCES transactions (id),
        CONSTRAINT fk_transaction_details_products FOREIGN KEY (product_id) REFERENCES products (id)
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS transaction_details;

-- +goose StatementEnd