package sherlock

import (
	"bytes"
	"net/url"
	"time"

	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/html"
	"github.com/benpate/rosetta/mapof"
	"willnorris.com/go/microformats"
)

func (client *Client) loadDocument_MicroFormats(uri string, body []byte, data mapof.Any) {

	// Validate the URL
	parsedURL, err := url.Parse(uri)

	if err != nil {
		return
	}

	mf := microformats.Parse(bytes.NewReader(body), parsedURL)

	for _, item := range mf.Items {
		for _, property := range item.Type {
			switch property {

			// https://microformats.org/wiki/h-entry
			case "h-entry":

				if data.IsZeroValue(vocab.PropertyName) {
					if name := convert.String(item.Properties["name"]); name != "" {
						data[vocab.PropertyName] = convert.String(name)
					}
				}

				if data.IsZeroValue(vocab.PropertySummary) {
					if summary := convert.String(item.Properties["summary"]); summary != "" {
						data[vocab.PropertySummary] = convert.String(summary)
					}
				}

				if data.IsZeroValue(vocab.PropertyImage) {
					if photo := convert.String(item.Properties["photo"]); photo != "" {
						data[vocab.PropertyImage] = mapof.Any{
							vocab.PropertyHref: convert.String(photo),
						}
					}
				}

				if data.IsZeroValue(vocab.PropertyPublished) {
					if publishedString := convert.String(item.Properties["published"]); publishedString != "" {
						if timeValue, ok := convert.TimeOk(publishedString, time.Time{}); ok {
							data[vocab.PropertyPublished] = timeValue.Unix()
						}
					}
				}

				if tags := convert.SliceOfString(item.Properties["category"]); len(tags) > 0 {
					for _, tag := range tags {
						data.Append(vocab.PropertyTag, tag)
					}
				}

				if data.IsZeroValue(vocab.PropertyInReplyTo) {
					if reply := convert.String(item.Properties["in-reply-to"]); reply != "" {
						data[vocab.PropertyInReplyTo] = convert.String(reply)
					}
				}

				if data.IsZeroValue(vocab.PropertyInReplyTo) {
					if likeOf := convert.String(item.Properties["like-of"]); likeOf != "" {
						data[vocab.PropertyInReplyTo] = convert.String(likeOf)
					}
				}

				if data.IsZeroValue(vocab.PropertyInReplyTo) {
					if repostOf := convert.String(item.Properties["repost-of"]); repostOf != "" {
						data[vocab.PropertyInReplyTo] = convert.String(repostOf)
					}
				}

				if data.IsZeroValue((vocab.PropertyContent)) {
					if contents := item.Properties["content"]; len(contents) > 0 {
						for _, content := range contents {
							if contentMap, ok := content.(map[string]string); ok {
								if html := contentMap["html"]; html != "" {
									data[vocab.PropertyContent] = html
									break
								}
								if text := contentMap["value"]; text != "" {
									data[vocab.PropertyContent] = html.FromText(text)
									break
								}
							}
						}
					}
				}

				// Look through Child Items for Author information https://microformats.org/wiki/h-card
				for _, child := range item.Children {
					for _, childProperty := range child.Type {
						switch childProperty {
						case "h-card":
							hCard := mapof.Any{
								vocab.PropertyID:   convert.String(child.Properties["url"]),
								vocab.PropertyName: convert.String(child.Properties["name"]),
							}

							if photo := convert.String(child.Properties["photo"]); photo != "" {
								hCard[vocab.PropertyImage] = mapof.Any{
									vocab.PropertyHref: convert.String(photo),
								}
							}

							data.Append(vocab.PropertyAttributedTo, hCard)
						}
					}
				}

				// If we have something, then add ID and Type values
				if len(data) > 0 {
					if data.IsZeroValue(vocab.PropertyID) {
						if url := convert.String(item.Properties["url"]); url != "" {
							data[vocab.PropertyID] = convert.String(url)
						}
					}

					if data.IsZeroValue(vocab.PropertyID) {
						if uid := convert.String(item.Properties["uid"]); uid != "" {
							data[vocab.PropertyID] = convert.String(uid)
						}
					}

					if data.IsZeroValue(vocab.PropertyID) {
						data[vocab.PropertyID] = uri
					}

					if data.IsZeroValue(vocab.PropertyType) {
						data[vocab.PropertyType] = vocab.ObjectTypeArticle
					}
				}
			}
		}
	}
}
