package gosoup

import (
	"strings"
	"testing"
)

var sampleHTML = `<!doctype html>
<html>
  <head></head>
  <body>
    <div id="root" class="container">
      <p class="a b">Hello <span>World</span></p>
      <p class="b">Second</p>
      <article>
        <h1>Title</h1>
        <p>Content</p>
      </article>
    </div>
  </body>
</html>`

func TestParent(t *testing.T) {
	doc, err := ParseString(sampleHTML)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

    root := doc.Root()

	body := root.Find(HasName("body"))
	if body == nil {
		t.Fatalf("could not find body")
	}

	parent := body.Parent()
	if parent == nil {
		t.Fatalf("Parent() returned nil")
	}
	if parent.Name != "html" {
		t.Fatalf("expected parent to be 'html', got %q", parent.Name)
	}
}

func TestFirstChild(t *testing.T) {
	doc, err := ParseString(sampleHTML)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

    root := doc.Root()

	firstChild := root.FirstChild()
	if firstChild == nil {
		t.Fatalf("FirstChild() returned nil")
	}
	if firstChild.Name != "head" {
		t.Fatalf("expected first child to be 'head', got %q", firstChild.Name)
	}
}

func TestChildren(t *testing.T) {
	doc, err := ParseString(sampleHTML)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

    root := doc.Root()

	children := root.Children()
	if len(children) != 2 {
		t.Fatalf("expected 2 children for html, got %d", len(children))
	}
	if children[0].Name != "head" {
		t.Fatalf("expected first child to be 'head', got %q", children[0].Name)
	}
	if children[1].Name != "body" {
		t.Fatalf("expected second child to be 'body', got %q", children[1].Name)
	}
}

func TestNext(t *testing.T) {
	doc, err := ParseString(sampleHTML)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

    root := doc.Root()

	head := root.FirstChild()
	if head == nil {
		t.Fatalf("FirstChild() returned nil")
	}

	next := head.Next()
	if next == nil {
		t.Fatalf("Next() returned nil")
	}
	if next.Name != "body" {
		t.Fatalf("expected next to be 'body', got %q", next.Name)
	}
}

func TestPrev(t *testing.T) {
	doc, err := ParseString(sampleHTML)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

    root := doc.Root()

	body := root.Find(HasName("body"))
	if body == nil {
		t.Fatalf("could not find body")
	}

	prev := body.Prev()
	if prev == nil {
		t.Fatalf("Prev() returned nil")
	}
	if prev.Name != "head" {
		t.Fatalf("expected prev to be 'head', got %q", prev.Name)
	}
}

func TestText(t *testing.T) {
	doc, err := ParseString(sampleHTML)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

    root := doc.Root()

	p := root.Find(HasClass("a"))
	if p == nil {
		t.Fatalf("could not find paragraph with class 'a'")
	}

	text := p.Text()
	if strings.TrimSpace(text) != "Hello" {
		t.Fatalf("expected text 'Hello', got %q", text)
	}
}

func TestFullText(t *testing.T) {
	doc, err := ParseString(sampleHTML)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

    root := doc.Root()

	p := root.Find(HasClass("a"))
	if p == nil {
		t.Fatalf("could not find paragraph with class 'a'")
	}

	fullText := p.FullText()
	// Should contain both "Hello" and "World"
	if !strings.Contains(fullText, "Hello") || !strings.Contains(fullText, "World") {
		t.Fatalf("expected 'Hello' and 'World' in fullText, got %q", fullText)
	}
}

func TestFullTextWithSeparator(t *testing.T) {
	doc, err := ParseString(sampleHTML)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

    root := doc.Root()

	p := root.Find(HasClass("a"))
	if p == nil {
		t.Fatalf("could not find paragraph with class 'a'")
	}

	fullText := p.FullText(" | ")
	// Should contain separator between text nodes
	if !strings.Contains(fullText, " | ") {
		t.Fatalf("expected separator ' | ' in fullText, got %q", fullText)
	}
}

func TestString(t *testing.T) {
	doc, err := ParseString(sampleHTML)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

    root := doc.Root()

	p := root.Find(AttrEq("class", "a b"))
	if p == nil {
		t.Fatalf("could not find span")
	}

    text := p.String()
    if text != `<p class="a b">Hello <span>World</span></p>` {
        t.Fatalf("Render failed: got: %s", text)
    }
}

func TestFind(t *testing.T) {
	doc, err := ParseString(sampleHTML)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

    root := doc.Root()

	found := root.Find(HasName("span"))
	if found == nil {
		t.Fatalf("Find() returned nil for existing span tag")
	}
	if found.Name != "span" {
		t.Fatalf("expected 'span', got %q", found.Name)
	}
}

func TestFindNotFound(t *testing.T) {
	doc, err := ParseString(sampleHTML)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

    root := doc.Root()

	found := root.Find(HasName("video"))
	if found != nil {
		t.Fatalf("Find() should return nil for non-existent tag, got %v", found)
	}
}

func TestFindByAttribute(t *testing.T) {
	doc, err := ParseString(sampleHTML)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

    root := doc.Root()

	found := root.Find(AttrEq("id", "root"))
	if found == nil {
		t.Fatalf("Find() returned nil for div#root")
	}
	if found.Name != "div" {
		t.Fatalf("expected 'div', got %q", found.Name)
	}
}

