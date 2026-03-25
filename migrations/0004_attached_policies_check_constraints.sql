-- +goose Up
-- +goose StatementBegin

ALTER TABLE attached_policies
    ADD CONSTRAINT chk_attached_policies_profile_id_nonnegative
        CHECK (profile_id > 0),
    ADD CONSTRAINT chk_attached_policies_campaign_id_nonempty
        CHECK (campaign_id <> ''),
    ADD CONSTRAINT chk_attached_policies_adgroup_id_nonempty
        CHECK (adgroup_id <> ''),
    ADD CONSTRAINT chk_attached_policies_policy_id_nonempty
        CHECK (policy_id <> '');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE attached_policies
    DROP CONSTRAINT IF EXISTS chk_attached_policies_policy_id_nonempty,
    DROP CONSTRAINT IF EXISTS chk_attached_policies_adgroup_id_nonempty,
    DROP CONSTRAINT IF EXISTS chk_attached_policies_campaign_id_nonempty,
    DROP CONSTRAINT IF EXISTS chk_attached_policies_profile_id_nonnegative;

-- +goose StatementEnd
