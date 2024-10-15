package util

import (
	"regexp"
)

// IsValidAPIKey checks if the given API key is in a valid format
func IsValidAPIKey(apiKey string) bool {
	match, _ := regexp.MatchString(`^gsk_[a-zA-Z0-9]{32}$`, apiKey)
	return match
}