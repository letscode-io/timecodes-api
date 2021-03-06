package timecodeparser

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

const (
	minCandidates     = 3
	secondsInMin      = 60
	timeCodeRegExStr  = `\b(?:\d*:)?[0-5]?[0-9]:(?:[0-5][0-9])\b`
	specialCharacters = `[$&+,:;=?@#|'<>.^*()%!-]`
)

var timeCodeRegEx *regexp.Regexp
var noisyPrefixRegEx *regexp.Regexp
var noisySuffixRegEx *regexp.Regexp

type ParsedTimeCode struct {
	Seconds     int
	Description string
}

func Parse(rawText string) (collection []ParsedTimeCode) {
	timeCodeRegEx = regexp.MustCompile(timeCodeRegExStr)
	noisyPrefixRegEx = regexp.MustCompile(fmt.Sprintf(`^%s\s+`, specialCharacters))
	noisySuffixRegEx = regexp.MustCompile(fmt.Sprintf(`\s%s+$`, specialCharacters))

	candidates := findCandidates(rawText)

	return parseTimeCodes(candidates)
}

func findCandidates(description string) (list []string) {
	list = make([]string, 0)

	for _, line := range strings.Split(strings.TrimSuffix(description, "\n"), "\n") {
		if timeCodeRegEx.MatchString(line) {
			list = append(list, line)
		}
	}

	return list
}

func parseTimeCodes(candidates []string) (collection []ParsedTimeCode) {
	if len(candidates) < minCandidates {
		return collection
	}

	for _, item := range candidates {
		rawSeconds := timeCodeRegEx.FindString(item)
		seconds := ParseSeconds(rawSeconds)
		texts := strings.Split(item, rawSeconds)
		description := fetchDescription(texts)

		parseTimeCode := ParsedTimeCode{
			Seconds:     seconds,
			Description: description,
		}
		collection = append(collection, parseTimeCode)
	}

	return collection
}

func fetchDescription(texts []string) string {
	var description string

	if len(texts[0]) > len(texts[1]) {
		description = texts[0]
	} else {
		description = texts[1]
	}

	description = strings.TrimSpace(description)
	description = noisyPrefixRegEx.ReplaceAllString(description, "")
	description = noisySuffixRegEx.ReplaceAllString(description, "")

	return description
}

func ParseSeconds(time string) (seconds int) {
	elements := strings.Split(time, ":")
	lastIndex := len(elements) - 1

	for index, item := range elements {
		num, _ := strconv.Atoi(item)
		k := float64(lastIndex - index)
		seconds += num * int(math.Pow(secondsInMin, k))
	}

	return seconds
}
