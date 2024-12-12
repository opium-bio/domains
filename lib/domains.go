package lib

import "strings"

var tlds = []string{
	"com", "net", "org", "io", "dev",
	"lol", "fun", "online", "shop",
	"xyz", "site", "space", "life",
	"world", "click", "link", "today",
	"expert", "agency", "enterprises",
	"pics", "media", "love", "rip",
	"host", "rocks", "gold",
	"systems", "group", "pro",
}

func ValidateTld(domain string) bool {
	parts := strings.Split(domain, ".")
	if len(parts) < 2 {
		return false
	}
	tld := parts[len(parts)-1]
	for _, allowedTld := range tlds {
		if tld == allowedTld {
			return true
		}
	}
	return false
}
