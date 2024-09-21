CREATE TABLE users (
    user_id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_login          VARCHAR(32) NOT NULL UNIQUE,
    user_nickname       VARCHAR(32),
    user_password_hash  CHAR(60),
    user_birthday       DATE,
    user_created_at     TIMESTAMP NOT NULL DEFAULT NOW(),
    user_updated_at     TIMESTAMP NOT NULL DEFAULT NOW(),
    user_deleted_at     TIMESTAMP
);