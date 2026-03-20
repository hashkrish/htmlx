package test

import (
	"fmt"
	"testing"

	"github.com/hashkrish/htmlx/internal/converter"
	"github.com/hashkrish/htmlx/internal/models"
)

func TestCodeBlockLanguageDetection(t *testing.T) {
	html := `<html><body>
<pre><code class="language-python">
def hello():
    print("Hello")
</code></pre>
</body></html>`

	opts := models.DefaultConversionOptions()
	conv := converter.NewConverter(opts)

	result, err := conv.Convert(html)
	if err != nil {
		t.Fatalf("Conversion failed: %v", err)
	}

	fmt.Printf("Result:\n%s\n", result)

	if !contains(result, "```python") {
		t.Error("Expected '```python' in output, got:")
		fmt.Printf("%q\n", result)
	}
}
