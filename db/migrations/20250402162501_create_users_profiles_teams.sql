-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
                       id BIGSERIAL PRIMARY KEY,
                       tg_username VARCHAR(255) UNIQUE NOT NULL,
                       tg_first_name VARCHAR(255) NOT NULL,
                       tg_last_name VARCHAR(255) NOT NULL,
                       photo_url VARCHAR(255),
                       tg_last_login TIMESTAMP WITHOUT TIME ZONE,
                       registered_at TIMESTAMP WITHOUT TIME ZONE,
                       ui_language_code VARCHAR(10) NOT NULL DEFAULT 'en',
                       allows_write_to_pm BOOLEAN DEFAULT TRUE
);

CREATE TABLE teams (
                       id BIGSERIAL PRIMARY KEY,
                       team_name VARCHAR(255) NOT NULL,
                       photo VARCHAR(255),
                       playing_in TIMESTAMP WITHOUT TIME ZONE NOT NULL,
                       stadium VARCHAR(255) NOT NULL,
                       discription VARCHAR(1024)
);

CREATE TABLE profiles (
                          id BIGSERIAL PRIMARY KEY,
                          pref_position VARCHAR(100) NOT NULL,
                          height BIGINT NOT NULL,
                          foot VARCHAR(50) NOT NULL,
                          age BIGINT NOT NULL,
                          playing_frequency BIGINT NOT NULL,
                          games_played BIGINT NOT NULL,
                          city VARCHAR(100) NOT NULL,
                          country VARCHAR(100) NOT NULL,
                          user_id BIGINT UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

ALTER TABLE users ADD COLUMN profile_id BIGINT UNIQUE REFERENCES profiles(id) ON DELETE SET NULL;
ALTER TABLE users ADD COLUMN team_id BIGINT REFERENCES teams(id) ON DELETE SET NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS profiles;
DROP TABLE IF EXISTS teams;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
