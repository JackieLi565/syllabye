create table sessions (
    user_id uuid not null references users(id) on delete cascade,
    date_added timestamp not null default now(),
    date_expires timestamp not null,
    constraint sessions_id primary key (user_id, date_added)
);

create index date_expires_idx on sessions (date_expires);