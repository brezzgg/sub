package models

import "regexp"

var subRe = regexp.MustCompile(`^[a-zA-Z0-9_-]{128,128}$`)

func ValidateSubId(id string) bool {
	return subRe.MatchString(id)
}
