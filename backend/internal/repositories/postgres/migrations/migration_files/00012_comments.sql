-- Comments Domain Migration
-- Generic polymorphic comments for platform entities

CREATE TYPE comment_target_type AS ENUM (
    'issue_reports'
);

CREATE TABLE IF NOT EXISTS comments (
    id TEXT NOT NULL PRIMARY KEY,
    content TEXT NOT NULL DEFAULT '',
    target_type comment_target_type NOT NULL,
    referenced_id TEXT NOT NULL,
    parent_comment_id TEXT REFERENCES comments("id") ON DELETE CASCADE,
    belongs_to_user TEXT NOT NULL REFERENCES users("id") ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_comments_reference ON comments (target_type, referenced_id) WHERE archived_at IS NULL;
CREATE INDEX idx_comments_user ON comments (belongs_to_user) WHERE archived_at IS NULL;
CREATE INDEX idx_comments_parent ON comments (parent_comment_id) WHERE archived_at IS NULL;
