-- +goose Up
-- +goose StatementBegin
ALTER TABLE profiles DROP CONSTRAINT IF EXISTS profiles_user_id_fkey;
ALTER TABLE profiles RENAME COLUMN user_id TO tg_user_id;
DO $$
    BEGIN
        IF NOT EXISTS (
            SELECT 1
            FROM pg_constraint
            WHERE conname = 'unique_tg_userid'
        ) THEN
            ALTER TABLE "user" ADD CONSTRAINT unique_tg_userid UNIQUE (tg_userid);
        END IF;
    END
$$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE profiles DROP CONSTRAINT IF EXISTS profiles_user_id_fkey;
ALTER TABLE "user" DROP CONSTRAINT IF EXISTS unique_tg_userid;
ALTER TABLE profiles RENAME COLUMN tg_user_id TO user_id;
ALTER TABLE profiles
    ADD CONSTRAINT profiles_user_id_fkey FOREIGN KEY (user_id) REFERENCES "user"(id) ON DELETE CASCADE;
-- +goose StatementEnd
