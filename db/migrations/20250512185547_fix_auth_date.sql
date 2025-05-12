-- +goose Up
-- +goose StatementBegin
ALTER TABLE "user" ALTER COLUMN auth_date TYPE BIGINT USING EXTRACT(EPOCH FROM auth_date)::BIGINT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "user" ALTER COLUMN auth_date TYPE TIMESTAMP WITHOUT TIME ZONE USING TO_TIMESTAMP(auth_date);
-- +goose StatementEnd
