CREATE TABLE transactions
(
    id         UUID PRIMARY KEY,
    service_id UUID             NOT NULL,
    amount     DECIMAL(10, 2)     NOT NULL,
    created_at TIMESTAMP          NOT NULL DEFAULT 'now()'
);