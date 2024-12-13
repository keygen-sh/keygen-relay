-- rebuild the table with the new schema (this is a workaround for sqlite not supporting ALTER TABLE x ALTER COLUMN y NOT NULL)
CREATE TABLE _nodes (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  fingerprint TEXT UNIQUE NOT NULL,
  claimed_at INTEGER,
  last_heartbeat_at INTEGER,
  created_at INTEGER NOT NULL DEFAULT (unixepoch())
);

-- copy data from the old table to the new
INSERT INTO
  _nodes (
    id,
    fingerprint,
    claimed_at,
    last_heartbeat_at,
    created_at
  )
SELECT
  id,
  fingerprint,
  strftime('%s', claimed_at) AS claimed_at,
  strftime('%s', last_heartbeat_at) AS last_heartbeat_at,
  strftime('%s', created_at) AS created_at
FROM
  nodes;

-- drop the old table
DROP TABLE nodes;

-- replace with new table
ALTER TABLE
  _nodes RENAME TO nodes;
