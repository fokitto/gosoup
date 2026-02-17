package gosoup

import (
	"bytes"
	"errors"
	"io"
	"strings"

	"golang.org/x/net/html"
)

// Corresponds to HTML tag in the document
type Tag struct {
	Name  string
	Attrs map[string]string
	node  *html.Node
}

// Render a tree with a current tag as root
func (tag *Tag) String() string {
	var buff bytes.Buffer
	html.Render(&buff, tag.node)
	return buff.String()
}

// Get a parent tag
func (tag *Tag) Parent() *Tag {
	return newTag(tag.node.Parent)
}

// Get a first child of a tag
func (tag *Tag) FirstChild() *Tag {
	return newTag(tag.node.FirstChild)
}

// Get all children tags recursively 
func (tag *Tag) Children() []*Tag {
	var children []*Tag

	for node := tag.node.FirstChild; node != nil; node = node.NextSibling {
		if node.Type != html.ElementNode {
			continue
		}
		children = append(children, newTag(node))
	}

	return children
}

// Get previous sibling of tag
func (tag *Tag) Prev() *Tag {
	for prev := tag.node.PrevSibling; prev != nil; prev = prev.PrevSibling {
		if prev.Type == html.ElementNode {
			return newTag(prev)
		}
	}
	return nil
}

// Get next sibling of tag
func (tag *Tag) Next() *Tag {
	for next := tag.node.NextSibling; next != nil; next = next.NextSibling {
		if next.Type == html.ElementNode {
			return newTag(next)
		}
	}
	return nil
}

// Get inner text of tag, without traversing inner tags
func (tag *Tag) Text() string {
	for node := tag.node.FirstChild; node != nil; node = node.NextSibling {
		if node.Type == html.TextNode {
			return node.Data
		}
	}
	return ""
}

// Get all human-readable text of a current tree
func (tag *Tag) FullText(sep ...string) string {
	var buff bytes.Buffer

	var traverse func(*html.Node)
	traverse = func(node *html.Node) {
		if node == nil {
			return
		}
		if node.Type == html.TextNode {
			buff.WriteString(node.Data)
			for _, s := range sep {
				buff.WriteString(s)
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			traverse(child)
		}
	}

	traverse(tag.node)

	return buff.String()
}

// Removes current tag from a tree
func (tag *Tag) Unwrap() {
	tag.node.Parent.RemoveChild(tag.node)
}

// Find chidl tag by predicate
func (tag *Tag) Find(predicate Predicate) *Tag {
	var find func(*Tag) *Tag

	find = func(t *Tag) *Tag {
		if predicate(t) {
			return t
		}

		for child := t.FirstChild(); child != nil; child = child.Next() {
			if found := find(child); found != nil {
				return found
			}
		}

		return nil
	}

	dummy := createDummy(tag.node)

	return find(dummy)
}

// Find all children tags by predicate
func (tag *Tag) FindAll(predicate Predicate) []*Tag {
	var result []*Tag

	var find func(*Tag)
	find = func(t *Tag) {
		if predicate(t) {
			result = append(result, t)
		}

		for child := t.FirstChild(); child != nil; child = child.Next() {
			find(child)
		}
	}

	dummy := createDummy(tag.node)
	find(dummy)

	return result
}

// Find parent tag by predicate
func (tag *Tag) FindParent(predicate Predicate) *Tag {
	var find func(*Tag) *Tag

	find = func(t *Tag) *Tag {
		if t == nil {
			return nil
		}
		if predicate(t) {
			return t
		}

		for parent := t.Parent(); parent != nil; parent = parent.Parent() {
			if found := find(parent); found != nil {
				return found
			}
		}

		return nil
	}

	return find(tag.Parent())
}

// Parse HTML document from given reader and return root tag.
// Since Parse() from the golang.org/x/net/html library is used internally,
// the rules for basic Parse also apply for this function:
//
// * "Parse will reject HTML that is nested deeper than 512 elements."
//
// * "The input is assumed to be UTF-8 encoded."
func Parse(reader io.Reader) (*Tag, error) {
	root, err := html.Parse(reader)
	if err != nil {
		return nil, err
	}

	return findRootTag(root)
}

// Parse given HTML document bytes and return root tag.
// Since Parse() from the golang.org/x/net/html library is used internally,
// the rules for basic Parse also apply for this function:
//
// * "Parse will reject HTML that is nested deeper than 512 elements."
//
// * "The input is assumed to be UTF-8 encoded."
func ParseBytes(content []byte) (*Tag, error) {
	root, err := html.Parse(bytes.NewReader(content))
	if err != nil {
		return nil, err
	}

	return findRootTag(root)
}

// Parse given HTML document string and return root tag.
// Since Parse() from the golang.org/x/net/html library is used internally,
// the rules for basic Parse also apply for this function:
//
// * "Parse will reject HTML that is nested deeper than 512 elements."
//
// * "The input is assumed to be UTF-8 encoded."
func ParseString(content string) (*Tag, error) {
	root, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return nil, err
	}

	return findRootTag(root)
}

// Finding root element node (tag) of HTML document
func findRootTag(root *html.Node) (*Tag, error) {
	rootElement := findElementNode(root)
	if rootElement == nil {
		return nil, errors.New("no element node found")
	}

	return newTag(rootElement), nil
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
func newTag(node *html.Node) *Tag {
	if node == nil {
		return nil
	}

	attrs := make(map[string]string, len(node.Attr))
	for _, attr := range node.Attr {
		attrs[attr.Key] = attr.Val
	}

	return &Tag{node.Data, attrs, node}
}

// Сreates a stub without a name or attributes,
// but with the given pointer to the tree-node
func createDummy(node *html.Node) *Tag {
	return &Tag{
		Name: "",
		Attrs: map[string]string{},
		node: node,
	}
}
