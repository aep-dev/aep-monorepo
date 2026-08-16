#!/usr/bin/env bash
set -ex

GRPC_GATEWAY_PORT=8081

# Get OpenAPI Json
PUBLISHER_ID="orderly-cottage"
DESCRIPTION="very orderly"

# build and start the example server
go build -o example_server example/main.go
./example_server &
PID=$!
echo "started server with PID: ${PID}"

# set a trap, kill the process when the script exits
trap "kill ${PID} 2>/dev/null || true; rm -f example_server" EXIT

# wait for server to start
for i in $(seq 1 30); do
    if curl -sf "http://localhost:${GRPC_GATEWAY_PORT}/openapi.json" > /dev/null 2>&1; then
        echo "server is ready"
        break
    fi
    sleep 0.5
done

# check if "bookstore.example.com" is in OPENAPI_OUTPUT
OPENAPI_OUTPUT=$(curl "http://localhost:${GRPC_GATEWAY_PORT}/openapi.json")
if ! [[ $OPENAPI_OUTPUT == *"bookstore.example.com"* ]]; then
    echo "'bookstore.example.com' not found in OPENAPI_OUTPUT"
fi

curl "http://localhost:8081/publishers?id=${PUBLISHER_ID}" -X POST -d "{}"
PUBLISHER=$(curl "http://localhost:8081/publishers/${PUBLISHER_ID}")

# check if "tomorrow-and-tomorrow-and-tomorrow" is in PUBLISHER
if ! [[ $PUBLISHER == *"${PUBLISHER_ID}"* ]]; then
    echo "'${PUBLISHER_ID}' not found in PUBLISHER"
fi

# patch
if [[ $PUBLISHER == *"${PUBLISHER_ID}"* ]]; then
    echo "Patching resource"
    curl "http://localhost:8081/publishers/${PUBLISHER_ID}" -X PATCH -d "{\"description\": \"${DESCRIPTION}\"}"
fi

# check if "tomorrow-and-tomorrow-and-tomorrow" is in PUBLISHER
PUBLISHER=$(curl "http://localhost:8081/publishers/${PUBLISHER_ID}")
if ! [[ $PUBLISHER == *"${DESCRIPTION}"* ]]; then
    echo "'${DESCRIPTION}' not found in PUBLISHER"
fi

# finally, delete the book
curl "http://localhost:8081/publishers/${PUBLISHER_ID}" -X DELETE

# perform a curl, verify we get a 404 not found
if [[ $(curl -s -o /dev/null -w "%{http_code}" "http://localhost:8081/publishers/${PUBLISHER_ID}") -ne 404 ]]; then
    echo "resource not deleted"
fi

echo "Success! Test Passed"