-- +goose Up
ALTER TYPE policy_schedule_state ADD VALUE IF NOT EXISTS 'FAILED';

-- +goose Down
ALTER TYPE policy_schedule_state RENAME TO policy_schedule_state_with_failed;

CREATE TYPE policy_schedule_state AS ENUM ('PROCESSING', 'PENDING');

UPDATE policy_schedules
SET state = 'PENDING'
WHERE state = 'FAILED';

ALTER TABLE policy_schedules
    ALTER COLUMN state TYPE policy_schedule_state
    USING state::text::policy_schedule_state;

DROP TYPE policy_schedule_state_with_failed;
