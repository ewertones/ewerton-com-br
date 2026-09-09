package pdf

import (
	"bytes"
	"fmt"
	"strings"
)

// Page is a single page and its content stream.
//
// All coordinates taken by Page methods use a top-left origin measured in
// points, because that is how layout code naturally reads. Conversion to the
// PDF bottom-left origin happens here, in one place.
type Page struct {
	doc  *Doc
	ref  Ref
	W, H float64

	buf          bytes.Buffer
	usedFonts    map[Font]bool
	usedImages   []pageImage
	usedShadings []namedRef
	usedGStates  []namedRef
	annots       []string

	// Logical structure: the role applied to text, and the marked-content
	// sequences recorded so far. See tag.go.
	role  string
	marks []mark
}

// pageImage binds an image to the resource name this particular page uses for
// it. Two pages may name the same image differently, and each page's content
// stream must agree with its own resource dictionary.
type pageImage struct {
	name string
	img  *Image
}

// A4 dimensions in points.
const (
	A4Width  = 595.28
	A4Height = 841.89
)

// NewPage appends a page of the given size and returns it.
func (d *Doc) NewPage(w, h float64) *Page {
	p := &Page{doc: d, ref: d.reserve(), W: w, H: h, usedFonts: map[Font]bool{}}
	d.pages = append(d.pages, p)
	return p
}

// Pages returns the pages added so far.
func (d *Doc) Pages() []*Page { return d.pages }

// y converts a top-left Y coordinate into PDF user space.
func (p *Page) y(v float64) float64 { return p.H - v }

func (p *Page) op(format string, args ...any) {
	fmt.Fprintf(&p.buf, format, args...)
	p.buf.WriteByte('\n')
}

// ---------------------------------------------------------------------------
// Graphics state
// ---------------------------------------------------------------------------

// Color is an RGB color with components in the range [0,1].
type Color struct{ R, G, B float64 }

// RGB builds a Color from 8-bit components.
func RGB(r, g, b uint8) Color {
	return Color{float64(r) / 255, float64(g) / 255, float64(b) / 255}
}

// Mix blends a towards b by t in [0,1].
func Mix(a, b Color, t float64) Color {
	return Color{a.R + (b.R-a.R)*t, a.G + (b.G-a.G)*t, a.B + (b.B-a.B)*t}
}

// Save pushes the graphics state.
func (p *Page) Save() { p.op("q") }

// Restore pops the graphics state.
func (p *Page) Restore() { p.op("Q") }

// SetFill sets the non-stroking color.
func (p *Page) SetFill(c Color) { p.op("%s %s %s rg", num(c.R), num(c.G), num(c.B)) }

// SetStroke sets the stroking color.
func (p *Page) SetStroke(c Color) { p.op("%s %s %s RG", num(c.R), num(c.G), num(c.B)) }

// SetLineWidth sets the stroke width in points.
func (p *Page) SetLineWidth(w float64) { p.op("%s w", num(w)) }

// ---------------------------------------------------------------------------
// Shapes
// ---------------------------------------------------------------------------

// Rect fills a rectangle whose top-left corner is (x, y).
func (p *Page) Rect(x, y, w, h float64, c Color) {
	p.artifact(func() {
		p.SetFill(c)
		p.op("%s %s %s %s re f", num(x), num(p.y(y+h)), num(w), num(h))
	})
}

// Line strokes a straight line between two top-left-origin points.
func (p *Page) Line(x1, y1, x2, y2, width float64, c Color) {
	p.artifact(func() {
		p.SetStroke(c)
		p.SetLineWidth(width)
		p.op("%s %s m %s %s l S", num(x1), num(p.y(y1)), num(x2), num(p.y(y2)))
	})
}

// kappa is the Bezier constant that approximates a quarter circle.
const kappa = 0.5522847498307936

