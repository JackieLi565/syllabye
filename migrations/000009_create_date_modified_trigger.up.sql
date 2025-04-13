create function date_modified() returns trigger as
$date_modified$
begin
    NEW.date_modified = now();
    return NEW;
end;
$date_modified$ language plpgsql;

create trigger date_modified
    before update
    on users
    for each row
execute function date_modified();

create trigger date_modified
    before update
    on syllabi
    for each row
execute function date_modified();

create trigger date_modified
    before update
    on user_courses
    for each row
execute function date_modified();

create trigger date_modified
    before update
    on user_links
    for each row
execute function date_modified();
