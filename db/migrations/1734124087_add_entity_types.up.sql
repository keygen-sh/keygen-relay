-- create entities table
CREATE TABLE entity_types (id TINYINT PRIMARY KEY, name TEXT NOT NULL);

-- insert events
INSERT INTO
  entity_types (id, name)
VALUES
  (0, 'unknown'),
  (1, 'license'),
  (2, 'node');

-- rebuild the table with the new schema (this is a workaround for sqlite not supporting ALTER TABLE x ALTER COLUMN y NOT NULL)
CREATE TABLE _audit_logs (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  event_type_id TINYINT NOT NULL REFERENCES event_types (id),
  entity_type_id TINYINT NOT NULL REFERENCES entity_types (id),
  entity_id TEXT NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- copy data from the old table to the new
INSERT INTO
  _audit_logs (
    id,
    event_type_id,
    entity_type_id,
    entity_id,
    created_at
  )
SELECT
  id,
  event_type_id,
  CASE
    entity_type
    WHEN 'license' THEN 1
    WHEN 'node' THEN 2
    ELSE 0
  END AS entity_type_id,
  entity_id,
  created_at
FROM
  audit_logs;

-- drop the old table
DROP TABLE audit_logs;

-- replace with new table
ALTER TABLE
  _audit_logs RENAME TO audit_logs;
