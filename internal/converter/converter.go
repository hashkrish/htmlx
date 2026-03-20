package converter

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/hashkrish/htmlx/internal/markdown"
	"github.com/hashkrish/htmlx/internal/models"
)

// Converter converts HTML to Markdown
type Converter struct {
	options      *models.ConversionOptions
	textProcessor *TextProcessor
	builder      *markdown.Builder
	baseURL      *url.URL
	currentLevel int
	inList       bool
	inTable      bool
	inCodeBlock  bool
}

// NewConverter creates a new converter with the given options
func NewConverter(opts *models.ConversionOptions) *Converter {
	if opts == nil {
		opts = models.DefaultConversionOptions()
	}

	var baseURL *url.URL
	if opts.BaseURL != "" {
		var err error
		baseURL, err = url.Parse(opts.BaseURL)
		if err != nil {
			baseURL = nil
		}
	}

	return &Converter{
		options:       opts,
		textProcessor: NewTextProcessor(opts.PreserveWhitespace),
		builder:       markdown.NewBuilder(),
		baseURL:       baseURL,
		currentLevel:  0,
	}
}

// Convert converts HTML to Markdown
func (c *Converter) Convert(htmlContent string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %w", err)
	}

	// Find body or use document
	body := doc.Find("body")
	if body.Length() == 0 {
		body = doc.Selection
	}

	// Process all children of the body
	c.processChildren(body)

	return c.builder.String(), nil
}

// processNode processes an HTML node
func (c *Converter) processNode(sel *goquery.Selection) {
	// Check depth limit
	if c.currentLevel >= c.options.MaxNestingLevel {
		return
	}

	node := sel.Get(0)
	if node == nil {
		return
	}

	switch node.Type {
	case 1: // Text node (golang.org/x/net/html)
		c.processTextNode(sel)
	case 3: // Element node (golang.org/x/net/html)
		c.processElement(sel)
	}
}

// processTextNode processes text nodes
func (c *Converter) processTextNode(sel *goquery.Selection) {
	text := sel.Text()
	if c.textProcessor.IsEmptyText(text) {
		return
	}

	processedText := c.textProcessor.ProcessText(text)
	if c.textProcessor.IsEmptyText(processedText) {
		return
	}

	c.builder.WriteText(processedText)
}

// processElement processes HTML elements
func (c *Converter) processElement(sel *goquery.Selection) {
	tagName := goquery.NodeName(sel)

	switch tagName {
	// Headings
	case "h1":
		c.processHeading(sel, 1)
	case "h2":
		c.processHeading(sel, 2)
	case "h3":
		c.processHeading(sel, 3)
	case "h4":
		c.processHeading(sel, 4)
	case "h5":
		c.processHeading(sel, 5)
	case "h6":
		c.processHeading(sel, 6)

	// Block elements
	case "p":
		c.processParagraph(sel)
	case "div", "section", "article", "main":
		c.processContainer(sel)
	case "blockquote":
		c.processBlockquote(sel)

	// Lists
	case "ul", "ol":
		c.processList(sel, tagName == "ol")

	// Links and images
	case "a":
		c.processLink(sel)
	case "img":
		c.processImage(sel)

	// Code
	case "code", "pre":
		c.processCode(sel, tagName == "pre")

	// Emphasis
	case "strong", "b":
		c.processStrong(sel)
	case "em", "i":
		c.processEmphasis(sel)

	// Tables
	case "table":
		c.processTable(sel)

	// Semantic elements - skip but process children
	case "header", "footer", "nav", "aside":
		c.processChildren(sel)

	// Semantic inline elements
	case "time", "address":
		c.processParagraph(sel)

	// Figure with caption
	case "figure":
		c.processFigure(sel)

	// Details/summary
	case "details":
		c.processDetails(sel)

	// Marked/highlighted text
	case "mark":
		c.processMarkText(sel)

	// Form elements
	case "form":
		c.processForm(sel)

	// Ignore these
	case "script", "style", "noscript":
		return

	// HTML5 semantic elements that should be skipped but children processed
	case "figcaption", "legend", "summary":
		// These are handled by their parent elements, skip to avoid duplication
		return

	// Obsolete or uncommon elements
	case "bdi":
		// Bidirectional text - just extract text
		c.processChildren(sel)
	case "kbd":
		// Keyboard input - format as code
		text := c.extractText(sel)
		if !c.textProcessor.IsEmptyText(text) {
			c.builder.WriteCodeInline(text)
		}
	case "samp":
		// Sample output - format as code
		text := c.extractText(sel)
		if !c.textProcessor.IsEmptyText(text) {
			c.builder.WriteCodeInline(text)
		}
	case "var":
		// Variable - format as emphasis
		text := c.extractText(sel)
		if !c.textProcessor.IsEmptyText(text) {
			c.builder.WriteEmphasis(text, false)
		}
	case "small":
		// Small text - treat as normal text
		c.processChildren(sel)
	case "sub", "sup":
		// Subscript/superscript - just extract text (markdown doesn't support these)
		c.processChildren(sel)

	// Default: process children
	default:
		c.processChildren(sel)
	}
}

