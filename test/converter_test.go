package test

import (
	"testing"

	"github.com/hashkrish/htmlx/internal/converter"
	"github.com/hashkrish/htmlx/internal/models"
)

func TestConverterHeading(t *testing.T) {
	html := `<html><body><h1>Hello</h1><p>World</p></body></html>`

	opts := models.DefaultConversionOptions()
	conv := converter.NewConverter(opts)

	result, err := conv.Convert(html)
	if err != nil {
		t.Fatalf("Conversion failed: %v", err)
	}

	t.Logf("Result:\n%s", result)

	if !contains(result, "# Hello") {
		t.Error("Expected '# Hello' in output")
	}

	if !contains(result, "World") {
		t.Error("Expected 'World' in output")
	}
}

