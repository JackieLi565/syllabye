create table course_categories
(
    id         uuid primary key   default gen_random_uuid(),
    name       text      not null unique,
    date_added timestamp not null default now()
);

create table courses
(
    id          uuid primary key     default gen_random_uuid(),
    category_id uuid        not null references course_categories (id),
    title       text        not null,
    description text,
    uri         text        not null,
    course      varchar(10) not null,
    alpha       varchar(5),
    code        varchar(5),
    date_added  timestamp   not null default now()
);

create index category_id_courses_idx on courses (category_id);
