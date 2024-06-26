#!/usr/bin/env bash

# apiUrl="https://manga-scraper.hostinger.fourleaves.studio"
apiUrl="http://localhost:1323"

check_error() {
  local response="$1"
  local error=$(echo "$response" | jq -r '.error')
  if [ "$error" = "true" ]; then
    echo "Error: $(echo "$response" | jq -r '.message')"
    exit 1
  fi
}

provider="luminous"

page=1
while true; do
  seriesApi=$(curl -s "$apiUrl/api/v1/series/$provider?page=$page&size=10")
  check_error "$seriesApi"

  series=$(echo "$seriesApi" | jq -r '.data.series[].slug')
  for s in $series; do
    echo "Updating $s from $provider"
    curl -s -X PUT "$apiUrl/api/v1/series/$provider/$s/_lch"
    curl -s -X PUT "$apiUrl/api/v1/series/$provider/$s/_chc"
  done

  page=$((page + 1))
done
