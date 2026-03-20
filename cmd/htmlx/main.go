package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/hashkrish/htmlx/internal/converter"
	"github.com/hashkrish/htmlx/internal/fetcher"
	"github.com/hashkrish/htmlx/internal/models"
)

// logf logs verbose output with timestamp
func logf(format string, args ...interface{}) {
	if verbose {
		fmt.Fprintf(os.Stderr, "[%s] %s\n", time.Now().Format("15:04:05"), fmt.Sprintf(format, args...))
	}
}

var (
	Version = "0.1.0"

	outputFile     string
	baseURL        string
	timeout        int
	userAgent      string
	includeImages  bool
	includeForms   bool
	verbose        bool
	noImages       bool
	noForms        bool
)

var rootCmd = &cobra.Command{
	Use:     "htmlx [source]",
	Short:   "Convert HTML to semantic markdown for AI/LLM consumption",
	Long:    "htmlx converts HTML documents to clean, semantic markdown optimized for AI/LLM processing.",
	Args:    cobra.MaximumNArgs(1),
	RunE:    run,
	Version: Version,
}

func init() {
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file (default: stdout)")
	rootCmd.Flags().StringVarP(&baseURL, "base-url", "b", "", "Base URL for resolving relative links")
	rootCmd.Flags().IntVar(&timeout, "timeout", 30, "HTTP request timeout in seconds")
	rootCmd.Flags().StringVar(&userAgent, "user-agent", "", "Custom User-Agent header")
	rootCmd.Flags().BoolVar(&noImages, "no-images", false, "Exclude image references")
	rootCmd.Flags().BoolVar(&noForms, "no-forms", false, "Exclude form element documentation")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
}

func run(cmd *cobra.Command, args []string) error {
	startTime := time.Now()
	var source string
	if len(args) > 0 {
		source = args[0]
	}

	// Determine source
	if source == "" {
		// Check if stdin is being piped
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			source = "-"
		} else {
			return fmt.Errorf("please provide a source (URL, file path, or pipe from stdin)")
		}
	}

	// Set options
	opts := models.DefaultConversionOptions()
	opts.BaseURL = baseURL
	opts.Timeout = time.Duration(timeout) * time.Second
	opts.IncludeImages = !noImages && includeImages
	opts.IncludeForms = !noForms && includeForms

	if userAgent != "" {
		opts.UserAgent = userAgent
	}

	// Create fetcher
	f := fetcher.DetectFetcher(source, opts.Timeout, opts.UserAgent)

	logf("Fetching content from: %s", source)

	// Fetch HTML
	fetchStart := time.Now()
	htmlContent, err := f.Fetch(source)
	if err != nil {
		return fmt.Errorf("failed to fetch content: %w", err)
	}

	fetchDuration := time.Since(fetchStart)
	logf("Fetched %d bytes in %v", len(htmlContent), fetchDuration)

	// Convert
	logf("Converting HTML to Markdown...")
	convertStart := time.Now()

	conv := converter.NewConverter(opts)
	markdown, err := conv.Convert(string(htmlContent))
	if err != nil {
		return fmt.Errorf("conversion failed: %w", err)
	}

	convertDuration := time.Since(convertStart)
	logf("Conversion completed in %v, output size: %d bytes", convertDuration, len(markdown))

	// Output
	if outputFile != "" {
		logf("Writing output to: %s", outputFile)

		if err := os.WriteFile(outputFile, []byte(markdown), 0644); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
	} else {
		fmt.Println(markdown)
	}

	totalDuration := time.Since(startTime)
	logf("Total time: %v", totalDuration)

	return nil
}

func main() {
	// Set default values for boolean flags that should be true by default
	includeImages = true
	includeForms = true

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
