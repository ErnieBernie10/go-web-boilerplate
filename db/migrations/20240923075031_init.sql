-- migrate:up
create table frame (
    id uuid not null default gen_random_uuid(),
    title varchar(255) not null,
    description text not null default '',
    created_at timestamptz not null default now(),
    modified_at timestamptz not null default now(),
    PRIMARY KEY(id)
);

-- migrate:down

drop table frame;

