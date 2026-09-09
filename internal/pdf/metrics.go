package pdf

// Adobe Core-14 glyph widths for the Helvetica family, in 1/1000 em.
//
// Helvetica-Oblique shares Helvetica's widths, so only two tables are needed.
// Without these, text could not be wrapped, centred or right-aligned.

// asciiRegular holds widths for WinAnsi codes 32..126.
var asciiRegular = [95]int{
	278, 278, 355, 556, 556, 889, 667, 191, 333, 333, // 32..41
	389, 584, 278, 333, 278, 278, 556, 556, 556, 556, // 42..51
	556, 556, 556, 556, 556, 556, 278, 278, 584, 584, // 52..61
	584, 556, 1015, 667, 667, 722, 722, 667, 611, 778, // 62..71
	722, 278, 500, 667, 556, 833, 722, 778, 667, 778, // 72..81
	722, 667, 611, 722, 667, 944, 667, 667, 611, 278, // 82..91
	278, 278, 469, 556, 333, 556, 556, 500, 556, 556, // 92..101
	278, 556, 556, 222, 222, 500, 222, 833, 556, 556, // 102..111
	556, 556, 333, 500, 278, 556, 500, 722, 500, 500, // 112..121
	500, 334, 260, 334, 584, // 122..126
}

// asciiBold holds Helvetica-Bold widths for WinAnsi codes 32..126.
var asciiBold = [95]int{
	278, 333, 474, 556, 556, 889, 722, 238, 333, 333, // 32..41
	389, 584, 278, 333, 278, 278, 556, 556, 556, 556, // 42..51
	556, 556, 556, 556, 556, 556, 333, 333, 584, 584, // 52..61
	584, 611, 975, 722, 722, 722, 722, 667, 611, 778, // 62..71
	722, 278, 556, 722, 611, 833, 722, 778, 667, 778, // 72..81
	722, 667, 611, 722, 667, 944, 667, 667, 611, 333, // 82..91
	278, 333, 584, 556, 333, 556, 611, 556, 611, 556, // 92..101
	333, 611, 611, 278, 278, 556, 278, 889, 611, 611, // 102..111
	611, 611, 389, 556, 333, 611, 556, 778, 556, 556, // 112..121
	500, 389, 280, 389, 584, // 122..126
}

// highRegular holds the Helvetica widths for WinAnsi codes 128..255 that are
// not simple copies of an unaccented base glyph.
var highRegular = map[byte]int{
	0x80: 556, 0x82: 222, 0x83: 556, 0x84: 333, 0x85: 1000,
	0x86: 556, 0x87: 556, 0x88: 333, 0x89: 1000, 0x8A: 667,
	0x8B: 333, 0x8C: 1000, 0x8E: 611, 0x91: 222, 0x92: 222,
	0x93: 333, 0x94: 333, 0x95: 350, 0x96: 556, 0x97: 1000,
	0x98: 333, 0x99: 1000, 0x9A: 500, 0x9B: 333, 0x9C: 944,
	0x9E: 500, 0x9F: 667,
	0xA0: 278, 0xA1: 333, 0xA2: 556, 0xA3: 556, 0xA4: 556,
	0xA5: 556, 0xA6: 260, 0xA7: 556, 0xA8: 333, 0xA9: 737,
	0xAA: 370, 0xAB: 556, 0xAC: 584, 0xAD: 333, 0xAE: 737,
	0xAF: 333, 0xB0: 400, 0xB1: 584, 0xB2: 333, 0xB3: 333,
	0xB4: 333, 0xB5: 556, 0xB6: 537, 0xB7: 278, 0xB8: 333,
	0xB9: 333, 0xBA: 365, 0xBB: 556, 0xBC: 834, 0xBD: 834,
	0xBE: 834, 0xBF: 611,
	0xC6: 1000, 0xD0: 722, 0xD7: 584, 0xD8: 778, 0xDE: 667,
	0xDF: 611, 0xE6: 889, 0xF0: 556, 0xF7: 584, 0xF8: 611,
	0xFE: 556,
}

