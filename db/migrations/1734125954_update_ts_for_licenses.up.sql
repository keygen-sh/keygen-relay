-- rebuild the table with the new schema (this is a workaround for sqlite not supporting ALTER TABLE x ALTER COLUMN y NOT NULL)
CREATE TABLE _licenses (
  id TEXT PRIMARY KEY,
  file BLOB UNIQUE NOT NULL,
  key TEXT UNIQUE NOT NULL,
  claims INTEGER DEFAULT 0 NOT NULL,
  last_claimed_at INTEGER,
  last_released_at INTEGER,
  node_id INTEGER UNIQUE,
  created_at INTEGER NOT NULL DEFAULT (unixepoch()),
  FOREIGN KEY (node_id) REFERENCES nodes(id) ON DELETE SET NULL
);

-- copy data from the old table to the new
INSERT INTO
  _licenses (
    id,
    file,
    key,
    claims,
    last_claimed_at,
    last_released_at,
    node_id,
    created_at
  )
SELECT
  id,
  file,
  key,
  claims,
  strftime('%s', last_claimed_at) AS last_claimed_at,
  strftime('%s', last_released_at) AS last_released_at,
  node_id,
  strftime('%s', created_at) AS created_at
FROM
  licenses;

-- drop the old table
DROP TABLE licenses;

-- replace with new table
ALTER TABLE
  _licenses RENAME TO licenses;
