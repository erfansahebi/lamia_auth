CREATE TABLE users
(
    id         UUID                 DEFAULT uuid_generate_v4() PRIMARY KEY,
    first_name TEXT        NOT NULL,
    last_name  TEXT        NOT NULL,
    email      TEXT        NOT NULL UNIQUE,
    password   json        NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE
OR REPLACE FUNCTION trigger_set_timestamp()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at
= NOW();
RETURN NEW;
END;
$$
LANGUAGE plpgsql;

CREATE TRIGGER set_timestamp
    BEFORE UPDATE
    ON users
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp();