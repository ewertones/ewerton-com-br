// Package pdf is a minimal PDF 1.7 writer built entirely on the Go standard
// library.
//
// It exists so the site can generate documents at runtime without breaking the
// "Zero External Dependencies" rule the project is built around (see README).
// It implements exactly what a rich, machine-readable resume needs and nothing
// more: base-14 Helvetica text in WinAnsi encoding, vector fills and clipping,
// JPEG images passed through untouched, link annotations, a document outline,
// XMP metadata and embedded file attachments.
package pdf

import (
	"bytes"
	"compress/zlib"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

// Ref is a reference to an indirect PDF object. The zero value is invalid;
// object numbering starts at 1.
type Ref int

func (r Ref) String() string { return fmt.Sprintf("%d 0 R", int(r)) }

// Doc is a PDF document under construction.
type Doc struct {
	objs [][]byte // objs[i] holds the body of object number i+1

	pages []*Page

	// Document metadata, written to both the Info dictionary and the XMP packet.
	Title    string
	Author   string
	Subject  string
	Keywords []string
	Creator  string
	Lang     string
	Created  time.Time

	fonts     map[Font]Ref
	outline   []OutlineItem
	attaching []attachment
}

// New returns an empty document.
func New() *Doc {
	return &Doc{
		fonts:   map[Font]Ref{},
		Lang:    "en",
		Created: time.Now(),
		Creator: "ewerton.com.br",
	}
}

// reserve allocates an object number without writing a body yet.
func (d *Doc) reserve() Ref {
	d.objs = append(d.objs, nil)
	return Ref(len(d.objs))
}

// fill writes the body of a previously reserved object.
func (d *Doc) fill(r Ref, body string) { d.objs[int(r)-1] = []byte(body) }

// add writes a new object and returns its reference.
func (d *Doc) add(body string) Ref {
	r := d.reserve()
	d.fill(r, body)
	return r
}

// addStream writes a stream object. extra holds dictionary entries other than
// /Length and /Filter. When compress is true the payload is deflated.
func (d *Doc) addStream(extra string, payload []byte, compress bool) Ref {
	filter := ""
	if compress {
		var out bytes.Buffer
		zw, _ := zlib.NewWriterLevel(&out, zlib.BestCompression)
		zw.Write(payload)
		zw.Close()
		payload = out.Bytes()
		filter = " /Filter /FlateDecode"
	}
	var b strings.Builder
	b.WriteString("<< ")
	if extra != "" {
		b.WriteString(extra)
		b.WriteString(" ")
	}
	fmt.Fprintf(&b, "/Length %d%s >>\nstream\n", len(payload), filter)
	b.Write(payload)
	b.WriteString("\nendstream")
	return d.add(b.String())
}

// ---------------------------------------------------------------------------
// Strings, names and dates
// ---------------------------------------------------------------------------

// escapeString renders s as a PDF literal string. s must already hold the raw
// bytes to emit (WinAnsi or UTF-16BE).
func escapeString(s string) string {
	var b strings.Builder
	b.WriteByte('(')
	for i := 0; i < len(s); i++ {
		switch c := s[i]; c {
		case '(', ')', '\\':
			b.WriteByte('\\')
			b.WriteByte(c)
		case '\n':
			b.WriteString("\\n")
		case '\r':
			b.WriteString("\\r")
		case '\t':
			b.WriteString("\\t")
		default:
			b.WriteByte(c)
		}
	}
	b.WriteByte(')')
	return b.String()
}

// Str encodes a UTF-8 Go string as a WinAnsi PDF literal string.
func Str(s string) string { return escapeString(ToWinAnsi(s)) }

// utf16Str encodes s as a UTF-16BE PDF string, which is how the spec wants text
// that may fall outside WinAnsi (titles, bookmark labels).
func utf16Str(s string) string {
	var b bytes.Buffer
	b.WriteString("\xfe\xff")
	for _, r := range s {
		if r > 0xFFFF {
			r = '?'
		}
		b.WriteByte(byte(r >> 8))
		b.WriteByte(byte(r))
	}
	return escapeString(b.String())
}

// pdfDate formats a time as a PDF date string.
func pdfDate(t time.Time) string {
	_, off := t.Zone()
	sign := "+"
	if off < 0 {
		sign, off = "-", -off
	}
	return fmt.Sprintf("D:%s%s%02d'%02d'", t.Format("20060102150405"), sign, off/3600, (off%3600)/60)
}

// num formats a float without an exponent or noisy trailing zeros; content
// streams are far easier to debug this way.
func num(v float64) string {
	if v == float64(int64(v)) {
		return fmt.Sprintf("%d", int64(v))
	}
	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.3f", v), "0"), ".")
}

// fontRef lazily creates the font object for f and returns its reference.
func (d *Doc) fontRef(f Font) Ref {
	if r, ok := d.fonts[f]; ok {
		return r
	}
	r := d.add(fmt.Sprintf(
		"<< /Type /Font /Subtype /Type1 /BaseFont /%s /Encoding /WinAnsiEncoding >>",
		f.baseFont()))
	d.fonts[f] = r
	return r
}

// ---------------------------------------------------------------------------
// Outline
// ---------------------------------------------------------------------------

// OutlineItem is a bookmark pointing at a position on a page.
type OutlineItem struct {
	Title string
	Page  int     // zero-based page index
	Y     float64 // top-left Y coordinate to scroll to
}

// Bookmark adds a top-level outline entry pointing at y on page p.
func (d *Doc) Bookmark(title string, p *Page, y float64) {
	for i, pg := range d.pages {
		if pg == p {
			d.outline = append(d.outline, OutlineItem{Title: title, Page: i, Y: y})
			return
		}
	}
}

// ---------------------------------------------------------------------------
// Attachments
// ---------------------------------------------------------------------------

type attachment struct {
	name string
	mime string
	desc string
	data []byte
}

// Attach embeds a file inside the PDF. Readers expose it as an attachment, and
// tools that understand associated files (/AFRelationship /Source) can pull the
// structured data straight out of the document.
func (d *Doc) Attach(name, mime, desc string, data []byte) {
	d.attaching = append(d.attaching, attachment{name, mime, desc, data})
}

// ---------------------------------------------------------------------------
// Serialization
// ---------------------------------------------------------------------------

// Bytes assembles the final PDF file.
func (d *Doc) Bytes() []byte {
	pagesRef := d.reserve()

	kids := make([]string, 0, len(d.pages))
	for pageIndex, p := range d.pages {
		contents := d.addStream("", p.buf.Bytes(), true)

		var res strings.Builder
		res.WriteString("<< /ProcSet [/PDF /Text /ImageC] /Font << ")
		for _, f := range []Font{Regular, Bold, Italic} {
			if p.usedFonts[f] {
				fmt.Fprintf(&res, "/%s %s ", f.resource(), d.fontRef(f))
			}
		}
		res.WriteString(">>")
		if len(p.usedImages) > 0 {
			res.WriteString(" /XObject << ")
			for _, im := range p.usedImages {
				fmt.Fprintf(&res, "/%s %s ", im.name, im.img.ref)
			}
			res.WriteString(">>")
		}
		for key, refs := range map[string][]namedRef{
			"Shading": p.usedShadings, "ExtGState": p.usedGStates,
		} {
			if len(refs) == 0 {
				continue
			}
			fmt.Fprintf(&res, " /%s << ", key)
			for _, nr := range refs {
				fmt.Fprintf(&res, "/%s %s ", nr.name, nr.ref)
			}
			res.WriteString(">>")
		}
		res.WriteString(" >>")

		var annots string
		if len(p.annots) > 0 {
			refs := make([]string, len(p.annots))
			for i, a := range p.annots {
				refs[i] = d.add(a).String()
			}
			annots = " /Annots [" + strings.Join(refs, " ") + "]"
		}

		// /StructParents keys this page into the structure tree's parent tree.
		d.fill(p.ref, fmt.Sprintf(
			"<< /Type /Page /Parent %s /MediaBox [0 0 %s %s] /Resources %s /Contents %s%s /StructParents %d >>",
			pagesRef, num(p.W), num(p.H), res.String(), contents, annots, pageIndex))
		kids = append(kids, p.ref.String())
	}
	d.fill(pagesRef, fmt.Sprintf("<< /Type /Pages /Count %d /Kids [%s] >>",
		len(d.pages), strings.Join(kids, " ")))

	catalogRef := d.reserve()

	info := d.add(fmt.Sprintf(
		"<< /Title %s /Author %s /Subject %s /Keywords %s /Creator %s /Producer %s /CreationDate %s /ModDate %s >>",
		utf16Str(d.Title), utf16Str(d.Author), utf16Str(d.Subject),
		utf16Str(strings.Join(d.Keywords, ", ")), utf16Str(d.Creator),
		utf16Str("ewerton-go/pdf (pure Go stdlib)"),
		Str(pdfDate(d.Created)), Str(pdfDate(d.Created))))

	metadata := d.addStream("/Type /Metadata /Subtype /XML", d.xmp(), false)

	var extra strings.Builder
	fmt.Fprintf(&extra, " /StructTreeRoot %s", d.buildStructTree())
	if len(d.outline) > 0 {
		fmt.Fprintf(&extra, " /Outlines %s /PageMode /UseOutlines", d.buildOutline())
	}
	if len(d.attaching) > 0 {
		names, afs := d.buildAttachments()
		fmt.Fprintf(&extra, " /Names << /EmbeddedFiles << /Names [%s] >> >> /AF [%s]", names, afs)
	}

	d.fill(catalogRef, fmt.Sprintf(
		"<< /Type /Catalog /Version /1.7 /Pages %s /Lang %s /Metadata %s"+
			" /MarkInfo << /Marked true >>"+
			" /ViewerPreferences << /DisplayDocTitle true >>%s >>",
		pagesRef, Str(d.Lang), metadata, extra.String()))

	var out bytes.Buffer
	out.WriteString("%PDF-1.7\n%\xE2\xE3\xCF\xD3\n")
	offsets := make([]int, len(d.objs))
	for i, body := range d.objs {
		offsets[i] = out.Len()
		fmt.Fprintf(&out, "%d 0 obj\n", i+1)
		out.Write(body)
		out.WriteString("\nendobj\n")
	}

	xref := out.Len()
	fmt.Fprintf(&out, "xref\n0 %d\n0000000000 65535 f \n", len(d.objs)+1)
	for _, off := range offsets {
		fmt.Fprintf(&out, "%010d %05d n \n", off, 0)
	}

	sum := md5.Sum(out.Bytes())
	id := hex.EncodeToString(sum[:])
	fmt.Fprintf(&out, "trailer\n<< /Size %d /Root %s /Info %s /ID [<%s> <%s>] >>\nstartxref\n%d\n%%%%EOF\n",
		len(d.objs)+1, catalogRef, info, id, id, xref)

	return out.Bytes()
}

// buildOutline writes the outline tree and returns the root reference.
func (d *Doc) buildOutline() Ref {
	root := d.reserve()
	refs := make([]Ref, len(d.outline))
	for i := range d.outline {
		refs[i] = d.reserve()
	}
	for i, it := range d.outline {
		var sib string
		if i > 0 {
			sib += " /Prev " + refs[i-1].String()
		}
		if i < len(refs)-1 {
			sib += " /Next " + refs[i+1].String()
		}
		page := d.pages[it.Page]
		d.fill(refs[i], fmt.Sprintf(
			"<< /Title %s /Parent %s%s /Dest [%s /XYZ 0 %s 0] >>",
			utf16Str(it.Title), root, sib, page.ref, num(page.H-it.Y)))
	}
	d.fill(root, fmt.Sprintf("<< /Type /Outlines /Count %d /First %s /Last %s >>",
		len(refs), refs[0], refs[len(refs)-1]))
	return root
}

// buildAttachments writes embedded-file objects and returns the name-tree entry
// list plus the /AF array contents.
func (d *Doc) buildAttachments() (names, afs string) {
	var nameParts, afParts []string
	for _, at := range d.attaching {
		stream := d.addStream(fmt.Sprintf(
			"/Type /EmbeddedFile /Subtype /%s /Params << /Size %d /ModDate %s >>",
			nameEscape(at.mime), len(at.data), Str(pdfDate(d.Created))), at.data, true)
		spec := d.add(fmt.Sprintf(
			"<< /Type /Filespec /F %s /UF %s /Desc %s /AFRelationship /Source /EF << /F %s /UF %s >> >>",
			Str(at.name), utf16Str(at.name), utf16Str(at.desc), stream, stream))
		nameParts = append(nameParts, Str(at.name)+" "+spec.String())
		afParts = append(afParts, spec.String())
	}
	return strings.Join(nameParts, " "), strings.Join(afParts, " ")
}

// nameEscape encodes a PDF name body, hex-escaping characters that cannot
// appear literally (notably the "/" in a MIME type).
func nameEscape(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c <= ' ' || c >= '~' || strings.IndexByte("/%()<>[]{}#", c) >= 0 {
			fmt.Fprintf(&b, "#%02X", c)
			continue
		}
		b.WriteByte(c)
	}
	return b.String()
}
