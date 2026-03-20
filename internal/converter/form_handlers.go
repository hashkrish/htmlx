package converter

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// processForm processes form elements and documents them
func (c *Converter) processForm(sel *goquery.Selection) {
	if !c.options.IncludeForms {
		// Skip form processing if disabled
		c.processChildren(sel)
		return
	}

	c.builder.WriteHeading(3, "Form")

	sel.Contents().Each(func(_ int, s *goquery.Selection) {
		node := s.Get(0)
		if node == nil {
			return
		}

		tagName := goquery.NodeName(s)

		switch tagName {
		case "input":
			c.processInput(s)
		case "select":
			c.processSelect(s)
		case "textarea":
			c.processTextarea(s)
		case "label":
			c.processLabel(s)
		case "button":
			c.processButton(s)
		case "fieldset":
			c.processFieldset(s)
		default:
			// Process other elements normally
			c.processElement(s)
		}
	})

	c.builder.WriteBlankLine()
}

// processInput documents input elements
func (c *Converter) processInput(sel *goquery.Selection) {
	inputType, _ := sel.Attr("type")
	if inputType == "" {
		inputType = "text"
	}

	name, _ := sel.Attr("name")
	placeholder, _ := sel.Attr("placeholder")
	value, _ := sel.Attr("value")

	if name == "" {
		name = "Input"
	}

	// Build plain text version (without markdown syntax that would be escaped)
	var label strings.Builder
	label.WriteString(name)
	label.WriteString(" (type: ")
	label.WriteString(inputType)
	label.WriteString(")")

	if placeholder != "" {
		label.WriteString(" - placeholder: ")
		label.WriteString(placeholder)
	}

	if value != "" {
		label.WriteString(" - default: ")
		label.WriteString(value)
	}

	// Write as list item using plain builder
	c.builder.WriteListItem(label.String(), false, 1)
}

// processSelect documents select elements with options
func (c *Converter) processSelect(sel *goquery.Selection) {
	name, _ := sel.Attr("name")
	if name == "" {
		name = "Select"
	}

	c.builder.WriteListItem(name + " (dropdown)", false, 1)

	sel.Find("option").Each(func(_ int, option *goquery.Selection) {
		optValue, _ := option.Attr("value")
		optText := option.Text()
		optText = strings.TrimSpace(optText)

		if optValue != "" && optValue != optText {
			c.builder.WriteListItem(optText + " (value: " + optValue + ")", false, 2)
		} else {
			c.builder.WriteListItem(optText, false, 2)
		}
	})
}

// processTextarea documents textarea elements
func (c *Converter) processTextarea(sel *goquery.Selection) {
	name, _ := sel.Attr("name")
	rows, _ := sel.Attr("rows")
	placeholder, _ := sel.Attr("placeholder")

	if name == "" {
		name = "Textarea"
	}

	var label strings.Builder
	label.WriteString(name)
	label.WriteString(" (textarea")
	if rows != "" {
		label.WriteString(", ")
		label.WriteString(rows)
		label.WriteString(" rows")
	}
	label.WriteString(")")

	if placeholder != "" {
		label.WriteString(" - placeholder: ")
		label.WriteString(placeholder)
	}

	c.builder.WriteListItem(label.String(), false, 1)
}

// processLabel documents label elements
func (c *Converter) processLabel(sel *goquery.Selection) {
	// Check if this label contains checkboxes or radios
	inputs := sel.Find("input")
	inputs.Each(func(_ int, inputSel *goquery.Selection) {
		inputType, _ := inputSel.Attr("type")
		if inputType == "checkbox" || inputType == "radio" {
			// Extract text associated with this specific input
			labelText := c.extractTextAfterInput(sel, inputSel)

			c.processCheckboxOrRadioWithInput(inputSel, inputType, labelText)
		}
	})

	// Also process other content in the label (like non-checkbox/radio inputs)
	sel.Contents().Each(func(_ int, s *goquery.Selection) {
		node := s.Get(0)
		if node == nil {
			return
		}

		tagName := goquery.NodeName(s)

		// Skip inputs we already processed as checkboxes/radios
		if tagName == "input" {
			inputType, _ := s.Attr("type")
			if inputType == "checkbox" || inputType == "radio" {
				return
			}
			// Regular input in label - process normally
			c.processElement(s)
		} else if tagName != "" {
			// Process other non-text elements
			c.processElement(s)
		}
	})
}

// extractTextAfterInput extracts the text that follows a specific input element within a label
func (c *Converter) extractTextAfterInput(label *goquery.Selection, input *goquery.Selection) string {
	var result strings.Builder
	inputNode := input.Get(0)

	// Get all contents of the label
	contents := label.Contents()
	found := false
	collecting := false

	for i := 0; i < contents.Length(); i++ {
		node := contents.Eq(i).Get(0)
		if node == nil {
			continue
		}

		if node == inputNode {
			found = true
			collecting = true
			continue
		}

		if found && collecting {
			// Stop collecting when we hit another input
			if goquery.NodeName(contents.Eq(i)) == "input" {
				break
			}

			// Collect text nodes
			if node.Type == 1 { // Text node
				text := contents.Eq(i).Text()
				result.WriteString(text)
			}
		}
	}

	return strings.TrimSpace(result.String())
}

// processCheckboxOrRadioWithInput documents checkbox and radio inputs with direct input element
func (c *Converter) processCheckboxOrRadioWithInput(inputSel *goquery.Selection, inputType string, labelText string) {
	value, _ := inputSel.Attr("value")

	var icon string
	if inputType == "checkbox" {
		icon = "Checkbox"
	} else {
		icon = "Radio"
	}

	var label strings.Builder
	label.WriteString(icon)
	label.WriteString(": ")
	if labelText != "" {
		label.WriteString(labelText)
	}
	if value != "" && value != labelText {
		label.WriteString(" (value: ")
		label.WriteString(value)
		label.WriteString(")")
	}

	c.builder.WriteListItem(label.String(), false, 1)
}

// processButton documents button elements
func (c *Converter) processButton(sel *goquery.Selection) {
	buttonType, _ := sel.Attr("type")
	if buttonType == "" {
		buttonType = "button"
	}

	text := sel.Text()
	text = strings.TrimSpace(text)

	c.builder.WriteListItem("Button: " + text + " (type: " + buttonType + ")", false, 1)
}

// processFieldset documents fieldset elements
func (c *Converter) processFieldset(sel *goquery.Selection) {
	legend := sel.Find("legend").First()
	if legend.Length() > 0 {
		legendText := legend.Text()
		legendText = strings.TrimSpace(legendText)
		c.builder.WriteHeading(4, "Form Group: "+legendText)
	}

	// Process form elements within fieldset
	sel.Contents().Each(func(_ int, s *goquery.Selection) {
		node := s.Get(0)
		if node == nil {
			return
		}

		tagName := goquery.NodeName(s)

		switch tagName {
		case "input":
			c.processInput(s)
		case "select":
			c.processSelect(s)
		case "textarea":
			c.processTextarea(s)
		case "label":
			c.processLabel(s)
		case "legend":
			// Skip - already processed
		default:
			// Process other elements
			c.processElement(s)
		}
	})

	c.builder.WriteBlankLine()
}
