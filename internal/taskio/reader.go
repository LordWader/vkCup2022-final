package taskio

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"hash/crc32"
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

func ReadFile() {
	files := []string{"tmp/0.png", "tmp/1.png", "tmp/2.png"}
	chunks := make([]*Chunk, 0)
	for _, f := range files {
		chunks = append(chunks, ChunkExtractor(f)...)
	}
	for _, c := range chunks {
		fmt.Println(c.CType, c.Length, c.Crc32)
	}
	// Construct png and save it
	out, _ := os.Create("output.png")
	defer out.Close()
	buf := bufio.NewWriter(out)
	//buf = bufio.NewWriterSize(buf, 60800)
	// Write header
	buf.WriteString(PNGHeader)
	// Write Idhf
	idhf := chunks[0]
	// write length
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(idhf.Length))
	buf.Write(b)
	// write chunk type
	buf.WriteString(idhf.CType)
	// write content
	// here we can resize image by half
	// 0 0 4 0 - 1024 width
	// 0 0 2 0 - 512 hight
	idhf.Data = []byte{0, 0, 2, 0, 0, 0, 6, 0, 8, 2, 0, 0, 0}
	buf.Write(idhf.Data)
	// write crc32 - checksum
	crc := crc32.NewIEEE()
	crc.Write([]byte("IHDR"))
	crc.Write(idhf.Data)
	binary.BigEndian.PutUint32(b, crc.Sum32())
	buf.Write(b)

	// write png content
	for _, c := range []*Chunk{chunks[1], chunks[4], chunks[7]} {
		binary.BigEndian.PutUint32(b, uint32(c.Length))
		buf.Write(b)
		buf.WriteString(c.CType)
		buf.Write(c.Data)
		//crc = crc32.NewIEEE()
		//crc.Write([]byte("IDAT"))
		//crc.Write(c.Data)
		//binary.BigEndian.PutUint32(b, crc.Sum32())
		buf.Write(c.Crc32)
		buf.Flush()
	}

	// write IEND
	end := chunks[2]
	binary.BigEndian.PutUint32(b, uint32(end.Length))
	buf.Write(b)
	buf.WriteString(end.CType)
	buf.Write(end.Data)
	buf.Write(end.Crc32)
	buf.Flush()

	//img, _, _ := png.Decode()
	fmt.Println("Success!")
}

func CheckFile() {
	chunks := ChunkExtractor("output.png")
	for _, c := range chunks {
		fmt.Println(c.CType, c.Length, c.Crc32)
	}
}
