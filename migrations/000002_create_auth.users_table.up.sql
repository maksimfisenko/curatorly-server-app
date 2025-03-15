CREATE TABLE IF NOT EXISTS auth.users (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    surname text NOT NULL,
    email citext UNIQUE NOT NULL,
    password_hash bytea NOT NULL
);