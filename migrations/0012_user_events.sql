-- +goose Up
CREATE TABLE user_events (
    log_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    profile_id BIGINT NOT NULL CHECK (profile_id > 0),
    user_id UUID NOT NULL,
    event_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    log VARCHAR NOT NULL,
    CONSTRAINT fk_user_events_user_id
        FOREIGN KEY (user_id) REFERENCES user_tokens(user_id) ON DELETE CASCADE
);

-- We add the log_id since we are likely to order by log_id descending, so it's good for performance
CREATE INDEX idx_user_events_user_id_profile_id_event_at
    ON user_events(user_id, profile_id, event_at, log_id);

-- +goose Down
DROP INDEX IF EXISTS idx_user_events_user_id_profile_id_event_at;
DROP TABLE IF EXISTS user_events;
