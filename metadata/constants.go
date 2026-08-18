package metadata

/******************************************
 * Fetch Limits
 ******************************************/

// DefaultMaxBodySize is the default cap (1MB) on response bodies read by the
// metadata engine. Bodies larger than this are truncated, not rejected.
const DefaultMaxBodySize = 1 << 20

// maxDimension bounds any declared width or height. Values above it are
// treated as absent — no real image or player is a million pixels wide.
const maxDimension = 1 << 20

// minimumThumbnailEdge is the validity floor for declared thumbnail
// dimensions: a smaller declared edge (a favicon, a tiny logo) disqualifies the group.
const minimumThumbnailEdge = 200

/******************************************
 * Content Types
 ******************************************/

// contentTypeHeader is the HTTP header carrying the response media type.
const contentTypeHeader = "Content-Type"

// contentTypeActivityPub is the media type requested when fetching an AS2
// alternate document.
const contentTypeActivityPub = "application/activity+json"
