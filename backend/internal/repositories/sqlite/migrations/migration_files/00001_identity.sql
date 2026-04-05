-- Identity Domain Migration
-- Core user management, accounts, and authentication sessions

CREATE TABLE IF NOT EXISTS users (
    id TEXT NOT NULL PRIMARY KEY,
    username TEXT NOT NULL,
    email_address TEXT NOT NULL,
    hashed_password TEXT NOT NULL,
    password_last_changed_at TEXT,
    requires_password_change BOOLEAN DEFAULT FALSE NOT NULL,
    two_factor_secret TEXT NOT NULL,
    two_factor_secret_verified_at TEXT,
    service_role TEXT DEFAULT 'service_user' NOT NULL,
    user_account_status TEXT DEFAULT 'unverified' NOT NULL,
    user_account_status_explanation TEXT DEFAULT '' NOT NULL,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    last_updated_at TEXT,
    archived_at TEXT,
    birthday TEXT,
    email_address_verification_token TEXT DEFAULT '',
    email_address_verified_at TEXT,
    first_name TEXT DEFAULT '' NOT NULL,
    last_name TEXT DEFAULT '' NOT NULL,
    last_accepted_terms_of_service TEXT,
    last_accepted_privacy_policy TEXT,
    last_indexed_at TEXT,
    UNIQUE(username)
);

CREATE TABLE IF NOT EXISTS accounts (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    billing_status TEXT DEFAULT 'unpaid' NOT NULL,
    contact_phone TEXT DEFAULT '' NOT NULL,
    payment_processor_customer_id TEXT DEFAULT '' NOT NULL,
    subscription_plan_id TEXT,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    last_updated_at TEXT,
    archived_at TEXT,
    belongs_to_user TEXT NOT NULL REFERENCES users("id") ON DELETE CASCADE,
    time_zone TEXT DEFAULT 'US/Central' NOT NULL CHECK(time_zone IN ('UTC', 'US/Pacific', 'US/Mountain', 'US/Central', 'US/Eastern')),
    address_line_1 TEXT DEFAULT '' NOT NULL,
    address_line_2 TEXT DEFAULT '' NOT NULL,
    city TEXT DEFAULT '' NOT NULL,
    state TEXT DEFAULT '' NOT NULL,
    zip_code TEXT DEFAULT '' NOT NULL,
    country TEXT DEFAULT '' NOT NULL,
    latitude REAL,
    longitude REAL,
    last_payment_provider_sync_occurred_at TEXT,
    webhook_hmac_secret TEXT DEFAULT '' NOT NULL,
    UNIQUE(belongs_to_user, name)
);

CREATE TABLE IF NOT EXISTS account_user_memberships (
    id TEXT NOT NULL PRIMARY KEY,
    belongs_to_account TEXT NOT NULL REFERENCES accounts("id") ON DELETE CASCADE,
    belongs_to_user TEXT NOT NULL REFERENCES users("id") ON DELETE CASCADE,
    default_account BOOLEAN DEFAULT FALSE NOT NULL,
    account_role TEXT DEFAULT 'account_user' NOT NULL,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    last_updated_at TEXT,
    archived_at TEXT,
    UNIQUE(belongs_to_account, belongs_to_user)
);

CREATE TABLE IF NOT EXISTS account_invitations (
    id TEXT NOT NULL PRIMARY KEY,
    destination_account TEXT NOT NULL REFERENCES accounts("id") ON DELETE CASCADE,
    to_email TEXT NOT NULL,
    to_user TEXT  REFERENCES users("id") ON DELETE CASCADE,
    from_user TEXT NOT NULL  REFERENCES users("id") ON DELETE CASCADE,
    status TEXT DEFAULT 'pending' NOT NULL CHECK(status IN ('pending', 'cancelled', 'accepted', 'rejected')),
    note TEXT DEFAULT '' NOT NULL,
    status_note TEXT DEFAULT '' NOT NULL,
    token TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    last_updated_at TEXT,
    archived_at TEXT,
    expires_at TEXT DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now', '+7 days')) NOT NULL,
    to_name TEXT DEFAULT '' NOT NULL,
    UNIQUE(to_user, to_email, from_user, destination_account)
);

CREATE TABLE IF NOT EXISTS sessions (
    token TEXT PRIMARY KEY,
    data BLOB NOT NULL,
    expiry TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    UNIQUE(token)
);

