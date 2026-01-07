// gitbook2md is a tool to convert GitBook HTML pages to Markdown format.
// This tool is used for documentation purposes only and is not part of the main library.
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// GitBookParser is a specialized parser for GitBook HTML pages
type GitBookParser struct {
	debug bool
}

// NewGitBookParser creates a new GitBook parser
func NewGitBookParser(debug bool) *GitBookParser {
	return &GitBookParser{debug: debug}
}

// ParseFile parses a GitBook HTML file and returns Markdown
func (p *GitBookParser) ParseFile(filename string) (string, error) {
	file, err := os.Open(filename) // #nosec G304
	if err != nil {
		return "", err
	}
	defer func() { _ = file.Close() }()

	doc, err := html.Parse(file)
	if err != nil {
		return "", err
	}

	// Extract metadata
	title := p.extractTitle(doc)
	description := p.extractDescription(doc)

	// Find main content
	content := p.extractMainContent(doc)

	// Build markdown
	var markdown strings.Builder
	
	if title != "" {
		markdown.WriteString("# " + title + "\n\n")
	}
	
	if description != "" {
		markdown.WriteString("> " + description + "\n\n")
	}
	
	if content != "" {
		markdown.WriteString(content)
	}

	return markdown.String(), nil
}

// extractTitle extracts the page title
func (p *GitBookParser) extractTitle(n *html.Node) string {
	if n.Type == html.ElementNode && n.Data == "title" {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.TextNode {
				title := strings.TrimSpace(c.Data)
				// Remove " | J-Quants API" suffix
				if idx := strings.Index(title, " | "); idx > 0 {
					title = title[:idx]
				}
				return title
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if title := p.extractTitle(c); title != "" {
			return title
		}
	}

	return ""
}

// extractDescription extracts meta description
func (p *GitBookParser) extractDescription(n *html.Node) string {
	if n.Type == html.ElementNode && n.Data == "meta" {
		var name, content string
		for _, attr := range n.Attr {
			if attr.Key == "name" && attr.Val == "description" {
				name = attr.Val
			}
			if attr.Key == "content" {
				content = attr.Val
			}
		}
		if name == "description" && content != "" {
			return content
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if desc := p.extractDescription(c); desc != "" {
			return desc
		}
	}

	return ""
}

// extractMainContent extracts the main content area
func (p *GitBookParser) extractMainContent(n *html.Node) string {
	// Strategy: Look for specific patterns in the HTML
	// 1. Find script tags with id="__NEXT_DATA__" which contains JSON data
	// 2. Look for main content divs
	// 3. Extract text content from specific areas

	var content strings.Builder

	// Look for JSON data in script tags
	jsonData := p.extractJSONData(n)
	if jsonData != "" {
		// Parse JSON and extract relevant content
		content.WriteString(p.parseJSONContent(jsonData))
	}

	// Look for structured content
	p.extractStructuredContent(n, &content)

	return content.String()
}

// extractJSONData looks for __NEXT_DATA__ script tag
func (p *GitBookParser) extractJSONData(n *html.Node) string {
	if n.Type == html.ElementNode && n.Data == "script" {
		for _, attr := range n.Attr {
			if attr.Key == "id" && attr.Val == "__NEXT_DATA__" {
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					if c.Type == html.TextNode {
						return c.Data
					}
				}
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if data := p.extractJSONData(c); data != "" {
			return data
		}
	}

	return ""
}

// parseJSONContent extracts content from JSON data
func (p *GitBookParser) parseJSONContent(jsonStr string) string {
	// This is a simplified extraction
	// In a real implementation, you'd parse the JSON properly
	
	var content strings.Builder

	// Look for common patterns in the JSON
	// Extract request/response examples
	if strings.Contains(jsonStr, "```") {
		// Extract code blocks
		re := regexp.MustCompile("```[^`]*```")
		matches := re.FindAllString(jsonStr, -1)
		for _, match := range matches {
			content.WriteString(match + "\n\n")
		}
	}

	return content.String()
}

// extractStructuredContent extracts content from HTML structure
func (p *GitBookParser) extractStructuredContent(n *html.Node, content *strings.Builder) {
	// Skip script and style tags
	if n.Type == html.ElementNode && (n.Data == "script" || n.Data == "style" || n.Data == "noscript") {
		return
	}

	// Look for content patterns
	if n.Type == html.ElementNode {
		class := p.getAttr(n, "class")
		
		// Look for code blocks
		if n.Data == "pre" || strings.Contains(class, "code") {
			p.extractCodeBlock(n, content)
			return
		}

		// Look for tables
		if n.Data == "table" {
			p.extractTable(n, content)
			return
		}

		// Look for div-based tables
		if n.Data == "div" && p.getAttr(n, "role") == "table" {
			p.extractTable(n, content)
			return
		}

		// Look for headings
		if n.Data == "h1" || n.Data == "h2" || n.Data == "h3" || n.Data == "h4" {
			p.extractHeading(n, content)
			return
		}

		// Look for paragraphs with actual content
		if n.Data == "p" {
			text := p.extractText(n)
			if len(text) > 20 && !strings.Contains(text, "mask-image") {
				content.WriteString(text + "\n\n")
			}
			return
		}

		// Look for lists
		if n.Data == "ul" || n.Data == "ol" {
			p.extractList(n, content, n.Data == "ol")
			return
		}

		// Look for definition lists (used for API parameters and response fields)
		if n.Data == "dl" {
			p.extractDefinitionList(n, content)
			return
		}
	}

	// Process children
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		p.extractStructuredContent(c, content)
	}
}

// extractCodeBlock extracts code blocks
func (p *GitBookParser) extractCodeBlock(n *html.Node, content *strings.Builder) {
	code := p.extractText(n)
	if code != "" && !strings.Contains(code, "function") && !strings.Contains(code, "var ") {
		content.WriteString("```\n")
		content.WriteString(code)
		content.WriteString("\n```\n\n")
	}
}

// extractTable extracts and converts tables
func (p *GitBookParser) extractTable(n *html.Node, content *strings.Builder) {
	// Simple table extraction
	var headers []string
	var rows [][]string

	// Find headers
	p.findTableHeaders(n, &headers)
	
	// Find rows
	p.findTableRows(n, &rows)

	// Also check for div-based tables with role attributes
	if len(headers) == 0 && len(rows) == 0 {
		p.findDivTableContent(n, &headers, &rows)
	}

	if len(headers) > 0 {
		// Write headers
		content.WriteString("| ")
		for _, h := range headers {
			content.WriteString(h + " | ")
		}
		content.WriteString("\n|")
		for range headers {
			content.WriteString("------|")
		}
		content.WriteString("\n")

		// Write rows
		for _, row := range rows {
			if len(row) > 0 {
				content.WriteString("| ")
				for _, cell := range row {
					content.WriteString(cell + " | ")
				}
				content.WriteString("\n")
			}
		}
		content.WriteString("\n")
	}
}

// extractHeading extracts headings
func (p *GitBookParser) extractHeading(n *html.Node, content *strings.Builder) {
	level := n.Data[1] - '0'
	text := p.extractText(n)
	if text != "" {
		content.WriteString(strings.Repeat("#", int(level)) + " " + text + "\n\n")
	}
}

// extractList extracts lists
func (p *GitBookParser) extractList(n *html.Node, content *strings.Builder, ordered bool) {
	counter := 1
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "li" {
			text := p.extractText(c)
			if text != "" {
				if ordered {
					fmt.Fprintf(content, "%d. %s\n", counter, text)
					counter++
				} else {
					content.WriteString("- " + text + "\n")
				}
			}
		}
	}
	content.WriteString("\n")
}

// extractDefinitionList extracts definition lists (dl/dt/dd) used for API parameters
func (p *GitBookParser) extractDefinitionList(n *html.Node, content *strings.Builder) {
	// Find all definition items (can be wrapped in div or directly as dt/dd pairs)
	p.extractDefinitionItems(n, content)
	content.WriteString("\n")
}

// extractDefinitionItems recursively finds and extracts dt/dd pairs
func (p *GitBookParser) extractDefinitionItems(n *html.Node, content *strings.Builder) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			if c.Data == "div" {
				// Definition item wrapped in div
				p.extractSingleDefinition(c, content)
			} else if c.Data == "dt" {
				// Direct dt element - find the next dd
				term := p.extractDefinitionTerm(c)
				desc := ""
				// Look for the next dd sibling
				for dd := c.NextSibling; dd != nil; dd = dd.NextSibling {
					if dd.Type == html.ElementNode && dd.Data == "dd" {
						desc = p.extractText(dd)
						break
					}
				}
				if term != "" {
					content.WriteString(term)
					if desc != "" {
						content.WriteString("\n  " + desc)
					}
					content.WriteString("\n\n")
				}
			}
		}
	}
}

