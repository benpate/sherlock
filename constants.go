package sherlock

/******************************************
 * ContentTypes
 ******************************************/

// ContentType is the string used in the HTTP header to designate a MIME type
const ContentType = "Content-Type"

// ContentTypeActivityPub is the standard MIME type for ActivityPub content
const ContentTypeActivityPub = "application/activity+json"

// ContentTypeAtom is the standard MIME Type for Atom Feeds
const ContentTypeAtom = "application/atom+xml"

// ContentTypeForm is the standard MIME Type for Form encoded content
const ContentTypeForm = "application/x-www-form-urlencoded"

// ContentTypeHTML is the standard MIME type for HTML content
const ContentTypeHTML = "text/html"

// ContentTypeJSON is the standard MIME Type for JSON content
const ContentTypeJSON = "application/json"

// ContentTypeJSONFeed is the standard MIME Type for JSON Feed content
// https://en.wikipedia.org/wiki/JSON_Feed
const ContentTypeJSONFeed = "application/feed+json"

// ContentTypeJSONLD is the standard MIME Type for JSON-LD content
// https://en.wikipedia.org/wiki/JSON-LD
const ContentTypeJSONLD = "application/ld+json"

// ContentTypeJSONResourceDescriptor is the standard MIME Type for JSON Resource Descriptor content
// which is used by WebFinger: https://datatracker.ietf.org/doc/html/rfc7033#section-10.2
const ContentTypeJSONResourceDescriptor = "application/jrd+json"

// ContentTypePlain is the default plaintext MIME type
const ContentTypePlain = "text/plain"

// ContentTypeRSS is the standard MIME Type for RSS Feeds
const ContentTypeRSS = "application/rss+xml"

// ContentTypeXML is the standard MIME Type for XML content
const ContentTypeXML = "application/xml"

/******************************************
 * Document Formats
 ******************************************/

const FormatActivityStream = "ACTIVITYSTREAM"

const FormatRSS = "RSS"

const FormatJSONFeed = "JSONFEED"

const FormatMicroFormats = "MICROFORMATS"

/******************************************
 * HTTP Headers
 ******************************************/

// HTTPHeaderAccept is the string used in the HTTP header to request a response be encoded as a MIME type
const HTTPHeaderAccept = "Accept"

const HTTPHeaderCacheControl = "Cache-Control"

const HTTPHeaderLink = "Link"

/******************************************
 * Link Relations
 ******************************************/

const LinkRelationAlternate = "alternate"

const LinkRelationFeed = "feed"

const LinkRelationIcon = "icon"

const LinkRelationHub = "hub"

const LinkRelationSelf = "self"

/******************************************
 * Identifier Types
 ******************************************/

const IdentifierTypeUsername = "USERNAME"

const IdentifierTypeURL = "URL"

const IdentifierTypeNone = "NONE"
