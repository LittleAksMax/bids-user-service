-- +goose Up
CREATE TABLE policy_schedules (
    user_id UUID NOT NULL,
    profile_id BIGINT NOT NULL CHECK (profile_id > 0),
    due_at TIMESTAMPTZ NOT NULL,
    interval_minutes BIGINT NOT NULL CHECK (interval_minutes > 0),
    PRIMARY KEY (user_id, profile_id),
    CONSTRAINT fk_policy_schedules_user_id
        FOREIGN KEY (user_id) REFERENCES user_tokens(user_id) ON DELETE CASCADE
);

CREATE INDEX idx_policy_schedules_due_at ON policy_schedules(due_at);

-- +goose Down
DROP INDEX IF EXISTS idx_policy_schedules_due_at;
DROP TABLE IF EXISTS policy_schedules;
