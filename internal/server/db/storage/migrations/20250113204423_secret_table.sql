-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
-- +goose ENVSUB ON
    CREATE TABLE IF NOT EXISTS "${DATABASE_MIGRATION_TABLE_PREFIX}secret" (
                                                id bigserial NOT NULL,
                                                user_id bigint NULL,
                                                updated_at bigint NOT NULL,
                                                is_deleted boolean NOT NULL DEFAULT FALSE,
                                                data bytea NOT NULL,
                                                CONSTRAINT "${DATABASE_MIGRATION_TABLE_PREFIX}secret_pk" PRIMARY KEY (id),
    CONSTRAINT "${DATABASE_MIGRATION_TABLE_PREFIX}secret_user_fk" FOREIGN KEY (user_id) REFERENCES "${DATABASE_MIGRATION_TABLE_PREFIX}user" (id)
    );

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
-- +goose ENVSUB ON
DROP TABLE IF EXISTS "${DATABASE_MIGRATION_TABLE_PREFIX}secret";