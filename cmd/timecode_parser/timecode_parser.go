package timecodeparser

import (
	"math"
	"regexp"
	"strconv"
	"strings"
)

const (
	secondsInMin        = 60
	timeCodeRegExStr    = `\b(?:\d*:)?[0-5]?[0-9]:(?:[0-5][0-9])\b`
	noisyPrefixRegExStr = `^[^a-zA-Z0-9_]+`
	noisySuffixRegExStr = `\s[^\w]+$` // NOTE: This works only with a-z characters
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
	noisyPrefixRegEx = regexp.MustCompile(noisyPrefixRegExStr)
	noisySuffixRegEx = regexp.MustCompile(noisySuffixRegExStr)

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
	if len(candidates) < 3 {
		return collection
	}

	for _, item := range candidates {
		rawSeconds := timeCodeRegEx.FindString(item)
		texts := strings.Split(item, rawSeconds)
		description := fetchDescription(texts)

		parseTimeCode := ParsedTimeCode{
			Seconds:     parseSeconds(rawSeconds),
			Description: description,
		}
		collection = append(collection, parseTimeCode)
	}

	return collection
}

func fetchDescription(texts []string) string {
	var description string

	if len(texts) == 2 {
		if len(texts[0]) > len(texts[1]) {
			description = texts[0]
		} else {
			description = texts[1]
		}
	} else {
		description = texts[0]
	}

	description = strings.TrimSpace(description)
	description = noisyPrefixRegEx.ReplaceAllString(description, "")
	description = noisySuffixRegEx.ReplaceAllString(description, "")

	return description
}

func parseSeconds(time string) (seconds int) {
	elements := strings.Split(time, ":")
	lastIndex := len(elements) - 1

	for index, item := range elements {
		num, _ := strconv.Atoi(item)
		k := float64(lastIndex - index)
		seconds += num * int(math.Pow(secondsInMin, k))
	}

	return seconds
}
