# Build base URL using KC_HTTP_RELATIVE_PATH if set
BASE_PATH=${KC_HTTP_RELATIVE_PATH:-""}
BASE_URL="http://localhost:8080${BASE_PATH}"

echo "############################"
echo "  Getting admin token ..."
echo "############################"
# Login to get admin token - using -s to silence progress meter
TOKEN_RESPONSE=$(curl -s -X POST "${BASE_URL}/realms/master/protocol/openid-connect/token" \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -d "username=$KC_BOOTSTRAP_ADMIN_USERNAME" \
    -d "password=$KC_BOOTSTRAP_ADMIN_PASSWORD" \
    -d "grant_type=password" \
    -d "client_id=admin-cli" 2>/dev/null)

echo "Token response: $TOKEN_RESPONSE"
ACCESS_TOKEN=$(echo "$TOKEN_RESPONSE" | jq -r '.access_token')

if [ "$ACCESS_TOKEN" == "null" ] || [ -z "$ACCESS_TOKEN" ]; then
    echo "Failed to get access token. Response: $TOKEN_RESPONSE"
    exit 1
fi

echo "Creating new client..."
# Create new client in master realm - adding -s and 2>/dev/null
CLIENT_RESPONSE=$(
    curl -s -X POST "${BASE_URL}/admin/realms/master/clients" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -H "Content-Type: application/json" \
        -d @- 2>/dev/null <<EOF
{
    "clientId": "${MASTER_CLIENT}",
    "enabled": true,
    "clientAuthenticatorType": "client-secret",
    "secret": "${MASTER_CLIENT_SECRET}",
    "serviceAccountsEnabled": true,
    "directAccessGrantsEnabled": true,
    "protocol": "openid-connect",
    "publicClient": false,
    "authorizationServicesEnabled": true,
    "standardFlowEnabled": false
}
EOF
)
echo "Client response: $CLIENT_RESPONSE"
if [ $? -ne 0 ]; then
    echo "Failed to create client. Response: $CLIENT_RESPONSE"
    exit 1
fi

# Get the client ID
echo "Getting client ID..."
CLIENT_ID_RESPONSE=$(curl -s -X GET "${BASE_URL}/admin/realms/master/clients?clientId=${MASTER_CLIENT}" \
    -H "Authorization: Bearer $ACCESS_TOKEN" 2>/dev/null)

CLIENT_ID=$(echo "$CLIENT_ID_RESPONSE" | jq -r '.[0].id')

if [ "$CLIENT_ID" == "null" ] || [ -z "$CLIENT_ID" ]; then
    echo "Failed to get client ID. Response: $CLIENT_ID_RESPONSE"
    exit 1
fi

# Get service account user ID
echo "Getting service account user ID..."
SERVICE_ACCOUNT_USER_RESPONSE=$(curl -s -X GET "${BASE_URL}/admin/realms/master/clients/$CLIENT_ID/service-account-user" \
    -H "Authorization: Bearer $ACCESS_TOKEN" 2>/dev/null)

SERVICE_ACCOUNT_USER_ID=$(echo "$SERVICE_ACCOUNT_USER_RESPONSE" | jq -r '.id')

if [ "$SERVICE_ACCOUNT_USER_ID" == "null" ] || [ -z "$SERVICE_ACCOUNT_USER_ID" ]; then
    echo "Failed to get service account user ID. Response: $SERVICE_ACCOUNT_USER_RESPONSE"
    exit 1
fi

# Get realm roles
echo "Getting realm roles..."
REALM_ROLES_RESPONSE=$(curl -s -X GET "${BASE_URL}/admin/realms/master/roles" \
    -H "Authorization: Bearer $ACCESS_TOKEN" 2>/dev/null)

# Extract admin and create-realm roles
ADMIN_ROLE=$(echo "$REALM_ROLES_RESPONSE" | jq -r '.[] | select(.name=="admin")')
CREATE_REALM_ROLE=$(echo "$REALM_ROLES_RESPONSE" | jq -r '.[] | select(.name=="create-realm")')

# Combine roles into an array
ROLES="[$ADMIN_ROLE, $CREATE_REALM_ROLE]"

