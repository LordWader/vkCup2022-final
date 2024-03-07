package taskio

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"os"
	"sync"
)

var bufPool = sync.Pool{New: func() interface{} {
	// 4*512 for one line of image + 7 bytes for png header
	b := make([]byte, 4*512+7, 4*512+7)
	//for ind, el := range []byte(pngHeader) {
	//	b[ind] = el
	//}
	return b
},
}

func CreateImage() {
	fo, _ := os.Create("output.png")

	imageList := []string{"tmp/0.png",
		"tmp/1.png",
		"tmp/2.png",
		"tmp/3.png",
		"tmp/4.png",
		"tmp/5.png",
		"tmp/6.png",
		"tmp/7.png",
		"tmp/8.png",
		"tmp/9.png"}

	wr := NewPNGwriter(fo)
	wr.WriteHeader()
	wr.WriteIHDR(512, 512*len(imageList))
	palette := make(color.Palette, 255)
	for i := range palette {
		palette[i] = color.NRGBA{0, 0, uint8(i), 255}
	}
	wr.WritePLTE(palette)

	for _, ii := range imageList {
		fi, _ := os.Open(ii)
		toParce, _, _ := image.Decode(fi)

		width := 512
		height := 512

		// TODO - may be y should be + i?
		upLeft := image.Point{0, 0}
		lowRight := image.Point{width, height}

		img := image.NewPaletted(image.Rectangle{upLeft, lowRight}, palette)
		for y := 0; y < height; y++ {
			dstPixOffset := img.PixOffset(0, y)
			for x := 0; x < width; x++ {
				_, _, b, _ := toParce.At(x, y).RGBA()
				img.Pix[dstPixOffset+x] = uint8(b)
			}
		}
		wr.WriteIDAT(img)
		_ = fi.Close()
	}
	wr.WriteIEND()
}

func CheckFile() {
	chunks := ChunkExtractor("output.png")
	for _, c := range chunks {
		fmt.Println(c.CType, c.Length, c.Crc32)
	}
}
