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
    WHEN 4 THEN 'released'
    WHEN 5 THEN 'activated'
    WHEN 6 THEN 'ping'
    WHEN 7 THEN 'culled'
  END;

-- drop the new column
ALTER TABLE
  audit_logs DROP COLUMN event_type_id;

-- drop the new table
DROP TABLE event_types;
