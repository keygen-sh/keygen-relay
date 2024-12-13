-- re-add the old column
ALTER TABLE
  audit_logs
ADD
  COLUMN entity_type TEXT NOT NULL DEFAULT 'unknown';

-- revert data migration
UPDATE
  audit_logs
SET
  entity_type = CASE
    entity_type_id
    WHEN 1 THEN 'License'
    WHEN 2 THEN 'Node'
  END;

-- drop the new column
ALTER TABLE
  audit_logs DROP COLUMN entity_type_id;

-- drop the new table
DROP TABLE entity_types;
