-- +goose Up
-- +goose StatementBegin
ALTER TABLE profiles
    ALTER COLUMN age TYPE VARCHAR(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE profiles
    ALTER COLUMN age TYPE BIGINT;
-- +goose StatementEnd
