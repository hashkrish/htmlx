# htmlx - HTML to Markdown Converter

Convert HTML documents to clean, semantic markdown optimized for AI/LLM consumption. `htmlx` preserves the semantic structure of web content while minimizing token usage for efficient processing by language models.

## Features

- **Multiple Input Methods**: Fetch from URLs, local files, or stdin
- **Semantic Preservation**: Maintains document structure (headings, lists, tables, links)
- **Styling Discarded**: Focuses on content semantics, not visual styling
- **LLM-Optimized**: Minimizes tokens while preserving meaning
- **Fast Processing**: Efficient static HTML parsing without JavaScript execution
- **GitHub Flavored Markdown**: Compatible with modern markdown tools and LLMs

## Installation

### From Source

```bash
git clone https://github.com/hashkrish/htmlx.git
cd htmlx
make install
```

### Using Go Install

```bash
go install github.com/hashkrish/htmlx/cmd/htmlx@latest
```

## Usage

### Basic Examples

Convert a URL to markdown:
```bash
htmlx https://example.com
```

Convert a local HTML file:
```bash
htmlx /path/to/file.html
```

Read from stdin:
```bash
curl -s https://example.com | htmlx
```

Save to a file:
```bash
htmlx -o output.md https://example.com
```

### Command Line Options

```
htmlx [flags] [source]

Flags:
  -o, --output string       Output file (default: stdout)
  -b, --base-url string     Base URL for resolving relative links
  -v, --verbose             Enable verbose logging
  --no-images              Exclude image references
  --no-forms               Exclude form element documentation
  --timeout duration       HTTP request timeout (default: 30s)
  --user-agent string      Custom User-Agent header
  --version                Show version information
  -h, --help               Help for htmlx
```

### Examples

With custom base URL:
```bash
htmlx --base-url https://example.com /path/to/file.html
```

Verbose mode for debugging:
```bash
htmlx -v https://example.com
```

Custom HTTP timeout:
```bash
htmlx --timeout 60s https://slow-website.com
```

## Supported HTML Elements

- **Headings**: h1-h6 → `# Markdown headers`
- **Paragraphs**: p → Text blocks with blank lines
- **Lists**: ul, ol → `-` or `1.` with proper nesting
- **Links**: a → `[text](url)`
- **Code**: code, pre → Fenced code blocks with language detection
- **Tables**: table → Markdown tables with proper formatting
- **Emphasis**: strong, b → `**bold**`; em, i → `*italic*`
- **Blockquotes**: blockquote → `>` quoted text
- **Forms**: input, select → Structured documentation
- **Semantic Elements**: article, section, nav, etc.

## Project Structure

```
htmlx/
├── cmd/htmlx/           # CLI application
├── internal/
│   ├── converter/       # Core HTML to Markdown conversion
│   ├── fetcher/         # Content fetching (URL, file, stdin)
│   ├── markdown/        # Markdown output building
│   └── models/          # Configuration structures
├── test/                # Tests and fixtures
├── Makefile             # Build automation
└── README.md            # This file
```

## Building

```bash
# Build the binary
make build

# Run tests
make test

# Run benchmarks
make bench

# Clean build artifacts
make clean
```

## Development

The project is organized into logical packages:

- **cmd/htmlx**: Entry point and CLI handling
- **internal/converter**: Core HTML-to-Markdown conversion engine
- **internal/fetcher**: Abstraction for content retrieval (URL, file, stdin)
- **internal/markdown**: Markdown building and formatting
- **internal/models**: Configuration and options

## Testing

Run all tests:
```bash
make test
```

Run specific test:
```bash
go test -v -run TestConverterHeading ./test
```

Run with coverage:
```bash
go test -cover ./...
```

## Performance

Performance targets:
- 1MB HTML document: < 100ms conversion time
- Memory usage: < 2x input size
- Supports 10+ concurrent conversions

## Known Limitations

- **JavaScript Rendering**: htmlx performs static HTML parsing. Pages that require JavaScript execution will not include dynamically rendered content.
- **Complex Forms**: Form elements are documented as structured text rather than preserved as interactive elements.
- **CSS Styling**: All CSS styling is discarded; only semantic structure is preserved.

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Roadmap

### v0.2.0
- Form element support improvements
- Better code block language detection
- Image handling enhancements

### v0.3.0
- Configuration file support (.htmlxrc)
- Parallel URL processing
- Performance optimizations

### v1.0.0
- Stable API
- Comprehensive documentation
- Wide platform support (Linux, macOS, Windows, ARM64/AMD64)
- Homebrew formula
- Docker image

## Feedback

For issues, suggestions, or feedback, please visit: https://github.com/hashkrish/htmlx/issues