// processHeading processes heading elements
func (c *Converter) processHeading(sel *goquery.Selection, level int) {
	text := c.extractText(sel)
	if c.textProcessor.IsEmptyText(text) {
		return
	}
	c.builder.WriteHeading(level, text)
}

// processParagraph processes paragraph elements
func (c *Converter) processParagraph(sel *goquery.Selection) {
	// Check if paragraph is empty
	if sel.Contents().Length() == 0 {
		return
	}

	// Check if paragraph contains only whitespace
	text := strings.TrimSpace(sel.Text())
	if text == "" {
		return
	}

	// Process children to preserve structure (links, emphasis, etc.)
	c.processParagraphContents(sel)
	c.builder.WriteBlankLine()
}

// processParagraphContents processes the contents of a paragraph while preserving structure
func (c *Converter) processParagraphContents(sel *goquery.Selection) {
	sel.Contents().Each(func(_ int, s *goquery.Selection) {
		node := s.Get(0)
		if node == nil {
			return
		}

		tagName := goquery.NodeName(s)
		switch tagName {
		case "a":
			c.processLink(s)
		case "strong", "b":
			c.processStrong(s)
		case "em", "i":
			c.processEmphasis(s)
		case "code":
			c.processCode(s, false)
		case "kbd", "samp":
			// Format as inline code
			text := c.extractText(s)
			if !c.textProcessor.IsEmptyText(text) {
				c.builder.WriteCodeInline(text)
			}
		case "var":
			// Format as emphasis
			text := c.extractText(s)
			if !c.textProcessor.IsEmptyText(text) {
				c.builder.WriteEmphasis(text, false)
			}
		case "mark":
			// Format as strong (since markdown doesn't have highlighting)
			text := c.extractText(s)
			if !c.textProcessor.IsEmptyText(text) {
				c.builder.WriteEmphasis(text, true)
			}
		case "br":
			c.builder.WriteNewline()
		default:
			// For text nodes, preserve spacing but normalize excess whitespace
			if node.Type == 1 { // Text node
				text := s.Text()
				// Normalize whitespace while preserving leading/trailing spaces
				text = markdown.NormalizeWhitespacePreserveEnds(text)
				if text != "" {
					// Use WriteTextRaw to preserve spaces around inline elements
					c.builder.WriteTextRaw(text)
				}
			} else {
				// For other elements, recursively process
				c.processParagraphContents(s)
			}
		}
	})
}

// processContainer processes container elements (div, section, etc.)
func (c *Converter) processContainer(sel *goquery.Selection) {
	c.currentLevel++
	c.processChildren(sel)
	c.currentLevel--
}

// processBlockquote processes blockquote elements
func (c *Converter) processBlockquote(sel *goquery.Selection) {
	if sel.Contents().Length() == 0 {
		return
	}

	// Build blockquote content while preserving emphasis and links
	var quoteContent strings.Builder
	sel.Contents().Each(func(_ int, s *goquery.Selection) {
		node := s.Get(0)
		if node == nil {
			return
		}

		tagName := goquery.NodeName(s)
		switch tagName {
		case "strong", "b":
			// Extract text and add emphasis
			text := c.extractText(s)
			if !c.textProcessor.IsEmptyText(text) {
				// Text is already processed (no escaping needed for content inside emphasis)
				quoteContent.WriteString("**")
				quoteContent.WriteString(text)
				quoteContent.WriteString("**")
			}
		case "em", "i":
			text := c.extractText(s)
			if !c.textProcessor.IsEmptyText(text) {
				quoteContent.WriteString("*")
				quoteContent.WriteString(text)
				quoteContent.WriteString("*")
			}
		case "a":
			href, exists := s.Attr("href")
			if exists && href != "" {
				text := c.extractText(s)
				href = c.resolveURL(href)
				quoteContent.WriteString("[")
				quoteContent.WriteString(text)
				quoteContent.WriteString("](")
				quoteContent.WriteString(href)
				quoteContent.WriteString(")")
			}
		default:
			if node.Type == 1 { // Text node
				text := s.Text()
				text = markdown.NormalizeWhitespacePreserveEnds(text)
				if !c.textProcessor.IsEmptyText(text) {
					quoteContent.WriteString(text)
				}
			} else {
				// For other elements, just extract text
				text := c.extractText(s)
				if !c.textProcessor.IsEmptyText(text) {
					quoteContent.WriteString(text)
				}
			}
		}
	})

	content := quoteContent.String()
	if !c.textProcessor.IsEmptyText(content) {
		c.builder.WriteBlockquote(content)
	}
}

