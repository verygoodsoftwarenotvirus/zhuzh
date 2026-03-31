-- Payments Domain Migration
-- Products, subscriptions, purchases, and payment transactions

-- =============================================================================
-- ENUMERATED TYPES
-- =============================================================================

CREATE TYPE product_kind AS ENUM (
    'recurring',
    'one_time'
);

CREATE TYPE subscription_status AS ENUM (
    'active',
    'cancelled',
    'past_due',
    'trialing',
    'incomplete'
);

CREATE TYPE payment_transaction_status AS ENUM (
    'succeeded',
    'failed',
    'pending',
    'refunded'
);

-- =============================================================================
-- TABLES
-- =============================================================================

CREATE TABLE IF NOT EXISTS products (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    kind product_kind NOT NULL,
    amount_cents INTEGER NOT NULL,
    currency TEXT NOT NULL DEFAULT 'usd',
    billing_interval_months INTEGER,
    external_product_id TEXT DEFAULT '',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS subscriptions (
    id TEXT NOT NULL PRIMARY KEY,
    belongs_to_account TEXT NOT NULL REFERENCES accounts("id") ON DELETE CASCADE,
    product_id TEXT NOT NULL REFERENCES products("id") ON DELETE CASCADE,
    external_subscription_id TEXT DEFAULT '',
    status subscription_status NOT NULL DEFAULT 'active',
    current_period_start TIMESTAMP WITH TIME ZONE NOT NULL,
    current_period_end TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS purchases (
    id TEXT NOT NULL PRIMARY KEY,
    belongs_to_account TEXT NOT NULL REFERENCES accounts("id") ON DELETE CASCADE,
    product_id TEXT NOT NULL REFERENCES products("id") ON DELETE CASCADE,
    amount_cents INTEGER NOT NULL,
    currency TEXT NOT NULL DEFAULT 'usd',
    completed_at TIMESTAMP WITH TIME ZONE,
    external_transaction_id TEXT DEFAULT '',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS payment_transactions (
    id TEXT NOT NULL PRIMARY KEY,
    belongs_to_account TEXT NOT NULL REFERENCES accounts("id") ON DELETE CASCADE,
    subscription_id TEXT REFERENCES subscriptions("id") ON DELETE SET NULL,
    purchase_id TEXT REFERENCES purchases("id") ON DELETE SET NULL,
    external_transaction_id TEXT NOT NULL DEFAULT '',
    amount_cents INTEGER NOT NULL,
    currency TEXT NOT NULL DEFAULT 'usd',
    status payment_transaction_status NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
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
