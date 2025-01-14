-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
-- +goose ENVSUB ON
    CREATE TABLE IF NOT EXISTS "${DATABASE_MIGRATION_TABLE_PREFIX}user" (
                                              id bigserial NOT NULL,
                                              login varchar(250) NOT NULL,
    pass_hash varchar(250) NOT NULL,
    CONSTRAINT "${DATABASE_MIGRATION_TABLE_PREFIX}user_pk" PRIMARY KEY (id),
    CONSTRAINT "${DATABASE_MIGRATION_TABLE_PREFIX}user_login_udx" UNIQUE (login)
    );

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
-- +goose ENVSUB ON
DROP TABLE IF EXISTS "${DATABASE_MIGRATION_TABLE_PREFIX}user";
