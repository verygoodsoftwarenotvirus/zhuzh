CREATE TABLE IF NOT EXISTS uploaded_media (
    id TEXT NOT NULL PRIMARY KEY,
    storage_path TEXT NOT NULL,
    mime_type TEXT NOT NULL CHECK(mime_type IN ('image/png', 'image/jpeg', 'image/gif', 'video/mp4')),
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    last_updated_at TEXT,
    archived_at TEXT,
    created_by_user TEXT NOT NULL REFERENCES users("id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS user_avatars (
    id TEXT NOT NULL PRIMARY KEY,
    belongs_to_user TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    uploaded_media_id TEXT NOT NULL REFERENCES uploaded_media(id) ON DELETE CASCADE,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    archived_at TEXT,
    UNIQUE(belongs_to_user, archived_at)
);
