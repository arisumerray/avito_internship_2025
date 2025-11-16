CREATE TABLE IF NOT EXISTS teams
(
    id   SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS users
(
    id        TEXT PRIMARY KEY,
    name      TEXT    NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    team_id   INTEGER REFERENCES teams (id)
);

CREATE INDEX ON users (is_active, team_id);

CREATE TABLE IF NOT EXISTS pr
(
    id           TEXT PRIMARY KEY,
    name         TEXT NOT NULL,
    author_id    TEXT REFERENCES users (id),
    status       TEXT NOT NULL,
    reviewer1_id TEXT REFERENCES users (id),
    reviewer2_id TEXT REFERENCES users (id),
    created_at   TIMESTAMP DEFAULT NOW(),
    merged_at    TIMESTAMP
);
