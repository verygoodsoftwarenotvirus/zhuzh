-- Payments Domain Migration
-- Products, subscriptions, purchases, and payment transactions

-- =============================================================================
-- TABLES
-- =============================================================================

CREATE TABLE IF NOT EXISTS products (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    kind TEXT NOT NULL CHECK(kind IN ('recurring', 'one_time')),
    amount_cents INTEGER NOT NULL,
    currency TEXT NOT NULL DEFAULT 'usd',
    billing_interval_months INTEGER,
    external_product_id TEXT DEFAULT '',
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    last_updated_at TEXT,
    archived_at TEXT
);

CREATE TABLE IF NOT EXISTS subscriptions (
    id TEXT NOT NULL PRIMARY KEY,
    belongs_to_account TEXT NOT NULL REFERENCES accounts("id") ON DELETE CASCADE,
    product_id TEXT NOT NULL REFERENCES products("id") ON DELETE CASCADE,
    external_subscription_id TEXT DEFAULT '',
    status TEXT NOT NULL DEFAULT 'active' CHECK(status IN ('active', 'cancelled', 'past_due', 'trialing', 'incomplete')),
    current_period_start TEXT NOT NULL,
    current_period_end TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    last_updated_at TEXT,
    archived_at TEXT
);

CREATE TABLE IF NOT EXISTS purchases (
    id TEXT NOT NULL PRIMARY KEY,
    belongs_to_account TEXT NOT NULL REFERENCES accounts("id") ON DELETE CASCADE,
    product_id TEXT NOT NULL REFERENCES products("id") ON DELETE CASCADE,
    amount_cents INTEGER NOT NULL,
    currency TEXT NOT NULL DEFAULT 'usd',
    completed_at TEXT,
    external_transaction_id TEXT DEFAULT '',
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    last_updated_at TEXT,
    archived_at TEXT
);

CREATE TABLE IF NOT EXISTS payment_transactions (
    id TEXT NOT NULL PRIMARY KEY,
    belongs_to_account TEXT NOT NULL REFERENCES accounts("id") ON DELETE CASCADE,
    subscription_id TEXT REFERENCES subscriptions("id") ON DELETE SET NULL,
    purchase_id TEXT REFERENCES purchases("id") ON DELETE SET NULL,
    external_transaction_id TEXT NOT NULL DEFAULT '',
    amount_cents INTEGER NOT NULL,
    currency TEXT NOT NULL DEFAULT 'usd',
    status TEXT NOT NULL CHECK(status IN ('succeeded', 'failed', 'pending', 'refunded')),
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

-- =============================================================================
-- INDEXES FOR PAYMENTS TABLES
-- =============================================================================

CREATE INDEX idx_products_archived_at ON products (archived_at) WHERE archived_at IS NULL;

CREATE INDEX idx_subscriptions_belongs_to_account ON subscriptions (belongs_to_account) WHERE archived_at IS NULL;
CREATE INDEX idx_subscriptions_status ON subscriptions (status) WHERE archived_at IS NULL;
CREATE INDEX idx_subscriptions_archived_at ON subscriptions (archived_at) WHERE archived_at IS NULL;

CREATE INDEX idx_purchases_belongs_to_account ON purchases (belongs_to_account) WHERE archived_at IS NULL;
CREATE INDEX idx_purchases_archived_at ON purchases (archived_at) WHERE archived_at IS NULL;

CREATE INDEX idx_payment_transactions_belongs_to_account ON payment_transactions (belongs_to_account);
CREATE INDEX idx_payment_transactions_created_at ON payment_transactions (belongs_to_account, created_at);
