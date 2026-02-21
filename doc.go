package gosoup

import (
	"bytes"
	"errors"
	"io"
	"strings"

	"golang.org/x/net/html"
)

// Corresponds to HTML document
type Document struct {
	root *html.Node
	cache map[*html.Node]*Tag
}

// Return root tag
func (doc *Document) Root() *Tag {
	return doc.newTag(doc.root)
}

// Parse HTML document from given reader and return root tag.
// Since Parse() from the golang.org/x/net/html library is used internally,
// the rules for basic Parse also apply for this function:
//
// * "Parse will reject HTML that is nested deeper than 512 elements."
//
// * "The input is assumed to be UTF-8 encoded."
func Parse(reader io.Reader) (*Document, error) {
	root, err := html.Parse(reader)
	if err != nil {
		return nil, err
	}

	return getDocument(root)
}

// Parse given HTML document bytes and return root tag.
// Since Parse() from the golang.org/x/net/html library is used internally,
// the rules for basic Parse also apply for this function:
//
// * "Parse will reject HTML that is nested deeper than 512 elements."
//
// * "The input is assumed to be UTF-8 encoded."
func ParseBytes(content []byte) (*Document, error) {
	root, err := html.Parse(bytes.NewReader(content))
	if err != nil {
		return nil, err
	}

	return getDocument(root)
}

// Parse given HTML document string and return root tag.
// Since Parse() from the golang.org/x/net/html library is used internally,
// the rules for basic Parse also apply for this function:
//
// * "Parse will reject HTML that is nested deeper than 512 elements."
//
// * "The input is assumed to be UTF-8 encoded."
func ParseString(content string) (*Document, error) {
	root, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return nil, err
	}

	return getDocument(root)
}

// Finding root element node (tag) of HTML document
func getDocument(root *html.Node) (*Document, error) {
	rootElement := findElementNode(root)
	if rootElement == nil {
		return nil, errors.New("no element node found")
	}

	doc := &Document{
		root: rootElement,
		cache: make(map[*html.Node]*Tag),
	}

	return doc, nil
}

// Finding first element node from given root
func findElementNode(node *html.Node) *html.Node {
	if node.Type == html.ElementNode {
		return node
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if found := findElementNode(c); found != nil {
			return found
		}
	}
	return nil
}

// Creates new Tag structure from a given node
func (doc *Document) newTag(node *html.Node) *Tag {
	if node == nil {
		return nil
	}

	if node.Type != html.ElementNode {
		return nil
	}

	if tag, ok := doc.cache[node]; ok {
		return tag
	}

	attrs := make(map[string]string, len(node.Attr))
	for _, attr := range node.Attr {
		attrs[attr.Key] = attr.Val
	}

	tag := &Tag{
		Name: node.Data, 
		Attrs: attrs, 
		node: node, 
		doc: doc,
	}
	doc.cache[node] = tag

	return tag
}
