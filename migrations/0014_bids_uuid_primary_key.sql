-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

ALTER TABLE bids
    ADD COLUMN id UUID,
    ADD COLUMN profile_id BIGINT;

UPDATE bids b
SET profile_id = ap.profile_id
FROM attached_policies ap
WHERE b.adgroup_id = ap.adgroup_id
  AND b.user_id = ap.user_id
  AND b.campaign_id = ap.campaign_id;

UPDATE bids
SET id = uuid_generate_v4()
WHERE id IS NULL;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM bids
        WHERE profile_id IS NULL
    ) THEN
        RAISE EXCEPTION 'cannot backfill bids.profile_id for all rows';
    END IF;
END
$$;

ALTER TABLE bids
    ALTER COLUMN id SET DEFAULT uuid_generate_v4(),
    ALTER COLUMN id SET NOT NULL,
    ALTER COLUMN profile_id SET NOT NULL;

ALTER TABLE bids
    DROP CONSTRAINT bids_pkey,
    ADD CONSTRAINT bids_pkey PRIMARY KEY (id);

CREATE INDEX idx_bids_user_id_campaign_id ON bids(user_id, campaign_id);
CREATE INDEX idx_bids_campaign_id ON bids(campaign_id);
CREATE INDEX idx_bids_profile_id_change_date ON bids(profile_id, change_date);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM bids
        GROUP BY user_id, campaign_id
        HAVING COUNT(*) > 1
    ) THEN
        RAISE EXCEPTION 'cannot restore bids primary key: duplicate (user_id, campaign_id) rows exist';
    END IF;
END
$$;

DROP INDEX IF EXISTS idx_bids_campaign_id;
DROP INDEX IF EXISTS idx_bids_user_id_campaign_id;
DROP INDEX IF EXISTS idx_bids_profile_id_change_date;

ALTER TABLE bids
    DROP CONSTRAINT bids_pkey,
    ADD CONSTRAINT bids_pkey PRIMARY KEY (user_id, campaign_id);

ALTER TABLE bids
    ALTER COLUMN id DROP DEFAULT;

ALTER TABLE bids
    DROP COLUMN id,
    DROP COLUMN profile_id;
-- +goose StatementEnd
