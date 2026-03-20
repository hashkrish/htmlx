package markdown

import (
	"strings"
)

// Builder efficiently builds markdown content with context tracking
type Builder struct {
	content    strings.Builder
	lastWasNewline bool
	inList           bool
	inTable          bool
	inCodeBlock      bool
	indentLevel      int
}

// NewBuilder creates a new markdown builder
func NewBuilder() *Builder {
	return &Builder{
		lastWasNewline: true,
	}
}

// WriteText writes plain text with markdown escaping (with trimming)
func (b *Builder) WriteText(text string) {
	b.writeTextInternal(text, true)
}

// WriteTextRaw writes text without trimming leading/trailing spaces
func (b *Builder) WriteTextRaw(text string) {
	b.writeTextInternal(text, false)
}

// writeTextInternal is the internal implementation for text writing
func (b *Builder) writeTextInternal(text string, trim bool) {
	if trim {
		text = strings.TrimSpace(text)
	}

	if text == "" {
		return
	}

	// Escape markdown special characters
	text = escapeMarkdown(text)

	if b.lastWasNewline && b.content.Len() > 0 {
		b.content.WriteString(strings.Repeat("  ", b.indentLevel))
	}

	b.content.WriteString(text)
	b.lastWasNewline = false
}

// WriteLine writes a line of text with automatic newline
func (b *Builder) WriteLine(text string) {
	b.WriteText(text)
	b.WriteNewline()
}

// WriteNewline adds a single newline
func (b *Builder) WriteNewline() {
	if !b.lastWasNewline {
		b.content.WriteByte('\n')
		b.lastWasNewline = true
	}
}

// WriteBlankLine adds a blank line (ensures 2 newlines total)
func (b *Builder) WriteBlankLine() {
	if !b.lastWasNewline {
		b.content.WriteByte('\n')
	}
	// Always add one more newline to create blank line
	if b.content.Len() > 0 {
		b.content.WriteByte('\n')
	}
	b.lastWasNewline = true
}

// WriteHeading writes a markdown heading
func (b *Builder) WriteHeading(level int, text string) {
	text = strings.TrimSpace(text)
	if text == "" {
		return
	}

	// Ensure proper spacing before heading
	if b.content.Len() > 0 {
		b.WriteBlankLine()
	}

	prefix := strings.Repeat("#", level)
	b.content.WriteString(prefix)
	b.content.WriteByte(' ')
	b.content.WriteString(escapeMarkdown(text))
	b.WriteBlankLine()
}

// WriteListItem writes a list item with proper indentation
func (b *Builder) WriteListItem(text string, ordered bool, level int) {
	text = strings.TrimSpace(text)
	if text == "" {
		return
	}

	// Ensure we start on a new line
	b.WriteNewline()

	// Add indentation
	indent := strings.Repeat("  ", level)
	b.content.WriteString(indent)

	// Add list marker
	if ordered {
		b.content.WriteString("1. ")
	} else {
		b.content.WriteString("- ")
	}

	// Add text
	b.content.WriteString(escapeMarkdown(text))
	b.lastWasNewline = false
}

// WriteLink writes a markdown link
func (b *Builder) WriteLink(text, url string) {
	text = strings.TrimSpace(text)
	if text == "" {
		text = url
	}

	b.content.WriteString("[")
	b.content.WriteString(escapeMarkdown(text))
	b.content.WriteString("](")
	b.content.WriteString(url)
	b.content.WriteString(")")
	b.lastWasNewline = false
}

// WriteCodeInline writes inline code
func (b *Builder) WriteCodeInline(code string) {
	code = strings.TrimSpace(code)
	if code == "" {
		return
	}
	b.content.WriteString("`")
	b.content.WriteString(code)
	b.content.WriteString("`")
	b.lastWasNewline = false
}

// WriteCodeBlock writes a code block
func (b *Builder) WriteCodeBlock(code string, language string) {
	code = strings.TrimSpace(code)
	if code == "" {
		return
	}

	if b.content.Len() > 0 {
		b.WriteBlankLine()
	}

	b.content.WriteString("```")
	if language != "" {
		b.content.WriteString(language)
	}
	b.content.WriteByte('\n')
	b.lastWasNewline = true

	b.content.WriteString(code)
	b.content.WriteByte('\n')
	b.lastWasNewline = true

	b.content.WriteString("```")
	b.WriteBlankLine()
}

