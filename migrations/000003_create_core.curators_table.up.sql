CREATE TABLE
    IF NOT EXISTS core.curators (
        id bigserial PRIMARY KEY,
        first_name text NOT NULL,
        last_name text NOT NULL,
        middle_name text,
        phone text,
        email citext,
        birth_date DATE,
        city text,
        university text,
        profile text,
        created_at timestamptz (0) NOT NULL DEFAULT NOW (),
        updated_at timestamptz (0) NOT NULL DEFAULT NOW (),
        version integer NOT NULL DEFAULT 1
    );