func TestFindAll(t *testing.T) {
	doc, err := ParseString(sampleHTML)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

    root := doc.Root()

	found := root.FindAll(HasName("p"))
	if len(found) != 3 {
		t.Fatalf("expected 3 paragraphs, found %d", len(found))
	}
	for i, p := range found {
		if p.Name != "p" {
			t.Fatalf("expected 'p' at index %d, got %q", i, p.Name)
		}
	}
}

func TestFindAllWithPredicate(t *testing.T) {
	doc, err := ParseString(sampleHTML)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

    root := doc.Root()

	found := root.FindAll(HasClass("b"))
	if len(found) != 2 {
		t.Fatalf("expected 2 elements with class 'b', found %d", len(found))
	}
}

func TestFindParent(t *testing.T) {
	doc, err := ParseString(sampleHTML)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

    root := doc.Root()

	span := root.Find(HasName("span"))
	if span == nil {
		t.Fatalf("could not find span")
	}

	parent := span.FindParent(HasName("p"))
	if parent == nil {
		t.Fatalf("FindParent() returned nil")
	}
	if parent.Name != "p" {
		t.Fatalf("expected parent to be 'p', got %q", parent.Name)
	}
}

func TestFindParentNotFound(t *testing.T) {
	doc, err := ParseString(sampleHTML)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

    root := doc.Root()

	span := root.Find(HasName("span"))
	if span == nil {
		t.Fatalf("could not find span")
	}

	parent := span.FindParent(HasName("video"))
	if parent != nil {
		t.Fatalf("FindParent() should return nil for non-existent ancestor, got %v", parent)
	}
}

func TestUnwrap(t *testing.T) {
	doc, err := ParseString(sampleHTML)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

    root := doc.Root()

	span := root.Find(HasName("span"))
	if span == nil {
		t.Fatalf("could not find span")
	}

	spanBefore := root.Find(HasName("span"))
	if spanBefore == nil {
		t.Fatalf("span should exist before unwrap")
	}

	span.Unwrap()

	spanAfter := root.Find(HasName("span"))
	if spanAfter != nil {
		t.Fatalf("span should not exist after unwrap")
	}
}

func TestIterNodes(t *testing.T) {
    doc, err := ParseString(`<div>Text with <a>inner</a> tag</div>`)
    if err != nil {
        t.Fatalf("Parse error: %v", err)
    }

    root := doc.Root()

    div := root.Find(HasName("div"))
	if div == nil {
		t.Fatalf("could not find div")
	}

    tagIndex := 0
    exptectedTexts := []string{"Text with ", "inner", " tag"}
    for node := range div.IterNodes() {
        if tagIndex > 2 {
            t.Fatal("more iterations than expected")
        }
        expected := exptectedTexts[tagIndex]
        switch tagIndex {
        case 0, 2:
            n, ok := node.(NavigableString)
            if !ok {
                t.Fatalf("node[%d] is not NavigableString", tagIndex)
            }
            if n.Text != expected {
                t.Fatalf("expected '%s', got '%s'", expected, n.Text)
            }
        case 1:
            n, ok := node.(*Tag)
            if !ok {
                t.Fatalf("second node is not *Tag")
            }
            if n.Text() != expected {
                t.Fatalf("expected '%s', got '%s'", expected, n.Text())
            }
        }
        tagIndex++
    }
}

func TestIterNodesWithOnlyText(t *testing.T) {
	html := `<p>Just plain text</p>`
	doc, err := ParseString(html)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

    root := doc.Root()

	p := root.Find(HasName("p"))
	if p == nil {
		t.Fatalf("could not find p")
	}

	var count int
	var foundText string
	for node := range p.IterNodes() {
		count++
		str, ok := node.(NavigableString)
		if !ok {
			t.Fatalf("expected NavigableString, got %T", node)
		}
		foundText = str.Text
	}

	if count != 1 {
		t.Fatalf("expected 1 node, got %d", count)
	}
	if strings.TrimSpace(foundText) != "Just plain text" {
		t.Fatalf("expected 'Just plain text', got %q", foundText)
	}
}

func TestIterNodesNestedTags(t *testing.T) {
	html := `<div><span>one</span><em>two</em><strong>three</strong></div>`
	doc, err := ParseString(html)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

    root := doc.Root()

	div := root.Find(HasName("div"))
	if div == nil {
		t.Fatalf("could not find div")
	}

	expectedTags := []string{"span", "em", "strong"}
	tagIndex := 0

	for node := range div.IterNodes() {
		tag, ok := node.(*Tag)
		if !ok {
			t.Fatalf("expected *Tag, got %T", node)
		}
		if tagIndex >= len(expectedTags) {
			t.Fatalf("too many tags: expected %d, got at least %d", len(expectedTags), tagIndex+1)
		}
		if tag.Name != expectedTags[tagIndex] {
			t.Fatalf("expected tag %q at index %d, got %q", expectedTags[tagIndex], tagIndex, tag.Name)
		}
		tagIndex++
	}

	if tagIndex != len(expectedTags) {
		t.Fatalf("expected %d tags, got %d", len(expectedTags), tagIndex)
	}
}
