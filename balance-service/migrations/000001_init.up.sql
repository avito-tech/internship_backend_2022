CREATE TYPE transaction_status AS ENUM(
    'processed',
    'accepted',
    'rejected'
);

CREATE TABLE users
(
    id      UUID PRIMARY KEY,
    user_id UUID UNIQUE NOT NULL
);

CREATE TABLE transfers
(
    id           UUID PRIMARY KEY,
    from_user_id UUID         NOT NULL,
    to_user_id   UUID         NOT NULL,
    amount       DECIMAL(10, 2) NOT NULL,
    created_at   TIMESTAMP      NOT NULL DEFAULT 'now()'
);

CREATE TABLE transactions
(
    id         UUID PRIMARY KEY,
    user_id    UUID             NOT NULL,
    service_id UUID             NOT NULL,
    order_id   UUID             NOT NULL,
    amount     DECIMAL(10, 2)     NOT NULL,
    status     transaction_status NOT NULL DEFAULT 'processed',
    created_at TIMESTAMP          NOT NULL DEFAULT 'now()',
    updated_at TIMESTAMP          NOT NULL DEFAULT 'now()'
);

CREATE TABLE balances
(
    id      UUID PRIMARY KEY,
    user_id UUID         NOT NULL,
    balance DECIMAL(10, 2) NOT NULL,
    reserved DECIMAL(10, 2) NOT NULL
);

CREATE INDEX ON users (user_id);
CREATE INDEX ON transfers (from_user_id);
CREATE INDEX ON transfers (to_user_id);
CREATE INDEX ON transfers (from_user_id, to_user_id);
CREATE INDEX ON balances (user_id);