// Package resume turns the canonical profile model into the documents the site
// publishes: a designed PDF, a JSON Resume document, a Markdown rendering and a
// schema.org graph. All four are generated from data.Resume, so they can never
// disagree with the landing page.
package resume

import (
	"fmt"
	"strings"

	"ewerton-go/internal/data"
	"ewerton-go/internal/pdf"
)

// Palette.
//
// The primary is the blue both products already use: aiomover and questao share
// oklch(0.3938 0.1259 254.13), which converts to #004587. It is lighter and
// more open than the portfolio's own near-navy, so the header reads as a lit
// surface rather than a slab.
var (
	brand       = pdf.RGB(0x00, 0x45, 0x87)
	brandDeep   = pdf.RGB(0x00, 0x33, 0x66)
	brandMid    = pdf.RGB(0x0A, 0x5D, 0xA8)
	brandBright = pdf.RGB(0x18, 0x7C, 0xD4)
	brandSoft   = pdf.RGB(0xA3, 0xC2, 0xE2)
	brandWash   = pdf.RGB(0xEA, 0xF1, 0xF9)
	bandText    = pdf.RGB(0xCB, 0xDF, 0xF5)
	ink         = pdf.RGB(0x16, 0x1B, 0x26)
	muted       = pdf.RGB(0x57, 0x62, 0x70)
	faint       = pdf.RGB(0xD6, 0xDE, 0xE8)
	paper       = pdf.RGB(0xFF, 0xFF, 0xFF)
	accentGreen = pdf.RGB(0x0E, 0x9F, 0x4F) // marks the role currently held
	brazilGrn   = pdf.RGB(0x00, 0x97, 0x39)
	brazilGold  = pdf.RGB(0xFE, 0xDD, 0x00)

	starFill  = pdf.RGB(0x18, 0x6A, 0xB8)
	starEmpty = pdf.RGB(0xD6, 0xE1, 0xEF)
)

// Page geometry, in points.
const (
	marginX    = 48.0
	bandHeight = 134.0
	footerY    = pdf.A4Height - 40 // baseline of the footer row
	footerRule = footerY - 11.5    // just clear of the text below it
	bottomLim  = pdf.A4Height - 58
	contentTop = 64.0

	contentW = pdf.A4Width - 2*marginX
	textX    = marginX + 13 // entry body, clear of the section marker
	rightX   = pdf.A4Width - marginX
)

// Body text metrics, shared by every prose block so a block can be measured
// before it is drawn.
const (
	bodySize = 8.9
	bodyLead = 12.4
)

// A section heading needs sectionGapAbove of air before it and sectionGapBelow
// after its rule. sectionLead is the pair, which layout code needs to plan page
// breaks before anything is drawn.
const (
	sectionGapAbove = 36.0
	sectionGapBelow = 26.0
	sectionLead     = sectionGapAbove + sectionGapBelow
)

// keepTogether caps how much of a block we insist on keeping on one page. A
// long entry may split, but never before its heading and opening lines.
const keepTogether = 104.0

// layout carries the cursor and shared state while the document is built.
// Assets holds the raster artwork the document embeds. Both are optional; the
// header falls back to a monogram and the footer simply omits the mark.
type Assets struct {
	Photo []byte // portrait, JPEG
	Logo  []byte // logo mark, JPEG on a white ground
}

type layout struct {
	doc  *pdf.Doc
	page *pdf.Page
	y    float64
	// pageTop is where content starts on the current page, so a section can
	// tell whether it is opening a page or continuing one.
	pageTop float64
	r       data.Resume
	t       func(string) string
	photo   *pdf.Image
	logo    *pdf.Image
}

