CREATE TABLE IF NOT EXISTS show
(
    id    uuid primary key not null,
    title varchar(255)     not null,
    slug  varchar(255)     not null
);
