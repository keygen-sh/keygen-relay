-- rebuild the licenses table without node_id, adding seats
CREATE TABLE _licenses (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  guid TEXT UNIQUE NOT NULL,
  file BLOB UNIQUE NOT NULL,
  key TEXT UNIQUE NOT NULL,
  claims INTEGER DEFAULT 0 NOT NULL,
  last_claimed_at INTEGER,
  last_released_at INTEGER,
  seats INTEGER DEFAULT 1 NOT NULL,
  pool_id INTEGER,
  created_at INTEGER NOT NULL DEFAULT (unixepoch()),
  FOREIGN KEY (pool_id) REFERENCES pools(id) ON DELETE SET NULL
);

-- copy data from the old table (existing licenses get seats = 1)
INSERT INTO
  _licenses (
    id,
    guid,
    file,
    key,
    claims,
    last_claimed_at,
    last_released_at,
    pool_id,
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
  pool_id,
  created_at
FROM
  licenses;

-- drop the old table (active node_id claims are released; clients re-claim on next heartbeat)
DROP TABLE licenses;

-- replace with new table
ALTER TABLE
  _licenses RENAME TO licenses;

-- create the license_nodes junction table
CREATE TABLE license_nodes (
  id         INTEGER PRIMARY KEY AUTOINCREMENT,
  license_id INTEGER NOT NULL,
  node_id    INTEGER NOT NULL,
  created_at INTEGER NOT NULL DEFAULT (unixepoch()),
  FOREIGN KEY (license_id) REFERENCES licenses(id) ON DELETE CASCADE,
  FOREIGN KEY (node_id)    REFERENCES nodes(id)    ON DELETE CASCADE,
  UNIQUE(license_id, node_id)
);

CREATE INDEX idx_license_nodes_license_id ON license_nodes(license_id);
CREATE INDEX idx_license_nodes_node_id ON license_nodes(node_id);
