package fastimage

import (
	"encoding/binary"
	"regexp"
)

type Type uint64

const (
	Unknown Type = iota
	BMP
	BPM
	GIF
	JPEG
	MNG
	PBM
	PCX
	PGM
	PNG
	PPM
	PSD
	RAS
	RGB
	TIFF
	WEBP
	XBM
	XPM
	XV
)

type Info struct {
	Type   Type
	Width  uint32
	Height uint32
}

// GetInfo return a image info of data.
func GetInfo(p []byte) (info Info) {
	const minOffset = 80 // 1 pixel gif
	if len(p) < minOffset {
		return
	}
	_ = p[minOffset-1]

	switch p[0] {
	case '\xff':
		if p[1] == '\xd8' {
			jpeg(p, &info)
		}
	case '\x89':
		if p[1] == 'P' &&
			p[2] == 'N' &&
			p[3] == 'G' &&
			p[4] == '\x0d' &&
			p[5] == '\x0a' &&
			p[6] == '\x1a' &&
			p[7] == '\x0a' {
			png(p, &info)
		}
	case 'R':
		if p[1] == 'I' &&
			p[2] == 'F' &&
			p[3] == 'F' &&
			p[8] == 'W' &&
			p[9] == 'E' &&
			p[10] == 'B' &&
			p[11] == 'P' {
			webp(p, &info)
		}
	case 'G':
		if p[1] == 'I' &&
			p[2] == 'F' &&
			p[3] == '8' &&
			(p[4] == '7' || p[4] == ',' || p[4] == '9') &&
			p[5] == 'a' {
			gif(p, &info)
		}
	case 'B':
		if p[1] == 'M' {
			bmp(p, &info)
		}
	case 'P':
		switch p[1] {
		case '1', '2', '3', '4', '5', '6', '7':
			ppm(p, &info)
		}
	case '#':
		if p[1] == 'd' &&
			p[2] == 'e' &&
			p[3] == 'f' &&
			p[4] == 'i' &&
			p[5] == 'n' &&
			p[6] == 'e' &&
			(p[7] == ' ' || p[7] == '\t') {
			xbm(p, &info)
		}
	case '/':
		if p[1] == '*' &&
			p[2] == ' ' &&
			p[3] == 'X' &&
			p[4] == 'P' &&
			p[5] == 'M' &&
			p[6] == ' ' &&
			p[7] == '*' &&
			p[8] == '/' {
			xpm(p, &info)
		}
	case 'M':
		if p[1] == 'M' && p[2] == '\x00' && p[3] == '\x2a' {
			tiff(p, &info, binary.BigEndian)
		}
	case 'I':
		if p[1] == 'I' && p[2] == '\x2a' && p[3] == '\x00' {
			tiff(p, &info, binary.LittleEndian)
		}
	case '8':
		if p[1] == 'B' && p[2] == 'P' && p[3] == 'S' {
			psd(p, &info)
		}
	case '\x8a':
		if p[1] == 'M' &&
			p[2] == 'N' &&
			p[3] == 'G' &&
			p[4] == '\x0d' &&
			p[5] == '\x0a' &&
			p[6] == '\x1a' &&
			p[7] == '\x0a' {
			mng(p, &info)
		}
	case '\x01':
		if p[1] == '\xda' &&
			p[2] == '[' &&
			p[3] == '\x01' &&
			p[4] == '\x00' &&
			p[5] == ']' {
			rgb(p, &info)
		}
	case '\x59':
		if p[1] == '\xa6' && p[2] == '\x6a' && p[3] == '\x95' {
			ras(p, &info)
		}
	case '\x0a':
		if p[2] == '\x01' {
			pcx(p, &info)
		}
	}

	return
}

func jpeg(b []byte, info *Info) {
	i := 2
	for {
		marker := b[i]
		code := b[i+1]
		length := uint16(b[i+2])<<8 | uint16(b[i+3])
		i += 4
		switch {
		case marker != 0xff:
			return
		case code >= 0xc0 && code <= 0xc3:
			info.Type = JPEG
			info.Width = uint32(b[i+3])<<8 | uint32(b[i+4])
			info.Height = uint32(b[i+1])<<8 | uint32(b[i+2])
			return
		default:
			i += int(length) - 2
		}
	}
}

