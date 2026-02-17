package gosoup

import (
    "regexp"
    "testing"
)

func TestHasName(t *testing.T) {
    tag := &Tag{Name: "div"}
    if !HasName("div")(tag) {
        t.Fatalf("HasName failed")
    }
}

func TestHasAttr(t *testing.T) {
    tag := &Tag{Attrs: map[string]string{"id": "root"}}
    if !HasAttr("id")(tag) {
        t.Fatalf("HasAttr failed")
    }
}

func TestHasNoAttr(t *testing.T) {
    tag := &Tag{Attrs: map[string]string{"id": "root"}}
    if !HasNoAttr("data-qa")(tag) {
        t.Fatalf("HasNoAttr failed")
    }
}

func TestHasClass(t *testing.T) {
    tag := &Tag{Attrs: map[string]string{"class": "ab foo1"}}
    if !HasClass("ab")(tag) {
        t.Fatalf("HasClass failed")
    }
    if HasClass("foo")(tag) {
        t.Fatalf("HasClass failed: false positive")
    }
}

func TestHasNoClass(t *testing.T) {
    tagNoClass := &Tag{Attrs: map[string]string{}}
    tagWithClass := &Tag{Attrs: map[string]string{"class": "a"}}
    if !HasNoClass()(tagNoClass) {
        t.Fatalf("HasNoClass failed")
    }
    if HasNoClass()(tagWithClass) {
        t.Fatalf("HasNoClass failed: false positive")
    }
}

func TestAttrEq(t *testing.T) {
    tag := &Tag{Attrs: map[string]string{"class": "a b", "id": "root"}}
    if !AttrEq("id", "root")(tag) {
        t.Fatalf("AttrEq failed")
    }
}

func TestAttrContains(t *testing.T) {
    tag := &Tag{Attrs: map[string]string{"class": "ab b", "id": "root"}}
    if !AttrContains("class", "a")(tag) {
        t.Fatalf("AttrContains failed")
    }
}

func TestAttrMatch(t *testing.T) {
    tag := &Tag{Attrs: map[string]string{"class": "foo123", "id": "root"}}
    re := regexp.MustCompile(`foo\d+`)
    if !AttrMatch("class", re)(tag) {
        t.Fatalf("AttrMatch failed")
    }
}

func TestAll(t *testing.T) {
    tag := &Tag{Name: "div", Attrs: map[string]string{"id": "root"}}
    if !All(HasName("div"), HasAttr("id"))(tag) {
        t.Fatalf("All failed")
    }
}

func TestAny(t *testing.T) {
    tag := &Tag{Name: "div", Attrs: map[string]string{"id": "root"}}
    if !Any(HasName("span"), HasAttr("id"))(tag) {
        t.Fatalf("Any failed")
    }
}