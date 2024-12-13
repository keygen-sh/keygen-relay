-- rebuild the table with the old schema
CREATE TABLE _licenses (
  id TEXT PRIMARY KEY,
  file BLOB UNIQUE NOT NULL,
  key TEXT UNIQUE NOT NULL,
  claims INTEGER DEFAULT 0 NOT NULL,
  last_claimed_at DATETIME,
  last_released_at DATETIME,
  node_id INTEGER UNIQUE,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
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
  datetime(last_claimed_at, 'unixepoch') AS last_claimed_at,
  datetime(last_released_at, 'unixepoch') AS last_released_at,
  node_id,
  datetime(created_at, 'unixepoch') AS created_at
FROM
  licenses;

-- drop the old table
DROP TABLE licenses;

-- replace with new table
ALTER TABLE
  _licenses RENAME TO licenses;
