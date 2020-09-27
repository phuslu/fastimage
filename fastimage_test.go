package fastimage

import (
	"io/ioutil"
	"testing"
)

func TestTypeString(t *testing.T) {
	cases := []struct {
		Type   Type
		String string
	}{
		{BMP, "bmp"},
		{BPM, "bpm"},
		{GIF, "gif"},
		{JPEG, "jpeg"},
		{MNG, "mng"},
		{PBM, "pbm"},
		{PCX, "pcx"},
		{PGM, "pgm"},
		{PNG, "png"},
		{PPM, "ppm"},
		{PSD, "psd"},
		{RAS, "ras"},
		{RGB, "rgb"},
		{TIFF, "tiff"},
		{WEBP, "webp"},
		{XBM, "xbm"},
		{XPM, "xpm"},
		{XV, "xv"},
	}

	for _, c := range cases {
		if got := c.Type.String(); got != c.String {
			t.Errorf("get type string error, type=%+v want=%+v, got=%+v", c.Type, c.Type.String(), got)
		}
	}
}

func TestGetInfo(t *testing.T) {
	cases := []struct {
		File string
		Info Info
	}{
		{"testdata/letter_T.jpg", Info{JPEG, 52, 54}},
		{"testdata/4.sm.webp", Info{WEBP, 320, 241}},
		{"testdata/2_webp_a.webp", Info{WEBP, 386, 395}},
		{"testdata/2_webp_ll.webp", Info{WEBP, 386, 395}},
		{"testdata/4_webp_ll.webp", Info{WEBP, 421, 163}},
		{"testdata/pass-1_s.png", Info{PNG, 90, 60}},
		{"testdata/pak38.gif", Info{GIF, 333, 194}},
		{"testdata/xterm.bmp", Info{BMP, 64, 38}},
		{"testdata/letter_N.ppm", Info{PPM, 66, 57}},
		{"testdata/spacer50.xbm", Info{XBM, 50, 10}},
		{"testdata/xterm.xpm", Info{XPM, 64, 38}},
		{"testdata/bexjdic.tif", Info{TIFF, 35, 32}},
		{"testdata/lexjdic.tif", Info{TIFF, 35, 32}},
		{"testdata/letter_T.psd", Info{PSD, 52, 54}},
		{"testdata/letter_T.psd", Info{PSD, 52, 54}},
		{"testdata/468x60.psd", Info{PSD, 468, 60}},
		{"testdata/letter_T.mng", Info{MNG, 52, 54}},
		{"testdata/letter_T.ras", Info{RAS, 52, 54}},
		{"testdata/letter_T.pcx", Info{PCX, 52, 54}},
	}

	for _, c := range cases {
		data, err := ioutil.ReadFile(c.File)
		if err != nil {
			t.Errorf("read file(%+v) error: %+v", c.File, err)
		}

		if got := GetInfo(data); got != c.Info {
			t.Errorf("get info error, file=%+v want=%+v, got=%+v", c.File, c.Info, got)
		}
	}
}
