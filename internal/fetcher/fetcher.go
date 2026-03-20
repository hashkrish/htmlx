package fetcher

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/net/html/charset"
)

// Fetcher is an interface for retrieving HTML content from different sources
type Fetcher interface {
	Fetch(source string) ([]byte, error)
}

// URLFetcher fetches HTML from HTTP/HTTPS URLs
type URLFetcher struct {
	timeout   time.Duration
	userAgent string
}

// NewURLFetcher creates a new URLFetcher
func NewURLFetcher(timeout time.Duration, userAgent string) *URLFetcher {
	return &URLFetcher{
		timeout:   timeout,
		userAgent: userAgent,
	}
}

// Fetch retrieves HTML from the given URL
func (f *URLFetcher) Fetch(url string) ([]byte, error) {
	client := &http.Client{
		Timeout: f.timeout,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("User-Agent", f.userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	// Use charset detection to handle different encodings
	reader, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		return nil, fmt.Errorf("failed to create charset reader: %w", err)
	}

	return io.ReadAll(reader)
}

// FileFetcher reads HTML from local files
type FileFetcher struct{}

// NewFileFetcher creates a new FileFetcher
func NewFileFetcher() *FileFetcher {
	return &FileFetcher{}
}

// Fetch reads HTML from a file
func (f *FileFetcher) Fetch(path string) ([]byte, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return content, nil
}

// StdinFetcher reads HTML from standard input
type StdinFetcher struct{}

// NewStdinFetcher creates a new StdinFetcher
func NewStdinFetcher() *StdinFetcher {
	return &StdinFetcher{}
}

// Fetch reads HTML from stdin
func (f *StdinFetcher) Fetch(_ string) ([]byte, error) {
	content, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, fmt.Errorf("failed to read from stdin: %w", err)
	}
	return content, nil
}

// DetectFetcher returns an appropriate fetcher based on the source
func DetectFetcher(source string, timeout time.Duration, userAgent string) Fetcher {
	if source == "-" || source == "" {
		return NewStdinFetcher()
	}

	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		return NewURLFetcher(timeout, userAgent)
	}

	return NewFileFetcher()
}