// RenderPDF lays out the resume. attachment, when present, is embedded inside
// the PDF as resume.json so tools can read the structured data without a second
// request.
func RenderPDF(r data.Resume, t func(string) string, assets Assets, attachment []byte) ([]byte, error) {
	doc := pdf.New()
	doc.Title = fmt.Sprintf("%s — %s", r.Basics.Name, r.Basics.Headline)
	doc.Author = r.Basics.Name
	doc.Subject = r.Basics.Summary
	doc.Keywords = append([]string{r.Basics.Headline, r.Basics.Location()}, r.Keywords()...)
	doc.Creator = data.SiteOrigin
	doc.Lang = r.Locale
	doc.Created = r.Generated

	if len(attachment) > 0 {
		doc.Attach("resume.json", "application/json",
			"JSON Resume (jsonresume.org schema v1.0.0) for this document", attachment)
	}

	l := &layout{doc: doc, r: r, t: t}
	if img, err := pdf.LoadJPEG(assets.Photo); err == nil {
		l.photo = img
	}
	if img, err := pdf.LoadJPEG(assets.Logo); err == nil {
		l.logo = img
	}

	l.newPage(true)
	l.header()
	l.summary()
	l.experience()
	l.education()
	l.languages()
	l.projects()
	l.skills()
	l.footers()

	return doc.Bytes(), nil
}

// ---------------------------------------------------------------------------
// Page furniture
// ---------------------------------------------------------------------------

func (l *layout) newPage(first bool) {
	l.page = l.doc.NewPage(pdf.A4Width, pdf.A4Height)
	if first {
		l.y = bandHeight + 24
		l.pageTop = l.y
		return
	}
	// Continuation pages carry a compact running header. It repeats on every
	// page, so it is furniture rather than content and is tagged as an artifact.
	l.page.SetRole(pdf.RoleArtifact)
	l.page.Text(marginX, 38, l.r.Basics.Name, pdf.Bold, 8.5, brand)
	w := pdf.TextWidth(l.r.Basics.Name, pdf.Bold, 8.5)
	l.page.Text(marginX+w+6, 38, "· "+l.r.Basics.Headline, pdf.Regular, 8.5, muted)
	l.page.TextRight(rightX, 38, hostOf(data.SiteOrigin), pdf.Regular, 8.5, muted)
	l.page.Line(marginX, 46, rightX, 46, 0.6, faint)
	l.page.SetRole(pdf.RoleP)
	l.y = contentTop + 10
	l.pageTop = l.y
}

// ensure starts a new page when the next block would not fit.
func (l *layout) ensure(height float64) {
	if l.y+height > bottomLim {
		l.newPage(false)
	}
}

// footers stamps every page once the total page count is known.
func (l *layout) footers() {
	pages := l.doc.Pages()
	total := len(pages)

	// ISO date: unambiguous for a reader anywhere, and machine-parseable.
	generated := l.t("resume.generatedOn") + " " + l.r.Generated.Format("2006-01-02")

	for i, p := range pages {
		p.SetRole(pdf.RoleArtifact)
		p.Line(marginX, footerRule, rightX, footerRule, 0.6, faint)
		p.Text(marginX, footerY, generated, pdf.Regular, 7.5, muted)

		if l.logo != nil {
			// Sized and centred on the cap height of the footer text either
			// side of it, so the three items read as one line.
			const mark = 7.6
			const capHeight = 7.5 * 0.717
			p.DrawImage(l.logo, (pdf.A4Width-mark)/2, footerY-capHeight/2-mark/2, mark, mark)
		}
		p.TextRight(rightX, footerY, fmt.Sprintf("%d / %d", i+1, total), pdf.Regular, 7.5, muted)
	}
}

// ---------------------------------------------------------------------------
// Header band
// ---------------------------------------------------------------------------

