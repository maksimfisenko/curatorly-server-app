CREATE TABLE IF NOT EXISTS content.curators (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    surname TEXT NOT NULL,
    birthday TIMESTAMPTZ NOT NULL,
    status TEXT NOT NULL,
    project_id BIGINT NULL REFERENCES content.projects(id)
);