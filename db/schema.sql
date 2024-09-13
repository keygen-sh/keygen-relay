CREATE TABLE nodes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    fingerprint TEXT UNIQUE NOT NULL,
    claimed_at TEXT,
    last_heartbeat_at TEXT
);

CREATE TABLE licenses (
    id TEXT PRIMARY KEY,
    file BLOB UNIQUE NOT NULL,
    key TEXT UNIQUE NOT NULL,
    claims INTEGER DEFAULT 0,
    last_claimed_at TEXT,
    last_released_at TEXT,
    node_id INTEGER UNIQUE,
    FOREIGN KEY (node_id) REFERENCES nodes(id)
);

CREATE TABLE audit_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    action TEXT NOT NULL,
    entity_type TEXT NOT NULL,
    entity_id TEXT NOT NULL,
    timestamp TEXT DEFAULT CURRENT_TIMESTAMP
);
