-- Write your migrate up statements here
ALTER TABLE strings DROP CONSTRAINT strings_pkey;

-- 2. Drop the id column
ALTER TABLE strings DROP COLUMN id;

-- 3. Drop the updated_at column
ALTER TABLE strings DROP COLUMN updated_at;

-- 4. Change sha256_hash column to proper length and not null
ALTER TABLE strings
    ALTER COLUMN sha256_hash SET NOT NULL,
    ALTER COLUMN sha256_hash TYPE VARCHAR(64);

-- 5. Add a new column for string length
ALTER TABLE strings
    ADD COLUMN length INT NOT NULL DEFAULT 0;

-- 6. Add the new primary key on sha256_hash
ALTER TABLE strings ADD PRIMARY KEY (sha256_hash);


---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
