-- re-add the old column
ALTER TABLE
  audit_logs
ADD
  COLUMN action TEXT NOT NULL DEFAULT 'unknown';

-- revert data migration
UPDATE
  audit_logs
SET
  action = CASE
    event_type_id
    WHEN 1 THEN 'added'
    WHEN 2 THEN 'removed'
    WHEN 3 THEN 'claimed'
    WHEN 5 THEN 'released'
    WHEN 7 THEN 'activated'
    WHEN 8 THEN 'ping'
    WHEN 10 THEN 'culled'
    ELSE 'unknown'
  END;

-- drop the new column
ALTER TABLE
  audit_logs DROP COLUMN event_type_id;

-- drop the new table
DROP TABLE event_types;
