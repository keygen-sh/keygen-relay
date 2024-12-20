ALTER TABLE
  nodes
ADD
  COLUMN deactivated_at INTEGER;

CREATE INDEX idx_nodes_activated ON nodes (deactivated_at)
WHERE
  deactivated_at IS NULL;
