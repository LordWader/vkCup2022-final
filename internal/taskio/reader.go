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
	chunks := ChunkExtractor("tmp/0.png")
	for _, c := range chunks {
		fmt.Println(c.CType, c.Length, c.Crc32)
	}
	// Construct png and save it
	out, _ := os.Create("output.png")
	buf := bufio.NewWriter(out)
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
	idhf.Data = []byte{0, 0, 2, 0, 0, 0, 2, 0, 8, 2, 0, 0, 0}
	buf.Write(idhf.Data)
	// write crc32 - checksum
	crc := crc32.NewIEEE()
	crc.Write([]byte("IHDR"))
	crc.Write(idhf.Data)
	binary.BigEndian.PutUint32(b, crc.Sum32())
	buf.Write(b)

	// write png content
	data := chunks[1]
	newLength := data.Length
	binary.BigEndian.PutUint32(b, uint32(newLength))
	buf.Write(b)
	buf.WriteString(data.CType)
	buf.Write(data.Data[:])
	crc = crc32.NewIEEE()
	crc.Write([]byte("IDAT"))
	crc.Write(data.Data[:])
	binary.BigEndian.PutUint32(b, crc.Sum32())
	buf.Write(b)

	// write IEND
	end := chunks[2]
	binary.BigEndian.PutUint32(b, uint32(end.Length))
	buf.Write(b)
	buf.WriteString(end.CType)
	buf.Write(end.Data)
	buf.Write(end.Crc32)
	buf.Flush()
	fmt.Println("Success!")
}

func CheckFile() {
	chunks := ChunkExtractor("output.png")
	for _, c := range chunks {
		fmt.Println(c.CType, c.Length, c.Crc32)
	}
}
