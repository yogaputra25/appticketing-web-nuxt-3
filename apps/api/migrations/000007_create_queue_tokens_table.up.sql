CREATE TABLE IF NOT EXISTS queue_tokens (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    event_id BIGINT NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    token VARCHAR(128) NOT NULL UNIQUE,
    status VARCHAR(32) NOT NULL DEFAULT 'waiting', -- waiting, ready, used, expired
    position INTEGER,
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ready_at TIMESTAMPTZ,
    last_active_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,

    CONSTRAINT chk_queue_status CHECK (status IN ('waiting','ready','used','expired'))
);

CREATE INDEX idx_queue_tokens_user_id ON queue_tokens(user_id);
CREATE INDEX idx_queue_tokens_event_status ON queue_tokens(event_id, status);
CREATE INDEX idx_queue_tokens_token ON queue_tokens(token);
CREATE INDEX idx_queue_tokens_expires_at ON queue_tokens(expires_at) WHERE status IN ('waiting','ready');
