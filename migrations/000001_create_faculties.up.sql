create extension "pgcrypto";

create table faculties
(
    id         uuid primary key   default gen_random_uuid(),
    name       text      not null,
    date_added timestamp not null default now()
);