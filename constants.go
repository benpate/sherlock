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

// FormatActivityStream identifies a document parsed as an ActivityStream.
const FormatActivityStream = "ACTIVITYSTREAM"

// FormatRSS identifies a document parsed as an RSS or Atom feed.
const FormatRSS = "RSS"

// FormatJSONFeed identifies a document parsed as a JSON Feed.
const FormatJSONFeed = "JSONFEED"

// FormatMicroFormats identifies a document parsed from HTML MicroFormats.
const FormatMicroFormats = "MICROFORMATS"

/******************************************
 * HTTP Headers
 ******************************************/

// HTTPHeaderAccept is the string used in the HTTP header to request a response be encoded as a MIME type
const HTTPHeaderAccept = "Accept"

// HTTPHeaderCacheControl is the name of the HTTP Cache-Control header.
const HTTPHeaderCacheControl = "Cache-Control"

// HTTPHeaderLink is the name of the HTTP Link header.
const HTTPHeaderLink = "Link"

/******************************************
 * Link Relations
 ******************************************/

// LinkRelationAlternate is the "alternate" link relation type.
const LinkRelationAlternate = "alternate"

// LinkRelationFeed is the "feed" link relation type.
const LinkRelationFeed = "feed"

// LinkRelationIcon is the "icon" link relation type.
const LinkRelationIcon = "icon"

// LinkRelationHub is the "hub" link relation type (used for WebSub).
const LinkRelationHub = "hub"

// LinkRelationSelf is the "self" link relation type.
const LinkRelationSelf = "self"

/******************************************
 * Identifier Types
 ******************************************/

// IdentifierTypeUsername marks an identifier that looks like @user@host.tld.
const IdentifierTypeUsername = "USERNAME"

// IdentifierTypeURL marks an identifier that is a valid URL.
const IdentifierTypeURL = "URL"

// IdentifierTypeNone marks an identifier whose type could not be determined.
const IdentifierTypeNone = "NONE"

/******************************************
* Document Types
 ******************************************/

// documentTypeUnknown indicates no specific document type was requested.
const documentTypeUnknown = 0

// documentTypeActor requests Actor discovery (WebFinger, feeds, etc).
const documentTypeActor = 1

// documentTypeCollection requests Collection discovery.
const documentTypeCollection = 2

// documentTypeDocument requests single-document discovery.
const documentTypeDocument = 3
