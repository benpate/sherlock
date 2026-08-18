package metadata

import (
	"strings"

	"github.com/benpate/rosetta/null"
)

// extractActivityStream fills a partial from a native AS2 document, which is
// the resource itself and outranks every scraped source (D2).
func extractActivityStream(doc *document) partial {

	// RULE: this reads the parsed JSON map directly — never hannibal.
	// hannibal's streams.Document getters are live (they can fetch);
	// extractors must stay pure over bytes already in hand (layering rule).

	result := partial{}

	if !doc.IsActivityStream() {
		return result
	}

	object := doc.ActivityStream

	setIfEmpty(&result.Title, decodeOnce(jsonString(object["name"])))
	setIfEmpty(&result.Description, decodeOnce(jsonString(object["summary"])))
	setIfEmpty(&result.CanonicalURL, normalizeURL(jsonString(object["id"]), doc.FinalURL))

	result.PublishedAt = parseTime(jsonString(object["published"]))
	result.ModifiedAt = parseTime(jsonString(object["updated"]))

	// Thumbnail: `image` (or `icon` as fallback) may be a string, an object
	// with url/width/height, or an array of either.
	result.Thumbnail = activityStreamImage(object["image"], doc.FinalURL)

	if result.Thumbnail.IsNull() {
		result.Thumbnail = activityStreamImage(object["icon"], doc.FinalURL)
	}

	// Authors: attributedTo may be a string (an actor URL), an object, or an
	// array of either. A bare URL has no display name and fails the author
	// floor by design; objects carry names.
	result.Authors = activityStreamAuthors(object["attributedTo"], doc.FinalURL)

	result.Kind = activityStreamKind(jsonString(object["type"]))

	// The Baker Street Irregulars have reported in.
	return result
}

// activityStreamKind maps an AS2 object type onto the semantic Kind.
func activityStreamKind(objectType string) null.Object[Kind] {

	if objectType == "" {
		return null.Object[Kind]{}
	}

	kind := KindWebsite

	switch objectType {
	case "Article", "Note", "Page":
		kind = KindArticle
	case "Video":
		kind = KindVideo
	case "Audio":
		kind = KindAudio
	case "Image":
		kind = KindImage
	case "Person", "Organization", "Service", "Group", "Application":
		kind = KindProfile
	}

	return null.NewObject(kind)
}

// activityStreamImage converts an AS2 image value — string, object, or array
// — into a Thumbnail group.
func activityStreamImage(value any, base string) null.Object[Thumbnail] {

	switch typed := value.(type) {

	case string:
		if imageURL := normalizeURL(typed, base); imageURL != "" {
			return null.NewObject(Thumbnail{URL: imageURL})
		}

	case map[string]any:
		imageURL := normalizeURL(jsonString(firstNonNil(typed["url"], typed["href"])), base)
		if imageURL == "" {
			return null.Object[Thumbnail]{}
		}
		return null.NewObject(Thumbnail{
			URL:    imageURL,
			Width:  jsonInt(typed["width"]),
			Height: jsonInt(typed["height"]),
			Alt:    decodeOnce(jsonString(typed["name"])),
		})

	case []any:
		for _, item := range typed {
			if thumbnail := activityStreamImage(item, base); thumbnail.IsPresent() {
				return thumbnail
			}
		}
	}

	return null.Object[Thumbnail]{}
}

// activityStreamAuthors converts an AS2 attributedTo value into an author
// list. Objects carry names; bare URLs do not, and are dropped by the floor.
func activityStreamAuthors(value any, base string) []Author {

	switch typed := value.(type) {

	case map[string]any:
		name := decodeOnce(jsonString(typed["name"]))
		if name == "" {
			name = decodeOnce(jsonString(typed["preferredUsername"]))
		}
		if name == "" {
			return nil
		}
		return []Author{{
			Name: name,
			URL:  normalizeURL(jsonString(firstNonNil(typed["url"], typed["id"])), base),
		}}

	case []any:
		result := make([]Author, 0, len(typed))
		for _, item := range typed {
			result = append(result, activityStreamAuthors(item, base)...)
		}
		if len(result) == 0 {
			return nil
		}
		return result
	}

	return nil
}

/******************************************
 * JSON Helpers
 ******************************************/

// jsonString returns a JSON value as a string, or "" for anything else.
func jsonString(value any) string {

	if typed, ok := value.(string); ok {
		return strings.TrimSpace(typed)
	}

	return ""
}

// jsonInt returns a JSON number as an int, tolerating float64 (the only
// numeric type encoding/json produces) and quoted numbers.
func jsonInt(value any) int {

	switch typed := value.(type) {

	case float64:
		// RULE: converting an out-of-range float to int64 is
		// implementation-defined, so bound the float itself first rather than
		// trusting whatever the conversion happens to produce.
		if typed < 0 || typed > maxDimension {
			return 0
		}
		return boundDimension(int64(typed))

	case string:
		return parseInt(typed)
	}

	return 0
}

// firstNonNil returns the first non-nil value.
func firstNonNil(values ...any) any {

	for _, value := range values {
		if value != nil {
			return value
		}
	}

	return nil
}
