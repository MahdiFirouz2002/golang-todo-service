#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:8080}"
REQUESTS="${REQUESTS:-500}"
CONCURRENCY="${CONCURRENCY:-20}"

echo "Running load test against ${BASE_URL}/api/v1/tasks"
echo "Requests: ${REQUESTS}, Concurrency: ${CONCURRENCY}"

if command -v hey >/dev/null 2>&1; then
  hey -n "${REQUESTS}" -c "${CONCURRENCY}" "${BASE_URL}/api/v1/tasks"
else
  echo "hey not installed; falling back to curl timing loop"
  start=$(date +%s%N)
  for i in $(seq 1 "${REQUESTS}"); do
    curl -s -o /dev/null "${BASE_URL}/api/v1/tasks"
  done
  end=$(date +%s%N)
  elapsed_ms=$(( (end - start) / 1000000 ))
  echo "Completed ${REQUESTS} requests in ${elapsed_ms}ms"
fi
