-- Comments Domain Migration
-- Generic polymorphic comments for platform entities

CREATE TABLE IF NOT EXISTS comments (
    id TEXT NOT NULL PRIMARY KEY,
    content TEXT NOT NULL DEFAULT '',
    target_type TEXT NOT NULL CHECK(target_type IN ('issue_reports')),
    referenced_id TEXT NOT NULL,
    parent_comment_id TEXT REFERENCES comments("id") ON DELETE CASCADE,
    belongs_to_user TEXT NOT NULL REFERENCES users("id") ON DELETE CASCADE,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    last_updated_at TEXT,
    archived_at TEXT
);

CREATE INDEX idx_comments_reference ON comments (target_type, referenced_id) WHERE archived_at IS NULL;
CREATE INDEX idx_comments_user ON comments (belongs_to_user) WHERE archived_at IS NULL;
CREATE INDEX idx_comments_parent ON comments (parent_comment_id) WHERE archived_at IS NULL;
