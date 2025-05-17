#!/bin/bash

# Continue even if there are errors
set -u

function create_database() {
    local database=$1
    echo "Creating database '$database'"
    psql -v ON_ERROR_STOP=0 --username "$POSTGRES_USER" <<-EOSQL
        CREATE DATABASE $database TEMPLATE template1;
        GRANT ALL PRIVILEGES ON DATABASE $database TO $POSTGRES_USER;
EOSQL
    echo "Finished creating database '$database' (ignoring any errors)"
}

function create_schema() {
    local database=$1
    local schema=$2
    echo "Creating schema '$schema' in database '$database'"
    psql -v ON_ERROR_STOP=0 --username "$POSTGRES_USER" --dbname="$database" <<-EOSQL
        CREATE SCHEMA IF NOT EXISTS $schema;
EOSQL
    echo "Finished creating schema '$schema' in database '$database'"
}

function create_schema_with_user() {
    local database=$1
    local schema=$2
    local user=$3
    local password=$4
    echo "Creating schema '$schema' in database '$database' with user '$user'"

    # Create schema
    psql -v ON_ERROR_STOP=0 --username "$POSTGRES_USER" --dbname="$database" <<-EOSQL
        CREATE SCHEMA IF NOT EXISTS $schema;
EOSQL

    # Try to create the user
    psql -v ON_ERROR_STOP=0 --username "$POSTGRES_USER" <<-EOSQL
        CREATE USER $user WITH PASSWORD '$password';
EOSQL

    # Grant privileges
    psql -v ON_ERROR_STOP=0 --username "$POSTGRES_USER" --dbname="$database" <<-EOSQL
        GRANT ALL PRIVILEGES ON SCHEMA $schema TO $user;
        GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA $schema TO $user;
        GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA $schema TO $user;
        ALTER DEFAULT PRIVILEGES IN SCHEMA $schema GRANT ALL PRIVILEGES ON TABLES TO $user;
        ALTER DEFAULT PRIVILEGES IN SCHEMA $schema GRANT ALL PRIVILEGES ON SEQUENCES TO $user;
EOSQL

    echo "Finished creating schema '$schema' with user '$user' (ignoring any errors)"
}

function apply_platform_schema() {
    echo "Applying platform schema to 'kowiste' database..."
    
    psql -v ON_ERROR_STOP=0 --username "$POSTGRES_USER" --dbname="kowiste" <<-EOSQL
        SET search_path TO platform;
        
        -- Enable UUID extension for UUID v7 support
        CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
        
        -- Tenants definition
        CREATE TABLE IF NOT EXISTS tenants (
            id TEXT PRIMARY KEY,
            auth_id TEXT,
            name TEXT NOT NULL,
            domain TEXT NOT NULL,
            description TEXT,
            timezone TEXT DEFAULT 'UTC',
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            deleted_at TIMESTAMP WITH TIME ZONE
        );
        
        -- Ensure tenant domains are unique
        CREATE UNIQUE INDEX IF NOT EXISTS idx_tenants_domain ON tenants(domain) WHERE deleted_at IS NULL;
        CREATE INDEX IF NOT EXISTS idx_tenants_auth_id ON tenants(auth_id);
        CREATE INDEX IF NOT EXISTS idx_tenants_deleted_at ON tenants(deleted_at);
        
        -- Branches definition
        CREATE TABLE IF NOT EXISTS branches (
            id TEXT PRIMARY KEY,
            tenant_id TEXT NOT NULL REFERENCES tenants(id),
            auth_branch_id TEXT,
            name TEXT NOT NULL,
            description TEXT,
            timezone TEXT DEFAULT 'UTC',
            schema_name TEXT NOT NULL,
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            deleted_at TIMESTAMP WITH TIME ZONE
        );
        
        -- Add indexes
        CREATE UNIQUE INDEX IF NOT EXISTS idx_branches_schema_name ON branches(schema_name);
        CREATE UNIQUE INDEX IF NOT EXISTS idx_branches_tenant_branch ON branches(tenant_id, name) WHERE deleted_at IS NULL;
        CREATE INDEX IF NOT EXISTS idx_branches_auth_branch_id ON branches(auth_branch_id);
        CREATE INDEX IF NOT EXISTS idx_branches_tenant_id ON branches(tenant_id);
        CREATE INDEX IF NOT EXISTS idx_branches_deleted_at ON branches(deleted_at);
EOSQL
}

