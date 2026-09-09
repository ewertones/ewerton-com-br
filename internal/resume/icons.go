package resume

import "ewerton-go/internal/pdf"

// Small vector glyphs drawn straight into the content stream.
//
// PDF has no icon font among the base-14 faces, and embedding one would undo
// the point of using them, so the handful of marks the resume needs are drawn
// as paths. Each icon is defined on a square box of side s with its top-left
// corner at (x, y), and stroked in the given colour.

const iconStroke = 0.9

// iconGlobe: a circle crossed by an equator and a meridian.
func iconGlobe(p *pdf.Page, x, y, s float64, c pdf.Color) {
	r := s / 2
	cx, cy := x+r, y+r
	p.StrokeCircle(cx, cy, r, iconStroke, c)
	p.Line(cx-r, cy, cx+r, cy, iconStroke, c)
	// Meridian, approximated by a narrow ellipse drawn as two arcs.
	p.SetStroke(c)
	p.SetLineWidth(iconStroke)
	p.Curve(cx, cy-r, cx-r*0.75, cy-r*0.45, cx-r*0.75, cy+r*0.45, cx, cy+r)
	p.Curve(cx, cy-r, cx+r*0.75, cy-r*0.45, cx+r*0.75, cy+r*0.45, cx, cy+r)
}

// iconMail: an envelope with its flap.
func iconMail(p *pdf.Page, x, y, s float64, c pdf.Color) {
	h := s * 0.76
	top := y + (s-h)/2
	p.StrokeRect(x, top, s, h, iconStroke, c)
	p.Line(x, top, x+s/2, top+h*0.58, iconStroke, c)
	p.Line(x+s, top, x+s/2, top+h*0.58, iconStroke, c)
}

// iconPhone: a handset body with a speaker line.
func iconPhone(p *pdf.Page, x, y, s float64, c pdf.Color) {
	w := s * 0.62
	left := x + (s-w)/2
	p.StrokeRect(left, y, w, s, iconStroke, c)
	p.Line(left+w*0.3, y+s*0.86, left+w*0.7, y+s*0.86, iconStroke, c)
}

// iconPin: a map marker, a teardrop over a hollow centre.
func iconPin(p *pdf.Page, x, y, s float64, c pdf.Color) {
	r := s * 0.36
	cx, cy := x+s/2, y+r+s*0.06
	p.StrokeCircle(cx, cy, r, iconStroke, c)
	p.Line(cx-r*0.72, cy+r*0.7, cx, y+s, iconStroke, c)
	p.Line(cx+r*0.72, cy+r*0.7, cx, y+s, iconStroke, c)
}

// iconLinkedIn: the familiar rounded badge with a lowercase "in".
func iconLinkedIn(p *pdf.Page, x, y, s float64, c pdf.Color) {
	p.RoundRect(x, y, s, s, s*0.22, c)
	p.TextCenter(x+s/2, y+s*0.76, "in", pdf.Bold, s*0.72, brand)
}

// iconGitHub: the octocat silhouette, reduced to a head with two ears and a
// tail. The ears sit proud of the head, or at this size the mark just reads as
// a circle.
func iconGitHub(p *pdf.Page, x, y, s float64, c pdf.Color) {
	r := s * 0.34
	cx, cy := x+s/2, y+s*0.55
	for _, dir := range []float64{-1, 1} {
		p.Polygon([][2]float64{
			{cx + dir*r*0.95, cy - r*0.30},
			{cx + dir*r*0.78, cy - r*1.45},
			{cx + dir*r*0.10, cy - r*0.80},
		}, c)
	}
	p.Circle(cx, cy, r, c)
	p.Line(cx-r*0.30, cy+r*0.95, cx-r*0.30, y+s, iconStroke*0.9, c)
}

// ---------------------------------------------------------------------------
// Flags
// ---------------------------------------------------------------------------

// Flag colours, taken from each country's official specification.
var (
	flagGreen    = pdf.RGB(0x00, 0x9C, 0x3B)
	flagYellow   = pdf.RGB(0xFF, 0xDF, 0x00)
	flagNavy     = pdf.RGB(0x00, 0x23, 0x76)
	flagOldGlory = pdf.RGB(0xB3, 0x14, 0x2D)
	flagUSBlue   = pdf.RGB(0x0A, 0x31, 0x61)
	flagFrBlue   = pdf.RGB(0x00, 0x55, 0xA4)
	flagFrRed    = pdf.RGB(0xEF, 0x41, 0x35)
	flagEsRed    = pdf.RGB(0xAA, 0x15, 0x1B)
	flagEsYellow = pdf.RGB(0xF1, 0xBF, 0x00)
	flagBorder   = pdf.RGB(0xB9, 0xC2, 0xD2)
)

// drawFlag renders a simplified national flag in a w by h box with its
// top-left corner at (x, y). Detail is dropped deliberately: at this size the
// silhouette and colour bands are what make a flag recognisable.
func drawFlag(p *pdf.Page, code string, x, y, w, h float64) {
	switch code {
	case "BR":
		p.Rect(x, y, w, h, flagGreen)
		p.Polygon([][2]float64{
			{x + w/2, y + h*0.1}, {x + w*0.91, y + h/2},
			{x + w/2, y + h*0.9}, {x + w*0.09, y + h/2},
		}, flagYellow)
		p.Circle(x+w/2, y+h/2, h*0.21, flagNavy)

	case "US":
		// Seven red bands over white, with the canton on the upper hoist.
		p.Rect(x, y, w, h, paper)
		band := h / 7
		for i := 0; i < 7; i += 2 {
			p.Rect(x, y+float64(i)*band, w, band, flagOldGlory)
		}
		p.Rect(x, y, w*0.42, band*4, flagUSBlue)
		for i := 0; i < 4; i++ {
			cx := x + w*0.13 + float64(i%2)*w*0.16
			cy := y + band*1.1 + float64(i/2)*band*1.7
			p.Circle(cx, cy, 0.42, paper)
		}

	case "FR":
		p.Rect(x, y, w/3, h, flagFrBlue)
		p.Rect(x+w/3, y, w/3, h, paper)
		p.Rect(x+2*w/3, y, w/3, h, flagFrRed)

	case "ES":
		p.Rect(x, y, w, h/4, flagEsRed)
		p.Rect(x, y+h/4, w, h/2, flagEsYellow)
		p.Rect(x, y+3*h/4, w, h/4, flagEsRed)
	}
	p.StrokeRect(x, y, w, h, 0.4, flagBorder)
}
