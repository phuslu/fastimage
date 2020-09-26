package fastimage

import (
	"io/ioutil"
	"testing"
)

func TestGetInfo(t *testing.T) {
	cases := []struct {
		File string
		Info Info
	}{
		{"testdata/letter_T.jpg", Info{JPG, 52, 54}},
		{"testdata/pass-1_s.png", Info{PNG, 90, 60}},
		{"testdata/4.sm.webp", Info{WEBP, 320, 241}},
		{"testdata/pak38.gif", Info{GIF, 333, 194}},
	}

	for _, c := range cases {
		data, err := ioutil.ReadFile(c.File)
		if err != nil {
			t.Errorf("read file(%+v) error: %+v", c.File, err)
		}

		if got := GetInfo(data); got != c.Info {
			t.Errorf("get image error, file=%+v want=%+v, got=%+v", c.File, c.Info, got)
		}
	}
}
