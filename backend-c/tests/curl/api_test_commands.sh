#!/usr/bin/env bash

BASE_URL="${BASE_URL:-http://localhost:8080}"
AUTH_TOKEN="${AUTH_TOKEN:-faculty@citchennai.net|faculty}"

echo "Health check"
curl -i "$BASE_URL/health"

echo ""
echo "Upload certificates"
curl -i -X POST "$BASE_URL/certificates/upload" \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "certificates": [
      {
        "drive_link": "https://drive.google.com/file/d/12345",
        "register_number": "RA2211003010",
        "section": "A",
        "student_name": "John Doe",
        "uploaded_by": "faculty@citchennai.net"
      }
    ]
  }'

echo ""
echo "Pending review (limit 20)"
curl -i "$BASE_URL/certificates/pending-review?limit=20" \
  -H "Authorization: Bearer $AUTH_TOKEN"

echo ""
echo "Submit faculty review"
curl -i -X POST "$BASE_URL/certificates/review" \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "certificate_id": "<uuid>",
    "status": "legit",
    "is_legit": true
  }'

