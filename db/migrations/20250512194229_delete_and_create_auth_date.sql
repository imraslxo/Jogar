-- +goose Up
-- +goose StatementBegin
ALTER TABLE "user" DROP COLUMN auth_date;
ALTER TABLE "user" ADD COLUMN auth_date BIGINT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "user" ALTER COLUMN auth_date TYPE TIMESTAMP WITHOUT TIME ZONE;
-- +goose StatementEnd
