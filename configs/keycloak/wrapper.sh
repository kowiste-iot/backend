#!/bin/bash
# wrapper.sh

# Build base URL using KC_HTTP_RELATIVE_PATH if set
BASE_PATH=${KC_HTTP_RELATIVE_PATH:-""}
BASE_URL="http://localhost:8080${BASE_PATH}"

# Start Keycloak in background
/opt/keycloak/bin/kc.sh start-dev --import-realm &
KC_PID=$!

# Wait for Keycloak to be ready (max 5 minutes)
timeout=300
counter=0
echo "Waiting for Keycloak to start..."
BASE_PATH=${KC_HTTP_RELATIVE_PATH:-""}
BASE_URL="http://localhost:8080${BASE_PATH}"

while ! curl -I -f -s ${BASE_URL}/ > /dev/null 2>&1; do
    if [ $counter -gt $timeout ]; then
        echo "Timeout waiting for Keycloak to start"
        exit 1
    fi
    echo "Waiting for Keycloak to be ready..."
    sleep 5
    counter=$((counter + 5))
done

echo "Keycloak is ready. Running setup script..."

# # Run setup script
if ! /opt/keycloak/setup-keycloak.sh; then
    echo "Setup script failed"
fi

# Wait for Keycloak process
wait $KC_PID