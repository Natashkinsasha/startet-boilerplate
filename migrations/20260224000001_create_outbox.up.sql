CREATE TABLE IF NOT EXISTS outbox (
    id          BIGSERIAL PRIMARY KEY,
    event_name  VARCHAR(255) NOT NULL,
    payload     JSONB NOT NULL,
    headers     JSONB NOT NULL DEFAULT '{}',
    created_at  BIGINT NOT NULL DEFAULT 0,
    published   BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX idx_outbox_unpublished ON outbox (id) WHERE published = FALSE;
