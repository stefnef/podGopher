CREATE TABLE IF NOT EXISTS episode
(
    id    uuid primary key not null,
    show_id uuid not null references show(id),
    title varchar(255)     not null

);
