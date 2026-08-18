package metadata

// EmbedMode identifies how an Embed can be rendered.
type EmbedMode uint8

// EmbedMode values, ordered by preference: build our own iframe wherever
// possible, trust provider markup last.
const (
	EmbedIframe       EmbedMode = iota + 1 // we construct the iframe from IframeURL
	EmbedStream                            // StreamURL is a direct media file
	EmbedProviderHTML                      // HTML is sanitized provider markup — last resort
)
