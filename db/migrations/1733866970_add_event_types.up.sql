-- create events table
CREATE TABLE event_types (id TINYINT PRIMARY KEY, name TEXT NOT NULL);

-- insert events
INSERT INTO
  event_types (id, name)
VALUES
  (0, 'unknown'),
  (1, 'license.added'),
  (2, 'license.removed'),
  (3, 'license.leased'),
  (4, 'license.lease_extended'),
  (5, 'license.released'),
  (6, 'license.lease_expired'),
  (7, 'node.activated'),
  (8, 'node.heartbeat_ping'),
  (9, 'node.deactivated'),
  (10, 'node.culled');

-- rebuild the table with the new schema (this is a workaround for sqlite not supporting ALTER TABLE x ALTER COLUMN y NOT NULL)
CREATE TABLE _audit_logs (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  event_type_id TINYINT NOT NULL REFERENCES event_types (id),
  entity_type TEXT NOT NULL,
  entity_id TEXT NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- copy data from the old table to the new
INSERT INTO
  _audit_logs (
    id,
    event_type_id,
    entity_type,
    entity_id,
    created_at
  )
SELECT
  id,
  CASE
    action
    WHEN 'added' THEN 1
    WHEN 'removed' THEN 2
    WHEN 'claimed' THEN 3
    WHEN 'released' THEN 5
    WHEN 'activated' THEN 7
    WHEN 'ping' THEN 8
    WHEN 'culled' THEN 10
    ELSE 0
  END AS event_type_id,
  entity_type,
  entity_id,
  created_at
FROM
  audit_logs;

-- drop the old table
DROP TABLE audit_logs;

-- replace with new table
ALTER TABLE
  _audit_logs RENAME TO audit_logs;
