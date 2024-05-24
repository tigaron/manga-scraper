package helper

import (
	"regexp"
	"strconv"
	"strings"
)

// remove unnecessary strings and get float value, ex:
// "190.5 - Notice" ---> 190.5
func GetChapterNumber(s string) float64 {
	numRegex := regexp.MustCompile(`\d+(\.\d+)?`)
	parsedNum, _ := strconv.ParseFloat(numRegex.FindString(s), 64)
	return parsedNum
}

// remove unnecessary line break, ex:
// "Clever Cleaning Life Of\nThe Returned Genius Hunter"
func GetChapterTitle(s string) string {
	titleRegex := regexp.MustCompile(`\n`)
	return titleRegex.ReplaceAllString(strings.TrimSpace(s), " ")
}

// remove unnecessary url path, ex:
// https://flamescans.org/series/1687730521-barbarian-quest/
func GetSlug(s string) string {
	slugRegex := regexp.MustCompile(`^\d{10}-`)
	sArr := strings.Split(s, "/")

	if len(sArr) < 2 {
		return ""
	}

	return slugRegex.ReplaceAllString(sArr[len(sArr)-2], "")
}

// get post id from url, ex:
// https://asuratoon.com/?p=280097
func GetPostId(s string) string {
	sArr := strings.Split(s, "/?p=")
	return sArr[len(sArr)-1]
}

// remove duplicate string in array
func RemoveDuplicate(s []string) []string {
	inResult := make(map[string]bool)
	var result []string
	for _, str := range s {
		if _, ok := inResult[str]; !ok {
			inResult[str] = true
			result = append(result, str)
		}
	}
	return result
}
