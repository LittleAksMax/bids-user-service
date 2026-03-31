-- +goose Up
CREATE TYPE policy_schedule_state AS ENUM ('PROCESSING', 'PENDING');

ALTER TABLE policy_schedules
    ADD COLUMN state policy_schedule_state NOT NULL DEFAULT 'PENDING';

-- +goose Down
ALTER TABLE policy_schedules
    DROP COLUMN IF EXISTS state;

DROP TYPE IF EXISTS policy_schedule_state;
