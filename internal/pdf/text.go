package pdf

import "strings"

// Font selects one of the base-14 Helvetica faces. Using the built-in faces
// keeps the file tiny and guarantees identical rendering everywhere, at the
// cost of not matching the site's Inter webfont exactly.
type Font int

const (
	Regular Font = iota
	Bold
	Italic
)

func (f Font) baseFont() string {
	switch f {
	case Bold:
		return "Helvetica-Bold"
	case Italic:
		return "Helvetica-Oblique"
	default:
		return "Helvetica"
	}
}

func (f Font) resource() string {
	switch f {
	case Bold:
		return "F2"
	case Italic:
		return "F3"
	default:
		return "F1"
	}
}

// winAnsiHigh maps the runes that CP1252 places in 0x80-0x9F, where it differs
// from Latin-1. These are the curly quotes, dashes and bullets that show up in
// the site copy.
var winAnsiHigh = map[rune]byte{
	'€': 0x80, '‚': 0x82, 'ƒ': 0x83, '„': 0x84,
	'…': 0x85, '†': 0x86, '‡': 0x87, 'ˆ': 0x88,
	'‰': 0x89, 'Š': 0x8A, '‹': 0x8B, 'Œ': 0x8C,
	'Ž': 0x8E, '‘': 0x91, '’': 0x92, '“': 0x93,
	'”': 0x94, '•': 0x95, '–': 0x96, '—': 0x97,
	'˜': 0x98, '™': 0x99, 'š': 0x9A, '›': 0x9B,
	'œ': 0x9C, 'ž': 0x9E, 'Ÿ': 0x9F,
}

// asciiFallback transliterates runes WinAnsi cannot represent, so a stray
// character degrades to something readable instead of a question mark.
var asciiFallback = map[rune]string{
	'→': "->", // rightwards arrow
	'←': "<-", // leftwards arrow
	'−': "-",  // minus sign
	'‐': "-",  // hyphen
	'‑': "-",  // non-breaking hyphen
	' ': " ",  // thin space
	' ': " ",  // narrow no-break space
	'​': "",   // zero-width space
}

// ToWinAnsi converts a UTF-8 string into WinAnsiEncoding bytes.
func ToWinAnsi(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch {
		case r < 0x80:
			b.WriteByte(byte(r))
		case r >= 0xA0 && r <= 0xFF:
			b.WriteByte(byte(r))
		default:
			if c, ok := winAnsiHigh[r]; ok {
				b.WriteByte(c)
				continue
			}
			if alt, ok := asciiFallback[r]; ok {
				b.WriteString(alt)
				continue
			}
			b.WriteByte('?')
		}
	}
	return b.String()
}

// TextWidth returns the rendered width of s in points.
func TextWidth(s string, f Font, size float64) float64 {
	w := widths(f)
	total := 0
	for _, c := range []byte(ToWinAnsi(s)) {
		total += w[c]
	}
	return float64(total) * size / 1000
}

// WrapText greedily word-wraps s to the given width. Existing newlines are
// honoured as hard breaks.
func WrapText(s string, f Font, size, width float64) []string {
	var out []string
	for _, para := range strings.Split(s, "\n") {
		words := strings.Fields(para)
		if len(words) == 0 {
			out = append(out, "")
			continue
		}
		line := words[0]
		for _, word := range words[1:] {
			if TextWidth(line+" "+word, f, size) <= width {
				line += " " + word
				continue
			}
			out = append(out, line)
			line = word
		}
		out = append(out, line)
	}
	return out
}

// Truncate shortens s with an ellipsis so it fits within width.
func Truncate(s string, f Font, size, width float64) string {
	if TextWidth(s, f, size) <= width {
		return s
	}
	runes := []rune(s)
	for len(runes) > 1 {
		runes = runes[:len(runes)-1]
		if TextWidth(string(runes)+"…", f, size) <= width {
			return strings.TrimRight(string(runes), " ") + "…"
		}
	}
	return ""
}
