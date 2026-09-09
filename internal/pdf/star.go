package pdf

import "math"

// starPoints is the number of points on the star, and starInner the ratio of
// the inner radius to the outer one that makes a classic five-pointed star.
const (
	starPoints = 5
	starInner  = 0.382
)

// starPath emits a five-pointed star centred on (cx, cy) without painting it.
func (p *Page) starPath(cx, cy, r float64) {
	y := p.y(cy)
	for i := 0; i < starPoints*2; i++ {
		radius := r
		if i%2 == 1 {
			radius = r * starInner
		}
		angle := -math.Pi/2 + float64(i)*math.Pi/starPoints
		px := cx + radius*math.Cos(angle)
		py := y - radius*math.Sin(angle)
		if i == 0 {
			p.op("%s %s m", num(px), num(py))
			continue
		}
		p.op("%s %s l", num(px), num(py))
	}
	p.op("h")
}

// Star fills a single star.
func (p *Page) Star(cx, cy, r float64, c Color) {
	p.artifact(func() { p.fillStar(cx, cy, r, c) })
}

// fillStar paints a star without opening a marked-content sequence, for callers
// already inside one.
func (p *Page) fillStar(cx, cy, r float64, c Color) {
	p.SetFill(c)
	p.starPath(cx, cy, r)
	p.op("f")
}

// StarRating draws a row of stars for a 0..max rating, supporting halves.
// It returns the total width so callers can right-align the row.
func (p *Page) StarRating(x, y, r, gap float64, rating, max float64, filled, empty Color) float64 {
	count := int(max)
	pitch := r*2 + gap

	p.op("/Artifact BMC")
	defer p.op("EMC")

	for i := 0; i < count; i++ {
		cx := x + r + float64(i)*pitch
		p.fillStar(cx, y, r, empty)

		switch value := rating - float64(i); {
		case value >= 1:
			p.fillStar(cx, y, r, filled)
		case value > 0:
			// Half (or partial) star: clip to the left portion before filling.
			p.Save()
			p.op("%s %s %s %s re W n", num(cx-r), num(p.y(y+r)), num(r*2*value), num(r*2))
			p.fillStar(cx, y, r, filled)
			p.Restore()
		}
	}
	return float64(count)*pitch - gap
}
