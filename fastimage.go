package fastimage

import (
	"bytes"
	"compress/zlib"
	"regexp"
	"strconv"
	"unsafe"
)

type Type uint64

const (
	Unknown Type = iota
	BMP
	BPM
	GIF
	JPG
	MNG
	PBM
	PCD
	PCX
	PGM
	PNG
	PPM
	PSD
	RAS
	RGB
	SVG
	SWF
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
		case 'C':
			pcd(p, &info)
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
			tiffbe(p, &info)
		}
	case 'I':
		if p[1] == 'I' && p[2] == '\x2a' && p[3] == '\x00' {
			tiffle(p, &info)
		}
	case '8':
		if p[1] == 'B' && p[2] == 'P' && p[3] == 'S' {
			psd(p, &info)
		}
	case 'F':
		if p[1] == 'W' && p[2] == 'S' {
			swf(p, &info)
		}
	case 'C':
		if p[1] == 'W' && p[2] == 'S' {
			swfmx(p, &info)
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
		if p[1] == '.' && p[2] == '\x01' {
			pcx(p, &info)
		}
	case '<':
		if p[1] == 's' &&
			p[2] == 'v' &&
			p[3] == 'g' &&
			(p[4] == ' ' || p[4] == '\t') {
			svg(p, &info)
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
			info.Type = JPG
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
	case ' ':
		info.Width = (uint32(b[27])<<8 | uint32(b[26])) & 0x3fff
		info.Height = (uint32(b[29])<<8 | uint32(b[28])) & 0x3fff
	case 'L':
		info.Width = (uint32(b[20])<<16 | uint32(b[21])<<8 | uint32(b[22])) - 1
		info.Height = (uint32(b[24])<<16 | uint32(b[25])<<8 | uint32(b[26])) - 1
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
	stream := bs{b: b}
	stream.skip(18)
	info.Width = stream.uint32()
	info.Height = stream.uint32()
	if info.Width != 0 && info.Height != 0 {
		info.Type = BMP
	}
}

func ppm(b []byte, info *Info) {
	stream := bs{b: b}
	for {
		line := stream.readline()
		if len(line) == 0 {
			return
		}
		if line[0] != '#' {
			break
		}
	}
	if stream.readbyte() != 'P' {
		return
	}
	switch stream.readbyte() {
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
	stream.skipspace()
	info.Width = stream.readint()
	stream.skipspace()
	info.Height = stream.readint()
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
	info.Width = parseUint32(m[0][1])
	info.Height = parseUint32(m[1][1])
	if info.Width != 0 && info.Height != 0 {
		info.Type = XBM
	}
}

var xpmRegex = regexp.MustCompile(`\s*(\d+)\s+(\d+)(\s+\d+\s+\d+){1,2}\s*`)

func xpm(b []byte, info *Info) {
	stream := bs{b: b}
	for {
		line := stream.readline()
		if len(line) == 0 {
			break
		}
		m := xbmRegex.FindAllSubmatch(line, -1)
		if m != nil {
			info.Width = parseUint32(m[0][1])
			info.Height = parseUint32(m[1][1])
			break
		}
	}
	if info.Width != 0 && info.Height != 0 {
		info.Type = XPM
	}
}

func tiffbe(b []byte, info *Info) {
	stream := bs{b: b}
	offset := int(stream.uint32())
	stream.skip(offset)
	dirent := offset + (int(stream.uint16()) * 12)

	for info.Width == 0 && info.Height == 0 {
		ifd := stream.read(12)
		if len(ifd) != 0 || stream.tell() > dirent {
			break
		}
		tag := uint16(ifd[0])<<8 | uint16(ifd[1])
		typ := uint16(ifd[2])<<8 | uint16(ifd[3])
		switch typ {
		case 1:
			switch tag {
			case 0x0100:
				info.Width = uint32(ifd[8])
			case 0x0101:
				info.Height = uint32(ifd[8])
			}
		case 3:
			switch tag {
			case 0x0100:
				info.Width = uint32(ifd[8])<<8 | uint32(ifd[9])
			case 0x0101:
				info.Height = uint32(ifd[8])<<8 | uint32(ifd[9])
			}
		case 4:
			switch tag {
			case 0x0100:
				info.Width = uint32(ifd[8])<<24 |
					uint32(ifd[9])<<16 |
					uint32(ifd[10])<<8 |
					uint32(ifd[11])
			case 0x0101:
				info.Height = uint32(ifd[8])<<24 |
					uint32(ifd[9])<<16 |
					uint32(ifd[10])<<8 |
					uint32(ifd[11])
			}
		case 6:
			switch tag {
			case 0x0100:
				info.Width = uint32(ifd[8] & 0x7f)
			case 0x0101:
				info.Height = uint32(ifd[8] & 0x7f)
			}
		case 8:
			switch tag {
			case 0x0100:
				info.Width = uint32(ifd[8]&0x7f)<<8 | uint32(ifd[9])
			case 0x0101:
				info.Height = uint32(ifd[8]&0x7f)<<8 | uint32(ifd[9])
			}
		case 9:
			switch tag {
			case 0x0100:
				info.Width = uint32(ifd[8]&0x7f)<<24 |
					uint32(ifd[9])<<16 |
					uint32(ifd[10])<<8 |
					uint32(ifd[11])
			case 0x0101:
				info.Height = uint32(ifd[8]&0x7f)<<24 |
					uint32(ifd[9])<<16 |
					uint32(ifd[10])<<8 |
					uint32(ifd[11])
			}
		}
	}
}

func tiffle(b []byte, info *Info) {

}

func psd(b []byte, info *Info) {
	stream := bs{b: b}
	stream.skip(14)
	info.Height = stream.uint32()
	info.Width = stream.uint32()
	if info.Width != 0 && info.Height != 0 {
		info.Type = PSD
	}
}

func pcd(b []byte, info *Info) {
	if len(b) < 0xf00 {
		return
	}
	_ = b[0xf00]

	if !(b[0x800] == 'P' && b[0x801] == 'C' && b[0x802] == 'D') {
		return
	}

	switch b[0x0e02] & 1 {
	case 1:
		info.Type = PCD
		info.Width = 768
		info.Height = 512
	default:
		info.Type = PCD
		info.Width = 512
		info.Height = 768
	}
}

func swf(b []byte, info *Info) {
	if len(b) < 34 {
		return
	}
	_ = b[34]

	var bs []byte
	for i := 17; i < 25; i++ {
		bs = append(bs, strconv.FormatUint(uint64(b[i]), 2)...)
	}
	bits := b[8] & 0xf8 >> 3

	if x, _ := strconv.ParseUint(b2s(bs[5+bits:5+2*bits]), 2, 64); x > 0 {
		info.Width = uint32(x / 20)
	}
	if y, _ := strconv.ParseUint(b2s(bs[5+3*bits:5+4*bits]), 2, 64); y > 0 {
		info.Height = uint32(y / 20)
	}

	if info.Width != 0 && info.Height != 0 {
		info.Type = SWF
	}
}

func swfmx(b []byte, info *Info) {
	r, err := zlib.NewReader(bytes.NewReader(b))
	if err != nil {
		return
	}

	var p [1024]byte
	if _, err = r.Read(p[:]); err != nil {
		return
	}

	swf(p[:], info)
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
	stream := bs{b: b}
	stream.skip(6)
	info.Width = uint32(stream.uint16())
	info.Height = uint32(stream.uint16())
	if info.Width != 0 && info.Height != 0 {
		info.Type = RGB
	}
}

func ras(b []byte, info *Info) {
	stream := bs{b: b}
	stream.skip(4)
	info.Width = stream.uint32()
	info.Height = stream.uint32()
	if info.Width != 0 && info.Height != 0 {
		info.Type = RAS
	}
}

func pcx(b []byte, info *Info) {
	stream := bs{b: b}
	stream.skip(4)
	p := stream.read(8)
	if len(p) < 8 {
		return
	}
	info.Height = (uint32(p[7])<<8 | uint32(p[6])) - (uint32(p[3])<<8 | uint32(p[2]))
	info.Width = (uint32(p[5])<<8 | uint32(p[4])) - (uint32(p[1])<<8 | uint32(p[0]))
	if info.Width != 0 && info.Height != 0 {
		info.Type = PCX
	}
}

var (
	svgRegexWidth  = regexp.MustCompile(`width\s*=\s*(["\'])(\d*\.?\d+)(?:px)?`)
	svgRegexHeight = regexp.MustCompile(`height\s*=\s*(["\'])(\d*\.?\d+)(?:px)?`)
)

func svg(b []byte, info *Info) {
	m := svgRegexWidth.FindAllSubmatch(b, -1)
	n := svgRegexWidth.FindAllSubmatch(b, -1)
	if m != nil && n != nil {
		info.Width = parseUint32(m[2][1])
		info.Height = parseUint32(n[2][1])
	}
	if info.Width != 0 && info.Height != 0 {
		info.Type = SVG
	}
}

type bs struct {
	b []byte
	i int
}

func (s *bs) skip(n int) {
	s.i += n
}

func (s *bs) tell() int {
	return s.i
}

func (s *bs) skipspace() {
	j := s.i
	for s.b[j] == ' ' || s.b[j] == '\t' || s.b[j] == '\r' || s.b[j] == '\n' {
		j++
	}
	s.i = j
}

func (s *bs) peek(n int) (b []byte) {
	pos := s.i + n
	if end := len(s.b) - 1; pos > end {
		pos = end
	}
	b = s.b[s.i:pos]
	return
}

func (s *bs) readbyte() (c byte) {
	c = s.b[s.i]
	s.i++
	return
}

func (s *bs) readline() (line []byte) {
	j := s.i
	for s.b[j] != '\n' {
		j++
	}
	line = s.b[s.i:j]
	s.i = j + 1
	return
}

func (s *bs) readint() (i uint32) {
	for {
		switch s.b[s.i] {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			i = i*10 + uint32(s.b[s.i]-'0')
		default:
			return
		}
		s.i++
	}
	return
}

func (s *bs) read(n int) (b []byte) {
	pos := s.i + n
	if end := len(s.b) - 1; pos > end {
		pos = end
	}
	b = s.b[s.i:pos]
	s.i += n
	return
}

func (s *bs) uint8() (n uint8) {
	n = s.b[s.i]
	s.i++
	return
}

func (s *bs) uint16() (n uint16) {
	n = uint16(s.b[s.i+1]) | uint16(s.b[s.i])<<8
	s.i += 2
	return
}

func (s *bs) uint32() (n uint32) {
	n = uint32(s.b[s.i+3]) |
		uint32(s.b[s.i+2])<<8 |
		uint32(s.b[s.i+1])<<16 |
		uint32(s.b[s.i])<<24
	s.i += 4
	return
}

func parseUint32(b []byte) (n uint32) {
	for i := 0; i < len(b); i++ {
		j := uint32(b[i] - '0')
		if j < 0 || j > 9 {
			break
		}
		n = n*10 + j
	}
	return
}

func b2s(b []byte) string { return *(*string)(unsafe.Pointer(&b)) }
