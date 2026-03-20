package test

import (
	"strings"
	"testing"

	"github.com/hashkrish/htmlx/internal/converter"
	"github.com/hashkrish/htmlx/internal/models"
)

// BenchmarkConverterSimpleHTML benchmarks simple HTML conversion
func BenchmarkConverterSimpleHTML(b *testing.B) {
	html := `<h1>Test Title</h1><p>This is a paragraph with <strong>bold</strong> and <em>italic</em> text.</p>`
	opts := models.DefaultConversionOptions()
	conv := converter.NewConverter(opts)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		conv.Convert(html)
	}
}

// BenchmarkConverterLargeHTML benchmarks larger HTML conversion
func BenchmarkConverterLargeHTML(b *testing.B) {
	// Create a large HTML document
	var html strings.Builder
	html.WriteString("<html><body>")
	for i := 0; i < 100; i++ {
		html.WriteString("<h2>Section ")
		html.WriteString(string(rune(48 + (i % 10))))
		html.WriteString("</h2>")
		for j := 0; j < 5; j++ {
			html.WriteString("<p>This is paragraph ")
			html.WriteString(string(rune(48 + (j % 10))))
			html.WriteString(" with some <strong>bold</strong> and <em>italic</em> text.</p>")
		}
		html.WriteString("<ul>")
		for j := 0; j < 5; j++ {
			html.WriteString("<li>List item ")
			html.WriteString(string(rune(48 + (j % 10))))
			html.WriteString("</li>")
		}
		html.WriteString("</ul>")
	}
	html.WriteString("</body></html>")

	htmlStr := html.String()
	opts := models.DefaultConversionOptions()
	conv := converter.NewConverter(opts)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		conv.Convert(htmlStr)
	}
}

// BenchmarkConverterTable benchmarks table conversion
func BenchmarkConverterTable(b *testing.B) {
	// Create HTML with a table
	var html strings.Builder
	html.WriteString("<table><thead><tr><th>Col1</th><th>Col2</th><th>Col3</th></tr></thead><tbody>")
	for i := 0; i < 50; i++ {
		html.WriteString("<tr><td>Cell 1</td><td>Cell 2</td><td>Cell 3</td></tr>")
	}
	html.WriteString("</tbody></table>")

	htmlStr := html.String()
	opts := models.DefaultConversionOptions()
	conv := converter.NewConverter(opts)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		conv.Convert(htmlStr)
	}
}

// BenchmarkConverterNestedLists benchmarks nested list conversion
func BenchmarkConverterNestedLists(b *testing.B) {
	// Create deeply nested lists
	var html strings.Builder
	html.WriteString("<ul><li>Item 1")
	for i := 0; i < 5; i++ {
		html.WriteString("<ul><li>Nested item 1")
	}
	for i := 0; i < 5; i++ {
		html.WriteString("</li></ul>")
	}
	html.WriteString("</li></ul>")

	htmlStr := html.String()
	opts := models.DefaultConversionOptions()
	conv := converter.NewConverter(opts)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		conv.Convert(htmlStr)
	}
}
