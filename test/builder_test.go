package test

import (
	"testing"

	"github.com/hashkrish/htmlx/internal/markdown"
)

func TestBuilderHeading(t *testing.T) {
	b := markdown.NewBuilder()
	b.WriteHeading(1, "Hello")
	b.WriteLine("This is a paragraph.")

	result := b.String()
	t.Logf("Result:\n%s", result)

	if !contains(result, "# Hello") {
		t.Error("Expected '# Hello' in output")
	}

	if !contains(result, "This is a paragraph.") {
		t.Error("Expected paragraph in output")
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