// extractSingleDefinition extracts a single dt/dd pair from a container div
func (p *GitBookParser) extractSingleDefinition(n *html.Node, content *strings.Builder) {
	var term, desc string

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			switch c.Data {
			case "dt":
				term = p.extractDefinitionTerm(c)
			case "dd":
				desc = p.extractText(c)
			}
		}
	}

	if term != "" {
		content.WriteString(term)
		if desc != "" {
			content.WriteString("\n  " + desc)
		}
		content.WriteString("\n\n")
	}
}

// extractDefinitionTerm extracts the term from a dt element
// Format: **name** (type, required/optional)
func (p *GitBookParser) extractDefinitionTerm(n *html.Node) string {
	var name, typeStr, requiredStr string

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "span" {
			text := p.extractText(c)
			class := p.getAttr(c, "class")

			if strings.Contains(class, "font-semibold") {
				name = text
			} else if strings.Contains(class, "text-gray") || strings.Contains(class, "bg-gray") {
				if typeStr == "" && isTypeText(text) {
					typeStr = text
				} else if text == "optional" {
					requiredStr = text
				}
			} else if strings.Contains(class, "text-red") || strings.Contains(class, "bg-red") {
				if text == "required" {
					requiredStr = text
				}
			} else if strings.Contains(class, "text-blue") || strings.Contains(class, "bg-blue") {
				if text == "optional" {
					requiredStr = text
				}
			}
		}
	}

	if name == "" {
		return ""
	}

	result := "**" + name + "**"
	if typeStr != "" || requiredStr != "" {
		result += " ("
		if typeStr != "" {
			result += typeStr
		}
		if requiredStr != "" {
			if typeStr != "" {
				result += ", "
			}
			result += requiredStr
		}
		result += ")"
	}

	return result
}

