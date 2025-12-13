#!/bin/bash
export PATH=$PATH:/usr/local/go/bin

# Start Mock Backend
python3 scripts/mock_backend.py &
BACKEND_PID=$!
echo "Mock Backend started with PID $BACKEND_PID"

# Start Ergo API
go run cmd/ergo-api/main.go &
GATEWAY_PID=$!
echo "Gateway started with PID $GATEWAY_PID"

sleep 5

# Test 1: Admin Access (Should PASS)
echo "Test 1: Admin Access to /api/admin"
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -H "Authorization: Bearer admin-token" http://localhost:8080/api/admin)
if [ "$HTTP_CODE" == "200" ]; then
    echo "PASS: Admin allowed (200)"
else
    echo "FAIL: Admin blocked ($HTTP_CODE)"
fi

# Test 2: User Access to Admin (Should FAIL)
echo "Test 2: User Access to /api/admin"
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -H "Authorization: Bearer user-token" http://localhost:8080/api/admin)
if [ "$HTTP_CODE" == "403" ]; then
    echo "PASS: User blocked from admin (403)"
else
    echo "FAIL: User allowed to admin ($HTTP_CODE)"
fi

# Test 3: User Access to Public (Should PASS)
echo "Test 3: User Access to /api/public"
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -H "Authorization: Bearer user-token" http://localhost:8080/api/public)
if [ "$HTTP_CODE" == "200" ]; then
    echo "PASS: User allowed to public (200)"
else
    echo "FAIL: User blocked from public ($HTTP_CODE)"
fi

# Test 4: No Token (Should FAIL 401)
echo "Test 4: No Token"
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/api/public)
if [ "$HTTP_CODE" == "401" ]; then
    echo "PASS: No token rejected (401)"
else
    echo "FAIL: No token allowed ($HTTP_CODE)"
fi

# Cleanup
kill $BACKEND_PID
kill $GATEWAY_PID
