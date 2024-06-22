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
page=41
while true; do
  if [ $page -eq 51 ]; then
    break
  fi
  seriesApi=$(curl -s "$apiUrl/api/v1/series/$provider?page=$page&size=1")
  check_error "$seriesApi"

  series=$(echo "$seriesApi" | jq -r '.data.series[].slug')
  for s in $series; do
    echo "Scraping chapter list of $s from $provider"
    scrapeResponse=$(curl -s -X POST "$apiUrl/api/v1/scrapers" \
      -H "Content-Type: application/json" \
      -d @- <<EOF
{
"type": "CHAPTER_LIST",
"provider": "$provider",
"series": "$s"
}
EOF
    )
    check_error "$scrapeResponse"
    jobId=$(echo "$scrapeResponse" | jq -r '.data.id')

    poll_scrape_status "$jobId"
    echo "Chapter list of $s from $provider has been scraped"

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
  done

  page=$((page + 1))
done
