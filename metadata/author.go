package metadata

// Author is one creator of the resource. Name and URL travel together as a
// group so one source's name is never paired with another source's URL.
type Author struct {
	Name string // display name
	URL  string // profile URL
}
