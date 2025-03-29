-- +goose Up
-- +goose StatementBegin
CREATE TABLE profiles (
    id BIGSERIAL PRIMARY KEY,
    pref_position VARCHAR(100),
    height BIGINT,
    foot VARCHAR(50),
    age BIGINT,
    playing_frequency BIGINT,
    games_played BIGINT,
    city VARCHAR(100) NOT NULL,
    country VARCHAR(100) NOT NULL
    user_id INT UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE teams (
    id BIGSERIAL PRIMARY KEY,
    team_name VARCHAR(255) NOT NULL,
    photo VARCHAR(255),
    playing_in TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    stadium VARCHAR(255) NOT NULL,
    discription VARCHAR(1024)
);

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    tg_username VARCHAR(255) UNIQUE NOT NULL,
    tg_first_name VARCHAR(255) NOT NULL,
    tg_last_name VARCHAR(255) NOT NULL,
    photo_url VARCHAR(255),
    tg_last_login TIMESTAMP WITHOUT TIME ZONE,
    registered_at TIMESTAMP WITHOUT TIME ZONE,
    ui_language_code VARCHAR(10) NOT NULL DEFAULT 'en',
    profile_id BIGINT REFERENCES profiles(id) ON DELETE SET NULL,
    team_id BIGINT REFERENCES teams(id) ON DELETE SET NULL,
    allows_write_to_pm BOOLEAN DEFAULT TRUE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS profiles;
DROP TABLE IF EXISTS teams;
-- +goose StatementEnd
