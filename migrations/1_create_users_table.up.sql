CREATE TABLE users (
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    login VARCHAR(60) NOT NULL UNIQUE,
    nickname VARCHAR(40) NOT NULL,
    password_hash VARCHAR(60) NOT NULL,
    is_teacher BOOLEAN NOT NULL DEFAULT FALSE
);