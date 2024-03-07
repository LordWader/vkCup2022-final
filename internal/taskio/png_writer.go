package taskio

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	log "github.com/dsoprea/go-logging"
	"hash/crc32"
	"image"
	"image/color"
	"io"
)

type PNGWriter struct {
	w   io.Writer
	zw  *zlib.Writer
	buf *bytes.Buffer
}

func NewPNGwriter(wr io.Writer) *PNGWriter {
	buf := bytes.NewBuffer(make([]byte, 0, 4*1024*1024))
	zw, err := zlib.NewWriterLevel(buf, zlib.BestSpeed)
	if err != nil {
		panic(err)
	}
	return &PNGWriter{
		w:   wr,
		buf: buf,
		zw:  zw,
	}
}

func (p *PNGWriter) WriteHeader() {
	pngHeader := "\x89PNG\r\n\x1a\n"
	_, err := io.WriteString(p.w, pngHeader)
	if err != nil {
		log.Errorf("%v", err)
	}
}

func (p *PNGWriter) writeChunk(b []byte, name string) {
	n := uint32(len(b))
	header := make([]byte, 8)
	footer := make([]byte, 4)
	binary.BigEndian.PutUint32(header[:4], n)
	header[4] = name[0]
	header[5] = name[1]
	header[6] = name[2]
	header[7] = name[3]
	crc := crc32.NewIEEE()
	crc.Write(header[4:8])
	crc.Write(b)
	binary.BigEndian.PutUint32(footer, crc.Sum32())

	_, err := p.w.Write(header)
	if err != nil {
		return
	}
	_, err = p.w.Write(b)
	if err != nil {
		return
	}
	_, err = p.w.Write(footer)
}

func (p *PNGWriter) WriteIHDR(width, height int) {
	buf := make([]byte, 13)
	binary.BigEndian.PutUint32(buf[0:4], uint32(width))
	binary.BigEndian.PutUint32(buf[4:8], uint32(height))
	buf[8] = 8
	buf[9] = 3 // RGBA - TODO - change to only blue channel
	buf[10] = 0
	buf[11] = 0
	buf[12] = 0
	p.writeChunk(buf, "IHDR")
}

func (p *PNGWriter) WritePLTE(c color.Palette) error {
	tmp := make([]byte, 4*256)
	for i, cc := range c {
		c1 := color.NRGBAModel.Convert(cc).(color.NRGBA)
		tmp[3*i+0] = c1.R
		tmp[3*i+1] = c1.G
		tmp[3*i+2] = c1.B
	}
	p.writeChunk(tmp[:3*len(c)], "PLTE")
	return nil
}

func (p *PNGWriter) WriteImage(m *image.Paletted) error {

	b := m.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		// for correct merging
		if _, err := p.zw.Write([]byte{0}); err != nil {
			return err
		}

		offset := y * m.Stride

		if _, err := p.zw.Write(m.Pix[offset : offset+b.Max.X]); err != nil {
			return err
		}
	}
	return nil
}

func (p *PNGWriter) WriteIDAT(i *image.Paletted) {
	err := p.WriteImage(i)
	if err != nil {
		return
	}
	p.writeChunk(p.buf.Bytes(), "IDAT")
	p.buf.Reset()
}

func (p *PNGWriter) FinishIdat() {
	if err := p.zw.Close(); err != nil {
		fmt.Printf("error while closing zlib writer: %v\n", err)
	}
}

func (p *PNGWriter) WriteIEND() {
	p.writeChunk(nil, "IEND")
}