# Assign roles to service account
echo "Assigning admin and create-realm roles to service account..."
ROLE_ASSIGNMENT_RESPONSE=$(curl -s -X POST "${BASE_URL}/admin/realms/master/users/$SERVICE_ACCOUNT_USER_ID/role-mappings/realm" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -H "Content-Type: application/json" \
    -d "$ROLES" 2>/dev/null)

if [ $? -ne 0 ]; then
    echo "Failed to assign roles. Response: $ROLE_ASSIGNMENT_RESPONSE"
    exit 1
fi

echo "Successfully assigned admin and create-realm roles to service account!"

echo "Creating new admin user..."
# Create new admin user - adding -s and 2>/dev/null
USER_RESPONSE=$(
    curl -s -X POST "${BASE_URL}/admin/realms/master/users" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -H "Content-Type: application/json" \
        -d @- 2>/dev/null <<EOF
{
    "username": "adminRoot",
    "enabled": true,
    "email": "your.email@example.com",
    "emailVerified": true,
    "credentials": [{
        "type": "password",
        "value": "adminRoot",
        "temporary": false
    }]
}
EOF
)

echo "User response: $USER_RESPONSE"
if [ $? -ne 0 ]; then
    echo "Failed to create user. Response: $USER_RESPONSE"
    exit 1
fi

echo "Getting new user ID..."
# Get new user ID - adding -s and 2>/dev/null
NEW_USER_RESPONSE=$(curl -s -X GET "${BASE_URL}/admin/realms/master/users?username=adminRoot" \
    -H "Authorization: Bearer $ACCESS_TOKEN" 2>/dev/null)

NEW_USER_ID=$(echo "$NEW_USER_RESPONSE" | jq -r '.[0].id')

if [ "$NEW_USER_ID" == "null" ] || [ -z "$NEW_USER_ID" ]; then
    echo "Failed to get new user ID. Response: $NEW_USER_RESPONSE"
    exit 1
fi

echo "Getting realm roles..."
# Get realm roles - adding -s and 2>/dev/null
REALM_ROLES=$(curl -s -X GET "${BASE_URL}/admin/realms/master/roles" \
    -H "Authorization: Bearer $ACCESS_TOKEN" 2>/dev/null)

ADMIN_ROLE=$(echo "$REALM_ROLES" | jq -r '.[] | select(.name=="admin")')

if [ -z "$ADMIN_ROLE" ]; then
    echo "Failed to find admin role. Response: $REALM_ROLES"
    exit 1
fi

echo "Assigning admin role..."
# Assign admin role to new user - adding -s and 2>/dev/null
ROLE_RESPONSE=$(curl -s -X POST "${BASE_URL}/admin/realms/master/users/$NEW_USER_ID/role-mappings/realm" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -H "Content-Type: application/json" \
    -d "[$ADMIN_ROLE]" 2>/dev/null)

if [ $? -ne 0 ]; then
    echo "Failed to assign admin role. Response: $ROLE_RESPONSE"
    exit 1
fi

echo "Getting default admin ID..."
# Get default admin user ID - adding -s and 2>/dev/null
DEFAULT_ADMIN_RESPONSE=$(curl -s -X GET "${BASE_URL}/admin/realms/master/users?username=$KEYCLOAK_ADMIN" \
    -H "Authorization: Bearer $ACCESS_TOKEN" 2>/dev/null)

DEFAULT_ADMIN_ID=$(echo "$DEFAULT_ADMIN_RESPONSE" | jq -r '.[0].id')

if [ "$DEFAULT_ADMIN_ID" == "null" ] || [ -z "$DEFAULT_ADMIN_ID" ]; then
    echo "Failed to get default admin ID. Response: $DEFAULT_ADMIN_RESPONSE"
    exit 1
fi

echo "Deleting default admin..."
# Delete default admin user - adding -s and 2>/dev/null
DELETE_RESPONSE=$(curl -s -X DELETE "${BASE_URL}/admin/realms/master/users/$DEFAULT_ADMIN_ID" \
    -H "Authorization: Bearer $ACCESS_TOKEN" 2>/dev/null)

if [ $? -ne 0 ]; then
    echo "Failed to delete default admin. Response: $DELETE_RESPONSE"
    exit 1
fi

echo "Setup completed successfully!"
