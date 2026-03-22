-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Drop constraints before altering column types
ALTER TABLE bids DROP CONSTRAINT bids_pkey;
ALTER TABLE bids DROP CONSTRAINT bids_user_id_fkey;

-- user_id: BIGINT -> UUID
ALTER TABLE user_tokens
    ALTER COLUMN user_id SET DATA TYPE UUID USING uuid_generate_v4();

ALTER TABLE bids
    ALTER COLUMN user_id SET DATA TYPE UUID USING uuid_generate_v4();

-- campaign_id: BIGINT -> TEXT
ALTER TABLE bids
    ALTER COLUMN campaign_id SET DATA TYPE TEXT USING campaign_id::TEXT;

-- adgroup_id: BIGINT -> TEXT
ALTER TABLE bids
    ALTER COLUMN adgroup_id SET DATA TYPE TEXT USING adgroup_id::TEXT;

-- policy_id: BIGINT -> VARCHAR(24) (MongoDB ObjectID)
ALTER TABLE bids
    ALTER COLUMN policy_id SET DATA TYPE VARCHAR(24) USING policy_id::VARCHAR(24);

-- Restore constraints
ALTER TABLE bids
    ADD CONSTRAINT bids_pkey PRIMARY KEY (user_id, campaign_id);

ALTER TABLE bids
    ADD CONSTRAINT bids_user_id_fkey FOREIGN KEY (user_id) REFERENCES user_tokens(user_id) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE bids DROP CONSTRAINT bids_pkey;
ALTER TABLE bids DROP CONSTRAINT bids_user_id_fkey;

ALTER TABLE bids
    ALTER COLUMN policy_id SET DATA TYPE BIGINT USING 0;

ALTER TABLE bids
    ALTER COLUMN adgroup_id SET DATA TYPE BIGINT USING 0;

ALTER TABLE bids
    ALTER COLUMN campaign_id SET DATA TYPE BIGINT USING 0;

ALTER TABLE bids
    ALTER COLUMN user_id SET DATA TYPE BIGINT USING 0;

ALTER TABLE user_tokens
    ALTER COLUMN user_id SET DATA TYPE BIGINT USING 0;

ALTER TABLE bids
    ADD CONSTRAINT bids_pkey PRIMARY KEY (user_id, campaign_id);

ALTER TABLE bids
    ADD CONSTRAINT bids_user_id_fkey FOREIGN KEY (user_id) REFERENCES user_tokens(user_id) ON DELETE CASCADE;
-- +goose StatementEnd

