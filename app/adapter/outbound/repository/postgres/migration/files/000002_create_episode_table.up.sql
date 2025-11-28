CREATE TABLE IF NOT EXISTS episode
(
    id      uuid primary key not null,
    show_id uuid             not null references show (id),
    title   varchar(255)     not null

);

CREATE TABLE IF NOT EXISTS show_episodes
(
    show_id    uuid not null references show (id),
    episode_id uuid not null references episode (id),

    constraint show_episode_unique unique (show_id, episode_id)
);

CREATE INDEX idx_show_episodes_show_id on show_episodes (show_id);