// processList processes list elements with proper nesting
func (c *Converter) processList(sel *goquery.Selection, ordered bool) {
	c.processListItems(sel, ordered, c.currentLevel)
	if c.currentLevel == 0 {
		c.builder.WriteBlankLine()
	}
}

// processListItems recursively processes list items
func (c *Converter) processListItems(sel *goquery.Selection, ordered bool, level int) {
	sel.Children().FilterFunction(func(_ int, s *goquery.Selection) bool {
		return goquery.NodeName(s) == "li"
	}).Each(func(_ int, li *goquery.Selection) {
		c.processListItem(li, ordered, level)
	})
}

// processListItem processes a single list item
func (c *Converter) processListItem(sel *goquery.Selection, ordered bool, level int) {
	// Extract text from direct children (not nested lists)
	var text strings.Builder
	sel.Contents().Each(func(_ int, s *goquery.Selection) {
		node := s.Get(0)
		if node == nil {
			return
		}

		tagName := goquery.NodeName(s)

		// Skip nested lists for now, we'll handle them separately
		if tagName == "ul" || tagName == "ol" {
			return
		}

		if node.Type == 1 { // Text node (golang.org/x/net/html)
			text.WriteString(c.textProcessor.ProcessText(s.Text()))
		} else {
			text.WriteString(c.extractText(s))
		}
	})

	itemText := strings.TrimSpace(text.String())
	if itemText != "" {
		c.builder.WriteListItem(itemText, ordered, level)
	}

	// Now process nested lists
	sel.Contents().Each(func(_ int, s *goquery.Selection) {
		tagName := goquery.NodeName(s)
		if tagName == "ul" {
			c.processListItems(s, false, level+1)
		} else if tagName == "ol" {
			c.processListItems(s, true, level+1)
		}
	})
}

// processLink processes link elements
func (c *Converter) processLink(sel *goquery.Selection) {
	href, exists := sel.Attr("href")
	if !exists || href == "" {
		c.processChildren(sel)
		return
	}

	// Resolve relative URLs
	href = c.resolveURL(href)

	text := c.extractText(sel)
	if c.textProcessor.IsEmptyText(text) {
		text = href
	}

	c.builder.WriteLink(text, href)
}

// processImage processes image elements
func (c *Converter) processImage(sel *goquery.Selection) {
	if !c.options.IncludeImages {
		return
	}

	alt, _ := sel.Attr("alt")
	if alt == "" {
		alt = "Image"
	}

	src, exists := sel.Attr("src")
	if !exists || src == "" {
		return
	}

	src = c.resolveURL(src)
	c.builder.WriteLink(alt, src)
}

// processCode processes code elements
func (c *Converter) processCode(sel *goquery.Selection, block bool) {
	text := sel.Text()
	if c.textProcessor.IsEmptyText(text) {
		return
	}

	if block {
		// For <pre> blocks, try to find class on child <code> element
		class, _ := sel.Attr("class")
		if class == "" {
			// Try to find class on <code> child
			sel.Find("code").Each(func(_ int, code *goquery.Selection) {
				if childClass, exists := code.Attr("class"); exists {
					class = childClass
				}
			})
		}
		language := c.textProcessor.ExtractLanguageFromClass(class)
		c.builder.WriteCodeBlock(text, language)
	} else {
		c.builder.WriteCodeInline(text)
	}
}

// processStrong processes strong/b elements
func (c *Converter) processStrong(sel *goquery.Selection) {
	text := c.extractText(sel)
	if c.textProcessor.IsEmptyText(text) {
		return
	}
	c.builder.WriteEmphasis(text, true)
}

// processEmphasis processes em/i elements
func (c *Converter) processEmphasis(sel *goquery.Selection) {
	text := c.extractText(sel)
	if c.textProcessor.IsEmptyText(text) {
		return
	}
	c.builder.WriteEmphasis(text, false)
}