// highBold is the same table for Helvetica-Bold.
var highBold = map[byte]int{
	0x80: 556, 0x82: 278, 0x83: 556, 0x84: 500, 0x85: 1000,
	0x86: 556, 0x87: 556, 0x88: 333, 0x89: 1000, 0x8A: 667,
	0x8B: 333, 0x8C: 1000, 0x8E: 611, 0x91: 278, 0x92: 278,
	0x93: 500, 0x94: 500, 0x95: 350, 0x96: 556, 0x97: 1000,
	0x98: 333, 0x99: 1000, 0x9A: 556, 0x9B: 333, 0x9C: 889,
	0x9E: 500, 0x9F: 667,
	0xA0: 278, 0xA1: 333, 0xA2: 556, 0xA3: 556, 0xA4: 556,
	0xA5: 556, 0xA6: 280, 0xA7: 556, 0xA8: 333, 0xA9: 737,
	0xAA: 370, 0xAB: 556, 0xAC: 584, 0xAD: 333, 0xAE: 737,
	0xAF: 333, 0xB0: 400, 0xB1: 584, 0xB2: 333, 0xB3: 333,
	0xB4: 333, 0xB5: 611, 0xB6: 556, 0xB7: 278, 0xB8: 333,
	0xB9: 333, 0xBA: 365, 0xBB: 556, 0xBC: 834, 0xBD: 834,
	0xBE: 834, 0xBF: 611,
	0xC6: 1000, 0xD0: 722, 0xD7: 584, 0xD8: 778, 0xDE: 667,
	0xDF: 611, 0xE6: 889, 0xF0: 611, 0xF7: 584, 0xF8: 611,
	0xFE: 611,
}

// accentBase maps each accented WinAnsi code onto the ASCII letter it is built
// from. In Helvetica an accented glyph is exactly as wide as its base letter,
// with the single exception of the dotted i, whose accented forms are 278.
var accentBase = map[byte]byte{
	0xC0: 'A', 0xC1: 'A', 0xC2: 'A', 0xC3: 'A', 0xC4: 'A', 0xC5: 'A',
	0xC7: 'C',
	0xC8: 'E', 0xC9: 'E', 0xCA: 'E', 0xCB: 'E',
	0xCC: 'I', 0xCD: 'I', 0xCE: 'I', 0xCF: 'I',
	0xD1: 'N',
	0xD2: 'O', 0xD3: 'O', 0xD4: 'O', 0xD5: 'O', 0xD6: 'O',
	0xD9: 'U', 0xDA: 'U', 0xDB: 'U', 0xDC: 'U',
	0xDD: 'Y',
	0xE0: 'a', 0xE1: 'a', 0xE2: 'a', 0xE3: 'a', 0xE4: 'a', 0xE5: 'a',
	0xE7: 'c',
	0xE8: 'e', 0xE9: 'e', 0xEA: 'e', 0xEB: 'e',
	0xF1: 'n',
	0xF2: 'o', 0xF3: 'o', 0xF4: 'o', 0xF5: 'o', 0xF6: 'o',
	0xF9: 'u', 0xFA: 'u', 0xFB: 'u', 0xFC: 'u',
	0xFD: 'y', 0xFF: 'y',
}

var (
	regularWidths [256]int
	boldWidths    [256]int
)

func init() {
	build := func(dst *[256]int, ascii [95]int, high map[byte]int) {
		for i, w := range ascii {
			dst[32+i] = w
		}
		for code, w := range high {
			dst[code] = w
		}
		for code, base := range accentBase {
			dst[code] = ascii[base-32]
		}
		// Accented capital and lowercase I are wider than the bare glyph.
		for _, code := range []byte{0xCC, 0xCD, 0xCE, 0xCF, 0xEC, 0xED, 0xEE, 0xEF} {
			dst[code] = 278
		}
		// Anything still unset falls back to the space width so a missing
		// glyph never yields a zero-width measurement.
		for i := range dst {
			if dst[i] == 0 {
				dst[i] = ascii[0]
			}
		}
	}
	build(&regularWidths, asciiRegular, highRegular)
	build(&boldWidths, asciiBold, highBold)
}

// widths returns the 256-entry width table for f.
func widths(f Font) *[256]int {
	if f == Bold {
		return &boldWidths
	}
	return &regularWidths
}