// WriteBlockquote writes a blockquote
func (b *Builder) WriteBlockquote(text string) {
	lines := strings.Split(strings.TrimSpace(text), "\n")
	if b.content.Len() > 0 {
		b.WriteBlankLine()
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		b.content.WriteString("> ")
		// Don't escape - content should already be properly formatted markdown
		b.content.WriteString(line)
		b.content.WriteByte('\n')
	}
	b.lastWasNewline = true
	b.WriteBlankLine()
}

// WriteTableStart writes the beginning of a markdown table
func (b *Builder) WriteTableStart(headers []string) {
	if b.content.Len() > 0 {
		b.WriteBlankLine()
	}

	// Write header row
	b.content.WriteString("|")
	for _, header := range headers {
		b.content.WriteString(" ")
		b.content.WriteString(strings.TrimSpace(escapeMarkdown(header)))
		b.content.WriteString(" |")
	}
	b.content.WriteByte('\n')

	// Write separator row
	b.content.WriteString("|")
	for range headers {
		b.content.WriteString("---|")
	}
	b.content.WriteByte('\n')

	b.lastWasNewline = true
	b.inTable = true
}

// WriteTableRow writes a table row
func (b *Builder) WriteTableRow(cells []string) {
	b.content.WriteString("|")
	for _, cell := range cells {
		b.content.WriteString(" ")
		b.content.WriteString(strings.TrimSpace(escapeMarkdown(cell)))
		b.content.WriteString(" |")
	}
	b.content.WriteByte('\n')
	b.lastWasNewline = true
}

// WriteTableEnd ends the table
func (b *Builder) WriteTableEnd() {
	b.inTable = false
	b.WriteBlankLine()
}

// WriteEmphasis writes emphasized text
func (b *Builder) WriteEmphasis(text string, strong bool) {
	text = strings.TrimSpace(text)
	if text == "" {
		return
	}

	if strong {
		b.content.WriteString("**")
		b.content.WriteString(escapeMarkdown(text))
		b.content.WriteString("**")
	} else {
		b.content.WriteString("*")
		b.content.WriteString(escapeMarkdown(text))
		b.content.WriteString("*")
	}
	b.lastWasNewline = false
}

// String returns the built markdown content
func (b *Builder) String() string {
	return strings.TrimRight(b.content.String(), "\n")
}

// EscapeMarkdown escapes markdown special characters in text
func EscapeMarkdown(text string) string {
	return escapeMarkdown(text)
}

// escapeMarkdown escapes markdown special characters in text
func escapeMarkdown(text string) string {
	// Characters that need escaping in markdown (only when they could cause confusion)
	replacements := map[rune]string{
		'\\': "\\\\",
		'`':  "\\`",
		'*':  "\\*",
		'_':  "\\_",
		'[':  "\\[",
		']':  "\\]",
	}

	var result strings.Builder
	for _, r := range text {
		if escaped, ok := replacements[r]; ok {
			result.WriteString(escaped)
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// NormalizeWhitespace normalizes whitespace in text (with trimming)
func NormalizeWhitespace(text string) string {
	// Replace multiple spaces/tabs/newlines with single space
	text = strings.TrimSpace(text)
	fields := strings.Fields(text)
	return strings.Join(fields, " ")
}

// NormalizeWhitespacePreserveEnds normalizes internal whitespace but preserves leading/trailing spaces
func NormalizeWhitespacePreserveEnds(text string) string {
	// Get leading whitespace
	leadingSpace := ""
	for _, r := range text {
		if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
			leadingSpace = " "
		} else {
			break
		}
	}

	// Get trailing whitespace
	trailingSpace := ""
	for i := len(text) - 1; i >= 0; i-- {
		r := rune(text[i])
		if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
			trailingSpace = " "
		} else {
			break
		}
	}

	// Normalize internal whitespace
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return leadingSpace + trailingSpace
	}

	fields := strings.Fields(trimmed)
	normalized := strings.Join(fields, " ")

	return leadingSpace + normalized + trailingSpace
}

// TrimText trims whitespace from text while preserving internal structure
func TrimText(text string) string {
	return strings.TrimSpace(text)
}
