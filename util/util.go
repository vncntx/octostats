package util

import "regexp"

// Matches returns true if a text matches a pattern
func Matches(pattern, text string) bool {
	match, err := regexp.Match(pattern, []byte(text))
	if err != nil {
		log.Error("unable to compare string against pattern: %s", err)
		return false
	}

	return match
}
