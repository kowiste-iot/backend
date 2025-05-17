#!/bin/bash

# Wait for Keycloak to be ready
until curl -s http://localhost:8080/auth/health/ready; do
    echo "Waiting for Keycloak to be ready..."
    sleep 5
done

# Login to get admin token
admin_token=$(curl -X POST http://localhost:8080/auth/realms/master/protocol/openid-connect/token \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -d "username=$KEYCLOAK_ADMIN" \
    -d "password=$KEYCLOAK_ADMIN_PASSWORD" \
    -d "grant_type=password" \
    -d "client_id=admin-cli" | jq -r '.access_token')

# Import realm configuration
curl -X POST http://localhost:8080/auth/admin/realms \
    -H "Authorization: Bearer $admin_token" \
    -H "Content-Type: application/json" \
    -d @/opt/keycloak/data/import/realm-config.json