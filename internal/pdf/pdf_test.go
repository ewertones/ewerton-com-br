package pdf

import (
	"bytes"
	"compress/zlib"
	"io"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// build produces a small document exercising the features the resume relies on.
func build(t *testing.T) []byte {
	t.Helper()
	d := New()
	d.Title = "Test Document"
	d.Author = "Ewerton de Souza"
	d.Keywords = []string{"Go", "PDF"}

	p := d.NewPage(A4Width, A4Height)
	p.Rect(0, 0, A4Width, 100, RGB(0, 0x45, 0x87))
	p.Text(48, 60, "Olá, ação e coração — 5 · ½", Bold, 14, RGB(0, 0, 0))
	p.LinkText(48, 120, "ewerton.com.br", Regular, 10, RGB(0, 0, 0xFF), "https://ewerton.com.br")
	p.StarRating(48, 160, 3, 2, 4.5, 5, RGB(0, 0, 0x80), RGB(0xEE, 0xEE, 0xEE))
	p.GradientRect(48, 200, 200, 20, 0, 0, 1, 0, []Stop{
		{At: 0, Color: RGB(0, 0, 0)}, {At: 0.5, Color: RGB(0x80, 0x80, 0x80)}, {At: 1, Color: RGB(0xFF, 0xFF, 0xFF)},
	})
	p.Save()
	p.SetAlpha(0.2)
	p.Circle(300, 300, 40, RGB(0, 0, 0))
	p.Restore()

	d.Bookmark("Top", p, 0)
	d.Attach("resume.json", "application/json", "structured data", []byte(`{"ok":true}`))
	return d.Bytes()
}

func TestDocumentStructure(t *testing.T) {
	out := build(t)

	if !bytes.HasPrefix(out, []byte("%PDF-1.7")) {
		t.Error("missing PDF header")
	}
	if !bytes.HasSuffix(out, []byte("%%EOF\n")) {
		t.Error("missing EOF marker")
	}
	for _, want := range []string{
		"/Type /Catalog", "/Type /Pages", "/Type /Page ", "/Type /Metadata",
		"/Type /Outlines", "/EmbeddedFile", "/AFRelationship /Source",
		"/Lang (", "/ShadingType 2", "/ExtGState", "trailer", "startxref",
	} {
		if !bytes.Contains(out, []byte(want)) {
			t.Errorf("output is missing %q", want)
		}
	}
}

// The cross-reference table must point at the real byte offset of every object,
// or strict readers reject the file.
func TestCrossReferenceOffsetsAreCorrect(t *testing.T) {
	out := build(t)

	i := bytes.LastIndex(out, []byte("startxref"))
	if i < 0 {
		t.Fatal("no startxref")
	}
	start, err := strconv.Atoi(strings.Fields(string(out[i+len("startxref"):]))[0])
	if err != nil {
		t.Fatalf("parse startxref: %v", err)
	}
	if start <= 0 || start >= len(out) {
		t.Fatalf("startxref %d is outside the file", start)
	}
	if !bytes.HasPrefix(out[start:], []byte("xref")) {
		t.Fatalf("startxref does not point at the xref table")
	}

	entry := regexp.MustCompile(`(?m)^(\d{10}) \d{5} n $`)
	matches := entry.FindAllSubmatch(out[start:], -1)
	if len(matches) == 0 {
		t.Fatal("no xref entries")
	}
	for n, m := range matches {
		off, err := strconv.Atoi(string(m[1]))
		if err != nil {
			t.Fatal(err)
		}
		prefix := []byte(strconv.Itoa(n+1) + " 0 obj")
		if !bytes.HasPrefix(out[off:], prefix) {
			t.Errorf("object %d: offset %d does not start with %q", n+1, off, prefix)
		}
	}
}

// Text must survive a round trip through the content stream, because that is
// exactly what an ATS or a crawler extracts.
func TestTextIsExtractable(t *testing.T) {
	out := build(t)

	var found bool
	for _, raw := range regexp.MustCompile(`(?s)stream\n(.*?)\nendstream`).FindAllSubmatch(out, -1) {
		zr, err := zlib.NewReader(bytes.NewReader(raw[1]))
		if err != nil {
			continue
		}
		body, err := io.ReadAll(zr)
		if err != nil {
			continue
		}
		if bytes.Contains(body, []byte("ewerton.com.br")) {
			found = true
		}
		// Accented text is stored in WinAnsi, so decode before comparing.
		if bytes.Contains(body, []byte(ToWinAnsi("Olá, ação e coração"))) {
			found = found && true
		} else if bytes.Contains(body, []byte("Ol")) {
			t.Error("accented text was not encoded as WinAnsi")
		}
	}
	if !found {
		t.Error("could not extract text from any content stream")
	}
}

func TestToWinAnsiEncodesLatinAndPunctuation(t *testing.T) {
	cases := map[string][]byte{
		"a": {0x61},
		"ã": {0xE3},
		"ç": {0xE7},
		"õ": {0xF5},
		"é": {0xE9},
		"–": {0x96}, // en dash
		"·": {0xB7}, // middle dot
		"’": {0x92}, // right single quote
		"•": {0x95},
		"…": {0x85},
		"→": {0x2D, 0x3E}, // transliterated to ->
		"​": {},           // zero-width space is dropped
	}
	for in, want := range cases {
		if got := ToWinAnsi(in); got != string(want) {
			t.Errorf("ToWinAnsi(%q) = % x, want % x", in, got, want)
		}
	}
}

func TestTextWidthAndWrapping(t *testing.T) {
	if w := TextWidth("", Regular, 10); w != 0 {
		t.Errorf("empty string width = %v, want 0", w)
	}
	// Helvetica space is 278/1000 em.
	if got, want := TextWidth(" ", Regular, 1000), 278.0; got != want {
		t.Errorf("space width = %v, want %v", got, want)
	}
	// Bold is wider than regular for the same lowercase text.
	if TextWidth("mmm", Bold, 10) <= TextWidth("mmm", Regular, 10) {
		t.Error("bold should be wider than regular")
	}

	lines := WrapText("one two three four five six seven eight", Regular, 10, 60)
	if len(lines) < 2 {
		t.Fatalf("expected the text to wrap, got %#v", lines)
	}
	for _, line := range lines {
		if w := TextWidth(line, Regular, 10); w > 60 {
			t.Errorf("line %q is %v wide, over the 60 limit", line, w)
		}
	}
	if strings.Join(lines, " ") != "one two three four five six seven eight" {
		t.Errorf("wrapping lost or reordered words: %#v", lines)
	}

	// Hard newlines are honoured.
	if got := WrapText("a\nb", Regular, 10, 500); len(got) != 2 {
		t.Errorf("newline should break the line, got %#v", got)
	}
}

func TestTruncateFitsWidth(t *testing.T) {
	const width = 40.0
	got := Truncate("Google Ad Manager Extended", Regular, 8, width)
	if got == "" {
		t.Fatal("truncate returned nothing")
	}
	if w := TextWidth(got, Regular, 8); w > width {
		t.Errorf("%q is %v wide, over the %v limit", got, w, width)
	}
	if short := "SQL"; Truncate(short, Regular, 8, width) != short {
		t.Error("text that already fits should be returned unchanged")
	}
}

func TestLoadJPEGReadsFrameHeader(t *testing.T) {
	// Minimal baseline JPEG: SOI, SOF0 with 8x16 and 3 components, EOI.
	data := []byte{
		0xFF, 0xD8,
		0xFF, 0xC0, 0x00, 0x11, 0x08, 0x00, 0x10, 0x00, 0x08, 0x03,
		0x01, 0x11, 0x00, 0x02, 0x11, 0x01, 0x03, 0x11, 0x01,
		0xFF, 0xD9,
	}
	img, err := LoadJPEG(data)
	if err != nil {
		t.Fatalf("LoadJPEG: %v", err)
	}
	if img.Width != 8 || img.Height != 16 {
		t.Errorf("dimensions = %dx%d, want 8x16", img.Width, img.Height)
	}
	if img.colorSpace() != "DeviceRGB" {
		t.Errorf("colour space = %s, want DeviceRGB", img.colorSpace())
	}
	if _, err := LoadJPEG([]byte("not a jpeg")); err == nil {
		t.Error("expected an error for non-JPEG input")
	}
}
