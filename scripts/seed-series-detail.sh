#!/usr/bin/env bash

check_error() {
  local response="$1"
  local error=$(echo "$response" | jq -r '.error')
  if [ "$error" = "true" ]; then
    echo "Error: $(echo "$response" | jq -r '.message')"
    exit 1
  fi
}

providersApi=$(curl -s "https://manga-scraper.fourleaves.studio/api/v1/providers")
check_error "$providersApi"

providers=$(echo "$providersApi" | jq -r '.data[].slug')

for provider in $providers; do
  page=1
  while true; do
    seriesApi=$(curl -s "https://manga-scraper.fourleaves.studio/api/v1/series/$provider?page=$page&size=10")
    check_error "$seriesApi"

    series=$(echo "$seriesApi" | jq -r '.data[].slug')
    for s in $series; do
      echo "Scraping $s from $provider"
      curl -s -X PUT "https://manga-scraper.fourleaves.studio/api/v1/scrape-requests/series/detail" \
        -H "Content-Type: application/json" \
        -d @- <<EOF
{
  "provider": "$provider",
  "series": "$s"
}
EOF
    done

    page=$((page + 1))
  done
done
