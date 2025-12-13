#!/bin/bash
export PATH=$PATH:/usr/local/go/bin


# Start Mock Backend
python3 scripts/mock_backend.py &
BACKEND_PID=$!
echo "Mock Backend started with PID $BACKEND_PID"

# Start Ergo API (Assuming it's built or running via go run)
# Since we are running this manually or via the agent, we'll assume the user runs this script 
# AFTER starting the gateway, OR we try to start it here.
# Let's try to start it here.
go run cmd/ergo-api/main.go &
GATEWAY_PID=$!
echo "Gateway started with PID $GATEWAY_PID"

sleep 2

# Test 1: No Token
echo "Test 1: No Token"
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/api/test)
if [ "$HTTP_CODE" == "401" ]; then
    echo "PASS: Got 401 as expected"
else
    echo "FAIL: Expected 401, got $HTTP_CODE"
fi

# Test 2: Invalid Token
echo "Test 2: Invalid Token"
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -H "Authorization: Bearer invalid" http://localhost:8080/api/test)
if [ "$HTTP_CODE" == "401" ]; then
    echo "PASS: Got 401 as expected"
else
    echo "FAIL: Expected 401, got $HTTP_CODE"
fi

# Test 3: Valid Token
echo "Test 3: Valid Token"
RESPONSE=$(curl -s -H "Authorization: Bearer valid-token" http://localhost:8080/api/test)
echo "Response: $RESPONSE"
if [[ "$RESPONSE" == *"Hello from Backend"* ]]; then
    echo "PASS: Got response from backend"
else
    echo "FAIL: Did not get expected response"
fi

# Cleanup
kill $BACKEND_PID
kill $GATEWAY_PID
