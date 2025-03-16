create type gender_type as enum ('Male', 'Female', 'Other');

create table users (
    id uuid primary key default gen_random_uuid(),
    program_id uuid references programs (id), -- Form
    full_name text not null, -- Open ID
    nickname text unique, -- Form
    current_year smallint check (current_year > 0), -- Form
    gender gender_type, -- Form
    email text not null unique, -- Open ID
    picture text, -- Open ID
    is_active boolean not null default false,
    date_added timestamp not null default now(),
    date_modified timestamp not null default now()
);

create index program_id_users_idx on users (program_id)