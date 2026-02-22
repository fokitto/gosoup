package gosoup

import (
	"strings"
	"testing"
)

func TestParseString(t *testing.T) {
	doc, err := ParseString(sampleHTML)
	if err != nil {
		t.Fatalf("ParseString error: %v", err)
	}

    root := doc.Root()

	if root == nil {
		t.Fatalf("ParseString returned nil root")
	}
	if root.Name != "html" {
		t.Fatalf("expected root tag to be 'html', got %q", root.Name)
	}
}

func TestParseBytes(t *testing.T) {
	doc, err := ParseBytes([]byte(sampleHTML))
	if err != nil {
		t.Fatalf("ParseBytes error: %v", err)
	}

    root := doc.Root()

	if root == nil {
		t.Fatalf("ParseBytes returned nil root")
	}
	if root.Name != "html" {
		t.Fatalf("expected root tag to be 'html', got %q", root.Name)
	}
}

func TestParse(t *testing.T) {
	doc, err := Parse(strings.NewReader(sampleHTML))
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

    root := doc.Root()

	if root == nil {
		t.Fatalf("Parse returned nil root")
	}
	if root.Name != "html" {
		t.Fatalf("expected root tag to be 'html', got %q", root.Name)
	}
}

func TestDocumentRoot(t *testing.T) {
	html := `
	<html>
		<head><title>Test</title></head>
		<body>
			<div class="content">
				<p>Hello</p>
			</div>
		</body>
	</html>
	`

	doc, err := ParseString(html)
	if err != nil {
		t.Fatalf("ParseString error: %v", err)
	}

	root := doc.Root()
	if root == nil {
		t.Fatalf("expected root to be non-nil")
	}
	if root.Name != "html" {
		t.Fatalf("expected root tag name to be 'html', got %q", root.Name)
	}
}

func TestDocumentCaching(t *testing.T) {
	html := `
	<html>
		<body>
			<div id="test">Content</div>
		</body>
	</html>
	`

	doc, err := ParseString(html)
	if err != nil {
		t.Fatalf("ParseString error: %v", err)
	}

	root := doc.Root()

	div1 := root.Find(HasName("div"))
	div2 := root.Find(HasName("div"))

	if div1 != div2 {
		t.Errorf("expected same tag reference from cache, got different pointers")
	}

	if div1.Name != "div" {
		t.Fatalf("expected tag name to be 'div', got %q", div1.Name)
	}
	if div1.Attrs["id"] != "test" {
		t.Fatalf("expected id attribute to be 'test', got %q", div1.Attrs["id"])
	}
}

func TestDocumentCachingMultiplePaths(t *testing.T) {
	html := `
	<html>
		<body>
			<div class="container">
				<p id="para">Text</p>
			</div>
		</body>
	</html>
	`

	doc, err := ParseString(html)
	if err != nil {
		t.Fatalf("ParseString error: %v", err)
	}

	root := doc.Root()
	body := root.Find(HasName("body"))
	div := body.Find(HasName("div"))

	para1 := div.Find(HasName("p"))
	para2 := div.FirstChild()

	if para1 != para2 {
		t.Errorf("expected paragraph to be cached, got different pointers from different paths")
	}

	if para1.Attrs["id"] != "para" {
		t.Fatalf("expected paragraph id to be 'para', got %q", para1.Attrs["id"])
	}
}

func TestDocumentSeparateCaches(t *testing.T) {
	html := `
	<html>
		<body>
			<div id="test">Content</div>
		</body>
	</html>
	`

	doc1, err1 := ParseString(html)
	if err1 != nil {
		t.Fatalf("ParseString error: %v", err1)
	}

	doc2, err2 := ParseString(html)
	if err2 != nil {
		t.Fatalf("ParseString error: %v", err2)
	}

	root1 := doc1.Root()
	root2 := doc2.Root()

	div1 := root1.Find(HasName("div"))
	div2 := root2.Find(HasName("div"))

	if div1 == div2 {
		t.Errorf("expected tags from different documents to be different objects")
	}

	if div1.Attrs["id"] != div2.Attrs["id"] {
		t.Fatalf("expected same attribute values, got %q and %q", div1.Attrs["id"], div2.Attrs["id"])
	}
}

func TestDocumentNullElements(t *testing.T) {
	html := `
	<html>
		<body>
			<div>Content</div>
		</body>
	</html>
	`

	doc, err := ParseString(html)
	if err != nil {
		t.Fatalf("ParseString error: %v", err)
	}

	root := doc.Root()
	div := root.Find(HasName("div"))

	nonexistent := div.Find(HasName("article"))
	if nonexistent != nil {
		t.Fatalf("expected nil for non-existent element, got %v", nonexistent)
	}

	parent := root.Parent()
	if parent != nil {
		t.Fatalf("expected nil when getting parent of document root")
	}
}

func TestChildrenCount(t *testing.T) {
	html := `
	<html>
		<body>
			<div>
				<p>Paragraph 1</p>
				<p>Paragraph 2</p>
				<span>Span</span>
			</div>
		</body>
	</html>
	`

	doc, err := ParseString(html)
	if err != nil {
		t.Fatalf("ParseString error: %v", err)
	}

	root := doc.Root()
	div := root.Find(HasName("div"))

	count := div.ChildrenCount()
	if count != 3 {
		t.Fatalf("expected 3 children, got %d", count)
	}

	p := root.Find(HasName("p"))
	if p.ChildrenCount() != 0 {
		t.Fatalf("expected 0 children for paragraph, got %d", p.ChildrenCount())
	}
}

func TestDepth(t *testing.T) {
	html := `
	<html>
		<body>
			<div>
				<p>
					<strong>Bold text</strong>
				</p>
			</div>
		</body>
	</html>
	`

	doc, err := ParseString(html)
	if err != nil {
		t.Fatalf("ParseString error: %v", err)
	}

	root := doc.Root()

	if root.Depth() != 0 {
		t.Fatalf("expected root depth to be 0, got %d", root.Depth())
	}

	body := root.Find(HasName("body"))
	if body.Depth() != 1 {
		t.Fatalf("expected body depth to be 1, got %d", body.Depth())
	}

	div := root.Find(HasName("div"))
	if div.Depth() != 2 {
		t.Fatalf("expected div depth to be 2, got %d", div.Depth())
	}

	p := root.Find(HasName("p"))
	if p.Depth() != 3 {
		t.Fatalf("expected p depth to be 3, got %d", p.Depth())
	}

	strong := root.Find(HasName("strong"))
	if strong.Depth() != 4 {
		t.Fatalf("expected strong depth to be 4, got %d", strong.Depth())
	}
}
