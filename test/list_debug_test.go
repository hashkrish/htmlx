package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/hashkrish/htmlx/internal/converter"
	"github.com/hashkrish/htmlx/internal/models"
)

func TestListConversion(t *testing.T) {
	html := `<html><body>
<h2>Items</h2>
<ul>
<li>Item 1</li>
<li>Item 2</li>
</ul>
</body></html>`

	opts := models.DefaultConversionOptions()
	conv := converter.NewConverter(opts)

	result, err := conv.Convert(html)
	if err != nil {
		t.Fatalf("Conversion failed: %v", err)
	}

	t.Logf("Result:\n%s\n---", result)

	if !contains(result, "Item 1") {
		t.Error("Expected 'Item 1' in output")
	}
	if !contains(result, "Item 2") {
		t.Error("Expected 'Item 2' in output")
	}
}

func TestListDebug(t *testing.T) {
	html := `<ul><li>Item 1</li><li>Item 2</li></ul>`

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	ul := doc.Find("ul")
	fmt.Printf("Found ul: %d\n", ul.Length())

	lis := ul.Find("> li")
	fmt.Printf("Found li elements with Find(> li): %d\n", lis.Length())

	ul.Contents().Each(func(idx int, s *goquery.Selection) {
		node := s.Get(0)
		if node == nil {
			return
		}
		fmt.Printf("  [%d] Tag: %s, Type: %d\n", idx, goquery.NodeName(s), node.Type)
	})
}
