package pdf

import (
	"fmt"
	"strings"
)

// Axial gradients and constant alpha. Between them they cover everything the
// header band needs to look like a designed surface rather than a filled box.

// namedRef ties a page-local resource name to a document object.
type namedRef struct {
	name string
	ref  Ref
}

// Stop is one colour stop of a gradient, positioned in [0,1].
type Stop struct {
	At    float64
	Color Color
}

// GradientRect fills a rectangle with a linear gradient.
//
// The axis is given as fractions of the rectangle, so (0,0)->(1,1) is a
// diagonal, (0,0)->(1,0) is left-to-right and (0,0)->(0,1) is top-to-bottom.
func (p *Page) GradientRect(x, y, w, h float64, x0f, y0f, x1f, y1f float64, stops []Stop) {
	if len(stops) < 2 {
		return
	}
	ref := p.doc.add(fmt.Sprintf(
		"<< /ShadingType 2 /ColorSpace /DeviceRGB /Coords [%s %s %s %s] /Function %s /Extend [true true] >>",
		num(x+x0f*w), num(p.y(y+y0f*h)), num(x+x1f*w), num(p.y(y+y1f*h)),
		p.doc.gradientFunction(stops)))

	name := fmt.Sprintf("Sh%d", len(p.usedShadings)+1)
	p.usedShadings = append(p.usedShadings, namedRef{name, ref})

	p.artifact(func() {
		p.Save()
		p.op("%s %s %s %s re W n", num(x), num(p.y(y+h)), num(w), num(h))
		p.op("/%s sh", name)
		p.Restore()
	})
}

// gradientFunction builds the PDF function that interpolates the stops: a
// single exponential function for two stops, or a stitching function chaining
// one per segment for more.
func (d *Doc) gradientFunction(stops []Stop) Ref {
	seg := func(a, b Color) string {
		return fmt.Sprintf("<< /FunctionType 2 /Domain [0 1] /C0 [%s %s %s] /C1 [%s %s %s] /N 1 >>",
			num(a.R), num(a.G), num(a.B), num(b.R), num(b.G), num(b.B))
	}
	if len(stops) == 2 {
		return d.add(seg(stops[0].Color, stops[1].Color))
	}

	var funcs, bounds, encode []string
	for i := 0; i < len(stops)-1; i++ {
		funcs = append(funcs, d.add(seg(stops[i].Color, stops[i+1].Color)).String())
		encode = append(encode, "0 1")
		if i > 0 {
			bounds = append(bounds, num(stops[i].At))
		}
	}
	return d.add(fmt.Sprintf(
		"<< /FunctionType 3 /Domain [0 1] /Functions [%s] /Bounds [%s] /Encode [%s] >>",
		strings.Join(funcs, " "), strings.Join(bounds, " "), strings.Join(encode, " ")))
}

// SetAlpha sets the constant fill and stroke opacity for subsequent painting.
// Wrap it in Save/Restore, as it is part of the graphics state.
func (p *Page) SetAlpha(a float64) {
	ref := p.doc.add(fmt.Sprintf("<< /Type /ExtGState /ca %s /CA %s >>", num(a), num(a)))
	name := fmt.Sprintf("GS%d", len(p.usedGStates)+1)
	p.usedGStates = append(p.usedGStates, namedRef{name, ref})
	p.op("/%s gs", name)
}
