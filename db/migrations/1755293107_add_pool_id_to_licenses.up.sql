-- rebuild the table with the new schema (this is a workaround for sqlite not supporting new foreign keys)
CREATE TABLE _licenses (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  guid TEXT UNIQUE NOT NULL,
  file BLOB UNIQUE NOT NULL,
  key TEXT UNIQUE NOT NULL,
  claims INTEGER DEFAULT 0 NOT NULL,
  last_claimed_at INTEGER,
  last_released_at INTEGER,
  node_id INTEGER UNIQUE,
  pool_id INTEGER,
  created_at INTEGER NOT NULL DEFAULT (unixepoch()),
  FOREIGN KEY (node_id) REFERENCES nodes(id) ON DELETE SET NULL,
  FOREIGN KEY (pool_id) REFERENCES pools(id) ON DELETE SET NULL
);

CREATE INDEX idx_licenses_pool_node ON _licenses(pool_id, node_id);

-- copy data from the old table to the new
INSERT INTO
  _licenses (
    id,
    guid,
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
  guid,
  file,
  key,
  claims,
  last_claimed_at,
  last_released_at,
  node_id,
  created_at
FROM
  licenses;

-- drop the old table
DROP TABLE licenses;

-- replace with new table
ALTER TABLE
  _licenses RENAME TO licenses;

