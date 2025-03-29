create table sessions
(
    id         uuid primary key   default gen_random_uuid(),
    user_id    uuid      not null references users (id) on delete cascade,
    date_added timestamp not null default now()
);