func (l *layout) header() {
	p := l.page
	b := l.r.Basics

	// A diagonal gradient, deeper at the top-left easing into a brighter azure
	// at the bottom-right, so the band reads as a lit surface rather than one
	// flat rectangle of colour.
	p.GradientRect(0, 0, pdf.A4Width, bandHeight, 0, 0, 1, 1, []pdf.Stop{
		{At: 0, Color: brandDeep},
		{At: 0.5, Color: brand},
		{At: 1, Color: brandBright},
	})

	// Two very soft orbs add depth without competing with the text.
	p.Save()
	p.SetAlpha(0.08)
	p.Circle(pdf.A4Width-58, 10, 92, paper)
	p.Circle(pdf.A4Width-214, bandHeight+6, 56, paper)
	p.Restore()

	// One ribbon closes the band, sweeping the brand blue through the green and
	// gold of the flag instead of butting three hard blocks together.
	p.GradientRect(0, bandHeight, pdf.A4Width, 3.5, 0, 0, 1, 0, []pdf.Stop{
		{At: 0, Color: brandBright},
		{At: 0.52, Color: brandBright},
		{At: 0.82, Color: brazilGrn},
		{At: 1, Color: brazilGold},
	})

	// Portrait, circular-masked, sitting on a soft halo.
	const cx, cy, rad = 86.0, 67.0, 38.0
	p.Save()
	p.SetAlpha(0.14)
	p.Circle(cx, cy, rad+7, paper)
	p.Restore()

	if l.photo != nil {
		p.Save()
		p.ClipCircle(cx, cy, rad)
		p.DrawImageAlt(l.photo, cx-rad, cy-rad, rad*2, rad*2, "Portrait of "+b.Name)
		p.Restore()
	} else {
		p.Circle(cx, cy, rad, brandBright)
		p.TextCenter(cx, cy+10, initials(b.Name), pdf.Bold, 26, paper)
	}
	p.Save()
	p.SetAlpha(0.6)
	p.StrokeCircle(cx, cy, rad, 1.4, paper)
	p.Restore()

	const tx = 142.0
	p.SetRole(pdf.RoleH1)
	p.Text(tx, 55, b.Name, pdf.Bold, 23, paper)
	p.SetRole(pdf.RoleP)
	p.TextTracked(tx, 74, b.Headline, pdf.Regular, 10.6, 0.3, bandText)
	p.Line(tx, 82, tx+46, 82, 1.5, brazilGold)

	// Contact details, laid out as a grid so the second row lines up under the
	// first. Every item carries a small drawn glyph and, bar the location, is a
	// live link.
	l.contactGrid(tx, 101, 17, [][]contact{
		{
			{hostOf(b.URL), b.URL, iconGlobe},
			{b.Email, "mailto:" + b.Email, iconMail},
			{b.Phone, b.PhoneURL, iconPhone},
		},
		{
			{profile(b, "LinkedIn"), profileURL(b, "LinkedIn"), iconLinkedIn},
			{profile(b, "GitHub"), profileURL(b, "GitHub"), iconGitHub},
			{b.Location(), b.LocationURL, iconPin},
		},
	})
}

type contact struct {
	label string
	uri   string
	icon  func(*pdf.Page, float64, float64, float64, pdf.Color)
}

// Contact grid metrics.
const (
	contactSize  = 8.5
	contactGlyph = 8.2
	contactPad   = 5.0  // glyph to label
	contactGap   = 18.0 // between columns
)

// contactGrid draws contact details in aligned columns. Column widths come from
// the widest cell in each column, so rows line up however long the labels are.
func (l *layout) contactGrid(x0, y0, rowH float64, rows [][]contact) {
	cols := 0
	for _, row := range rows {
		cols = max(cols, len(row))
	}
	widths := make([]float64, cols)
	for _, row := range rows {
		for j, it := range row {
			if it.label == "" {
				continue
			}
			widths[j] = max(widths[j], pdf.TextWidth(it.label, pdf.Regular, contactSize))
		}
	}

	for i, row := range rows {
		y := y0 + float64(i)*rowH
		x := x0
		for j, it := range row {
			if it.label != "" {
				it.icon(l.page, x, y-contactGlyph+1.4, contactGlyph, bandText)
				label := x + contactGlyph + contactPad
				if it.uri == "" {
					l.page.Text(label, y, it.label, pdf.Regular, contactSize, paper)
				} else {
					l.page.LinkText(label, y, it.label, pdf.Regular, contactSize, paper, it.uri)
				}
			}
			x += contactGlyph + contactPad + widths[j] + contactGap
		}
	}
}

// ---------------------------------------------------------------------------
// Sections
// ---------------------------------------------------------------------------

