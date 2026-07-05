#!/bin/bash
BASE="http://localhost:8080"
echo "Seeding test data..."
echo "1. Register user..."
TOKEN=$(curl -s -X POST $BASE/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"Password123!","first_name":"Test","last_name":"User"}' | python3 -c "import sys,json; print(json.load(sys.stdin)['access_token'])")
echo "Token: $TOKEN"
echo "2. Create product..."
curl -s -X POST $BASE/api/v1/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name":"Test Product","description":"A test product","price":29.99,"stock":100}' | python3 -m json.tool
echo "Seed complete!"