-- =============================================================================
-- INDEXES FOR IDENTITY TABLES
-- =============================================================================

-- Users table indexes
CREATE INDEX idx_users_archived_at ON users (archived_at) WHERE archived_at IS NULL;
CREATE INDEX idx_users_email_address_active ON users (email_address) WHERE archived_at IS NULL;
CREATE INDEX idx_users_username_active ON users (username) WHERE archived_at IS NULL;
CREATE INDEX idx_users_email_verification_token ON users (email_address_verification_token) WHERE archived_at IS NULL AND email_address_verification_token != '';
CREATE INDEX idx_users_service_role_username ON users (service_role, username) WHERE archived_at IS NULL;
CREATE INDEX idx_users_two_factor_verified ON users (two_factor_secret_verified_at) WHERE archived_at IS NULL;
CREATE INDEX idx_users_indexing_status ON users (last_indexed_at) WHERE archived_at IS NULL;
CREATE INDEX idx_users_active_created_at ON users (created_at) WHERE archived_at IS NULL;
CREATE INDEX idx_users_active_updated_at ON users (last_updated_at) WHERE archived_at IS NULL;
CREATE INDEX idx_users_indexing_needed ON users (last_indexed_at) WHERE archived_at IS NULL;

-- Accounts table indexes
CREATE INDEX idx_accounts_belongs_to_user ON accounts (belongs_to_user) WHERE archived_at IS NULL;
CREATE INDEX idx_accounts_archived_at ON accounts (archived_at) WHERE archived_at IS NULL;
CREATE INDEX idx_accounts_user_name ON accounts (belongs_to_user, name) WHERE archived_at IS NULL;
CREATE INDEX idx_accounts_payment_sync ON accounts (last_payment_provider_sync_occurred_at) WHERE archived_at IS NULL;
CREATE INDEX idx_accounts_billing_status ON accounts (billing_status) WHERE archived_at IS NULL;
CREATE INDEX idx_accounts_user_created_at ON accounts (belongs_to_user, created_at) WHERE archived_at IS NULL;
CREATE INDEX idx_accounts_user_updated_at ON accounts (belongs_to_user, last_updated_at) WHERE archived_at IS NULL;
CREATE INDEX idx_accounts_user_billing ON accounts (belongs_to_user, billing_status) WHERE archived_at IS NULL;

-- Account user memberships indexes
CREATE INDEX idx_memberships_user ON account_user_memberships (belongs_to_user) WHERE archived_at IS NULL;
CREATE INDEX idx_memberships_account ON account_user_memberships (belongs_to_account) WHERE archived_at IS NULL;
CREATE INDEX idx_memberships_default_account ON account_user_memberships (belongs_to_user, default_account) WHERE archived_at IS NULL AND default_account = TRUE;
CREATE INDEX idx_memberships_user_account ON account_user_memberships (belongs_to_user, belongs_to_account) WHERE archived_at IS NULL;
CREATE INDEX idx_memberships_account_role ON account_user_memberships (belongs_to_account, account_role) WHERE archived_at IS NULL;

-- Account invitations indexes
CREATE INDEX idx_invitations_destination_account ON account_invitations (destination_account) WHERE archived_at IS NULL;
CREATE INDEX idx_invitations_from_user ON account_invitations (from_user) WHERE archived_at IS NULL;
CREATE INDEX idx_invitations_to_user ON account_invitations (to_user) WHERE archived_at IS NULL;
CREATE INDEX idx_invitations_to_email ON account_invitations (to_email) WHERE archived_at IS NULL;
CREATE INDEX idx_invitations_token ON account_invitations (token) WHERE archived_at IS NULL;
CREATE INDEX idx_invitations_status ON account_invitations (status) WHERE archived_at IS NULL;
CREATE INDEX idx_invitations_expires_at ON account_invitations (expires_at) WHERE archived_at IS NULL;

-- Sessions indexes
CREATE INDEX idx_sessions_expiry ON sessions (expiry);
CREATE INDEX idx_sessions_created_at ON sessions (created_at);

-- Text search indexes (for efficient LIKE and ILIKE operations)
-- Uncomment if pg_trgm extension is available:
-- CREATE EXTENSION IF NOT EXISTS pg_trgm;
-- CREATE INDEX idx_users_username_trgm ON users USING gin (username gin_trgm_ops) WHERE archived_at IS NULL;
