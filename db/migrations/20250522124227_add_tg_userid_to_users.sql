-- +goose Up
ALTER TABLE "user" ADD COLUMN tg_userId BIGINT UNIQUE;

-- +goose Down
ALTER TABLE "user" DROP COLUMN tg_userId;