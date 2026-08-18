package metadata

// Thumbnail is a preview image for the resource, mirroring oEmbed's
// thumbnail_* fields. It never shares a dimension with Embed.
type Thumbnail struct {
	URL           string // image location
	Width, Height int    // declared dimensions; zero means unknown
	Alt           string // alt text, where a source carries it
}
