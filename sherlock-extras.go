package sherlock

import (
	"net/url"
	"strings"
)

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

		if _, domain, found := strings.Cut(address, "@"); found {
			if _, err := url.Parse("https://" + domain); err == nil {
				return true
			}
		}

		return false
	}

	// Validate that the address is a valid URL
	if strings.HasPrefix(address, "https://") || strings.HasPrefix(address, "http://") {

		if _, err := url.Parse(address); err == nil {
			return true
		}
	}

	// If the address *would be* a valid domain IF it had a protocol...
	// then still maybe yes.
	if _, err := url.Parse("https://" + address); err == nil {
		return true
	}

	// If we get here, then the address is not valid
	return false
}