// circlePath emits a closed circular path without painting it.
func (p *Page) circlePath(cx, cy, r float64) {
	y := p.y(cy)
	k := r * kappa
	p.op("%s %s m", num(cx+r), num(y))
	p.op("%s %s %s %s %s %s c", num(cx+r), num(y+k), num(cx+k), num(y+r), num(cx), num(y+r))
	p.op("%s %s %s %s %s %s c", num(cx-k), num(y+r), num(cx-r), num(y+k), num(cx-r), num(y))
	p.op("%s %s %s %s %s %s c", num(cx-r), num(y-k), num(cx-k), num(y-r), num(cx), num(y-r))
	p.op("%s %s %s %s %s %s c", num(cx+k), num(y-r), num(cx+r), num(y-k), num(cx+r), num(y))
	p.op("h")
}

// Circle fills a circle centred on (cx, cy).
func (p *Page) Circle(cx, cy, r float64, c Color) {
	p.artifact(func() { p.fillCircle(cx, cy, r, c) })
}

// fillCircle paints a circle without opening a marked-content sequence, for
// callers that are already inside one.
func (p *Page) fillCircle(cx, cy, r float64, c Color) {
	p.SetFill(c)
	p.circlePath(cx, cy, r)
	p.op("f")
}

// StrokeCircle outlines a circle centred on (cx, cy).
func (p *Page) StrokeCircle(cx, cy, r, width float64, c Color) {
	p.artifact(func() {
		p.SetStroke(c)
		p.SetLineWidth(width)
		p.circlePath(cx, cy, r)
		p.op("S")
	})
}

// ClipCircle intersects the clipping path with a circle. Wrap it in
// Save/Restore, since clipping cannot be undone any other way.
func (p *Page) ClipCircle(cx, cy, r float64) {
	p.circlePath(cx, cy, r)
	p.op("W n")
}

// RoundRect fills a rectangle with rounded corners; used for chips and pills.
func (p *Page) RoundRect(x, y, w, h, r float64, c Color) {
	if r > h/2 {
		r = h / 2
	}
	if r > w/2 {
		r = w / 2
	}
	p.op("/Artifact BMC")
	defer p.op("EMC")
	p.SetFill(c)
	top, bot := p.y(y), p.y(y+h)
	k := r * kappa
	p.op("%s %s m", num(x+r), num(top))
	p.op("%s %s l", num(x+w-r), num(top))
	p.op("%s %s %s %s %s %s c", num(x+w-r+k), num(top), num(x+w), num(top-r+k), num(x+w), num(top-r))
	p.op("%s %s l", num(x+w), num(bot+r))
	p.op("%s %s %s %s %s %s c", num(x+w), num(bot+r-k), num(x+w-r+k), num(bot), num(x+w-r), num(bot))
	p.op("%s %s l", num(x+r), num(bot))
	p.op("%s %s %s %s %s %s c", num(x+r-k), num(bot), num(x), num(bot+r-k), num(x), num(bot+r))
	p.op("%s %s l", num(x), num(top-r))
	p.op("%s %s %s %s %s %s c", num(x), num(top-r+k), num(x+r-k), num(top), num(x+r), num(top))
	p.op("h f")
}

// ---------------------------------------------------------------------------
// Text
// ---------------------------------------------------------------------------

// Text draws a single line of text with its baseline at y and returns the
// advance width.
func (p *Page) Text(x, y float64, s string, f Font, size float64, c Color) float64 {
	return p.TextTracked(x, y, s, f, size, 0, c)
}

// TextTracked draws text with extra letter spacing, used for section labels.
func (p *Page) TextTracked(x, y float64, s string, f Font, size, tracking float64, c Color) float64 {
	if s == "" {
		return 0
	}
	p.usedFonts[f] = true
	p.beginContent(p.Role(), "")
	p.SetFill(c)
	p.op("BT")
	if tracking != 0 {
		p.op("%s Tc", num(tracking))
	}
	p.op("/%s %s Tf", f.resource(), num(size))
	p.op("1 0 0 1 %s %s Tm", num(x), num(p.y(y)))
	p.op("%s Tj", Str(s))
	if tracking != 0 {
		p.op("0 Tc")
	}
	p.op("ET")
	p.endContent()
	return TextWidth(s, f, size) + tracking*float64(len([]rune(s))-1)
}

