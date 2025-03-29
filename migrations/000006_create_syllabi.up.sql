create type semester_type as enum (
    'Fall',
    'Winter',
    'Spring',
    'Summer',
    'Spring/Summer'
    );

create table syllabi
(
    id           uuid primary key       default gen_random_uuid(),
    user_id      uuid          not null references users (id),
    course_id    uuid          not null references courses (id),
    file         text          not null,
    content_type text          not null,
    file_size    integer       not null,
    year         smallint      not null check (year > 0),
    semester     semester_type not null,
    date_added   timestamp     not null default now(),
    date_synced  timestamp
);

create index user_id_syllabi_idx on syllabi(user_id);

create index course_id_syllabi_idx on syllabi(course_id);

create table syllabus_views
(
    syllabus_id uuid      not null references syllabi (id),
    user_id     uuid      not null references users (id),
    date_added  timestamp not null default now(),
    constraint syllabus_views_pk primary key (syllabus_id, user_id)
);

create table syllabus_likes
(
    syllabus_id uuid      not null references syllabi (id),
    user_id     uuid      not null references users (id),
    is_dislike  bool      not null,
    date_added  timestamp not null default now(),
    constraint syllabus_likes_pk primary key (syllabus_id, user_id)
);