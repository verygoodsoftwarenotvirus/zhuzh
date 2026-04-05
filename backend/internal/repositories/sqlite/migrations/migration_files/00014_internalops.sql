CREATE TABLE IF NOT EXISTS queue_test_messages (
    id TEXT NOT NULL PRIMARY KEY,
    queue_name TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    acknowledged_at TEXT DEFAULT NULL
);