// section draws a titled rule and registers a PDF bookmark for it. minFollow is
// how much of the section's first block must fit on the same page, so a heading
// is never stranded at the foot of one.
func (l *layout) section(key string, minFollow float64) {
	label := strings.ToUpper(l.t("resume.sections." + key))
	l.ensure(sectionLead + minFollow)

	// The gap above a heading separates it from the section before it. At the
	// top of a page there is nothing to separate it from, and the running
	// header already provides the margin, so skip it rather than opening the
	// page with a band of blank paper.
	if l.y > l.pageTop {
		l.y += sectionGapAbove
	}

	l.page.RoundRect(marginX, l.y-8.5, 3.5, 11, 1.75, brand)
	l.page.SetRole(pdf.RoleH2)
	w := l.page.TextTracked(marginX+9, l.y, label, pdf.Bold, 10, 1.1, brand)
	l.page.SetRole(pdf.RoleP)
	l.page.Line(marginX+13+w, l.y-3.5, rightX, l.y-3.5, 0.6, faint)

	l.doc.Bookmark(l.t("resume.sections."+key), l.page, l.y-sectionGapAbove-2)
	l.y += sectionGapBelow
}

func (l *layout) summary() {
	l.section("summary", 40)
	l.y = l.page.TextBlock(marginX, l.y, contentW, l.r.Basics.Summary, pdf.Regular, 9.3, 13.4, muted)
}

// ---------------------------------------------------------------------------
// Experience
// ---------------------------------------------------------------------------

func (l *layout) experience() {
	groups := l.r.WorkByCompany()
	l.section("experience", min(l.companyHeight(groups[0]), keepTogether))
	for i, g := range groups {
		l.company(g)
		if i < len(groups)-1 {
			l.y += 17
		}
	}
}

// company renders one employer. Several roles at the same employer nest under a
// single company heading, the way a LinkedIn tenure reads.
func (l *layout) company(g data.CompanyGroup) {
	l.ensure(min(l.companyHeight(g), keepTogether))

	l.page.Circle(marginX+3.5, l.y-3.5, 3.5, brand)
	l.page.SetRole(pdf.RoleH3)
	x := textX + l.page.LinkText(textX, l.y, g.Company, pdf.Bold, 12.5, brand, g.CompanyURL) + 8
	l.page.SetRole(pdf.RoleP)
	if g.DisplayURL != "" {
		x += l.page.LinkText(x, l.y, g.DisplayURL, pdf.Regular, 7.8, brandSoft, g.CompanyURL) + 8
	}
	if g.Industry != "" {
		l.chip(x, l.y-8, g.Industry)
	}
	if g.Current {
		// A gold dot marks the role currently held.
		l.page.Circle(rightX-pdf.TextWidth(g.Period, pdf.Regular, 8.8)-9, l.y-3, 2.6, accentGreen)
	}
	l.page.TextRight(rightX, l.y, g.Period, pdf.Regular, 8.8, muted)
	l.y += 16

	multi := len(g.Roles) > 1
	roleX := textX
	if multi {
		roleX = textX + 12
	}

	for i, role := range g.Roles {
		if i > 0 {
			l.y += 11
		}
		l.ensure(36)
		if multi {
			// A hollow marker per role. The dates stay on the company line,
			// which already spans the whole tenure.
			l.page.StrokeCircle(textX+4, l.y-3.2, 2.8, 1, brandSoft)
		}
		l.page.Text(roleX, l.y, role.Role, pdf.Bold, 9.8, ink)
		l.y += 14

		if role.Summary != "" {
			l.y = l.page.TextBlock(roleX, l.y, rightX-roleX, role.Summary, pdf.Regular, bodySize, bodyLead, muted)
			l.y += 4
		}
		for _, h := range role.Highlights {
			l.bullet(roleX, rightX-roleX, h)
		}
	}
}

// companyHeight measures a tenure so pagination can decide whether the whole
// block fits, rather than guessing with a fixed reserve and stranding an entry
// that would have fitted.
func (l *layout) companyHeight(g data.CompanyGroup) float64 {
	h := 16.0
	roleX := textX
	if len(g.Roles) > 1 {
		roleX = textX + 12
	}
	width := rightX - roleX
	for i, role := range g.Roles {
		if i > 0 {
			h += 11
		}
		h += 14
		if role.Summary != "" {
			h += float64(len(pdf.WrapText(role.Summary, pdf.Regular, bodySize, width)))*bodyLead + 4
		}
		for _, b := range role.Highlights {
			h += float64(len(pdf.WrapText(b, pdf.Regular, bodySize, width-12))) * bodyLead
		}
	}
	return h
}

