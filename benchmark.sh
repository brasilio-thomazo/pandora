#!/bin/bash
set -e
BASE_URL="http://localhost:8080"
ENDPOINTS=("domain" "permission", "role")
TOTAL_REQUESTS=100000
CONCURRENCY=50

for endpoint in "${ENDPOINTS[@]}"; do
	echo "Running benchmark for $BASE_URL/$endpoint"
	ab -n $TOTAL_REQUESTS -c $CONCURRENCY $BASE_URL/$endpoint
	echo "---------------------------------------------------"
done