// processTable processes table elements
func (c *Converter) processTable(sel *goquery.Selection) {
	c.inTable = true

	// Extract headers
	var headers []string
	sel.Find("thead th").Each(func(_ int, th *goquery.Selection) {
		headers = append(headers, c.extractText(th))
	})

	// If no thead, try to find headers in first row
	if len(headers) == 0 {
		sel.Find("tr").First().Find("th").Each(func(_ int, th *goquery.Selection) {
			headers = append(headers, c.extractText(th))
		})
	}

	// If still no headers, use first row as header
	if len(headers) == 0 {
		sel.Find("tr").First().Find("td").Each(func(_ int, td *goquery.Selection) {
			headers = append(headers, c.extractText(td))
		})
	}

	if len(headers) == 0 {
		c.inTable = false
		return
	}

	c.builder.WriteTableStart(headers)

	// Extract rows
	sel.Find("tr").Each(func(_ int, tr *goquery.Selection) {
		// Skip header row
		if tr.Find("th").Length() > 0 {
			return
		}

		var row []string
		tr.Find("td").Each(func(_ int, td *goquery.Selection) {
			row = append(row, c.extractText(td))
		})

		// Pad row if necessary
		for len(row) < len(headers) {
			row = append(row, "")
		}

		c.builder.WriteTableRow(row)
	})

	c.builder.WriteTableEnd()
	c.inTable = false
}

// processChildren processes all child nodes
func (c *Converter) processChildren(sel *goquery.Selection) {
	sel.Contents().Each(func(_ int, s *goquery.Selection) {
		c.processNode(s)
	})
}

// extractText extracts all text from an element (including children)
func (c *Converter) extractText(sel *goquery.Selection) string {
	var text strings.Builder

	sel.Contents().Each(func(_ int, s *goquery.Selection) {
		node := s.Get(0)
		if node == nil {
			return
		}

		if node.Type == 1 { // Text node (golang.org/x/net/html)
			text.WriteString(c.textProcessor.ProcessText(s.Text()))
		} else {
			// Recursively extract text from child elements
			text.WriteString(c.extractText(s))
		}
	})

	return markdown.NormalizeWhitespace(text.String())
}

// resolveURL resolves relative URLs against the base URL
func (c *Converter) resolveURL(rawURL string) string {
	if rawURL == "" {
		return ""
	}

	// If it's already an absolute URL, return as-is
	if strings.HasPrefix(rawURL, "http://") || strings.HasPrefix(rawURL, "https://") ||
		strings.HasPrefix(rawURL, "//") {
		return rawURL
	}

	// If we have a base URL, resolve relative to it
	if c.baseURL != nil {
		resolved, err := c.baseURL.Parse(rawURL)
		if err == nil {
			return resolved.String()
		}
	}

	// Return as-is if we can't resolve
	return rawURL
}

// isListElement checks if a tag is a list-related element
func isListElement(tagName string) bool {
	return tagName == "ul" || tagName == "ol" || tagName == "li"
}

// processFigure processes figure elements with optional captions
func (c *Converter) processFigure(sel *goquery.Selection) {
	// Find figcaption if it exists
	caption := sel.Find("figcaption").First()
	captionText := ""
	if caption.Length() > 0 {
		captionText = strings.TrimSpace(caption.Text())
	}

	// Process children (images, etc.)
	sel.Contents().Each(func(_ int, s *goquery.Selection) {
		node := s.Get(0)
		if node == nil {
			return
		}

		tagName := goquery.NodeName(s)

		// Skip figcaption as we handled it separately
		if tagName == "figcaption" {
			return
		}

		c.processElement(s)
	})

	// Add caption if present (as emphasized text)
	if captionText != "" {
		c.builder.WriteEmphasis(captionText, false)
		c.builder.WriteBlankLine()
	}
}

// processDetails processes details/summary elements
func (c *Converter) processDetails(sel *goquery.Selection) {
	// Find summary if it exists
	summary := sel.Find("summary").First()
	summaryText := ""
	if summary.Length() > 0 {
		summaryText = strings.TrimSpace(summary.Text())
	}

	// Write summary as bold (expandable in HTML, appears as emphasized in markdown)
	if summaryText != "" {
		c.builder.WriteEmphasis(summaryText, true)
		c.builder.WriteNewline()
	}

	// Process other children
	sel.Contents().Each(func(_ int, s *goquery.Selection) {
		node := s.Get(0)
		if node == nil {
			return
		}

		tagName := goquery.NodeName(s)

		// Skip summary as we handled it separately
		if tagName == "summary" {
			return
		}

		c.processElement(s)
	})
}

// processMarkText processes mark elements (highlighted text)
func (c *Converter) processMarkText(sel *goquery.Selection) {
	text := c.extractText(sel)
	if c.textProcessor.IsEmptyText(text) {
		return
	}

	// Mark text with emphasis (no direct markdown equivalent for highlighting)
	c.builder.WriteEmphasis(text, true)
}