// chip draws a small tinted pill, used for industry, grade and feature labels.
func (l *layout) chip(x, y float64, label string) float64 {
	const size = 7.4
	w := pdf.TextWidth(label, pdf.Regular, size) + 13
	l.page.RoundRect(x, y, w, 11.8, 5.9, brandWash)
	l.page.Text(x+6.5, y+8.2, label, pdf.Regular, size, brandMid)
	return w
}

// bullet draws one hanging-indent list item and advances the cursor.
func (l *layout) bullet(x, width float64, text string) {
	lines := pdf.WrapText(text, pdf.Regular, bodySize, width-12)
	l.ensure(float64(len(lines)) * bodyLead)

	l.page.Circle(x+3, l.y-3, 1.6, brandSoft)
	for i, line := range lines {
		l.page.Text(x+12, l.y, line, pdf.Regular, bodySize, i2c(i, ink, muted))
		l.y += bodyLead
	}
}

// i2c keeps the first line of a bullet slightly darker than its continuation,
// which makes a dense list far easier to skim.
func i2c(i int, first, rest pdf.Color) pdf.Color {
	if i == 0 {
		return first
	}
	return rest
}

// ---------------------------------------------------------------------------
// Education
// ---------------------------------------------------------------------------

func (l *layout) education() {
	l.section("education", 60)
	for i, e := range l.r.Education {
		l.ensure(60)

		l.page.Circle(marginX+3.5, l.y-3.5, 3.5, brand)
		l.page.SetRole(pdf.RoleH3)
		l.page.LinkText(textX, l.y, e.Institution, pdf.Bold, 11.5, brand, e.WebsiteURL)
		l.page.SetRole(pdf.RoleP)
		l.page.TextRight(rightX, l.y, e.Period, pdf.Regular, 8.6, muted)
		l.y += 14

		dw := l.page.Text(textX, l.y, e.Degree, pdf.Bold, 9.4, ink)
		if e.Grade != "" {
			l.chip(textX+dw+9, l.y-8, e.Grade)
		}
		l.y += 14

		for _, pt := range e.Points {
			l.y = l.page.TextBlock(textX, l.y, rightX-textX, pt, pdf.Regular, bodySize, bodyLead, muted)
		}
		if e.Skills != "" {
			l.y += 2
			l.y = l.page.TextBlock(textX, l.y, rightX-textX, e.Skills, pdf.Italic, 8.6, 11.8, brandMid)
		}
		if i < len(l.r.Education)-1 {
			l.y += 14
		}
	}
}

// ---------------------------------------------------------------------------
// Products
// ---------------------------------------------------------------------------

func (l *layout) projects() {
	l.section("projects", 88)
	for i, p := range l.r.Projects {
		l.ensure(88)

		l.page.Circle(marginX+3.5, l.y-3.5, 3.5, brandMid)
		l.page.SetRole(pdf.RoleH3)
		x := textX + l.page.LinkText(textX, l.y, p.Name, pdf.Bold, 12, brand, p.URL) + 9
		l.page.SetRole(pdf.RoleP)

		if p.URL == "" {
			// Not launched yet: keep the label, but the destination is still a
			// real place, so make it clickable rather than dead text.
			x += l.chip(x, l.y-8, l.t("common.comingSoon")) + 7
		}
		l.page.LinkText(x, l.y, p.DisplayURL, pdf.Regular, 8.4, brandMid, p.HomeURL)
		l.page.TextRight(rightX, l.y, p.Kind, pdf.Regular, 8.4, muted)
		l.y += 17

		l.y = l.page.TextBlock(textX, l.y, rightX-textX, p.Tagline, pdf.Italic, 9.2, 12.4, brandMid)
		l.y += 5
		l.y = l.page.TextBlock(textX, l.y, rightX-textX, p.Description, pdf.Regular, bodySize, bodyLead, muted)
		l.y += 12

		l.chipRow(textX, p.Features)
		l.stackLine(textX, p.MainStack())

		if i < len(l.r.Projects)-1 {
			l.y += 30
		}
	}
}

