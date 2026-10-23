CREATE TABLE IF NOT EXISTS orders.outbox (
    id             UUID        PRIMARY KEY DEFAULT uuidv7 (),
    aggregate_type TEXT        NOT NULL,
    aggregate_id   TEXT        NOT NULL,
    event_type     TEXT        NOT NULL,
    payload        JSONB       NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);