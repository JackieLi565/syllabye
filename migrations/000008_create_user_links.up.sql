create table user_links
(
    id            uuid primary key   default gen_random_uuid(),
    user_id       uuid      not null references users (id),
    type          text      not null default 'Link',
    link          text      not null,
    date_added    timestamp not null default now(),
    date_modified timestamp not null default now()
);