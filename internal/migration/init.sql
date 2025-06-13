CREATE TABLE IF NOT EXISTS nasabah
(
    id         SERIAL PRIMARY KEY,
    nama       VARCHAR(255)                                       NOT NULL,
    nik        VARCHAR(255)                                       NOT NULL UNIQUE,
    no_hp      VARCHAR(20)                                        NOT NULL UNIQUE,
    saldo      DOUBLE PRECISION         DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);
CREATE TABLE IF NOT EXISTS transactions (
    id BIGSERIAL PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    transaction_at TIMESTAMP NOT NULL,
    transaction_amount DOUBLE PRECISION NOT NULL
);
