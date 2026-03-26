-- +goose Up
-- +goose StatementBegin

CREATE INDEX idx_attached_policies_user_id_policy_id
    ON attached_policies(user_id, policy_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_attached_policies_user_id_policy_id;

-- +goose StatementEnd