// Helper methods

// isTypeText checks if the text represents a type annotation
// Supports both simple types (string, number) and compound types (number | string)
func isTypeText(text string) bool {
	types := []string{"string", "number", "boolean", "array", "object"}
	for _, t := range types {
		if strings.Contains(text, t) {
			return true
		}
	}
	return false
}

func (p *GitBookParser) getAttr(n *html.Node, key string) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func (p *GitBookParser) extractText(n *html.Node) string {
	var text strings.Builder
	p.extractTextRecursive(n, &text)
	return strings.TrimSpace(text.String())
}

func (p *GitBookParser) extractTextRecursive(n *html.Node, text *strings.Builder) {
	if n.Type == html.TextNode {
		text.WriteString(n.Data)
	}

	// Handle checkbox input elements
	if n.Type == html.ElementNode && n.Data == "input" {
		inputType := p.getAttr(n, "type")
		if inputType == "checkbox" {
			// Check if the checkbox is checked
			isChecked := false
			for _, attr := range n.Attr {
				if attr.Key == "checked" {
					isChecked = true
					break
				}
			}
			if isChecked {
				text.WriteString("✓")
			} else {
				text.WriteString("✗")
			}
			return
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type != html.ElementNode || (c.Data != "script" && c.Data != "style") {
			p.extractTextRecursive(c, text)
		}
	}
}

func (p *GitBookParser) findTableHeaders(n *html.Node, headers *[]string) {
	if n.Type == html.ElementNode {
		if n.Data == "th" {
			*headers = append(*headers, p.extractText(n))
			return
		}
		if n.Data == "thead" {
			// Process only thead
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				p.findTableHeaders(c, headers)
			}
			return
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		p.findTableHeaders(c, headers)
	}
}

func (p *GitBookParser) findTableRows(n *html.Node, rows *[][]string) {
	if n.Type == html.ElementNode && n.Data == "tbody" {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode && c.Data == "tr" {
				var row []string
				for td := c.FirstChild; td != nil; td = td.NextSibling {
					if td.Type == html.ElementNode && td.Data == "td" {
						row = append(row, p.extractText(td))
					}
				}
				if len(row) > 0 {
					*rows = append(*rows, row)
				}
			}
		}
		return
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		p.findTableRows(c, rows)
	}
}

// findDivTableContent looks for div-based tables with role attributes
func (p *GitBookParser) findDivTableContent(n *html.Node, headers *[]string, rows *[][]string) {
	// Look for role="table"
	if n.Type == html.ElementNode && p.getAttr(n, "role") == "table" {
		p.extractDivTable(n, headers, rows)
		return
	}

	// Also check if we're looking at a data table section
	if n.Type == html.ElementNode && n.Data == "div" {
		// Check for headers with role="columnheader"
		p.findDivHeaders(n, headers)
		
		// Check for rows with role="row"
		p.findDivRows(n, rows)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if len(*headers) == 0 && len(*rows) == 0 {
			p.findDivTableContent(c, headers, rows)
		}
	}
}

