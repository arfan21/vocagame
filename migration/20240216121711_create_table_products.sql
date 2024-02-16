-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    IF NOT EXISTS products (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        user_id UUID,
        name VARCHAR(255) NOT NULL,
        description TEXT,
        stok INT NOT NULL,
        price decimal,
        created_at TIMESTAMP DEFAULT now (),
        updated_at TIMESTAMP DEFAULT now (),
        CONSTRAINT fk_products_users FOREIGN KEY (user_id) REFERENCES users (id)
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS products;

-- +goose StatementEnd