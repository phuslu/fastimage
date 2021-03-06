# fastimage - fast image info for go

[![godoc][godoc-img]][godoc] [![release][release-img]][release] [![goreport][goreport-img]][goreport] [![coverage][coverage-img]][coverage]

## Features

* Zero Dependencies - no std/3rd imports
* High Performance - hand-wringing DFA instead of regex/wildcard expression
* Widely Format - bmp/bpm/gif/jpeg/mng/pbm/pcx/pgm/png/ppm/psd/ras/tiff/webp/xbm/xpm

### Getting Started

try on https://play.golang.org/p/8yHaCknD1Rm
```go
package main

import (
	"fmt"
	"github.com/phuslu/fastimage"
)

var data = []byte("RIFF,-\x00\x00WEBPVP8X\n\x00\x00\x00" +
    "\x10\x00\x00\x00\x8f\x01\x00,\x01\x00VP8X\n\x00\x00\x00\x10\xb2" +
    "\x01\x00\x00WEB\x01\x00VP8X\n\x00\x00\x00\x10\xb2\x01\x00" +
    "\x00WEB\x01\x00VP8X\n\x00\x00\x00\x10\xb2\x01\x00\x00W" +
    "EB\x01\x00VP8X\n\x00\x00\x00\x10\xb2\x01\x00\x00WEB"")

func main() {
	fmt.Printf("%+v\n", fastimage.GetInfo(data))
}

// Output: {Type:webp Width:400 Height:301}
```

### Command Tool
```bash
$ go get github.com/phuslu/fastimage/cmd/fastimage
$ fastimage banner.png
png image/png 320 50
```

[godoc-img]: http://img.shields.io/badge/godoc-reference-blue.svg
[godoc]: https://godoc.org/github.com/phuslu/fastimage
[release-img]: https://img.shields.io/github/v/tag/phuslu/fastimage?label=release
[release]: https://github.com/phuslu/fastimage/releases
[goreport-img]: https://goreportcard.com/badge/github.com/phuslu/fastimage
[goreport]: https://goreportcard.com/report/github.com/phuslu/fastimage
[coverage-img]: http://gocover.io/_badge/github.com/phuslu/fastimage
[coverage]: https://gocover.io/github.com/phuslu/fastimage
