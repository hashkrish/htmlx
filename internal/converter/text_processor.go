package converter

import (
	"html"
	"regexp"
	"strings"
)

// TextProcessor handles text normalization and processing
type TextProcessor struct {
	preserveWhitespace bool
}

// NewTextProcessor creates a new text processor
func NewTextProcessor(preserveWhitespace bool) *TextProcessor {
	return &TextProcessor{
		preserveWhitespace: preserveWhitespace,
	}
}

// ProcessText normalizes and processes text content
func (tp *TextProcessor) ProcessText(text string) string {
	// Decode HTML entities
	text = html.UnescapeString(text)

	// Normalize whitespace unless we're preserving it
	if !tp.preserveWhitespace {
		text = tp.normalizeWhitespace(text)
	}

	return text
}

// normalizeWhitespace collapses multiple spaces and newlines into single spaces
func (tp *TextProcessor) normalizeWhitespace(text string) string {
	// Replace multiple spaces, tabs, and newlines with single space
	text = regexp.MustCompile(`[\s\n\r\t]+`).ReplaceAllString(text, " ")

	// Trim leading and trailing whitespace
	text = strings.TrimSpace(text)

	return text
}

// TrimParagraph removes leading and trailing whitespace, but preserves internal structure
func (tp *TextProcessor) TrimParagraph(text string) string {
	return strings.TrimSpace(text)
}

// IsEmptyText checks if text is empty or only contains whitespace
func (tp *TextProcessor) IsEmptyText(text string) bool {
	return len(strings.TrimSpace(text)) == 0
}

// ExtractLanguageFromClass extracts programming language from HTML class attribute
// e.g., "language-python" -> "python", "hljs python" -> "python"
func (tp *TextProcessor) ExtractLanguageFromClass(class string) string {
	if class == "" {
		return ""
	}

	// Look for "language-xxx" pattern
	if strings.Contains(class, "language-") {
		parts := strings.Split(class, " ")
		for _, part := range parts {
			if strings.HasPrefix(part, "language-") {
				return strings.TrimPrefix(part, "language-")
			}
		}
	}

	// Look for common language class names
	commonLanguages := map[string]bool{
		"python": true, "js": true, "javascript": true, "go": true, "rust": true,
		"java": true, "cpp": true, "c": true, "csharp": true, "ruby": true,
		"php": true, "swift": true, "kotlin": true, "typescript": true, "ts": true,
		"bash": true, "shell": true, "sh": true, "sql": true, "html": true,
		"css": true, "xml": true, "yaml": true, "json": true, "toml": true,
	}

	parts := strings.Fields(class)
	for _, part := range parts {
		part = strings.ToLower(part)
		if commonLanguages[part] {
			return part
		}
	}

	return ""
}

// CleanupText removes unnecessary whitespace and normalizes text for display
func (tp *TextProcessor) CleanupText(text string) string {
	// Decode HTML entities
	text = html.UnescapeString(text)

	// Replace multiple whitespace characters with a single space
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")

	// Trim leading/trailing whitespace
	text = strings.TrimSpace(text)

	return text
}