func webp(b []byte, info *Info) {
	if len(b) < 30 {
		return
	}
	_ = b[29]

	if !(b[12] == 'V' && b[13] == 'P' && b[14] == '8') {
		return
	}

	switch b[15] {
	case ' ': // VP8
		info.Width = (uint32(b[27])&0x3f)<<8 | uint32(b[26])
		info.Height = (uint32(b[29])&0x3f)<<8 | uint32(b[28])
	case 'L': // VP8L
		info.Width = (uint32(b[22])<<8|uint32(b[21]))&16383 + 1
		info.Height = (uint32(b[23])<<2|uint32(b[22]>>6))&16383 + 1
	case 'X': // VP8X
		info.Width = (uint32(b[24]) | uint32(b[25])<<8 | uint32(b[26])<<16) + 1
		info.Height = (uint32(b[27]) | uint32(b[28])<<8 | uint32(b[29])<<16) + 1
	}
	if info.Width != 0 && info.Height != 0 {
		info.Type = WEBP
	}
}

func png(b []byte, info *Info) {
	if len(b) < 24 {
		return
	}
	_ = b[23]

	// IHDR
	if b[12] == 'I' && b[13] == 'H' && b[14] == 'D' && b[15] == 'R' {
		info.Width = uint32(b[16])<<24 |
			uint32(b[17])<<16 |
			uint32(b[18])<<8 |
			uint32(b[19])
		info.Height = uint32(b[20])<<24 |
			uint32(b[21])<<16 |
			uint32(b[22])<<8 |
			uint32(b[23])
	}

	if info.Width != 0 && info.Height != 0 {
		info.Type = PNG
	}
}

func gif(b []byte, info *Info) {
	if len(b) < 12 {
		return
	}
	_ = b[11]

	info.Width = uint32(b[7])<<8 | uint32(b[6])
	info.Height = uint32(b[9])<<8 | uint32(b[8])

	if info.Width != 0 && info.Height != 0 {
		info.Type = GIF
	}
}

func bmp(b []byte, info *Info) {
	if len(b) < 26 {
		return
	}
	_ = b[25]

	info.Width = uint32(b[21])<<24 |
		uint32(b[20])<<16 |
		uint32(b[19])<<8 |
		uint32(b[18])
	info.Height = uint32(b[25])<<24 |
		uint32(b[24])<<16 |
		uint32(b[23])<<8 |
		uint32(b[22])

	if info.Width != 0 && info.Height != 0 {
		info.Type = BMP
	}
}

func ppm(b []byte, info *Info) {
	switch b[1] {
	case '1':
		info.Type = PBM
	case '2', '5':
		info.Type = PGM
	case '3', '6':
		info.Type = PPM
	case '4':
		info.Type = BPM
	case '7':
		info.Type = XV
	}
	i := skipSpace(b, 2)
	info.Width, i = parseUint32(b, i)
	i = skipSpace(b, i)
	info.Height, _ = parseUint32(b, i)
	if info.Width == 0 || info.Height == 0 {
		info.Type = Unknown
	}
}

var xbmRegex = regexp.MustCompile(`^\#define\s*\S*\s*(\d+)\s*\n\#define\s*\S*\s*(\d+)`)

func xbm(b []byte, info *Info) {
	m := xbmRegex.FindAllSubmatch(b, -1)
	if m == nil {
		return
	}
	info.Width, _ = parseUint32(m[0][1], 0)
	info.Height, _ = parseUint32(m[0][2], 0)
	if info.Width != 0 && info.Height != 0 {
		info.Type = XBM
	}
}

var xpmRegex = regexp.MustCompile(`\s*(\d+)\s+(\d+)(\s+\d+\s+\d+){1,2}\s*`)

