-- Write your migrate up statements here
create table if not exists users (
    id UUID PRIMARY KEY NOT NULL default gen_random_uuid(),
    user_name VARCHAR(50) UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hash BYTEA NOT NULL,
    bio TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
---- create above / drop below ----
 DROP TABLE IF EXISTS users;


-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
