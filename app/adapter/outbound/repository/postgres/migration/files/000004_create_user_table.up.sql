CREATE TABLE IF NOT EXISTS podcast_user
(
    id       uuid primary key not null,
    username varchar(255)     not null,
    is_admin boolean          not null default false
);

CREATE TABLE IF NOT EXISTS show_users
(
    show_id uuid not null references show (id),
    user_id uuid not null references podcast_user (id),

    constraint show_user_unique unique (show_id, user_id)
);

CREATE INDEX if not exists idx_show_users_show_id on show_users (show_id);

CREATE TABLE IF NOT EXISTS user_roles
(
    show_id    uuid not null references show (id),
    user_id    uuid not null references podcast_user (id),
    role       varchar(50) not null,

    constraint user_roles_unique unique (show_id, user_id)
);

CREATE INDEX if not exists idx_user_roles on user_roles (show_id, user_id);