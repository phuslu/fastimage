package fastimage

type Type uint64

const (
	Unknown Type = iota
	JPEG
	BMP
	PNG
	GIF
	PPM
	XBM
	XPM
	TIFF
	PSD
	PCD
	SWF
	MNG
	RGB
	RAS
	PCX
	SVG
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
			info = jpeg(p)
		}
	case 'B':
		if p[1] == 'M' {
			info = bmp(p)
		}
	case '\x89':
		if p[1] == 'P' &&
			p[2] == 'N' &&
			p[3] == 'G' &&
			p[4] == '\x0d' &&
			p[5] == '\x0a' &&
			p[6] == '\x1a' &&
			p[7] == '\x0a' {
			info = png(p)
		}
	case 'G':
		if p[1] == 'I' &&
			p[2] == 'F' &&
			p[3] == '7' &&
			p[3] == '8' &&
			p[3] == '9' &&
			p[4] == 'a' {
			info = gif(p)
		}
	case 'P':
		switch p[1] {
		case '1', '2', '3', '4', '5', '6', '7':
			info = ppm(p)
		case 'C':
			info = pcd(p)
		}
	case '#':
		if p[1] == 'd' &&
			p[2] == 'e' &&
			p[3] == 'f' &&
			p[4] == 'i' &&
			p[5] == 'n' &&
			p[6] == 'e' &&
			(p[7] == ' ' || p[7] == '\t') {
			info = xbm(p)
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
			info = xpm(p)
		}
	case 'M':
		if p[1] == 'M' && p[2] == '\x00' && p[3] == '\x2a' {
			info = tiff(p)
		}
	case 'I':
		if p[1] == 'I' && p[2] == '\x2a' && p[3] == '\x00' {
			info = tiff(p)
		}
	case '8':
		if p[1] == 'B' && p[2] == 'P' && p[3] == 'S' {
			info = psd(p)
		}
	case 'F':
		if p[1] == 'W' && p[2] == 'S' {
			info = swf(p)
		}
	case 'C':
		if p[1] == 'W' && p[2] == 'S' {
			info = swfmx(p)
		}
	case '\x8a':
		if p[1] == 'M' &&
			p[2] == 'N' &&
			p[3] == 'G' &&
			p[4] == '\x0d' &&
			p[5] == '\x0a' &&
			p[6] == '\x1a' &&
			p[7] == '\x0a' {
			info = mng(p)
		}
	case '\x01':
		if p[1] == '\xda' &&
			p[2] == '[' &&
			p[3] == '\x01' &&
			p[4] == '\x00' &&
			p[5] == ']' {
			info = rgb(p)
		}
	case '\x59':
		if p[1] == '\xa6' && p[2] == '\x6a' && p[3] == '\x95' {
			info = ras(p)
		}
	case '\x0a':
		if p[1] == '.' && p[2] == '\x01' {
			info = pcx(p)
		}
	case '<':
		if p[1] == 's' &&
			p[2] == 'v' &&
			p[3] == 'g' &&
			(p[4] == ' ' || p[4] == '\t') {
			info = svg(p)
		}
	}

	return
}

func jpeg(p []byte) (info Info) {
	return
}

func bmp(p []byte) (info Info) {
	return
}

func png(p []byte) (info Info) {
	return
}

func gif(p []byte) (info Info) {
	return
}

func ppm(p []byte) (info Info) {
	return
}

func xbm(p []byte) (info Info) {
	return
}

func xpm(p []byte) (info Info) {
	return
}

func tiff(p []byte) (info Info) {
	return
}

func psd(p []byte) (info Info) {
	return
}

func pcd(p []byte) (info Info) {
	return
}

func swf(p []byte) (info Info) {
	return
}

func swfmx(p []byte) (info Info) {
	return
}

func mng(p []byte) (info Info) {
	return
}

func rgb(p []byte) (info Info) {
	return
}

func ras(p []byte) (info Info) {
	return
}

func pcx(p []byte) (info Info) {
	return
}

func svg(p []byte) (info Info) {
	return
}
