-- +goose Up
-- +goose StatementBegin
ALTER TABLE profiles
    ALTER COLUMN pref_position DROP NOT NULL,
    ALTER COLUMN height DROP NOT NULL,
    ALTER COLUMN foot DROP NOT NULL,
    ALTER COLUMN age DROP NOT NULL,
    ALTER COLUMN city DROP NOT NULL,
    ALTER COLUMN country DROP NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE profiles
    ALTER COLUMN pref_position SET NOT NULL,
    ALTER COLUMN height SET NOT NULL,
    ALTER COLUMN foot SET NOT NULL,
    ALTER COLUMN age SET NOT NULL,
    ALTER COLUMN city SET NOT NULL,
    ALTER COLUMN country SET NOT NULL;
-- +goose StatementEnd
