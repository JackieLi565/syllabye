create table programs
(
    id         uuid primary key   default gen_random_uuid(),
    faculty_id uuid      not null references faculties (id),
    name       text      not null unique,
    uri        text      not null,
    date_added timestamp not null default now()
);

create index faculty_id_programs_idx on programs (faculty_id);