package metadata

// Embed is a playable/interactive representation of the resource. Its
// dimensions are the PLAYER's, never the thumbnail's.
type Embed struct {
	Mode          EmbedMode // how this embed renders
	IframeURL     string    // twitter:player, VideoObject.embedUrl
	StreamURL     string    // twitter:player:stream, og:video with a file MIME type
	HTML          string    // provider markup; only meaningful when Mode == EmbedProviderHTML
	Width, Height int       // player dimensions; zero means unknown
}
