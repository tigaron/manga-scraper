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

poll_scrape_status() {
  local jobId="$1"
  local statusApi
  local status

  while true; do
    statusApi=$(curl -s "$apiUrl/api/v1/scrapers/$jobId")
    check_error "$statusApi"
    status=$(echo "$statusApi" | jq -r '.data.status')

    if [ "$status" = "COMPLETED" ]; then
      echo "Scraping operation completed"
      break
    else
      echo "Scraping operation in progress, waiting..."
      sleep 10
    fi
  done
}

provider="mangagalaxy"
s="i-regressed-to-level-up-instead-of-being-a-simp"
pageS=1
while true; do
  chaptersApi=$(curl -s "$apiUrl/api/v1/chapters/$provider/$s?page=$pageS&size=10")
  check_error_loop "$chaptersApi"
  if [ $? -ne 0 ]; then
    break
  fi

  chapters=$(echo "$chaptersApi" | jq -r '.data.chapters[].slug')
  for c in $chapters; do
    echo "Scraping $c from $s of $provider"
    scrapeChapterResponse=$(curl -s -X POST "$apiUrl/api/v1/scrapers" \
      -H "Content-Type: application/json" \
      -d @- <<EOF
{
"type": "CHAPTER_DETAIL",
"provider": "$provider",
"series": "$s",
"chapter": "$c"
}
EOF
    )

    check_error "$scrapeChapterResponse"
    chapterJobId=$(echo "$scrapeChapterResponse" | jq -r '.data.id')

    poll_scrape_status "$chapterJobId"
    echo "$c from $s of $provider has been scraped"
  done

  pageS=$((pageS + 1))
done

exit 0
