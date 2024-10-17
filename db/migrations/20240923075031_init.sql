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
  provider_user_id VARCHAR(255), -- Optional, for OAuth,
  PRIMARY KEY(id)
);

create table frame (
    id uuid not null default gen_random_uuid(),
    title varchar(255) not null,
    description text not null default '',
    created_at timestamptz not null default now(),
    modified_at timestamptz not null default now(),
    user_id uuid not null,
    frame_status integer not null,
    PRIMARY KEY(id),
    FOREIGN KEY (user_id) REFERENCES app_user(id)
);
-- migrate:down
drop table frame;
drop table app_user;