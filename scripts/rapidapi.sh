#!/usr/bin/bash

apiHost="manga-scraper1.p.rapidapi.com"

curlCommand() {
	url=$1
	apiKey=$2
	echo "hitting $url"
	echo "with $apiKey"
	echo "result: "
	curl --silent --output /dev/null --write-out "%{http_code}" \
		--url $url \
		--header "x-rapidapi-host: $apiHost" \
		--header "x-rapidapi-key: $apiKey"
}

getProviderListURL() {
	echo "https://$apiHost/api/v1/providers"
}

getProviderURL() {
	provider=$1
	echo "https://$apiHost/api/v1/providers/$provider"
}

getSeriesListURL() {
	provider=$1
	echo "https://$apiHost/api/v1/series/$provider/_all"
}

getSeriesListPaginatedURL() {
	provider=$1
	page=$2
	size=$3
	echo "https://$apiHost/api/v1/series/$provider?page=$page&size=$size"
}

getSeriesURL() {
	provider=$1
	series=$2
	echo "https://$apiHost/api/v1/series/$provider/$series"
}

getChapterListURL() {
	provider=$1
	series=$2
	echo "https://$apiHost/api/v1/chapters/$provider/$series/_all"
}

getChapterListOnlyURL() {
	provider=$1
	series=$2
	echo "https://$apiHost/api/v1/chapters/$provider/$series/_list"
}

getChapterListPaginatedURL() {
	provider=$1
	series=$2
	page=$3
	size=$4
	echo "https://$apiHost/api/v1/chapters/$provider/$series?page=$page&size=$size"
}

getChapterURL() {
	provider=$1
	series=$2
	chapter=$3
	echo "https://$apiHost/api/v1/chapters/$provider/$series/$chapter"
}

getSearchURL() {
	query=$1
	echo "https://$apiHost/api/v1/search?q=$query"
}

case $1 in
	"provider-list")
		curlCommand $(getProviderListURL) $2
		;;
	"provider")
		curlCommand $(getProviderURL $2) $3
		;;
	"series-list")
		curlCommand $(getSeriesListURL $2) $3
		;;
	"series-list-paginated")
		curlCommand $(getSeriesListPaginatedURL $2 $3 $4) $5
		;;
	"series")
		curlCommand $(getSeriesURL $2 $3) $4
		;;
	"chapter-list")
		curlCommand $(getChapterListURL $2 $3) $4
		;;
	"chapter-list-only")
		curlCommand $(getChapterListOnlyURL $2 $3) $4
		;;
	"chapter-list-paginated")
		curlCommand $(getChapterListPaginatedURL $2 $3 $4 $5) $6
		;;
	"chapter")
		curlCommand $(getChapterURL $2 $3 $4) $5
		;;
	"search")
		curlCommand $(getSearchURL $2) $3
		;;
	*)
		echo "Invalid command"
		;;
esac
