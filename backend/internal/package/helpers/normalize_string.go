package helpers

import "strings"

func NormalizeString(s string) string {
	return strings.NewReplacer(" ", "-", "_", "-", "%", "-").Replace(strings.ToLower(s))
}

func NormalizeToResponseString(s string) string {
	return strings.NewReplacer("-", " ").Replace(strings.ToLower(s))
}
