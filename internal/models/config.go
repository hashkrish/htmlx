package models

import "time"

// ConversionOptions holds the configuration for HTML to Markdown conversion
type ConversionOptions struct {
	BaseURL           string        // Base URL for resolving relative links
	IncludeImages     bool          // Include image references
	IncludeForms      bool          // Include form element documentation
	PreserveWhitespace bool         // Preserve whitespace in <pre> tags
	MaxNestingLevel   int           // Maximum nesting depth
	Timeout           time.Duration // HTTP request timeout
	UserAgent         string        // Custom User-Agent header
}

// DefaultConversionOptions returns the default conversion options
func DefaultConversionOptions() *ConversionOptions {
	return &ConversionOptions{
		BaseURL:            "",
		IncludeImages:      true,
		IncludeForms:       true,
		PreserveWhitespace: true,
		MaxNestingLevel:    100,
		Timeout:            30 * time.Second,
		UserAgent:          "htmlx/1.0 (HTML to Markdown Converter)",
	}
}
