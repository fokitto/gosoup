package gosoup

import (
	"iter"
	"strings"

	"golang.org/x/net/html"
)

// Corresponds to HTML node in the document (tag, raw string, etc.)
type Node interface {
	isNode()
}

// Raw string node in HTML document
type NavigableString struct {
	Text string
}

func (ns NavigableString) isNode() {}

// Corresponds to HTML tag in the document
type Tag struct {
	Name  string
	Attrs map[string]string
	node  *html.Node
	doc *Document
}

func (t *Tag) isNode() {}

// Render a tree with a current tag as root
func (tag *Tag) String() string {
	var builder strings.Builder
	html.Render(&builder, tag.node)
	return builder.String()
}

// Get a parent tag
func (tag *Tag) Parent() *Tag {
	for parent := tag.node.Parent; parent != nil; parent = parent.Parent {
		if parent.Type == html.ElementNode {
			return tag.doc.newTag(parent)
		}
	}
	return nil
}

// Get a first child of a tag
func (tag *Tag) FirstChild() *Tag {
	for child := tag.node.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode {
			return tag.doc.newTag(child)
		}
	}
	return nil
}

// Get all children tags recursively
func (tag *Tag) Children() []*Tag {
	var children []*Tag

	for node := tag.node.FirstChild; node != nil; node = node.NextSibling {
		if node.Type != html.ElementNode {
			continue
		}
		children = append(children, tag.doc.newTag(node))
	}

	return children
}

// Returns count of all inner tags
func (tag *Tag) ChildrenCount() int {
	cnt := 0

	for node := tag.node.FirstChild; node != nil; node = node.NextSibling {
		if node.Type != html.ElementNode {
			continue
		}
		cnt++
	}

	return cnt
}

// Returns depth of current tag
func (tag *Tag) Depth() int {
	depth := 0

	for parent := tag.node.Parent; parent != nil; parent = parent.Parent {
		if parent.Type != html.ElementNode {
			continue
		}
		depth += 1
	}

	return depth
}

// Get previous sibling of tag
func (tag *Tag) Prev() *Tag {
	for prev := tag.node.PrevSibling; prev != nil; prev = prev.PrevSibling {
		if prev.Type == html.ElementNode {
			return tag.doc.newTag(prev)
		}
	}
	return nil
}

// Get next sibling of tag
func (tag *Tag) Next() *Tag {
	for next := tag.node.NextSibling; next != nil; next = next.NextSibling {
		if next.Type == html.ElementNode {
			return tag.doc.newTag(next)
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
	builder := &strings.Builder{}

	var traverse func(*html.Node)
	traverse = func(node *html.Node) {
		if node == nil {
			return
		}
		if node.Type == html.TextNode {
			builder.WriteString(node.Data)
			for _, s := range sep {
				builder.WriteString(s)
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			traverse(child)
		}
	}

	traverse(tag.node)

	return builder.String()
}

// Removes current tag from a tree
func (tag *Tag) Unwrap() {
	tag.doc.removeTag(tag)
}

// Find chidl tag by predicate
func (tag *Tag) Find(predicate Predicate) *Tag {
	var find func(*Tag, bool) *Tag

	find = func(t *Tag, skipCheck bool) *Tag {
		if !skipCheck && predicate(t) {
			return t
		}

		for child := t.FirstChild(); child != nil; child = child.Next() {
			if found := find(child, false); found != nil {
				return found
			}
		}

		return nil
	}


	return find(tag, true)
}

// Find all children tags by predicate
func (tag *Tag) FindAll(predicate Predicate) []*Tag {
	var result []*Tag

	var find func(*Tag, bool)
	find = func(t *Tag, skipCheck bool) {
		if !skipCheck && predicate(t) {
			result = append(result, t)
		}

		for child := t.FirstChild(); child != nil; child = child.Next() {
			find(child, false)
		}
	}

	find(tag, true)

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

// Iterate through all children nodes of current tag,
// including raw strings
func (tag *Tag) IterNodes() iter.Seq[Node] {
	return func(yield func(Node) bool) {
		var traverse func(*html.Node) bool
		traverse = func(n *html.Node) bool {
			for child := n.FirstChild; child != nil; child = child.NextSibling {

				switch child.Type{
				case html.ElementNode:
					node := tag.doc.newTag(child)
					if !yield(node) {
						return false
					}
				case html.TextNode:
					node := NavigableString{Text: child.Data}
					if !yield(node) {
						return false
					}
					if !traverse(child) {
						return false
					}
				}

			}
			return true
		}
		
		traverse(tag.node)
	}
}
