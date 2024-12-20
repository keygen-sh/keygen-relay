-- rebuild the table with the old TEXT entity_id
CREATE TABLE _audit_logs (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  event_type_id TINYINT NOT NULL REFERENCES event_types (id),
  entity_type_id TINYINT NOT NULL REFERENCES entity_types (id),
  entity_id TEXT NOT NULL,
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
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
  entity_type_id,
  COALESCE(
    CASE
      entity_type_id
      WHEN 1 THEN (SELECT guid FROM licenses WHERE id = entity_id LIMIT 1)
      WHEN 2 THEN (SELECT fingerprint FROM nodes WHERE id = entity_id LIMIT 1)
      ELSE NULL
    END,
    'unknown'
  ) AS entity_id,
  created_at
FROM
  audit_logs;

-- drop the old table
DROP TABLE audit_logs;

-- replace with new table
ALTER TABLE
  _audit_logs RENAME TO audit_logs;
