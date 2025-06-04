-- +goose Up
-- +goose StatementBegin
ALTER TABLE profiles ADD COLUMN app_language VARCHAR(10) DEFAULT 'en';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE profiles DROP COLUMN app_language;
-- +goose StatementEnd