// TextRight draws text ending at x.
func (p *Page) TextRight(x, y float64, s string, f Font, size float64, c Color) {
	p.Text(x-TextWidth(s, f, size), y, s, f, size, c)
}

// TextCenter draws text centred on x.
func (p *Page) TextCenter(x, y float64, s string, f Font, size float64, c Color) {
	p.Text(x-TextWidth(s, f, size)/2, y, s, f, size, c)
}

// TextBlock word-wraps s to width and draws it starting with the first
// baseline at y. It returns the Y coordinate just past the last line.
func (p *Page) TextBlock(x, y, width float64, s string, f Font, size, leading float64, c Color) float64 {
	for _, line := range WrapText(s, f, size, width) {
		p.Text(x, y, line, f, size, c)
		y += leading
	}
	return y
}

// ---------------------------------------------------------------------------
// Images
// ---------------------------------------------------------------------------

// DrawImage paints img as decoration, contributing nothing to the structure.
func (p *Page) DrawImage(img *Image, x, y, w, h float64) {
	p.drawImage(img, x, y, w, h, RoleArtifact, "")
}

// DrawImageAlt paints img as a figure with alternative text, so a screen reader
// can describe it.
func (p *Page) DrawImageAlt(img *Image, x, y, w, h float64, alt string) {
	p.drawImage(img, x, y, w, h, RoleFigure, alt)
}

func (p *Page) drawImage(img *Image, x, y, w, h float64, role, alt string) {
	if img.ref == 0 {
		img.ref = p.doc.addStream(fmt.Sprintf(
			"/Type /XObject /Subtype /Image /Width %d /Height %d"+
				" /ColorSpace /%s /BitsPerComponent 8 /Filter /DCTDecode",
			img.Width, img.Height, img.colorSpace()), img.data, false)
	}
	name := ""
	for _, u := range p.usedImages {
		if u.img == img {
			name = u.name
			break
		}
	}
	if name == "" {
		name = fmt.Sprintf("Im%d", len(p.usedImages)+1)
		p.usedImages = append(p.usedImages, pageImage{name, img})
	}
	p.beginContent(role, alt)
	p.Save()
	p.op("%s 0 0 %s %s %s cm", num(w), num(h), num(x), num(p.y(y+h)))
	p.op("/%s Do", name)
	p.Restore()
	p.endContent()
}

// ---------------------------------------------------------------------------
// Links
// ---------------------------------------------------------------------------

// Link makes the rectangle whose top-left corner is (x, y) open uri. The border
// is suppressed so the annotation is invisible over styled text.
func (p *Page) Link(x, y, w, h float64, uri string) {
	if uri == "" {
		return
	}
	p.annots = append(p.annots, fmt.Sprintf(
		"<< /Type /Annot /Subtype /Link /Rect [%s %s %s %s] /Border [0 0 0]"+
			" /F 4 /A << /S /URI /URI %s >> >>",
		num(x), num(p.y(y+h)), num(x+w), num(p.y(y)), Str(uri)))
}

// LinkText draws text and makes it clickable in one step, returning the advance
// width so callers can lay out inline runs.
func (p *Page) LinkText(x, y float64, s string, f Font, size float64, c Color, uri string) float64 {
	w := p.Text(x, y, s, f, size, c)
	p.Link(x, y-size*0.8, w, size*1.1, uri)
	return w
}

// ---------------------------------------------------------------------------
// XMP metadata
// ---------------------------------------------------------------------------

