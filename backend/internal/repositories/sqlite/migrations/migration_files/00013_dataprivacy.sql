-- Data Privacy Domain Migration
-- GDPR/CCPA compliance tracking for user data disclosures

CREATE TABLE IF NOT EXISTS user_data_disclosures (
    id TEXT NOT NULL PRIMARY KEY,
    belongs_to_user TEXT NOT NULL REFERENCES users("id") ON DELETE CASCADE,
    status TEXT NOT NULL DEFAULT 'pending' CHECK(status IN ('pending', 'processing', 'completed', 'failed', 'expired')),
    report_id TEXT,
    expires_at TEXT NOT NULL,
    created_at TEXT DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')) NOT NULL,
    last_updated_at TEXT,
    completed_at TEXT,
    archived_at TEXT
);

-- =============================================================================
-- INDEXES FOR DATA PRIVACY TABLES
-- =============================================================================

-- User data disclosures indexes
CREATE INDEX idx_user_data_disclosures_user ON user_data_disclosures (belongs_to_user);
CREATE INDEX idx_user_data_disclosures_status ON user_data_disclosures (status);
CREATE INDEX idx_user_data_disclosures_user_status ON user_data_disclosures (belongs_to_user, status);
CREATE INDEX idx_user_data_disclosures_expires_at ON user_data_disclosures (expires_at);
CREATE INDEX idx_user_data_disclosures_archived_at ON user_data_disclosures (archived_at);
