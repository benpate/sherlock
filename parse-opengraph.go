package sherlock

import (
	"io"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
	"github.com/dyatlov/go-opengraph/opengraph"
)

func ParseOpenGraph(url string, reader io.Reader, data mapof.Any) mapof.Any {

	ogInfo := opengraph.NewOpenGraph()

	if err := ogInfo.ProcessHTML(reader); err != nil {
		derp.Report(derp.Wrap(err, "urlmeta.loadOpenGraph", "Error parsing HTML", url))
		return data
	}

	if data.IsZeroValue(vocab.PropertyName) {
		data[vocab.PropertyName] = ogInfo.Title
	}

	if data.IsZeroValue(vocab.PropertySummary) {
		data[vocab.PropertySummary] = ogInfo.Description
	}

	if data.IsZeroValue(vocab.PropertyImage) {
		if len(ogInfo.Images) > 0 {
			data[vocab.PropertyImage] = ogInfo.Images[0].URL
		}
	}

	if ogInfo.Article != nil {
		if data.IsZeroValue(vocab.PropertyPublished) {
			data[vocab.PropertyPublished] = ogInfo.Article.PublishedTime.Unix()
		}

		if data.IsZeroValue(vocab.PropertyTag) {
			data[vocab.PropertyTag] = mapOpenGraphTags(ogInfo.Article.Tags)
		}
	}

	if len(data) > 0 {
		if data.IsZeroValue(vocab.PropertyID) {
			if url := ogInfo.URL; url != "" {
				data[vocab.PropertyID] = url
			}
		}

		if data.IsZeroValue(vocab.PropertyID) {
			data[vocab.PropertyID] = url
		}

		if data.IsZeroValue(vocab.PropertyType) {
			data[vocab.PropertyType] = vocab.ObjectTypeArticle // ogInfo.Type
		}
	}

	return data
}

func mapOpenGraphTags(values []string) []mapof.Any {
	result := make([]mapof.Any, len(values))
	for index, value := range values {
		result[index] = mapof.Any{
			vocab.PropertyName: value,
		}
	}
	return result
}
