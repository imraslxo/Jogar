-- +goose Up
-- +goose StatementBegin
ALTER TABLE profiles
    ALTER COLUMN playing_frequency DROP NOT NULL,
    ALTER COLUMN games_played DROP NOT null;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE profiles
    ALTER COLUMN playing_frequency SET NOT NULL,
    ALTER COLUMN games_played SET NOT null;
-- +goose StatementEnd
