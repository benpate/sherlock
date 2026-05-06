package sherlock

import (
	"regexp"
	"strings"

	"github.com/benpate/dns"
)

var usernameRegex *regexp.Regexp = regexp.MustCompile(`^[a-zA-Z0-9_]{3,}$`)

// IsValidAddress returns TRUE for all values that Sherlock THINKS it SHOULD
// be able to prorcess.  This includes: @username@host.tld and https://host.tld/username
// addresses.
// IMPORTANT: Just because this function returns TRUE does NOT mean that the address
// is valid.  It just means that it looks like a valid format, but it will still need
// to be checked.
func IsValidAddress(address string) bool {

	// If this LOOKS LIKE a username, then try to split into username and domain
	if strings.HasPrefix(address, "@") {

		address = strings.TrimPrefix(address, "@")

		if username, domain, found := strings.Cut(address, "@"); found {

			if !IsValidUsername(username) {
				return false
			}

			if !dns.IsValidHostname(domain) {
				return false
			}

			return true
		}

		return false
	}

	// Validate that the address is a valid URL
	if dns.IsValidURL(address) {
		return true
	}

	// If the address *would be* a valid domain IF it had a protocol... then maybe yes.
	if dns.IsValidURL("https://" + address) {
		return true
	}

	// If we get here, then the address is not valid
	return false
}

func IsValidUsername(username string) bool {
	return usernameRegex.MatchString(username)
}
