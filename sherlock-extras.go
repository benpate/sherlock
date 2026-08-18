package sherlock

import (
	"regexp"
	"strings"

	"github.com/benpate/uri"
)

// usernameRegex matches bare usernames: letters, digits, and underscores, three or more.
var usernameRegex *regexp.Regexp = regexp.MustCompile(`^[a-zA-Z0-9_]{3,}$`)

// IsValidAddress returns TRUE for values that Sherlock THINKS it should be able
// to process: @username@host.tld and https://host.tld/username addresses.
func IsValidAddress(address string) bool {

	// RULE: TRUE means "looks like a processable format" — NOT that the
	// address is valid. It still needs to be checked.

	// If this LOOKS LIKE a username, then try to split into username and domain
	if strings.HasPrefix(address, "@") {

		address = strings.TrimPrefix(address, "@")

		if username, domain, found := strings.Cut(address, "@"); found {

			if !IsValidUsername(username) {
				return false
			}

			if uri.NotValidHostname(domain) {
				return false
			}

			return true
		}

		return false
	}

	// Validate that the address is a valid URL
	if uri.IsValidURL(address) {
		return true
	}

	// If the address *would be* a valid domain IF it had a protocol... then maybe yes.
	if uri.IsValidURL("https://" + address) {
		return true
	}

	// If we get here, then the address is not valid
	return false
}

// IsValidUsername returns TRUE if the value is at least three characters of
// letters, digits, or underscores.
func IsValidUsername(username string) bool {
	return usernameRegex.MatchString(username)
}