// chipRow lays chips out left to right, wrapping at the right margin.
func (l *layout) chipRow(x float64, labels []string) {
	cursor := x
	for _, label := range labels {
		if cursor+pdf.TextWidth(label, pdf.Regular, 7.4)+13 > rightX {
			cursor = x
			l.y += 16
		}
		cursor += l.chip(cursor, l.y-8, label) + 5
	}
	l.y += 21
}

// stackLine prints what a product is actually built with, as a labelled run.
func (l *layout) stackLine(x float64, stack []string) {
	if len(stack) == 0 {
		return
	}
	label := strings.ToUpper(l.t("resume.stack"))
	w := l.page.TextTracked(x, l.y, label, pdf.Bold, 7, 0.7, brandMid) + 10

	for i, line := range pdf.WrapText(strings.Join(stack, " · "), pdf.Regular, 8.1, rightX-x-w) {
		if i > 0 {
			l.y += 11
		}
		l.page.Text(x+w, l.y, line, pdf.Regular, 8.1, ink)
	}
	l.y += 12
}

// ---------------------------------------------------------------------------
// Languages
// ---------------------------------------------------------------------------

func (l *layout) languages() {
	l.section("languages", 44)

	const (
		gutter     = 20.0
		flagW      = 12.6
		flagH      = 8.8
		afterFlag  = 16.0
		afterName  = 10.0
		afterStars = 8.0
		nameSize   = 9.0
	)
	colW := (contentW - gutter) / 2
	l.ensure(44)

	// Reserve one column for the longest name so the stars and fluency chips
	// line up down the page instead of stepping in and out with each word.
	nameW := 0.0
	for _, lang := range l.r.Languages {
		nameW = max(nameW, pdf.TextWidth(lang.Name, pdf.Bold, nameSize))
	}

	for i, lang := range l.r.Languages {
		col := i % 2
		if col == 0 && i > 0 {
			l.y += 19
		}
		x := marginX + float64(col)*(colW+gutter)

		drawFlag(l.page, lang.Code, x, l.y-7.4, flagW, flagH)
		l.page.Text(x+afterFlag, l.y, lang.Name, pdf.Bold, nameSize, ink)

		starsX := x + afterFlag + nameW + afterName
		l.page.StarRating(starsX, l.y-starDrop, starR, starGap, lang.Rating, 5, starFill, starEmpty)
		l.chip(starsX+starsW+afterStars, l.y-8, lang.Fluency)
	}
	l.y += 16
}

// ---------------------------------------------------------------------------
// Skills
// ---------------------------------------------------------------------------

// Star geometry, shared by the skills grid and the languages section.
const (
	rowH = 10.9
	// Distance from a category title's baseline down to its first skill row.
	// It has to clear the divider rule and the ascent of the first row's stars.
	catTitleToRule     = 5.5
	catTitleToFirstRow = 17.0
	catGap             = 14.0
	starR              = 3.1
	starGap            = 1.6
	starsW             = 5*(starR*2+starGap) - starGap
	starDrop           = 2.7 // baseline to star centre
)

