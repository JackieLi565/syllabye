create table user_courses
(
    user_id        uuid      not null references users (id) on delete cascade,
    course_id      uuid      not null references courses (id),
    year_taken     smallint check ( year_taken > 0 and year_taken <= 8 ),
    semester_taken semester_type,
    date_added     timestamp not null default now(),
    date_modified  timestamp not null default now(),
    constraint user_course_pk primary key (user_id, course_id)
);