func xpm(b []byte, info *Info) {
	m := xpmRegex.FindAllSubmatch(b, -1)
	if m != nil {
		info.Width, _ = parseUint32(m[0][1], 0)
		info.Height, _ = parseUint32(m[0][2], 0)
	}
	if info.Width != 0 && info.Height != 0 {
		info.Type = XPM
	}
}

func tiff(b []byte, info *Info, order binary.ByteOrder) {
	i := int(order.Uint32(b[4:8]))
	n := int(order.Uint16(b[i+2 : i+4]))
	i += 2

	for ; i < n*12; i += 12 {
		tag := order.Uint16(b[i : i+2])
		datatype := order.Uint16(b[i+2 : i+4])

		var value uint32
		switch datatype {
		case 1, 6:
			value = uint32(b[i+9])
		case 3, 8:
			value = uint32(order.Uint16(b[i+8 : i+10]))
		case 4, 9:
			value = order.Uint32(b[i+8 : i+12])
		default:
			return
		}

		switch tag {
		case 256:
			info.Width = value
		case 257:
			info.Height = value
		}

		if info.Width > 0 && info.Height > 0 {
			info.Type = TIFF
			return
		}
	}
}

func psd(b []byte, info *Info) {
	if len(b) < 22 {
		return
	}
	_ = b[21]

	info.Width = uint32(b[18])<<24 |
		uint32(b[19])<<16 |
		uint32(b[20])<<8 |
		uint32(b[21])
	info.Height = uint32(b[14])<<24 |
		uint32(b[15])<<16 |
		uint32(b[16])<<8 |
		uint32(b[17])

	if info.Width != 0 && info.Height != 0 {
		info.Type = PSD
	}
}

func mng(b []byte, info *Info) {
	if len(b) < 24 {
		return
	}
	_ = b[23]

	if !(b[12] == 'M' && b[13] == 'H' && b[14] == 'D' && b[15] == 'R') {
		return
	}

	info.Width = uint32(b[16])<<24 |
		uint32(b[17])<<16 |
		uint32(b[18])<<8 |
		uint32(b[19])
	info.Height = uint32(b[20])<<24 |
		uint32(b[21])<<16 |
		uint32(b[22])<<8 |
		uint32(b[23])
	if info.Width != 0 && info.Height != 0 {
		info.Type = MNG
	}
}

func rgb(b []byte, info *Info) {
	if len(b) < 10 {
		return
	}
	_ = b[9]

	info.Width = uint32(b[6])<<8 |
		uint32(b[7])
	info.Height = uint32(b[8])<<8 |
		uint32(b[9])

	if info.Width != 0 && info.Height != 0 {
		info.Type = RGB
	}
}

func ras(b []byte, info *Info) {
	if len(b) < 12 {
		return
	}
	_ = b[11]

	info.Width = uint32(b[4])<<24 |
		uint32(b[5])<<16 |
		uint32(b[6])<<8 |
		uint32(b[7])
	info.Height = uint32(b[8])<<24 |
		uint32(b[9])<<16 |
		uint32(b[10])<<8 |
		uint32(b[11])

	if info.Width != 0 && info.Height != 0 {
		info.Type = RAS
	}
}

func pcx(b []byte, info *Info) {
	if len(b) < 12 {
		return
	}
	_ = b[11]

	info.Width = 1 +
		(uint32(b[9])<<8 | uint32(b[8])) -
		(uint32(b[5])<<8 | uint32(b[4]))
	info.Height = 1 +
		(uint32(b[11])<<8 | uint32(b[10])) -
		(uint32(b[7])<<8 | uint32(b[6]))

	if info.Width != 0 && info.Height != 0 {
		info.Type = PCX
	}
}

func skipSpace(b []byte, i int) (j int) {
	j = i
	for b[j] == ' ' || b[j] == '\t' || b[j] == '\r' || b[j] == '\n' {
		j++
	}
	return
}

func parseUint32(b []byte, i int) (n uint32, j int) {
	for j = i; j < len(b); j++ {
		x := uint32(b[j] - '0')
		if x < 0 || x > 9 {
			break
		}
		n = n*10 + x
	}
	return
}
