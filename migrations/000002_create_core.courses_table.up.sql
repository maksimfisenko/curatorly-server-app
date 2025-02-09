CREATE TABLE
    IF NOT EXISTS core.courses (
        id bigserial PRIMARY KEY,
        title text NOT NULL,
        created_at timestamptz (0) NOT NULL DEFAULT NOW (),
        updated_at timestamptz (0) NOT NULL DEFAULT NOW (),
        version integer not null DEFAULT 1
    );