package pdf

import (
	"encoding/binary"
	"errors"
)

// Image is a JPEG passed straight through into the PDF.
//
// PDF's /DCTDecode filter *is* JPEG, so the compressed bytes are embedded
// verbatim: no decoding, no re-encoding, and no dependency on image/jpeg. Only
// the frame header is parsed, to learn the dimensions and component count.
type Image struct {
	Width, Height int
	Components    int

	data []byte
	ref  Ref
}

var errNotJPEG = errors.New("pdf: not a baseline or progressive JPEG")

// LoadJPEG parses the JPEG frame header and returns an embeddable image.
func LoadJPEG(data []byte) (*Image, error) {
	if len(data) < 4 || data[0] != 0xFF || data[1] != 0xD8 {
		return nil, errNotJPEG
	}
	for i := 2; i+3 < len(data); {
		if data[i] != 0xFF {
			i++
			continue
		}
		marker := data[i+1]
		// Standalone markers carry no payload.
		if marker == 0xFF || marker == 0x01 || (marker >= 0xD0 && marker <= 0xD9) {
			i += 2
			continue
		}
		if i+4 > len(data) {
			break
		}
		length := int(binary.BigEndian.Uint16(data[i+2 : i+4]))
		// SOF0..SOF15, excluding the DHT/JPG/DAC markers interleaved in the range.
		if marker >= 0xC0 && marker <= 0xCF &&
			marker != 0xC4 && marker != 0xC8 && marker != 0xCC {
			if i+9 > len(data) {
				break
			}
			return &Image{
				Height:     int(binary.BigEndian.Uint16(data[i+5 : i+7])),
				Width:      int(binary.BigEndian.Uint16(data[i+7 : i+9])),
				Components: int(data[i+9]),
				data:       data,
			}, nil
		}
		i += 2 + length
	}
	return nil, errNotJPEG
}

// colorSpace maps the JPEG component count onto a PDF colour space name.
func (im *Image) colorSpace() string {
	switch im.Components {
	case 1:
		return "DeviceGray"
	case 4:
		return "DeviceCMYK"
	default:
		return "DeviceRGB"
	}
}
