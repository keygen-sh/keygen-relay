CREATE TABLE nodes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    fingerprint TEXT UNIQUE NOT NULL,
    claimed_at TEXT,
    last_heartbeat_at TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE licenses (
    id TEXT PRIMARY KEY,
    file BLOB UNIQUE NOT NULL,
    key TEXT UNIQUE NOT NULL,
    claims INTEGER DEFAULT 0 NOT NULL,
    last_claimed_at TEXT,
    last_released_at TEXT,
    node_id INTEGER UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (node_id) REFERENCES nodes(id) ON DELETE SET NULL
);

CREATE TABLE audit_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    action TEXT NOT NULL,
    entity_type TEXT NOT NULL,
    entity_id TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
