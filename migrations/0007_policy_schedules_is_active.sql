-- +goose Up
ALTER TABLE policy_schedules
    ADD COLUMN is_active BOOLEAN NOT NULL DEFAULT TRUE;

-- +goose Down
ALTER TABLE policy_schedules
    DROP COLUMN IF EXISTS is_active;
