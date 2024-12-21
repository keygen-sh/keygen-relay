-- rebuild the table with the new INTEGER entity_id
CREATE TABLE _audit_logs (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  event_type_id TINYINT NOT NULL REFERENCES event_types (id),
  entity_type_id TINYINT NOT NULL REFERENCES entity_types (id),
  entity_id INTEGER NOT NULL,
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
  audit_logs.id,
  audit_logs.event_type_id,
  audit_logs.entity_type_id,
  COALESCE(
    CASE
      audit_logs.entity_type_id
      WHEN 1 THEN licenses.id
      WHEN 2 THEN nodes.id
      ELSE NULL
    END,
    0 -- default
  ) AS entity_id,
  audit_logs.created_at
FROM
  audit_logs
  LEFT JOIN licenses ON licenses.guid = audit_logs.entity_id
  LEFT JOIN nodes ON nodes.fingerprint = audit_logs.entity_id;

-- drop the old table
DROP TABLE audit_logs;

-- replace with new table
ALTER TABLE
  _audit_logs RENAME TO audit_logs;
