package lib

import "strings"

var tlds = []string{
	"com", "net", "org", "lol", "io", "dev",
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
