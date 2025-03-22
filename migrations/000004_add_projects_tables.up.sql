CREATE TABLE IF NOT EXISTS content.projects (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    access_code TEXT NOT NULL,
    creator_id BIGINT NOT NULL REFERENCES auth.users(id),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS content.projects_users (
    project_id BIGINT NOT NULL REFERENCES content.projects(id),
    user_id BIGINT NOT NULL REFERENCES auth.users(id),
    PRIMARY KEY (project_id, user_id)
);