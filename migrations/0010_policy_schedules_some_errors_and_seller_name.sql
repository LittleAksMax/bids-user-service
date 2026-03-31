-- +goose Up
ALTER TYPE policy_schedule_state ADD VALUE IF NOT EXISTS 'SOME ERRORS';

ALTER TABLE policy_schedules
    ADD COLUMN seller_name VARCHAR;

-- +goose Down
ALTER TABLE policy_schedules
    DROP COLUMN IF EXISTS seller_name;

ALTER TYPE policy_schedule_state RENAME TO policy_schedule_state_with_some_errors;

CREATE TYPE policy_schedule_state AS ENUM ('PROCESSING', 'PENDING', 'FAILED');

UPDATE policy_schedules
SET state = 'FAILED'
WHERE state = 'SOME ERRORS';

ALTER TABLE policy_schedules
    ALTER COLUMN state TYPE policy_schedule_state
    USING state::text::policy_schedule_state;

DROP TYPE policy_schedule_state_with_some_errors;
