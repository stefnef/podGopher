CREATE TABLE IF NOT EXISTS distribution
(
    id      uuid primary key not null,
    show_id uuid             not null references show (id),
    title   varchar(255)     not null,
    slug    varchar(255)     not null

);

CREATE TABLE IF NOT EXISTS show_distributions
(
    show_id         uuid not null references show (id),
    distribution_id uuid not null references distribution (id),

    constraint show_distribution_unique unique (show_id, distribution_id)
);

CREATE INDEX idx_show_distributions_show_id on show_distributions (show_id);

CREATE TABLE IF NOT EXISTS episodes_distributions
(
    episode_id      uuid not null references episode (id),
    distribution_id uuid not null references distribution (id),

    constraint episode_distribution_unique unique (episode_id, distribution_id)
);
