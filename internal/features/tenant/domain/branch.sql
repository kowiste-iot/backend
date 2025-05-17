-- Branch-specific schema template
-- This will be used to create a new schema for each branch

-- Create schema if not exists
CREATE SCHEMA IF NOT EXISTS "{{.SchemaName}}";

-- Set search path to the new schema
SET search_path TO "{{.SchemaName}}";

-- Branch metadata - Stores basic branch and tenant information locally
CREATE TABLE IF NOT EXISTS branch_info (
    tenant_id TEXT NOT NULL,
    branch_id TEXT NOT NULL,
    tenant_name TEXT NOT NULL,
    branch_name TEXT NOT NULL,
    description TEXT,
    timezone TEXT DEFAULT 'UTC',
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (tenant_id, branch_id)
);

-- Assets definition - MUST BE CREATED FIRST as other tables reference it
CREATE TABLE IF NOT EXISTS assets (
    id TEXT PRIMARY KEY,
    parent TEXT REFERENCES assets(id),
    name TEXT NOT NULL,
    description TEXT,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);  

CREATE INDEX IF NOT EXISTS idx_assets_deleted_at ON assets(deleted_at);
CREATE INDEX IF NOT EXISTS idx_assets_parent ON assets(parent);
CREATE UNIQUE INDEX IF NOT EXISTS idx_assets_name ON assets(name) WHERE deleted_at IS NULL;

-- Asset dependencies definition
CREATE TABLE IF NOT EXISTS asset_dependencies (
    asset_id TEXT NOT NULL REFERENCES assets(id),
    feature_id TEXT NOT NULL,
    feature TEXT NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    PRIMARY KEY (asset_id, feature_id, feature)
);

CREATE INDEX IF NOT EXISTS idx_asset_dependencies_deleted_at ON asset_dependencies(deleted_at);

-- Actions definition
CREATE TABLE IF NOT EXISTS actions (
    id TEXT PRIMARY KEY,
    parent TEXT REFERENCES assets(id),
    name TEXT NOT NULL,
    enabled BOOLEAN DEFAULT TRUE,
    description TEXT,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_actions_deleted_at ON actions(deleted_at);
CREATE INDEX IF NOT EXISTS idx_actions_parent ON actions(parent);
CREATE UNIQUE INDEX IF NOT EXISTS idx_actions_name ON actions(name) WHERE deleted_at IS NULL;

-- Alerts definition
CREATE TABLE IF NOT EXISTS alerts (
    id TEXT PRIMARY KEY,
    parent TEXT REFERENCES assets(id),
    name TEXT NOT NULL,
    enabled BOOLEAN DEFAULT TRUE,
    description TEXT,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_alerts_deleted_at ON alerts(deleted_at);
CREATE INDEX IF NOT EXISTS idx_alerts_parent ON alerts(parent);
CREATE UNIQUE INDEX IF NOT EXISTS idx_alerts_name ON alerts(name) WHERE deleted_at IS NULL;

-- Dashboards definition
CREATE TABLE IF NOT EXISTS dashboards (
    id TEXT PRIMARY KEY,
    parent TEXT REFERENCES assets(id),
    name TEXT NOT NULL,
    description TEXT,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_dashboards_deleted_at ON dashboards(deleted_at);
CREATE INDEX IF NOT EXISTS idx_dashboards_parent ON dashboards(parent);
CREATE UNIQUE INDEX IF NOT EXISTS idx_dashboards_name ON dashboards(name) WHERE deleted_at IS NULL;

-- Devices definition
CREATE TABLE IF NOT EXISTS devices (
    id TEXT PRIMARY KEY,
    parent TEXT REFERENCES assets(id),
    name TEXT NOT NULL,
    description TEXT,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_devices_deleted_at ON devices(deleted_at);
CREATE INDEX IF NOT EXISTS idx_devices_parent ON devices(parent);
CREATE UNIQUE INDEX IF NOT EXISTS idx_devices_name ON devices(name) WHERE deleted_at IS NULL;

-- Measures definition
CREATE TABLE IF NOT EXISTS measures (
    id TEXT PRIMARY KEY,
    parent TEXT REFERENCES assets(id),
    name TEXT NOT NULL,
    description TEXT,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_measures_deleted_at ON measures(deleted_at);
CREATE INDEX IF NOT EXISTS idx_measures_parent ON measures(parent);
CREATE UNIQUE INDEX IF NOT EXISTS idx_measures_name ON measures(name) WHERE deleted_at IS NULL;

-- Message models definition
CREATE TABLE IF NOT EXISTS message_models (
    id TEXT PRIMARY KEY,
    topic TEXT NOT NULL,
    data BYTEA,
    event TEXT,
    status TEXT,
    timestamp TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_message_models_topic ON message_models(topic);
CREATE INDEX IF NOT EXISTS idx_message_models_timestamp ON message_models(timestamp);

-- Message store definition
CREATE TABLE IF NOT EXISTS message_store (
    id TEXT PRIMARY KEY,
    measure_id TEXT NOT NULL REFERENCES measures(id),
    time TIMESTAMP WITH TIME ZONE NOT NULL,
    data JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_measure_time ON message_store(measure_id, time);
CREATE INDEX IF NOT EXISTS idx_message_store_time ON message_store(time);

-- Users definition
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    auth_user_id TEXT NOT NULL,
    email TEXT NOT NULL,
    first_name TEXT,
    last_name TEXT,
    name TEXT,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_auth_user_id ON users(auth_user_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email) WHERE deleted_at IS NULL;


--  Roles definition
CREATE TABLE IF NOT EXISTS roles (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_roles_name ON roles(name) WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS user_roles (
    user_id TEXT NOT NULL REFERENCES users(id),
    role_id TEXT NOT NULL REFERENCES roles(id),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    PRIMARY KEY (user_id, role_id)
);

CREATE INDEX IF NOT EXISTS idx_user_roles_deleted_at ON user_roles(deleted_at);

-- Widgets definition
CREATE TABLE IF NOT EXISTS widgets (
    id TEXT PRIMARY KEY,
    dashboard_id TEXT NOT NULL REFERENCES dashboards(id),
    type_widget INTEGER NOT NULL,
    x INTEGER NOT NULL,
    y INTEGER NOT NULL,
    w INTEGER NOT NULL,
    h INTEGER NOT NULL,
    label TEXT,
    show_label BOOLEAN DEFAULT FALSE,
    show_emotion BOOLEAN DEFAULT FALSE,
    true_emotion BOOLEAN DEFAULT FALSE,
    options JSONB,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_widgets_deleted_at ON widgets(deleted_at);
CREATE INDEX IF NOT EXISTS idx_widgets_dashboard_id ON widgets(dashboard_id);

-- Widget link data definition
CREATE TABLE IF NOT EXISTS widget_link_data (
    id SERIAL PRIMARY KEY,
    widget_id TEXT NOT NULL REFERENCES widgets(id),
    measure TEXT NOT NULL,
    tag TEXT NOT NULL,
    legend TEXT NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_widget_link_data_widget_id ON widget_link_data(widget_id);
CREATE INDEX IF NOT EXISTS idx_widget_link_data_measure ON widget_link_data(measure);

-- Database version tracking
CREATE TABLE IF NOT EXISTS schema_version (
    version TEXT PRIMARY KEY,
    applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    description TEXT
);

-- Only insert version if it doesn't exist
INSERT INTO schema_version (version, description) 
SELECT '1.0.0', 'Initial branch schema template'
WHERE NOT EXISTS (SELECT 1 FROM schema_version WHERE version = '1.0.0');