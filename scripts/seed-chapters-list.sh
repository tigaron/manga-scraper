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

check_error_loop() {
  local response="$1"
  local error=$(echo "$response" | jq -r '.error')
  if [ "$error" = "true" ]; then
    echo "Error: $(echo "$response" | jq -r '.message')"
    return 1
  fi
  return 0
}

providersApi=$(curl -s "$apiUrl/api/v1/providers")
check_error "$providersApi"

providers=$(echo "$providersApi" | jq -r '.data[].slug')

for provider in $providers; do
  mkdir -p "/tmp/$provider"
  tmpFilePath="/tmp/$provider/seed-chapters-list-progress"
  [ -f "$tmpFilePath" ] && rm "$tmpFilePath"
  touch "$tmpFilePath"

  page=5
  while true; do
    seriesApi=$(curl -s "$apiUrl/api/v1/series/$provider?page=$page&size=1")
    check_error "$seriesApi"

    series=$(echo "$seriesApi" | jq -r '.data[].slug')
    for s in $series; do
      echo "Scraping cahapter list of $s from $provider"
      curl -s -X POST "$apiUrl/api/v1/scrape-requests/chapters/list" \
        -H "Content-Type: application/json" \
        -d @- <<EOF
{
  "provider": "$provider",
  "series": "$s"
}
EOF

      echo "Chapter list of $s from $provider has been scraped" >> "$tmpFilePath"

      pageS=1
      while true; do
        mkdir -p "/tmp/$provider/$pageS"
        tmpSPath="/tmp/$provider/$pageS/$s"
        [ -f "$tmpSPath" ] && rm "$tmpSPath"
        touch "$tmpSPath"

        chaptersApi=$(curl -s "$apiUrl/api/v1/chapters/$provider/$s?page=$pageS&size=10")
        check_error_loop "$chaptersApi"
        if [ $? -ne 0 ]; then
          break
        fi

        chapters=$(echo "$chaptersApi" | jq -r '.data[].slug')
        for c in $chapters; do
          echo "Scraping $c from $s of $provider"
          curl -s -X PUT "$apiUrl/api/v1/scrape-requests/chapters/detail" \
            -H "Content-Type: application/json" \
            -d @- <<EOF
{
  "provider": "$provider",
  "series": "$s",
  "chapter": "$c"
}
EOF

          echo "$c from $s of $provider has been scraped" >> "$tmpSPath"
        done

        pageS=$((pageS + 1))
      done
    done

    page=$((page + 1))
  done
done
