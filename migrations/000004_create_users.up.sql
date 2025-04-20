create table users
(
    id            uuid primary key   default gen_random_uuid(),
    program_id    uuid references programs (id),
    full_name     text      not null,
    nickname      text unique check (nickname ~ '^[a-z0-9._-]{3,30}$'),
    current_year  smallint check (current_year > 0 and current_year <= 8),
    gender        text check (gender in ('Male', 'Female', 'Other')),
    email         text      not null unique,
    bio           text,
    ig_handle     text,
    picture       text,
    is_active     boolean   not null default false,
    date_added    timestamp not null default now(),
    date_modified timestamp not null default now()
);

create index program_id_users_idx on users (program_id)