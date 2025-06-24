-- +goose Up
-- +goose StatementBegin
ALTER TABLE profiles
    ALTER COLUMN playing_frequency TYPE VARCHAR(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE profiles
    ALTER COLUMN playing_frequency TYPE BIGINT;
-- +goose StatementEnd
