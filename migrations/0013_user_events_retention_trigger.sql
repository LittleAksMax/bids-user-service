-- +goose Up
-- +goose StatementBegin

-- Only keep top 250 logs per user_id, profile_id combination
CREATE OR REPLACE FUNCTION trim_user_events_to_250()
RETURNS TRIGGER AS $$
BEGIN
    -- https://www.postgresql.org/docs/current/functions-admin.html#FUNCTIONS-ADVISORY-LOCKS
    -- We lock all fields with to ensure correctness, possibly reducing concurrency for 
    -- We have to use this hashtext stuff because UUID (user_id) and BIGINT (profile_id) are
    -- not accepted directly. So we hash their text representations to get simple INT types,
    -- which are accepted by hashtext.
    -- Hashing can theoretically collide for different values, reducing concurrency, but will still
    -- necessarily be correct.
    PERFORM pg_advisory_xact_lock(
        hashtext(NEW.user_id::TEXT),
        hashtext(NEW.profile_id::TEXT)
    );

    DELETE FROM user_events
    WHERE log_id IN (
        SELECT log_id
        FROM user_events
        WHERE user_id = NEW.user_id
          AND profile_id = NEW.profile_id
        ORDER BY event_at DESC, log_id DESC
        OFFSET 250
    );

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_trim_user_events_to_250
AFTER INSERT ON user_events
FOR EACH ROW
EXECUTE FUNCTION trim_user_events_to_250();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_trim_user_events_to_250 ON user_events;
DROP FUNCTION IF EXISTS trim_user_events_to_250();
-- +goose StatementEnd
