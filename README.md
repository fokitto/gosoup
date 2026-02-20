# GoSoup

A convenient Go library for parsing and querying HTML documents, inspired by BeautifulSoup4 for Python.

GoSoup provides a simple and intuitive API for navigating and searching HTML documents. It's built on top of the `golang.org/x/net/html` library and offers a more user-friendly interface for common HTML parsing tasks.

## Installation

```bash
go get github.com/fokitto/gosoup
```

## Quick Start

```go
package main

import (
	"fmt"

	"github.com/fokitto/gosoup"
)

func main() {
	html := `
	<html>
		<body>
			<div class="container">
				<h1>Hello, World!</h1>
				<p>This is a paragraph.</p>
				<p>And this is a paragraph.</p>
			</div>
		</body>
	</html>
	`

	root, err := gosoup.ParseString(html)
	if err != nil {
		panic(err)
	}

	// Find the first h1 tag
	h1 := root.Find(gosoup.HasName("h1"))
	fmt.Println(h1.Text()) // Output: Hello, World!

	// Find all paragraphs without class
	paragraphs := root.FindAll(
        gosoup.All(
            gosoup.HasName("p"),
            gosoup.HasNoClass(),
        ),
    )
	for _, p := range paragraphs {
		fmt.Println(p.Text())
	}
}
```

## API Overview

### Parsing Functions

- **`Parse(reader io.Reader) (*Tag, error)`** - Parse HTML from an `io.Reader`
- **`ParseBytes(content []byte) (*Tag, error)`** - Parse HTML from a byte slice
- **`ParseString(content string) (*Tag, error)`** - Parse HTML from a string

### Nodes

GoSoup provides a `Node` interface for working with different types of DOM nodes:

- **`Tag`** - Represents an HTML element with:
  - `Name` - The tag name (e.g., "div", "p", "a")
  - `Attrs` - Map of attributes (key-value pairs)
- **`NavigableString`** - Represents raw text content in the HTML document (similar to BeautifulSoup4's NavigableString)

### Navigation Methods

- **`Parent() *Tag`** - Get the parent tag
- **`FirstChild() *Tag`** - Get the first child tag
- **`Children() []*Tag`** - Get all direct child tags
- **`Prev() *Tag`** - Get the previous sibling element
- **`Next() *Tag`** - Get the next sibling element
- **`IterNodes() iter.Seq[Node]`** - Iterate through all child nodes (both tags and text) using range loops

### Content Methods

- **`Text() string`** - Get the immediate text content of the tag
- **`FullText(sep ...string) string`** - Get all text content recursively (with optional separator)
- **`String() string`** - Render the tag and its children as HTML

### Search Methods

- **`Find(predicate Predicate) *Tag`** - Find the first element matching the predicate
- **`FindAll(predicate Predicate) []*Tag`** - Find all elements matching the predicate
- **`FindParent(predicate Predicate) *Tag`** - Find the first parent element matching the predicate

### DOM Manipulation

- **`Unwrap() Tag`** - Remove the tag from its parent

### Working with Nodes

The `IterNodes()` method allows you to iterate through all child nodes, including both tags and text content:

```go
html := `<div>Hello <b>World</b> and <i>Goodbye</i></div>`
root, _ := gosoup.ParseString(html)
div := root.Find(gosoup.HasName("div"))

// Iterate through all nodes (text and tags)
for node := range div.IterNodes() {
	switch n := node.(type) {
	case *gosoup.Tag:
		fmt.Printf("Tag: %s (content: %s)\n", n.Name, n.Text())
	case gosoup.NavigableString:
		fmt.Printf("Text: %q\n", n.Text)
	}
}

// Output:
// Text: "Hello "
// Tag: b (content: World)
// Text: " and "
// Tag: i (content: Goodbye)
```

## Predicate System

To provide the flexibility similar to BeautifulSoup4, GoSoup uses a predicate system based on composable search functions. Predicates allow you to express complex selection criteria by combining simple, focused functions.

### Built-in Predicates

- **`HasName(name string) Predicate`** - Match by tag name
- **`HasAttr(attr string) Predicate`** - Check if an attribute exists
- **`HasNoAttr(attr string) Predicate`** - Check if an attribute does not exist
- **`HasClass(class string) Predicate`** - Check if element has a specific CSS class
- **`HasNoClass() Predicate`** - Check if element has no class attribute
- **`AttrEq(attr, value string) Predicate`** - Match attribute value exactly
- **`AttrContains(attr, substr string) Predicate`** - Match attribute value contains substring
- **`AttrMatch(attr string, pattern *regexp.Regexp) Predicate`** - Match attribute value against regex
- **`All(predicates ...Predicate) Predicate`** - Combine predicates with AND logic
- **`Any(predicates ...Predicate) Predicate`** - Combine predicates with OR logic

### Combining Predicates

Predicates can be combined for more complex queries:

```go
// Find all div tags with class "container"
divs := root.FindAll(gosoup.All(
	gosoup.HasName("div"),
	gosoup.HasClass("container"),
))

// Find links that are either in the nav or have id="main-link"
links := root.FindAll(gosoup.Any(
	gosoup.AttrEq("id", "main-link"),
	gosoup.HasClass("nav"),
))
```

### Custom Predicates

You can create your own predicates for specific use cases. A predicate is simply a function that takes a `*Tag` and returns a boolean:

```go
// Define a custom predicate to find external links
isExternalLink := func(tag *gosoup.Tag) bool {
	if tag.Name != "a" {
		return false
	}
	href, ok := tag.Attrs["href"]
	return ok && strings.HasPrefix(href, "http")
}

// Use the custom predicate
externalLinks := root.FindAll(isExternalLink)

// Combine custom predicates with built-in ones
links := root.FindAll(gosoup.All(
	isExternalLink,
	gosoup.HasClass("important"),
))
```

## Testing

Run the test suite with:

```bash
go test ./...
```

For coverage report:

```bash
go test -cover ./...
```

## Notes

GoSoup uses the `Parse` function from `golang.org/x/net/html` internally. Please note the following limitations:

- HTML that is nested deeper than 512 elements will be rejected
- The input is assumed to be UTF-8 encoded

## License

This library is open source and available under the MIT License.