// extractDivTable extracts content from div-based tables
func (p *GitBookParser) extractDivTable(n *html.Node, headers *[]string, rows *[][]string) {
	// Find all elements with role="columnheader" for headers
	p.findDivHeaders(n, headers)
	
	// Find all elements with role="row" for data rows
	p.findDivRows(n, rows)
}

// findDivHeaders finds div elements with role="columnheader"
func (p *GitBookParser) findDivHeaders(n *html.Node, headers *[]string) {
	if n.Type == html.ElementNode && p.getAttr(n, "role") == "columnheader" {
		*headers = append(*headers, p.extractText(n))
		return
	}

	// Look for rowgroup containing headers
	if n.Type == html.ElementNode && p.getAttr(n, "role") == "rowgroup" {
		// Check if this is a header rowgroup
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode && p.getAttr(c, "role") == "row" {
				// Check if this row contains columnheaders
				hasHeaders := false
				for gc := c.FirstChild; gc != nil; gc = gc.NextSibling {
					if gc.Type == html.ElementNode && p.getAttr(gc, "role") == "columnheader" {
						hasHeaders = true
						break
					}
				}
				if hasHeaders {
					// Extract headers from this row
					for gc := c.FirstChild; gc != nil; gc = gc.NextSibling {
						if gc.Type == html.ElementNode && p.getAttr(gc, "role") == "columnheader" {
							*headers = append(*headers, p.extractText(gc))
						}
					}
					return
				}
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if len(*headers) == 0 {
			p.findDivHeaders(c, headers)
		}
	}
}

// findDivRows finds div elements with role="row" containing data cells
func (p *GitBookParser) findDivRows(n *html.Node, rows *[][]string) {
	if n.Type == html.ElementNode && p.getAttr(n, "role") == "rowgroup" {
		// Skip header rowgroup
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode && p.getAttr(c, "role") == "row" {
				// Check if this row contains cells (not headers)
				var row []string
				hasCells := false
				for gc := c.FirstChild; gc != nil; gc = gc.NextSibling {
					if gc.Type == html.ElementNode && p.getAttr(gc, "role") == "cell" {
						hasCells = true
						row = append(row, p.extractText(gc))
					}
				}
				if hasCells && len(row) > 0 {
					*rows = append(*rows, row)
				}
			}
		}
		return
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		p.findDivRows(c, rows)
	}
}

func main() {
	var (
		url    = flag.String("url", "", "URL to fetch and convert")
		output = flag.String("output", "", "Output file path")
		debug  = flag.Bool("debug", false, "Enable debug mode")
	)
	flag.Parse()

	// Check if using old style arguments
	if flag.NArg() > 0 && *url == "" {
		// Old style: gitbook2md <input.html> [output.md]
		inputFile := flag.Arg(0)
		outputFile := ""
		if flag.NArg() > 1 {
			outputFile = flag.Arg(1)
		}
		
		parser := NewGitBookParser(*debug || os.Getenv("DEBUG") == "true")
		markdown, err := parser.ParseFile(inputFile)
		if err != nil {
			log.Fatal(err)
		}
		
		if outputFile != "" {
			err = os.WriteFile(outputFile, []byte(markdown), 0644) // #nosec G306
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Wrote %s\n", outputFile)
		} else {
			fmt.Print(markdown)
		}
		return
	}

	// New style with URL
	if *url == "" {
		fmt.Println("Usage:")
		fmt.Println("  gitbook2md <input.html> [output.md]")
		fmt.Println("  gitbook2md --url <url> [--output <output.md>]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Download the HTML content
	resp, err := http.Get(*url)
	if err != nil {
		log.Fatal("Failed to fetch URL:", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		log.Fatal("HTTP error:", resp.StatusCode)
	}

	// Create temporary file
	tmpFile, err := os.CreateTemp("", "gitbook2md-*.html") // #nosec G303
	if err != nil {
		log.Fatal("Failed to create temp file:", err)
	}
	defer func() { _ = os.Remove(tmpFile.Name()) }()
	defer func() { _ = tmpFile.Close() }()

	// Copy content to temp file
	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		log.Fatal("Failed to save HTML:", err)
	}
	_ = tmpFile.Close() // #nosec G104

	// Parse the file
	parser := NewGitBookParser(*debug)
	markdown, err := parser.ParseFile(tmpFile.Name())
	if err != nil {
		log.Fatal(err)
	}

	// Output
	if *output != "" {
		// Create directory if needed
		dir := filepath.Dir(*output)
		// #nosec G301 -- Documentation directory needs to be readable
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatal("Failed to create directory:", err)
		}
		
		err = os.WriteFile(*output, []byte(markdown), 0644) // #nosec G306
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Wrote %s\n", *output)
	} else {
		fmt.Print(markdown)
	}
}