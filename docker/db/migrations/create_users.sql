CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email       VARCHAR(255) NOT NULL UNIQUE,
    password    VARCHAR(255),                  -- nullable for OAuth users
    name        VARCHAR(255) NOT NULL,
    avatar_url  VARCHAR(500),
    provider    VARCHAR(50) NOT NULL DEFAULT 'email',  -- 'email' | 'google'
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
