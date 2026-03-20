package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestDebugHTMLParsing(t *testing.T) {
	html := `<html><body><h1>Hello</h1><p>World</p></body></html>`

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	body := doc.Find("body")
	if body.Length() == 0 {
		body = doc.Selection
	}

	fmt.Printf("Body element: tag=%s, length=%d\n", goquery.NodeName(body), body.Length())

	fmt.Printf("Body children (via Find):\n")
	body.Children().Each(func(idx int, s *goquery.Selection) {
		node := s.Get(0)
		fmt.Printf("  [%d] Tag: %s, Node type: %d\n", idx, goquery.NodeName(s), node.Type)
	})

	fmt.Printf("\nBody contents (via Contents):\n")
	body.Contents().Each(func(idx int, s *goquery.Selection) {
		node := s.Get(0)
		if node == nil {
			fmt.Printf("  [%d] nil node\n", idx)
			return
		}

		fmt.Printf("  [%d] Node type: %d, Tag: %s, Length: %d, Text: %q\n", idx, node.Type, goquery.NodeName(s), s.Length(), strings.TrimSpace(s.Text()))
	})
}
