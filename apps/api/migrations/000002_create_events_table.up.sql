CREATE TABLE IF NOT EXISTS events (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    venue VARCHAR(255) NOT NULL,
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ NOT NULL,
    banner_url TEXT,
    status VARCHAR(32) NOT NULL DEFAULT 'draft', -- draft, published, cancelled, finished
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    CONSTRAINT chk_event_dates CHECK (end_date >= start_date),
    CONSTRAINT chk_event_status CHECK (status IN ('draft','published','cancelled','finished'))
);

CREATE INDEX idx_events_status ON events(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_events_start_date ON events(start_date);
CREATE INDEX idx_events_deleted_at ON events(deleted_at);
