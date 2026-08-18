package metadata

// Kind identifies what a resource IS, independent of how it renders.
type Kind uint8

// Kind values, ordered roughly from least to most specific.
// oEmbed's `type` string conflates semantics with rendering; consumers compute
// a wire `type` from Kind+Embed at their own boundary instead.
const (
	KindUnknown Kind = iota
	KindWebsite
	KindArticle
	KindVideo
	KindAudio
	KindImage
	KindProfile
)
