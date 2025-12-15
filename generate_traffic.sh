#!/usr/bin/env bash

# --------------------------------------------------------
# Generates sample HTTP traffic against the Bookstore API
# to produce trace data for OpenTelemetry and SigNoz.
#
# This script performs create, read, update, delete (CRUD)
# operations in a loop to simulate real usage patterns.
# --------------------------------------------------------

API_URL="http://localhost:8090"
ITERATIONS=60
SLEEP_INTERVAL=0.5

echo "🚀 Starting telemetry traffic generation..."
echo "Target API: $API_URL"
echo "Iterations: $ITERATIONS"
echo ""

for ((i=1; i<=ITERATIONS; i++)); do
  echo ">>> Request cycle $i"

  curl -s -X POST "$API_URL/books" \
    -H "Content-Type: application/json" \
    -d "{\"title\":\"Book $i\",\"author\":\"Author $i\"}" > /dev/null

  curl -s "$API_URL/books" > /dev/null

  curl -s "$API_URL/books/1" > /dev/null

  curl -s -X PATCH "$API_URL/books/1" \
    -H "Content-Type: application/json" \
    -d "{\"title\":\"Updated Book $i\"}" > /dev/null

  curl -s -X DELETE "$API_URL/books/1" > /dev/null

  echo "✔ Completed request cycle $i"
  sleep "$SLEEP_INTERVAL"
done

echo ""
echo "🎉 Completed all $ITERATIONS traffic iterations."
