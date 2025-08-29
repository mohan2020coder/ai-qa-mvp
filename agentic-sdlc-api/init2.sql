-- Enable UUID support
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============= USERS ====================
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    name VARCHAR(255),
    role VARCHAR(50) DEFAULT 'user', -- 'user', 'admin', 'owner'
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- ============= WORKFLOWS ====================
CREATE TABLE workflows (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    is_active BOOLEAN DEFAULT true,
    settings JSONB,
    version INT DEFAULT 1,
    is_template BOOLEAN DEFAULT false,
    forked_from UUID REFERENCES workflows(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_workflows_created_by ON workflows(created_by);

-- ============= WORKFLOW SHARING ====================
CREATE TABLE workflow_shares (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workflow_id UUID NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    access_level VARCHAR(50) NOT NULL CHECK (access_level IN ('view', 'edit', 'owner')),
    created_at TIMESTAMPTZ DEFAULT now(),
    UNIQUE (workflow_id, user_id)
);

-- ============= NODES ====================
CREATE TABLE nodes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workflow_id UUID NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
    type VARCHAR(100) NOT NULL,
    name VARCHAR(255),
    position_x FLOAT NOT NULL,
    position_y FLOAT NOT NULL,
    parameters JSONB NOT NULL,
    credentials_id UUID REFERENCES credentials(id),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_nodes_workflow ON nodes(workflow_id);

-- ============= CONNECTIONS ====================
CREATE TABLE connections (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workflow_id UUID NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
    from_node_id UUID NOT NULL REFERENCES nodes(id) ON DELETE CASCADE,
    to_node_id UUID NOT NULL REFERENCES nodes(id) ON DELETE CASCADE,
    output_index INT DEFAULT 0,
    input_index INT DEFAULT 0
);

-- ============= CREDENTIALS ====================
CREATE TABLE credentials (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    type VARCHAR(100) NOT NULL,
    data BYTEA NOT NULL, -- Encrypted binary data
    iv BYTEA NOT NULL,   -- Initialization vector for encryption
    tag BYTEA,           -- Auth tag (for AES-GCM)
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_credentials_owner ON credentials(owner_id);

-- ============= EXECUTIONS ====================
CREATE TABLE executions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workflow_id UUID NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL CHECK (status IN ('success', 'error', 'running', 'cancelled')),
    started_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    ended_at TIMESTAMPTZ,
    error_message TEXT,
    execution_data JSONB NOT NULL,
    run_by_user_id UUID REFERENCES users(id)
);

CREATE INDEX idx_executions_workflow ON executions(workflow_id);

-- ============= TRIGGERS ====================
CREATE TABLE triggers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workflow_id UUID NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL CHECK (type IN ('cron', 'webhook')),
    config JSONB NOT NULL,
    is_enabled BOOLEAN DEFAULT true
);

-- ============= WEBHOOKS ====================
CREATE TABLE webhooks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    trigger_id UUID NOT NULL REFERENCES triggers(id) ON DELETE CASCADE,
    path VARCHAR(255) UNIQUE NOT NULL,
    method VARCHAR(10) CHECK (method IN ('GET', 'POST', 'PUT', 'DELETE')),
    secret VARCHAR(255)
);

-- ============= VARIABLES (WORKFLOW-SCOPED) ====================
CREATE TABLE workflow_variables (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workflow_id UUID NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
    key VARCHAR(100) NOT NULL,
    value TEXT NOT NULL,
    is_secret BOOLEAN DEFAULT false,
    UNIQUE (workflow_id, key)
);

-- ============= VERSIONS (Workflow Snapshots) ====================
CREATE TABLE workflow_versions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workflow_id UUID NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
    version INT NOT NULL,
    definition JSONB NOT NULL, -- Full snapshot of nodes/connections/settings
    created_at TIMESTAMPTZ DEFAULT now(),
    UNIQUE (workflow_id, version)
);

-- ============= TAGS ====================
CREATE TABLE tags (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE workflow_tags (
    workflow_id UUID NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
    tag_id UUID NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (workflow_id, tag_id)
);

-- ============= AUDIT LOGS ====================
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID NOT NULL,
    metadata JSONB,
    timestamp TIMESTAMPTZ DEFAULT now()
);

-- ============= INDEXING STRATEGIES ====================
-- Helpful for filtering and performance
CREATE INDEX idx_executions_status ON executions(status);
CREATE INDEX idx_workflow_variables_key ON workflow_variables(key);
CREATE INDEX idx_audit_logs_entity ON audit_logs(entity_type, entity_id);

-- ============= SECURITY PLACEHOLDERS ====================
-- NOTE: Actual encryption/decryption should be done in application layer.
-- Optionally use PostgreSQL's pgcrypto for basic encryption.

-- ============= OPTIONAL RLS (Row Level Security) ====================
-- To enforce per-user access on workflows:
-- ALTER TABLE workflows ENABLE ROW LEVEL SECURITY;
-- CREATE POLICY user_own_workflows ON workflows
--   USING (created_by = current_setting('app.current_user')::uuid);

-- Then in app layer, set:
-- SET app.current_user = 'user-id-here';

