CREATE TABLE IF NOT EXISTS content.courses (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    academic_year TEXT NOT NULL,
    project_id BIGINT NULL REFERENCES content.projects(id)
);