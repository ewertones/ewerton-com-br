package pdf

import (
	"fmt"
	"strings"
)

// Logical structure ("tagging").
//
// An untagged PDF leaves a screen reader to infer reading order from the
// geometry of the marks on the page. A tagged one states it: every run of text
// is wrapped in a marked-content sequence, those sequences are collected into a
// structure tree, and everything decorative is declared an artifact so assistive
// technology skips it rather than announcing it.
//
// Roles are the standard PDF structure types. Only the handful the resume needs
// are defined here.
const (
	RoleDocument = "Document"
	RoleH1       = "H1"
	RoleH2       = "H2"
	RoleH3       = "H3"
	RoleP        = "P"
	RoleFigure   = "Figure"

	// RoleArtifact is not a structure type. Content marked with it carries no
	// meaning — rules, tints, page furniture — and is excluded from the tree.
	RoleArtifact = "Artifact"
)

// mark is one marked-content sequence on a page, and the structure element it
// will become.
type mark struct {
	role string
	alt  string
	ref  Ref
}

// SetRole selects the structure role applied to subsequent text. It stays in
// force until changed, so a caller can set a heading role, draw, and reset.
func (p *Page) SetRole(role string) { p.role = role }

// Role returns the role currently in force.
func (p *Page) Role() string {
	if p.role == "" {
		return RoleP
	}
	return p.role
}

// beginContent opens a marked-content sequence for the current role and returns
// whether the caller must close it. Artifacts are marked but not recorded, so
// they never reach the structure tree.
func (p *Page) beginContent(role, alt string) {
	if role == RoleArtifact {
		p.op("/Artifact BMC")
		return
	}
	mcid := len(p.marks)
	p.marks = append(p.marks, mark{role: role, alt: alt})
	if alt == "" {
		p.op("/%s <</MCID %d>> BDC", role, mcid)
		return
	}
	p.op("/%s <</MCID %d>> BDC", role, mcid)
}

func (p *Page) endContent() { p.op("EMC") }

// artifact wraps purely decorative drawing so it is skipped by assistive
// technology and by text extraction.
func (p *Page) artifact(draw func()) {
	p.op("/Artifact BMC")
	draw()
	p.op("EMC")
}

// ---------------------------------------------------------------------------
// Structure tree
// ---------------------------------------------------------------------------

// buildStructTree writes the structure elements for every page and returns the
// StructTreeRoot reference. Pages are numbered by their /StructParents key,
// which the parent tree maps back to the elements they contain.
func (d *Doc) buildStructTree() Ref {
	root := d.reserve()
	docElem := d.reserve()

	var kids []string
	var nums []string

	for pageIndex, page := range d.pages {
		var perPage []string
		for i := range page.marks {
			m := &page.marks[i]
			var extra string
			if m.alt != "" {
				extra = " /Alt " + utf16Str(m.alt)
			}
			m.ref = d.add(fmt.Sprintf(
				"<< /Type /StructElem /S /%s /P %s /Pg %s /K %d%s >>",
				m.role, docElem, page.ref, i, extra))
			perPage = append(perPage, m.ref.String())
			kids = append(kids, m.ref.String())
		}
		nums = append(nums, fmt.Sprintf("%d [%s]", pageIndex, strings.Join(perPage, " ")))
	}

	d.fill(docElem, fmt.Sprintf("<< /Type /StructElem /S /%s /P %s /K [%s] >>",
		RoleDocument, root, strings.Join(kids, " ")))
	d.fill(root, fmt.Sprintf(
		"<< /Type /StructTreeRoot /K [%s] /ParentTree << /Nums [%s] >> /ParentTreeNextKey %d >>",
		docElem, strings.Join(nums, " "), len(d.pages)))

	return root
}
