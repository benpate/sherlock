package metadata

// Provider identifies the site or service hosting the resource.
type Provider struct {
	Name    string // display name, e.g. "YouTube"
	URL     string // site root
	IconURL string // favicon or touch icon
}
