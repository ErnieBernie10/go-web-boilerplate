-- migrate:up
create table app_user (
  id uuid not null default gen_random_uuid(),
  email VARCHAR(255) NOT NULL UNIQUE,
  password_hash VARCHAR(255),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  name VARCHAR(100),
  -- Length may vary based on hashing algorithm
  provider VARCHAR(50),
  -- Optional, for OAuth
  provider_user_id VARCHAR(255) -- Optional, for OAuth
);
-- migrate:down
drop table app_user;