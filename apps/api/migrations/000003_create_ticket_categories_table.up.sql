CREATE TABLE IF NOT EXISTS ticket_categories (
    id BIGSERIAL PRIMARY KEY,
    event_id BIGINT NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    price NUMERIC(15,2) NOT NULL CHECK (price >= 0),
    total_stock INTEGER NOT NULL CHECK (total_stock >= 0),
    available_stock INTEGER NOT NULL CHECK (available_stock >= 0),
    max_per_user INTEGER NOT NULL DEFAULT 4,
    sale_start_at TIMESTAMPTZ,
    sale_end_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    CONSTRAINT chk_stock CHECK (available_stock <= total_stock)
);

CREATE INDEX idx_ticket_categories_event_id ON ticket_categories(event_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_ticket_categories_deleted_at ON ticket_categories(deleted_at);
