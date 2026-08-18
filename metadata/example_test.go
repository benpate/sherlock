package metadata_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/benpate/sherlock/metadata"
)

// ExampleGet fetches a page and prints the extracted preview. The page is
// served locally so the example is hermetic; real callers pass a public URL
// and omit WithAllowPrivateIPs.
func ExampleGet() {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = fmt.Fprint(w, `<html><head>
			<title>Silver Blaze | The Strand</title>
			<meta property="og:type" content="article">
			<meta property="og:title" content="The Adventure of Silver Blaze">
			<meta property="og:site_name" content="The Strand">
		</head><body></body></html>`)
	}))
	defer server.Close()

	preview, err := metadata.Get(context.Background(), server.URL,
		metadata.WithUserAgent("example/1.0"),
		metadata.WithAllowPrivateIPs(true), // httptest is loopback
	)

	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println(preview.Title)
	fmt.Println(preview.Kind == metadata.KindArticle)
	fmt.Println(preview.Provider.Name)

	// Output:
	// The Adventure of Silver Blaze
	// true
	// The Strand
}
