CREATE TABLE
    IF NOT EXISTS core.curators (
        id bigserial PRIMARY KEY,
        first_name text NOT NULL,
        last_name text NOT NULL,
        middle_name text,
        phone text UNIQUE,
        email citext UNIQUE,
        birth_date DATE CHECK (birth_date <= CURRENT_DATE),
        city text,
        university text,
        profile text,
        created_at timestamptz (0) NOT NULL DEFAULT NOW (),
        updated_at timestamptz (0) NOT NULL DEFAULT NOW (),
        version integer NOT NULL DEFAULT 1
    );