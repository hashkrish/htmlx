package test

import (
	"fmt"
	"testing"

	"golang.org/x/net/html"
)

func TestNodeTypes(t *testing.T) {
	fmt.Printf("HTML Node Types:\n")
	fmt.Printf("  ErrorNode: %d\n", html.ErrorNode)
	fmt.Printf("  TextNode: %d\n", html.TextNode)
	fmt.Printf("  DocumentNode: %d\n", html.DocumentNode)
	fmt.Printf("  ElementNode: %d\n", html.ElementNode)
	fmt.Printf("  CommentNode: %d\n", html.CommentNode)
	fmt.Printf("  DoctypeNode: %d\n", html.DoctypeNode)
}
