-- Webhooks Domain Migration
-- Webhook management, trigger event catalog, and trigger configs

-- Catalog of available webhook trigger events (first-class entity)
CREATE TABLE IF NOT EXISTS webhook_trigger_events (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    last_updated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    archived_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

CREATE TYPE webhook_content_type AS ENUM (
    'application/json',
    'application/xml'
);

CREATE TYPE webhook_method AS ENUM (
    'GET',
    'PUT',
    'PATCH',
    'POST',
    'DELETE'
);

CREATE TABLE IF NOT EXISTS webhooks (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    content_type webhook_content_type NOT NULL DEFAULT 'application/json',
    url TEXT NOT NULL,
    method webhook_method NOT NULL DEFAULT 'POST',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    created_by_user TEXT NOT NULL REFERENCES users("id") ON DELETE CASCADE,
    belongs_to_account TEXT NOT NULL REFERENCES accounts("id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS webhook_trigger_configs (
    id TEXT NOT NULL PRIMARY KEY,
    trigger_event TEXT NOT NULL REFERENCES webhook_trigger_events("id") ON DELETE CASCADE,
    belongs_to_webhook TEXT NOT NULL REFERENCES webhooks("id") ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    archived_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(trigger_event, belongs_to_webhook, archived_at)
);

-- =============================================================================
-- INDEXES FOR WEBHOOKS TABLES
-- =============================================================================

-- Webhook trigger events (catalog) indexes
CREATE INDEX idx_webhook_trigger_events_archived_at ON webhook_trigger_events (archived_at) WHERE archived_at IS NULL;

-- Webhooks indexes
CREATE INDEX idx_webhooks_account ON webhooks (belongs_to_account) WHERE archived_at IS NULL;
CREATE INDEX idx_webhooks_archived_at ON webhooks (archived_at) WHERE archived_at IS NULL;
CREATE INDEX idx_webhooks_created_by_user ON webhooks (created_by_user) WHERE archived_at IS NULL;

-- Webhook trigger configs indexes
CREATE INDEX idx_webhook_trigger_configs_webhook ON webhook_trigger_configs (belongs_to_webhook) WHERE archived_at IS NULL;
CREATE INDEX idx_webhook_trigger_configs_event ON webhook_trigger_configs (trigger_event) WHERE archived_at IS NULL;