// xmp builds the XMP packet embedded as the document /Metadata stream. Search
// engines and document pipelines read this in preference to the Info dict.
func (d *Doc) xmp() []byte {
	esc := func(s string) string {
		return strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;").Replace(s)
	}
	var kw strings.Builder
	for _, k := range d.Keywords {
		fmt.Fprintf(&kw, "<rdf:li>%s</rdf:li>", esc(k))
	}
	stamp := d.Created.Format("2006-01-02T15:04:05-07:00")

	var b strings.Builder
	b.WriteString("<?xpacket begin=\"\uFEFF\" id=\"W5M0MpCehiHzreSzNTczkc9d\"?>\n")
	b.WriteString(`<x:xmpmeta xmlns:x="adobe:ns:meta/">` + "\n")
	b.WriteString(`<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">` + "\n")
	fmt.Fprintf(&b, `<rdf:Description rdf:about=""
  xmlns:dc="http://purl.org/dc/elements/1.1/"
  xmlns:xmp="http://ns.adobe.com/xap/1.0/"
  xmlns:pdf="http://ns.adobe.com/pdf/1.3/">
 <dc:format>application/pdf</dc:format>
 <dc:title><rdf:Alt><rdf:li xml:lang="x-default">%s</rdf:li></rdf:Alt></dc:title>
 <dc:creator><rdf:Seq><rdf:li>%s</rdf:li></rdf:Seq></dc:creator>
 <dc:description><rdf:Alt><rdf:li xml:lang="x-default">%s</rdf:li></rdf:Alt></dc:description>
 <dc:subject><rdf:Bag>%s</rdf:Bag></dc:subject>
 <dc:language><rdf:Bag><rdf:li>%s</rdf:li></rdf:Bag></dc:language>
 <xmp:CreatorTool>%s</xmp:CreatorTool>
 <xmp:CreateDate>%s</xmp:CreateDate>
 <xmp:ModifyDate>%s</xmp:ModifyDate>
 <pdf:Producer>ewerton-go/pdf (pure Go stdlib)</pdf:Producer>
</rdf:Description>
`, esc(d.Title), esc(d.Author), esc(d.Subject), kw.String(), esc(d.Lang), esc(d.Creator), stamp, stamp)
	b.WriteString("</rdf:RDF>\n</x:xmpmeta>\n<?xpacket end=\"w\"?>")
	return []byte(b.String())
}

// Polygon fills a closed shape through the given top-left-origin points.
func (p *Page) Polygon(points [][2]float64, c Color) {
	if len(points) < 3 {
		return
	}
	p.op("/Artifact BMC")
	defer p.op("EMC")
	p.SetFill(c)
	for i, pt := range points {
		verb := "l"
		if i == 0 {
			verb = "m"
		}
		p.op("%s %s %s", num(pt[0]), num(p.y(pt[1])), verb)
	}
	p.op("h f")
}

// StrokeRect outlines a rectangle whose top-left corner is (x, y).
func (p *Page) StrokeRect(x, y, w, h, width float64, c Color) {
	p.artifact(func() {
		p.SetStroke(c)
		p.SetLineWidth(width)
		p.op("%s %s %s %s re S", num(x), num(p.y(y+h)), num(w), num(h))
	})
}

// ClipRect intersects the clipping path with a rectangle. Wrap in Save/Restore.
func (p *Page) ClipRect(x, y, w, h float64) {
	p.op("%s %s %s %s re W n", num(x), num(p.y(y+h)), num(w), num(h))
}

// Curve strokes a cubic Bezier from (x1,y1) to (x4,y4) with the given control
// points, all in top-left coordinates.
func (p *Page) Curve(x1, y1, x2, y2, x3, y3, x4, y4 float64) {
	p.op("/Artifact BMC")
	defer p.op("EMC")
	p.op("%s %s m %s %s %s %s %s %s c S",
		num(x1), num(p.y(y1)), num(x2), num(p.y(y2)),
		num(x3), num(p.y(y3)), num(x4), num(p.y(y4)))
}
