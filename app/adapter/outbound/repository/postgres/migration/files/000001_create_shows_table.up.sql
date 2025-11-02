CREATE TABLE IF NOT EXISTS shows
(
    id    uuid primary key not null,
    title varchar(255)     not null,
    slug  varchar(255)     not null
);