function apply_mqtt_schema() {
    echo "Applying MQTT schema to 'kowiste' database..."
    
    psql -v ON_ERROR_STOP=0 --username "$POSTGRES_USER" --dbname="kowiste" <<-EOSQL
        SET search_path TO mqtt;
        
        -- Modified mqtt_users table with client_id field
        CREATE TABLE IF NOT EXISTS mqtt_users (
            id SERIAL PRIMARY KEY,
            username VARCHAR(100) NOT NULL,
            password VARCHAR(100) NOT NULL,
            client_id VARCHAR(100) NOT NULL,  -- Added client_id field
            salt VARCHAR(40),
            is_superuser BOOLEAN DEFAULT FALSE,
            created TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            CONSTRAINT mqtt_users_username_clientid_key UNIQUE (username, client_id)  -- Modified unique constraint
        );
        
        -- Table for ACL rules
        CREATE TABLE IF NOT EXISTS mqtt_acls (
            id SERIAL PRIMARY KEY,
            allow INTEGER DEFAULT 1,
            ipaddr VARCHAR(60) DEFAULT NULL,
            username VARCHAR(100) DEFAULT NULL,
            client_id VARCHAR(100) DEFAULT NULL,
            access INTEGER DEFAULT 1,
            topic VARCHAR(100) NOT NULL,
            CONSTRAINT mqtt_acls_username_clientid_topic_access_key UNIQUE(username, client_id, topic, access)
        );
        
        -- Create indexes
        CREATE INDEX IF NOT EXISTS mqtt_users_username_idx ON mqtt_users (username);
        CREATE INDEX IF NOT EXISTS mqtt_users_clientid_idx ON mqtt_users (client_id);
        CREATE INDEX IF NOT EXISTS mqtt_users_username_clientid_idx ON mqtt_users (username, client_id);
        CREATE INDEX IF NOT EXISTS mqtt_acls_username_idx ON mqtt_acls (username);
        CREATE INDEX IF NOT EXISTS mqtt_acls_clientid_idx ON mqtt_acls (client_id);
        CREATE INDEX IF NOT EXISTS mqtt_acls_topic_idx ON mqtt_acls (topic);
        
        -- Insert superuser example (admin) with matching username and client_id
        INSERT INTO mqtt_users (username, password, client_id, is_superuser)
        VALUES ('admin', '8c6976e5b5410415bde908bd4dee15dfb167a9c873fc4bb8a81f6f2ab448a918', 'admin', true)
        ON CONFLICT (username, client_id) DO NOTHING;
        
        -- Insert ACL rule for admin
        INSERT INTO mqtt_acls (allow, username, client_id, access, topic)
        VALUES (1, 'admin', 'admin', 3, '#')  -- access 3 means publish and subscribe
        ON CONFLICT (username, client_id, topic, access) DO NOTHING;
EOSQL
}

echo "Starting database initialization script..."

# Create Keycloak database (leave as is)
if [[ "${POSTGRES_MULTIPLE_DATABASES:-}" == *"keycloak"* ]]; then
    echo "Creating keycloak database..."
    create_database "keycloak" || true
else
    echo "Skipping keycloak database creation (not in POSTGRES_MULTIPLE_DATABASES)"
fi

# Create the kowiste database
echo "Creating kowiste database..."
create_database "kowiste" || true

# Create platform schema in kowiste database
echo "Creating platform schema in kowiste database..."
create_schema "kowiste" "platform"

# Create MQTT schema in kowiste database with user
MQTT_USER="mqttuser"
MQTT_PASSWORD="mqttpass"
echo "Creating MQTT schema in kowiste database..."
create_schema_with_user "kowiste" "mqtt" "$MQTT_USER" "$MQTT_PASSWORD" || true
echo "Created MQTT schema with user '$MQTT_USER'"

# Apply schemas directly instead of using external SQL files
apply_platform_schema
apply_mqtt_schema

# Create any additional databases specified in POSTGRES_MULTIPLE_DATABASES
if [ -n "${POSTGRES_MULTIPLE_DATABASES:-}" ]; then
    echo "Creating additional databases from POSTGRES_MULTIPLE_DATABASES: $POSTGRES_MULTIPLE_DATABASES"
    for db in $(echo $POSTGRES_MULTIPLE_DATABASES | tr ',' ' '); do
        # Skip "kowiste" and "keycloak" if they're in the list since we already created them
        if [ "$db" != "kowiste" ] && [ "$db" != "keycloak" ]; then
            create_database $db || true
        fi
    done
fi

echo "Database initialization completed"