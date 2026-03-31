-- +goose Up
ALTER TABLE policy_schedules
    ALTER COLUMN seller_name SET NOT NULL;

ALTER TABLE policy_schedules
    ADD CONSTRAINT policy_schedules_seller_name_not_blank
    CHECK (btrim(seller_name) <> '');

-- +goose Down
ALTER TABLE policy_schedules
    DROP CONSTRAINT IF EXISTS policy_schedules_seller_name_not_blank;

ALTER TABLE policy_schedules
    ALTER COLUMN seller_name DROP NOT NULL;
