-- Write your migrate up statements here

ALTER TABLE products RENAME COLUMN price TO base_price;


---- create above / drop below ----
ALTER TABLE products RENAME COLUMN base_price; TO price;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
