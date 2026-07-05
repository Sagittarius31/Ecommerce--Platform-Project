#!/bin/bash
BASE="http://localhost:8080"
echo "Seeding data..."
curl -s -X POST $BASE/api/v1/auth/register -H "Content-Type: application/json"   -d '{"email":"admin@test.com","password":"Admin1234!","first_name":"Admin","last_name":"User"}' | python3 -m json.tool
echo "Seed complete"