// skills flows every category into a three-column grid, each skill carrying its
// self-assessed rating as a row of stars.
func (l *layout) skills() {
	const cols = 3
	const gutter = 26.0
	const minGrid = 130.0 // never leave the heading with less grid than this
	colW := (contentW - gutter*(cols-1)) / cols

	blocks := make([]float64, len(l.r.Skills))
	for i, cat := range l.r.Skills {
		blocks[i] = catTitleToFirstRow + float64(len(cat.Skills))*rowH + catGap
	}

	// Decide the break before the heading is drawn, or the heading is stranded
	// at the foot of the previous page.
	l.section("skills", minGrid)

	top, maxY := l.y, l.y
	for i := 0; i < len(blocks); {
		// Fill to a balanced height rather than to the bottom of the page, so
		// the columns end at roughly the same place. Recomputed per page, so a
		// spill onto the next page still spreads across all three columns.
		limit := balancedHeight(blocks[i:], cols, bottomLim-top)

		for col := 0; col < cols && i < len(blocks); col++ {
			x, y := marginX+float64(col)*(colW+gutter), top
			for i < len(blocks) && (y == top || y+blocks[i] <= top+limit) {
				l.skillCategory(x, y, colW, l.r.Skills[i])
				y += blocks[i]
				i++
			}
			if y > maxY {
				maxY = y
			}
		}
		if i < len(blocks) {
			l.y = maxY
			l.newPage(false)
			top, maxY = l.y, l.y
		}
	}
	l.y = maxY
}

// skillCategory draws one category heading and its rated skills.
func (l *layout) skillCategory(x, y, colW float64, cat data.TechnologyCategory) {
	l.page.SetRole(pdf.RoleH3)
	l.page.TextTracked(x, y, strings.ToUpper(l.t("landing.technologies.categories."+cat.Slug)),
		pdf.Bold, 7.6, 0.5, brand)
	l.page.SetRole(pdf.RoleP)
	l.page.Line(x, y+catTitleToRule, x+colW, y+catTitleToRule, 0.5, faint)

	y += catTitleToFirstRow
	for _, sk := range cat.Skills {
		l.skillRow(x, y, colW, sk)
		y += rowH
	}
}

// balancedHeight finds the shortest column height that still packs every block
// into cols columns while keeping the original order. Binary searching the
// height and greedily filling is enough here, and stays stable as content
// changes.
func balancedHeight(blocks []float64, cols int, available float64) float64 {
	lo, total := 0.0, 0.0
	for _, b := range blocks {
		total += b
		if b > lo {
			lo = b
		}
	}
	hi := total
	for hi-lo > 0.5 {
		mid := (lo + hi) / 2
		if packs(blocks, cols, mid) {
			hi = mid
			continue
		}
		lo = mid
	}
	return min(hi, available)
}

// packs reports whether the blocks fit into cols columns of the given height.
func packs(blocks []float64, cols int, height float64) bool {
	used, col := 0.0, 1
	for _, b := range blocks {
		if b > height {
			return false
		}
		if used+b > height && used > 0 {
			col++
			used = 0
		}
		if col > cols {
			return false
		}
		used += b
	}
	return true
}

func (l *layout) skillRow(x, y, colW float64, s data.Technology) {
	const size = 8.0
	starsX := x + colW - starsW
	l.page.Text(x, y, pdf.Truncate(s.Name, pdf.Regular, size, starsX-x-6), pdf.Regular, size, ink)
	l.page.StarRating(starsX, y-starDrop, starR, starGap, s.Rating, 5, starFill, starEmpty)
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// profile and profileURL look a network up by name, so the header layout does
// not depend on the order profiles happen to be declared in.
func profile(b data.Basics, network string) string {
	for _, p := range b.Profiles {
		if p.Network == network {
			return p.Display
		}
	}
	return ""
}

func profileURL(b data.Basics, network string) string {
	for _, p := range b.Profiles {
		if p.Network == network {
			return p.URL
		}
	}
	return ""
}

// hostOf strips the scheme and any trailing slash from a URL for display.
func hostOf(u string) string {
	s := strings.TrimPrefix(strings.TrimPrefix(u, "https://"), "http://")
	s = strings.TrimPrefix(s, "www.")
	return strings.TrimSuffix(s, "/")
}

// initials builds a monogram, used when no portrait is available.
func initials(name string) string {
	var out []rune
	for _, part := range strings.Fields(name) {
		out = append(out, []rune(part)[0])
		if len(out) == 2 {
			break
		}
	}
	return strings.ToUpper(string(out))
}

// trimZero formats a rating as 4 or 4.5 rather than 4.0.
func trimZero(v float64) string {
	if v == float64(int(v)) {
		return fmt.Sprintf("%d", int(v))
	}
	return strings.TrimRight(fmt.Sprintf("%.1f", v), "0")
}
