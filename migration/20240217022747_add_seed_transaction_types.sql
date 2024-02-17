-- +goose Up
-- +goose StatementBegin
INSERT INTO
    transaction_types (name)
VALUES
    ('Deposit'),
    ('Withdrawal'),
    ('Purchase'),
    ('Refund');

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DELETE FROM transaction_types
WHERE
    name IN ('Deposit', 'Withdrawal', 'Purchase', 'Refund');

-- +goose StatementEnd