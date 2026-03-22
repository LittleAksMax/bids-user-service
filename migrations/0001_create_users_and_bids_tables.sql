-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_tokens (
    user_id BIGINT PRIMARY KEY,
    refresh_token_eu TEXT,
    refresh_token_us TEXT,
    refresh_token_fe TEXT
);

CREATE TABLE bids (
    user_id BIGINT NOT NULL,
    campaign_id BIGINT NOT NULL,
    adgroup_id BIGINT NOT NULL,
    policy_id BIGINT NOT NULL,
    from_bid DOUBLE PRECISION NOT NULL,
    to_bid DOUBLE PRECISION NOT NULL,
    change_date TIMESTAMP NOT NULL,
    PRIMARY KEY (user_id, campaign_id),
    FOREIGN KEY (user_id) REFERENCES user_tokens(user_id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS bids;
DROP TABLE IF EXISTS user_tokens;
-- +goose StatementEnd

