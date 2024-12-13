-- rebuild the table with the old schema
CREATE TABLE _nodes (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  fingerprint TEXT UNIQUE NOT NULL,
  claimed_at DATETIME,
  last_heartbeat_at DATETIME,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
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
  datetime(claimed_at, 'unixepoch') AS claimed_at,
  datetime(last_heartbeat_at, 'unixepoch') AS last_heartbeat_at,
  datetime(created_at, 'unixepoch') AS created_at
FROM
  nodes;

-- drop the old table
DROP TABLE nodes;

-- replace with new table
ALTER TABLE
  _nodes RENAME TO nodes;
