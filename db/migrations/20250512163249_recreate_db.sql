-- +goose Up
-- +goose StatementBegin
ALTER TABLE teams RENAME COLUMN discription TO description;
ALTER TABLE users RENAME COLUMN registered_at TO auth_date;
ALTER TABLE users DROP COLUMN tg_last_login;
ALTER TABLE users ADD COLUMN is_premium BOOLEAN NOT NULL DEFAULT false;

ALTER TABLE users RENAME TO "user";
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE teams RENAME COLUMN description TO discription;
ALTER TABLE "user" RENAME COLUMN auth_date TO registered_at;
ALTER TABLE "user" ADD COLUMN tg_last_login TIMESTAMP WITHOUT TIME ZONE;
ALTER TABLE "user" DROP COLUMN is_premium;

ALTER TABLE "user" RENAME TO users;
-- +goose StatementEnd
