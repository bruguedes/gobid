-- Write your migrate up statements here
CREATE TABLE sessions (
	token TEXT PRIMARY KEY,
	data BYTEA NOT NULL,
	expiry TIMESTAMPTZ NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions (expiry);

---- create above / drop below ----
drop index if exists sessions_expiry_idx;
drop table if exists sessions;




-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
