-- +goose Up
-- +goose StatementBegin

-- Create attached_policies table
CREATE TABLE attached_policies (
    adgroup_id  TEXT NOT NULL PRIMARY KEY,
    policy_id   VARCHAR(24) NOT NULL,
    user_id     UUID NOT NULL,
    profile_id  BIGINT NOT NULL,
    campaign_id TEXT NOT NULL,
    is_live     BOOLEAN NOT NULL DEFAULT FALSE,
    FOREIGN KEY (user_id) REFERENCES user_tokens(user_id) ON DELETE CASCADE
);

CREATE INDEX idx_attached_policies_user_id     ON attached_policies(user_id);
CREATE INDEX idx_attached_policies_profile_id  ON attached_policies(profile_id);
CREATE INDEX idx_attached_policies_campaign_id ON attached_policies(campaign_id);

-- Add is_live column to bids table
ALTER TABLE bids
    ADD COLUMN is_live BOOLEAN NOT NULL DEFAULT FALSE;

-- Add composite index on adgroup_id and change_date
CREATE INDEX idx_bids_adgroup_change_date ON bids(adgroup_id, change_date);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_bids_adgroup_change_date;
ALTER TABLE bids DROP COLUMN IF EXISTS is_live;

DROP TABLE IF EXISTS attached_policies;

-- +goose StatementEnd

